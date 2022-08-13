package common

// Key is the AES key type used
type Key [32]byte

// User is the client of the secret store service
type User struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

// Account holds account data for some resource
type Account struct {
	URL      string `json:"url"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// Note holds text data
type Note struct {
	Text string `json:"text"`
}

// Card holds creditcard data
type Card struct {
	Holder   string `json:"holder"`
	Number   string `json:"number"`
	ExpMonth int    `json:"exp_month"`
	ExpYear  int    `json:"exp_year"`
	CVC      string `json:"cvc"`
}

// Binary holds base64-encoded binary data
type Binary struct {
	Data string `json:"data"`
}

// Record can hold any record that could be stored
type Record struct {
	Name   string     `json:"name"`
	Type   RecordType `json:"record_type"`
	Opaque string     `json:"opaque"`
	Meta   string     `json:"meta"`
}

// Records can hold the map of any record that could be stored
type Records map[int64]Record

// RecordType is the type of record conveyed
type RecordType string

const (
	// AccountRecord is the Account record type
	AccountRecord RecordType = "acc"
	// NoteRecord is the Note record type
	NoteRecord RecordType = "note"
	// CardRecord is the Card record type
	CardRecord RecordType = "card"
	// BinaryRecord is the Binary record type
	BinaryRecord RecordType = "bin"
	// UnspecifiedRecord is the unspecified record type
	UnspecifiedRecord RecordType = ""
)

// AddUserResponse is the response for AddUser request
type AddUserResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	ID     int64  `json:"id"`
}

// StoreRecordResponse is the responce for store account
type StoreRecordResponse struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	ID     int64  `json:"id"`
}
