// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// Amount is an arbitrary-precision decimal.
type Amount decimal.Decimal

func (a Amount) MarshalJSON() ([]byte, error) {
	return []byte(decimal.Decimal(a).StringFixed(2)), nil
}

func (a *Amount) UnmarshalJSON(data []byte) error {
	d, err := decimal.NewFromString(trimUnescapeQuotes(data))
	if err != nil {
		return fmt.Errorf("failed to parse Amount: %w", err)
	}
	*a = Amount(d)
	return nil
}

type Content struct {
	Content      string `json:"content"`
	LanguageCode string `json:"languageCode"`
}

type Coordinate float64

func (c *Coordinate) UnmarshalJSON(data []byte) error {
	f, err := strconv.ParseFloat(trimUnescapeQuotes(data), 64)
	if err != nil {
		return fmt.Errorf("failed to parse Coordinate: %w", err)
	}
	*c = Coordinate(f)
	return nil
}

// Timestamp is time with "2006-01-02 15:04:05.000" layout.
type Timestamp time.Time

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(t).Format("2006-01-02 15:04:05.000") + "\""), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	v, err := time.Parse("2006-01-02 15:04:05.000", trimUnescapeQuotes(data))
	if err != nil {
		return fmt.Errorf("failed to parse Timestamp: %w", err)
	}
	*t = Timestamp(v)
	return nil
}

// TimestampTZ is time with "2006-01-02T15:04:05Z07:00" layout.
type TimestampTZ time.Time

func (t TimestampTZ) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(t).Format(time.RFC3339) + "\""), nil
}

func (t *TimestampTZ) UnmarshalJSON(data []byte) error {
	v, err := time.Parse(time.RFC3339, trimUnescapeQuotes(data))
	if err != nil {
		return fmt.Errorf("failed to parse Timestamp: %w", err)
	}
	*t = TimestampTZ(v)
	return nil
}

// Datetime is time with "2006-01-02" layout.
type Datetime time.Time

func (d Datetime) IsZero() bool {
	return time.Time(d).IsZero()
}

func (d Datetime) String() string {
	return time.Time(d).Format("2006-01-02")
}

func (t Datetime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.String() + "\""), nil
}

func (t *Datetime) UnmarshalJSON(data []byte) error {
	v, err := time.Parse("2006-01-02", trimUnescapeQuotes(data))
	if err != nil {
		return fmt.Errorf("failed to parse DateTime: %w", err)
	}
	*t = Datetime(v)
	return nil
}

type Order int

func (o *Order) UnmarshalJSON(data []byte) error {
	n, err := strconv.Atoi(trimUnescapeQuotes(data))
	if err != nil {
		return fmt.Errorf("failed to parse Order: %w", err)
	}
	*o = Order(n)
	return nil
}

type Distance float64

func (d *Distance) UnmarshalJSON(data []byte) error {
	f, err := strconv.ParseFloat(trimUnescapeQuotes(data), 64)
	if err != nil {
		return fmt.Errorf("failed to parse Distance: %w", err)
	}
	*d = Distance(f)
	return nil
}

type Radius int

func (r Radius) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strconv.Itoa(int(r)) + "\""), nil
}

func (r Radius) Int() int {
	return int(r)
}

type Unit string

const (
	UnitMiles      Unit = "mi"
	UnitKilometers Unit = "km"
)

func (u Unit) String() string {
	return string(u)
}

type FloatRate float64

func (r *FloatRate) UnmarshalJSON(data []byte) error {
	f, err := strconv.ParseFloat(trimUnescapeQuotes(data), 64)
	if err != nil {
		return fmt.Errorf("failed to parse FloatRate: %w", err)
	}
	*r = FloatRate(f)
	return nil
}

func (r FloatRate) MarshalJSON() ([]byte, error) {
	return []byte("\"" + strconv.FormatFloat(float64(r), 'f', 2, 64) + "\""), nil
}

func (r FloatRate) Float() float64 {
	return float64(r)
}

func trimUnescapeQuotes(data []byte) string {
	str, err := strconv.Unquote(string(data))
	if err != nil {
		str = string(data)
	}
	if str[0] == '"' {
		return str[1 : len(str)-1]
	}
	return str
}
