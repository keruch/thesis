package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type Node struct {
	host       host.Host
	discovery  *NodeDiscovery
	partyMgr   *PartyManager
	msgRouter  *MessageRouter
	secLayer   *SecurityLayer
	tssHandler *TSSHandler
}

func NewNode(ctx context.Context, privKey crypto.PrivKey) (*Node, error) {
	h, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"),
		libp2p.Ping(false),
	)
	if err != nil {
		return nil, err
	}

	discovery, err := NewNodeDiscovery(ctx, h)
	if err != nil {
		return nil, fmt.Errorf("failed to create node discovery: %w", err)
	}

	msgRouter := NewMessageRouter(h)
	partyMgr := NewPartyManager(msgRouter)
	secLayer := NewSecurityLayer(privKey)
	tssHandler := NewTSSHandler(partyMgr)

	node := &Node{
		host:       h,
		discovery:  discovery,
		partyMgr:   partyMgr,
		msgRouter:  msgRouter,
		secLayer:   secLayer,
		tssHandler: tssHandler,
	}

	msgRouter.RegisterHandler(MessageTypePartyFormation, node.handlePartyFormation)
	msgRouter.RegisterHandler(MessageTypeKeyGeneration, node.handleKeyGeneration)
	msgRouter.RegisterHandler(MessageTypeSigning, node.handleSigning)

	return node, nil
}

func (n *Node) Start(ctx context.Context) error {
	if err := n.discovery.Start(); err != nil {
		return fmt.Errorf("failed to start node discovery: %w", err)
	}

	if err := n.msgRouter.Start(ctx); err != nil {
		return fmt.Errorf("failed to start message router: %w", err)
	}

	go n.handleDiscoveredPeers(ctx)

	return nil
}

func (n *Node) Stop() error {
	if err := n.discovery.Stop(); err != nil {
		return fmt.Errorf("failed to stop node discovery: %w", err)
	}

	if err := n.msgRouter.Stop(); err != nil {
		return fmt.Errorf("failed to stop message router: %w", err)
	}

	return nil
}

func (n *Node) handleDiscoveredPeers(ctx context.Context) {
	for {
		select {
		case peer := <-n.discovery.PeerChan():
			fmt.Printf("Discovered new peer: %s\n", peer.ID)
			// TODO: Implement logic to decide whether to form a party with the new peer
		case <-ctx.Done():
			return
		}
	}
}

func (n *Node) CreateParty(ctx context.Context, members []peer.ID, threshold int, operation TSSOperation) (*Party, error) {
	return n.partyMgr.CreateParty(ctx, n.host.ID(), members, threshold, operation)
}

func (n *Node) GetParty(partyID string) (*Party, error) {
	return n.partyMgr.GetParty(partyID)
}

func (n *Node) GetPeerParties(peerID peer.ID) []string {
	return n.partyMgr.GetPeerParties(peerID)
}

func (n *Node) UpdatePartyStatus(partyID string, status PartyStatus) error {
	return n.partyMgr.UpdatePartyStatus(partyID, status)
}

func (n *Node) SendSecureMessage(ctx context.Context, msg *Message, recipientPubKey crypto.PubKey) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	encryptedMsg, err := n.secLayer.EncryptMessage(msgBytes, recipientPubKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt message: %w", err)
	}

	signature, err := n.secLayer.SignMessage(encryptedMsg)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	secureMsg := &Message{
		Type:    msg.Type,
		PartyID: msg.PartyID,
		From:    n.secLayer.GetPeerID(),
		To:      msg.To,
		Payload: json.RawMessage(encryptedMsg),
	}

	secureMsg.Payload = append(secureMsg.Payload, []byte(":"+string(signature))...)

	return n.msgRouter.SendMessage(ctx, secureMsg)
}

