package alertmanager

// Payload Struct for the received alert
type Payload struct {
	Alerts      []alert     `json:"alerts"`
	Status      string      `json:"status"`
	GroupLabels groupLabels `json:"groupLabels"`
}

type alert struct {
	Annotations annotations `json:"annotations"`
	Labels      labels      `json:"labels"`
}

type labels struct {
	AlertName string `json:"alertname"`
}

type annotations struct {
	Message string `json:"message"`
}

type groupLabels struct {
	AlertName string `json:"alertname"`
}
