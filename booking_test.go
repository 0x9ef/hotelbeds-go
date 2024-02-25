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
	assert.Equal(t, 6619, resp.Hotels.Hotels[0].Code)
	assert.Equal(t, "2024-04-02", resp.Hotels.CheckIn.String())
	assert.Equal(t, "2024-04-03", resp.Hotels.CheckOut.String())
	assert.Equal(t, 1, resp.Hotels.Total)
	assert.Equal(t, 1, len(resp.Hotels.Hotels))
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
	assert.Equal(t, 6619, resp.Hotel.Code)
	assert.Equal(t, "Millennium  Hotel London Knightsbridge", resp.Hotel.Name)
	assert.Equal(t, 1, len(resp.Hotel.Rooms))
}

func TestConfirmBooking(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Post("/hotel-api/1.2/bookings").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-confirm-booking.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ConfirmBooking(context.TODO(), &ConfirmBookingInput{
		Holder: Holder{
			Name:    "HolderFirstName",
			Surname: "HolderLastName",
		},
		Rooms: []ConfirmBookingRoom{
			{
				RateKey: "20240402|20240403|W|164|6619|TWN.ST|BAR BB FLEX 14|BB||1~1~0||N@06~~21e12c~1630615603~S~~~NOR~5F05A4B7D40E44A170871765642600AADE00000010000000006248118",
				Paxes: []Pax{
					{
						RoomID:  1,
						Type:    "AD",
						Name:    "HolderFirstName",
						Surname: "HolderLastName",
					},
				},
			},
		},
		ClientReference: "IntegrationAgency",
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 712986, resp.Booking.Hotel.Code)
	assert.Equal(t, "INTEGRATIONAGENCY", resp.Booking.ClientReference)
	assert.Equal(t, BookingStatus("CONFIRMED"), resp.Booking.Status)
	assert.Equal(t, 1, len(resp.Booking.Hotel.Rooms))
}
