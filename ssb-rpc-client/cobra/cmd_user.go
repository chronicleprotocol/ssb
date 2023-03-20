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

// User returns a cobra command that prints the history of a feed.
func User(opts *Options) *cobra.Command {
	var limit, lt, gt int64
	var live, reverse, keys, values, private bool
	cmd := &cobra.Command{
		Use:  "user FEED_ID [--live] [--gte ts] [--lte ts] [--reverse] [--keys] [--limit n]",
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
			ch, err := c.UserStream(args[0], limit, lt, gt, live, reverse, keys, values, private)
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
		"Returns a different format than the default.",
	)
	cmd.Flags().BoolVar(
		&values,
		"values",
		false,
		"DOES NOT SEEM TO BE DOING ANYTHING IN THIS COMMAND.",
	)
	cmd.Flags().BoolVar(
		&private,
		"private",
		false,
		"DOES NOT SEEM TO BE DOING ANYTHING IN THIS COMMAND.",
	)
	cmd.Flags().BoolVar(
		&live,
		"live",
		false,
		"Keep the stream open and emit new messages as they are received.",
	)
	cmd.Flags().BoolVar(
		&reverse,
		"reverse",
		false,
		"Reverse stream output. Beware that due to the way LevelDB works, a reverse seek will be slower than a forward seek.",
	)
	cmd.Flags().Int64Var(
		&lt,
		"lt",
		0,
		"Timestamp is less than. When `--reverse` the order will be reversed, but the records streamed will be the same.",
	)
	cmd.Flags().Int64Var(
		&gt,
		"gt",
		0,
		"Timestamp is greater than. When `--reverse` the order will be reversed, but the records streamed will be the same.",
	)
	cmd.Flags().Int64Var(
		&limit,
		"limit",
		-1,
		"This number represents a maximum number of results and may not be reached if you get to the end of the data first. A value of -1 means there is no limit. When `--reverse` the highest keys will be returned instead of the lowest keys.",
	)
	return cmd
}
