package main

import (
	"context"
	"saskara/rad-blog-golang/global"
	"saskara/rad-blog-golang/proto"
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
	// t.Error(res.GetToken())

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

func Test_authServer_AuthUser(t *testing.T) {
	server := authServer{}
	res, err := server.AuthUser(context.Background(), &proto.AuthUserRequest{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjoie1wiSURcIjpcIjYwMzM3NjA2YTg2YjNjOTVhMThkODk3NFwiLFwiVXNlcm5hbWVcIjpcInNhc2thdGVzdFwiLFwiRW1haWxcIjpcInJhZHRlc3RAcmFkLmNvbVwiLFwiUGFzc3dvcmRcIjpcIiQyYSQxMCRMYmtQUnN0UE1XcnhTT1l0U0VkRlplV0YycVhwSEsxd1hHUVVoZkg2ZzliM3hiLlFUdkdkYVwifSJ9.b0w6FHnf8oUxdvduwcu4v2F8l4tk1g9mO7OWOsZOuDI"})
	if err != nil {
		t.Error("there was an error")
	}
	if res.GetID() != "60337606a86b3c95a18d8974" || res.GetUsername() != "saskatest" || res.GetEmail() != "radtest@rad.com" {
		t.Error("There was an error", res)
	}

}
