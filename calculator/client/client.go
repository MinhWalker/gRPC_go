package main

import (
	"context"
	"gRPC-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalln("err when dial %v", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	//log.Printf("services client %f", client)
	//callSum(client)
	//callPND(client)
	//callAverage(client)
	callMax(client)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 5,
		Num2: 5,
	})
	if err != nil {
		log.Fatalln("call sum api err %v", err)
	}

	log.Printf("sum api response %v\n", resp.GetResult())
}

func callPND(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling PND api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{Number: 120})

	if err != nil {
		log.Fatalln("call PND api err %v", err)
	}

	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Fatalln("server finish streaming", err)
		}

		log.Printf("prime number %v", resp.GetResult())
	}
}

func callAverage(c calculatorpb.CalculatorServiceClient)  {
	log.Println("calling average api")
	stream, err := c.Average(context.Background())
	if err != nil {
		log.Fatalln("call average err %v", err)
	}

	listReq := []calculatorpb.AverageRequest{
		calculatorpb.AverageRequest{
			Num1: 5,
		},
		calculatorpb.AverageRequest{
			Num1: 3,
		},
		calculatorpb.AverageRequest{
			Num1: 8,
		},
		calculatorpb.AverageRequest{
			Num1: 10,
		},
	}

	for _, req := range listReq {
		err := stream.Send(&req)
		if err != nil {
			log.Fatalln("send average request err %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalln("receive average request err %v", err)
	}

	log.Printf("average number %+v ", resp)
}

func callMax(c calculatorpb.CalculatorServiceClient)  {
	log.Println("calling max api")
	stream, err := c.Max(context.Background())
	if err != nil {
		log.Fatalln("call average err %v", err)
	}

	waitc := make(chan struct{})
	//defer close(waitc)

	go func() {
		// Send multi request
		listReq := []calculatorpb.MaxRequest{
			calculatorpb.MaxRequest{
				Num: 1,
			},
			calculatorpb.MaxRequest{
				Num: 44,
			},
			calculatorpb.MaxRequest{
				Num: 0,
			},
			calculatorpb.MaxRequest{
				Num: 10,
			},
		}

		for _, req := range listReq {
			err := stream.Send(&req)
			if err != nil {
				log.Fatalln("send max request err %v", err)
				break
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("ending find max api ...")

				break
			}
			if err != nil {
				log.Fatalln("receive find max err %v ", err)
				break
			}

			log.Printf("max: %v", resp.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

