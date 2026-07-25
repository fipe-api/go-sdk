package fipe

// VehicleType identifies the category of vehicle being queried.
type VehicleType string

const (
	Cars        VehicleType = "cars"
	Motorcycles VehicleType = "motorcycles"
	Trucks      VehicleType = "trucks"
)

// Reference is a FIPE monthly reference table. Prices are published per
// reference month; pass a Reference code via WithReference to query
// historical tables.
type Reference struct {
	Code  string `json:"code"`
	Month string `json:"month"`
}

// Brand is a vehicle manufacturer, e.g. "VW - VolksWagen".
type Brand struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Model is a vehicle model within a brand.
type Model struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Year is a model year variant. Code combines year and fuel, e.g. "2022-3".
type Year struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// Vehicle holds the FIPE price details for a specific model year.
type Vehicle struct {
	Brand          string         `json:"brand"`
	CodeFipe       string         `json:"codeFipe"`
	Fuel           string         `json:"fuel"`
	FuelAcronym    string         `json:"fuelAcronym"`
	Model          string         `json:"model"`
	ModelYear      int            `json:"modelYear"`
	Price          string         `json:"price"`
	PriceHistory   []PriceHistory `json:"priceHistory,omitempty"`
	ReferenceMonth string         `json:"referenceMonth"`
	VehicleType    int            `json:"vehicleType"`
}

// PriceHistory is one entry of a vehicle's price across reference months.
type PriceHistory struct {
	Month     string `json:"month"`
	Price     string `json:"price"`
	Reference string `json:"reference"`
}
