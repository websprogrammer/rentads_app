package schemas

type Feedback struct {
	City    string `json:"city"`
	PostId  int    `json:"post_id"`
	Type    int    `json:"type"`
	Message string `json:"message,omitempty"`
}

type FeedbackResponse struct {
	Id           string  `json:"id"`
	City         string  `json:"city"`
	PostId       int     `json:"post_id"`
	FeedbackType int     `json:"feedback_type"`
	Message      *string `json:"message,omitempty"`
}
