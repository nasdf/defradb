// Copyright 2024 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package tests

import (
	"github.com/sourcenetwork/defradb/net"
	"github.com/sourcenetwork/defradb/tests/clients"
)

// setupClient returns the client implementation for the current
// testing state. The client type on the test state is used to
// select the client implementation to use.
func setupClient(s *state, node *net.Node) (impl clients.Client, err error) {
	return node, nil
}
