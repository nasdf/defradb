// Copyright 2023 Democratized Data Foundation
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package test

import (
	"github.com/stretchr/testify/suite"

	"github.com/sourcenetwork/defradb/api"
)

type TestSuite struct {
	suite.Suite
	impl api.API
}

func NewTestSuite(impl api.API) *TestSuite {
	return &TestSuite{impl: impl}
}
