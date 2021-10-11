package main

import (
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

	err = orm.RunSyncdb("default", false, false)
	if err != nil {
		log.Panicf("run migrate DB err %v", err)
	}

	fmt.Println("Connect db successfully!")
}

type server struct {}

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
