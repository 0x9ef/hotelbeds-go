# HotelBeds client written in Go
[![Go Tests](https://github.com/0x9ef/hotelbeds-go/actions/workflows/ci.yaml/badge.svg)](https://github.com/0x9ef/hotelbeds-go/actions/workflows/ci.yaml)

Go client written for [HotelBeds](https://hotelbeds.com/). It's unofficial client, currently supported only by my [0x9ef](https://github.com/0x9ef). I will tack version changes and changelogs as soon as possible.

## Installation
> NOTE: Requires at least Go 1.18 since we use generics

To get latest version use:
```
go get github.com/0x9ef/hotelbeds-go@latest
```

To specify version use:
```
go get github.com/0x9ef/clientx@1.24.4 # version
```

## Usage Examples
See [examples/](https://github.com/0x9ef/hotelbeds-go/tree/master/examples) folder or `_test.go` files.

## Getting Started
The client is built on [ClientX](https://github.com/0x9ef/clientx) library, so if there is no way to implement some functionality through because of ClientX, please submit [issue](https://github.com/0x9ef/clientx/issues) or create [pull request](https://github.com/0x9ef/clientx/pulls) in ClientX repository.


## Initialization
The client automatically builds authorization headers and calculates X-Signature from current timestamp, so you don't need to do it manually.

```go
api := hotelbeds.New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))
```

## Useful information
Useful articles to stay tuned:
- [Booking API](https://developer.hotelbeds.com/documentation/hotels/booking-api/workflow/)
- [Content API](https://developer.hotelbeds.com/documentation/hotels/content-api/how-use-content-api/) 
- [Cache API](https://developer.hotelbeds.com/documentation/hotels/cache-api/workflows/)

## Get Availability
```go
func main() {
	api := hotelbeds.New(os.Getenv("HOTELBEDS_API_KEY"), os.Getenv("HOTELBEDS_API_SECRET"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := api.ListAvailableHotels(ctx, &hotelbeds.ListAvailableHotelsInput{
		Stay: hotelbeds.Stay{
			CheckIn:  "2024-04-02",
			CheckOut: "2024-04-03",
		},
		Occupancies: []hotelbeds.Occupancy{
			{
				Rooms:  1,
				Adults: 1,
			},
		},
		Hotels: hotelbeds.FilterHotel{
			HotelCodes: []int{6619, 6613},
		},
	})
	if err != nil {
		panic(err)
	}

	... // do something with resp
}
```

## Implementation status

Internal:
- [x] API Client
- [x] Rate Limiting
- [x] Retry Mechanism
- [x] Error Handling

Hotel APIs:
- [x] Availability
- [x] Check rates
- [x] Booking Confirmation
- [x] Booking List
- [x] Booking Detail
- [x] Booking Change
- [x] Booking Cancellation
- [ ] Booking Reconfirmation

Content APIs:
- [x] Hotels List
- [x] Hotel Details
- [x] Countries
- [x] Destinations
- [x] Acommodations
- [x] Boards
- [x] BoardGroups
- [x] Categories
- [x] Chains
- [x] Classifications
- [x] Currencies
- [x] Facilities
- [x] Facility Groups
- [x] Facility Typologies
- [x] Image Types
- [x] Issues
- [x] Languages
- [x] Promotions
- [x] Rate Comments
- [ ] Rate Comment Details
- [x] Rooms
- [x] Segments

## License

This source code is licensed under the MIT license found
in the LICENSE file in the root directory of this source tree.
