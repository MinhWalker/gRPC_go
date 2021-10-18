package main

import (
	"context"
	"gRPC-demo/contact/contactpb"
	"google.golang.org/grpc"
	"log"
)

func main() {
	cc, err := grpc.Dial("localhost:50070", grpc.WithInsecure())

	if err != nil {
		log.Fatalln("err when dial %v", err)
	}
	defer cc.Close()

	client := contactpb.NewContactServiceClient(cc)

	insertContact(client, "0987232332", "Contact2", "address 2")
}

func insertContact(cli contactpb.ContactServiceClient, phone, name, addr string) {
	req := &contactpb.InsertRequest{
		Contact: &contactpb.Contact{
			PhoneNumber: phone,
			Name:        name,
			Address:     addr,
		}}
	resp, err := cli.Insert(context.Background(), req)
	if err != nil {
		log.Printf("call insert err %v\n", err)
		return
	}

	log.Printf("insert response %+v", resp)
}
