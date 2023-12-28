package main

import (
	ints "cloudbees/interceptor"
	train "cloudbees/train"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-memdb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type server struct {
	train.TrainServiceServer
}

var db *memdb.MemDB

func (s *server) GetTrain(ctx context.Context, req *train.TrainRequest) (*train.TrainReply, error) {
	var seats []*train.TrainReply_Seat = make([]*train.TrainReply_Seat, 0)

	md, _ := metadata.FromIncomingContext(ctx)
	isAuth := GetMetadataKey(md, "isAuth", "false")
	isAdmin := GetMetadataKey(md, "isAdmin", "false")

	txn := db.Txn(false)
	defer txn.Abort()

	it, err := txn.LowerBound("Seat", "id", req.TrainId+":")
	if err != nil {
		panic(err)
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(SeatS)

		if p.TrainId == req.TrainId {
			pid := "redacted"
			if isAuth == "true" && isAdmin == "true" {
				pid = p.PassengerId
			}

			st := &train.TrainReply_Seat{
				Id:          p.Id,
				Row:         p.Row,
				Position:    p.Position,
				Status:      p.Status,
				PassengerId: pid,
			}

			seats = append(seats, st)
		}
	}

	reply := train.TrainReply{
		TrainId: req.TrainId,
		Day:     req.Day,
		Month:   req.Month,
		Year:    req.Year,
		Seats:   seats,
	}

	return &reply, nil
}

func (s *server) GetSeat(ctx context.Context, req *train.SeatRequest) (*train.SeatReply, error) {
	token := GetAdminToken()

	md, _ := metadata.FromIncomingContext(ctx)
	isAuth := GetMetadataKey(md, "isAuth", "false")
	isCustomer := GetMetadataKey(md, "isCustomer", "false")

	username := GetMetadataKey(md, "user", "")
	if isAuth == "true" && isCustomer == "true" && (username != req.UserFirstName+"."+req.UserLastName) {
		return &train.SeatReply{
			Id:         req.Id,
			Status:     "failure: Login User does not match customer.",
			CustomerId: req.UserEmail,
			Seat:       &train.SeatReply_Seat{},
		}, nil
	}

	if isAuth != "true" || (isAuth == "true" && isCustomer == "true") {
		txn := db.Txn(true)
		defer txn.Abort()

		o, err := txn.First("Seat", "id", req.SeatValue)
		if err != nil {
			panic(err)
		}

		var st SeatS = o.(SeatS)

		//if seat is occupied, reject
		if st.Status == "occupied" {
			return &train.SeatReply{
				Id:         req.Id,
				Status:     st.Status,
				CustomerId: "redacted",
			}, err
		} else {
			success := false
			if isAuth == "false" {
				userid := CreateUser(token, req.UserFirstName, req.UserLastName, req.UserEmail)
				if len(userid) > 0 {
					success = AssignUserAsClientRole(token, userid, "travel_customer")
				}
			} else {
				success = true
			}

			if success {
				st.Status = "occupied"
				st.PassengerId = req.UserFirstName + "." + req.UserLastName

				txn.Insert("Seat", st)
				txn.Commit()

				return &train.SeatReply{
					Id:         req.Id,
					Status:     "success",
					CustomerId: req.UserFirstName + "." + req.UserLastName,
					Seat: &train.SeatReply_Seat{
						Id:          st.Id,
						Car:         st.Car,
						Row:         st.Row,
						Position:    st.Position,
						Status:      "occupied",
						PassengerId: req.UserFirstName + "." + req.UserLastName,
					},
				}, err
			}
		}
	} else {
		return &train.SeatReply{
			Id:         req.Id,
			CustomerId: req.UserEmail,
			Status:     fmt.Sprintf("failure: user '%s' not recognized ", username),
			Seat:       &train.SeatReply_Seat{},
		}, nil
	}
	return &train.SeatReply{
		Id:         req.Id,
		CustomerId: req.UserEmail,
		Status:     "failure: unknown",
		Seat:       &train.SeatReply_Seat{},
	}, nil
}

func (s *server) GetReceipt(ctx context.Context, req *train.ReceiptRequest) (*train.ReceiptReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	isAuth := GetMetadataKey(md, "isAuth", "false")
	isCustomer := GetMetadataKey(md, "isCustomer", "false")

	username := GetMetadataKey(md, "user", "")
	seats := make([]*train.ReceiptReply_Seat, 0)

	if isAuth != "true" || isCustomer != true {
		return &train.ReceiptReply{
			Id:      req.Id,
			TrainId: req.TrainId,
			Seats:   seats,
		}, nil
	}

	txn := db.Txn(false)

	it, _ := txn.LowerBound("seat", "id", req.TrainId)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(SeatS)

		if p.TrainId == req.TrainId && p.PassengerId == req.CustomerId && req.CustomerId == username {
			seats = append(seats, &train.ReceiptReply_Seat{
				Id:          p.Id,
				Car:         p.Car,
				Row:         p.Row,
				Position:    p.Position,
				Status:      p.Status,
				PassengerId: p.PassengerId,
			})
		}
	}

	return &train.ReceiptReply{
		Id:      req.Id,
		TrainId: req.TrainId,
		Seats:   seats,
	}, nil
}

