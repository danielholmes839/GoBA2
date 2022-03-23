// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func dial(ctx context.Context, n int) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	url := fmt.Sprintf("ws://localhost:3000/game/connect?token=%d", n)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				err := conn.WriteMessage(websocket.TextMessage, []byte(t.String()+"\n"))
				if err != nil {
					log.Println("write:", err)
					return
				}
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s\n", message)
	}
}

func main() {
	ctx := context.Background()
	// go dial(ctx, 0)
	// for i := 0; i < 5; i++ {
	// 	go dial(ctx, i)
	// 	time.Sleep(time.Millisecond * 50)
	// }

	dial(ctx, 100)
	done := make(chan struct{})
	<-done

}
