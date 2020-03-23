package main

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
	. "gophercises/phone/db"
	. "gophercises/phone/normalizer"
)

var PhoneNumbers = []string{
	"1234567890",
	"123 456 7891",
	"(123) 456 7892",
	"(123) 456-7893",
	"123-456-7894",
	"123-456-7890",
	"1234567892",
	"(123)456-7892",
}

func main() {
	db := PhoneDB{}
	db.Create(PhoneNumber{})
	defer db.Close()

	for _, number := range PhoneNumbers {
		phoneNumber := PhoneNumber{
			Number:     number,
			Owner:      "EE",
			NumberType: "Mobile",
		}
		db.Insert(&phoneNumber)
	}

	formattedNumbers := make(map[string]bool)
	retrievedNumbers := db.QueryAll()
	for _, number := range retrievedNumbers {
		number.Normalize()
		if formattedNumbers[number.Number] == false {
			db.Update(&number)
			formattedNumbers[number.Number] = true
		} else {
			db.Delete(&number)
		}
	}
}
