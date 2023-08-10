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

	"github.com/sourcenetwork/defradb/api"
	"github.com/spf13/cobra"
)

func PeerInfo(a api.API) *cobra.Command {
	return &cobra.Command{
		Use:   "peerid",
		Short: "Get the PeerID of the node",
		Long:  `Get the PeerID of the node.`,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			info, err := a.PeerInfo(cmd.Context())
			if err != nil {
				return err
			}

			encoder := json.NewEncoder(cmd.OutOrStdout())
			return encoder.Encode(info)
		},
	}
}
