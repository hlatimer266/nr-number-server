package server

import (
	"context"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	maxClients    = 5
	port          = ":4000"
	validInput    = "114157263$\n"
	twoValidInput = "225317119$\n222222225$\n"
	allValid      = []string{validInput, twoValidInput}
)

func TestServer(t *testing.T) {

	// test 5 connections can be made
	t.Run("make 5 connections and write valid input", func(t *testing.T) {
		// ctx := context.Background()
		go Run(port, context.Background())
		time.Sleep(time.Second * 2)

		wg := sync.WaitGroup{}
		wg.Add(maxClients)
		for i := 0; i < maxClients; i++ {
			client, err := net.Dial("tcp", port)
			if err != nil {
				t.Fatalf("unable to connect to server - %s\n", err.Error())
			}
			defer client.Close()

			go func() {
				for _, input := range allValid {
					_, err := client.Write([]byte(input))
					if err != nil {
						break
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()
		os.Remove("number.log")
	})

}
