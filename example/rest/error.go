package rest

type ErrorOutput struct {
	status  int
	Message string `json:"message"`
}

func (o ErrorOutput) Status() int { return o.status }

func (o ErrorOutput) Error() string { return o.Message }
