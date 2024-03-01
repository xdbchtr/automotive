package inputters

type CarInput struct {
	Name        string  `json:"name" binding:"required"`
	BrandID     int     `json:"brand_id" binding:"required"`
	YearMade    int     `json:"year_made" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	Stock       int     `json:"stock" binding:"required"`
	TypeID      *int    `json:"type_id"`
	ImageUrl    *string `json:"image_url"`
	Description *string `json:"description"`
}
