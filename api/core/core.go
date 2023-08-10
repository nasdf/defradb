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
	"strings"

	dag "github.com/ipfs/boxo/ipld/merkledag"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/sourcenetwork/defradb/api"
	"github.com/sourcenetwork/defradb/client"
	corecrdt "github.com/sourcenetwork/defradb/core/crdt"
	"github.com/sourcenetwork/defradb/net"
)

var _ api.API = (*Core)(nil)

type Core struct {
	db   client.DB
	node *net.Node
}

func New(db client.DB, node *net.Node) *Core {
	return &Core{db, node}
}

func (c *Core) LoadSchema(ctx context.Context, schema string) ([]client.CollectionDescription, error) {
	return c.db.AddSchema(ctx, schema)
}

func (c *Core) PatchSchema(ctx context.Context, patch string) ([]client.CollectionDescription, error) {
	if err := c.db.PatchSchema(ctx, patch); err != nil {
		return nil, err
	}
	cols, err := c.db.GetAllCollections(ctx)
	if err != nil {
		return nil, err
	}
	var colDescs []client.CollectionDescription
	for _, col := range cols {
		colDescs = append(colDescs, col.Description())
	}
	return colDescs, nil
}

func (c *Core) ListSchemas(ctx context.Context) ([]client.CollectionDescription, error) {
	cols, err := c.db.GetAllCollections(ctx)
	if err != nil {
		return nil, err
	}
	var colDescs []client.CollectionDescription
	for _, col := range cols {
		colDescs = append(colDescs, col.Description())
	}
	return colDescs, nil
}

func (c *Core) SetMigration(ctx context.Context, config client.LensConfig) error {
	return c.db.LensRegistry().SetMigration(ctx, config)
}

func (c *Core) GetMigration(ctx context.Context) ([]client.LensConfig, error) {
	return c.db.LensRegistry().Config(ctx)
}

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

func (c *Core) CreateIndex(ctx context.Context, index api.CreateIndexRequest) (*client.IndexDescription, error) {
	col, err := c.db.GetCollectionByName(ctx, index.Collection)
	if err != nil {
		return nil, err
	}
	fields := strings.Split(index.Fields, ",")
	fieldDescriptions := make([]client.IndexedFieldDescription, 0, len(fields))
	for _, field := range fields {
		fieldDescriptions = append(fieldDescriptions, client.IndexedFieldDescription{Name: field})
	}
	indexDesc := client.IndexDescription{
		Name:   index.Name,
		Fields: fieldDescriptions,
	}
	colDesc, err := col.CreateIndex(ctx, indexDesc)
	if err != nil {
		return nil, err
	}
	return &colDesc, nil
}

func (c *Core) DropIndex(ctx context.Context, index api.DropIndexRequest) error {
	col, err := c.db.GetCollectionByName(ctx, index.Collection)
	if err != nil {
		return err
	}
	return col.DropIndex(ctx, index.Name)
}

func (c *Core) ListIndexes(ctx context.Context) (map[string][]client.IndexDescription, error) {
	return c.db.GetAllIndexes(ctx)
}

func (c *Core) ListIndexesForCollection(ctx context.Context, collection string) ([]client.IndexDescription, error) {
	col, err := c.db.GetCollectionByName(ctx, collection)
	if err != nil {
		return nil, err
	}
	return col.GetIndexes(ctx)
}

func (c *Core) Dump(ctx context.Context) error {
	return c.db.PrintDump(ctx)
}

func (c *Core) GetBlock(ctx context.Context, blockCID string) (*api.GetBlockResponse, error) {
	cID, err := parseBlockCID(blockCID)
	if err != nil {
		return nil, err
	}

	block, err := c.db.Blockstore().Get(ctx, cID)
	if err != nil {
		return nil, err
	}

	nd, err := dag.DecodeProtobuf(block.RawData())
	if err != nil {
		return nil, err
	}

	buf, err := nd.MarshalJSON()
	if err != nil {
		return nil, err
	}

	delta, err := corecrdt.LWWRegister{}.DeltaDecode(nd)
	if err != nil {
		return nil, err
	}

	data, err := delta.Marshal()
	if err != nil {
		return nil, err
	}

	return &api.GetBlockResponse{
		Block: string(buf),
		Delta: string(data),
		Value: delta.Value(),
	}, nil
}

func (c *Core) ExportBackup(ctx context.Context, config client.BackupConfig) error {
	err := validateBackupConfig(ctx, &config, c.db)
	if err != nil {
		return err
	}
	return c.db.BasicExport(ctx, &config)
}

func (c *Core) ImportBackup(ctx context.Context, config client.BackupConfig) error {
	err := validateBackupConfig(ctx, &config, c.db)
	if err != nil {
		return err
	}
	return c.db.BasicImport(ctx, config.Filepath)
}
