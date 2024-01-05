package main

import (
	"cloudbees/train"
	"context"
	"crypto/tls"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func TestGetTrains(t *testing.T) {
	ctx := context.Background()

	creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
	conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := train.NewTrainServiceClient(conn)

	resp, err := client.GetTrains(ctx, &train.TrainsRequest{})
	if len(resp.Trains) != 30 {
		t.Fatalf("Did not receive all trains")
	}
	if err != nil {
		t.Fatalf("Error receiving trains: %v", err)
	}
}
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

	log.Println("CurrentSeat:", currentSeat, testCustomer)
	rqst := &train.ReceiptRequest{
		Id:         "someid",
		SeatId:     currentSeat.Id,
		CustomerId: testCustomer,
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

var adminToken string = "bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJCVGRGME5EY2UyZWhsQkl6RE43VGR2UnU3QUlfNVVFeWJnbXI2aGRLTjcwIn0.eyJleHAiOjE3MDQ0MjY0MTIsImlhdCI6MTcwNDM5MDQxMywiYXV0aF90aW1lIjoxNzA0MzkwNDEyLCJqdGkiOiI5NWI4YWVhZi00NDY1LTQ2YTctYjg3Zi0wODFlZjhlN2Y1MjkiLCJpc3MiOiJodHRwczovL2Nsb3VkYmVlcy5kZXY6ODQ0My9yZWFsbXMvY2xvdWRiZWVzIiwiYXVkIjpbImNsb3VkYmVlcy1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjkzZDU1YmRhLWY2NzEtNGFhNy1hNmJlLWFhNjg0NTMyZDAxNyIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNsb3VkYmVlcy1jbGllbnQiLCJzZXNzaW9uX3N0YXRlIjoiZjlkMTNkNGUtYzYwZS00ZTgxLWI4ZGQtYTkyODg0ZGYxYzg4IiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwczovL2Nsb3VkYmVlcy5kZXY6MzQ0MyJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1jbG91ZGJlZXMiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2xvdWRiZWVzLWNsaWVudCI6eyJyb2xlcyI6WyJ0cmF2ZWxfYWRtaW4iXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoib3BlbmlkIENsb3VkYmVlc0NvbW1vbiBlbWFpbCBwcm9maWxlIiwic2lkIjoiZjlkMTNkNGUtYzYwZS00ZTgxLWI4ZGQtYTkyODg0ZGYxYzg4IiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJuYW1lIjoiQWRtaW4gQ2xvdWRiZWVzIiwicHJlZmVycmVkX3VzZXJuYW1lIjoiYWRtaW4iLCJnaXZlbl9uYW1lIjoiQWRtaW4iLCJmYW1pbHlfbmFtZSI6IkNsb3VkYmVlcyIsImVtYWlsIjoiYWRtaW5AY2xvdWRiZWVzLmRldiJ9.D_y1Gc-_w2Ep4I_TsKanypcPOcZ95G92KfkqliWTDbrp3GqDJ-o9TtyXRY40ZziygQLnnbvRwyUAI9TeHTD2gcWtdnSw8_JKeuF8aq3P_jopjxl05Y8pLPFu2Naa_8zG6UMPkF6GYf0x-SpCaXxS-ASXl30EXr3aR9XPW9JOiFYlc6SRtWT3qrDOihFpeC2L0uGnthAAbuV1zPWMCZ5mhklp4DdY3bgjePddgxB4nJPkmZmaKUfCKOKObq1NAEuHvnHtzI8IbUJ-9osUAEVyXqG1736FYWzQQERgBlLm7L3JkFlhJ67EH-E31sJjTonSFbDPKWudxZbRDnCZRUbvDw"

var customerToken string = "bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJCVGRGME5EY2UyZWhsQkl6RE43VGR2UnU3QUlfNVVFeWJnbXI2aGRLTjcwIn0.eyJleHAiOjE3MDQ0MjY1MTUsImlhdCI6MTcwNDM5MDU3NywiYXV0aF90aW1lIjoxNzA0MzkwNTE1LCJqdGkiOiI4YTEwOGQwNC03MTVmLTQ5YjctOGUzMS0xNmI2ZGY5Mjk3YTUiLCJpc3MiOiJodHRwczovL2Nsb3VkYmVlcy5kZXY6ODQ0My9yZWFsbXMvY2xvdWRiZWVzIiwiYXVkIjpbImNsb3VkYmVlcy1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjJjNWFjNjI5LTIwNTktNDM0OS1hMmRjLWIzODI3NjE2MjI0MiIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNsb3VkYmVlcy1jbGllbnQiLCJzZXNzaW9uX3N0YXRlIjoiZGNjNjdkOGUtNzEyYS00ZjE1LWJmNjktNmExYjcxOTljNTBhIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyJodHRwczovL2Nsb3VkYmVlcy5kZXY6MzQ0MyJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy1jbG91ZGJlZXMiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2xvdWRiZWVzLWNsaWVudCI6eyJyb2xlcyI6WyJ0cmF2ZWxfY3VzdG9tZXIiXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoib3BlbmlkIENsb3VkYmVlc0NvbW1vbiBlbWFpbCBwcm9maWxlIiwic2lkIjoiZGNjNjdkOGUtNzEyYS00ZjE1LWJmNjktNmExYjcxOTljNTBhIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJuYW1lIjoiU2hlbGRvbiBDYXJiYWphbCIsInByZWZlcnJlZF91c2VybmFtZSI6InNoZWxkb24uY2FyYmFqYWwiLCJnaXZlbl9uYW1lIjoiU2hlbGRvbiIsImZhbWlseV9uYW1lIjoiQ2FyYmFqYWwiLCJlbWFpbCI6InNoZWxkb24uY2FyYmFqYWxAdHJhdmVsZXIuZGV2In0.Q_hKtRv7ct3Fkt5a7h0MBa2Njl01iOioHWBUKSknM1mZHUKM3qf1T12xoatkuK-cm_WrLslO1W-535FDmr2gOJtGY1uYSFZhjJPRkIYd4BhTWgw0-eEACbkDchNW6pvSJb7c_dtSbV7f6He1i8dRyr5eokGz6bDfKn2wE44pfBiclOwBgvCNQ-vNfHWPsmrfCJ3PesWQuxA28M9MtiMlsVzDdSKB4VZn2wRud6TLtiOA-qz3Qct28M-ky2cVG3vYPQkhNYCPxw1sRBaIINc5u6yxIu29iwxm60KvMtD7tFY1gXsMTLxjEUeusCZRpG0A0a7hUhqc853eZnldyMYFGg"
