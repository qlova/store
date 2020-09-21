package db

import (
	"bytes"
	"reflect"
	"sort"
)

type sortable struct {
	indicies []int

	table  *storage
	tables []*storage

	sorter  Sorter
	sorters []Sorter

	linker  Linker
	linkers []Linker
}

func (s sortable) Len() int {
	return len(s.indicies)
}

func (s sortable) compare(a, b interface{}) bool {
	switch a.(type) {
	case uint:
		return a.(uint) < b.(uint)
	case uint8:
		return a.(uint8) < b.(uint8)
	case uint16:
		return a.(uint16) < b.(uint16)
	case uint32:
		return a.(uint32) < b.(uint32)
	case uint64:
		return a.(uint64) < b.(uint64)

	case int:
		return a.(int) < b.(int)
	case int8:
		return a.(int8) < b.(int8)
	case int16:
		return a.(int16) < b.(int16)
	case int32:
		return a.(int32) < b.(int32)
	case int64:
		return a.(int64) < b.(int64)

	case float32:
		return a.(float32) < b.(float32)
	case float64:
		return a.(float64) < b.(float64)

	case string:
		return a.(string) < b.(string)

	case bool:
		return a == false

	case []byte:
		return bytes.Compare(a.([]byte), b.([]byte)) == -1

	}
	panic("unsortable type: " + reflect.TypeOf(a).String())
}

func (s sortable) Less(i, j int) bool {
	a := s.table.slice.Index(i)
	b := s.table.slice.Index(j)

	less := s.compare(a.FieldByName(s.sorter.Column).Interface(), b.FieldByName(s.sorter.Column).Interface())
	if s.sorter.Decreasing {
		less = !less
	}

	if !less {
		return false
	}

	for _, sorter := range s.sorters {
		less = s.compare(a.FieldByName(sorter.Column).Interface(), b.FieldByName(sorter.Column).Interface())
		if sorter.Decreasing {
			less = !less
		}

		if !less {
			return false
		}
	}

	return true
}

func (s sortable) Swap(i, j int) {
	s.indicies[i], s.indicies[j] = s.indicies[j], s.indicies[i]
}

func (s selection) query(table *storage) []int {
	var results []int

	for i := 0; i < table.slice.Len(); i++ {
		row := table.slice.Index(i)

		var matches bool = true
		for _, condition := range s.conditions {
			if !condition(row) {
				matches = false
				break
			}
		}

		if matches {
			results = append(results, i)
		}
	}

	if s.sort.Column != "" {
		sort.Sort(sortable{
			results,
			table,
			nil,
			s.sort,
			s.sorts,
			s.link,
			s.links,
		})
	}

	return results
}
