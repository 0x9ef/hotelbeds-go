// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/0x9ef/clientx"
)

type ContentClient interface {
	ListHotels(ctx context.Context, inp *ListHotelsInput) (*ListHotelsResponse, error)
	GetHotelDetails(ctx context.Context, codes []int, inp *GetHotelDetailsInput) (*GetHotelDetailsResponse, error)
}

type (
	Hotel struct {
		Code                 int                `json:"code"`
		Name                 Content            `json:"name"`
		CountryCode          string             `json:"countryCode"`
		StateCode            string             `json:"stateCode"`
		DestinationCode      string             `json:"destinationCode"`
		ZoneCode             int                `json:"zoneCode"`
		Coordinates          Coordinates        `json:"coordinates"`
		CategoryCode         string             `json:"categoryCode"`
		CategoryGroupCode    string             `json:"categoryGroupCode"`
		ChainCode            string             `json:"chainCode"`
		AccommodationType    *HotelAccomodation `json:"accommodationType"`
		AccomodationTypeCode string             `json:"accomodationTypeCode,omitempty"`
		BoardCodes           []string           `json:"boardCodes"`
		SegmentCodes         []int              `json:"segmentCodes"`
		Address              Address            `json:"address"`
		PostalCode           string             `json:"postalCode"`
		City                 Content            `json:"city"`
		Email                string             `json:"email"`
		License              string             `json:"license,omitempty"`
		URL                  string             `json:"web"`
		LastUpdate           Datetime           `json:"lastUpdate"`
		S2C                  string             `json:"S2C"`
		Ranking              int                `json:"ranking"`
		Phones               []Phone            `json:"phones"`
		Rooms                []HotelRoom        `json:"rooms"`
		Facilities           []Facility         `json:"facilities"`
	}

	HotelAccomodation struct {
		Code        string `json:"code"`
		Description string `json:"typeDescription"`
	}

	HotelRoom struct {
		Code               string              `json:"roomCode"`
		IsParentRoom       bool                `json:"isParentRoom"`
		MinPax             int                 `json:"minPax"`
		MaxPax             int                 `json:"maxPax"`
		MinAdults          int                 `json:"minAdults"`
		MaxAdults          int                 `json:"maxAdults"`
		MinChildren        int                 `json:"minChildren"`
		MaxChildren        int                 `json:"maxChildren"`
		Type               string              `json:"roomType"`
		CharacteristicCode string              `json:"characteristicCode"`
		Facilities         []HotelRoomFacility `json:"roomFacilities"`
	}

	HotelRoomFacility struct {
		Code          int  `json:"facilityCode"`
		GroupCode     int  `json:"facilityGroupCode"`
		IndicateLogic bool `json:"indLogic"`
		Number        int  `json:"number"`
		Voucher       bool `json:"voucher"`
	}

	HotelRoomStay struct {
		Type        string              `json:"stayType"`
		Order       Order               `json:"order"`
		Description string              `json:"description"`
		Facilities  []HotelRoomFacility `json:"roomStayFacilities"`
	}

	Terminal struct {
		Code     string   `json:"terminalCode"`
		Distance Distance `json:"distance"`
	}

	Issue struct {
		Code          string    `json:"issueCode"`
		Type          string    `json:"issueType"`
		From          time.Time `json:"dateFrom"`
		To            time.Time `json:"dateTo"`
		Order         Order     `json:"order"`
		IsAlternative bool      `json:"alternative"`
	}

	InterestPoint struct {
		FacilityCode      int      `json:"facilityCode"`
		FacilityGroupCode int      `json:"facilityGroupCode"`
		Order             Order    `json:"order"`
		Name              string   `json:"poiName"`
		Distance          Distance `json:"distance"`
	}

	Facility struct {
		Code          int      `json:"facilityCode"`
		GroupCode     int      `json:"facilityGroupCode"`
		Order         Order    `json:"order"`
		IndicateLogic bool     `json:"indLogic"`
		IndicateFee   bool     `json:"indFee"`
		Number        int      `json:"number"`
		Voucher       bool     `json:"voucher"`
		Distance      Distance `json:"distance"`
	}

	Image struct {
		TypeCode           string `json:"imageTypeCode"`
		Path               string `json:"path"`
		Order              Order  `json:"order"`
		VisualOrder        int    `json:"visualOrder"`
		RoomCode           string `json:"roomCode"`
		RoomType           string `json:"roomType"`
		CharacteristicCode string `json:"characteristicCode"`
	}

	WildCard struct {
		RoomType           string  `json:"roomType"`
		RoomCode           string  `json:"roomCode"`
		CharacteristicCode string  `json:"characteristicCode"`
		Description        Content `json:"hotelRoomDescription"`
	}

	ListHotelsInput struct {
		// Filter for a specific hotel or list of hotels.
		Codes []int `url:"hotelCode"`
		// Filter to limit the results for an specific country.
		CountryCode string `url:"countryCode"`
		// Filter to limit the results for an specific destination.
		DestinationCode string `url:"destinationCode"`
		// Use "webOnly" to include in the response hotels sellable only to websites.
		// Use "notOnSale" to include in the response hotels without rates on sale.
		// By default non of them is included in the response.
		IncludeHotels IncludeHotels `url:"includeHotels"`
		// The list of fields to be received in the response. To retrieve all the fields use ‘all’.
		// If nothing is specified, all fields are returned. See the complete list of available fields in the response.
		Fields []string `url:"fields"`
		// The language code for the language in which you want the descriptions to be returned.
		// If language is not specified, English will be used as default language.
		Language string `url:"language"`
		// The number of the initial record to receive. If nothing is specified, 1 is the default value.
		From int `url:"from"`
		// The number of the final record to receive. If nothing is indicated, 100 is the default value.
		To int `url:"to"`
		// Defines if you want to receive the descriptions in English if the description
		// is not available in the language requested.
		UseSecondaryLanguage *bool `url:"useSecondaryLanguage"`
		// Specifying this parameter limits the results to those modified or added
		// after the date specified. The required format is YYYY-MM-DD.
		LastUpdateTime Datetime `url:"lastUpdateTime"`
		// Sending this parameter as true in the /hotels operations will only return
		// the hotels which possess at least one PMSRoomCode (useful when mapping against the original property codes).
		OnlyPMSRoomCode *bool `url:"PMSRoomCode"`
	}

	ListHotelsResponse struct {
		From   int        `json:"from"`
		To     int        `json:"to"`
		Total  int        `json:"total"`
		Audit  *AuditData `json:"auditData"`
		Hotels []Hotel    `json:"hotels"`
	}

	GetHotelDetailsInput struct {
		Language             string `url:"language"`
		UseSecondaryLanguage *bool  `url:"useSecondaryLanguage"`
	}

	GetHotelDetailsResponse struct {
		Audit  *AuditData `json:"auditData"`
		Hotels []Hotel    `json:"hotels"`
	}
)

