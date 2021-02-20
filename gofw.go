package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

// Our connections
type connection struct {
	client net.Conn
	target net.Conn
}

// src and dest ports
var (
	portFrom int
	portTo   int
)

// Parse arguments
func init() {
	flag.IntVar(&portFrom, "from", 0, "Port to forward tcp traffic from")
	flag.IntVar(&portTo, "to", 0, "Port to forward tcp traffic to")
}

func main() {
	// Simple greeting message
	fmt.Println("GO! Forward 1.0")

	// Parse the input flags
	flag.Parse()

	// Check if ports are specified
	if portFrom == 0 || portTo == 0 {
		println("Error: Both ports must be specified")
		os.Exit(1)
	}

	// Create a new signals channel
	signals := make(chan os.Signal)
	// When the Interrupt signal was noticed, it will be send to the 'signals' channel
	signal.Notify(signals, os.Interrupt)

	// Create listening tcp server
	incomingAddress := "127.0.0.1:" + strconv.Itoa(portFrom)
	incoming, err := net.Listen("tcp", incomingAddress)
	if err != nil {
		println("Error: Could not bind on local port:", portFrom)
		os.Exit(1)
	}

	// Pass connections around
	connections := make(chan connection)

	// Receive quit from forwarder
	rQuit := make(chan string)

	// Send quit to forwarder
	sQuit := make(chan string)

	// Accept clients
	// Accepted client will be passed to the connections channel
	go acceptClient(incoming, connections)

	// Keep track of running tcp connections
	connNumber := 0

	// Endless control loop
	for {
		select {
		case conn := <-connections:
			// Receives connection from acceptClient() and sends it's information to the forwarder
			go forward(conn, rQuit, sQuit)
			connNumber++
		case <-rQuit:
			// We receive a quit signal if the forwarder ends
			connNumber--
		case <-signals:
			// We receive a interrupt signal from the os
			// Send quit signal to forwarder
			for i := 0; i < connNumber; i++ {
				sQuit <- "quit"
			}

			// Wait till forwarder quit
			for i := 0; i < connNumber; i++ {
				<-rQuit
			}

			os.Exit(0)
		}
	}
}

func acceptClient(incomingServer net.Listener, connections chan<- connection) {
	// Accept any number of incoming clients
	for {
		// Accept client who wants to connect
		client, err := incomingServer.Accept()
		if err != nil {
			println("Error: Could not accept client:", err.Error())
			continue
		}

		// Build tcp session to forward data to the target
		targetAddress := "127.0.0.1:" + strconv.Itoa(portTo)
		target, err := net.Dial("tcp", targetAddress)
		if err != nil {
			println("Error: Could not establish tcp connection to:", targetAddress)
			client.Close()
			continue
		}

		// Send client and target to the control loop in the main function
		connections <- connection{client, target}
	}
}

func forward(conn connection, sQuit chan<- string, rQuit <-chan string) {
	// Setup WaitGroup for the active tcp session
	var wg sync.WaitGroup
	waitChannel := make(chan struct{})
	wg.Add(2)

	// Wrapper function to enclose the WaitGroup
	go func() {
		// Forward from the target to the client
		go func() {
			defer wg.Done()
			io.Copy(conn.target, conn.client)
		}()

		// Forward from the client to the target
		go func() {
			defer wg.Done()
			io.Copy(conn.client, conn.target)
		}()

		// Wait for the connection.
		wg.Wait()

		// If the connection ended we close the wait channel. This will trigger
		// an action in the select statement below
		close(waitChannel)
	}()

	// Wait for a quit signal. Either from the WaitGroup or the main control loop
	select {
	case <-waitChannel:
		break
	case <-rQuit:
		break
	}

	// Close connections
	conn.client.Close()
	conn.target.Close()

	// Send a notification that one forwarder quits
	sQuit <- "quit"

}
