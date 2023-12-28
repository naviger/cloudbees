package main

import (
	"cloudbees/train"
	"context"
	"crypto/tls"
	"math/rand"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func TestGetTrain_Insecure(t *testing.T) {
	ctx := context.Background()

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainServiceClient(conn)

	dt := time.Now()
	resp, err := client.GetTrain(ctx, &train.TrainRequest{
		TrainId: dt.Format("20060102"),
		Year:    int32(dt.Year()),
		Month:   int32(dt.Month()),
		Day:     int32(dt.Day()),
	})

	if err != nil {
		t.Fatalf("GetTrainInsecure failed: %v", err)
	}
	if len(resp.Seats) != 20 {
		t.Fatalf("GetTrainSecure failed: Seat Count Wrong  (%d) - %v", len(resp.Seats), err)
	}
}

func TestGetTrain_Secure(t *testing.T) {
	md := metadata.New(map[string]string{"authorization": adminToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainServiceClient(conn)

	dt := time.Now()
	resp, err := client.GetTrain(ctx, &train.TrainRequest{
		TrainId: dt.Format("20060102"),
		Year:    int32(dt.Year()),
		Month:   int32(dt.Month()),
		Day:     int32(dt.Day()),
	})

	if err != nil {
		t.Fatalf("GetTrainSecure failed: %v", err)
	}
	if len(resp.Seats) != 20 {
		t.Fatalf("GetTrainSecure failed: Seat Count Wrong (%d) - %v", len(resp.Seats), err)
	}
}

func TestGetSeatUnsecured(t *testing.T) {
	//md := metadata.New(map[string]string{"authorization": adminToken})
	ctx := context.Background() //metadata.NewOutgoingContext(context.Background(), md)

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainServiceClient(conn)

	dt := time.Now()
	usr := SeedNames[rand.Intn(580)]
	usrArray := strings.Split(usr, " ")
	currentCustomer = usrArray[0] + "." + usrArray[1]

	resp, err := client.GetSeat(ctx, &train.SeatRequest{
		Id:            "someid",
		UserFirstName: usrArray[0],
		UserLastName:  usrArray[1],
		UserEmail:     usrArray[0] + "." + usrArray[1] + "@traveler.dev",
		TrainId:       dt.Format("20060102"),
		SeatValue:     dt.Format("20060102") + ":A1A",
	})

	if strings.ToLower(resp.Status) == "occupied" {
		freeSeat := GetFirstOpenSeat(dt)
		resp, err = client.GetSeat(ctx, &train.SeatRequest{
			Id:            "someid",
			UserFirstName: usrArray[0],
			UserLastName:  usrArray[1],
			UserEmail:     usrArray[0] + "." + usrArray[1] + "@traveler.dev",
			TrainId:       freeSeat.TrainId,
			SeatValue:     freeSeat.Id,
		})
	}

	if resp.Status != "success" {
		t.Fatalf("Failed to get seat: %v", err)
	}
}

func TestGetSeatSecured(t *testing.T) {
	md := metadata.New(map[string]string{"authorization": customerToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainServiceClient(conn)

	dt := time.Now()
	usr := testCustomer
	usrArray := strings.Split(usr, ".")

	resp, err := client.GetSeat(ctx, &train.SeatRequest{
		Id:            "someid",
		UserFirstName: usrArray[0],
		UserLastName:  usrArray[1],
		UserEmail:     usrArray[0] + "." + usrArray[1] + "@traveler.dev",
		TrainId:       dt.Format("20060102"),
		SeatValue:     dt.Format("20060102") + ":A1A",
	})

	if strings.ToLower(resp.Status) == "occupied" {
		currentSeat = GetFirstOpenSeat(dt)
		resp, err = client.GetSeat(ctx, &train.SeatRequest{
			Id:            "someid",
			UserFirstName: usrArray[0],
			UserLastName:  usrArray[1],
			UserEmail:     usrArray[0] + "." + usrArray[1] + "@traveler.dev",
			TrainId:       currentSeat.TrainId,
			SeatValue:     currentSeat.Id,
		})
	}

	if resp.Status != "success" {
		t.Fatalf("Failed to buy seat w/secured user: %v", err)
	}
}

func TestGetReceipt(t *testing.T) {
	md := metadata.New(map[string]string{"authorization": customerToken})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := train.NewTrainServiceClient(conn)

	rqst := &train.ReceiptRequest{
		Id:         "someid",
		TrainId:    currentSeat.TrainId,
		CustomerId: currentCustomer,
	}

	resp, err := client.GetReceipt(ctx, rqst)

	if resp.Id != rqst.Id ||
		resp.TrainId != currentSeat.TrainId {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
}

func TestCustomerCancelSeat(t *testing.T) {}

func TestAdminCancelSeat(t *testing.T) {}

func TestCustomerChangeSeat(t *testing.T) {}

func TestAdminChangeSeat(t *testing.T) {}

// variables
var testCustomer = "sheldon.carbajal"
var currentCustomer = ""
var currentSeat = SeatS{}

var adminToken string = "bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJCVGRGME5EY2UyZWhsQkl6RE43VGR2UnU3QUlfNVVFeWJnbXI2aGRLTjcwIn0.eyJleHAiOjE3MDM4MDQyMzcsImlhdCI6MTcwMzc2ODIzNywiYXV0aF90aW1lIjoxNzAzNzY4MjM3LCJqdGkiOiI3ZTdmMTVmNi1kNDI2LTQ0MDUtYjc0Zi01NzkxYzAwN2U5YWEiLCJpc3MiOiJodHRwczovL2Nsb3VkYmVlcy5kZXY6ODQ0My9yZWFsbXMvY2xvdWRiZWVzIiwiYXVkIjpbImNsb3VkYmVlcy1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjkzZDU1YmRhLWY2NzEtNGFhNy1hNmJlLWFhNjg0NTMyZDAxNyIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNsb3VkYmVlcy1jbGllbnQiLCJzZXNzaW9uX3N0YXRlIjoiODU2OWQ4MDYtNDUyMS00Yjg0LTgyZmQtNjljYTg0MjFhOTZmIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwczovL2Nsb3VkYmVlcy5kZXY6MzQ0MyJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1jbG91ZGJlZXMiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2xvdWRiZWVzLWNsaWVudCI6eyJyb2xlcyI6WyJ0cmF2ZWxfYWRtaW4iXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoib3BlbmlkIENsb3VkYmVlc0NvbW1vbiBlbWFpbCBwcm9maWxlIiwic2lkIjoiODU2OWQ4MDYtNDUyMS00Yjg0LTgyZmQtNjljYTg0MjFhOTZmIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJuYW1lIjoiQWRtaW4gQ2xvdWRiZWVzIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiYWRtaW4iLCJnaXZlbl9uYW1lIjoiQWRtaW4iLCJmYW1pbHlfbmFtZSI6IkNsb3VkYmVlcyIsImVtYWlsIjoiYWRtaW5AY2xvdWRiZWVzLmRldiJ9.Q1d6URcovocznZLkRMGNVWBOAnM4i4KqH8ln5osIFhRMotlpUukEmWnH8NOcUflLX_hfweos0BTsKIbifKa1q_--e6ssLdfOueOok3E_s3rn1SlL1TrTvPDw97-65lNnrwXW0Hm-cr5ZGQBhZjlcgj3VeOkTZqLCkh9A36jSJoldbQVKJV1jCrBsws_7sTCiWUHZlrpu2gUEjl2VPbDElxwQaq6Fe-OupyN8rSmgIexPT0fAzq2qIMOZ7xI5Dbv0hXWrfepvjTzLonzcFfw38UfeR_KHQ5TcHjtBV5RAr0Q5Q7eLPel_cLP__1N-OIiqueLSqPlQ-WvYeWnCLpqqLw"

var customerToken string = "bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJCVGRGME5EY2UyZWhsQkl6RE43VGR2UnU3QUlfNVVFeWJnbXI2aGRLTjcwIn0.eyJleHAiOjE3MDM4MjM1NTAsImlhdCI6MTcwMzc4NzU1MSwiYXV0aF90aW1lIjoxNzAzNzg3NTUwLCJqdGkiOiI1ZjMyYmY0Mi02ZDUwLTQyZDEtOTIzNi01ZWI1NTkzMDVlZjUiLCJpc3MiOiJodHRwczovL2Nsb3VkYmVlcy5kZXY6ODQ0My9yZWFsbXMvY2xvdWRiZWVzIiwiYXVkIjpbImNsb3VkYmVlcy1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjJjNWFjNjI5LTIwNTktNDM0OS1hMmRjLWIzODI3NjE2MjI0MiIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNsb3VkYmVlcy1jbGllbnQiLCJzZXNzaW9uX3N0YXRlIjoiOWRlM2VlZTMtZWNmMS00YzEwLWFhOGUtMGMyNDYxODI3MDRhIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwczovL2Nsb3VkYmVlcy5kZXY6MzQ0MyJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1jbG91ZGJlZXMiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2xvdWRiZWVzLWNsaWVudCI6eyJyb2xlcyI6WyJ0cmF2ZWxfY3VzdG9tZXIiXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoib3BlbmlkIENsb3VkYmVlc0NvbW1vbiBlbWFpbCBwcm9maWxlIiwic2lkIjoiOWRlM2VlZTMtZWNmMS00YzEwLWFhOGUtMGMyNDYxODI3MDRhIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJuYW1lIjoiU2hlbGRvbiBDYXJiYWphbCIsInByZWZlcnJlZF91c2VybmFtZSI6InNoZWxkb24uY2FyYmFqYWwiLCJnaXZlbl9uYW1lIjoiU2hlbGRvbiIsImZhbWlseV9uYW1lIjoiQ2FyYmFqYWwiLCJlbWFpbCI6InNoZWxkb24uY2FyYmFqYWxAdHJhdmVsZXIuZGV2In0.JHZGsIc_clnoAI7-PUlf7Mp4XBWllqUA7sjAMj3JTEMbgAnwL5XE4mJS721AiC0Q8ADxgFOGJ-E32Qgr1o4f-EQFEqG-uj9OKYfHK722PobFGdOxl9Oyrve3Gyi5P-74pB5ObqRuTeG0773PBv0b9BJ2BXo9OybP2dqV2SHaHlNtfVipvzYlCJDtZsF7RBcYOlsJz8bI9O00CBM8Y_CCtuMxjVR-0K9nQvivIN4CoVq9YwBK9ktU1K1bGAmYW5ijRbV1F-rXZsrrxCEY7cPe-fEJDUQ6QumQTZYB-N4E-SXqUxjcV90IPAGYa0vlypin4K3mihFeuOsQI-o9kIsaQw"
