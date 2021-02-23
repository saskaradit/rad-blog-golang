package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"
	"saskara/rad-blog-golang/global"
	"saskara/rad-blog-golang/proto"

	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type authServer struct{}

var userCollection mongo.Collection

func (authServer) Login(_ context.Context, in *proto.LoginRequest) (*proto.AuthResponse, error) {
	// return &proto.AuthResponse{}, nil
	login, password := in.GetLogin(), in.GetPassword()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var user global.User
	userCollection.FindOne(ctx, bson.M{"$or": []bson.M{bson.M{"username": login}, bson.M{"email": login}}}).Decode(&user)
	// userCollection.FindOne(ctx, bson.M{"$or": []bson.M{"username": login, "email": login}}).Decode(&user)
	if user == global.NilUser {
		return &proto.AuthResponse{}, errors.New("Invalid Credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return &proto.AuthResponse{}, errors.New("Invalid Credentials")
	}
	return &proto.AuthResponse{Token: user.GetToken()}, nil
}

func (server authServer) Signup(_ context.Context, in *proto.SignupRequest) (*proto.AuthResponse, error) {
	username, email, password := in.GetUsername(), in.GetEmail(), in.GetPassword()
	match, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", email)
	if len(username) < 4 || len(username) > 20 || len(email) < 8 || len(email) > 50 || len(password) < 8 || len(password) > 50 || !match {
		return &proto.AuthResponse{}, errors.New("Validation error")
	}

	res, err := server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: username})
	if err != nil {
		log.Fatal("Error signing up from username", err.Error())
		return &proto.AuthResponse{}, errors.New("Something went wrong")
	}
	if res.GetUsed() {
		// log.Fatal("Usernam", err.Error())
		return &proto.AuthResponse{}, errors.New("Username is used")
	}

	res, err = server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: email})
	if err != nil {
		log.Fatal("Error signing up from email", err.Error())
		return &proto.AuthResponse{}, errors.New("Something went wrong")
	}
	if res.GetUsed() {
		// log.Fatal("Email", err.Error())
		return &proto.AuthResponse{}, errors.New("Email is used")
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	newUser := global.User{ID: primitive.NewObjectID(), Username: username, Email: email, Password: string(pass)}

	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	_, err = userCollection.InsertOne(ctx, newUser)
	if err != nil {
		log.Println("Error inserting user to db")
		return &proto.AuthResponse{}, errors.New("Something went wrong")
	}
	return &proto.AuthResponse{Token: newUser.GetToken()}, nil
}

func (authServer) UsernameUsed(_ context.Context, in *proto.UsernameUsedRequest) (*proto.UsedResponse, error) {
	username := in.GetUsername()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}
func (authServer) EmailUsed(_ context.Context, in *proto.EmailUsedRequest) (*proto.UsedResponse, error) {
	email := in.GetEmail()
	ctx, cancel := global.NewDBContext(5 * time.Second)
	defer cancel()
	var result global.User
	userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&result)
	return &proto.UsedResponse{Used: result != global.NilUser}, nil
}

func (authServer) AuthUser(_ context.Context, in *proto.AuthUserRequest) (*proto.AuthUserResponse, error) {
	token := in.GetToken()
	user := global.UserFromToken(token)
	return &proto.AuthUserResponse{ID: user.ID.Hex(), Username: user.Username, Email: user.Email}, nil
}

func main() {
	userCollection = *global.DB.Collection("user")
	server := grpc.NewServer()
	proto.RegisterAuthServiceServer(server, authServer{})
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error", err)
	}

	go func() {
		log.Fatal("serving gRPC: ", server.Serve(listener).Error())
		// server.Serve(listener)
	}()

	grpcWebServer := grpcweb.WrapServer(server)
	httpServer := &http.Server{
		Addr: ":9001",
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}
	log.Fatal("Serving Proxy: ", httpServer.ListenAndServe().Error())
}
