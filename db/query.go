package db

type Operator int

//Operators
const (
	Equals Operator = iota
	Contains
	NotEquals
	StartsWith
)

type Condition struct {
	Column
	Operator
	Value interface{}

	Or *Condition
}

func Either(a, b Condition) Condition {
	a.Or = &b
	return a
}

type Querier interface {
	Where(Condition, ...Condition) Query
	SortBy(Column, ...Column) Query

	Link(Linker, ...Linker) Query

	Update(Update, ...Update) (int, error)

	Slice(int, int, ...Value) Slicer

	Count(Value) (int, error)
	Average(Value) (float64, error)
}

type Query interface {
	Querier
	Delete() (int, error)

	Get(MutableValue, ...MutableValue) error
	Read(Connectable) error
}

func Where(condition Condition, conditions ...Condition) Query {
	return connections[condition.Connection].Where(condition, conditions...)
}

func SortBy(column Column, columns ...Column) Query {
	return connections[column.Connection].SortBy(column, columns...)
}
