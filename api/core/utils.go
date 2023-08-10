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
	"os"
	"strings"

	"github.com/ipfs/boxo/datastore/dshelp"
	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	"github.com/sourcenetwork/defradb/client"
	"github.com/sourcenetwork/defradb/errors"
)

func validateBackupConfig(ctx context.Context, cfg *client.BackupConfig, db client.DB) error {
	if !isValidPath(cfg.Filepath) {
		return errors.New("invalid file path")
	}

	if cfg.Format != "" && strings.ToLower(cfg.Format) != "json" {
		return errors.New("only JSON format is supported at the moment")
	}
	for _, colName := range cfg.Collections {
		_, err := db.GetCollectionByName(ctx, colName)
		if err != nil {
			return errors.Wrap("collection does not exist", err)
		}
	}
	return nil
}

func isValidPath(filepath string) bool {
	// if a file exists, return true
	if _, err := os.Stat(filepath); err == nil {
		return true
	}

	// if not, attempt to write to the path and if successful,
	// remove the file and return true
	var d []byte
	if err := os.WriteFile(filepath, d, 0o644); err == nil {
		_ = os.Remove(filepath)
		return true
	}

	return false
}

func parseBlockCID(blockCID string) (cid.Cid, error) {
	cID, err := cid.Decode(blockCID)
	if err == nil {
		return cID, nil
	}
	// If we can't try to parse DSKeyToCID
	// return error if we still can't
	hash, err := dshelp.DsKeyToMultihash(ds.NewKey(blockCID))
	if err != nil {
		return cid.Cid{}, err
	}
	return cid.NewCidV1(cid.Raw, hash), nil
}
