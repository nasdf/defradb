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
	"github.com/sourcenetwork/defradb/client"
	"github.com/sourcenetwork/defradb/db"
)

// setupDatabase returns the database implementation for the current
// testing state. The database type on the test state is used to
// select the datastore implementation to use.
func setupDatabase(s *state) (client.DB, string, error) {
	dbopts := []db.Option{
		db.WithUpdateEvents(),
		db.WithLensPoolSize(lensPoolSize),
	}

	impl, err := NewInMemoryDB(s.ctx, dbopts...)
	if err != nil {
		return nil, "", err
	}

	return impl, "", nil
}
