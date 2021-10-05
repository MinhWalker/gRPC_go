package main

import (
	"context"
	"fmt"
	"gRPC-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"math"
	"net"
	"time"
)

type server struct{}

// TODO: Handle error
func (s *server) Square(ctx context.Context, req *calculatorpb.SquareRequest) (*calculatorpb.SquareResponse, error) {
	log.Println("Square called ...")
	num := req.GetNum()
	if num < 0 {
		log.Printf("req num < 0, num = %v, return InvalidArgument", num)
		return nil, status.Errorf(codes.InvalidArgument, "Expect num larger than 0, request num was %v", num)
	}

	return &calculatorpb.SquareResponse{SquareRoot: math.Sqrt(float64(num))}, nil
}

// TODO: BI-Direction streaming API
func (s *server) Max(stream calculatorpb.CalculatorService_MaxServer) error {
	log.Println("Max called ...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("client finish streaming", err)
			return nil
		}
		if err != nil {
			log.Fatalln("err while receive request", err)
			return err
		}
		log.Printf("receive num %d", req.GetNum())
		if req.GetNum() >= max {
			max = req.GetNum()
		}
		err = stream.Send(&calculatorpb.MaxResponse{Result: max})
		if err != nil {
			log.Fatalln("send max err %v", err)
			return err
		}
		//time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

// TODO: Client streaming API
func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average called ...")
	var average float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			resq := &calculatorpb.AverageResponse{Result: average / float32(count)}

			return stream.SendAndClose(resq)
		}
		if err != nil {
			log.Fatalln("err while receive request", err)
		}
		log.Println("receive num %v", req.GetNum1())
		average += req.GetNum1()
		count++
	}

	return nil
}

// TODO: Server streaming API
func (*server) PrimeNumberDecomposition(req *calculatorpb.PNDRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k
			//sent to client
			stream.Send(&calculatorpb.PNDResponse{Result: k})
			time.Sleep(500 * time.Millisecond)
		} else {
			k++
			log.Printf("k increase to %v", k)
		}
	}
	return nil
}

// TODO: unary API
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("Sum called ...")
	resq := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}
	return resq, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50069")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}
	certFile := "ssl/server.crt"
	keyFile := "ssl/server.pem"

	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("create creds ssl err %v\n", sslErr)
		return
	}
	opts := grpc.Creds(creds)

	s := grpc.NewServer(opts)

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("calculator is running")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalln("err while serve %v", err)
	}
}
