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
	var seq, limit, lt, gt int64
	var live, reverse, keys, values, private, json bool
	cmd := &cobra.Command{
		Use:  "hist FEED_ID [--limit n] [--seq n] [--live] [--keys] [--values] [--private]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			conf, err := opts.SSBConfig()
			if err != nil {
				return err
			}
			c, err := conf.Client(cmd.Context())
			if err != nil {
				return err
			}
			ch, err := c.HistStream(args[0], seq, limit, lt, gt, live, reverse, keys, values, private, json)
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
	cmd.Flags().BoolVar(
		&live,
		"live",
		false,
		"Keep the stream open and emit new messages as they are received.",
	)
	cmd.Flags().Int64Var(
		&seq,
		"seq",
		0,
		"(default: 0): If seq > 0, then only stream messages with sequence numbers greater than or equal to `seq`.",
	)
	cmd.Flags().Int64Var(
		&limit,
		"limit",
		-1,
		"max message count",
	)
	return cmd
}
