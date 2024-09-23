// Copyright 2024 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package crypto

import (
	"crypto/rand"
	"errors"
	"io"
	"os"
	"strings"
)

const AESNonceSize = 12

var generateNonceFunc = generateNonce

func generateNonce() ([]byte, error) {
	nonce := make([]byte, AESNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return nonce, nil
}

// generateTestNonce generates a deterministic nonce for testing.
func generateTestNonce() ([]byte, error) {
	nonce := []byte("deterministic nonce for testing")

	if len(nonce) < AESNonceSize {
		return nil, errors.New("nonce length is longer than available deterministic nonce")
	}

	return nonce[:AESNonceSize], nil
}

func init() {
	arg := os.Args[0]
	// If the binary is a test binary, use a deterministic nonce.
	// TODO: We should try to find a better way to detect this https://github.com/sourcenetwork/defradb/issues/2801
	if strings.HasSuffix(arg, ".test") || strings.Contains(arg, "/defradb/tests/") {
		generateNonceFunc = generateTestNonce
	}
}
