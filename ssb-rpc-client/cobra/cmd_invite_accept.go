/*
 * SSB Tools
 *     Copyright (C) 2023 Chronicle Labs, Inc.
 *
 *     This program is free software: you can redistribute it and/or modify
 *     it under the terms of the GNU Affero General Public License as published
 *     by the Free Software Foundation, either version 3 of the License, or
 *     (at your option) any later version.
 *
 *     This program is distributed in the hope that it will be useful,
 *     but WITHOUT ANY WARRANTY; without even the implied warranty of
 *     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *     GNU Affero General Public License for more details.
 *
 *     You should have received a copy of the GNU Affero General Public License
 *     along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package cobra

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InviteAccept(opts *Options) *cobra.Command {
	return &cobra.Command{
		Use:  "invite.accept",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := opts.SSBConfig()
			if err != nil {
				return err
			}
			c, err := conf.Client(cmd.Context())
			if err != nil {
				return err
			}

			invites, err := c.InviteAccept(args[0])
			if err != nil {
				return err
			}
			fmt.Println(string(invites))
			return nil
		},
	}
}
