package controllers

import (
	"cloudbees-api/train"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

func RequestSeat() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()

		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
		conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

		if err != nil {
			log.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		var isAuth bool = c.GetBool("isAuth")
		if isAuth {
			token := c.Request.Header["Authorization"][0]
			log.Println("SET CONTEXT:", isAuth, token)
			md := metadata.New(map[string]string{"authorization": token})
			ctx = metadata.NewOutgoingContext(context.Background(), md)
		}

		trainId := c.Param("trainId")

		var requestBody RequestSeatBody
		if err := c.BindJSON(&requestBody); err != nil {
			log.Fatal("Error binding JSON")
		}

		client := train.NewTrainServiceClient(conn)
		resp, err := client.GetSeat(ctx, &train.SeatRequest{
			Id:            time.Now().Local().String(),
			UserFirstName: strings.ToLower(requestBody.FirstName),
			UserLastName:  strings.ToLower(requestBody.LastName),
			UserEmail:     requestBody.Email,
			TrainId:       trainId,
			SeatValue:     trainId + ":" + requestBody.Seat,
		})

		msg, _ := protojson.Marshal(resp)
		c.JSON(http.StatusOK, gin.H{"data": string(msg[:])})
	}
}

func RequestReceipt() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()

		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
		conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

		if err != nil {
			log.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		var isAuth bool = c.GetBool("isAuth")
		if isAuth {
			token := c.Request.Header["Authorization"][0]
			log.Println("SET CONTEXT:", isAuth, token)
			md := metadata.New(map[string]string{"authorization": token})
			ctx = metadata.NewOutgoingContext(context.Background(), md)
		}

		trainId := c.Param("trainId")
		seatId := c.Param("seatId")
		customerId := c.Query("customerId")

		client := train.NewTrainServiceClient(conn)
		resp, err := client.GetReceipt(ctx, &train.ReceiptRequest{
			Id:         time.Now().Local().String(),
			CustomerId: customerId,
			SeatId:     trainId + ":" + seatId,
		})

		//obj, _ := protojson.Marshal(resp)

		msg := "ID:       " + resp.Id + "\n" +
			"TRAIN:    " + resp.TrainId + "\n" +
			"SEAT:     " + resp.Seat.Id + "\n" +
			"FROM:     " + resp.From + "\n" +
			"TO:       " + resp.To + "\n" +
			"PRICE:    " + fmt.Sprintf("%.2f", resp.Price) + "\n" +
			"CUSTOMER: \n" +
			"   USERNAME:  " + resp.User.UserId + "\n" +
			"   FIRST:     " + resp.User.FirstName + "\n" +
			"   LAST:      " + resp.User.LastName + "\n" +
			"   EMAIL:     " + resp.User.EmailAddress + "\n"
		c.JSON(http.StatusOK, gin.H{"data": string(msg[:])})
	}
}

func ChangeSeat() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()

		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
		conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

		if err != nil {
			log.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		var isAuth bool = c.GetBool("isAuth")
		if isAuth {
			token := c.Request.Header["Authorization"][0]
			log.Println("SET CONTEXT:", isAuth, token)
			md := metadata.New(map[string]string{"authorization": token})
			ctx = metadata.NewOutgoingContext(context.Background(), md)
		}

		trainId := c.Param("trainId")
		//seatId := c.Param("seatId")

		var requestBody RequestChangeBody
		if err := c.BindJSON(&requestBody); err != nil {
			log.Fatalf("ERROR BINDING JSON %v", err)
		}

		client := train.NewTrainServiceClient(conn)
		resp, err := client.ChangeSeat(ctx, &train.ChangeSeatRequest{
			Id:         time.Now().Local().String(),
			CustomerId: requestBody.CustomerId,
			SourceSeat: trainId + ":" + requestBody.Source,
			DestSeat:   trainId + ":" + requestBody.Dest,
		})

		msg, _ := protojson.Marshal(resp)
		log.Println("RESPONSE: ", msg, resp)
		c.JSON(http.StatusOK, gin.H{"data": string(msg[:])})
	}
}

func CancelSeat() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()

		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
		conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

		if err != nil {
			log.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		var isAuth bool = c.GetBool("isAuth")
		if isAuth {
			token := c.Request.Header["Authorization"][0]
			log.Println("SET CONTEXT:", isAuth, token)
			md := metadata.New(map[string]string{"authorization": token})
			ctx = metadata.NewOutgoingContext(context.Background(), md)
		}

		seatId := c.Param("seatId")
		trainId := c.Param("trainId")

		var requestBody RequestCancelBody
		if err := c.BindJSON(&requestBody); err != nil {
			log.Fatalf("ERROR BINDING JSON %v", err)
		}

		client := train.NewTrainServiceClient(conn)
		resp, err := client.CancelSeat(ctx, &train.CancelSeatRequest{
			Id:         time.Now().Local().String(),
			CustomerId: requestBody.CustomerId,
			SeatId:     trainId + ":" + seatId,
		})

		msg, _ := protojson.Marshal(resp)
		c.JSON(http.StatusOK, gin.H{"data": string(msg[:])})
	}
}

type RequestSeatBody struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Seat      string `json:"seat"`
}

type RequestChangeBody struct {
	Source     string `json:"source"`
	Dest       string `json:"dest"`
	CustomerId string `json:"customerId"`
}

type RequestCancelBody struct {
	Seat       string `json:"seat"`
	CustomerId string `json:"customerId"`
}
