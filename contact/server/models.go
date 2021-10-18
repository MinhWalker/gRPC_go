package main

import (
	"gRPC-demo/contact/contactpb"
	"github.com/beego/beego/v2/client/orm"
	"log"
)

type ContactInfo struct {
	PhoneNumber string `orm:"size(15);pk"`
	Name        string
	Address     string `orm:"type(text)"`
}

func ConvertPbContactToContactInfo(contact *contactpb.Contact) *ContactInfo {
	return &ContactInfo{
		PhoneNumber: contact.PhoneNumber,
		Name:        contact.Name,
		Address:     contact.Address,
	}
}

func (c *ContactInfo) insert() error {
	o := orm.NewOrm()

	_, err := o.Insert(c)
	if err != nil {
		log.Printf("Insert contact %+v err %v\n", c, err)
		return err
	}

	log.Printf("Insert %+v successfully\n", c)
	return nil
}
