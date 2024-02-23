// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestListAvailableHotels(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Post("/hotel-api/1.0/hotels").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-available-hotels.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListAvailableHotels(context.TODO(), &ListAvailableHotelsInput{
		Stay: Stay{
			CheckIn:  "2024-04-02",
			CheckOut: "2024-04-03",
		},
		Occupancies: []Occupancy{
			{
				Rooms:  1,
				Adults: 1,
			},
		},
		Hotels: FilterHotel{
			HotelCodes: []int{6619, 6613},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Hotels.Hotels[0].Code, 6619)
	assert.Equal(t, resp.Hotels.CheckIn.String(), "2024-04-02")
	assert.Equal(t, resp.Hotels.CheckOut.String(), "2024-04-03")
	assert.Equal(t, resp.Hotels.Total, 1)
	assert.Equal(t, len(resp.Hotels.Hotels), 1)
}

func TestListCheckRates(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Post("/hotel-api/1.0/checkrates").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-checkrates.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListCheckRates(context.TODO(), &ListCheckRatesInput{
		Rooms: []ListCheckRatesRoom{
			{
				RateKey: "20240402|20240403|W|164|6619|TWN.ST|BAR BB FLEX 14|BB||1~1~0||N@06~~21e12c~1630615603~S~~~NOR~5F05A4B7D40E44A170871765642600AADE00000010000000006248118",
			},
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Hotel.Code, 6619)
	assert.Equal(t, resp.Hotel.Name, "Millennium  Hotel London Knightsbridge")
	assert.Equal(t, len(resp.Hotel.Rooms), 1)
}
