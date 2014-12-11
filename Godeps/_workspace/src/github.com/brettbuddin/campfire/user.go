package campfire

type User struct {
	*Connection `json:"-"`

	ID           int    `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	AvatarURL    string `json:"avatar_url"`
	EmailAddress string `json:"email_address"`
	Admin        bool   `json:"admin"`
	ApiAuthToken string `json:"api_auth_token"`
	CreatedAt    string `json:"created_at"`
}

type UserResult struct {
	User *User `json:"user"`
}
