package main

import (
	"fmt"

	"github.com/UpLiftL1f3/hotel-reservation/types"
)

func main() {
	hotel := types.Hotel{
		Name:     "Bellucia",
		Location: "Texas",
	}

	room := types.Room{
		Type:      types.SingleRoomType,
		BasePrice: 99.9,
	}

	fmt.Println("seeding the database")
}
