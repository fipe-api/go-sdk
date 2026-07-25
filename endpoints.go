package fipe

import (
	"context"
	"fmt"
	"net/url"
)

// References lists the FIPE monthly reference tables, newest first.
func (c *Client) References(ctx context.Context) ([]Reference, error) {
	var out []Reference
	if err := c.get(ctx, "/references", &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Brands lists the brands available for a vehicle type.
func (c *Client) Brands(ctx context.Context, vt VehicleType, opts ...RequestOption) ([]Brand, error) {
	var out []Brand
	path := fmt.Sprintf("/%s/brands", vt)
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// Models lists the models of a brand.
func (c *Client) Models(ctx context.Context, vt VehicleType, brandID string, opts ...RequestOption) ([]Model, error) {
	var out []Model
	path := fmt.Sprintf("/%s/brands/%s/models", vt, url.PathEscape(brandID))
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// Years lists the model-year variants of a model.
func (c *Client) Years(ctx context.Context, vt VehicleType, brandID, modelID string, opts ...RequestOption) ([]Year, error) {
	var out []Year
	path := fmt.Sprintf("/%s/brands/%s/models/%s/years", vt, url.PathEscape(brandID), url.PathEscape(modelID))
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// Vehicle returns the FIPE price details for a brand, model and year.
func (c *Client) Vehicle(ctx context.Context, vt VehicleType, brandID, modelID, yearID string, opts ...RequestOption) (*Vehicle, error) {
	var out Vehicle
	path := fmt.Sprintf("/%s/brands/%s/models/%s/years/%s", vt, url.PathEscape(brandID), url.PathEscape(modelID), url.PathEscape(yearID))
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// YearsByBrand lists all model years available for a brand.
func (c *Client) YearsByBrand(ctx context.Context, vt VehicleType, brandID string, opts ...RequestOption) ([]Year, error) {
	var out []Year
	path := fmt.Sprintf("/%s/brands/%s/years", vt, url.PathEscape(brandID))
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// ModelsByBrandYear lists the models of a brand available for a given year.
func (c *Client) ModelsByBrandYear(ctx context.Context, vt VehicleType, brandID, yearID string, opts ...RequestOption) ([]Model, error) {
	var out []Model
	path := fmt.Sprintf("/%s/brands/%s/years/%s/models", vt, url.PathEscape(brandID), url.PathEscape(yearID))
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// YearsByFipeCode lists the model-year variants of a vehicle by its FIPE code
// (e.g. "005340-6").
func (c *Client) YearsByFipeCode(ctx context.Context, vt VehicleType, fipeCode string, opts ...RequestOption) ([]Year, error) {
	var out []Year
	path := fmt.Sprintf("/%s/%s/years", vt, url.PathEscape(fipeCode))
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return out, nil
}

// VehicleByFipeCode returns the FIPE price details for a vehicle by its FIPE
// code and year.
func (c *Client) VehicleByFipeCode(ctx context.Context, vt VehicleType, fipeCode, yearID string, opts ...RequestOption) (*Vehicle, error) {
	var out Vehicle
	path := fmt.Sprintf("/%s/%s/years/%s", vt, url.PathEscape(fipeCode), url.PathEscape(yearID))
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}

// HistoryByFipeCode returns the vehicle details including its price history
// across reference months.
func (c *Client) HistoryByFipeCode(ctx context.Context, vt VehicleType, fipeCode, yearID string, opts ...RequestOption) (*Vehicle, error) {
	var out Vehicle
	path := fmt.Sprintf("/%s/%s/years/%s/history", vt, url.PathEscape(fipeCode), url.PathEscape(yearID))
	if err := c.get(ctx, path, &out, opts...); err != nil {
		return nil, err
	}
	return &out, nil
}
