// Copyright 2023 Democratized Data Foundation
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
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/sourcenetwork/defradb/api"
)

func ListSchemas(a api.API) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List schema types with their respective fields",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cols, err := a.ListSchemas(cmd.Context())
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			return encoder.Encode(cols)
		},
	}
}

func PatchSchema(a api.API) *cobra.Command {
	return &cobra.Command{
		Use:   "patch [schema]",
		Short: "Patch an existing schema type",
		Long: `Patch an existing schema.

Uses JSON Patch to modify schema types.

Example: patch from string:
  defradb client schema patch '[{ "op": "add", "path": "...", "value": {...} }]'

Example: patch from file:
  cat patch.json | defradb client schema patch -

To learn more about the DefraDB GraphQL Schema Language, refer to https://docs.source.network.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cols, err := a.PatchSchema(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			return encoder.Encode(cols)
		},
	}
}

func LoadSchema(a api.API) *cobra.Command {
	return &cobra.Command{
		Use:   "add [schema]",
		Short: "Add new schema",
		Long: `Add new schema.
	
Example: add from string:
	defradb client schema add 'type Foo { ... }'

Example: add from file:
	cat schema.graphql | defradb client schema add

Learn more about the DefraDB GraphQL Schema Language on https://docs.source.network.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cols, err := a.LoadSchema(cmd.Context(), args[0])
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			return encoder.Encode(cols)
		},
	}
}
