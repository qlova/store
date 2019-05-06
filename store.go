package store

import "io"
import "io/ioutil"

//A location can be created or deleted.
type Location interface {
	//Return the full location.
	String() string
	
	//Delete at this location.
	Delete() error
	
	//Does this location exist?
	Exists() bool
}

//A store can store data and stores.
type Store interface {
	Location
	
	//List the files and stores in this directory.
	List() []string //TODO replace this with cursor
	
	//Create at this location.
	Create() error
	
	//Return the store at the given relative path.
	Goto(name string) Store
	
	//Return the data at the given relative path.
	Data(name string) Data
	
	//Holds the latest error.
	Error() error
}

//Data can be opened for reading or writing.
type Data struct {
	err error

	Raw interface{
		Location

		//Create, and open the data for writing.
		Create() io.WriteCloser
		
		//Open the data for reading.
		Reader() io.ReadCloser
		
		//Return the size of the file.
		Size() int64
		
		//Holds the latest error.
		Error() error
	}
}

func (data Data) SetString(s string) error {
	var raw = data.Raw.Create()
	if err := data.Raw.Error(); err != nil {
		return err
	}
	
	_, err := raw.Write([]byte(s))
	if err != nil {
		return err
	}	
	
	raw.Close()
	return nil
}

func (data Data) String() string {
	var raw = data.Raw.Reader()
	if err := data.Raw.Error(); err != nil {
		data.err = err
		return ""
	}
	
	binary, err := ioutil.ReadAll(raw)
	if err != nil {
		data.err = err
		return ""
	}
	
	var result = string(binary)
	
	err = raw.Close()
	if err != nil {
		data.err = err
	}

	return result
}
