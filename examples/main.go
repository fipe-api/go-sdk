// Command examples demonstrates the FIPE SDK against the live API:
// it lists car brands, drills into VW Amarok models and years, and
// prints the current FIPE price. Then it looks up the plate in FIPE_PLATE
// (default ABC1234). Set FIPE_TOKEN to use a subscription token.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	fipe "github.com/fipe-api/go-sdk"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := fipe.New(fipe.WithSubscriptionToken(os.Getenv("FIPE_TOKEN")))

	brands, err := client.Brands(ctx, fipe.Cars)
	if err != nil {
		log.Fatalf("list brands: %v", err)
	}
	fmt.Printf("%d car brands available. First 5:\n", len(brands))
	for _, b := range brands[:min(5, len(brands))] {
		fmt.Printf("  %s: %s\n", b.Code, b.Name)
	}

	var vw fipe.Brand
	for _, b := range brands {
		if strings.Contains(b.Name, "VolksWagen") {
			vw = b
			break
		}
	}
	if vw.Code == "" {
		log.Fatal("VolksWagen not found in brand list")
	}

	models, err := client.Models(ctx, fipe.Cars, vw.Code)
	if err != nil {
		log.Fatalf("list models: %v", err)
	}
	var amarok fipe.Model
	for _, m := range models {
		if strings.Contains(m.Name, "AMAROK") {
			amarok = m
			break
		}
	}
	if amarok.Code == "" {
		log.Fatal("AMAROK not found in model list")
	}
	fmt.Printf("\nModel: %s (code %s)\n", amarok.Name, amarok.Code)

	years, err := client.Years(ctx, fipe.Cars, vw.Code, amarok.Code)
	if err != nil {
		log.Fatalf("list years: %v", err)
	}
	fmt.Printf("Years available: %d\n", len(years))

	vehicle, err := client.Vehicle(ctx, fipe.Cars, vw.Code, amarok.Code, years[0].Code)
	if err != nil {
		log.Fatalf("get vehicle: %v", err)
	}
	fmt.Printf("\n%s %d (%s)\nFIPE code: %s\nPrice: %s (%s)\n",
		vehicle.Model, vehicle.ModelYear, vehicle.Fuel,
		vehicle.CodeFipe, vehicle.Price, vehicle.ReferenceMonth)

	printPlate(ctx, client, envOrDefault("FIPE_PLATE", "ABC1234"))
}

func printPlate(ctx context.Context, client *fipe.Client, plate string) {
	vehicle, err := client.VehicleByPlate(ctx, plate)
	if err != nil {
		log.Fatalf("look up plate: %v", err)
	}
	fmt.Printf("\nPlate %s: %s %s %s (%s, %s/%s)\n",
		vehicle.Plate, vehicle.Brand, vehicle.Model, vehicle.Year,
		vehicle.Color, vehicle.City, vehicle.State)
	for _, p := range vehicle.FipePrices {
		fmt.Printf("  %s %s: %s\n", p.CodeFipe, p.Model, p.Price)
	}
	fmt.Printf("Plate quota: %d of %d left this month\n", vehicle.Quota.Remaining, vehicle.Quota.Limit)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
