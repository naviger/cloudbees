package controllers

import (
	"cloudbees-api/train"
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
)

func GetTrains() func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := context.Background()

		creds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: false})
		conn, err := grpc.Dial("cloudbees.dev:5443", grpc.WithTransportCredentials(creds))

		if err != nil {
			log.Fatalf("Failed to dial bufnet: %v", err)
		}
		defer conn.Close()

		client := train.NewTrainServiceClient(conn)
		resp, err := client.GetTrains(ctx, &train.TrainsRequest{})
		c.JSON(http.StatusOK, gin.H{"data": resp.Trains})
	}
}

func GetTrain() func(c *gin.Context) {
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
		dateValue, _ := time.Parse("20060102", trainId)

		client := train.NewTrainServiceClient(conn)
		resp, err := client.GetTrain(ctx, &train.TrainRequest{
			TrainId: trainId,
			Year:    int32(dateValue.Year()),
			Month:   int32(dateValue.Month()),
			Day:     int32(dateValue.Day()),
		})

		msg, _ := protojson.Marshal(resp)
		c.JSON(http.StatusOK, gin.H{"data": string(msg[:])})
	}
}
