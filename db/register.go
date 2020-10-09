package db

import (
	"errors"
	"reflect"
	"strings"
	"unsafe"
)

//Connect initialises and connects the given viewer.
func Connect(viewer Viewer, driver Driver) error {
	if viewer.Master() {
		panic("invalid register of master viewer, datarace")
	}

	var vtype = reflect.TypeOf(viewer).Elem()
	var rvalue = reflect.ValueOf(viewer).Elem()

	var table table

	var EmbeddedView reflect.Value

	for i := 0; i < rvalue.NumField(); i++ {
		var rvalue = rvalue.Field(i)
		var field = vtype.Field(i)

		if field.Type == reflect.TypeOf(View{}) {
			if name, ok := field.Tag.Lookup("db"); ok {
				EmbeddedView = rvalue
				table.name = name
				break
			}
			return errors.New("db.Register: embedded View must have a db tag `db:\"table name\"` ")
		}
	}

	var columns []Column

	for i := 0; i < rvalue.NumField(); i++ {
		var rvalue = rvalue.Field(i)
		var field = vtype.Field(i)

		if setter, ok := rvalue.Addr().Interface().(value); ok {
			var name = field.Name
			var key bool

			//The name can be overriden in the tag.
			//Further tags include key and unique
			if tag, ok := field.Tag.Lookup("db"); ok {
				args := strings.Split(tag, ",")
				if args[0] != "" {
					name = args[0]
				}
				if len(args) > 0 {
					if args[1] == "key" {
						key = true
					}
				}
			}

			setter.setprivate(
				table.name, name,
				field.Offset,
				key,
				driver,
				viewer,
			)

			columns = append(columns, setter)

			//Special case for Text.
			if t, ok := rvalue.Addr().Interface().(*Text); ok {

				t.WordIndex.setprivate(
					table.name, name+"_index",
					field.Offset+unsafe.Offsetof(t.WordIndex),
					false,
					driver,
					viewer,
				)

				columns = append(columns, &t.WordIndex)
			}
		}
	}

	table.columns = columns

	/*if table.Name == "" {
		return errors.New("table name not detected")
	}*/

	EmbeddedView.Set(reflect.ValueOf(View{
		table:  table,
		master: viewer,
		vtype:  vtype,
	}))

	viewer.SetDriver(driver)

	/*return connections[model.getModel().Connection].Verify(Schema{
		Table:   model.GetTable(),
		Columns: columns,
	})*/

	return nil
}
