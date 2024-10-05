package product

type Service interface {
	CreateProduct(input CreateProductInput) (Product, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateProduct(input CreateProductInput) (Product, error) {
	product := Product{}
	product.Name = input.Name
	product.Price = input.Price
	product.Stock = input.Stock

	newProduct, err := s.repository.Save(product)
	if err != nil {
		return newProduct, err
	}

	return newProduct, nil
}
