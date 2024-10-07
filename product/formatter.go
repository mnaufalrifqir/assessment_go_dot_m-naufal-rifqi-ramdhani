package product

type ProductFormatter struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func FormatProduct(product Product) ProductFormatter {
	formatter := ProductFormatter{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
		Stock: product.Stock,
	}

	return formatter
}

func FormatProducts(products []Product) []ProductFormatter {
	productsFormatter := []ProductFormatter{}

	for _, product := range products {
		formatter := FormatProduct(product)
		productsFormatter = append(productsFormatter, formatter)
	}

	return productsFormatter
}
