package responses

type ResponseList struct {
	Count int64       `json:"count"`
	Docs  interface{} `json:"documents"`
}
