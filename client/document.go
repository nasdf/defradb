// Copyright 2022 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package client

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/fxamacker/cbor/v2"
	"github.com/sourcenetwork/defradb/client/request"
)

// This is the main implementation starting point for accessing the internal Document API
// which provides API access to the various operations available for Documents, i.e. CRUD.
//
// Documents in this case refer to the core database type of DefraDB which is a
// "NoSQL Document Datastore".
//
// This section is not concerned with the outer request layer used to interact with the
// Document API, but instead is solely concerned with carrying out the internal API
// operations.
//
// Note: These actions on the outside are deceivingly simple, but require a number
// of complex interactions with the underlying KV Datastore, as well as the
// Merkle CRDT semantics.

// Document is a generalized struct referring to a stored document in the database.
//
// It *can* have a reference to a enforced schema, which is enforced at the time
// of an operation.
//
// Documents are similar to JSON Objects stored in MongoDB, which are collections
// of Fields and Values.
//
// Fields are Key names that point to values.
// Values are literal or complex objects such as strings, integers, or sub documents (objects).
//
// Note: Documents represent the serialized state of the underlying MerkleCRDTs
type Document struct {
	// key is the unique identifier of the document
	key DocKey
	// mutex is used to synchronize document access
	mutex sync.RWMutex
	// fields is a mapping of field names to field kinds
	fields map[string]FieldDescription
	// values is a mapping of field names to values
	values map[string]Value
}

var (
	_ json.Marshaler = (*Document)(nil)
	_ cbor.Marshaler = (*Document)(nil)
)

// NewDocument creates a document which adheres to the given schema.
func NewDocument(schema SchemaDescription) *Document {
	fields := make(map[string]FieldDescription)
	values := make(map[string]Value)

	for _, field := range schema.Fields {
		fields[field.Name] = field

		switch field.Kind {
		case FieldKind_FOREIGN_OBJECT, FieldKind_FOREIGN_OBJECT_ARRAY:
			values[field.Name] = &ReferenceValue{ctype: field.Typ}
		default:
			values[field.Name] = &PrimitiveValue{ctype: field.Typ}
		}
	}

	return &Document{
		fields: fields,
		values: values,
	}
}

// Key returns the unique id for this document.
func (doc *Document) Key() DocKey {
	doc.mutex.RLock()
	defer doc.mutex.RUnlock()
	return doc.key
}

// Set sets the value for the field with the given name and marks it as dirty.
func (doc *Document) Set(name string, newValue any) error {
	doc.mutex.Lock()
	defer doc.mutex.Unlock()

	desc, ok := doc.fields[name]
	if !ok {
		desc, ok = doc.fields[name+request.RelatedObjectID]
	}
	if !ok {
		return fmt.Errorf("field not found: %s", name)
	}
	value, err := ValueToFieldKind(newValue, desc.Kind)
	if err != nil {
		return err
	}
	// value can be a primitive or reference type
	// make sure to chose the correct type
	// so that the document is serialized correctly
	switch desc.Kind {
	case FieldKind_FOREIGN_OBJECT, FieldKind_FOREIGN_OBJECT_ARRAY:
		doc.values[name] = &ReferenceValue{
			value: value,
			ctype: desc.Typ,
			dirty: true,
		}
	default:
		doc.values[name] = &PrimitiveValue{
			value: value,
			ctype: desc.Typ,
			dirty: true,
		}
	}
	// generate the document key when fields are changed
	// this ensures that the document key is immutable
	key, err := GenerateDocKey(doc)
	if err != nil {
		return err
	}
	doc.key = key
	return nil
}

// SetWithJSON sets all the fields of a document using the given json.
//
// Fields set to nil will be marked as deleted.
func (doc *Document) SetWithJSON(patch []byte) error {
	var values map[string]any
	if err := json.Unmarshal([]byte(patch), &values); err != nil {
		return err
	}
	return doc.SetWithMap(values)
}

// SetWithMap sets the fields of a document using the given map.
//
// Fields set to nil will be marked as deleted.
func (doc *Document) SetWithMap(values map[string]any) (err error) {
	for name, value := range values {
		if value == nil {
			err = doc.Delete(name)
		} else {
			err = doc.Set(name, value)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete marks the field with the given name as deleted.
func (doc *Document) Delete(name string) error {
	doc.mutex.Lock()
	defer doc.mutex.Unlock()

	val, ok := doc.values[name]
	if !ok {
		return fmt.Errorf("field not found: %s", name)
	}
	val.Delete()
	// generate the document key when fields are changed
	// this ensures that the document key is immutable
	key, err := GenerateDocKey(doc)
	if err != nil {
		return err
	}
	doc.key = key
	return nil
}

// Value returns the value for the field with the given name.
func (doc *Document) Value(name string) (any, error) {
	doc.mutex.RLock()
	defer doc.mutex.RUnlock()

	val, ok := doc.values[name]
	if !ok {
		return nil, fmt.Errorf("field not found: %s", name)
	}
	return val.Value(), nil
}

// Values gets the document values as a map.
func (doc *Document) Values() map[string]any {
	doc.mutex.RLock()
	defer doc.mutex.RUnlock()

	var values map[string]any
	for name, value := range doc.values {
		values[name] = value.Value()
	}
	return values
}

// Bytes returns the document as a serialzed byte array using CBOR encoding.
func (doc *Document) Bytes() ([]byte, error) {
	doc.mutex.RLock()
	defer doc.mutex.RUnlock()
	return cbor.Marshal(doc)
}

// Clean removes dirty and deleted marks from all values.
func (doc *Document) Clean() {
	doc.mutex.Lock()
	defer doc.mutex.Unlock()

	for _, value := range doc.values {
		value.Clean()
	}
}

func (doc *Document) MarshalCBOR() ([]byte, error) {
	// Important: CanonicalEncOptions ensures consistent serialization of
	// indeterministic datastructures, like Go Maps
	em, err := cbor.CanonicalEncOptions().EncMode()
	if err != nil {
		return nil, err
	}
	return em.Marshal(doc.values)
}

func (doc *Document) MarshalJSON() ([]byte, error) {
	return json.Marshal(doc.values)
}

// DocumentStatus represent the state of the document in the DAG store.
// It can either be `Active“ or `Deleted`.
type DocumentStatus uint8

const (
	// Active is the default state of a document.
	Active DocumentStatus = 1
	// Deleted represents a document that has been marked as deleted. This means that the document
	// can still be in the datastore but a normal request won't return it. The DAG store will still have all
	// the associated links.
	Deleted DocumentStatus = 2
)

var DocumentStatusToString = map[DocumentStatus]string{
	Active:  "Active",
	Deleted: "Deleted",
}

func (dStatus DocumentStatus) UInt8() uint8 {
	return uint8(dStatus)
}

func (dStatus DocumentStatus) IsDeleted() bool {
	return dStatus > 1
}
