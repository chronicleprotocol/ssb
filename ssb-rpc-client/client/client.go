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

package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ssbc/go-muxrpc/v2"
	"github.com/ssbc/go-ssb-refs"
	"github.com/ssbc/go-ssb/message"
	"github.com/ssbc/go-ssb/plugins/legacyinvites"
)

const methodPublish = "publish"
const methodWhoAmI = "whoami"
const methodCreateLogStream = "createLogStream"

// createLogStream
// log
// [--live]
// [--gt index]
// [--gte index]
// [--lt index]
// [--lte index]
// [--reverse]
// [--keys]
// [--values]
// [--limit n]

const methodCreateUserStream = "createUserStream"

// createUserStream
//
//	--id {feedid}
//	[--live]
//	[--gt index]
//	[--gte index]
//	[--lt index]
//	[--lte index]
//	[--reverse]
//	[--keys]
//	[--values]
//	[--limit n]
const methodCreateHistoryStream = "createHistoryStream"

const methodInvite = "invite"
const methodCreate = "create"
const methodAccept = "accept"

type Client struct {
	ctx context.Context
	rpc Endpoint
}

type Endpoint interface {
	Async(context.Context, interface{}, muxrpc.RequestEncoding, muxrpc.Method, ...interface{}) error
	Source(context.Context, muxrpc.RequestEncoding, muxrpc.Method, ...interface{}) (*muxrpc.ByteSource, error)
}

func (c *Client) WhoAmI() ([]byte, error) {
	var ret []byte
	return ret, c.rpc.Async(c.ctx, &ret, muxrpc.TypeBinary, muxrpc.Method{methodWhoAmI})
}

func (c *Client) InviteCreate(n uint) ([]byte, error) {
	var ret []byte
	var args legacyinvites.CreateArguments
	args.Uses = n
	return ret, c.rpc.Async(c.ctx, &ret, muxrpc.TypeBinary, muxrpc.Method{methodInvite, methodCreate}, args)
}

func (c *Client) InviteAccept(invite string) ([]byte, error) {
	var ret []byte
	var args = struct {
		Invite string `json:"invite"`
	}{invite}

	return ret, c.rpc.Async(c.ctx, &ret, muxrpc.TypeBinary, muxrpc.Method{methodInvite, methodAccept}, args)
}

func (c *Client) Transmit(v interface{}) ([]byte, error) {
	var ret []byte
	// TODO Add Rate Limiter
	return ret, c.rpc.Async(c.ctx, &ret, muxrpc.TypeBinary, muxrpc.Method{methodPublish}, v)
}

func (c *Client) ReceiveLast(id, contentType string, limit int64) ([]byte, error) {
	feedRef, err := refs.ParseFeedRef(id)
	if err != nil {
		return nil, err
	}
	ch, err := c.callSSB(methodCreateUserStream, message.CreateHistArgs{
		CommonArgs: message.CommonArgs{
			Keys: true,
		},
		StreamArgs: message.StreamArgs{
			Limit:   limit,
			Reverse: true,
		},
		ID: feedRef,
	}, false)
	if err != nil {
		return nil, err
	}
	var data struct {
		Value struct {
			Content FeedAssetPrice `json:"content"`
		} `json:"value"`
	}
	for b := range ch {
		if err = json.Unmarshal(b, &data); err != nil {
			return nil, err
		}
		if contentType == "" || data.Value.Content.Type == contentType {
			return b, nil
		}
	}
	if contentType != "" {
		return nil, fmt.Errorf("no data of type %s in the stream for ref: %s", contentType, feedRef.String())
	}
	return nil, fmt.Errorf("no data in the stream for ref: %s", feedRef.String())
}

func (c *Client) LogStream(seq, limit, lt, gt int64, live, reverse, keys, values, private bool) (chan []byte, error) {
	return c.callSSB(methodCreateLogStream, message.CreateLogArgs{
		CommonArgs: message.CommonArgs{
			Keys:    keys,
			Values:  values,
			Private: private,
			Live:    live,
		},
		StreamArgs: message.StreamArgs{
			Limit:   limit,
			Reverse: reverse,
			Lt:      message.RoundedInteger(lt),
			Gt:      message.RoundedInteger(gt),
		},
		Seq: seq,
	}, live)
}

func (c *Client) HistStream(id string, seq, limit, lt, gt int64, live, reverse, keys, values, private, json bool) (chan []byte, error) {
	feedRef, err := refs.ParseFeedRef(id)
	if err != nil {
		return nil, err
	}
	return c.callSSB(methodCreateHistoryStream, message.CreateHistArgs{
		CommonArgs: message.CommonArgs{
			Keys:    keys,
			Values:  values,
			Private: private,
			Live:    live,
		},
		StreamArgs: message.StreamArgs{
			Limit:   limit,
			Reverse: reverse,
			Lt:      message.RoundedInteger(lt),
			Gt:      message.RoundedInteger(gt),
		},
		ID:     feedRef,
		Seq:    seq,
		AsJSON: json,
	}, live)
}

func (c *Client) callSSB(method string, arg interface{}, live bool) (chan []byte, error) {
	var ctx context.Context
	var cancel context.CancelFunc
	if live {
		ctx, cancel = context.WithCancel(c.ctx)
	} else {
		ctx, cancel = context.WithTimeout(c.ctx, time.Second)
	}
	src, err := c.rpc.Source(ctx, muxrpc.TypeBinary, muxrpc.Method{method}, arg)
	if err != nil {
		cancel()
		return nil, err
	}
	ch := make(chan []byte)
	go func() {
		defer close(ch)
		defer cancel()
		defer func() {
			if r := recover(); r != nil {
				log.Println("recovered:", r)
			}
		}()
		for nxt := src.Next(ctx); nxt; nxt = src.Next(ctx) {
			b, err := src.Bytes()
			if err != nil {
				log.Println(err)
				return
			}
			ch <- b
		}
	}()
	return ch, nil
}
