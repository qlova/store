package db

//Text is searchable text.
type Text struct {
	String

	//WordIndex is used to speed up search.
	WordIndex Int64

	//Dictionary used for indexing.
	Dictionary Dictionary
}

var _ Variable = new(Text)

//Has does a word search on the text.
func (t *Text) Has(word string) Condition {
	var index = []int64{0}
	if err := t.Dictionary.LookupWords(index, word); err != nil || index[0] == 0 {
		return Condition{}
	}
	return Both(t.WordIndex.NotEquals(0), t.WordIndex.DivisibleBy(index[0]))
}

//Search does a search checking if the tokens in the query string are also found in the text.
func (t *Text) Search(query string) Condition {
	words := Tokenise(query)

	var index int64 = 1

	var indicies = make([]int64, len(words))
	if err := t.Dictionary.LookupWords(indicies, words...); err != nil {
		return Condition{}
	}

	for _, i := range indicies {
		if i > 0 {
			index *= i
		}
	}

	if index == 1 {
		return False
	}

	return Both(t.WordIndex.NotEquals(0), t.WordIndex.DivisibleBy(index))
}

func (t *Text) To(new string) Update {
	if t.Dictionary == nil {
		return t.String.To(new).And(t.WordIndex.To(0))
	}

	var index int64 = 1

	words := Tokenise(new)

	var indicies = make([]int64, len(words))

	if err := t.Dictionary.LookupWords(indicies, words...); err != nil {
		return t.String.To(new).And(t.WordIndex.To(0))
	}

	for _, i := range indicies {
		if i > 0 {
			index *= i
		}
	}

	if index == 1 {
		index = 0
	}

	return t.String.To(new).And(t.WordIndex.To(index))
}

//Set text to the given string.
func (t *Text) Set(to string) {
	defer t.String.Set(to)

	if t.Dictionary == nil {
		t.WordIndex.Set(0)
		return
	}

	var index int64 = 1

	words := Tokenise(to)

	var indicies = make([]int64, len(words))

	if err := t.Dictionary.LookupWords(indicies, words...); err != nil {
		t.WordIndex.Set(0)
		return
	}

	for _, i := range indicies {
		if i > 0 {
			index *= i
		}
	}

	t.WordIndex.Set(index)
	return
}