func (s *server) ChangeSeat(ctx context.Context, req *train.ChangeSeatRequest) (*train.ChangeSeatReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	isAuth := GetMetadataKey(md, "isAuth", "false")
	isCustomer := GetMetadataKey(md, "isCustomer", "false")

	username := GetMetadataKey(md, "user", "")

	if isAuth != "true" {
		return &train.ChangeSeatReply{
			Id:         req.Id,
			Status:     "Unauthorized: User unknown",
			SourceSeat: &train.TrainReply_Seat{},
			DestSeat:   &train.TrainReply_Seat{},
		}, nil
	} else if isCustomer == "true" && (req.CustomerId != username) {
		return &train.ChangeSeatReply{
			Id:         req.Id,
			Status:     "Unauthorized: User unknown",
			SourceSeat: &train.TrainReply_Seat{},
			DestSeat:   &train.TrainReply_Seat{},
		}, nil
	} else {
		txn := db.Txn(true)
		od, err := txn.First("Seat", "id", req.DestSeat)
		if err != nil {
			panic(err)
		}
		sd := od.(SeatS)

		os, err := txn.First("Seat", "id", req.DestSeat)
		if err != nil {
			panic(err)
		}
		ss := os.(SeatS)

		if sd.Status == "occupied" {
			return &train.ChangeSeatReply{
				Id:         req.Id,
				Status:     "Rejected: Target Seat is occupied",
				SourceSeat: &train.TrainReply_Seat{},
				DestSeat:   &train.TrainReply_Seat{},
			}, nil
		} else if ss.PassengerId != username && isCustomer == "true" {
			return &train.ChangeSeatReply{
				Id:         req.Id,
				Status:     "Rejected: Requesting customer does not have source seat permissions",
				SourceSeat: &train.TrainReply_Seat{},
				DestSeat:   &train.TrainReply_Seat{},
			}, nil
		} else {
			sd.PassengerId = ss.PassengerId
			sd.Status = "occupied"
			ss.PassengerId = ""
			ss.Status = "available"

			txn.Insert("Seat", ss)
			txn.Insert("Seat", sd)

			return &train.ChangeSeatReply{
				Id:     req.Id,
				Status: "Success",
				SourceSeat: &train.TrainReply_Seat{
					Id:          ss.Id,
					Car:         ss.Car,
					Row:         ss.Row,
					Position:    ss.Position,
					Status:      ss.Status,
					PassengerId: ss.PassengerId,
				},
				DestSeat: &train.TrainReply_Seat{
					Id:          sd.Id,
					Car:         sd.Car,
					Row:         sd.Row,
					Position:    sd.Position,
					Status:      sd.Status,
					PassengerId: sd.PassengerId,
				},
			}, nil
		}
	}
}

