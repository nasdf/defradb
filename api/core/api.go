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
	"github.com/sourcenetwork/defradb/api"
	"github.com/sourcenetwork/defradb/client"
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
