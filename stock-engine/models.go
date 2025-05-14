package engine

type ScreenerJob struct {
	JobID    string                 `json:"job_id"`
	Rules    map[string]interface{} `json:"rules"`
	Username string                 `json:"username"`
}

type Stock struct {
	Symbol     string
	Indicators map[string]interface{}
}
