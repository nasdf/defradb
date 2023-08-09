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

	"github.com/sourcenetwork/defradb/api"
	"github.com/sourcenetwork/defradb/client"
)

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
