//  Copyright (C) 2020 Maker Ecosystem Growth Holdings, INC.
//
//  This program is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Affero General Public License as
//  published by the Free Software Foundation, either version 3 of the
//  License, or (at your option) any later version.
//
//  This program is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Affero General Public License for more details.
//
//  You should have received a copy of the GNU Affero General Public License
//  along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cobra

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Hist returns a cobra command that prints the history of a feed.
func Hist(opts *Options) *cobra.Command {
	var id string
	var seq, limit, lt, gt int64
	var live, reverse, keys, values, private bool
	cmd := &cobra.Command{
		Use: "hist --id {feedId} [--seq n] [--live]",
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := opts.SSBConfig()
			if err != nil {
				return err
			}
			c, err := conf.Client(cmd.Context())
			if err != nil {
				return err
			}
			ch, err := c.HistStream(id, seq, limit, lt, gt, live, reverse, keys, values, private)
			if err != nil {
				return err
			}
			for b := range ch {
				fmt.Println(string(b))
			}
			return err
		},
	}
	cmd.Flags().BoolVar(
		&keys,
		"keys",
		false,
		"",
	)
	cmd.Flags().BoolVar(
		&values,
		"values",
		false,
		"",
	)
	cmd.Flags().BoolVar(
		&private,
		"private",
		false,
		"",
	)
	cmd.Flags().StringVar(
		&id,
		"id",
		"",
		"feed id",
	)
	cmd.Flags().BoolVar(
		&live,
		"live",
		false,
		"live stream",
	)
	cmd.Flags().BoolVar(
		&reverse,
		"reverse",
		false,
		"reverse stream",
	)
	cmd.Flags().Int64Var(
		&seq,
		"seq",
		0,
		"sequence number",
	)
	cmd.Flags().Int64Var(
		&lt,
		"lt",
		0,
		"less than",
	)
	cmd.Flags().Int64Var(
		&gt,
		"gt",
		0,
		"greater than",
	)
	cmd.Flags().Int64Var(
		&limit,
		"limit",
		-1,
		"max message count",
	)
	return cmd
}
