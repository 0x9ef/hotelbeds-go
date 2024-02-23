package gethoteldetails

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

	resp, err := api.GetHotelDetails(ctx, []int{6613, 6619}, &hotelbeds.GetHotelDetailsInput{})
	if err != nil {
		panic(err)
	}

	r, err := json.Marshal(resp)
	fmt.Println(string(r), err)
}
