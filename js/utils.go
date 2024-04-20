// Copyright 2024 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

//go:build js

package js

import (
	"errors"
	"syscall/js"

	"github.com/sourcenetwork/goji"
)

var errInvalidArgs = errors.New("invalid arguments")

func async(fn func(this js.Value, args []js.Value) (js.Value, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		return goji.PromiseOf(func(resolve, reject func(value js.Value)) {
			res, err := fn(this, args)
			if err != nil {
				reject(js.Value(goji.Error.New(err.Error())))
			} else {
				resolve(res)
			}
		})
	})
}

func asyncIterator(fn func(this js.Value, args []js.Value) (js.Value, error)) js.Value {
	return js.ValueOf(map[string]any{
		"next": async(fn),
	})
}
