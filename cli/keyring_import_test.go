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
	"encoding/hex"
	"path/filepath"
	"testing"

	"github.com/sourcenetwork/defradb/crypto"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyringImport(t *testing.T) {
	rootdir := t.TempDir()
	readPassword = func(_ *cobra.Command, _ string) ([]byte, error) {
		return []byte("secret"), nil
	}

	keyBytes, err := crypto.GenerateAES256()
	require.NoError(t, err)
	keyHex := hex.EncodeToString(keyBytes)

	cmd := NewDefraCommand()
	cmd.SetArgs([]string{"keyring", "import", "--rootdir", rootdir, encryptionKeyName, keyHex})

	err = cmd.Execute()
	require.NoError(t, err)

	assert.FileExists(t, filepath.Join(rootdir, "keys", encryptionKeyName))
}
