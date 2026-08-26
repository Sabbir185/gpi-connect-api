package country

import "database/sql"

type CountryRepository interface {
	Create(country *CountryPayload) (*ResponseCountry, error)
	FindOneByName(name string) (*ResponseCountry, error)
	FindOneByCode(code string) (*ResponseCountry, error)
}

type countryRepository struct {
	db *sql.DB
}

func NewCountryRepository(db *sql.DB) CountryRepository {
	return &countryRepository{
		db: db,
	}
}

// Implicit Implementation
func (r *countryRepository) Create(country *CountryPayload) (*ResponseCountry, error) {
	return nil, nil
}

func (r *countryRepository) FindOneByName(name string) (*ResponseCountry, error) {
	return nil, nil
}

func (r *countryRepository) FindOneByCode(code string) (*ResponseCountry, error) {
	return nil, nil
}
