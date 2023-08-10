// Copyright 2023 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package test

import (
	"context"
)

func (s *TestSuite) TestPeerInfo() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	info, err := s.impl.PeerInfo(ctx)

	s.Require().NoError(err)
	s.Assert().Len(info.PeerID, 52)
}
