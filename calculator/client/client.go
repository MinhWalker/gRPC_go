package main

import (
	"context"
	"gRPC-demo/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main()  {
	cc, err := grpc.Dial("localhost:50069", grpc.WithInsecure())

	if err != nil {
		log.Fatalln("err when dial %v", err)
	}
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	//log.Printf("services client %f", client)
	callSum(client)
}

func callSum(c calculatorpb.CalculatorServiceClient)  {
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
