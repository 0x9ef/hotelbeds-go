// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"context"
	"errors"
	"fmt"
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
		Code                 int                  `json:"code"`
		Name                 Content              `json:"name"`
		CountryCode          string               `json:"countryCode"`
		StateCode            string               `json:"stateCode"`
		DestinationCode      string               `json:"destinationCode"`
		ZoneCode             int                  `json:"zoneCode"`
		Coordinates          Coordinates          `json:"coordinates"`
		CategoryCode         string               `json:"categoryCode"`
		CategoryGroupCode    string               `json:"categoryGroupCode"`
		ChainCode            string               `json:"chainCode"`
		AccommodationType    *HotelAccomodation   `json:"accommodationType"`
		AccomodationTypeCode string               `json:"accomodationTypeCode,omitempty"`
		BoardCodes           []string             `json:"boardCodes"`
		SegmentCodes         []int                `json:"segmentCodes"`
		Address              Address              `json:"address"`
		PostalCode           string               `json:"postalCode"`
		City                 Content              `json:"city"`
		Email                string               `json:"email"`
		License              string               `json:"license,omitempty"`
		URL                  string               `json:"web"`
		LastUpdate           Datetime             `json:"lastUpdate"`
		S2C                  string               `json:"S2C"`
		Ranking              int                  `json:"ranking"`
		Phones               []Phone              `json:"phones"`
		Rooms                []HotelRoom          `json:"rooms"`
		Facilities           []HotelFacility      `json:"facilities"`
		Issues               []HotelIssue         `json:"issues"`
		Wildcards            []HotelWildCard      `json:"wildCards"`
		Terminals            []HotelTerminal      `json:"terminals"`
		InterestPoints       []HotelInterestPoint `json:"interestPoints"`
		Images               []HotelImage         `json:"images,omitempty"`
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

	HotelTerminal struct {
		Code     string   `json:"terminalCode"`
		Distance Distance `json:"distance"`
	}

	HotelIssue struct {
		Code          string    `json:"issueCode"`
		Type          string    `json:"issueType"`
		From          time.Time `json:"dateFrom"`
		To            time.Time `json:"dateTo"`
		Order         Order     `json:"order"`
		IsAlternative bool      `json:"alternative"`
	}

	HotelInterestPoint struct {
		FacilityCode      int      `json:"facilityCode"`
		FacilityGroupCode int      `json:"facilityGroupCode"`
		Order             Order    `json:"order"`
		Name              string   `json:"poiName"`
		Distance          Distance `json:"distance"`
	}

	HotelFacility struct {
		Code          int      `json:"facilityCode"`
		GroupCode     int      `json:"facilityGroupCode"`
		Order         Order    `json:"order"`
		IndicateLogic bool     `json:"indLogic"`
		IndicateFee   bool     `json:"indFee"`
		Number        int      `json:"number"`
		Voucher       bool     `json:"voucher"`
		Distance      Distance `json:"distance"`
	}

	HotelImage struct {
		TypeCode           string `json:"imageTypeCode"`
		Path               string `json:"path"`
		Order              Order  `json:"order"`
		VisualOrder        int    `json:"visualOrder"`
		RoomCode           string `json:"roomCode"`
		RoomType           string `json:"roomType"`
		CharacteristicCode string `json:"characteristicCode"`
	}

	HotelWildCard struct {
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

	ListInput struct {
		Fields               []string `url:"fields"`
		Codes                []string `url:"codes"`
		Language             string   `url:"language"`
		From                 int      `url:"from"`
		To                   int      `url:"to"`
		UseSecondaryLanguage bool     `url:"useSecondaryLanguage"`
		LastUpdateTime       Datetime `url:"lastUpdateTime"`
	}

	ListCountriesInput struct {
		ListInput
	}

	ListCountriesResp struct {
		Audit     *AuditData `json:"auditData"`
		Countries []Country  `json:"countries"`
	}

	Country struct {
		Code    string  `json:"code"`
		IsoCode string  `json:"isoCode"`
		States  []State `json:"states"`
	}

	ListDestinationsInput struct {
		ListInput
	}

	ListDestinationsResponse struct {
		Audit        *AuditData    `json:"auditData"`
		Destinations []Destination `json:"destinations"`
	}

	Destination struct {
		Code        string      `json:"code"`
		CountryCode string      `json:"countryCode"`
		Zones       []Zone      `json:"zones"`
		GroupZones  []GroupZone `json:"groupZones"`
	}

	Zone struct {
		Code        int     `json:"zoneCode"`
		Name        string  `json:"name"`
		Description Content `json:"description"`
	}

	GroupZone struct {
		Code string  `json:"groupZoneCode"`
		Name Content `json:"content"`
	}

	State struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Name        string `json:"name"`
		ZoneCode    int    `json:"zoneCode"`
	}

	ListAccommodationsInput struct {
		ListInput
	}

	ListAccommodationsResponse struct {
		Audit          *AuditData `json:"auditData"`
		Accommodations []Accommodation
	}

	Accommodation struct {
		Code            string `json:"code"`
		TypeDescription string `json:"typeDescription"`
	}

	ListBoardsInput struct {
		ListInput
	}

	ListBoardsResponse struct {
		Audit  *AuditData `json:"auditData"`
		Boards []Board    ` json:"boards"`
	}

	Board struct {
		Code             string  `json:"code"`
		Description      Content `json:"description"`
		MultiLingualCode string  `json:"multiLingualCode"`
	}

	ListBoardGroupsInput struct {
		Codes []string `json:"codes"`
		ListInput
	}

	BoardGroup struct {
		Code             string  `json:"code"`
		Description      Content `json:"description"`
		MultiLingualCode string  `json:"multiLingualCode"`
	}

	ListBoardGroupsResponse struct {
		Audit  *AuditData   `json:"auditData"`
		Groups []BoardGroup `json:"boards"`
	}

	ListCategoriesInput struct {
		ListInput
	}

	Category struct {
		Code        string     `json:"code"`
		SimpleCode  SimpleCode `json:"simpleCode"`
		Group       string     `json:"group"`
		Description Content    `json:"description"`
	}

	ListCategoriesResponse struct {
		Audit      *AuditData `json:"audit"`
		Categories []Category `json:"categories"`
	}

	ListClassificationsInput struct {
		ListInput
	}

	Classification struct {
		Code        string  `json:"code"`
		Description Content `json:"description"`
	}

	ListClassificationsResponse struct {
		Audit           *AuditData       `json:"auditData"`
		Classifications []Classification `json:"classifications"`
	}

	ListChainsInput struct {
		ListInput
	}

	Chain struct {
		Code        string  `json:"code"`
		Description Content `json:"description"`
	}

	ListChainsResponse struct {
		Audit  *AuditData `json:"auditData"`
		Chains []Chain    `json:"chains"`
	}

	ListCurrenciesInput struct {
		ListInput
	}

	Currency struct {
		Code        string  `json:"code"`
		Type        string  `json:"currencyType"`
		Description Content `json:"description"`
	}

	ListCurrenciesResponse struct {
		Audit      *AuditData `json:"auditData"`
		Currencies []Currency `json:"currencies"`
	}

	ListFacilitiesInput struct {
		ListInput
	}

	Facility struct {
		Code         int     `json:"code"`
		GroupCode    int     `json:"facilityGroupCode"`
		TopologyCode int     `json:"topologyCode"`
		Description  Content `json:"description"`
	}

	ListFacilitiesResponse struct {
		Audit      *AuditData `json:"auditData"`
		Facilities []Facility `json:"facilities"`
	}

	ListFacilityGroupsInput struct {
		ListInput
	}

	FacilityGroup struct {
		Code        int     `json:"code"`
		Description Content `json:"description"`
	}

	ListFacilityGroupsResponse struct {
		Audit  *AuditData      `json:"auditData"`
		Groups []FacilityGroup `json:"facilityGroups"`
	}

	ListFacilityTypologiesInput struct {
		ListInput
	}

	FacilityTypology struct {
		Code                int  `json:"code"`
		HasNumber           bool `json:"numberFlag"`
		HasLogic            bool `json:"logicFlag"`
		HasDistance         bool `json:"distanceFlag"`
		HasAgeFrom          bool `json:"ageFromFlag"`
		HasAgeTo            bool `json:"ageToFlag"`
		HasTimeFrom         bool `json:"timeFromFlag"`
		HasTimeTo           bool `json:"timeToFlag"`
		HasIndicatesYesOrNo bool `json:"indYesOrNoFlag"`
		HasAmount           bool `json:"amountFlag"`
		HasCurrency         bool `json:"currencyFlag"`
		HasApplicationType  bool `json:"appTypeFlag"`
		HasText             bool `json:"textFlag"`
	}

	ListFacilityTypologiesResponse struct {
		Audit      *AuditData         `json:"auditData"`
		Typologies []FacilityTypology `json:"facilityTypologies"`
	}

	ListImageTypesInput struct {
		ListInput
	}

	ImageType struct {
		Code        string  `json:"code"`
		Description Content `json:"description"`
	}

	ListImageTypesResponse struct {
		Audit *AuditData  `json:"auditData"`
		Types []ImageType `json:"imageTypes"`
	}

	ListIssuesInput struct {
		ListInput
	}

	Issue struct {
		Code          string  `json:"code"`
		Type          string  `json:"type"`
		Description   Content `json:"description"`
		Name          Content `json:"name"`
		IsAlternative bool    `json:"alternative"`
	}

	ListIssuesResponse struct {
		Audit  *AuditData `json:"auditData"`
		Issues []Issue    `json:"issues"`
	}

	ListLanguagesInput struct {
		ListInput
	}

	Language struct {
		Code        string  `json:"code"`
		Name        string  `json:"name"`
		Description Content `json:"description"`
	}

	ListLanguagesResponse struct {
		Audit     *AuditData `json:"auditData"`
		Languages []Language `json:"languages"`
	}

	ListPromotionsInput struct {
		ListInput
	}

	Promotion struct {
		Code        string  `json:"code"`
		Name        Content `json:"name"`
		Description Content `json:"description"`
	}

	ListPromotionsResponse struct {
		Audit      *AuditData  `json:"auditData"`
		Promotions []Promotion `json:"promotions"`
	}

	ListTerminalsInput struct {
		ListInput
	}

	ListRoomsInput struct {
		Codes []string `json:"codes"`
		ListInput
	}

	Room struct {
		Code                      string  `json:"code"`
		Type                      string  `json:"type"`
		TypeDescription           Content `json:"typeDescription"`
		Characteristic            string  `json:"characteristic"`
		CharacteristicDescription Content `json:"characteristicDescription"`
		Description               string  `json:"description"`
		MinPax                    int     `json:"minPax"`
		MaxPax                    int     `json:"maxPax"`
		MinAdults                 int     `json:"minAdults"`
		MaxAdults                 int     `json:"maxAdults"`
		MinChildren               int     `json:"minChildren"`
		MaxChildren               int     `json:"maxChildren"`
	}

	ListRoomsResponse struct {
		Audit *AuditData `json:"auditData"`
		Rooms []Room     `json:"rooms"`
	}

	Terminal struct {
		Code        string  `json:"code"`
		Type        string  `json:"type"`
		Country     string  `json:"country"`
		Description Content `json:"description"`
		Name        Content `json:"name"`
	}

	ListTerminalsResponse struct {
		Audit     *AuditData `json:"auditData"`
		Terminals []Terminal `json:"terminals"`
	}

	ListSegmentsInput struct {
		ListInput
	}

	Segment struct {
		Code        int     `json:"code"`
		Description Content `json:"description"`
	}

	ListSegmentsResponse struct {
		Audit    *AuditData `json:"auditData"`
		Segments []Segment  `json:"segments"`
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

type SimpleCode int

const (
	SimpleCode1Star SimpleCode = iota + 1
	SimpleCode2Stars
	SimpleCode3Stars
	SimpleCode4Stars
	SimpleCode5Stars
)

func (sc SimpleCode) Int() int {
	return int(sc)
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

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/countriesUsingGET
func (api *API) ListCountries(ctx context.Context, inp *ListCountriesInput) (*ListCountriesResp, error) {
	return clientx.NewRequestBuilder[ListCountriesInput, ListCountriesResp](api.API).
		Get("/hotel-content-api/1.0/locations/countries", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/destinationsUsingGET
func (api *API) ListDestinations(ctx context.Context, inp *ListDestinationsInput) (*ListDestinationsResponse, error) {
	return clientx.NewRequestBuilder[ListDestinationsInput, ListDestinationsResponse](api.API).
		Get("/hotel-content-api/1.0/locations/destinations", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

func (api *API) ListAccommodations(ctx context.Context, inp *ListAccommodationsInput) (*ListAccommodationsResponse, error) {
	return clientx.NewRequestBuilder[ListAccommodationsInput, ListAccommodationsResponse](api.API).
		Get("/hotel-content-api/1.0/types/accommodations", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/boardsUsingGET
func (api *API) ListBoards(ctx context.Context, inp *ListBoardsInput) (*ListBoardsResponse, error) {
	return clientx.NewRequestBuilder[ListBoardsInput, ListBoardsResponse](api.API).
		Get("/hotel-content-api/1.0/types/boards", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/boardGroupsUsingGET
func (api *API) ListBoardGroups(ctx context.Context, inp *ListBoardGroupsInput) (*ListBoardGroupsResponse, error) {
	return clientx.NewRequestBuilder[ListBoardGroupsInput, ListBoardGroupsResponse](api.API).
		Get("/hotel-content-api/1.0/types/boardgroups", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/categoriesUsingGET
func (api *API) ListCategories(ctx context.Context, inp *ListCategoriesInput) (*ListCategoriesResponse, error) {
	return clientx.NewRequestBuilder[ListCategoriesInput, ListCategoriesResponse](api.API).
		Get("/hotel-content-api/1.0/types/categories", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/chainsUsingGET
func (api *API) ListChains(ctx context.Context, inp *ListChainsInput) (*ListChainsResponse, error) {
	return clientx.NewRequestBuilder[ListChainsInput, ListChainsResponse](api.API).
		Get("/hotel-content-api/1.0/types/chains", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/classificationsUsingGET
func (api *API) ListClassifications(ctx context.Context, inp *ListClassificationsInput) (*ListClassificationsResponse, error) {
	return clientx.NewRequestBuilder[ListClassificationsInput, ListClassificationsResponse](api.API).
		Get("/hotel-content-api/1.0/types/classifications", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/currenciesUsingGET
func (api *API) ListCurrencies(ctx context.Context, inp *ListCurrenciesInput) (*ListCurrenciesResponse, error) {
	return clientx.NewRequestBuilder[ListCurrenciesInput, ListCurrenciesResponse](api.API).
		Get("/hotel-content-api/1.0/types/currencies", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/facilitiesUsingGET
func (api *API) ListFacilities(ctx context.Context, inp *ListFacilitiesInput) (*ListFacilitiesResponse, error) {
	return clientx.NewRequestBuilder[ListFacilitiesInput, ListFacilitiesResponse](api.API).
		Get("/hotel-content-api/1.0/types/facilities", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/facilitygroupsUsingGET
func (api *API) ListFacilityGroups(ctx context.Context, inp *ListFacilityGroupsInput) (*ListFacilityGroupsResponse, error) {
	return clientx.NewRequestBuilder[ListFacilityGroupsInput, ListFacilityGroupsResponse](api.API).
		Get("/hotel-content-api/1.0/types/facilitygroups", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/facilitytypologiesUsingGET
func (api *API) ListFacilityTypologies(ctx context.Context, inp *ListFacilityTypologiesInput) (*ListFacilityTypologiesResponse, error) {
	return clientx.NewRequestBuilder[ListFacilityTypologiesInput, ListFacilityTypologiesResponse](api.API).
		Get("/hotel-content-api/1.0/types/facilitytypologies", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/imagetypesUsingGET
func (api *API) ListImageTypes(ctx context.Context, inp *ListImageTypesInput) (*ListImageTypesResponse, error) {
	return clientx.NewRequestBuilder[ListImageTypesInput, ListImageTypesResponse](api.API).
		Get("/hotel-content-api/1.0/types/imagetypes", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/issuesUsingGET
func (api *API) ListIssues(ctx context.Context, inp *ListIssuesInput) (*ListIssuesResponse, error) {
	return clientx.NewRequestBuilder[ListIssuesInput, ListIssuesResponse](api.API).
		Get("/hotel-content-api/1.0/types/issues", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/languagesUsingGET
func (api *API) ListLanguages(ctx context.Context, inp *ListLanguagesInput) (*ListLanguagesResponse, error) {
	return clientx.NewRequestBuilder[ListLanguagesInput, ListLanguagesResponse](api.API).
		Get("/hotel-content-api/1.0/types/languages", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/promotionsUsingGET
func (api *API) ListPromotions(ctx context.Context, inp *ListPromotionsInput) (*ListPromotionsResponse, error) {
	return clientx.NewRequestBuilder[ListPromotionsInput, ListPromotionsResponse](api.API).
		Get("/hotel-content-api/1.0/types/promotions", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/roomsUsingGET
func (api *API) ListRooms(ctx context.Context, inp *ListRoomsInput) (*ListRoomsResponse, error) {
	return clientx.NewRequestBuilder[ListRoomsInput, ListRoomsResponse](api.API).
		Get("/hotel-content-api/1.0/types/rooms", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/segmentsUsingGET
func (api *API) ListSegments(ctx context.Context, inp *ListSegmentsInput) (*ListSegmentsResponse, error) {
	return clientx.NewRequestBuilder[ListSegmentsInput, ListSegmentsResponse](api.API).
		Get("/hotel-content-api/1.0/types/segments", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/content-api/api-reference/#operation/terminalsUsingGET
func (api *API) ListTerminals(ctx context.Context, inp *ListTerminalsInput) (*ListTerminalsResponse, error) {
	return clientx.NewRequestBuilder[ListTerminalsInput, ListTerminalsResponse](api.API).
		Get("/hotel-content-api/1.0/types/terminals", clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}
