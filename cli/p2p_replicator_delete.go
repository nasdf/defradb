// Copyright 2022 Democratized Data Foundation
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
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"

	"github.com/sourcenetwork/defradb/client"
	"github.com/sourcenetwork/defradb/config"
)

func MakeP2PReplicatorDeleteCommand(cfg *config.Config) *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "delete <peer>",
		Short: "Delete a replicator. It will stop synchronizing",
		Long:  `Delete a replicator. It will stop synchronizing.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store := cmd.Context().Value(storeContextKey).(client.Store)

			addr, err := peer.AddrInfoFromString(args[0])
			if err != nil {
				return err
			}
			return store.DeleteReplicator(cmd.Context(), client.Replicator{Info: *addr})
		},
	}
	return cmd
}
