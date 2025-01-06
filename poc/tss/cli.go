package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
)

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tssd",
		Short: "TSS daemon for threshold signature operations",
		Long: `TSS daemon implements a transport layer for Threshold Signature Schemes (TSS).
It provides peer discovery, party management, and secure communication for TSS operations.`,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())
			return nil
		},
	}

	initRootCmd(rootCmd)
	return rootCmd
}

func initRootCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(
		NewStartCmd(),
		NewPartyCmd(),
		NewKeygenCmd(),
		NewSignCmd(),
	)
}

func NewStartCmd() *cobra.Command {
	var keyFile string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the TSS node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return startNode(cmd.Context(), keyFile)
		},
	}

	cmd.Flags().StringVarP(&keyFile, "key", "k", "node_key", "Path to the node's private key file")
	return cmd
}

func NewPartyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "party",
		Short: "Party management commands",
	}

	cmd.AddCommand(
		NewPartyCreateCmd(),
		NewPartyListCmd(),
		NewPartyInfoCmd(),
	)

	return cmd
}

func NewPartyCreateCmd() *cobra.Command {
	var (
		members   string
		threshold int
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new TSS party",
		RunE: func(cmd *cobra.Command, args []string) error {
			return createParty(members, threshold)
		},
	}

	cmd.Flags().StringVarP(&members, "members", "m", "", "Comma-separated list of peer IDs")
	cmd.Flags().IntVarP(&threshold, "threshold", "t", 2, "Threshold for the party")
	cmd.MarkFlagRequired("members")

	return cmd
}

func NewPartyListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all parties",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listParties()
		},
	}
}

func NewPartyInfoCmd() *cobra.Command {
	var partyID string

	cmd := &cobra.Command{
		Use:   "info",
		Short: "Get information about a specific party",
		RunE: func(cmd *cobra.Command, args []string) error {
			return getPartyInfo(partyID)
		},
	}

	cmd.Flags().StringVarP(&partyID, "party-id", "p", "", "Party ID to get information for")
	cmd.MarkFlagRequired("party-id")

	return cmd
}

func NewKeygenCmd() *cobra.Command {
	var partyID string

	cmd := &cobra.Command{
		Use:   "keygen",
		Short: "Initiate key generation for a party",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initiateKeyGeneration(partyID)
		},
	}

	cmd.Flags().StringVarP(&partyID, "party-id", "p", "", "Party ID for key generation")
	cmd.MarkFlagRequired("party-id")

	return cmd
}

func NewSignCmd() *cobra.Command {
	var (
		partyID string
		message string
	)

	cmd := &cobra.Command{
		Use:   "sign",
		Short: "Initiate signing process for a party",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initiateSigningProcess(partyID, message)
		},
	}

	cmd.Flags().StringVarP(&partyID, "party-id", "p", "", "Party ID for signing")
	cmd.Flags().StringVarP(&message, "message", "m", "", "Message to sign")
	cmd.MarkFlagRequired("party-id")
	cmd.MarkFlagRequired("message")

	return cmd
}

// Command execution functions

func startNode(ctx context.Context, keyFile string) error {
	privKey, err := loadOrCreatePrivateKey(keyFile)
	if err != nil {
		return fmt.Errorf("failed to load or create private key: %w", err)
	}

	node, err := NewNode(ctx, privKey)
	if err != nil {
		return fmt.Errorf("failed to create node: %w", err)
	}

	if err := node.Start(ctx); err != nil {
		return fmt.Errorf("failed to start node: %w", err)
	}

	fmt.Printf("Node started with ID: %s\n", node.host.ID().String())

	// Wait for interrupt signal
	<-ctx.Done()
	return node.Stop()
}

func createParty(membersStr string, threshold int) error {
	members := strings.Split(membersStr, ",")
	peerIDs := make([]peer.ID, len(members))
	for i, m := range members {
		peerID, err := peer.Decode(m)
		if err != nil {
			return fmt.Errorf("invalid peer ID %s: %w", m, err)
		}
		peerIDs[i] = peerID
	}

	party, err := globalNode.CreateParty(context.Background(), peerIDs, threshold, TSSOperationKeyGen)
	if err != nil {
		return fmt.Errorf("failed to create party: %w", err)
	}

	fmt.Printf("Party created with ID: %s\n", party.ID)
	return nil
}

func listParties() error {
	parties := globalNode.GetAllParties()
	fmt.Println("Parties:")
	for _, party := range parties {
		fmt.Printf("- ID: %s, Members: %d, Threshold: %d, Status: %s\n",
			party.ID, len(party.Members), party.Threshold, party.Status)
	}
	return nil
}

func getPartyInfo(partyID string) error {
	party, err := globalNode.GetParty(partyID)
	if err != nil {
		return fmt.Errorf("failed to get party info: %w", err)
	}

	fmt.Printf("Party ID: %s\n", party.ID)
	fmt.Printf("Members: %d\n", len(party.Members))
	fmt.Printf("Threshold: %d\n", party.Threshold)
	fmt.Printf("Status: %s\n", party.Status)
	fmt.Println("Member IDs:")
	for _, member := range party.Members {
		fmt.Printf("- %s\n", member.String())
	}
	return nil
}

func initiateKeyGeneration(partyID string) error {
	msg := &Message{
		Type:    MessageTypeKeyGeneration,
		PartyID: partyID,
		From:    globalNode.host.ID(),
		Payload: json.RawMessage(`{"action": "start_keygen"}`),
	}

	if err := globalNode.msgRouter.SendMessage(context.Background(), msg); err != nil {
		return fmt.Errorf("failed to initiate key generation: %w", err)
	}

	fmt.Printf("Key generation initiated for party: %s\n", partyID)
	return nil
}

func initiateSigningProcess(partyID string, message string) error {
	msg := &Message{
		Type:    MessageTypeSigning,
		PartyID: partyID,
		From:    globalNode.host.ID(),
		Payload: json.RawMessage(fmt.Sprintf(`{"message": "%s"}`, message)),
	}

	if err := globalNode.msgRouter.SendMessage(context.Background(), msg); err != nil {
		return fmt.Errorf("failed to initiate signing process: %w", err)
	}

	fmt.Printf("Signing process initiated for party: %s\n", partyID)
	return nil
}

// Helper functions

func loadOrCreatePrivateKey(keyFile string) (crypto.PrivKey, error) {
	privKey, err := loadPrivateKey(keyFile)
	if err == nil {
		return privKey, nil
	}

	privKey, _, err = crypto.GenerateKeyPair(crypto.Ed25519, -1)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	if err := savePrivateKey(privKey, keyFile); err != nil {
		return nil, fmt.Errorf("failed to save private key: %w", err)
	}

	return privKey, nil
}

func loadPrivateKey(keyFile string) (crypto.PrivKey, error) {
	keyBytes, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(keyBytes)
}

func savePrivateKey(privKey crypto.PrivKey, keyFile string) error {
	keyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return err
	}

	return os.WriteFile(keyFile, keyBytes, 0600)
}

var globalNode *Node // Global node instance for CLI commands

func init() {
	// Initialize the global node instance
	ctx := context.Background()
	privKey, err := loadOrCreatePrivateKey("node_key")
	if err != nil {
		log.Fatalf("Failed to initialize node: %v", err)
	}

	node, err := NewNode(ctx, privKey)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}

	globalNode = node
}
