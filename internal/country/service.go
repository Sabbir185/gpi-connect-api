package country

type CountryService interface {
	Create(payload *CountryPayload) (*ResponseCountry, error)
}

type countryService struct {
	repo CountryRepository
}

func (s *countryService) Create(payload *CountryPayload) (*ResponseCountry, error) {
	return nil, nil
}
