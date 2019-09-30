/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/kcollasarundell/grpc-demo-me/step3/helloworld"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

type user struct {
	name         string
	betterColour string
}

var users = []user{
	user{name: "a", betterColour: "#0000FF"},
	user{name: "b", betterColour: "#008000"},
	user{name: "c", betterColour: "#00FF00"},
	user{name: "d", betterColour: "#FFA500"},
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancel()
	}()

	rl := rate.NewLimiter(rate.Limit(2), 2)
	// Set up a connection to the server.
	for {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("did not connect: %v", err)
			conn.Close()

			time.Sleep(2 * time.Second)
			continue
		}

		c := pb.NewGreeterClient(conn)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				ctx, cancel := context.WithTimeout(ctx, time.Second)
				err := rl.Wait(ctx)
				if err != nil {
					cancel()
					return
				}

				n := rand.Int63n(int64(len(users)))
				r, err := c.SayHello(ctx, &pb.HelloRequest{
					Name:         users[n].name,
					BetterColour: users[n].betterColour,
				})

				if err != nil {
					log.Printf("could not greet: %v", err)
					continue
				}

				log.Printf("%s, auth = %v ", r.Message, r.Authed)
				cancel()
			}
		}
	}

}
