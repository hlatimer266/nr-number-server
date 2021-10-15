package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

const (
	maxClients   = 5
	port         = ":4000"
	waitPeriod   = 10
	clientWrites = 2000000
	batchSize    = 1000
)

func inputGenerator(numbersPerClient, numbersPerWrite int) (inputsPerClient map[int][][]byte) {
	var unique = 1
	var numsPerWrite string
	inputsPerClient = make(map[int][][]byte, maxClients)

	for i := 0; i < maxClients; i++ {
		inputsClient := [][]byte{}
		for ; unique <= numbersPerClient; unique++ {
			numsPerWrite = fmt.Sprintf("%s%09d"+"$"+"\n", numsPerWrite, unique)

			if unique%numbersPerWrite == 0 || unique == numbersPerClient {
				inputsClient = append(inputsClient, []byte(numsPerWrite))
				numsPerWrite = ""
			}
		}
		inputsPerClient[i] = inputsClient
	}

	return
}

func main() {

	testPeriod := time.Duration(waitPeriod) * time.Second

	fmt.Printf("***** STEP 1 *****\nGenerate inputs to send\n")

	inputs := inputGenerator(clientWrites, batchSize)

	testStart := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(maxClients)

	fmt.Printf("***** STEP 2 *****\nConnect %v clients and write\n", maxClients)
	for _, clientInput := range inputs {
		go func(clientInput [][]byte) {
			client, _ := net.Dial("tcp", port)
			defer client.Close()

			// wait for test start
			wg.Done()
			<-testStart

			for _, input := range clientInput {
				client.Write([]byte(input))
			}
		}(clientInput)
	}

	wg.Wait()
	close(testStart)

	timer := time.NewTimer(testPeriod)
	<-timer.C
	fmt.Println("***** STEP 3 *****\nTest Complete")
}
