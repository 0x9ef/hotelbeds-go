package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/0x9ef/hotelbeds-go"
)

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

	r, err := json.Marshal(resp)
	fmt.Println(string(r), err)
}
