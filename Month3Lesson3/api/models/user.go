package models

type User struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Cash     uint   `json:"cash"`
	UserType string `json:"user_type"`
}

type CreateUser struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Cash     uint   `json:"cash"`
	UserType string `json:"user_type"`
}

type UpdateUser struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Cash     uint   `json:"cash"`
}

type UsersResponse struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

type UpdateUserPassword struct {
	ID          string `json:"id"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}
