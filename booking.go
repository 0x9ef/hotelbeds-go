// Copyright (c) 2024 0x9ef. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.
package hotelbeds

import (
	"context"
	"errors"
	"net/http"

	"github.com/0x9ef/clientx"
)

type BookingClient interface {
	ListAvailableHotels(ctx context.Context, inp *ListAvailableHotelsInput) (*ListAvailableHotelsResponse, error)
	ListCheckRates(ctx context.Context, inp *ListCheckRatesInput) (*ListCheckRatesResponse, error)
	GetBooking(ctx context.Context, id string) (*GetBookingResponse, error)
	ListBookings(ctx context.Context, inp *CancelBookingInput) (*CancelBookingResponse, error)
	ConfirmBooking(ctx context.Context, inp *ConfirmBookingInput) (*ConfirmBookingResponse, error)
	ChangeBooking(ctx context.Context, id string, inp *ChangeBookingInput) (*ChangeBookingResponse, error)
	CancelBooking(ctx context.Context, id string, inp *CancelBookingInput) (*CancelBookingResponse, error)
}

type (
	ListAvailableHotelsInput struct {
		Stay        Stay          `json:"stay"`
		Occupancies []Occupancy   `json:"occupancies"`
		Keywords    []Keyword     `json:"keywords,omitempty"`
		Geolocation *Geolocation  `json:"geolocation,omitempty"`
		Filter      *Filter       `json:"filter,omitempty"`
		Boards      *FilterBoards `json:"boards,omitempty"`
		Rooms       *FilterRooms  `json:"rooms,omitempty"`
		Hotels      FilterHotel   `json:"hotels"`
		// Displays price breakdown per each day of the hotel stay.
		WithDailyRate bool `json:"dailyRate"`
		// Hotelbeds Group client source market.
		SourceMarket string `json:"sourceMarket,omitempty"`
		// Defines the platform for multiclient developer platforms.
		Platform int `json:"platform,omitempty"`
		// Language code that defines the language of the response.
		// English will be used by default if this field is not informed.
		Language string `json:"language,omitempty"`
		// Filter for accomodation type.
		Accomodations []string `json:"accomodations,omitempty"`
	}

	AvailableHotel struct {
		Code            int                  `json:"code"`
		Name            string               `json:"name"`
		CategoryCode    string               `json:"categoryCode"`
		CategoryName    string               `json:"categoryName"`
		DestinationCode string               `json:"destinationCode"`
		DestinationName string               `json:"destinationName"`
		ZoneCode        int                  `json:"zoneCode"`
		ZoneName        string               `json:"zoneName"`
		Latitude        Coordinate           `json:"latitude"`
		Longitude       Coordinate           `json:"longitude"`
		Rooms           []AvailableHotelRoom `json:"rooms"`
		MinRate         FloatRate            `json:"minRate"`
		MaxRate         FloatRate            `json:"maxRate"`
		Currency        string               `json:"currency"`
	}

	AvailableHotelRoom struct {
		Code  string `json:"code"`
		Name  string `json:"name"`
		Rates []Rate `json:"rates"`
	}

	Rate struct {
		RateKey              string               `json:"rateKey"`
		RateClass            string               `json:"rateClass"`
		RateType             string               `json:"rateType"`
		Net                  Amount               `json:"net"`
		Selling              Amount               `json:"sellingRate"`
		Allotment            int                  `json:"allotment"`
		RateCommentdsID      string               `json:"rateCommentsId,omitempty"`
		PaymentType          PaymentType          `json:"paymentType"`
		Packaging            bool                 `json:"packaging"`
		BoardCode            string               `json:"boardCode"`
		BoardName            string               `json:"boardName"`
		CancellationPolicies []CancellationPolicy `json:"cancellationPolicies"`
		Rooms                int                  `json:"rooms"`
		Adults               int                  `json:"adults"`
		Children             int                  `json:"children"`
		Offers               []Offer              `json:"offers,omitempty"`
	}

	ShiftRate struct {
		RateKey   string   `json:"rateKey"`
		RateClass string   `json:"rateClass"`
		RateType  string   `json:"rateType"`
		Net       Amount   `json:"net"`
		Selling   Amount   `json:"sellingRate"`
		Allotment int      `json:"allotment"`
		CheckIn   Datetime `json:"checkIn"`
		CheckOut  Datetime `json:"checkOut"`
	}

	CancellationPolicy struct {
		Amount Amount      `json:"amount"`
		From   TimestampTZ `json:"from"`
	}

	Offer struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Amount Amount `json:"amount"`
	}

	ListAvailableHotelsResponse struct {
		Audit  *AuditData `json:"auditData"`
		Hotels struct {
			CheckIn  Datetime         `json:"checkIn"`
			CheckOut Datetime         `json:"checkOut"`
			Total    int              `json:"total"`
			Hotels   []AvailableHotel `json:"hotels"`
		} `json:"hotels"`
	}

	// CheckRates.
	ListCheckRatesInput struct {
		// Parameter to add or remove the upSelling options node to the response.
		Upselling bool `json:"upselling"`
		// When true, it will add either the percent or the numberOfnights to the cancellation policies.
		ExpandCXL bool `json:"expandCXL"`
		// Language code that defines the language of the response.
		// English will be used by default if this field is not informed.
		Language string `json:"language"`
		// List of rooms to be checked/valuated.
		Rooms []ListCheckRatesRoom `json:"rooms"`
	}

	ListCheckRatesRoom struct {
		// Internal key that represents a combination of room type, category, board and occupancy.
		// Is returned in Availability and used to valuate a rate and confirm a booking.
		RateKey string `json:"rateKey"`
		// Data of the passengers assigned to this room.
		Paxes []Pax `json:"paxes"`
	}

	CheckRateHotel struct {
		Code                 int                `json:"code"`
		Name                 string             `json:"name"`
		CategoryCode         string             `json:"categoryCode"`
		CategoryName         string             `json:"categoryName"`
		DestinationCode      string             `json:"destinationCode"`
		DestinationName      string             `json:"destinationName"`
		ZoneCode             int                `json:"zoneCode"`
		ZoneName             string             `json:"zoneName"`
		Latitude             Coordinate         `json:"latitude"`
		Longitude            Coordinate         `json:"longitude"`
		Rooms                []CheckRateRoom    `json:"rooms"`
		MinRate              *FloatRate         `json:"minRate"`
		MaxRate              *FloatRate         `json:"maxRate"`
		Currency             string             `json:"currency"`
		CheckIn              string             `json:"checkIn"`
		CheckOut             string             `json:"checkOut"`
		TotalNet             Amount             `json:"totalNet"`
		PaymentDataRequired  bool               `json:"paymentDataRequired"`
		ModificationPolicies ModificationPolicy `json:"modificationPolicies"`
	}

	CheckRateRoom struct {
		Code  string      `json:"code"`
		Name  string      `json:"name"`
		Rates []CheckRate `json:"rates"`
	}

	CheckRate struct {
		Rate
		RateComments string `json:"rateComments"`
		Taxes        *struct {
			Taxes []Tax `json:"taxes"`
		} `json:"taxes,omitempty"`
		BreakDown *BreakDown `json:"breakDown,omitempty"`
	}

	Tax struct {
		Included       bool   `json:"included"`
		Amount         Amount `json:"amount"`
		Currency       string `json:"currency"`
		ClientAmount   Amount `json:"clientAmount"`
		ClientCurrency string `json:"clientCurrency"`
	}

	BreakDown struct {
		Discounts []Discount `json:"rateDiscounts"`
	}

	Discount struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Amount Amount `json:"amount"`
	}

	ModificationPolicy struct {
		IsCancellationAllowed bool `json:"cancellation"`
		IsModificationAllowed bool `json:"modification"`
	}

	ListCheckRatesResponse struct {
		Audit *AuditData      `json:"auditData"`
		Hotel *CheckRateHotel `json:"hotel"`
	}

	// Reservations.
	ConfirmBookingInput struct {
		Holder  Holder       `json:"holder"`
		Payment *PaymentData `json:"paymentData,omitempty"`
		// Internal booking reference.
		ClientReference string `json:"clientReference"`
		// Identifies the agent name of the booking.
		// If empty, by default it will be the same than the apikey performing the booking .
		CreationUser string `json:"creationUser,omitempty"`
		// Free text sent to the hotelier. It can be used to request or inform of special requests to hotelier like:
		// “Non-smoking room preferred”, “Twin bed please”, “Upper floor preferred”, “Late arrival”….
		Remark  string   `json:"remark,omitempty"`
		Voucher *Voucher `json:"voucher,omitempty"`
		// Margin of price difference (as percentage) accepted when a price difference occurs between
		// Availability/CheckRate and Booking operations. Do not use more than two decimal characters when
		// defining tolerance.
		// Example: to input a tolerance of 5%, you should input 5.00.
		// NOTE: Authorisation for the use of this tag is subject to prior written agreement with your sales manager
		Tolerance Amount `json:"tolerance,omitempty"`
		// Language code that defines the language of the response.
		// English will be used by default if this field is not informed.
		Language string               `json:"language,omitempty"`
		Rooms    []ConfirmBookingRoom `json:"rooms"`
	}

	ConfirmBookingRoom struct {
		RateKey string `json:"rateKey"`
		Paxes   []Pax  `json:"paxes"`
	}

	Holder struct {
		Name    string `json:"name"`
		Surname string `json:"surname"`
	}

	PaymentData struct {
		Card       PaymentCard     `json:"paymentCard"`
		Contact    ContactData     `json:"contactData"`
		Billing    *BillingAddress `json:"billingAddress"`
		ThreeDS    *ThreeDSData    `json:"threeDsData"`
		WebPartner int             `json:"webPartner,omitempty"`
		Device     *Device         `json:"device"`
	}

	PaymentCard struct {
		Type   string `json:"cardType"`
		Number string `json:"cardNumber"`
		Expiry string `json:"expiryDate"`
		CVC    string `json:"cardCVC"`
	}

	ContactData struct {
		Email       string `json:"email"`
		PhoneNumber string `json:"phoneNumber"`
	}

	BillingAddress struct {
		Address1    string `json:"address1"`
		Address2    string `json:"address2"`
		City        string `json:"city"`
		State       string `json:"state"`
		PostalCode  string `json:"postalCode"`
		CountryCode string `json:"countryCode"`
	}

	ThreeDSData struct {
		// Enum: "PROVIDED" "REQUESTED"
		// Use PROVIDED if you have your own integration with a PSP
		Option string `json:"option"`
		// This field contains the 3DS version used for authentication
		Version      string                  `json:"version"`
		InfoProvided ThreeDSDataInfoProvided `json:"infoProvided"`
	}

	ThreeDSDataInfoProvided struct {
		// Transaction identifier resulting from authentication processing.
		ID string `json:"id"`
		// Cardholder Authentication Verification Value (CAVV)
		CAVV string `json:"cavv"`
		// Electronic Commerce Indicator (ECI). The ECI value is part of the two data elements
		// that indicate the transaction was processed electronically.
		ECI string `json:"eci"`
	}

	Device struct {
		ID        string `json:"id"`
		IP        string `json:"ip"`
		UserAgent string `json:"userAgent"`
	}

	Voucher struct {
		Language string `json:"language,omitempty"`
		Email    Email  `json:"email"`
		// Path of the picture you want to be featured in the header of the voucher as logo.
		// Accepted formats: PNG, JPG, GIF & BMP. Max size: 100px x 100px.
		// The filename part of the picture path cannot exceed 35 characters length.
		Logo string `json:"logo"`
	}

	Email struct {
		To    string `json:"to"`
		From  string `json:"from"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	ConfirmBookingResponse struct {
		Audit   *AuditData `json:"auditData"`
		Booking *Booking   `json:"booking"`
	}

	Booking struct {
		Reference             string             `json:"reference"`
		CancellationReference string             `json:"cancellationReference"`
		ClientReference       string             `json:"clientReference"`
		CreationDate          Datetime           `json:"creationDate"`
		CreationUser          string             `json:"creationUser"`
		Status                BookingStatus      `json:"status"`
		ModificationPolicy    ModificationPolicy `json:"modificationPolicies"`
		AgComission           string             `json:"agComission"`
		VATComission          string             `json:"comissionVAT"`
		Holder                Holder             `json:"holder"`
		Remark                string             `json:"remark"`
		InvoiceCompany        InvoiceCompany     `json:"invoiceCompany"`
		TotalSellingRate      Amount             `json:"totalSellingRate"`
		TotalNet              Amount             `json:"totalNet"`
		PendingAmount         Amount             `json:"pendingAmount"`
		Currency              string             `json:"currency"`
		Hotel                 BookingHotel       `json:"hotel"`
	}

	BookingHotel struct {
		AvailableHotel
		CheckIn            string          `json:"checkIn"`
		CheckOut           string          `json:"checkOut"`
		TotalSellingRate   Amount          `json:"totalSellingRate"`
		TotalNet           Amount          `json:"totalNet"`
		PendingAmount      Amount          `json:"pendingAmount"`
		Currency           string          `json:"currency"`
		Supplier           *Supplier       `json:"supplier"`
		ClientComments     string          `json:"clientComments"`
		CancellationAmount Amount          `json:"cancellationAmount"`
		Upselling          []UpsellingRate `json:"upselling,omitempty"`
		Keywords           []Keyword       `json:"keywords,omitempty"`
		Reviews            []Review        `json:"reviews,omitempty"`
		Rooms              []BookingRoom   `json:"rooms"`
	}

	BookingRoom struct {
		Code              string        `json:"code"`
		Name              string        `json:"name"`
		Status            BookingStatus `json:"status"`
		ID                int           `json:"id"`
		SupplierReference string        `json:"supplierReference,omitempty"`
		Paxes             []Pax         `json:"paxes"`
		Rates             []Rate        `json:"rates"`
	}

	InvoiceCompany struct {
		RegistrationNumber string `json:"registrationNumber"`
		Code               string `json:"code"`
		Name               string `json:"name"`
	}

	Supplier struct {
		Name      string `json:"name"`
		VATNumber string `json:"vatNumber"`
	}

	UpsellingRate struct {
		Rate
		Discount         string      `json:"discount"`
		DiscountPercent  string      `json:"discountPCT"`
		HotelMandatory   bool        `json:"hotelMandatory"`
		Comission        Amount      `json:"comission"`
		ComissionVAT     string      `json:"comissionVAT"`
		ComissionPercent string      `json:"comissionPCT"`
		Rateup           Amount      `json:"rateup"`
		Brand            string      `json:"brand"`
		Taxes            []Tax       `json:"taxes"`
		DailyRates       []DailyRate `json:"dailyRates"`
	}

	DailyRate struct {
		// Day number of the stay for which the breakdown is made.
		Offset  int    `json:"offset"`
		Net     Amount `json:"net"`
		Selling Amount `json:"selling"`
	}

	Review struct {
		Rate  string `json:"rate"`
		Count int    `json:"reviewCount"`
		Type  string `json:"type"`
	}

	ListBookingsInput struct {
		ListInput
		FilterType            string `json:"filterType"`
		FilterClientReference string `json:"clientReference"`
		FilterCreationUser    string `json:"creationUser"`
		// Parameter to filter the results by the country code of the hotel. Can include multiple values separated by commas.
		FilterCountires    CommaSliceString `json:"country"`
		FilterDestinations CommaSliceString `json:"destination"`
		FilterHotels       CommaSliceInt    `json:"hotel"`
		// Defines the starting date of the range of bookings to be returned.
		FilterStart Datetime `json:"start"`
		// Defines the ending date of the range of bookings to be returned. value.
		FilterEnd Datetime `json:"end"`
	}

	ListBookingsResponse struct {
		Audit    *AuditData `json:"auditData"`
		Bookings []Booking  `json:"bookings"`
	}

	GetBookingResponse struct {
		Audit   *AuditData `json:"auditData"`
		Booking *Booking   `json:"booking"`
	}

	ChangeBookingInput struct {
		// Enum: "SIMULATION" "UPDATE"
		// Defines if the operation will be a simulation or an actual modification.
		Mode     Mode         `json:"mode"`
		Payment  *PaymentData `json:"paymentData"`
		Language string       `json:"language"`
		Booking  *Booking     `json:"booking"`
	}

	ChangeBookingResponse struct {
		Audit   *AuditData `json:"auditData"`
		Booking *Booking   `json:"booking"`
	}

	CancelBookingInput struct {
		// Enum: "CANCELLATION" "SIMULATION"
		// Defines if the operation will be a simulation or an actual cancellation.
		Mode     Mode   `url:"cancellationFlag"`
		Language string `url:"language"`
	}

	CancelBookingResponse struct {
		Audit   *AuditData `json:"auditData"`
		Booking *Booking   `json:"booking"`
	}
)

type BookingStatus string

const (
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCancelled BookingStatus = "CANCELLED"
)

type Mode string

const (
	ModeUpdate       Mode = "UPDATE"
	ModeCancellation Mode = "CANCELLATION"
	ModeSimulation   Mode = "SIMULATION"
)

func (inp *ListAvailableHotelsInput) Validate() error {
	if err := inp.Stay.Validate(); err != nil {
		return err
	}
	if inp.Filter != nil {
		if err := inp.Filter.Validate(); err != nil {
			return err
		}
	}
	if err := inp.Hotels.Validate(); err != nil {
		return err
	}
	return nil
}

type Stay struct {
	CheckIn  string `json:"checkIn"`
	CheckOut string `json:"checkOut"`
	// Amount of days after and before the check-in to check availability, keeping the same stay duration.
	ShiftDays int `json:"shiftDays,omitempty"`
	// Defines if results are returned for shiftDays even if there's no results for the searched dates.
	AllowOnlyShift *bool `json:"allowOnlyShift,omitempty"`
}

func (stay *Stay) Validate() error {
	if stay.ShiftDays > 5 {
		return errors.New("ShiftDays is invalid (should <=5)")
	}
	return nil
}

type Occupancy struct {
	Rooms    int `json:"rooms"`
	Adults   int `json:"adults"`
	Children int `json:"children"`
	// Use Paxes only when has children.
	Paxes []Pax `json:"paxes,omitempty"`
}

type Pax struct {
	Type    PaxType `json:"type"`
	Age     int     `json:"age"`
	Name    string  `json:"name,omitempty"`
	Surname string  `json:"surname,omitempty"`
	RoomID  int     `json:"roomId,omitempty"`
}

type PaxType string

const (
	PaxTypeAdult    PaxType = "AD"
	PaxTypeChildren PaxType = "CH"
)

type Keyword struct {
	Keywords    []int `json:"keywords,omitempty"`
	AllIncluded bool  `json:"allIncluded"`
}

type Geolocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    Radius  `json:"radius"`
	Unit      Unit    `json:"unit"`
}

func (geo *Geolocation) Validate() error {
	if geo.Latitude == 0 {
		return &ValidationError{
			FieldName: "Latitude",
			Required:  true,
		}
	}
	if geo.Longitude == 0 {
		return &ValidationError{
			FieldName: "Longitude",
			Required:  true,
		}
	}
	if geo.Radius > 200 {
		return &ValidationError{
			FieldName: "Radius",
			Min:       0,
			Max:       200,
		}
	}
	if geo.Unit != "" && geo.Unit != UnitMiles && geo.Unit != UnitKilometers {
		return &ValidationError{
			FieldName: "Unit",
			Allow:     []string{UnitMiles.String(), UnitKilometers.String()},
		}
	}
	return nil
}

type Filter struct {
	MaxHotels       int       `json:"maxHotels,omitempty"`
	MaxRooms        int       `json:"maxRooms,omitempty"`
	MinRate         FloatRate `json:"minRate,omitempty"`
	MaxRate         FloatRate `json:"maxRate,omitempty"`
	MaxRatesPerRoom int       `json:"maxRatesPerRoom"`
	MinCategory     int       `json:"minCategory,omitempty"`
	MaxCategory     int       `json:"maxCategory,omitempty"`
}

func (filter *Filter) Validate() error {
	if filter.MaxHotels < 1 || filter.MaxHotels > 2000 {
		return &ValidationError{
			FieldName: "MaxHotels",
			Min:       1,
			Max:       2000,
		}
	}
	if filter.MaxRooms < 1 || filter.MaxRooms > 50 {
		return &ValidationError{
			FieldName: "MaxRooms",
			Min:       1,
			Max:       50,
		}
	}
	if filter.MinCategory < 1 || filter.MinCategory > 5 {
		return &ValidationError{
			FieldName: "MinCategory",
			Min:       1,
			Max:       5,
		}
	}
	if filter.MaxCategory < 1 || filter.MaxCategory > 5 {
		return &ValidationError{
			FieldName: "MaxCategory",
			Min:       1,
			Max:       5,
		}
	}
	return nil
}

type FilterBoards struct {
	Boards   []string `json:"boards"`
	Included bool     `json:"included"`
}

type FilterRooms struct {
	Codes    []string `json:"room"`
	Included bool     `json:"included"`
}

type FilterHotel struct {
	HotelCodes []int `json:"hotel"`
}

func (f *FilterHotel) Validate() error {
	if len(f.HotelCodes) > 2000 {
		return &ValidationError{
			FieldName: "FilterHotel.Hotel",
			Max:       2000,
		}
	}
	return nil
}

type PaymentType string

const (
	PaymentTypeAtWeb   PaymentType = "AT_WEB"
	PaymentTypeAtHotel PaymentType = "AT_HOTEL"
)

// Ref - https://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/#operation/availability
func (api *API) ListAvailableHotels(ctx context.Context, inp *ListAvailableHotelsInput) (*ListAvailableHotelsResponse, error) {
	if err := inp.Validate(); err != nil {
		return nil, err
	}
	return clientx.NewRequestBuilder[ListAvailableHotelsInput, ListAvailableHotelsResponse](api.API).
		Post("/hotel-api/1.0/hotels", inp, clientx.WithRequestHeaders(api.buildHeaders())).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/#operation/checkRate
func (api *API) ListCheckRates(ctx context.Context, inp *ListCheckRatesInput) (*ListCheckRatesResponse, error) {
	return clientx.NewRequestBuilder[ListCheckRatesInput, ListCheckRatesResponse](api.API).
		Post("/hotel-api/1.0/checkrates", inp, clientx.WithRequestHeaders(api.buildHeaders())).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// https://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/#operation/bookingDetail
func (api *API) GetBooking(ctx context.Context, id string) (*GetBookingResponse, error) {
	return clientx.NewRequestBuilder[struct{}, GetBookingResponse](api.API).
		Get("/hotel-api/1.0/bookings/"+id, clientx.WithRequestHeaders(api.buildHeaders())).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/#operation/bookingList
func (api *API) ListBookings(ctx context.Context, inp *CancelBookingInput) (*CancelBookingResponse, error) {
	return clientx.NewRequestBuilder[CancelBookingInput, CancelBookingResponse](api.API).
		Delete("/hotel-api/1.0/bookings", inp, clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/#operation/booking
func (api *API) ConfirmBooking(ctx context.Context, inp *ConfirmBookingInput) (*ConfirmBookingResponse, error) {
	return clientx.NewRequestBuilder[ConfirmBookingInput, ConfirmBookingResponse](api.API).
		Post("/hotel-api/1.2/bookings", inp, clientx.WithRequestHeaders(api.buildHeaders())).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/#operation/bookingChange
func (api *API) ChangeBooking(ctx context.Context, id string, inp *ChangeBookingInput) (*ChangeBookingResponse, error) {
	return clientx.NewRequestBuilder[ChangeBookingInput, ChangeBookingResponse](api.API).
		Put("/hotel-api/1.0/bookings/"+id, inp, clientx.WithRequestHeaders(api.buildHeaders())).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}

// Ref - https://developer.hotelbeds.com/documentation/hotels/booking-api/api-reference/#operation/bookingCancellation
func (api *API) CancelBooking(ctx context.Context, id string, inp *CancelBookingInput) (*CancelBookingResponse, error) {
	return clientx.NewRequestBuilder[CancelBookingInput, CancelBookingResponse](api.API).
		Delete("/hotel-api/1.0/bookings/"+id, nil, clientx.WithRequestHeaders(api.buildHeaders())).
		WithQueryParams("url", *inp).
		WithErrorDecode(func(resp *http.Response) (bool, error) {
			return resp.StatusCode > 399, decodeError(resp)
		}).
		DoWithDecode(ctx)
}
