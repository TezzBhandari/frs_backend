package frs

type Auth struct {
	ID          int64  `json:"id"`
	AccessToken string `json:"access_token"`
}
