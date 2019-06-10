// Go ddfplus API
//
// Copyright 2019 Barchart.com, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the GNU license
// available at https://github.com/barchart/go-ddfplus-api/blob/master/LICENSE.
package ddf

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"time"
)

func ParseTimestamp(ba []byte, etxpos int) (time.Time, error) {
	var (
		t time.Time
	)

	st := etxpos + 1
	xlen := len(ba) - st
	if xlen != 7 && xlen != 9 {
		return t, fmt.Errorf("invalid time length %d", len(ba)-st)
	}

	if ba[etxpos] != 3 {
		return t, fmt.Errorf("invalid char when <ETX> expected. %d", ba[etxpos])
	}

	if ba[st] != 20 {
		return t, fmt.Errorf("invalid century %d. Should be 20", ba[st])
	}
	year := int(ba[st])*100 + int(ba[st+1]-64)
	month := int(ba[st+2] - 64)
	date := int(ba[st+3] - 64)
	hour := int(ba[st+4] - 64)
	minute := int(ba[st+5] - 64)
	second := int(ba[st+6] - 64)
	ms := 0
	if xlen == 9 {
		ms = int((0xFF & ba[st+7])) + ((0xFF & int(ba[st+8])) << 8)
	}
	t = time.Date(year, time.Month(month), date, hour, minute, second, ms*1000, time.UTC)

	return t, nil
}

