package dtos

type UserInfo struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Status         string `json:"status"`
	Login_attempts int    `json:"loginAttempts"`
}

type UserForm struct {
	Email string `json:"email" binding:"required"`
	Pass  string `json:"pass" binding:"required"`
}

