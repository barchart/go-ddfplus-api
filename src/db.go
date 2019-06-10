// Go ddfplus API
//
// Copyright 2019 Barchart.com, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the GNU license
// available at https://github.com/barchart/go-ddfplus-api/blob/master/LICENSE.
package ddf

import (
	"fmt"
	"log"
	"time"
)

type Quote struct {
	Symbol string `json:"symbol"`
	Info   struct {
		Name          string  `json:"name"`
		Exchange      string  `json:"exchange"`
		DDFExchange   string  `json:"ddfexchange"`
		BaseCode      string  `json:"basecode"`
		TickIncrement int     `json:"tickincrement"`
		PointValue    float64 `json:"pointvalue"`
	} `json:":info"`
	Data struct {
		CurrentSession struct {
			Bid       float64   `json:"bid"`
			BidSize   int64     `json:"bidsize"`
			Ask       float64   `json:"ask"`
			AskSize   int64     `json:"asksize"`
			Open      float64   `json:"open"`
			High      float64   `json:"high"`
			Low       float64   `json:"low"`
			Last      float64   `json:"last"`
			LastSize  int64     `json:"lastsize"`
			TradeTime time.Time `json:"tradetime"`
			Timestamp time.Time `json:timestamp"`
		} `json:"current"`
	} `json:"data"`
	LastUpdate time.Time `json:"lastupdate"`
}

type DB struct {
	data                map[string]*Quote
	listeners           map[string][]chan *Quote
	timestamp           time.Time
	marketUpdateChannel chan Message
}

func (db *DB) Connect(conn *Connection) {
	if db.marketUpdateChannel != nil {
		return
	}

	ch1 := make(chan MessageTimestamp)
	go func() {
		for m := range ch1 {
			db.timestamp = m.Timestamp
		}
	}()
	conn.RegisterTimestamp(ch1)

	db.marketUpdateChannel = make(chan Message)

	go func() {
		for m := range db.marketUpdateChannel {
			err := db.Process(m)
			if err != nil {
				log.Printf("Error processing message. %v", err)
				log.Println(m)
			}
		}
	}()
	conn.RegisterMarketUpdateAll(db.marketUpdateChannel)
}

func (db DB) GetQuote(symbol string) *Quote {
	return db.data[symbol]
}

func (db *DB) Process(m Message) error {
	switch m.Type() {
	case BidAsk:
		ba := m.(MessageBidAsk)
		q := db.data[ba.Symbol]
		if q == nil {
			return nil
		}

		q.Data.CurrentSession.Bid = ba.Bid
		q.Data.CurrentSession.BidSize = ba.BidSize
		q.Data.CurrentSession.Ask = ba.Ask
		q.Data.CurrentSession.AskSize = ba.AskSize
		q.Data.CurrentSession.Timestamp = ba.Timestamp

	case Refresh:
		rf := m.(MessageRefresh)
		q := db.data[rf.Symbol]
		if q == nil {
			q = &Quote{}
			q.Symbol = rf.Symbol
		}
		q.Info.Name = rf.Name
		q.Info.BaseCode = rf.BaseCode
		q.Info.Exchange = rf.Exchange
		q.Info.DDFExchange = rf.DDFExchange
		q.Info.TickIncrement = rf.TickIncrement
		q.Info.PointValue = rf.PointValue
		q.Data.CurrentSession.Open = rf.CurrentSession.Open
		q.Data.CurrentSession.High = rf.CurrentSession.High
		q.Data.CurrentSession.Low = rf.CurrentSession.Low
		q.Data.CurrentSession.Last = rf.CurrentSession.Last
		q.Data.CurrentSession.Ask = rf.Ask
		q.Data.CurrentSession.AskSize = rf.AskSize
		q.Data.CurrentSession.Bid = rf.Bid
		q.Data.CurrentSession.BidSize = rf.BidSize
		q.LastUpdate = rf.LastUpdate
		db.data[q.Symbol] = q

	case Trade:
		tr := m.(MessageTrade)
		q := db.data[tr.Symbol]
		if q == nil {
			return nil
		}

		q.Data.CurrentSession.Last = tr.Trade
		q.Data.CurrentSession.LastSize = tr.TradeSize
		q.Data.CurrentSession.TradeTime = tr.Timestamp
		q.Data.CurrentSession.Timestamp = q.Data.CurrentSession.TradeTime

	default:
		return fmt.Errorf("unhandled type %v", m.Type())
	}
	fmt.Println("DB Msg", m)
	return nil
}

func (db *DB) Register(symbols []string, ch chan *Quote) {
	for _, s := range symbols {
		channels := db.listeners[s]
		if channels == nil {
			channels = make([]chan *Quote, 0)
		}

		add := true
		for i := range channels {
			if channels[i] == ch {
				add = false
				break
			}
		}

		if add {
			channels = append(channels, ch)
			db.listeners[s] = channels
		}
	}
}

func (db *DB) Timestamp() time.Time {
	return db.timestamp
}

func InitDB() *DB {
	var db DB
	db.data = make(map[string]*Quote)
	db.listeners = make(map[string][]chan *Quote)

	return &db
}
