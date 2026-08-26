package country

type CountryPayload struct {
	Name      string `json:"name" validate:"required,min=3,max=5"`
	Code      string `json:"code" validate:"required,min=2"`
	Currency  string `json:"currency" validate:"required,min=2"`
	PhoneCode string `json:"phone_code" validate:"required,min=1"`
	TimeZone  string `json:"time_zone" validate:"required,min=3"`
	Status    string `json:"status"`
}

type ResponseCountry struct {
	ID int64 `json:"id"`
	CountryPayload
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
