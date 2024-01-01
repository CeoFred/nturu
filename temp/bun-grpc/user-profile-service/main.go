package main

import (
	"context"
	"flag"
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"google.golang.org/grpc"

	pb "github.com/CeoFred/nturu-bun-grpc/shared/grpc"
)

var (
	port = flag.String("port", ":50051", "The server port")
)

type server struct {
	pb.UserServer
}

func (s *server) GetProfile(ctx context.Context, in *pb.UserID) (*pb.UserProfile, error) {
	log.Printf("Received: %v", in.GetID())
	return &pb.UserProfile{
		FirstName: "Alfred",
		LastName:  "Johnson-Awah",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServer(s, &server{})

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
	go func() {
		log.Println("listening on http://localhost:9999")
		log.Fatal(http.ListenAndServe(":9999", router))
	}()

	log.Println("gRPC server listening on", *port)
	log.Fatal(s.Serve(lis))

}

func indexHandler(w http.ResponseWriter, req bunrouter.Request) error {
	return indexTemplate().Execute(w, nil)
}

func debugHandler(w http.ResponseWriter, req bunrouter.Request) error {
	return bunrouter.JSON(w, bunrouter.H{
		"route":  req.Route(),
		"params": req.Params().Map(),
	})
}

var indexTmpl = `
<html>
  <h1>Welcome</h1>
  <ul>
    <li><a href="/api/users/123">/api/users/123</a></li>
    <li><a href="/api/users/current">/api/users/current</a></li>
    <li><a href="/api/users/foo/bar">/api/users/foo/bar</a></li>
  </ul>
</html>
`

func indexTemplate() *template.Template {
	return template.Must(template.New("index").Parse(indexTmpl))
}
