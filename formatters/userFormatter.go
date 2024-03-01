package formatters

import (
	"automotive/models"
)

type UserFormatter struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Identity    string `json:"identity"`
	ImageUrl    string `json:"image_url"`
	Age         int    `json:"age"`
	TelephoneNo string `json:"telephone_no"`
	Username    string `json:"username"`
	RoleName    string `json:"role_name"`
}

func FormatUser(user models.User) UserFormatter {
	userFormatter := UserFormatter{}

	userFormatter.ID = user.ID
	userFormatter.Name = user.Name
	userFormatter.Identity = user.Identity
	userFormatter.ImageUrl = *user.ImageUrl
	userFormatter.Age = *user.Age
	userFormatter.TelephoneNo = user.TelephoneNo
	userFormatter.Username = user.Username
	userFormatter.RoleName = user.Role.Name

	return userFormatter
}

func FormatUsers(users []models.User) []UserFormatter {
	usersFormatter := []UserFormatter{}

	for _, user := range users {
		userFormatter := FormatUser(user)
		usersFormatter = append(usersFormatter, userFormatter)
	}

	return usersFormatter
}
