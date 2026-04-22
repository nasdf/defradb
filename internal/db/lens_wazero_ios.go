// Copyright 2026 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

// On iOS, wasmtime cannot be linked because its C library doesn't ship for
// the iOS SDK. This file mirrors lens_wazero_android.go by registering
// wazero (pure Go) as both the default lens runtime and the named "wazero"
// runtime.

//go:build ios

package db

import (
	"github.com/sourcenetwork/lens/host-go/engine/module"
	"github.com/sourcenetwork/lens/host-go/runtimes/wazero"
)

const Wazero LensRuntimeType = "wazero"

func init() {
	runtimeConstructors[DefaultLens] = func() module.Runtime { return wazero.New() }
	runtimeConstructors[Wazero] = func() module.Runtime { return wazero.New() }
}
