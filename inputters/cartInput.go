package inputters

type (
	CartInput struct {
		CarID    int `json:"car_id" binding:"required"`
		Quantity int `json:"quantity"`
		SaledID  int `json:"sales_id" binding:"required"`
	}

	TransactionInput struct {
	}
)
