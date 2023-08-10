// Copyright 2023 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package http

import (
	"context"
	"net/http"

	"github.com/sourcenetwork/defradb/api"
	"github.com/sourcenetwork/defradb/client"
)

var _ api.API = (*Client)(nil)

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient() *Client {
	client := http.DefaultClient
	return &Client{client, ""}
}

func (c *Client) ExportBackup(ctx context.Context, config client.BackupConfig) error {
	return nil
}

func (c *Client) ImportBackup(ctx context.Context, config client.BackupConfig) error {
	return nil
}

func (c *Client) Dump(ctx context.Context) error {
	return nil
}

func (c *Client) GetBlock(ctx context.Context, blockCID string) (*api.GetBlockResponse, error) {
	return nil, nil
}

func (c *Client) CreateIndex(ctx context.Context, index api.CreateIndexRequest) (*client.IndexDescription, error) {
	return nil, nil
}

func (c *Client) DropIndex(ctx context.Context, index api.DropIndexRequest) error {
	return nil
}

func (c *Client) ListIndexes(ctx context.Context) (map[string][]client.IndexDescription, error) {
	return nil, nil
}

func (c *Client) ListIndexesForCollection(ctx context.Context, collection string) ([]client.IndexDescription, error) {
	return nil, nil
}

func (c *Client) PeerInfo(ctx context.Context) (*api.PeerInfoResponse, error) {
	return nil, nil
}

func (c *Client) SetReplicator(ctx context.Context, req api.ReplicatorRequest) error {
	return nil
}

func (c *Client) DeleteReplicator(ctx context.Context, req api.ReplicatorRequest) error {
	return nil
}

func (c *Client) ListReplicators(ctx context.Context) ([]client.Replicator, error) {
	return nil, nil
}

func (c *Client) AddPeerCollection(ctx context.Context, collectionId string) error {
	return nil
}

func (c *Client) RemovePeerCollection(ctx context.Context, collectionId string) error {
	return nil
}

func (c *Client) ListPeerCollections(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (c *Client) LoadSchema(ctx context.Context, schema string) ([]client.CollectionDescription, error) {
	return nil, nil
}

func (c *Client) PatchSchema(ctx context.Context, patch string) ([]client.CollectionDescription, error) {
	return nil, nil
}

func (c *Client) ListSchemas(ctx context.Context) ([]client.CollectionDescription, error) {
	return nil, nil
}

func (c *Client) SetMigration(ctx context.Context, config client.LensConfig) error {
	return nil
}

func (c *Client) GetMigration(ctx context.Context) ([]client.LensConfig, error) {
	return nil, nil
}
