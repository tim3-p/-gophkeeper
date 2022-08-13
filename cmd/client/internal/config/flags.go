package config

import (
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/tim3-p/gophkeeper/internal/common"
)

type (
	// OpType is the operation type
	OpType int
	// OpSubtype is the operation subtype
	OpSubtype int
)

const (
	// OpTypeUser is for user operations
	OpTypeUser OpType = iota
	// OpTypeCache is for user operations
	OpTypeCache
	// OpTypeAccount is for account operations
	OpTypeAccount
	// OpTypeNote is for note operations
	OpTypeNote
	// OpTypeCard is for note operations
	OpTypeCard
	// OpTypeBinary is for binary operations
	OpTypeBinary
)

const (
	// OpSubtypeUserRegister is the user registration
	OpSubtypeUserRegister OpSubtype = iota
	// OpSubtypeUserVerify is the user auth verification
	OpSubtypeUserVerify
	// OpSubtypeUserPasswordChange is for changing the password
	OpSubtypeUserPasswordChange

	// OpSubtypeCacheSync is the cache sync
	OpSubtypeCacheSync OpSubtype = iota
	// OpSubtypeCacheClean is the cache cleaning
	OpSubtypeCacheClean

	// OpSubtypeRecordStore is the record creation
	OpSubtypeRecordStore
	// OpSubtypeRecordGet is the regord retrieval
	OpSubtypeRecordGet
	// OpSubtypeRecordList is the listing of records
	OpSubtypeRecordList
	// OpSubtypeRecordUpdate is the record update
	OpSubtypeRecordUpdate
	// OpSubtypeRecordDelete is the removal of the record
	OpSubtypeRecordDelete
	// OpSubtypeOther is unknown operation
	OpSubtypeOther
)

// RequestedChange indicates if name, opaque or meta of the record
// is requested to change
type RequestedChange struct {
	Name   bool
	Opaque bool
	Meta   bool
}

// Operation describes the current operation type
type Operation struct {
	Op           OpType
	Subop        OpSubtype
	User         common.User
	Account      common.Account
	accountFlags []string
	Note         common.Note
	noteFlags    []string
	Card         common.Card
	cardFlags    []string
	Binary       common.Binary
	binaryFlags  []string
	RecordChange RequestedChange
	RecordID     int64
	RecordName   string
	RecordMeta   string
	RecordType   common.RecordType
	FileName     string
}

