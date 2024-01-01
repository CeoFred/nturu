package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/CeoFred/nturu-bun-grpc/shared/grpc"
)

var (
	port = flag.String("port", ":50052", "The server port")
)

var (
	addr    = flag.String("addr", "localhost:50051", "the address to connect to")
	user_id = flag.String("user_id", "1", "User ID")
)

func main() {

	// Setup HTTP server
	router := bunrouter.New(
		bunrouter.Use(reqlog.NewMiddleware()),
	)

	router.GET("/", indexHandler)

	router.WithGroup("/api", func(g *bunrouter.Group) {
		g.GET("/users/:id", debugHandler)
		g.GET("/users/current", debugHandler)
		g.GET("/users/*path", debugHandler)
	})

	// Run both servers concurrently
	log.Println("user service listening on http://localhost" + *port)
	log.Fatal(http.ListenAndServe(*port, router))

}

func indexHandler(w http.ResponseWriter, req bunrouter.Request) error {

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewUserClient(conn)
	ctx := context.Background()
	pr, err := c.GetProfile(ctx, &pb.UserID{ID: *user_id})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return bunrouter.JSON(w, bunrouter.H{
		"first_name": pr.FirstName,
		"last_name":  pr.LastName,
	})
}

func debugHandler(w http.ResponseWriter, req bunrouter.Request) error {
	return bunrouter.JSON(w, bunrouter.H{
		"route":  req.Route(),
		"params": req.Params().Map(),
	})
}
