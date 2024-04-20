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
	"encoding/json"
	"errors"
	"syscall/js"
)

var (
	errInvalidArgs = errors.New("invalid arguments")
)

// async returns a js.Func that wraps the given func in a promise.
func async(fn func(this js.Value, args []js.Value) (js.Value, error)) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		var (
			callback js.Func
			resolve  func(js.Value)
			reject   func(js.Value)
		)
		// callback is called with resolve and reject funcs as arguments
		callback = js.FuncOf(func(this js.Value, args []js.Value) any {
			callback.Release()
			resolve = func(val js.Value) {
				args[0].Invoke(val)
			}
			reject = func(val js.Value) {
				args[1].Invoke(val)
			}
			return js.Undefined()
		})
		promise := newPromise(callback)
		// inner func is run in a go routine so that
		// the main JavaScript thread is not blocked
		go func() {
			val, err := fn(this, args)
			if err != nil {
				reject(newError(err.Error()))
			} else {
				resolve(val)
			}
		}()
		return promise
	})
}

func asyncIterator(fn func(this js.Value, args []js.Value) (js.Value, error)) js.Value {
	return js.ValueOf(map[string]any{
		"next": async(fn),
	})
}

// encodeJS marshals the given value into a js.Value.
func encodeJS(v any) (js.Value, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return js.Undefined(), err
	}
	return jsonParse(js.ValueOf(string(data))), nil
}

// decodeJS unmarshals the given js.Value into the given out value.
func decodeJS(out any, value js.Value) error {
	res := jsonStringify(value)
	return json.Unmarshal([]byte(res), out)
}

// newPromise returns a new JavaScript promise.
func newPromise(callback js.Func) js.Value {
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/Promise
	return js.Global().Get("Promise").New(callback)
}

// newError returns a new JavaScript error.
func newError(message string) js.Value {
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error/Error
	return js.Global().Get("Error").New(message)
}

func jsonParse(value js.Value) js.Value {
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON/parse
	return js.Global().Get("JSON").Call("parse", value)
}

// jsonStringify calls JSON.stringify for the given value
func jsonStringify(value js.Value) string {
	// https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON/stringify
	return js.Global().Get("JSON").Call("stringify", value).String()
}
