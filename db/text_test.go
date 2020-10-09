package db_test

import (
	"testing"

	"qlova.org/should"

	"qlova.store/db"
)

type TextViewer struct {
	db.View `db:"text"`

	Body db.Text
}

func Test_Text(t *testing.T) {
	var TextMaster = TextViewer{}
	defer db.Open().Connect(&TextMaster).Close()

	should.NotError(
		db.Sync(&TextMaster),
	).Test(t)

	d := db.NewDictionary()
	d.Add("Hello", "World", "Test")

	TextMaster.Body.Dictionary = d

	var Text = TextMaster

	Text.Body.Set("Hello World")

	should.NotError(
		db.Insert(Text),
	).Test(t)

	Text = TextMaster

	should.NotError(
		db.If(Text.Body.Has("Hello")).Get(&Text),
	).Test(t)

	should.Be("Hello World")(Text.Body.Value()).Test(t)

	Text = TextMaster

	should.Error(
		db.If(Text.Body.Has("Test")).Get(&Text),
	).Test(t)
}
