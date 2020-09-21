package db

//LessThan returns a condition that is true if i is less then val.
func (i Int64) LessThan(val int64) Condition {
	return Condition{
		Column:   i.column,
		Operator: OpLessThan,
		Value:    val,
	}
}

//Contains returns a condition that is true if s contains val.
func (s String) Contains(val string) Condition {
	return Condition{
		Column:   s.column,
		Operator: OpContains,
		Value:    val,
	}
}

//HasPrefix returns a condition that is true if s starts with val.
func (s String) HasPrefix(val string) Condition {
	return Condition{
		Column:   s.column,
		Operator: OpHasPrefix,
		Value:    val,
	}
}
