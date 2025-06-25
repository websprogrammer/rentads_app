package schemas

type Notification struct {
	DeviceToken string `json:"token"`
	City        string `json:"city"`
	KeyWords    string `json:"key_words,omitempty"`
	RentType    int    `json:"rent_type,omitempty"`
	RoomType    int    `json:"room_type,omitempty"`
	Districts   string `json:"districts,omitempty"`
	Enabled     int    `json:"enabled"`
}
