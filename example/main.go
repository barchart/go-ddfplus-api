// Go ddfplus API Examples
//
// Copyright 2019 Barchart.com, Inc. All rights reserved.
//
// This Source Code Form is subject to the terms of the GNU license
// available at https://github.com/barchart/go-ddfplus-api/blob/master/LICENSE.

package main

import (
	ddf "barchart/go-ddfpus-api/src"
	"flag"
	"fmt"
	"log"
)

type Listener struct{}

func (l Listener) OnBidAsk(m ddf.MessageBidAsk) {
	fmt.Println("BA!!!", m)
}

func (l Listener) OnTimestamp(m ddf.MessageTimestamp) {
	fmt.Println("TS!!!", m)
}

func (l Listener) OnTrade(m ddf.MessageTrade) {
	fmt.Println("TR!!!", m)
}

func main() {
	var user = flag.String("u", "", "Username")
	var pass = flag.String("p", "", "password")
	flag.Parse()

	conn, err := ddf.NewConnection(&ddf.Credentials{
		Username: *user,
		Password: *pass,
	})

	if err != nil {
		log.Printf("Error creating connection. %v\n", err)
	} else {
		var l Listener

		conn.AddListener([]string{"ESM9"}, &l)
		conn.Start()
	}
}
