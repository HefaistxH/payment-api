package main

import (
	"fmt"
	"mnc-techtest/delivery"
)

func main() {
	fmt.Println("mnc-merchant-api")
	delivery.NewServer().Run()
}
