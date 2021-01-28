package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"time"

	socketio "github.com/gaconkzk/socket.io-client-go"
	"github.com/gaconkzk/socket.io-client-go/websocket"
)

// // Route to fly.
// type Route struct {
// 	To   string
// 	From string
// }

// // HotelReservation of a room at a hotel nearby the airport.
// type HotelReservation struct {
// 	Name     string
// 	Location string
// 	Room     string
// 	Price    int
// }

// // Airports clique.
// var Airports = []string{"JFK", "KEF", "ATL", "MIA", "DAO", "FCO"}

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	var u = url.URL{
		Scheme: "wss",
		Host:   "manage.vw3.cc:8443",
	}

	var query = u.Query()
	query.Add("refresh_token", "8502bf8db18b8bfe1b95a03d6128309891275575")
	u.RawQuery = query.Encode()

	c := socketio.NewClient(u, websocket.NewTransport())
	c.On(socketio.OnConnection, func() {
		fmt.Print("Connected")
	})

	if err := c.On(socketio.OnError, errorHandler); err != nil {
		log.Fatalf("error %v", err)
		panic(err)
	}

	if err := c.On("connect_error", errorHandler); err != nil {
		log.Fatalf("error %v", err)
		panic(err)
	}

	err := c.Connect()
	if err != nil {
		log.Fatalf("error, %v", err)
		// panic(err) // you should prefer returning errors than panicking
	}

	nsp, err := c.Of("accountant")
	if err := nsp.On(socketio.OnError, errorHandler); err != nil {
		log.Fatalf("error, namespace %v", err)
		panic(err)
	}

	if err := nsp.On(socketio.OnDisconnect, disconnectHandler); err != nil {
		panic(err)
	}

	// if err := nsp.On("flight", flightHandler); err != nil {
	// 	panic(err)
	// }

	if err := nsp.On("ready", readyHandler); err != nil {
		panic(err)
	}

	if err := nsp.On("skip", skipHandler); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Second)

	g := goodbye{c, cancel}

	if err := nsp.On("goodbye", g.Handler); err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			// doSomething(c)
		}
	}
}

// func doSomething(c *socketio.Client) {
// 	index1 := rand.Intn(len(Airports))
// 	index2 := rand.Intn(len(Airports))

// 	if index1 == index2 {
// 		bookHotelRoom(c, Airports[index1])
// 	}

// 	if err := c.Emit("find_tickets", Route{
// 		From: Airports[index1],
// 		To:   Airports[index2],
// 	}); err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
// }

// func bookHotelRoom(c *socketio.Client, hotel string) {
// 	var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	var res HotelReservation
// 	switch err := c.Ack(ctx, "book_hotel_for_tonight", hotel, &res); {
// 	case err == nil:
// 		fmt.Printf("%s has booked a %s bedroom for $%d near %s.\n", res.Name, res.Room, res.Price, res.Location)
// 	case err == context.DeadlineExceeded || err == context.Canceled:
// 		fmt.Fprintf(os.Stderr, "Couldn't complete a booking at %s.\n", hotel)
// 	case err != nil:
// 		fmt.Fprintln(os.Stderr, err)
// 		os.Exit(1)
// 	}
// }

func errorHandler(err error) {
	fmt.Fprintf(os.Stderr, "error received: %v\n", err.Error())
	os.Exit(1)
}

func disconnectHandler() {
	fmt.Println("Disconnecting.")
	os.Exit(0)
}

// func flightHandler(vehicle string, route Route) {
// 	fmt.Printf("The %s is flying from %s to %s.\n", vehicle, route.From, route.To)
// }

func skipHandler(vehicle string) {
	fmt.Printf("The %s is not in use.\n", vehicle)
}

func readyHandler() {
	log.Print("should ok???.\n")
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
