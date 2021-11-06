// Package main implements a server for message service.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/ccwutw/grpchub/message"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	certFile   string
	keyFile    string
	caCertFile string
	address    string
	insecure   bool
)

func init() {
	flag.StringVar(&address, "address", ":50051", "listen address")
	flag.StringVar(&certFile, "cert", "certificates/server.crt", "certificate file")
	flag.StringVar(&keyFile, "key", "certificates/server.key", "key file")
	flag.StringVar(&caCertFile, "cacert", "certificates/ca.crt", "ca certificate file")
	flag.BoolVar(&insecure, "insecure", false, "run insecure")
	flag.Parse()
}

// server is used to implement message.GreeterServer.
type greeterServer struct {
	message.UnimplementedGreeterServer
}

// echoServer is used to implement a message.EchoServer
type echoServer struct {
	message.UnimplementedEchoServer
}

// SayHello implements message.GreeterServer
func (s *greeterServer) SayHello(ctx context.Context, in *message.HelloRequest) (*message.HelloReply, error) {
	log.Printf("Received HelloRequest: %v", in.GetName())
	return &message.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// LotsOfReplices repeaths greeting 5 times
func (s *greeterServer) LotsOfReplies(in *message.HelloRequest, streamRes message.Greeter_LotsOfRepliesServer) error {
	log.Printf("Received LotsOfReplies for HelloRequest: %v", in.GetName())

	// Return greeting 5 times
	for i := 0; i < 5; i++ {
		log.Printf("--> sending reply %d for %v", i, in.GetName())
		res := message.HelloReply{Message: fmt.Sprintf("Hello %v - %d", in.GetName(), i)}
		streamRes.Send(&res)
	}
	return nil
}

// ManyHellos sends hello to many names
func (s *greeterServer) ManyHellos(in message.Greeter_ManyHellosServer) error {
	log.Printf("Received ManyHellos for HelloRequest")
	for {
		req, err := in.Recv()
		if err == io.EOF {
			// end of receiving requests
			log.Print("--> end receiving requests and send response")
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("--> sending reply for %v", req.GetName())
		in.Send(&message.HelloReply{Message: "Hello " + req.GetName()})
	}
}

// Replay impelments message.EchoServer
func (s *echoServer) Replay(ctx context.Context, in *message.EchoRequest) (*message.EchoReply, error) {
	log.Printf("Received EchoRequest: %v", in.GetMessage())
	return &message.EchoReply{Message: in.GetMessage()}, nil
}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var s *grpc.Server

	if insecure {
		log.Print("running insecure!")
		s = grpc.NewServer()
	} else {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatal(err)
		}

		// Create client CA pool
		caPool, err := message.NewCAPool(caCertFile)
		if err != nil {
			log.Fatal(err)
		}

		// Create server TLS config
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientCAs:    caPool,
			ClientAuth:   tls.RequireAndVerifyClientCert,
		}

		creds := credentials.NewTLS(tlsConfig)
		s = grpc.NewServer(grpc.Creds(creds))
	}

	// Register API
	message.RegisterGreeterServer(s, &greeterServer{})
	message.RegisterEchoServer(s, &echoServer{})

	// Enable reflection for grpcurl
	reflection.Register(s)

	log.Printf("Listening on %v", address)
	log.Fatal(s.Serve(lis))
}