func Parse(ba []byte) (Message, error) {
	if len(ba) == 0 {
		return nil, nil
	}

	if len(ba) < 10 {
		log.Printf("Malformed message? %s\n", string(ba))
		log.Println(ba)
		return nil, nil
	}

	switch ba[0] {
	case 1: // DDF message
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
			info.BaseCode = string(ba[pos+3])
			info.Exchange = string(ba[pos+4])
			info.Delay, _ = strconv.Atoi(string(ba[pos+5 : pos+7]))

			switch ba[pos+1] {
			case '7': // Trades
				m := MessageTrade{}
				m.Symbol = sym

				pos += 7
				i := bytes.IndexByte(ba[pos:], ',')
				if i == -1 {
					return nil, fmt.Errorf("missing comma for trade")
				}
				m.Trade, _ = ParseFloat(string(ba[pos:pos+i]), info.BaseCode)

				pos += i + 1
				i = bytes.IndexByte(ba[pos:], ',')
				if i == -1 {
					return nil, fmt.Errorf("missing comma for trade size")
				}
				m.TradeSize, _ = strconv.ParseInt(string(ba[pos:pos+i]), 10, 64)

				pos += i + 1

				info.DayCode = ba[pos]
				info.Session = ba[pos+1]

				m.Timestamp, _ = ParseTimestamp(ba, pos+2)

				m.Info = info
				return m, nil

			case '8': // Bid/Ask
				m := MessageBidAsk{}
				m.Symbol = sym

				pos += 7
				i := bytes.IndexByte(ba[pos:], ',')
				if i == -1 {
					return nil, fmt.Errorf("missing comma for bid")
				}
				m.Bid, _ = ParseFloat(string(ba[pos:pos+i]), info.BaseCode)

				pos += i + 1
				i = bytes.IndexByte(ba[pos:], ',')
				if i == -1 {
					return nil, fmt.Errorf("missing comma for bid size")
				}
				m.BidSize, _ = strconv.ParseInt(string(ba[pos:pos+i]), 10, 64)

				pos += i + 1
				i = bytes.IndexByte(ba[pos:], ',')
				if i == -1 {
					return nil, fmt.Errorf("missing comma for ask")
				}
				m.Ask, _ = ParseFloat(string(ba[pos:pos+i]), info.BaseCode)

				pos += i + 1
				i = bytes.IndexByte(ba[pos:], ',')
				if i == -1 {
					return nil, fmt.Errorf("missing comma for ask size")
				}
				m.AskSize, _ = strconv.ParseInt(string(ba[pos:pos+i]), 10, 64)
				pos += i + 1

				info.DayCode = ba[pos]
				info.Session = ba[pos+1]

				m.Timestamp, _ = ParseTimestamp(ba, pos+2)

				m.Info = info
				return m, nil
			}
		}
	case 37: // '%' Refresh Message
		if ba[1] == '<' {
			type XMLSession struct {
				XMLName      xml.Name `xml:"SESSION"`
				ID           string   `xml:"id,attr"`
				Day          string   `xml:"day,attr"`
				Session      string   `xml:"session,attr"`
				Timestamp    string   `xml:"timestamp,attr"`
				Open         string   `xml:"open,attr"`
				High         string   `xml:"high,attr"`
				Low          string   `xml:"low,attr"`
				Last         string   `xml:"last,attr"`
				Previous     string   `xml:"previous,attr"`
				Settlement   string   `xml:"settlement,attr"`
				TradeSize    string   `xml:"tradesize,attr"`
				Volume       string   `xml:"volume,attr"`
				OpenInterest string   `xml:"openinterest,attr"`
				NumTrades    string   `xml:"numtrades,attr"`
				PriceVolume  string   `xml:"pricevolume,attr"`
				TradeTime    string   `xml:"tradetime,attr"`
				Ticks        string   `xml:"ticks,attr"`
			}

			type XMLQuote struct {
				XMLName       xml.Name     `xml:"QUOTE"`
				Sessions      []XMLSession `xml:"SESSION"`
				Symbol        string       `xml:"symbol,attr"`
				Name          string       `xml:"name,attr"`
				Exchange      string       `xml:"exchange,attr"`
				BaseCode      string       `xml:"basecode,attr"`
				PointValue    string       `xml:"pointvalue,attr"`
				TickIncrement string       `xml:"tickincrement,attr"`
				DDFExchange   string       `xml:"ddfexchange,attr"`
				Flag          string       `xml:"flag,attr"`
				LastUpdate    string       `xml:"lastupdate,attr"`
				Bid           string       `xml:"bid,attr"`
				BidSize       string       `xml:"bidsize,attr"`
				Ask           string       `xml:"ask,attr"`
				AskSize       string       `xml:"asksize,attr"`
				Mode          string       `xml:"mode,attr"`
			}

			var q XMLQuote
			err := xml.Unmarshal(ba[1:], &q)
			if err != nil {
				return nil, err
			}

			m := MessageRefresh{}
			m.Symbol = q.Symbol
			m.Name = q.Name
			m.Exchange = q.Exchange
			m.DDFExchange = q.DDFExchange
			m.BaseCode = q.BaseCode
			m.TickIncrement, err = strconv.Atoi(q.TickIncrement)
			m.PointValue, err = strconv.ParseFloat(q.PointValue, 32)
			m.LastUpdate, err = time.Parse("20060102150405", q.LastUpdate)

			m.Bid, err = ParseFloat(q.Bid, q.BaseCode)
			m.BidSize, err = strconv.ParseInt(q.BidSize, 10, 64)
			m.Ask, err = ParseFloat(q.Ask, q.BaseCode)
			m.AskSize, err = strconv.ParseInt(q.AskSize, 10, 64)

			for _, session := range q.Sessions {
				var ptr *struct {
					Day          string
					Session      string
					Timestamp    time.Time
					Open         float64
					High         float64
					Low          float64
					Last         float64
					Previous     float64
					Settlement   float64
					TradeSize    int64
					Volume       int64
					OpenInterest int64
					NumTrades    int64
					PriceVolume  float64
					TradeTime    time.Time
					Ticks        string
				}

				switch session.ID {
				case "combined":
					ptr = &m.CurrentSession
				case "previous":
					ptr = &m.PreviousSession
				default:
					continue
				}

				ptr.Day = session.Day
				ptr.Session = session.Session
				ptr.Timestamp, err = time.Parse("20060102150405", session.Timestamp)
				ptr.Open, err = ParseFloat(session.Open, q.BaseCode)
				ptr.High, err = ParseFloat(session.High, q.BaseCode)
				ptr.Low, err = ParseFloat(session.Low, q.BaseCode)
				ptr.Last, err = ParseFloat(session.Last, q.BaseCode)
				ptr.Previous, err = ParseFloat(session.Previous, q.BaseCode)
				ptr.Settlement, err = ParseFloat(session.Settlement, q.BaseCode)
				ptr.TradeSize, err = strconv.ParseInt(session.TradeSize, 10, 64)
				ptr.Volume, err = strconv.ParseInt(session.Volume, 10, 64)
				ptr.OpenInterest, err = strconv.ParseInt(session.OpenInterest, 10, 64)
				ptr.NumTrades, err = strconv.ParseInt(session.NumTrades, 10, 64)
				ptr.PriceVolume, err = strconv.ParseFloat(session.OpenInterest, 64)
				ptr.TradeTime, err = time.Parse("20060102150405", session.TradeTime)
				ptr.Ticks = session.Ticks
			}

			return m, nil

		} else {
			return nil, fmt.Errorf("unsupported refresh message")
		}
	default:
		return nil, fmt.Errorf("malformed message: %d \"%s\"", ba[0], string(ba))
	}

	return nil, nil
}
