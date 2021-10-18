package main

import (
	"context"
	"fmt"
	"gRPC-demo/contact/contactpb"
	"github.com/beego/beego/v2/client/orm"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
)

func init()  {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	connectStr := "root:123456@tcp(127.0.0.1:3306)/contact?charset:utf8"
	err := orm.RegisterDataBase("default", "mysql", connectStr, orm.MaxIdleConnections(100), orm.MaxOpenConnections(100))
	if err != nil {
		log.Panicf("register DB err %v", err)
	}

	orm.RegisterModel(new(ContactInfo))

	err = orm.RunSyncdb("default", true, false)
	if err != nil {
		log.Panicf("run migrate DB err %v", err)
	}

	fmt.Println("Connect db successfully!")
}

type server struct {}

func (s server) Insert(ctx context.Context, req *contactpb.InsertRequest) (*contactpb.InsertResponse, error) {
	log.Printf("calling insert %+v\n", req.Contact)
	ci := ConvertPbContactToContactInfo(req.Contact)

	err := ci.insert()
	if err != nil {
		resp := &contactpb.InsertResponse{
			StatusCode: -1,
			Message:    fmt.Sprintf("insert err %v", err),
		}
		return resp, nil
		//return nil, status.Errorf(codes.InvalidArgument, "Insert %+v err %v", ci, err)
	}

	resp := &contactpb.InsertResponse{
		StatusCode: 1,
		Message:    "OK",
	}

	return resp, nil
}

func main()  {
	lis, err := net.Listen("tcp", "0.0.0.0:50070")
	if err != nil {
		log.Fatalf("err while create listen %v", err)
	}

	s := grpc.NewServer()

	contactpb.RegisterContactServiceServer(s, &server{})

	fmt.Println("calculator is running ...")
	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err while serve %v", err)
	}
}
