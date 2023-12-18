// Copyright 2022 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package mapper

type MutationType int

const (
	NoneMutationType MutationType = iota
	CreateObjects
	UpdateObjects
	DeleteObjects
)

// Mutation represents a request to mutate data stored in Defra.
type Mutation struct {
	// The underlying Select, defining the information requested upon return.
	Select

	// The type of mutation. For example a create request.
	Type MutationType

	// The data to be used for the mutation.  For example, during a create this
	// will be the json representation of the object to be inserted.
	Data map[string]any
}

func (m *Mutation) CloneTo(index int) Requestable {
	return m.cloneTo(index)
}

func (m *Mutation) cloneTo(index int) *Mutation {
	return &Mutation{
		Select: *m.Select.cloneTo(index),
		Type:   m.Type,
		Data:   m.Data,
	}
}