func (s *server) CancelSeat(ctx context.Context, req *train.CancelSeatRequest) (*train.CancelSeatReply, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	isAuth := GetMetadataKey(md, "isAuth", "false")
	isCustomer := GetMetadataKey(md, "isCustomer", "false")

	username := GetMetadataKey(md, "user", "")

	if isAuth != "true" {
		return &train.CancelSeatReply{
			Id:     req.Id,
			Status: "Unauthorized: User unknown",
		}, nil
	} else if isCustomer == "true" && (req.CustomerId != username) {
		return &train.CancelSeatReply{
			Id:     req.Id,
			Status: "Unauthorized: Requesting customer does not have seat permissions",
		}, nil
	} else {
		txn := db.Txn(true)
		o, err := txn.First("Seat", "id", req.SeatId)
		if err != nil {
			panic(err)
		}
		st := o.(SeatS)
		st.PassengerId = ""
		st.Status = "vacant"

		txn.Insert("Seat", st)

		return &train.CancelSeatReply{
			Id:     req.Id,
			Status: "Success",
		}, nil
	}
}

func Server() {
	caPem, err := os.ReadFile("./certs/RootCA.pem")
	if err != nil {
		log.Fatal(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		log.Fatal(err)
	}

	cert, _ := tls.LoadX509KeyPair("./certs/localhost.crt", "./certs/localhost.key")

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(ints.JwtInterceptor),
	}

	address := "cloudbees.dev:5443"
	lis, err := net.Listen("tcp", address)

	db, err = memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	SetupTrainDB(db)

	if err != nil {
		log.Fatalf("Error %v", err)
	}

	log.Println("Server Started on ", lis.Addr())

	srv := grpc.NewServer(opts...)

	train.RegisterTrainServiceServer(srv, &server{})
	reflection.Register(srv)
	log.Fatalln(srv.Serve(lis))
	fmt.Printf("Listening . . . ")
}

func SetupTrainDB(db *memdb.MemDB) {
	//set up for 30 days in advance
	dt := time.Now()
	log.Println("Setting up train . . .", dt)
	for i := 0; i < 30; i++ {
		txn := *db.Txn(true)

		var t TrainS = TrainS{
			Id:    dt.Format("20060102"),
			Year:  int32(dt.Year()),
			Month: int32(dt.Month()),
			Day:   int32(dt.Day()),
		}

		if err := txn.Insert("Train", t); err != nil {
			panic(err)
		}
		txn.Commit()

		SetupSeats(db, t.Id, i)

		dt = dt.AddDate(0, 0, 1)
	}
}

func SetupSeats(db *memdb.MemDB, id string, day int) {
	max := 90
	min := 10
	step := (max - min) / 30
	cur := min

	var car string = "A"
	var r int32 = 1
	var p string = "A"

	txn := db.Txn(true)

	for i := 0; i < 20; i++ {

		// determine if seat is free . . . as the day (0 - 29) gets bigger, the more likely a seat is free
		seed := rand.Intn(100)
		var vacant bool = seed < cur
		cur += step

		if i == 10 {
			car = "B"
			r = 1
			p = "A"
		}

		var s SeatS = SeatS{
			Id:          id + ":" + car + strconv.FormatInt(int64(r), 10) + p,
			TrainId:     id,
			Car:         car,
			Row:         r,
			Position:    p,
			Status:      "vacant",
			PassengerId: "-",
		}

		if !vacant {
			s.Status = "occupied"
			s.PassengerId = strings.ReplaceAll(SeedNames[rand.Intn(500)], " ", ".") + "@customer.dev"
		}

		if err := txn.Insert("Seat", s); err != nil {
			panic(err)
		}

		switch p {
		case "A":
			p = "B"
			break
		case "B":
			p = "C"
			break
		case "C":
			p = "D"
			break
		case "D":
			p = "A"
			r++
			break
		}
	}
	txn.Commit()
}

type TrainS struct {
	Id    string
	Year  int32
	Month int32
	Day   int32
}

type SeatS struct {
	Id          string
	TrainId     string
	Car         string
	Row         int32
	Position    string
	Status      string
	PassengerId string
}
