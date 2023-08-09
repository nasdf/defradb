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

	"github.com/ipfs/boxo/datastore/dshelp"
	dag "github.com/ipfs/boxo/ipld/merkledag"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"

	"github.com/sourcenetwork/defradb/api"
	corecrdt "github.com/sourcenetwork/defradb/core/crdt"
)

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

func parseBlockCID(blockCID string) (cid.Cid, error) {
	cID, err := cid.Decode(blockCID)
	if err == nil {
		return cID, nil
	}
	// If we can't try to parse DSKeyToCID
	// return error if we still can't
	hash, err := dshelp.DsKeyToMultihash(ds.NewKey(blockCID))
	if err != nil {
		return cid.Cid{}, err
	}
	return cid.NewCidV1(cid.Raw, hash), nil
}
