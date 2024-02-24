package domain

type User struct {
	ID int `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
	Verified bool `json:"verified"`
	VerifCode string `json:"verif_code"`
	RefreshToken string `json:"refresh_token"`
}
