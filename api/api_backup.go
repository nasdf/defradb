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

type BackupAPI interface {
	ExportBackup(ctx context.Context, config client.BackupConfig) error
	ImportBackup(ctx context.Context, config client.BackupConfig) error
}
