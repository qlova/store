package db_test

import (
	"fmt"

	"qlova.store/db"
)

type UserViewer struct {
	db.View `db:"users"`

	ID db.UUID `db:",key"`

	Name db.String
}

func Example() {
	var User UserViewer

	defer db.Open().Connect(&User).Close()

	db.Sync(User)

	var bob = User
	bob.Name.Set("Bob")

	db.Insert(bob)

	var user = User
	db.If(user.Name.Equals("Bob")).Get(&user)

	fmt.Println(user.Name)

	// Output: Bob
}
