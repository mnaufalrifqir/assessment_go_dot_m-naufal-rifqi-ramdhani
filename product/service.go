package product

type Service interface {
	GetAllProduct() ([]Product, error)
	GetProduct(input GetProductDetailInput) (Product, error)
	CreateProduct(input CreateProductInput) (Product, error)
	UpdateProduct(inputID GetProductDetailInput, inputData CreateProductInput) (Product, error)
	DeleteProduct(input GetProductDetailInput) error
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

func (s *service) GetProduct(input GetProductDetailInput) (Product, error) {
	product, err := s.repository.FindByID(input.ID)
	if err != nil {
		return product, err
	}

	return product, nil
}

func (s *service) GetAllProduct() ([]Product, error) {
	products, err := s.repository.FindAll()
	if err != nil {
		return products, err
	}

	return products, nil
}

func (s *service) UpdateProduct(inputID GetProductDetailInput, inputData CreateProductInput) (Product, error) {
	product, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return product, err
	}

	product.Name = inputData.Name
	product.Price = inputData.Price
	product.Stock = inputData.Stock

	updatedProduct, err := s.repository.Update(product)
	if err != nil {
		return updatedProduct, err
	}

	return updatedProduct, nil
}

func (s *service) DeleteProduct(input GetProductDetailInput) error {
	_, err := s.repository.DeleteByID(input.ID)
	if err != nil {
		return err
	}

	return nil
}

