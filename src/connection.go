// Go ddfplus API
//
// Copyright 2019 Barchart.com, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the GNU license
// available at https://github.com/barchart/go-ddfplus-api/blob/master/LICENSE.
package ddf

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

type Connection struct {
	credentials        *Credentials
	settings           UserSettings
	symbolListeners    map[string][]ConnectionListener
	timestampListeners []ConnectionListener
}

type ConnectionListener interface {
	OnBidAsk(MessageBidAsk)
	OnTimestamp(MessageTimestamp)
	OnTrade(MessageTrade)
}

func (this *Connection) connect() {

}

func (c *Connection) Server() {
}

func (c *Connection) AddListener(symbols []string, l ConnectionListener) {
	c.timestampListeners = append(c.timestampListeners, l)

	for _, s := range symbols {
		lst := c.symbolListeners[s]
		if lst == nil {
			lst = make([]ConnectionListener, 0)
		}
		lst = append(lst, l)
		c.symbolListeners[s] = lst
	}
}

func (c *Connection) Start() {
	// Dial the tcp
	conn, err := net.Dial("tcp", "qs01.ddfplus.com:7500")
	if err != nil {
		log.Printf("Error connecting. %v", err)
		return
	}

	reader := bufio.NewReader(conn)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		s := string(line)

		if s[0] == '+' {
			break
		}
	}

	fmt.Fprintf(conn, "LOGIN %s:%s VERSION 4\r\n", c.credentials.Username, c.credentials.Password)
	reader.ReadLine()

	fmt.Fprintf(conn, "GO YMM9=Ss,ESM9=Ss\r\n")
	for {
		line, _, _ := reader.ReadLine()
		m, err := Parse(line)
		if err == nil {
			if m != nil {
				switch m.Type() {
				case BidAsk:
					var ba MessageBidAsk
					ba = m.(MessageBidAsk)
					for _, l := range c.symbolListeners[ba.Symbol] {
						l.OnBidAsk(ba)
					}
				case Timestamp:
					var ts MessageTimestamp
					ts = m.(MessageTimestamp)
					for _, l := range c.timestampListeners {
						l.OnTimestamp(ts)
					}
				case Trade:
					var tr MessageTrade
					tr = m.(MessageTrade)
					for _, l := range c.symbolListeners[tr.Symbol] {
						l.OnTrade(tr)
					}
				}
			}
		}
	}

	conn.Close()
}

func NewConnection(credentials *Credentials) (*Connection, error) {
	conn := &Connection{
		credentials: credentials,
	}

	conn.timestampListeners = make([]ConnectionListener, 0)
	conn.symbolListeners = make(map[string][]ConnectionListener)

	settings, err := GetUserSettings(credentials)
	if err != nil {
		return nil, err
	}

	conn.settings = settings

	return conn, nil
}
