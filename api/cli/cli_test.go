// Copyright 2023 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package cli

import (
	"context"
	"testing"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/sourcenetwork/defradb/api/core"
	"github.com/sourcenetwork/defradb/api/test"
	badgerds "github.com/sourcenetwork/defradb/datastore/badger/v4"
	"github.com/sourcenetwork/defradb/db"
	"github.com/sourcenetwork/defradb/net"
)

func TestCLI(t *testing.T) {
	badgerOpts := badgerds.Options{
		Options: badger.DefaultOptions("").WithInMemory(true),
	}

	rootstore, err := badgerds.NewDatastore("", &badgerOpts)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defra, err := db.NewDB(ctx, rootstore, db.WithUpdateEvents())
	require.NoError(t, err)
	t.Cleanup(func() { defra.Close(ctx) })

	node, err := net.NewNode(ctx, defra, net.WithDataPath(t.TempDir()))
	require.NoError(t, err)
	t.Cleanup(func() { node.Close() })

	impl := NewWrapper(core.New(defra, node))
	suite.Run(t, test.NewTestSuite(impl))
}
