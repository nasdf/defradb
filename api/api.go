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

type SetMigrationResponse struct {
	Result string `json:"result"`
}

type CreateIndexRequest struct {
	Collection string `json:"collection"`
	Fields     string `json:"fields"`
	Name       string `json:"name"`
}

type DropIndexRequest struct {
	Collection string `json:"collection"`
	Name       string `json:"name"`
}

type GetBlockResponse struct {
	Block string `json:"block"`
	Delta string `json:"delta"`
	Value any    `json:"value"`
}

type API interface {
	LoadSchema(ctx context.Context, schema string) ([]client.CollectionDescription, error)
	PatchSchema(ctx context.Context, patch string) ([]client.CollectionDescription, error)
	ListSchemas(ctx context.Context) ([]client.CollectionDescription, error)

	SetMigration(ctx context.Context, config client.LensConfig) error
	GetMigration(ctx context.Context) ([]client.LensConfig, error)

	PeerInfo(ctx context.Context) (*PeerInfoResponse, error)

	SetReplicator(ctx context.Context, req ReplicatorRequest) error
	DeleteReplicator(ctx context.Context, req ReplicatorRequest) error
	ListReplicators(ctx context.Context) ([]client.Replicator, error)

	AddPeerCollection(ctx context.Context, collectionId string) error
	RemovePeerCollection(ctx context.Context, collectionId string) error
	ListPeerCollections(ctx context.Context) ([]string, error)

	CreateIndex(ctx context.Context, index CreateIndexRequest) (*client.IndexDescription, error)
	DropIndex(ctx context.Context, index DropIndexRequest) error
	ListIndexes(ctx context.Context) (map[string][]client.IndexDescription, error)
	ListIndexesForCollection(ctx context.Context, collection string) ([]client.IndexDescription, error)

	Dump(ctx context.Context) error
	GetBlock(ctx context.Context, blockCID string) (*GetBlockResponse, error)

	ExportBackup(ctx context.Context, config client.BackupConfig) error
	ImportBackup(ctx context.Context, config client.BackupConfig) error
}
