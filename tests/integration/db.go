// Copyright 2023 Democratized Data Foundation
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
	"context"
	"os"
	"strconv"

	"github.com/sourcenetwork/defradb/client"
	"github.com/sourcenetwork/defradb/datastore/memory"
	"github.com/sourcenetwork/defradb/db"
	changeDetector "github.com/sourcenetwork/defradb/tests/change_detector"
)

type DatabaseType string

const (
	memoryBadgerEnvName   = "DEFRA_BADGER_MEMORY"
	fileBadgerEnvName     = "DEFRA_BADGER_FILE"
	fileBadgerPathEnvName = "DEFRA_BADGER_FILE_PATH"
	inMemoryEnvName       = "DEFRA_IN_MEMORY"
)

const (
	badgerIMType   DatabaseType = "badger-in-memory"
	defraIMType    DatabaseType = "defra-memory-datastore"
	badgerFileType DatabaseType = "badger-file-system"
)

var (
	badgerInMemory bool
	badgerFile     bool
	inMemoryStore  bool
	databaseDir    string
)

func init() {
	// We use environment variables instead of flags `go test ./...` throws for all packages
	// that don't have the flag defined
	badgerFile, _ = strconv.ParseBool(os.Getenv(fileBadgerEnvName))
	badgerInMemory, _ = strconv.ParseBool(os.Getenv(memoryBadgerEnvName))
	inMemoryStore, _ = strconv.ParseBool(os.Getenv(inMemoryEnvName))

	if changeDetector.Enabled {
		// Change detector only uses badger file db type.
		badgerFile = true
		badgerInMemory = false
		inMemoryStore = false
	} else if !badgerInMemory && !badgerFile && !inMemoryStore {
		// Default is to test all but filesystem db types.
		badgerFile = false
		badgerInMemory = true
		inMemoryStore = true
	}
}

func NewInMemoryDB(ctx context.Context, dbopts ...db.Option) (client.DB, error) {
	dbopts = append(dbopts, db.WithACPInMemory())
	db, err := db.NewDB(ctx, memory.NewDatastore(ctx), dbopts...)
	if err != nil {
		return nil, err
	}
	return db, nil
}
