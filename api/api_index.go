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

type CreateIndexRequest struct {
	Collection string `json:"collection"`
	Fields     string `json:"fields"`
	Name       string `json:"name"`
}

type DropIndexRequest struct {
	Collection string `json:"collection"`
	Name       string `json:"name"`
}

type IndexAPI interface {
	CreateIndex(ctx context.Context, index CreateIndexRequest) (*client.IndexDescription, error)
	DropIndex(ctx context.Context, index DropIndexRequest) error
	ListIndexes(ctx context.Context) (map[string][]client.IndexDescription, error)
	ListIndexesForCollection(ctx context.Context, collection string) ([]client.IndexDescription, error)
}
