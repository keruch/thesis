package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/bnb-chain/tss-lib/v2/ecdsa/keygen"
	"github.com/bnb-chain/tss-lib/v2/tss"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize TSS configs",
		RunE: func(cmd *cobra.Command, args []string) error {
			return prepare()
		},
	}
	return cmd
}

// TODO: make home dir configurable
// TODO: make generation timeout configurable
// When using the keygen party it is recommended that you pre-compute the "safe
// primes" and Paillier secret beforehand because this can take some time. This
// code will generate those parameters using a concurrency limit equal to the
// number of available CPU cores. The generated parameters will be saved to a
// file named "pre-params.json" in the current directory.
func prepare() error {
	preParams, err := keygen.GeneratePreParams(time.Minute)
	if err != nil {
		return fmt.Errorf("generate pre-parameters: %v", err)
	}
	file, err := os.Create("pre-params.json")
	if err != nil {
		return fmt.Errorf("create pre-params file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(preParams); err != nil {
		return fmt.Errorf("encode pre-params: %v", err)
	}
	return nil
}

// TODO: make home dir configurable
func readPreParams() (keygen.LocalPreParams, error) {
	file, err := os.Open("pre-params.json")
	if err != nil {
		return keygen.LocalPreParams{}, fmt.Errorf("open pre-params file: %v", err)
	}
	defer file.Close()

	var preParams keygen.LocalPreParams
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&preParams); err != nil {
		return keygen.LocalPreParams{}, fmt.Errorf("decode pre-params: %v", err)
	}

	return preParams, nil
}

func generatePartyIDs(number int) []*tss.PartyID {
	ids := make([]*tss.PartyID, number)
	for i := 0; i < number; i++ {
		ids[i] = generatePartyID(i)
	}
	return tss.SortPartyIDs(ids)
}

func generatePartyID(index int) *tss.PartyID {
	partyID := fmt.Sprintf("%s-%d", tssPartyID, index)
	moniker := fmt.Sprintf("%s-%d", tssPartyMoniker, index)
	uniqueKey := big.NewInt(int64(index + 1))

	// Set up the parameters
	//
	// Note: The `id` and `moniker` fields are for convenience to allow you to easily
	// track participants. The `id` should be a unique string representing this party
	// in the network and `moniker` can be anything (even left blank). The `uniqueKey`
	// is a unique identifying key for this peer (such as its p2p public key) as a big.Int.
	return tss.NewPartyID(partyID, moniker, uniqueKey)
}
