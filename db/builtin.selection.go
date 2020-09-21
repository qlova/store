package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"sort"
)

type selection struct {
	b Builtin

	table string

	conditions []func(reflect.Value) bool
	updates    []func(reflect.Value) error

	sort  Sorter
	sorts []Sorter

	link  Linker
	links []Linker

	offset, length int

	columns []Variable
}

func (s *selection) addCondition(c Condition) {
	switch c.Operator {
	case NoOperator:
		return
	case OpEquals:
		s.conditions = append(s.conditions, func(v reflect.Value) bool {
			return reflect.DeepEqual(v.FieldByName(c.Column).Interface(), c.Value)
		})
	case OpNotEquals:
		s.conditions = append(s.conditions, func(v reflect.Value) bool {
			return !reflect.DeepEqual(v.FieldByName(c.Column).Interface(), c.Value)
		})
	default:
		panic("not implemented") // TODO: Implement
	}
}

func (s *selection) addUpdate(u Update) {
	s.updates = append(s.updates, func(v reflect.Value) error {
		v.FieldByName(u.Column).Set(reflect.ValueOf(u.Value))
		return nil
	})
}

func (s selection) MarshalJSON() ([]byte, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	var table = database[index(s.b, s.table)]

	if table == nil {
		return nil, ErrTableNotFound
	}

	var results = s.query(table)

	var buffer bytes.Buffer

	buffer.WriteByte('[')

	for i, index := range results {
		b, err := json.Marshal(table.slice.Index(index).Interface())
		if err != nil {
			return nil, err
		}
		buffer.Write(b)
		if i < len(results)-1 {
			buffer.WriteByte(',')
		}
	}

	buffer.WriteByte(']')

	return buffer.Bytes(), nil
}

//Count returns the number of results.
func (s selection) Count(v Viewable) (int, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	var table = database[index(s.b, s.table)]

	if table == nil {
		return 0, ErrTableNotFound
	}

	return len(s.query(table)), nil
}

//Sum returns the sum amount of the values of all selected rows.
func (s selection) Sum(v Variable) error {
	mutex.RLock()
	defer mutex.RUnlock()

	var table = database[index(s.b, s.table)]

	if table == nil {
		return ErrTableNotFound
	}

	var results = s.query(table)

	var sum = v.Pointer()

	reflect.ValueOf(sum).Elem().Set(reflect.Zero(reflect.ValueOf(sum).Elem().Type()))

	for _, index := range results {
		var value = table.slice.Index(index).FieldByName(v.Column()).Interface()
		switch val := value.(type) {
		case uint:
			*(sum.(*uint)) += val
		case uint8:
			*(sum.(*uint8)) += val
		case uint16:
			*(sum.(*uint16)) += val
		case uint32:
			*(sum.(*uint32)) += val
		case uint64:
			*(sum.(*uint64)) += val

		case int:
			*(sum.(*int)) += val
		case int8:
			*(sum.(*int8)) += val
		case int16:
			*(sum.(*int16)) += val
		case int32:
			*(sum.(*int32)) += val
		case int64:
			*(sum.(*int64)) += val

		case float32:
			*(sum.(*float32)) += val
		case float64:
			*(sum.(*float64)) += val

		default:
			return errors.New("cannot sum type: " + reflect.TypeOf(sum).Elem().String())
		}
	}

	return nil
}

//Average returns the average value in the given column for all results.
func (s selection) Average(v Viewable) (float64, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	var table = database[index(s.b, s.table)]

	if table == nil {
		return 0, ErrTableNotFound
	}

	var results = s.query(table)

	var avg float64

	for _, index := range results {
		var value = table.slice.Index(index).FieldByName(v.Column()).Interface()
		switch val := value.(type) {
		case int:
			avg += float64(val)
		case int8:
			avg += float64(val)
		case int16:
			avg += float64(val)
		case int32:
			avg += float64(val)
		case int64:
			avg += float64(val)
		case float32:
			avg += float64(val)
		case float64:
			avg += float64(val)
		default:
			return 0, errors.New("cannot average type: " + reflect.TypeOf(value).Elem().String())
		}
	}

	return avg / float64(len(results)), nil
}

//Update updates the selected items with the given updates.
//Returns the number of items updated (or -1 if the statistic is unavailable).
func (s selection) Update(update Update, updates ...Update) (int, error) {

	mutex.Lock()
	defer mutex.Unlock()

	if update.Column != "" {
		s.addUpdate(update)
	}
	for _, update := range updates {
		s.addUpdate(update)
	}

	var table = database[index(s.b, update.Table)]

	if table == nil {
		return 0, ErrTableNotFound
	}

	var results = s.query(table)

	for _, index := range results {
		row := table.slice.Index(index)
		for _, update := range s.updates {
			update(row)
		}
	}

	return len(results), nil
}

//Delete deletes all the selected items.
func (s selection) Delete() (int, error) {
	mutex.Lock()
	defer mutex.Unlock()

	var table = database[index(s.b, s.table)]

	if table == nil {
		return 0, ErrTableNotFound
	}

	var results = s.query(table)

	//Sort indicies so that our indicies stay correct during removal.
	sort.Sort(sort.Reverse(sort.IntSlice(results)))

	for _, index := range results {
		var last = table.slice.Len() - 1
		table.slice.Index(index).Set(table.slice.Index(last))
		table.slice.Set(table.slice.Slice(0, last))
	}

	return len(results), nil
}

func get(variable, value reflect.Value, column string) {
	variable.Elem().Set(value.FieldByName(column))
}

//Get loads the given columns of the selection into those columns.
func (s selection) Get(v Variable, vs ...Variable) (int, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	if v.Master() {
		return 0, ErrIllegalMaster
	}

	if v.Database() == nil {

		return 0, ErrDisconnectedViewer
	}

	var table = database[index(s.b, v.Table())]

	if table == nil {
		return 0, ErrTableNotFound
	}

	var results = s.query(table)

	if len(results) == 0 {
		return 0, ErrNotFound
	}

	if s.length == 1 {
		row := table.slice.Index(results[0])

		get(reflect.ValueOf(v.Pointer()), row, v.Column())
		for _, v := range vs {
			get(reflect.ValueOf(v.Pointer()), row, v.Column())
		}

		return 1, nil
	}

	v.Make(len(results))
	for _, v := range vs {
		v.Make(len(results))
	}

	var count int

	for _, index := range results {
		row := table.slice.Index(index)

		get(reflect.ValueOf(v.Slice(count)), row, v.Column())
		for _, v := range vs {
			get(reflect.ValueOf(v.Slice(count)), row, v.Column())
		}

		count++
	}

	return len(results), nil
}

func into(viewer, value reflect.Value) {
	var vtype = viewer.Type()

	for i := 1; i < vtype.NumField(); i++ {
		viewer.Field(i).Addr().MethodByName("Set").Call([]reflect.Value{value.FieldByName(vtype.Field(i).Name)})
	}
}

//Search with the given filter and return the results.
func (b Builtin) Search(f Filter) Results {
	var s selection
	s.b = b

	s.table = f.Table

	if f.Condition.Operator != 0 {
		s.addCondition(f.Condition)
	}
	for _, condition := range f.Conditions {
		s.addCondition(condition)
	}

	s.sort = f.Sort
	s.sorts = f.Sorts

	s.link = f.Link
	s.links = f.Links

	s.length = f.Length
	s.offset = f.Offset

	return s
}
