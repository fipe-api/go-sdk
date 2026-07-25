package fipe

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc, opts ...Option) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return New(append([]Option{WithBaseURL(srv.URL)}, opts...)...)
}

func TestEndpoints(t *testing.T) {
	vehicleJSON := `{
		"brand": "VW - VolksWagen",
		"codeFipe": "005340-6",
		"fuel": "Diesel",
		"fuelAcronym": "D",
		"model": "AMAROK High.CD 2.0 16V TDI 4x4 Dies. Aut",
		"modelYear": 2014,
		"price": "R$ 10.000,00",
		"referenceMonth": "abril de 2024",
		"vehicleType": 1
	}`
	wantVehicle := &Vehicle{
		Brand:          "VW - VolksWagen",
		CodeFipe:       "005340-6",
		Fuel:           "Diesel",
		FuelAcronym:    "D",
		Model:          "AMAROK High.CD 2.0 16V TDI 4x4 Dies. Aut",
		ModelYear:      2014,
		Price:          "R$ 10.000,00",
		ReferenceMonth: "abril de 2024",
		VehicleType:    1,
	}

	tests := []struct {
		name      string
		call      func(*Client) (any, error)
		wantPath  string
		wantQuery string
		body      string
		want      any
	}{
		{
			name:     "References",
			call:     func(c *Client) (any, error) { return c.References(context.Background()) },
			wantPath: "/references",
			body:     `[{"code": "308", "month": "abril de 2024"}]`,
			want:     []Reference{{Code: "308", Month: "abril de 2024"}},
		},
		{
			name:     "Brands",
			call:     func(c *Client) (any, error) { return c.Brands(context.Background(), Cars) },
			wantPath: "/cars/brands",
			body:     `[{"code": "23", "name": "VW - VolksWagen"}]`,
			want:     []Brand{{Code: "23", Name: "VW - VolksWagen"}},
		},
		{
			name:      "Brands with reference",
			call:      func(c *Client) (any, error) { return c.Brands(context.Background(), Trucks, WithReference(308)) },
			wantPath:  "/trucks/brands",
			wantQuery: "reference=308",
			body:      `[]`,
			want:      []Brand{},
		},
		{
			name:     "Models",
			call:     func(c *Client) (any, error) { return c.Models(context.Background(), Cars, "23") },
			wantPath: "/cars/brands/23/models",
			body:     `[{"code": "5585", "name": "AMAROK CD2.0 16V/S CD2.0 16V TDI 4x2 Die"}]`,
			want:     []Model{{Code: "5585", Name: "AMAROK CD2.0 16V/S CD2.0 16V TDI 4x2 Die"}},
		},
		{
			name:     "Years",
			call:     func(c *Client) (any, error) { return c.Years(context.Background(), Cars, "23", "5585") },
			wantPath: "/cars/brands/23/models/5585/years",
			body:     `[{"code": "2022-3", "name": "2022 Diesel"}]`,
			want:     []Year{{Code: "2022-3", Name: "2022 Diesel"}},
		},
		{
			name:     "Vehicle",
			call:     func(c *Client) (any, error) { return c.Vehicle(context.Background(), Cars, "23", "5585", "2022-3") },
			wantPath: "/cars/brands/23/models/5585/years/2022-3",
			body:     vehicleJSON,
			want:     wantVehicle,
		},
		{
			name:     "YearsByBrand",
			call:     func(c *Client) (any, error) { return c.YearsByBrand(context.Background(), Motorcycles, "23") },
			wantPath: "/motorcycles/brands/23/years",
			body:     `[{"code": "2022-3", "name": "2022 Diesel"}]`,
			want:     []Year{{Code: "2022-3", Name: "2022 Diesel"}},
		},
		{
			name:     "ModelsByBrandYear",
			call:     func(c *Client) (any, error) { return c.ModelsByBrandYear(context.Background(), Cars, "23", "2022-3") },
			wantPath: "/cars/brands/23/years/2022-3/models",
			body:     `[{"code": "5585", "name": "AMAROK"}]`,
			want:     []Model{{Code: "5585", Name: "AMAROK"}},
		},
		{
			name:     "YearsByFipeCode",
			call:     func(c *Client) (any, error) { return c.YearsByFipeCode(context.Background(), Cars, "005340-6") },
			wantPath: "/cars/005340-6/years",
			body:     `[{"code": "2022-3", "name": "2022 Diesel"}]`,
			want:     []Year{{Code: "2022-3", Name: "2022 Diesel"}},
		},
		{
			name: "VehicleByFipeCode",
			call: func(c *Client) (any, error) {
				return c.VehicleByFipeCode(context.Background(), Cars, "005340-6", "2022-3")
			},
			wantPath: "/cars/005340-6/years/2022-3",
			body:     vehicleJSON,
			want:     wantVehicle,
		},
		{
			name: "HistoryByFipeCode",
			call: func(c *Client) (any, error) {
				return c.HistoryByFipeCode(context.Background(), Cars, "005340-6", "2022-3")
			},
			wantPath: "/cars/005340-6/years/2022-3/history",
			body:     vehicleJSON,
			want:     wantVehicle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.wantPath {
					t.Errorf("path = %q, want %q", r.URL.Path, tt.wantPath)
				}
				if r.URL.RawQuery != tt.wantQuery {
					t.Errorf("query = %q, want %q", r.URL.RawQuery, tt.wantQuery)
				}
				if got := r.Header.Get("Accept"); got != "application/json" {
					t.Errorf("Accept header = %q, want %q", got, "application/json")
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(tt.body))
			})

			got, err := tt.call(c)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("result = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSubscriptionToken(t *testing.T) {
	t.Run("set", func(t *testing.T) {
		c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if got := r.Header.Get("X-Subscription-Token"); got != "secret" {
				t.Errorf("X-Subscription-Token = %q, want %q", got, "secret")
			}
			w.Write([]byte(`[]`))
		}, WithSubscriptionToken("secret"))
		if _, err := c.References(context.Background()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("unset", func(t *testing.T) {
		c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
			if _, ok := r.Header["X-Subscription-Token"]; ok {
				t.Error("X-Subscription-Token header should not be sent")
			}
			w.Write([]byte(`[]`))
		})
		if _, err := c.References(context.Background()); err != nil {
			t.Fatal(err)
		}
	})
}

func TestErrorMapping(t *testing.T) {
	tests := []struct {
		status   int
		sentinel error
	}{
		{http.StatusNotFound, ErrNotFound},
		{http.StatusTooManyRequests, ErrTooManyRequests},
		{http.StatusInternalServerError, nil},
	}

	for _, tt := range tests {
		t.Run(http.StatusText(tt.status), func(t *testing.T) {
			c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, http.StatusText(tt.status), tt.status)
			})

			_, err := c.Brands(context.Background(), Cars)
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var apiErr *APIError
			if !errors.As(err, &apiErr) {
				t.Fatalf("error %v is not an *APIError", err)
			}
			if apiErr.StatusCode != tt.status {
				t.Errorf("StatusCode = %d, want %d", apiErr.StatusCode, tt.status)
			}
			if tt.sentinel != nil && !errors.Is(err, tt.sentinel) {
				t.Errorf("errors.Is(err, %v) = false, want true", tt.sentinel)
			}
		})
	}
}

func TestContextCancellation(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := c.References(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestInvalidJSON(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not json`))
	})

	if _, err := c.References(context.Background()); err == nil {
		t.Fatal("expected decode error, got nil")
	}
}
