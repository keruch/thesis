package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	TSSTopicName = "tss-messages"
)

type MessageType int

const (
	MessageTypePartyFormation MessageType = iota
	MessageTypeKeyGeneration
	MessageTypeSigning
)

type Message struct {
	Type    MessageType     `json:"type"`
	PartyID string          `json:"party_id"`
	From    peer.ID         `json:"from"`
	To      peer.ID         `json:"to"`
	Payload json.RawMessage `json:"payload"`
}

type MessageHandler func(msg *Message) error

type MessageRouter struct {
	host         host.Host
	pubsub       *pubsub.PubSub
	topic        *pubsub.Topic
	subscription *pubsub.Subscription
	handlers     map[MessageType]MessageHandler
	mu           sync.RWMutex
}

func NewMessageRouter(h host.Host) *MessageRouter {
	return &MessageRouter{
		host:     h,
		handlers: make(map[MessageType]MessageHandler),
	}
}

func (mr *MessageRouter) Start(ctx context.Context) error {
	ps, err := pubsub.NewGossipSub(ctx, mr.host)
	if err != nil {
		return fmt.Errorf("failed to create pubsub: %w", err)
	}
	mr.pubsub = ps

	topic, err := mr.pubsub.Join(TSSTopicName)
	if err != nil {
		return fmt.Errorf("failed to join topic: %w", err)
	}
	mr.topic = topic

	subscription, err := topic.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}
	mr.subscription = subscription

	go mr.handleMessages(ctx)

	return nil
}

func (mr *MessageRouter) Stop() error {
	if mr.subscription != nil {
		mr.subscription.Cancel()
	}
	if mr.topic != nil {
		return mr.topic.Close()
	}
	return nil
}

func (mr *MessageRouter) RegisterHandler(msgType MessageType, handler MessageHandler) {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.handlers[msgType] = handler
}

func (mr *MessageRouter) SendMessage(ctx context.Context, msg *Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return mr.topic.Publish(ctx, data)
}

func (mr *MessageRouter) handleMessages(ctx context.Context) {
	for {
		msg, err := mr.subscription.Next(ctx)
		if err != nil {
			fmt.Printf("Error receiving message: %v\n", err)
			continue
		}

		// Ignore messages from self
		if msg.ReceivedFrom == mr.host.ID() {
			continue
		}

		var message Message
		if err := json.Unmarshal(msg.Data, &message); err != nil {
			fmt.Printf("Error unmarshaling message: %v\n", err)
			continue
		}

		mr.mu.RLock()
		handler, exists := mr.handlers[message.Type]
		mr.mu.RUnlock()

		if !exists {
			fmt.Printf("No handler registered for message type: %d\n", message.Type)
			continue
		}

		if err := handler(&message); err != nil {
			fmt.Printf("Error handling message: %v\n", err)
		}
	}
}
