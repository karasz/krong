package krong

// User is a Krong user
type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

func NewUser() *User {
	return &User{
		ID:      0,
		Name:    "",
		Email:   "",
		Type:    "",
		Address: "",
	}
}
