# FIPE Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/fipe-api/go-sdk.svg)](https://pkg.go.dev/github.com/fipe-api/go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/fipe-api/go-sdk)](https://goreportcard.com/report/github.com/fipe-api/go-sdk)
[![CI](https://github.com/fipe-api/go-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/fipe-api/go-sdk/actions/workflows/ci.yml)

A zero-dependency Go client for the [FIPE API](https://fipe.api.br) (`/api/v2`), which provides average vehicle prices in the Brazilian market from Fundação Instituto de Pesquisas Econômicas (FIPE). Prices are updated monthly.

## Install

```sh
go get github.com/fipe-api/go-sdk
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	fipe "github.com/fipe-api/go-sdk"
)

func main() {
	client := fipe.New()
	ctx := context.Background()

	brands, err := client.Brands(ctx, fipe.Cars)
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range brands {
		fmt.Println(b.Code, b.Name)
	}
}
```

Vehicle types: `fipe.Cars`, `fipe.Motorcycles`, `fipe.Trucks`.

## Authentication

The free tier works without a token but is rate limited. With a [subscription token](https://fipe.api.parallelum.com.br), you can make more requests per minute:

```go
client := fipe.New(fipe.WithSubscriptionToken("your-token"))
```

Other client options: `fipe.WithHTTPClient(*http.Client)`, `fipe.WithBaseURL(string)`.

## Endpoints

Drill down brand → model → year → price:

```go
brands, _ := client.Brands(ctx, fipe.Cars)                          // GET /cars/brands
models, _ := client.Models(ctx, fipe.Cars, "59")                    // GET /cars/brands/59/models
years, _  := client.Years(ctx, fipe.Cars, "59", "5940")             // GET /cars/brands/59/models/5940/years
vehicle, _ := client.Vehicle(ctx, fipe.Cars, "59", "5940", "2014-3") // GET /cars/brands/59/models/5940/years/2014-3

fmt.Println(vehicle.Model, vehicle.Price) // "AMAROK High.CD 2.0 16V TDI 4x4 Dies. Aut" "R$ 10.000,00"
```

Browse by year:

```go
years, _  := client.YearsByBrand(ctx, fipe.Cars, "59")                    // GET /cars/brands/59/years
models, _ := client.ModelsByBrandYear(ctx, fipe.Cars, "59", "2014-3")     // GET /cars/brands/59/years/2014-3/models
```

Look up by FIPE code:

```go
years, _   := client.YearsByFipeCode(ctx, fipe.Cars, "005340-6")                     // GET /cars/005340-6/years
vehicle, _ := client.VehicleByFipeCode(ctx, fipe.Cars, "005340-6", "2014-3")         // GET /cars/005340-6/years/2014-3
history, _ := client.HistoryByFipeCode(ctx, fipe.Cars, "005340-6", "2014-3")         // GET /cars/005340-6/years/2014-3/history

for _, h := range history.PriceHistory {
	fmt.Println(h.Month, h.Price)
}
```

### Reference months

Prices are published per monthly reference table. Every endpoint accepts `fipe.WithReference` to query a past table:

```go
refs, _ := client.References(ctx) // GET /references — e.g. {Code: "308", Month: "abril de 2024"}

brands, _ := client.Brands(ctx, fipe.Cars, fipe.WithReference(308))
```

## Error handling

Non-2xx responses return an `*fipe.APIError` carrying the status code and body. 404 and 429 also match sentinel errors:

```go
vehicle, err := client.Vehicle(ctx, fipe.Cars, "59", "5940", "1900-1")
switch {
case errors.Is(err, fipe.ErrNotFound):
	// unknown brand/model/year
case errors.Is(err, fipe.ErrTooManyRequests):
	// rate limited — back off or use a subscription token
case err != nil:
	var apiErr *fipe.APIError
	if errors.As(err, &apiErr) {
		log.Printf("API returned %d: %s", apiErr.StatusCode, apiErr.Body)
	}
}
```

## Example program

A runnable example that lists brands and prints an Amarok price from the live API:

```sh
go run ./examples
```

## License

MIT