type Address struct {
	Content string `json:"content"`
	Street  string `json:"street"`
	Number  string `json:"number,omitempty"`
}

type Coordinates struct {
	Long float64 `json:"longitude"`
	Lat  float64 `json:"latitude"`
}

type Phone struct {
	Number string    `json:"phoneNumber"`
	Type   PhoneType `json:"phoneType"`
}

type PhoneType string

const (
	PhoneTypeHotel      PhoneType = "PHONEHOTEL"
	PhoneTypeBooking    PhoneType = "PHONEBOOKING"
	PhoneTypeFax        PhoneType = "FAXNUMBER"
	PhoneTypeManagement PhoneType = "PHONEMANAGEMENT"
)

type IncludeHotels string

const (
	IncludeHotelsWebOnly   IncludeHotels = "webOnly"
	IncludeHotelsNotOnSale IncludeHotels = "notOnSale"
)

func (ih IncludeHotels) String() string {
	return string(ih)
}

const minFromParam = 1
const maxToParam = 1000

func (inp *ListHotelsInput) Validate() error {
	if inp.From != 0 && inp.From < 1 {
		return errors.New("From param < 1")
	}
	if inp.To != 0 && inp.To > maxToParam {
		return errors.New("To param > 1000")
	}
	if inp.IncludeHotels != "" && inp.IncludeHotels != IncludeHotelsWebOnly && inp.IncludeHotels != IncludeHotelsNotOnSale {
		return errors.New("IncludeHotels invalid, only webOnly, notOnSale supported")
	}
	return nil
}

func (inp ListHotelsInput) Encode(v url.Values) error {
	if len(inp.Codes) != 0 {
		var sb strings.Builder
		for i := range inp.Codes {
			sb.WriteString(strconv.Itoa(inp.Codes[i]) + ",")
		}
		v.Set("codes", sb.String()[:sb.Len()-1])
	}
	if inp.CountryCode != "" {
		v.Set("countryCode", inp.CountryCode)
	}
	if inp.DestinationCode != "" {
		v.Set("destinationCode", inp.DestinationCode)
	}
	if inp.IncludeHotels != "" {
		v.Set("includeHotels", inp.IncludeHotels.String())
	}
	if len(inp.Fields) != 0 {
		v.Set("fields", strings.Join(inp.Fields, ","))
	}
	if inp.Language != "" {
		v.Set("language", inp.Language)
	}
	if inp.From != 0 {
		v.Set("from", strconv.Itoa(inp.From))
	}
	if inp.To != 0 {
		v.Set("to", strconv.Itoa(inp.To))
	}
	if inp.UseSecondaryLanguage != nil {
		v.Set("useSecondaryLanguage", strconv.FormatBool(*inp.UseSecondaryLanguage))
	}
	if !inp.LastUpdateTime.IsZero() {
		v.Set("lastUpdateTime", inp.LastUpdateTime.String())
	}
	if inp.OnlyPMSRoomCode != nil {
		v.Set("PMSRoomCode", strconv.FormatBool(*inp.OnlyPMSRoomCode))
	}
	return nil
}

func (inp GetHotelDetailsInput) Encode(v url.Values) error {
	if inp.Language != "" {
		v.Set("language", inp.Language)
	}
	if inp.UseSecondaryLanguage != nil {
		v.Set("useSecondaryLanguage", strconv.FormatBool(*inp.UseSecondaryLanguage))
	}
	return nil
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/hotelsUsingGET
func (api *API) ListHotels(ctx context.Context, inp *ListHotelsInput) (*ListHotelsResponse, error) {
	if err := inp.Validate(); err != nil {
		return nil, err
	}

	return clientx.NewRequestBuilder[ListHotelsInput, ListHotelsResponse](api.API).
		Get("/hotel-content-api/1.0/hotels", clientx.WithRequestHeaders(api.buildHeaders())).
		WithEncodableQueryParams(inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			b, err := io.ReadAll(resp.Body)
			fmt.Println(string(b), err)
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/hotelWithIdDetailsUsingGET
func (api *API) GetHotelDetails(ctx context.Context, codes []int, inp *GetHotelDetailsInput) (*GetHotelDetailsResponse, error) {
	return clientx.NewRequestBuilder[GetHotelDetailsInput, GetHotelDetailsResponse](api.API).
		Get(fmt.Sprintf("/hotel-content-api/1.0/hotels/%s/details", joinInts[int](codes)), clientx.WithRequestHeaders(api.buildHeaders())).
		WithEncodableQueryParams(inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}
