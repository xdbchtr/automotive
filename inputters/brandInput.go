package inputters

type BrandInput struct {
	Name    string `json:"name" binding:"required"`
	Country string `json:"country" binding:"required"`
}
