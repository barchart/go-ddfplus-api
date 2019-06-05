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
	Timestamp MessageType = iota
	Trade
	BidAsk
)

type DDFMessageInfo struct {
	BaseCode  byte
	Exchange  byte
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
	Symbol  string
	Info    DDFMessageInfo
	Bid     string
	BidSize string
	Ask     string
	AskSize string
}

func (m MessageBidAsk) Type() MessageType {
	return BidAsk
}

type MessageTimestamp struct {
	Timestamp time.Time
}

func (this MessageTimestamp) Type() MessageType {
	return Timestamp
}

type MessageTrade struct {
	Symbol    string
	Info      DDFMessageInfo
	Trade     string
	TradeSize string
}

func (m MessageTrade) Type() MessageType {
	return Trade
}
