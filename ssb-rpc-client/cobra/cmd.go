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
	"net"

	"github.com/spf13/cobra"
	"github.com/ssbc/go-netwrap"
	"github.com/ssbc/go-secretstream"
	ssbServer "github.com/ssbc/go-ssb"
	"github.com/ssbc/go-ssb/invite"

	"github.com/chronicleprotocol/ssb/ssb-rpc-client/config"

	"github.com/chronicleprotocol/ssb/ssb-rpc-client/client"
)

type Options struct {
	ConfigPath string
	SecretPath string
	SsbHost    string
	SsbPort    int
	Verbose    bool
}

func (opts *Options) SSBConfig() (*client.Config, error) {
	keys, err := ssbServer.LoadKeyPair(opts.SecretPath)
	if err != nil {
		return nil, err
	}
	caps, err := config.LoadCapsFile(opts.ConfigPath)
	if err != nil {
		return nil, err
	}
	if caps.Shs == "" || caps.Sign == "" {
		caps, err = config.LoadCapsFromConfigFile(opts.ConfigPath)
		if err != nil {
			return nil, err
		}
	}
	if caps.Invite != "" {
		inv, err := invite.ParseLegacyToken(caps.Invite)
		if err != nil {
			return nil, err
		}
		return &client.Config{
			Keys: keys,
			Shs:  caps.Shs,
			Addr: inv.Address,
		}, nil
	}
	ip := net.ParseIP(opts.SsbHost)
	if ip == nil {
		resolvedAddr, err := net.ResolveIPAddr("ip", opts.SsbHost)
		if err != nil {
			return nil, err
		}
		ip = resolvedAddr.IP
	}
	return &client.Config{
		Keys: keys,
		Shs:  caps.Shs,
		Addr: netwrap.WrapAddr(
			&net.TCPAddr{
				IP:   ip,
				Port: opts.SsbPort,
			},
			secretstream.Addr{PubKey: keys.ID().PubKey()},
		),
	}, nil
}

func Root() (*Options, *cobra.Command) {
	return &Options{}, &cobra.Command{
		Use: "ssb-rpc-client",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		DisableAutoGenTag: true,
		SilenceUsage:      true,
	}
}
