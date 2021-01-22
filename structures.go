package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Advert struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	PostId      uint64        `bson:"PostId"`
	Description string        `bson:"Description"`
	Date        uint64        `bson:"Date"`
	Photos      []Photo       `bson:"Photos"`
	ProfileLink string        `bson:"ProfileLink"`
	ProfileName string        `bson:"ProfileName"`
	City        string        `bson:"City"`
	District    string        `bson:"District"`
	SubDistrict string        `bson:"SubDistrict"`
	Metro       string        `bson:"Metro"`
	RentType    uint          `bson:"RentType"`
	RoomType    uint          `bson:"RoomType"`
}

type DB struct {
	session   *mgo.Session
	adverts   *mgo.Collection
	tokens    *mgo.Collection
	feedbacks *mgo.Collection
}

type Feedback struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	PostId  uint64        `bson:"PostId"`
	City    string        `bson:"City"`
	Type    uint          `bson:"Type"`
	Message string        `bson:"Message"`
}

type Photo struct {
	Average string `json:"average"`
	Small   string `json:"small"`
}

type Results struct {
	Adverts  []Advert `bson:"Adverts"`
	LastDate uint64   `bson:"LastDate"`
}
