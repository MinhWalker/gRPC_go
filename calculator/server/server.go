package main

import (
	"context"
	"fmt"
	"gRPC-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)

type server struct{}

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
		log.Fatalln("err while create listen %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("calculator is running")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalln("err while serve %v", err)
	}
}
