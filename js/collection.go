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
	"errors"
	"syscall/js"

	"github.com/sourcenetwork/defradb/client"
	"github.com/sourcenetwork/goji"
)

type collectionFuncs struct {
	col client.Collection
}

func (f collectionFuncs) new(ctx context.Context) js.Value {
	return js.ValueOf(map[string]any{
		"name":         f.nameFunc(),
		"id":           f.idFunc(),
		"schemaRoot":   f.schemaRootFunc(),
		"definition":   f.definitionFunc(),
		"schema":       f.schemaFunc(),
		"create":       f.createFunc(ctx),
		"createMany":   f.createManyFunc(ctx),
		"update":       f.updateFunc(ctx),
		"save":         f.save(ctx),
		"delete":       f.delete(ctx),
		"exists":       f.exists(ctx),
		"get":          f.get(ctx),
		"getAllDocIDs": f.getAllDocIDs(ctx),
		"createIndex":  f.createIndex(ctx),
		"dropIndex":    f.dropIndex(ctx),
		"getIndexes":   f.getIndexes(ctx),
	})
}

func (f collectionFuncs) nameFunc() js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if f.col.Name().HasValue() {
			return js.ValueOf(f.col.Name().Value()), nil
		}
		return js.Undefined(), nil
	})
}

func (f collectionFuncs) idFunc() js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		return js.ValueOf(f.col.ID()), nil
	})
}

func (f collectionFuncs) schemaRootFunc() js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		return js.ValueOf(f.col.SchemaRoot()), nil
	})
}

func (f collectionFuncs) definitionFunc() js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		return goji.MarshalJS(f.col.Definition())
	})
}

func (f collectionFuncs) schemaFunc() js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		return goji.MarshalJS(f.col.Schema())
	})
}

func (f collectionFuncs) createFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeObject {
			return js.Undefined(), errInvalidArgs
		}
		docData := goji.JSON.Stringify(args[0])
		doc, err := client.NewDocFromJSON([]byte(docData), f.col.Definition())
		if err != nil {
			return js.Undefined(), err
		}
		err = f.col.Create(ctx, doc)
		return js.Undefined(), err
	})
}

func (f collectionFuncs) createManyFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeObject {
			return js.Undefined(), errInvalidArgs
		}
		docsData := goji.JSON.Stringify(args[0])
		docs, err := client.NewDocsFromJSON([]byte(docsData), f.col.Definition())
		if err != nil {
			return js.Undefined(), err
		}
		err = f.col.CreateMany(ctx, docs)
		return js.Undefined(), err
	})
}

func (f collectionFuncs) updateFunc(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeObject {
			return js.Undefined(), errInvalidArgs
		}
		docIDValue := args[0].Get("_docID")
		if docIDValue.Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		docID, err := client.NewDocIDFromString(docIDValue.String())
		if err != nil {
			return js.Undefined(), err
		}
		doc, err := f.col.Get(ctx, docID, false)
		if err != nil {
			return js.Undefined(), err
		}
		docData := goji.JSON.Stringify(args[0])
		if err := doc.SetWithJSON([]byte(docData)); err != nil {
			return js.Undefined(), err
		}
		err = f.col.Update(ctx, doc)
		return js.Undefined(), err
	})
}

func (f collectionFuncs) save(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeObject {
			return js.Undefined(), errInvalidArgs
		}
		docIDValue := args[0].Get("_docID")
		if docIDValue.Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		docID, err := client.NewDocIDFromString(docIDValue.String())
		if err != nil {
			return js.Undefined(), err
		}
		_, err = f.col.Get(ctx, docID, true)
		if err == nil {
			return this.Call("update", args[0]), nil
		}
		if errors.Is(err, client.ErrDocumentNotFoundOrNotAuthorized) {
			return this.Call("create", args[0]), nil
		}
		return js.Undefined(), err
	})
}

func (f collectionFuncs) delete(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.ValueOf(false), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.ValueOf(false), errInvalidArgs
		}
		docID, err := client.NewDocIDFromString(args[0].String())
		if err != nil {
			return js.ValueOf(false), err
		}
		exists, err := f.col.Delete(ctx, docID)
		return js.ValueOf(exists), err
	})
}

func (f collectionFuncs) exists(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.ValueOf(false), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.ValueOf(false), errInvalidArgs
		}
		docID, err := client.NewDocIDFromString(args[0].String())
		if err != nil {
			return js.ValueOf(false), err
		}
		exists, err := f.col.Exists(ctx, docID)
		return js.ValueOf(exists), err
	})
}

func (f collectionFuncs) get(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 2 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		if args[1].Type() != js.TypeBoolean {
			return js.Undefined(), errInvalidArgs
		}
		docID, err := client.NewDocIDFromString(args[0].String())
		if err != nil {
			return js.Undefined(), err
		}
		doc, err := f.col.Get(ctx, docID, args[1].Bool())
		if err != nil {
			return js.Undefined(), err
		}
		docMap, err := doc.ToMap()
		if err != nil {
			return js.Undefined(), err
		}
		return goji.MarshalJS(docMap)
	})
}

func (f collectionFuncs) getAllDocIDs(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		resCh, err := f.col.GetAllDocIDs(ctx)
		if err != nil {
			return js.Undefined(), err
		}
		return asyncIterator(func(this js.Value, args []js.Value) (js.Value, error) {
			res, ok := <-resCh
			if res.Err != nil {
				return js.Undefined(), res.Err
			}
			val := map[string]any{
				"done": !ok,
			}
			if !ok {
				return js.ValueOf(val), nil
			}
			val["value"] = res.ID.String()
			return js.ValueOf(val), nil
		}), nil
	})
}

func (f collectionFuncs) createIndex(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeObject {
			return js.Undefined(), errInvalidArgs
		}
		var desc client.IndexDescription
		if err := goji.UnmarshalJS(args[0], &desc); err != nil {
			return js.Undefined(), err
		}
		out, err := f.col.CreateIndex(ctx, desc)
		if err != nil {
			return js.Undefined(), err
		}
		return goji.MarshalJS(out)
	})
}

func (f collectionFuncs) dropIndex(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		if len(args) < 1 {
			return js.Undefined(), errInvalidArgs
		}
		if args[0].Type() != js.TypeString {
			return js.Undefined(), errInvalidArgs
		}
		err := f.col.DropIndex(ctx, args[0].String())
		return js.Undefined(), err
	})
}

func (f collectionFuncs) getIndexes(ctx context.Context) js.Func {
	return async(func(this js.Value, args []js.Value) (js.Value, error) {
		indexes, err := f.col.GetIndexes(ctx)
		if err != nil {
			return js.Undefined(), err
		}
		return goji.MarshalJS(indexes)
	})
}
