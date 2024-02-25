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

func TestListCountries(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/locations/countries").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-locations-countries.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListCountries(context.TODO(), &ListCountriesInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Countries))
	assert.Equal(t, resp.Countries[0].Code, "AD")
	assert.Equal(t, resp.Countries[1].Code, "AE")
	assert.Equal(t, resp.Countries[0].IsoCode, "AD")
	assert.Equal(t, resp.Countries[1].IsoCode, "AE")
}

func TestListDestinations(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/locations/destinations").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-locations-destinations.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListDestinations(context.TODO(), &ListDestinationsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Destinations))
	assert.Equal(t, resp.Destinations[0].Code, "01H")
	assert.Equal(t, resp.Destinations[1].Code, "01N")
	assert.Equal(t, resp.Destinations[0].Zones[0].Code, 2)
	assert.Equal(t, resp.Destinations[1].Zones[0].Code, 3)
}

func TestListAccommodations(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/accommodations").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-accommodations.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListAccommodations(context.TODO(), &ListAccommodationsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Accommodations))
	assert.Equal(t, resp.Accommodations[0].Code, "G")
	assert.Equal(t, resp.Accommodations[1].Code, "Q")
}

func TestListBoards(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/boards").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-boards.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListBoards(context.TODO(), &ListBoardsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Boards))
	assert.Equal(t, resp.Boards[0].Code, "AB")
	assert.Equal(t, resp.Boards[1].Code, "AI")
}

func TestListBoardGroups(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/boardgroups").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-board-groups.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListBoardGroups(context.TODO(), &ListBoardGroupsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Groups))
	assert.Equal(t, resp.Groups[0].Code, "AB")
	assert.Equal(t, resp.Groups[1].Code, "AI")
}

func TestListCategories(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/categories").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-categories.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListCategories(context.TODO(), &ListCategoriesInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Categories))
	assert.Equal(t, resp.Categories[0].Code, "1EST")
	assert.Equal(t, resp.Categories[1].Code, "1LL")
}

func TestListChains(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/chains").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-chains.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListChains(context.TODO(), &ListChainsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Chains))
	assert.Equal(t, resp.Chains[0].Code, "007")
	assert.Equal(t, resp.Chains[1].Code, "13CO")
}

func TestListClassifications(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/classifications").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-classifications.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListClassifications(context.TODO(), &ListClassificationsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Classifications))
	assert.Equal(t, resp.Classifications[0].Code, "AUS")
	assert.Equal(t, resp.Classifications[1].Code, "BAL")
}

func TestListCurrencies(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/currencies").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-currencies.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListCurrencies(context.TODO(), &ListCurrenciesInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Currencies))
	assert.Equal(t, resp.Currencies[0].Code, "AED")
	assert.Equal(t, resp.Currencies[1].Code, "AFA")
}

func TestListFacilities(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/facilities").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-facilities.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListFacilities(context.TODO(), &ListFacilitiesInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Facilities))
	assert.Equal(t, resp.Facilities[0].Code, 1)
	assert.Equal(t, resp.Facilities[0].GroupCode, 61)
	assert.Equal(t, resp.Facilities[1].Code, 1)
	assert.Equal(t, resp.Facilities[1].GroupCode, 62)
}

func TestListFacilityGroups(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/facilitygroups").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-facllity-groups.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListFacilityGroups(context.TODO(), &ListFacilityGroupsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Groups))
	assert.Equal(t, resp.Groups[0].Code, 10)
	assert.Equal(t, resp.Groups[1].Code, 100)
}

func TestListFacilityTypologies(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/facilitytypologies").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-facility-typologies.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListFacilityTypologies(context.TODO(), &ListFacilityTypologiesInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Typologies))
	assert.Equal(t, resp.Typologies[0].Code, 12)
	assert.Equal(t, resp.Typologies[1].Code, 14)
}

func TestListIssues(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/issues").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-issues.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListIssues(context.TODO(), &ListIssuesInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Issues))
	assert.Equal(t, resp.Issues[0].Code, "PROHIBITED")
	assert.Equal(t, resp.Issues[1].Code, "ARRIVALTIME")
}

func TestListLanguages(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/languages").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-languages.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListLanguages(context.TODO(), &ListLanguagesInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Languages))
	assert.Equal(t, resp.Languages[0].Code, "ALE")
	assert.Equal(t, resp.Languages[1].Code, "ARA")
}

func TestListPromotions(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/promotions").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-promotions.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListPromotions(context.TODO(), &ListPromotionsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Promotions))
	assert.Equal(t, resp.Promotions[0].Code, "013")
	assert.Equal(t, resp.Promotions[1].Code, "014")
}

func TestListRooms(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/rooms").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-rooms.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListRooms(context.TODO(), &ListRoomsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Rooms))
	assert.Equal(t, resp.Rooms[0].Code, "APT.0E")
	assert.Equal(t, resp.Rooms[1].Code, "APT.1B")
	assert.Equal(t, resp.Rooms[0].Type, "APT")
	assert.Equal(t, resp.Rooms[1].Type, "APT")
}

func TestListSegments(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/segments").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-segments.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListSegments(context.TODO(), &ListSegmentsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Segments))
	assert.Equal(t, resp.Segments[0].Code, 100)
	assert.Equal(t, resp.Segments[1].Code, 102)
}

func TestListTerminals(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.test.hotelbeds.com").
		Get("/hotel-content-api/1.0/types/terminals").
		Reply(200).
		SetHeader("X-Ratelimit-Limit: 50000", "100").
		SetHeader("X-Ratelimit-Remaining", "100").
		File("fixtures/200-list-types-terminals.json")

	client := New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
	resp, err := client.ListTerminals(context.TODO(), &ListTerminalsInput{
		ListInput: ListInput{
			From: 1,
			To:   2,
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Terminals))
	assert.Equal(t, resp.Terminals[0].Code, "AAE")
	assert.Equal(t, resp.Terminals[1].Code, "AAGT")
	assert.Equal(t, resp.Terminals[0].Type, "A")
	assert.Equal(t, resp.Terminals[1].Type, "A")
}
