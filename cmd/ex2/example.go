package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	socketio "github.com/gaconkzk/socket.io-client-go"
	"github.com/gaconkzk/socket.io-client-go/websocket"
)

func main() {
	var u = url.URL{
		Scheme: "wss",
		Host:   "manage.vw3.cc:8443",
	}

	var query = u.Query()
	query.Add("refresh_token", "96078fa8181831718b404d70c0f3d78d76fae8cc")
	u.RawQuery = query.Encode()

	co := socketio.NewClient(u, websocket.NewTransport())
	c, err := co.Of("")
	if err != nil {
		log.Fatalf("error, %v", err)
	}

	nsp, err := co.Of("accountant")
	if err != nil {
		log.Fatalf("error, namespace %v", err)
		panic(err)
	}
	if err := nsp.On(socketio.OnConnection, connectionHandler); err != nil {
		log.Fatalf("error, namespace %v", err)
		panic(err)
	}
	if err := nsp.On(socketio.OnError, errorHandler); err != nil {
		log.Fatalf("error, namespace %v", err)
		panic(err)
	}

	if err := nsp.On(socketio.OnDisconnect, disconnectHandler); err != nil {
		panic(err)
	}

	if err := nsp.On("ready", readyHandler); err != nil {
		panic(err)
	}

	err = co.Connect()
	if err != nil {
		log.Fatalf("error, %v", err)
		// panic(err) // you should prefer returning errors than panicking
	}
	// on methods - for default namespace should after connect
	co.On(socketio.OnConnection, connectionHandler)

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Second)

	g := goodbye{co, cancel}

	if err := c.On("goodbye", g.Handler); err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			// doSomething(c)
		case <-c.Ready():
			log.Print("Ready for ns connect")
			co.NamespaceConnect("accountant")
		}
	}
}

type goodbye struct {
	client *socketio.Client
	cancel context.CancelFunc
}

func (g *goodbye) Handler() {
	fmt.Print(`Oops! This program is exiting in 5s to demonstrate a clean termination approach.
Comment the "goodbye" event listener in the Go code example to avoid this from happening.
The server sends this "goodbye" message 120 seconds after the connection has been established.
`)
	time.Sleep(5 * time.Second)
	g.cancel()

	err := g.client.Close()

	if err != nil {
		panic(err)
	}
}

func errorHandler(err error) {
	fmt.Fprintf(os.Stderr, "error received: %v\n", err.Error())
	os.Exit(1)
}

func disconnectHandler() {
	fmt.Println("Disconnecting.")
	os.Exit(0)
}

func skipHandler(vehicle string) {
	fmt.Printf("The %s is not in use.\n", vehicle)
}

func connectionHandler() {
	log.Print("connected")
}

func readyHandler() {
	log.Print("should ok???.\n")
}
