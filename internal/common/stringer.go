package common

import "fmt"

func (r Record) String() string {
	repr := ""
	repr += fmt.Sprintf("\n  Type: %s", r.Type)
	repr += fmt.Sprintf("\n  Name: %s", r.Name)
	if r.Meta != "" {
		repr += fmt.Sprintf("\n  Meta info: %s", r.Meta)
	}
	if r.Type != BinaryRecord && r.Opaque != "" {
		repr += fmt.Sprintf("\n  Data: %s", r.Opaque)
	}

	return repr
}

func (rr Records) String() string {
	repr := ""
	for n, r := range rr {
		repr += fmt.Sprintf("\nId: %d%s", n, r.String())
	}
	return repr
}
