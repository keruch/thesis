package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/spf13/cobra"
)

const (
	tssPartyID      = "poc-party-id"
	tssPartyMoniker = "poc-moniker"
)

func NewKeygenSimulateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keygen-simulate",
		Short: "Simulate TSS keygen",
		RunE: func(cmd *cobra.Command, args []string) error {
			keygenSimulate()
			return nil
		},
	}
	return cmd
}

func keygenSimulate() {
	partyIDs := generatePartyIDs(4)
	ctx := tss.NewPeerContext(partyIDs)

	// Select an elliptic curve
	curve := tss.S256()

	// Sey up channels for communication
	errCh := make(chan *tss.Error, len(partyIDs))
	outCh := make(chan tss.Message, len(partyIDs))
	endCh := make(chan *keygen.LocalPartySaveData, len(partyIDs))

	// TODO: make threshold configurable
	const threshold = 3

	parties := make([]*keygen.LocalParty, len(partyIDs))
	for i := range partyIDs {
		// Create the party
		params := tss.NewParameters(curve, ctx, partyIDs[i], len(partyIDs), threshold)
		party := keygen.NewLocalParty(params, outCh, endCh)
		parties[i] = party.(*keygen.LocalParty)
		fmt.Printf("Created party: Index %d, Moniker %s, PartyID %s\n", party.PartyID().Index, party.PartyID().Moniker, party.PartyID().Id)
	}

	for i := range parties {
		go func() {
			if err := parties[i].Start(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var keyShares []*keygen.LocalPartySaveData
outer:
	for {
		select {
		case err := <-errCh:
			log.Fatal(err)

		case msg := <-outCh:
			to := msg.GetTo()
			if to == nil { // broadcast
				to = partyIDs
			}

			// log format:
			// Type: binance.tsslib.ecdsa.keygen.KGRound2Message2, From: {1,P[2]}, To: all
			fmt.Printf("%s\n", msg.String())

			for _, p := range to {
				go handlePartyMessage(parties[p.Index], msg, errCh)
			}

		case keyShare := <-endCh:
			// in future every party will need to persist the signature
			// for now we just print it
			partyIdx, err := keyShare.OriginalIndex()
			if err != nil {
				panic(err)
			}
			partyID := parties[partyIdx].PartyID()
			fmt.Printf("Patry got signature: moniker %s\n", partyID.Moniker)
			keyShares = append(keyShares, keyShare)

			if len(keyShares) == len(parties) {
				// all parties have finished
				break outer
			}
		}
	}
	fmt.Println("All parties have finished, saving key shares...")
	for i := range keyShares {
		err := saveKeyShare(keyShares[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handlePartyMessage(to tss.Party, msg tss.Message, errCh chan<- *tss.Error) {
	// do not send a message from this party back to itself
	if to.PartyID() == msg.GetFrom() {
		return
	}

	bz, _, err := msg.WireBytes()
	if err != nil {
		errCh <- to.WrapError(err)
		return
	}
	_, wErr := to.UpdateFromBytes(bz, msg.GetFrom(), msg.IsBroadcast())
	if wErr != nil {
		errCh <- wErr
	}
}

func saveKeyShare(keyShare *keygen.LocalPartySaveData) error {
	partyIdx, err := keyShare.OriginalIndex()
	if err != nil {
		return fmt.Errorf("signature party original index: %v", err)
	}

	file, err := os.Create(fmt.Sprintf("data/key-share-%d.json", partyIdx))
	if err != nil {
		return fmt.Errorf("create signature file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(keyShare); err != nil {
		return fmt.Errorf("encode signature: %v", err)
	}

	return nil
}

func getKeyShare(partyIdx int) (*keygen.LocalPartySaveData, error) {
	file, err := os.Open(fmt.Sprintf("data/key-share-%d.json", partyIdx))
	if err != nil {
		return nil, fmt.Errorf("open signature file: %v", err)
	}
	defer file.Close()

	var keyShare keygen.LocalPartySaveData
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&keyShare); err != nil {
		return nil, fmt.Errorf("decode signature: %v", err)
	}

	return &keyShare, nil
}
