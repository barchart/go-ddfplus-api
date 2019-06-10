// Go ddfplus API
//
// Copyright 2019 Barchart.com, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the GNU license
// available at https://github.com/barchart/go-ddfplus-api/blob/master/LICENSE.
package ddf

import "time"

type MessageType int

const (
	BidAsk MessageType = iota
	Refresh
	Timestamp
	Trade
)

type DDFMessageInfo struct {
	BaseCode  string
	Exchange  string
	Delay     int
	Record    byte
	Subrecord byte
	DayCode   byte
	Session   byte
}

type Message interface {
	Type() MessageType
}

type MessageBidAsk struct {
	Symbol    string
	Info      DDFMessageInfo
	Bid       float64
	BidSize   int64
	Ask       float64
	AskSize   int64
	Timestamp time.Time
}

func (m MessageBidAsk) Type() MessageType {
	return BidAsk
}

type MessageRefresh struct {
	Symbol         string
	Name           string
	Exchange       string
	BaseCode       string
	PointValue     float64
	TickIncrement  int
	DDFExchange    string
	LastUpdate     time.Time
	Bid            float64
	BidSize        int64
	Ask            float64
	AskSize        int64
	CurrentSession struct {
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
	PreviousSession struct {
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
}

func (m MessageRefresh) Type() MessageType {
	return Refresh
}

type MessageTimestamp struct {
	Timestamp time.Time
}

func (m MessageTimestamp) Type() MessageType {
	return Timestamp
}

type MessageTrade struct {
	Symbol    string
	Info      DDFMessageInfo
	Trade     float64
	TradeSize int64
	Timestamp time.Time
}

func (m MessageTrade) Type() MessageType {
	return Trade
}
