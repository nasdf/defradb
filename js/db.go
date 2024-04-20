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
	"context"
	"syscall/js"

	"github.com/lens-vm/lens/host-go/config/model"
	"github.com/sourcenetwork/immutable"

	"github.com/sourcenetwork/defradb/client"
)

type dbFuncs struct {
	db client.DB
}

func (f dbFuncs) new(ctx context.Context) js.Value {
	return js.ValueOf(map[string]any{
		"addSchema":              f.addSchemaFunc(ctx),
		"patchSchem":             f.patchSchemaFunc(ctx),
		"patchCollection":        f.patchCollectionFunc(ctx),
		"setActiveSchemaVersion": f.setActiveSchemaVersionFunc(ctx),
		"addView":                f.addViewFunc(ctx),
		"setMigration":           f.setMigrationFunc(ctx),
		"getCollectionByName":    f.getCollectionByName(ctx),
		"getCollections":         f.getCollections(ctx),
		"getSchemaByVersionID":   f.getSchemaByVersionID(ctx),
		"getSchemas":             f.getSchemas(ctx),
		"getAllIndexes":          f.getAllIndexes(ctx),
		"execRequest":            f.execRequestFunc(ctx),
		"close":                  f.closeFunc(),
	})
}

func (f dbFuncs) addSchemaFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		cols, err := f.db.AddSchema(ctx, args[0].String())
		if err != nil {
			return js.Undefined(), err
		}
		return encodeJS(cols)
	})
}

func (f dbFuncs) patchSchemaFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 3 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		if args[1].Type() != js.TypeObject || args[1].Type() != js.TypeUndefined {
			return js.Undefined(), errInvalidArgs
		}
		if args[2].Type() != js.TypeBoolean {
			return js.Undefined(), errInvalidArgs
		}
		var lens immutable.Option[model.Lens]
		if err := decodeJS(&lens, args[1]); err != nil {
			return js.Undefined(), err
		}
		err := f.db.PatchSchema(ctx, args[0].String(), lens, args[2].Bool())
		return js.Undefined(), err
	})
}

func (f dbFuncs) patchCollectionFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		err := f.db.PatchCollection(ctx, args[0].String())
		return js.Undefined(), err
	})
}

func (f dbFuncs) setActiveSchemaVersionFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		err := f.db.SetActiveSchemaVersion(ctx, args[0].String())
		return js.Undefined(), err
	})
}

func (f dbFuncs) addViewFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 3 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		if args[1].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		if args[2].Type() != js.TypeObject && args[2].Type() != js.TypeUndefined {
			return js.Undefined(), errInvalidArgs
		}
		var lens immutable.Option[model.Lens]
		if err := decodeJS(&lens, args[2]); err != nil {
			return js.Undefined(), err
		}
		cols, err := f.db.AddView(ctx, args[0].String(), args[1].String(), lens)
		if err != nil {
			return js.Undefined(), err
		}
		return encodeJS(cols)
	})
}

func (f dbFuncs) setMigrationFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeObject {
			return js.Undefined(), errInvalidArgs
		}
		var cfg client.LensConfig
		if err := decodeJS(&cfg, args[0]); err != nil {
			return js.Undefined(), err
		}
		err := f.db.SetMigration(ctx, cfg)
		return js.Undefined(), err
	})
}

func (f dbFuncs) getCollectionByName(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		col, err := f.db.GetCollectionByName(ctx, client.CollectionName(args[0].String()))
		if err != nil {
			return js.Undefined(), err
		}
		return collectionFuncs{col}.new(ctx), nil
	})
}

func (f dbFuncs) getCollections(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeObject {
			return js.Undefined(), errInvalidArgs
		}
		var opts client.CollectionFetchOptions
		if err := decodeJS(&opts, args[0]); err != nil {
			return js.Undefined(), err
		}
		cols, err := f.db.GetCollections(ctx, opts)
		if err != nil {
			return js.Undefined(), err
		}
		out := make([]any, len(cols))
		for i, v := range cols {
			out[i] = collectionFuncs{v}.new(ctx)
		}
		return js.ValueOf(out), nil
	})
}

func (f dbFuncs) getSchemaByVersionID(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		schema, err := f.db.GetSchemaByVersionID(ctx, args[0].String())
		if err != nil {
			return js.Undefined(), err
		}
		return encodeJS(schema)
	})
}

func (f dbFuncs) getSchemas(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeObject {
			return js.Undefined(), errInvalidArgs
		}
		var opts client.SchemaFetchOptions
		if err := decodeJS(&opts, args[0]); err != nil {
			return js.Undefined(), err
		}
		schemas, err := f.db.GetSchemas(ctx, opts)
		if err != nil {
			return js.Undefined(), err
		}
		return encodeJS(schemas)
	})
}

func (f dbFuncs) getAllIndexes(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		indexes, err := f.db.GetAllIndexes(ctx)
		if err != nil {
			return js.Undefined(), err
		}
		return encodeJS(indexes)
	})
}

func (f dbFuncs) execRequestFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		res := f.db.ExecRequest(ctx, args[0].String())
		if res.Pub == nil {
			return encodeJS(res.GQL)
		}
		return asyncIterator(func(this js.Value, args []js.Value) (js.Value, error) {
			out, ok := <-res.Pub.Stream()
			val, err := encodeJS(out)
			if err != nil {
				return js.Undefined(), err
			}
			return js.ValueOf(map[string]any{
				"done":  !ok,
				"value": val,
			}), nil
		}), nil
	})
}

func (f dbFuncs) closeFunc() js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		f.db.Close()
		return js.Undefined(), nil
	})
}
