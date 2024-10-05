package product

type CreateProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Price 	 float64 `json:"price" binding:"required"`
	Stock 	 int     `json:"stock" binding:"required"`
}

