package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"saskara/blog-app-go/global"
	"saskara/blog-app-go/proto"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type authServer struct{}

func (authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error) {
	// return &proto.AuthResponse{}, nil
	login, password := in.GetLogin(), in.GetPassword()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var user global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": login}, bson.M{"email": login}}}).Decode(&user)
	// global.DB.Collection("user").FindOne(ctx, bson.M{"$or": []bson.M{"username": login, "email": login}}).Decode(&user)
	if user == global.NilUser {
		return &proto.AuthResponse{}, errors.New("Invalid Credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &proto.AuthResponse{}, errors.New("Invalid Credentials")
	}
	return &proto.AuthResponse{Token: user.GetToken()}, nil
}

func (authServer) Signup(_ context.Context, in *proto.SignupRequest) (*proto.AuthResponse, error) {
	username, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	match, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if len(username) < 4 || len(username) > 20 {
		return &proto.AuthResponse{}, errors.New("Username needs to be more than 4 and less than 20 characters")
	}
	if len(email) < 8 || len(email) > 50 || !match {
		return &proto.AuthResponse{}, errors.New("email needs to be valid")
	}
	if len(password) < 8 || len(password) > 50 {
		return &proto.AuthResponse{}, errors.New("password needs to be more than 8 characters")
	}
	return &proto.AuthResponse{}, nil
}

func (authServer) UsernameUsed(_ context.Context, in *proto.UsernameUsedRequest) (*proto.UsedResponse, error) {
	username := in.GetUsername()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	global.DB.Collection("user").FindOne(ctx, bson.M{"username": username}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}
func (authServer) EmailUsed(_ context.Context, in *proto.EmailUsedRequest) (*proto.UsedResponse, error) {
	return &proto.UsedResponse{}, nil
}

func main() {
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, authServer{})
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error", err)
	}
	server.Serve(listener)
}
