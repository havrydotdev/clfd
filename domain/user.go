package domain

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	Verified     bool   `json:"verified"`
	VerifCode    string `json:"verif_code"`
}

type UpdateUserDTO struct {
	RefreshToken *string `json:"refresh_token"`
	Verified     *bool   `json:"verified"`
}
