package main

import (
	"context"
	"log"
	"saskara/rad-blog-golang/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	client := proto.NewAuthServiceClient(conn)
	client.Signup(context.Background(), &proto.SignupRequest{Username: "saskara", Email: "rad@gmail.com", Password: "saskarad"})
}