func (n *Node) handleSecureMessage(msg *Message) (*Message, error) {
	parts := bytes.Split(msg.Payload, []byte(":"))
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid secure message format")
	}

	encryptedMsg, signature := parts[0], parts[1]

	senderPubKey := n.host.Peerstore().PubKey(msg.From)

	valid, err := n.secLayer.VerifySignature(encryptedMsg, signature, senderPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify signature: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid signature")
	}

	decryptedMsg, err := n.secLayer.DecryptMessage(encryptedMsg, senderPubKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt message: %w", err)
	}

	var originalMsg Message
	if err := json.Unmarshal(decryptedMsg, &originalMsg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted message: %w", err)
	}

	return &originalMsg, nil
}

func (n *Node) handlePartyFormation(msg *Message) error {
	secureMsg, err := n.handleSecureMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to handle secure message: %w", err)
	}
	// TODO: Implement party formation logic
	fmt.Printf("Received party formation message: %+v\n", secureMsg)
	return nil
}

func (n *Node) handleKeyGeneration(msg *Message) error {
	secureMsg, err := n.handleSecureMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to handle secure message: %w", err)
	}

	var keyGenRequest struct {
		PartyID string `json:"party_id"`
	}
	if err := json.Unmarshal(secureMsg.Payload, &keyGenRequest); err != nil {
		return fmt.Errorf("failed to unmarshal key generation request: %w", err)
	}

	keyShares, err := n.tssHandler.GenerateKeyShares(keyGenRequest.PartyID)
	if err != nil {
		return fmt.Errorf("failed to generate key shares: %w", err)
	}

	party, err := n.partyMgr.GetParty(keyGenRequest.PartyID)
	if err != nil {
		return fmt.Errorf("failed to get party: %w", err)
	}

	for _, member := range party.Members {
		if member == n.host.ID() {
			continue
		}

		keyShareMsg := &Message{
			Type:    MessageTypeKeyGeneration,
			PartyID: keyGenRequest.PartyID,
			From:    n.host.ID(),
			To:      member,
			Payload: json.RawMessage(fmt.Sprintf(`{"key_share": %v}`, keyShares[0])),
		}

		if err := n.SendSecureMessage(context.Background(), keyShareMsg, n.host.Peerstore().PubKey(member)); err != nil {
			return fmt.Errorf("failed to send key share to %s: %w", member, err)
		}
	}

	return nil
}

func (n *Node) handleSigning(msg *Message) error {
	secureMsg, err := n.handleSecureMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to handle secure message: %w", err)
	}

	var signRequest struct {
		PartyID string `json:"party_id"`
		Message []byte `json:"message"`
	}
	if err := json.Unmarshal(secureMsg.Payload, &signRequest); err != nil {
		return fmt.Errorf("failed to unmarshal signing request: %w", err)
	}

	keyShares, err := n.tssHandler.GenerateKeyShares(signRequest.PartyID)
	if err != nil {
		return fmt.Errorf("failed to generate key shares: %w", err)
	}

	signatureShares, err := n.tssHandler.SignMessage(signRequest.PartyID, signRequest.Message, keyShares)
	if err != nil {
		return fmt.Errorf("failed to generate signature shares: %w", err)
	}

	party, err := n.partyMgr.GetParty(signRequest.PartyID)
	if err != nil {
		return fmt.Errorf("failed to get party: %w", err)
	}

	for _, member := range party.Members {
		if member == n.host.ID() {
			continue
		}

		signShareMsg := &Message{
			Type:    MessageTypeSigning,
			PartyID: signRequest.PartyID,
			From:    n.host.ID(),
			To:      member,
			Payload: json.RawMessage(fmt.Sprintf(`{"signature_share": %v}`, signatureShares[0])),
		}

		if err := n.SendSecureMessage(context.Background(), signShareMsg, n.host.Peerstore().PubKey(member)); err != nil {
			return fmt.Errorf("failed to send signature share to %s: %w", member, err)
		}
	}

	return nil
}

// Add this method to the Node struct
func (n *Node) GetAllParties() []*Party {
	return n.partyMgr.GetAllParties()
}
