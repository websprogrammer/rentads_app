package schemas

type Feedback struct {
	City    string `json:"city"`
	PostId  int    `json:"post_id"`
	Type    int    `json:"type"`
	Message string `json:"message,omitempty"`
}
