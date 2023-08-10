// Copyright 2023 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package core

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/sourcenetwork/defradb/api"
	"github.com/sourcenetwork/defradb/client"
)

func (c *Core) PeerInfo(ctx context.Context) (*api.PeerInfoResponse, error) {
	return &api.PeerInfoResponse{
		PeerID: c.node.PeerID().Pretty(),
	}, nil
}

func (c *Core) SetReplicator(ctx context.Context, req api.ReplicatorRequest) error {
	info, err := peer.AddrInfoFromString(req.Address)
	if err != nil {
		return err
	}
	return c.db.SetReplicator(ctx, client.Replicator{
		Info:    *info,
		Schemas: req.Collections,
	})
}

func (c *Core) DeleteReplicator(ctx context.Context, req api.ReplicatorRequest) error {
	info, err := peer.AddrInfoFromString(req.Address)
	if err != nil {
		return err
	}
	return c.db.DeleteReplicator(ctx, client.Replicator{
		Info:    *info,
		Schemas: req.Collections,
	})
}

func (c *Core) ListReplicators(ctx context.Context) ([]client.Replicator, error) {
	return c.db.GetAllReplicators(ctx)
}

func (c *Core) AddPeerCollection(ctx context.Context, collectionId string) error {
	return c.db.AddP2PCollection(ctx, collectionId)
}

func (c *Core) RemovePeerCollection(ctx context.Context, collectionId string) error {
	return c.db.RemoveP2PCollection(ctx, collectionId)
}

func (c *Core) ListPeerCollections(ctx context.Context) ([]string, error) {
	return c.db.GetAllP2PCollections(ctx)
}
