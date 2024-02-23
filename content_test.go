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

func TestListHotels(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/hotels").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-hotels.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListHotels(context.TODO(), &ListHotelsInput{
		Codes: []int{6619, 6613},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, resp.Total, 2)
	assert.Equal(t, len(resp.Hotels), 2)
	assert.Equal(t, resp.Hotels[0].Code, 6613)
	assert.Equal(t, resp.Hotels[1].Code, 6619)
}

func TestGetHotelDetails(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/hotels/6613,6619/details").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-get-hotel-details.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.GetHotelDetails(context.TODO(), []int{6613, 6619}, &GetHotelDetailsInput{})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, len(resp.Hotels), 2)
	assert.Equal(t, resp.Hotels[0].Code, 6613)
	assert.Equal(t, resp.Hotels[1].Code, 6619)
}
