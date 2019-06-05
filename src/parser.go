// Go ddfplus API
//
// Copyright 2019 Barchart.com, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the GNU license
// available at https://github.com/barchart/go-ddfplus-api/blob/master/LICENSE.
package ddf

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

func Parse(ba []byte) (Message, error) {
	if len(ba) < 10 || ba[0] != 1 {
		return nil, fmt.Errorf("malformed message: \"%s\".", string(ba))
	}

	switch ba[1] { // Record
	case '#': // Timestamp
		t, err := time.Parse("20060102150405", string(ba[2:16]))
		if err == nil {
			m := MessageTimestamp{}
			m.Timestamp = t
			return m, nil

		}
	case '2': //
		info := DDFMessageInfo{}
		info.Record = ba[1] // Record

		// Find the comma
		i := bytes.IndexByte(ba, ',')
		if i == -1 {
			return nil, fmt.Errorf("no comma in type 2")
		}

		pos := i

		sym := string(ba[2:pos])
		info.Subrecord = ba[pos+1]
		info.BaseCode = ba[pos+3]
		info.Exchange = ba[pos+4]
		info.Delay, _ = strconv.Atoi(string(ba[pos+5 : pos+7]))

		switch ba[pos+1] {
		case '7': // Trades
			m := MessageTrade{}
			m.Symbol = sym
			info.BaseCode = ba[pos+3]

			pos += 7
			i := bytes.IndexByte(ba[pos:], ',')
			if i == -1 {
				return nil, fmt.Errorf("missing comma for trade")
			}
			m.Trade = string(ba[pos : pos+i])

			pos += i + 1
			i = bytes.IndexByte(ba[pos:], ',')
			if i == -1 {
				return nil, fmt.Errorf("missing comma for trade size")
			}
			m.TradeSize = string(ba[pos : pos+i])

			pos += i + 1

			info.DayCode = ba[pos]
			info.Session = ba[pos+1]

			m.Info = info
			return m, nil

		case '8': // Bid/Ask
			m := MessageBidAsk{}
			m.Symbol = sym
			info.BaseCode = ba[pos+3]

			pos += 7
			i := bytes.IndexByte(ba[pos:], ',')
			if i == -1 {
				return nil, fmt.Errorf("missing comma for bid")
			}
			m.Bid = string(ba[pos : pos+i])

			pos += i + 1
			i = bytes.IndexByte(ba[pos:], ',')
			if i == -1 {
				return nil, fmt.Errorf("missing comma for bid size")
			}
			m.BidSize = string(ba[pos : pos+i])

			pos += i + 1
			i = bytes.IndexByte(ba[pos:], ',')
			if i == -1 {
				return nil, fmt.Errorf("missing comma for ask")
			}
			m.Ask = string(ba[pos : pos+i])

			pos += i + 1
			i = bytes.IndexByte(ba[pos:], ',')
			if i == -1 {
				return nil, fmt.Errorf("missing comma for ask size")
			}
			m.AskSize = string(ba[pos : pos+i])
			pos += i + 1

			info.DayCode = ba[pos]
			info.Session = ba[pos+1]

			m.Info = info
			return m, nil
		}
	}

	return nil, nil
}
