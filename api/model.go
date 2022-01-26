package api

type Merchant struct {
	Code    string   `json:"code"`
	Members []Member `json:"members"`
}

type Member struct {
	Email string `json:"email"`
}
