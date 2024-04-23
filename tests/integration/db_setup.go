// Copyright 2024 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

//go:build !js

package tests

import (
	"context"
	"fmt"
	"testing"

	badger "github.com/sourcenetwork/badger/v4"
	"github.com/sourcenetwork/defradb/client"
	badgerds "github.com/sourcenetwork/defradb/datastore/badger/v4"
	"github.com/sourcenetwork/defradb/db"
	changeDetector "github.com/sourcenetwork/defradb/tests/change_detector"
)

func NewBadgerMemoryDB(ctx context.Context, dbopts ...db.Option) (client.DB, error) {
	opts := badgerds.Options{
		Options: badger.DefaultOptions("").WithInMemory(true),
	}
	rootstore, err := badgerds.NewDatastore("", &opts)
	if err != nil {
		return nil, err
	}
	dbopts = append(dbopts, db.WithACPInMemory())
	db, err := db.NewDB(ctx, rootstore, dbopts...)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewBadgerFileDB(ctx context.Context, t testing.TB, dbopts ...db.Option) (client.DB, string, error) {
	var dbPath string
	switch {
	case databaseDir != "":
		// restarting database
		dbPath = databaseDir

	case changeDetector.Enabled:
		// change detector
		dbPath = changeDetector.DatabaseDir(t)

	default:
		// default test case
		dbPath = t.TempDir()
	}

	opts := &badgerds.Options{
		Options: badger.DefaultOptions(dbPath),
	}

	rootstore, err := badgerds.NewDatastore(dbPath, opts)
	if err != nil {
		return nil, "", err
	}

	dbopts = append(dbopts, db.WithACP(dbPath))
	db, err := db.NewDB(ctx, rootstore, dbopts...)
	if err != nil {
		return nil, "", err
	}

	return db, dbPath, err
}

// setupDatabase returns the database implementation for the current
// testing state. The database type on the test state is used to
// select the datastore implementation to use.
func setupDatabase(s *state) (impl client.DB, path string, err error) {
	dbopts := []db.Option{
		db.WithUpdateEvents(),
		db.WithLensPoolSize(lensPoolSize),
	}

	switch s.dbt {
	case badgerIMType:
		impl, err = NewBadgerMemoryDB(s.ctx, dbopts...)

	case badgerFileType:
		impl, path, err = NewBadgerFileDB(s.ctx, s.t, dbopts...)

	case defraIMType:
		impl, err = NewInMemoryDB(s.ctx, dbopts...)

	default:
		err = fmt.Errorf("invalid database type: %v", s.dbt)
	}

	if err != nil {
		return nil, "", err
	}
	return
}
