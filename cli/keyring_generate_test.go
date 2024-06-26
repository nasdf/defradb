// Copyright 2024 Democratized Data Foundation
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
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyringGenerate(t *testing.T) {
	rootdir := t.TempDir()
	readPassword = func(_ *cobra.Command, _ string) ([]byte, error) {
		return []byte("secret"), nil
	}

	cmd := NewDefraCommand()
	cmd.SetArgs([]string{"keyring", "generate", "--rootdir", rootdir})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(rootdir, "keys", encryptionKeyName))
	assert.FileExists(t, filepath.Join(rootdir, "keys", peerKeyName))
}

func TestKeyringGenerateNoEncryptionKey(t *testing.T) {
	rootdir := t.TempDir()
	readPassword = func(_ *cobra.Command, _ string) ([]byte, error) {
		return []byte("secret"), nil
	}

	cmd := NewDefraCommand()
	cmd.SetArgs([]string{"keyring", "generate", "--no-encryption-key", "--rootdir", rootdir})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.NoFileExists(t, filepath.Join(rootdir, "keys", encryptionKeyName))
	assert.FileExists(t, filepath.Join(rootdir, "keys", peerKeyName))
}

func TestKeyringGenerateNoPeerKey(t *testing.T) {
	rootdir := t.TempDir()
	readPassword = func(_ *cobra.Command, _ string) ([]byte, error) {
		return []byte("secret"), nil
	}

	cmd := NewDefraCommand()
	cmd.SetArgs([]string{"keyring", "generate", "--no-peer-key", "--rootdir", rootdir})

	err := cmd.Execute()
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(rootdir, "keys", encryptionKeyName))
	assert.NoFileExists(t, filepath.Join(rootdir, "keys", peerKeyName))
}
