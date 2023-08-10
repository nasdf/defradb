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
	"encoding/json"
	"io"

	"github.com/sourcenetwork/defradb/api"
	"github.com/sourcenetwork/defradb/client"
	"github.com/spf13/cobra"
)

var _ api.API = (*Wrapper)(nil)

type Wrapper struct {
	api api.API
}

func NewWrapper(api api.API) *Wrapper {
	return &Wrapper{api}
}

func (w *Wrapper) execute(ctx context.Context, cmd *cobra.Command) ([]byte, []byte, error) {
	var (
		stdErr []byte
		stdOut []byte
	)

	stdErrReader, stdErrWriter := io.Pipe()
	stdOutReader, stdOutWriter := io.Pipe()

	cmd.SetOut(stdOutWriter)
	cmd.SetErr(stdErrWriter)

	go func() {
		err := cmd.ExecuteContext(ctx)
		stdErrWriter.CloseWithError(err)
		stdOutWriter.CloseWithError(err)
	}()

	stdOut, err := io.ReadAll(stdOutReader)
	if err != nil {
		return nil, nil, err
	}

	stdErr, err = io.ReadAll(stdErrReader)
	if err != nil {
		return nil, nil, err
	}

	return stdOut, stdErr, nil
}

func (w *Wrapper) LoadSchema(ctx context.Context, schema string) ([]client.CollectionDescription, error) {
	return nil, nil
}

func (w *Wrapper) PatchSchema(ctx context.Context, patch string) ([]client.CollectionDescription, error) {
	return nil, nil
}

func (w *Wrapper) ListSchemas(ctx context.Context) ([]client.CollectionDescription, error) {
	return nil, nil
}

func (w *Wrapper) SetMigration(ctx context.Context, config client.LensConfig) error {
	return nil
}

func (w *Wrapper) GetMigration(ctx context.Context) ([]client.LensConfig, error) {
	return nil, nil
}

func (w *Wrapper) PeerInfo(ctx context.Context) (*api.PeerInfoResponse, error) {
	stdOut, _, err := w.execute(ctx, PeerInfo(w.api))
	if err != nil {
		return nil, err
	}
	var res api.PeerInfoResponse
	if err := json.Unmarshal(stdOut, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (w *Wrapper) SetReplicator(ctx context.Context, req api.ReplicatorRequest) error {
	return nil
}

func (w *Wrapper) DeleteReplicator(ctx context.Context, req api.ReplicatorRequest) error {
	return nil
}

func (w *Wrapper) ListReplicators(ctx context.Context) ([]client.Replicator, error) {
	return nil, nil
}

func (w *Wrapper) AddPeerCollection(ctx context.Context, collectionId string) error {
	return nil
}

func (w *Wrapper) RemovePeerCollection(ctx context.Context, collectionId string) error {
	return nil
}

func (w *Wrapper) ListPeerCollections(ctx context.Context) ([]string, error) {
	return nil, nil
}

func (w *Wrapper) CreateIndex(ctx context.Context, index api.CreateIndexRequest) (*client.IndexDescription, error) {
	return nil, nil
}

func (w *Wrapper) DropIndex(ctx context.Context, index api.DropIndexRequest) error {
	return nil
}

func (w *Wrapper) ListIndexes(ctx context.Context) (map[string][]client.IndexDescription, error) {
	return nil, nil
}

func (w *Wrapper) ListIndexesForCollection(ctx context.Context, collection string) ([]client.IndexDescription, error) {
	return nil, nil
}

func (w *Wrapper) Dump(ctx context.Context) error {
	return nil
}

func (w *Wrapper) GetBlock(ctx context.Context, blockCID string) (*api.GetBlockResponse, error) {
	return nil, nil
}

func (w *Wrapper) ExportBackup(ctx context.Context, config client.BackupConfig) error {
	return nil
}

func (w *Wrapper) ImportBackup(ctx context.Context, config client.BackupConfig) error {
	return nil
}
