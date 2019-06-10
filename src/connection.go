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
	"strings"
)

const (
	JerqVersion = 4
)

type Connection struct {
	connected               bool
	credentials             *Credentials
	settings                UserSettings
	marketDepthChannels     map[string][]chan Message
	marketUpdateChannels    map[string][]chan Message
	marketUpdateAllChannels []chan Message
	timestampChannels       []chan MessageTimestamp
}

func (c *Connection) connect() {

}

func (c *Connection) buildRequest() string {
	list := make(map[string]string)
	for s := range c.marketUpdateChannels {
		s2 := list[s]
		s2 += "Ss"
		list[s] = s2
	}

	for s := range c.marketDepthChannels {
		s2 := list[s]
		s2 += "Ss"
		list[s] = s2
	}

	sb := strings.Builder{}
	count := 0
	for k, v := range list {
		if count > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(k + "=" + v)
		count++
	}

	return sb.String()
}

func (c *Connection) Server() {
}

func (c *Connection) RegisterMarketUpdateAll(ch chan Message) {
	add := true
	for i, _ := range c.marketUpdateAllChannels {
		if c.marketUpdateAllChannels[i] == ch {
			fmt.Println("DON'T ADD", ch)

			add = false
			break
		}
	}

	if add {
		c.marketUpdateAllChannels = append(c.marketUpdateAllChannels, ch)
	}
}

func (c *Connection) RegisterMarketUpdate(symbols []string, ch chan Message) {
	newlist := make([]string, 0)

	for _, s := range symbols {
		channels := c.marketUpdateChannels[s]
		if channels == nil {
			channels = make([]chan Message, 0)
			newlist = append(newlist, s)
		}

		add := true
		for i, _ := range channels {
			if channels[i] == ch {
				add = false
				break
			}
		}

		if add {
			channels = append(channels, ch)
			c.marketUpdateChannels[s] = channels
		}

		if c.connected {
			// TO DO: Maybe we need to reset the connection
		}
	}
}

func (c *Connection) RegisterTimestamp(ch chan MessageTimestamp) {
	add := true
	for i, _ := range c.timestampChannels {
		if c.timestampChannels[i] == ch {
			add = false
			break
		}
	}

	if add {
		c.timestampChannels = append(c.timestampChannels, ch)
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

	fmt.Fprintf(conn, "LOGIN %s:%s\r\n", c.credentials.Username, c.credentials.Password)
	line, _, err := reader.ReadLine()
	if err != nil {
		log.Printf("Network error. %v", err)
		return
	}

	if line[0] != '+' {
		log.Printf("Error logging in. Server said: \"%s\"", string(line))
		return
	}

	fmt.Fprintf(conn, "VERSION %d\r\n", JerqVersion)
	line, _, err = reader.ReadLine()
	if err != nil {
		log.Printf("Network error. %v", err)
		return
	}

	command := c.buildRequest()

	fmt.Printf("Sending \"%s\" to server.\n", command)
	fmt.Fprintf(conn, "GO %s\r\n", command)
	for {
		line, _, _ := reader.ReadLine()
		m, err := Parse(line)
		if err == nil {
			if m != nil {
				switch m.Type() {
				case Timestamp:
					var ts MessageTimestamp
					ts = m.(MessageTimestamp)
					for _, ch := range c.timestampChannels {
						ch <- ts
					}
				case BidAsk, Refresh, Trade:
					var symbol string
					switch m.Type() {
					case BidAsk:
						symbol = m.(MessageBidAsk).Symbol
					case Refresh:
						symbol = m.(MessageRefresh).Symbol
					case Trade:
						symbol = m.(MessageTrade).Symbol
					}

					if symbol != "" {
						for _, ch := range c.marketUpdateAllChannels {
							ch <- m
						}

						for _, ch := range c.marketUpdateChannels[symbol] {
							ch <- m
						}
					}
				}
			}
		} else {
			log.Printf("Error parsing jerq data. %v", err)
		}
	}

	conn.Close()
}

func NewConnection(credentials *Credentials) (*Connection, error) {
	conn := &Connection{
		credentials: credentials,
	}

	conn.marketDepthChannels = make(map[string][]chan Message)
	conn.marketUpdateChannels = make(map[string][]chan Message)
	conn.marketUpdateAllChannels = make([]chan Message, 0)
	conn.timestampChannels = make([]chan MessageTimestamp, 0)

	settings, err := GetUserSettings(credentials)
	if err != nil {
		return nil, err
	}

	conn.settings = settings

	return conn, nil
}
