package bolt

type Cursor struct{}

//Return a unique string representing this cursor location.
func (cursor *Cursor) Name() string {
	return ""
}

//Goto the cursor location previously returned by Name()
func (cursor *Cursor) Goto(name string) {

}

func (cursor *Cursor) Stores(count int) []string {
	return nil
}

func (cursor *Cursor) Data(count int) []string {
	return nil
}
