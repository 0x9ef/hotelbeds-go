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

	resp, err := api.ConfirmBooking(ctx, &hotelbeds.ConfirmBookingInput{
		Holder: hotelbeds.Holder{
			Name:    "HolderFirstName",
			Surname: "HolderLastName",
		},
		Rooms: []hotelbeds.ConfirmBookingRoom{
			{
				RateKey: "20240402|20240403|W|164|6619|TWN.ST|BAR BB FLEX 14|BB||1~1~0||N@06~~21e12c~1630615603~S~~~NOR~5F05A4B7D40E44A170871765642600AADE00000010000000006248118",
				Paxes: []hotelbeds.Pax{
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
	if err != nil {
		panic(err)
	}

	r, err := json.Marshal(resp)
	fmt.Println(string(r), err)
}
