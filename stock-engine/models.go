package engine

type ScreenerJob struct {
	JobID  string                 `json:"job_id"`
	Rules  map[string]interface{} `json:"rules"`
	UserID string                 `json:"user_id"`
}

type Stock struct {
	Symbol string
	PE     float64
	Volume int
	Sector string
	// add more attributes as needed
}
