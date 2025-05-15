package engine

type Rule struct {
	Type               string      `json:"type"`
	Condition          string      `json:"condition"`
	ComparisonType     string      `json:"comparisonType"`
	ComparisonValue    interface{} `json:"comparisonValue"`
	TechnicalIndicator string      `json:"technicalIndicator"`
}

type ScreenerJob struct {
	JobID    string `json:"job_id"`
	Rules    []Rule `json:"rules"`
	Username string `json:"username"`
}

type Stock struct {
	Symbol     string
	Indicators map[string]interface{}
}
