package domain

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model             string    `json:"model"`
	Messages          []Message `json:"messages"`
	N                 int       `json:"n,omitempty"`
	Stream            bool      `json:"stream,omitempty"`
	MaxTokens         int       `json:"max_tokens,omitempty"`
	RepetitionPenalty float64   `json:"repetition_penalty,omitempty"`
	UpdateInterval    int       `json:"update_interval,omitempty"`
}
type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
type Brigade struct {
	ID     int64
	Name   string
	Lat    float64
	Lon    float64
	Status string
}
