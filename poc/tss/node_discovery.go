package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const (
	DiscoveryInterval = 1 * time.Minute
	DiscoveryTimeout  = 10 * time.Second
)

type NodeDiscovery struct {
	host       host.Host
	dht        *dht.IpfsDHT
	pubsub     *pubsub.PubSub
	peerChan   chan peer.AddrInfo
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

func NewNodeDiscovery(ctx context.Context, h host.Host) (*NodeDiscovery, error) {
	ctx, cancel := context.WithCancel(ctx)

	// Initialize DHT
	dht, err := dht.New(ctx, h)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create DHT: %w", err)
	}

	// Initialize PubSub
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create PubSub: %w", err)
	}

	nd := &NodeDiscovery{
		host:       h,
		dht:        dht,
		pubsub:     ps,
		peerChan:   make(chan peer.AddrInfo),
		ctx:        ctx,
		cancelFunc: cancel,
	}

	return nd, nil
}

func (nd *NodeDiscovery) Start() error {
	// Start DHT
	if err := nd.dht.Bootstrap(nd.ctx); err != nil {
		return fmt.Errorf("failed to bootstrap DHT: %w", err)
	}

	// Start mDNS discovery
	if err := nd.setupMDNS(); err != nil {
		return fmt.Errorf("failed to setup mDNS: %w", err)
	}

	// Start continuous peer discovery
	nd.wg.Add(1)
	go nd.discoverPeers()

	return nil
}

func (nd *NodeDiscovery) Stop() error {
	nd.cancelFunc()
	nd.wg.Wait()
	close(nd.peerChan)
	return nil
}

func (nd *NodeDiscovery) setupMDNS() error {
	mdnsService := mdns.NewMdnsService(nd.host, "tss-network", nil)
	if err := mdnsService.Start(); err != nil {
		return fmt.Errorf("failed to start mDNS service: %w", err)
	}
	// TODO ???
	//nd.host.Network().Notify(mdnsService.EventHandler())
	return nil
}

func (nd *NodeDiscovery) discoverPeers() {
	defer nd.wg.Done()

	ticker := time.NewTicker(DiscoveryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-nd.ctx.Done():
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(nd.ctx, DiscoveryTimeout)
			peers, err := nd.dht.GetClosestPeers(ctx, "tss-network")
			cancel()

			if err != nil {
				fmt.Printf("Error finding peers: %v\n", err)
				continue
			}

			for _, p := range peers {
				if p == nd.host.ID() {
					continue // Skip self
				}
				select {
				case nd.peerChan <- peer.AddrInfo{ID: p}:
				case <-nd.ctx.Done():
					return
				}
			}
		}
	}
}

func (nd *NodeDiscovery) PeerChan() <-chan peer.AddrInfo {
	return nd.peerChan
}