func isFlagPassed(set *flag.FlagSet, name string) bool {
	found := false
	set.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Op describes the current operation
var Op Operation

// ErrUnknownMode returned when unknown operation mode is requested
var ErrUnknownMode = errors.New("unknown mode")

// Usage prints the usage info: flag format, etc.
func Usage(msg string) {
	if msg != "" {
		fmt.Println(msg)
	}
	fmt.Println("usage: 'client MODE -a ACTION flags'")
	fmt.Println("  where MODE is one of user, cache, acc, note, card or bin")
	fmt.Println("  run 'client MODE -h' for further help")
}

func checkChanges(set *flag.FlagSet, opaqueFlags []string) RequestedChange {
	var r RequestedChange
	r.Name = isFlagPassed(set, "n")
	for _, fl := range opaqueFlags {
		if isFlagPassed(set, fl) {
			r.Opaque = true
		}
	}
	r.Meta = isFlagPassed(set, "m")
	return r
}

func actionType(a *string) OpSubtype {
	switch *a {
	case "store":
		return OpSubtypeRecordStore
	case "get":
		return OpSubtypeRecordGet
	case "list":
		return OpSubtypeRecordList
	case "update":
		return OpSubtypeRecordUpdate
	case "delete":
		return OpSubtypeRecordDelete
	}
	return OpSubtypeOther
}

// ParseFlags parses cmd line arguments
func ParseFlags() error {
	userFlags := flag.NewFlagSet("user", flag.ExitOnError)
	cacheFlags := flag.NewFlagSet("cache", flag.ExitOnError)
	accFlags := flag.NewFlagSet(string(common.AccountRecord), flag.ExitOnError)
	noteFlags := flag.NewFlagSet(string(common.NoteRecord), flag.ExitOnError)
	cardFlags := flag.NewFlagSet(string(common.CardRecord), flag.ExitOnError)
	binFlags := flag.NewFlagSet(string(common.BinaryRecord), flag.ExitOnError)

	userAction := userFlags.String("a", "verify", "action: verify|register|password")
	userPass := userFlags.String("p", "", "new password")

	cacheAction := cacheFlags.String("a", "sync", "action: sync|clean")

	accAction := accFlags.String("a",
		"list",
		"action: list|store|get|update|delete",
	)
	accName := accFlags.String("n", "", "account name")
	// opaque flags
	accUserName := accFlags.String("u", "", "account user name")
	accPassword := accFlags.String("p", "", "account password")
	accURL := accFlags.String("l", "", "account URL")
	Op.accountFlags = []string{"u", "p", "l"}

	accMeta := accFlags.String("m", "", "account metainfo")
	accID := accFlags.Int64("i", 0, "account ID")

	noteAction := noteFlags.String("a",
		"list",
		"action: list|store|get|update|delete",
	)
	noteName := noteFlags.String("n", "", "note name")
	// opaque flags
	noteText := noteFlags.String("t", "", "note text")
	Op.noteFlags = []string{"t"}

	noteMeta := noteFlags.String("m", "", "note metainfo")
	noteID := noteFlags.Int64("i", 0, "note ID")

	cardAction := cardFlags.String("a",
		"list",
		"action: list|store|get|update|delete",
	)
	cardName := cardFlags.String("n", "", "card name")
	// opaque flags
	cardHolder := cardFlags.String("ch", "", "card holder")
	cardNumber := cardFlags.String("num", "", "card number")
	cardExpMonth := cardFlags.Int("em", 0, "card expiry month")
	cardExpYear := cardFlags.Int("ey", 0, "card expiry year")
	cardCVC := cardFlags.String("c", "", "card CVC code")
	Op.cardFlags = []string{"ch", "num", "em", "ey", "c"}

	cardMeta := cardFlags.String("m", "", "card metainfo")
	cardID := cardFlags.Int64("i", 0, "card ID")

	binAction := binFlags.String("a",
		"list",
		"action: list|store|get|update|delete",
	)
	binName := binFlags.String("n", "", "binary record name")
	// opaque flags
	binFile := binFlags.String("f", "", "file name")
	Op.binaryFlags = []string{"f"}

	binID := binFlags.Int64("i", 0, "binary record ID")

	if len(os.Args) < 2 {
		return errors.New("mode is not set")
	}

	switch os.Args[1] {
	case "user":
		userFlags.Parse(os.Args[2:])
	case "cache":
		cacheFlags.Parse(os.Args[2:])
	case string(common.AccountRecord):
		accFlags.Parse(os.Args[2:])
	case string(common.NoteRecord):
		noteFlags.Parse(os.Args[2:])
	case string(common.CardRecord):
		cardFlags.Parse(os.Args[2:])
	case string(common.BinaryRecord):
		binFlags.Parse(os.Args[2:])
	default:
		return ErrUnknownMode
	}

	if userFlags.Parsed() {
		Op.Op = OpTypeUser
		switch *userAction {
		case "verify":
			Op.Subop = OpSubtypeUserVerify
		case "register":
			Op.Subop = OpSubtypeUserRegister
		case "password":
			Op.Subop = OpSubtypeUserPasswordChange
		default:
			return errors.New("unknown user action")
		}
		Op.User.Password = *userPass
	} else if cacheFlags.Parsed() {
		Op.Op = OpTypeCache
		switch *cacheAction {
		case "sync":
			Op.Subop = OpSubtypeCacheSync
		case "clean":
			Op.Subop = OpSubtypeCacheClean
		}
	} else if accFlags.Parsed() {
		Op.Op = OpTypeAccount
		Op.RecordType = common.AccountRecord
		Op.Subop = actionType(accAction)

		Op.RecordName = *accName
		Op.Account.UserName = *accUserName
		Op.Account.Password = *accPassword
		Op.Account.URL = *accURL
		Op.RecordMeta = *accMeta
		Op.RecordID = *accID
		Op.RecordChange = checkChanges(accFlags, Op.accountFlags)
	} else if noteFlags.Parsed() {
		Op.Op = OpTypeNote
		Op.RecordType = common.NoteRecord
		Op.Subop = actionType(noteAction)

		Op.RecordName = *noteName
		Op.Note.Text = *noteText
		Op.RecordMeta = *noteMeta
		Op.RecordID = *noteID
		Op.RecordChange = checkChanges(noteFlags, Op.noteFlags)
	} else if cardFlags.Parsed() {
		Op.Op = OpTypeCard
		Op.RecordType = common.CardRecord
		Op.Subop = actionType(cardAction)

		Op.RecordName = *cardName
		Op.Card.Holder = *cardHolder
		Op.Card.Number = *cardNumber
		Op.Card.ExpMonth = *cardExpMonth
		Op.Card.ExpYear = *cardExpYear
		Op.Card.CVC = *cardCVC
		Op.RecordMeta = *cardMeta
		Op.RecordID = *cardID
		Op.RecordChange = checkChanges(cardFlags, Op.cardFlags)
	} else if binFlags.Parsed() {
		Op.Op = OpTypeBinary
		Op.RecordType = common.BinaryRecord
		Op.Subop = actionType(binAction)

		Op.RecordName = *binName
		var err error
		if Op.Subop == OpSubtypeRecordStore || Op.Subop == OpSubtypeRecordUpdate {
			Op.Binary.Data, err = readEncodeFile(*binFile)
		}
		Op.FileName = *binFile
		if err != nil {
			return err
		}
		Op.RecordID = *binID
		Op.RecordChange = checkChanges(binFlags, Op.binaryFlags)
	}

	return nil
}

func readEncodeFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(data)
	return str, nil
}
