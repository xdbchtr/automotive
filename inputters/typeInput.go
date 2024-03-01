package inputters

type TypeInput struct {
	Name string `json:"name" binding:"required"`
}
