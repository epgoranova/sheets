package components

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// ownerRWPermissions marks read and write permissions for the owner.
const ownerRWPermissions = 0600

// WriteSliceToFile writes a slice to the file at the filepath. It writes each
// element of the slice on a separate line.
func WriteSliceToFile(filepath string, values []string) error {
	data := strings.Join(values, "\n")
	return ioutil.WriteFile(filepath, []byte(data), ownerRWPermissions)
}

// WriteSliceToStdout prints the elements of the slice. Each element is on a new
// line.
func WriteSliceToStdout(values []string) {
	for _, value := range values {
		fmt.Println(value)
	}
}
