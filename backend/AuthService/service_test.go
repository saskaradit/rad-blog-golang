package main

import (
	"context"
	"saskara/blog-app-go/global"
	"saskara/blog-app-go/proto"
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Test_authServer_Login(t *testing.T) {
	global.ConnectToTestDB()
	pass, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	global.DB.Collection("user").InsertOne(context.Background(), global.User{ID: primitive.NewObjectID(), Email: "radtest@rad.com", Username: "saskatest", Password: string(pass)})

	server := authServer{}
	_, err := server.Login(context.Background(), &proto.LoginRequest{Login: "radtest@rad.com", Password: "test"})
	if err != nil {
		t.Error("An Error Occurred ", err.Error())
	}
	_, err = server.Login(context.Background(), &proto.LoginRequest{Login: "jengjet", Password: "Jengjet"})
	if err == nil {
		t.Error("Error was nil")
	}

	_, err = server.Login(context.Background(), &proto.LoginRequest{Login: "saskatest", Password: "test"})
	if err != nil {
		t.Error("Error was nil")
	}
}

func Test_authServer_UsernameUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "Rad"})
	server := authServer{}
	res, err := server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: "Rado"})
	if err != nil {
		t.Error("An error was returned", err.Error())
	}
	if res.GetUsed() {
		t.Error("Wrong result")
	}
	res, err = server.UsernameUsed(context.Background(), &proto.UsernameUsedRequest{Username: "Rad"})
	if err != nil {
		t.Error("No Error", err.Error())
	}
	if !res.GetUsed() {
		t.Error("Wrong result")
	}

}

func Test_authServer_EmailUsed(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Email: "radtest@rad.com"})
	server := authServer{}
	res, err := server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "radtest@google.com"})
	if err != nil {
		t.Error("An error was returned", err.Error())
	}
	if res.GetUsed() {
		t.Error("Wrong result")
	}
	res, err = server.EmailUsed(context.Background(), &proto.EmailUsedRequest{Email: "radtest@rad.com"})
	if err != nil {
		t.Error("No Error", err.Error())
	}
	if !res.GetUsed() {
		t.Error("Wrong result")
	}
}

func Test_authServer_Signup(t *testing.T) {
	global.ConnectToTestDB()
	global.DB.Collection("user").InsertOne(context.Background(), global.User{Username: "saska", Email: "radtest@rad.com"})
	server := authServer{}
	_, err := server.Signup(context.Background(), &proto.SignupRequest{Username: "saska", Email: "radtest@gmail.com", Password: "jengjeta"})
	if err.Error() != "Username is used" {
		t.Error("1. Wrong error was returned")
	}
	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "saskara", Email: "radtest@rad.com", Password: "jengjeta"})
	if err.Error() != "Email is used" {
		// fmt.Println(err.Error())
		t.Error(err.Error())
		// t.Error("2. Wrong error was returned")
	}
	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "saskara", Email: "rad@rad.com", Password: "jengjeta"})
	if err != nil {
		t.Error("3. error was returned")
	}
	_, err = server.Signup(context.Background(), &proto.SignupRequest{Username: "saskara", Email: "rad@rad.com", Password: "rad"})
	if err.Error() != "Validation error" {
		t.Error("4. Wrong error was returned")
	}
}
