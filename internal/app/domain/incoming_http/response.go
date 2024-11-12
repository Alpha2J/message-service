package incoming_http

type EmailVO struct {
	Id      string   `json:"id"`
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}
