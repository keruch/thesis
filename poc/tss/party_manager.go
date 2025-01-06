package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	PartyFormationTimeout = 2 * time.Minute
	MinPartySize          = 3
	MaxPartySize          = 10
)

type PartyStatus int

const (
	PartyStatusForming PartyStatus = iota
	PartyStatusReady
	PartyStatusActive
	PartyStatusCompleted
	PartyStatusFailed
)

type Party struct {
	ID        string
	Members   []peer.ID
	Threshold int
	Status    PartyStatus
	Operation TSSOperation
}

type TSSOperation int

const (
	TSSOperationKeyGen TSSOperation = iota
	TSSOperationSigning
)

type PartyManager struct {
	parties     map[string]*Party
	peerParties map[peer.ID]map[string]struct{}
	mu          sync.RWMutex
	msgRouter   *MessageRouter
}

func NewPartyManager(msgRouter *MessageRouter) *PartyManager {
	return &PartyManager{
		parties:     make(map[string]*Party),
		peerParties: make(map[peer.ID]map[string]struct{}),
		msgRouter:   msgRouter,
	}
}

func (pm *PartyManager) CreateParty(ctx context.Context, initiator peer.ID, members []peer.ID, threshold int, operation TSSOperation) (*Party, error) {
	if len(members) < MinPartySize || len(members) > MaxPartySize {
		return nil, fmt.Errorf("invalid party size: %d (min: %d, max: %d)", len(members), MinPartySize, MaxPartySize)
	}

	if threshold < 2 || threshold > len(members) {
		return nil, fmt.Errorf("invalid threshold: %d (must be between 2 and party size)", threshold)
	}

	partyID := generatePartyID()
	party := &Party{
		ID:        partyID,
		Members:   members,
		Threshold: threshold,
		Status:    PartyStatusForming,
		Operation: operation,
	}

	pm.mu.Lock()
	pm.parties[partyID] = party
	for _, member := range members {
		if pm.peerParties[member] == nil {
			pm.peerParties[member] = make(map[string]struct{})
		}
		pm.peerParties[member][partyID] = struct{}{}
	}
	pm.mu.Unlock()

	go pm.formParty(ctx, party)

	return party, nil
}

func (pm *PartyManager) formParty(ctx context.Context, party *Party) {
	formationCtx, cancel := context.WithTimeout(ctx, PartyFormationTimeout)
	defer cancel()

	// Simulate party formation process
	time.Sleep(2 * time.Second)

	select {
	case <-formationCtx.Done():
		pm.mu.Lock()
		if party.Status == PartyStatusForming {
			party.Status = PartyStatusFailed
			pm.cleanupParty(party.ID)
		}
		pm.mu.Unlock()
	default:
		pm.mu.Lock()
		party.Status = PartyStatusReady
		pm.mu.Unlock()

		pm.notifyPartyMembers(party)
	}
}

func (pm *PartyManager) GetParty(partyID string) (*Party, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	party, exists := pm.parties[partyID]
	if !exists {
		return nil, fmt.Errorf("party not found: %s", partyID)
	}

	return party, nil
}

func (pm *PartyManager) GetPeerParties(peerID peer.ID) []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	parties := make([]string, 0)
	for partyID := range pm.peerParties[peerID] {
		parties = append(parties, partyID)
	}

	return parties
}

func (pm *PartyManager) UpdatePartyStatus(partyID string, status PartyStatus) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	party, exists := pm.parties[partyID]
	if !exists {
		return fmt.Errorf("party not found: %s", partyID)
	}

	party.Status = status

	if status == PartyStatusCompleted || status == PartyStatusFailed {
		pm.cleanupParty(partyID)
	}

	return nil
}

func (pm *PartyManager) cleanupParty(partyID string) {
	party, exists := pm.parties[partyID]
	if !exists {
		return
	}

	for _, member := range party.Members {
		delete(pm.peerParties[member], partyID)
		if len(pm.peerParties[member]) == 0 {
			delete(pm.peerParties, member)
		}
	}

	delete(pm.parties, partyID)
}

func (pm *PartyManager) notifyPartyMembers(party *Party) {
	// TODO: Implement actual notification logic using MessageRouter
	for _, member := range party.Members {
		fmt.Printf("Notifying peer %s about party %s formation\n", member, party.ID)
	}
}

func generatePartyID() string {
	// TODO: Implement a proper unique ID generation
	return fmt.Sprintf("party-%d", time.Now().UnixNano())
}

// Add this method to the PartyManager struct
func (pm *PartyManager) GetAllParties() []*Party {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	parties := make([]*Party, 0, len(pm.parties))
	for _, party := range pm.parties {
		parties = append(parties, party)
	}
	return parties
}
