package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"reflect"

	"github.com/bnb-chain/tss-lib/v2/common"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/ecdsa/signing"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/spf13/cobra"
)

func NewKeysignSimulateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keysign-simulate",
		Short: "Simulate TSS keysign",
		RunE: func(cmd *cobra.Command, args []string) error {
			keysignSimulate()
			return nil
		},
	}
	return cmd
}

func keysignSimulate() {
	partyIDs := generatePartyIDs(4)
	ctx := tss.NewPeerContext(partyIDs)

	// Select an elliptic curve
	curve := tss.S256()

	// Sey up channels for communication
	errCh := make(chan *tss.Error, len(partyIDs))
	outCh := make(chan tss.Message, len(partyIDs))
	endCh := make(chan *common.SignatureData, len(partyIDs))

	// TODO: make threshold configurable
	const threshold = 3

	keyShares := make(map[int]*keygen.LocalPartySaveData, len(partyIDs))
	for i := range partyIDs {
		ks, err := getKeyShare(i)
		if err != nil {
			log.Fatal(err)
		}
		keyShares[i] = ks
	}

	rawMsg := rand.Int63()
	msg := big.NewInt(rawMsg)
	fmt.Printf("Message: %s\n", msg)

	parties := make([]*signing.LocalParty, len(partyIDs))
	for i := range partyIDs {
		// Create the party
		params := tss.NewParameters(curve, ctx, partyIDs[i], len(partyIDs), threshold)
		party := signing.NewLocalParty(msg, params, *keyShares[i], outCh, endCh)
		parties[i] = party.(*signing.LocalParty)
		fmt.Printf("Created party: Index %d, Moniker %s, PartyID %s\n", party.PartyID().Index, party.PartyID().Moniker, party.PartyID().Id)
	}

	for i := range parties {
		go func() {
			if err := parties[i].Start(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var signatures []*common.SignatureData
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

		case signature := <-endCh:
			// in future every party will need to persist the signature
			// for now we just print it
			fmt.Printf("Patry got signature\n")
			signatures = append(signatures, signature)

			if len(signatures) == len(parties) {
				// all parties have finished
				break outer
			}
		}
	}

	fmt.Println("All parties have finished, validating signatures...")
	for i := range signatures {
		eq := reflect.DeepEqual(signatures[0], signatures[i])
		if !eq {
			log.Fatal("Signatures are not equal")
		}
	}

	sig := signatures[0]
	ks := keyShares[0]
	pk := ecdsa.PublicKey{
		Curve: curve,
		X:     ks.ECDSAPub.X(),
		Y:     ks.ECDSAPub.Y(),
	}
	r := new(big.Int).SetBytes(sig.GetR())
	s := new(big.Int).SetBytes(sig.GetS())

	ok := ecdsa.Verify(&pk, msg.Bytes(), r, s)
	if !ok {
		log.Fatal("ECDSA signature verification did not pass")
	}

	fmt.Println("Signature is valid")

	fmt.Println("X: ", ks.ECDSAPub.X().String())
	fmt.Println("Y: ", ks.ECDSAPub.Y().String())
	fmt.Println("M: ", msg.String())
	fmt.Println("R: ", r.String())
	fmt.Println("S: ", s.String())
}
