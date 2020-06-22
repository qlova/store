package fs

import (
	"path"
)

//Path is a file-system path.
type Path string

func (p Path) String() string {
	return string(p)
}

//Base returns the last element of path.
//Trailing slashes are removed before extracting the last element.
//If the path is empty, Base returns ".". If the path consists entirely of slashes, Base returns "/".
func (p Path) Base() string {
	return path.Base(p.String())
}

//Join joins any number of path elements into a single path, separating them with slashes. Empty elements are ignored. The result is Cleaned. However, if the argument list is empty or all its elements are empty, Join returns an empty string.
func (p Path) Join(others ...Path) Path {
	var strings = make([]string, len(others)+1)
	strings[0] = p.String()
	for i := range others {
		strings[i+1] = others[i].String()
	}
	return Path(path.Join(strings...))
}

//Dir returns all but the last element of path, typically the path's directory.
//After dropping the final element using Split, the path is Cleaned and trailing slashes are removed.
//If the path is empty, Dir returns ".". If the path consists entirely of slashes followed by non-slash bytes, Dir returns a single slash.
//In any other case, the returned path does not end in a slash.
func (p Path) Dir() Path {
	return Path(path.Dir(p.String()))
}
