package config

import (
	"errors"
	"os"

	"github.com/tim3-p/gophkeeper/internal/common"
	"github.com/tim3-p/gophkeeper/internal/crypt"
)

// CheckFileMode returns true if the named file is not readable
// or writeable by anyone except the owner
func CheckFileMode(file string) (bool, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return false, err
	}
	perms := fileInfo.Mode().Perm()
	if (perms & 0077) != 0 {
		return false, nil
	}
	return true, nil
}

const minPhraseLen = 10

// GetKey returns a key computed from the given file content.
// It checks proper file permission and can check secret phrase strength.
func GetKey(file string) (*common.Key, error) {
	modeOK, err := CheckFileMode(file)
	if err != nil {
		return nil, err
	}
	if !modeOK {
		return nil, errors.New("key phrase file mode incorrect")
	}
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if len(buf) < minPhraseLen {
		return nil, errors.New("key phrase is too short")
	}
	key := crypt.MakeKey(string(buf))
	return &key, nil
}
