package main

import (
	"fmt"
	"mnc-techtest/delivery"

	"golang.org/x/crypto/bcrypt"
)

func main() {

	password := "Password123"
	hashedPassword := "$2a$12$7fb6RY3cGve07FAjf51gsOlfw0bbowS2JZ5WSCDmJKDtkwZ/JoOx6"
	// Compare the hashed password with the plain-text password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println("Invalid credentials")
	} else {
		fmt.Println("Password matched")
	}
	fmt.Println("mnc-merchant-api")
	server := delivery.NewServer()
	server.Run()
}
