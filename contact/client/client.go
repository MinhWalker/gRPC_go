package main

import (
	"gRPC-demo/contact/contactpb"
	"google.golang.org/grpc"
	"log"
)

func main()  {
	cc, err := grpc.Dial("localhost:50070", grpc.WithInsecure())

	if err != nil {
		log.Fatalln("err when dial %v", err)
	}
	defer cc.Close()

	client := contactpb.NewContactServiceClient(cc)

	log.Printf("services client %f", client)
}
