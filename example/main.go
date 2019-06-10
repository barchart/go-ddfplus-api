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

func main() {
	var user = flag.String("u", "", "Username")
	var pass = flag.String("p", "", "password")
	flag.Parse()

	var symbols = []string{"^EURUSD"}

	conn, err := ddf.NewConnection(&ddf.Credentials{
		Username: *user,
		Password: *pass,
	})

	if err != nil {
		log.Printf("Error creating connection. %v\n", err)
	} else {

		db := ddf.InitDB()
		db.Connect(conn)

		ch := make(chan ddf.MessageTimestamp)
		go func() {
			for m := range ch {
				fmt.Println(">>TS", m)
			}
		}()
		conn.RegisterTimestamp(ch)

		// Connection needs market update registration, but
		// we null op it here, since we listen for the processed
		// Quote message
		chMU := make(chan ddf.Message)
		go func() {
			for _ = range chMU {

			}
		}()
		conn.RegisterMarketUpdate(symbols, chMU)

		chq := make(chan *ddf.Quote)
		go func() {
			for q := range chq {
				fmt.Println(q)
			}
		}()
		db.Register(symbols, chq)

		conn.Start()
	}
}
