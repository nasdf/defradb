// Copyright 2024 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

//go:build js

package js

import (
	"context"
	"log"
	"syscall/js"

	"github.com/sourcenetwork/defradb/datastore/memory"
	"github.com/sourcenetwork/defradb/db"
	"github.com/sourcenetwork/defradb/net"
)

// Open returns a new DefraDB client instance.
func Open(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		d, err := db.NewDB(ctx, memory.NewDatastore(ctx), db.WithUpdateEvents())
		if err != nil {
			return js.Undefined(), err
		}
		node, err := net.NewNode(ctx, d)
		if err != nil {
			return js.Undefined(), err
		}
		go func() {
			if err := node.Start(); err != nil {
				log.Fatal("failed to start node: " + err.Error())
			}
		}()
		return dbFuncs{node}.new(ctx), nil
	})
}
