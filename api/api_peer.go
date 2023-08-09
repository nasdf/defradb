// Copyright 2023 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package api

import (
	"context"

	"github.com/sourcenetwork/defradb/client"
)

type PeerInfoResponse struct {
	PeerID string `json:"peerID"`
}

type ReplicatorRequest struct {
	Collections []string `json:"collections"`
	Address     string   `json:"address"`
}

type PeerAPI interface {
	PeerInfo(ctx context.Context) (*PeerInfoResponse, error)

	SetReplicator(ctx context.Context, req ReplicatorRequest) error
	DeleteReplicator(ctx context.Context, req ReplicatorRequest) error
	ListReplicators(ctx context.Context) ([]client.Replicator, error)

	AddPeerCollection(ctx context.Context, collectionId string) error
	RemovePeerCollection(ctx context.Context, collectionId string) error
	ListPeerCollections(ctx context.Context) ([]string, error)
}
