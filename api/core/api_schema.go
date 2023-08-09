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

	"github.com/sourcenetwork/defradb/client"
)

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
