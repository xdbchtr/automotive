package inputters

type UserInput struct {
	Name        string  `json:"name" validate:"required"`
	Identity    string  `json:"identity" validate:"required"`
	ImageUrl    *string `json:"image_url"`
	Age         *int    `json:"age"`
	TelephoneNo string  `json:"telephone_no" validate:"required"`
	RoleId      uint    `json:"role_id" validate:"required"`
	Username    string  `json:"username" validate:"required"`
	Password    string  `json:"password" validate:"required"`
}

type UserUpdateInput struct {
	Name        string  `json:"name" validate:"required"`
	Identity    string  `json:"identity" validate:"required"`
	ImageUrl    *string `json:"image_url"`
	Age         *int    `json:"age"`
	TelephoneNo string  `json:"telephone_no" validate:"required"`
	RoleId      uint    `json:"role_id" validate:"required"`
}
