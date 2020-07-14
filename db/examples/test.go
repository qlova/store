package main

import (
	"fmt"

	"qlova.store/db"
	"qlova.store/db/driver/postgres"
)

//Customer is the customers table
var Customer struct {
	db.Model `db:"Customers"`

	ID         db.Int    `db:"CustomerID"`
	Name       db.String `db:"CustomerName"`
	Contact    db.String `db:"ContactName"`
	Address    db.String
	City       db.String
	PostalCode db.String
	Country    db.String

	Token db.String
}

func main() {
	err := postgres.Open("", "host=localhost sslmode=disable user=postgres dbname=postgres password=eeuMjXaurD port=5433")
	if err != nil {
		panic(err.Error())
	}

	if err := db.Register(&Customer); err != nil {
		panic(err.Error())
	}

	if err := db.Truncate(Customer.Table); err != nil {
		panic(err.Error())
	}

	if err := db.Insert(Customer.With(
		Customer.ID.SetTo(0),
		Customer.Name.SetTo("Cardinal"),
		Customer.Contact.SetTo("Tom B. Erichsen"),
		Customer.Address.SetTo("Skagen 21"),
		Customer.City.SetTo("Stavanger"),
		Customer.PostalCode.SetTo("4006"),
		Customer.Country.SetTo("Norway"),
	)); err != nil {
		panic(err.Error())
	}

	if err := db.Insert(Customer.With(
		Customer.ID.SetTo(1),
		Customer.Name.SetTo("Quentin"),
		Customer.Contact.SetTo("Nathalie"),
		Customer.Address.SetTo("10 Rosemary Place"),
		Customer.City.SetTo("Katikati"),
		Customer.PostalCode.SetTo("3129"),
		Customer.Country.SetTo("New Zealand"),
	)); err != nil {
		panic(err.Error())
	}

	if _, err := db.Where(Customer.Name.Equals("Quentin")).Update(
		Customer.City.To("Auckland"),
	); err != nil {
		panic(err.Error())
	}

	var Address = Customer.Address
	if err := db.Where(Customer.Name.Equals("Quentin")).Get(&Address); err != nil {
		panic(err.Error())
	}
	fmt.Println(Address)

	var customer = Customer
	if err := db.SortBy(Customer.ID.Column).Slice(1, 1).Into(&customer); err != nil {
		panic(err.Error())
	}

	for range customer.Range() {
		db.Next(&customer)

		fmt.Println(customer)
	}
}
