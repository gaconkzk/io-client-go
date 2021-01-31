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
	query.Add("refresh_token", "3f9f611a258a13de23c61e5c5fe832d6dc43f351")
	u.RawQuery = query.Encode()

	co := socketio.NewClient(u, websocket.NewTransport())
	c, err := co.Of("")

	// nc, err := c.Of("accountant")
	c.On(socketio.OnConnection, func() {
		log.Print("Connected")
		fmt.Println("Connected")
	})

	if err := c.On(socketio.OnError, errorHandler); err != nil {
		log.Fatalf("error %v", err)
		panic(err)
	}

	if err := c.On("connect_error", errorHandler); err != nil {
		log.Fatalf("error %v", err)
		panic(err)
	}

	if err := c.On("ready", readyHandler); err != nil {
		panic(err)
	}

	err = co.Connect()
	if err != nil {
		log.Fatalf("error, %v", err)
		// panic(err) // you should prefer returning errors than panicking
	}

	// nsp, err := c.Of("accountant")
	// if err := nsp.On(socketio.OnError, errorHandler); err != nil {
	// 	log.Fatalf("error, namespace %v", err)
	// 	panic(err)
	// }

	// if err := nsp.On(socketio.OnDisconnect, disconnectHandler); err != nil {
	// 	panic(err)
	// }

	// if err := nsp.On("flight", flightHandler); err != nil {
	// 	panic(err)
	// }

	// if err := nsp.On("ready", readyHandler); err != nil {
	// 	panic(err)
	// }

	// if err := nsp.On("skip", skipHandler); err != nil {
	// 	panic(err)
	// }

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
			fmt.Printf("??? ready?")
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

func readyHandler() {
	log.Print("should ok???.\n")
}
