package schemas

type Advert struct {
	PostId      uint64  `bson:"PostId"`
	Description *string `bson:"Description"`
	Date        uint64  `bson:"Date"`
	Photos      []Photo `bson:"Photos"`
	ProfileLink *string `bson:"ProfileLink"`
	ProfileName *string `bson:"ProfileName"`
	City        *string `bson:"City"`
	District    *string `bson:"District"`
	SubDistrict *string `bson:"SubDistrict"`
	Metro       *string `bson:"Metro"`
	RentType    uint    `bson:"RentType"`
	RoomType    uint    `bson:"RoomType"`
	Price       int64   `bson:"Price"`
}

type Photo struct {
	Average string `json:"average"`
	Small   string `json:"small"`
}
