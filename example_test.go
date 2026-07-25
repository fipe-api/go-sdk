package fipe_test

import (
	"context"
	"errors"
	"fmt"
	"log"

	fipe "github.com/fipe-api/go-sdk"
)

func ExampleNew() {
	// The free tier needs no token; pass one for higher rate limits.
	client := fipe.New(fipe.WithSubscriptionToken("your-token"))

	brands, err := client.Brands(context.Background(), fipe.Cars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(brands[0].Name)
}

func ExampleClient_Brands() {
	client := fipe.New()

	brands, err := client.Brands(context.Background(), fipe.Cars)
	if err != nil {
		log.Fatal(err)
	}
	for _, b := range brands {
		fmt.Println(b.Code, b.Name)
	}
}

func ExampleClient_Vehicle() {
	client := fipe.New()
	ctx := context.Background()

	// Drill down brand -> model -> year -> price.
	models, err := client.Models(ctx, fipe.Cars, "59")
	if err != nil {
		log.Fatal(err)
	}
	years, err := client.Years(ctx, fipe.Cars, "59", models[0].Code)
	if err != nil {
		log.Fatal(err)
	}
	vehicle, err := client.Vehicle(ctx, fipe.Cars, "59", models[0].Code, years[0].Code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(vehicle.Model, vehicle.Price)
}

func ExampleClient_VehicleByFipeCode() {
	client := fipe.New()

	vehicle, err := client.VehicleByFipeCode(context.Background(), fipe.Cars, "005340-6", "2014-3")
	if errors.Is(err, fipe.ErrNotFound) {
		log.Fatal("unknown FIPE code or year")
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println(vehicle.Price, vehicle.ReferenceMonth)
}

func ExampleClient_HistoryByFipeCode() {
	client := fipe.New()

	vehicle, err := client.HistoryByFipeCode(context.Background(), fipe.Cars, "005340-6", "2014-3")
	if err != nil {
		log.Fatal(err)
	}
	for _, h := range vehicle.PriceHistory {
		fmt.Println(h.Month, h.Price)
	}
}

func ExampleWithReference() {
	client := fipe.New()
	ctx := context.Background()

	// Query prices from a past monthly reference table.
	refs, err := client.References(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("latest table:", refs[0].Month)

	brands, err := client.Brands(ctx, fipe.Cars, fipe.WithReference(308))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(brands))
}
