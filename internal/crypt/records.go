package crypt

import (
	"github.com/tim3-p/gophkeeper/internal/common"
)

// EncryptRecord encrypts sensitive fields in record
func EncryptRecord(key common.Key, a common.Record) (common.Record, error) {
	e := common.Record{
		Name: a.Name,
		Type: a.Type,
	}
	eOpaque, err := EncryptString(key, a.Opaque)
	if err != nil {
		return e, err
	}
	eMeta, err := EncryptString(key, a.Meta)
	if err != nil {
		return e, err
	}
	e.Opaque = eOpaque
	e.Meta = eMeta
	return e, nil
}

// DecryptRecord decrypts sensitive fields in record
func DecryptRecord(key common.Key, e common.Record) (common.Record, error) {
	a := common.Record{
		Name: e.Name,
		Type: e.Type,
	}
	Opaque, err := DecryptString(key, e.Opaque)
	if err != nil {
		return e, err
	}
	Meta, err := DecryptString(key, e.Meta)
	if err != nil {
		return e, err
	}
	a.Opaque = Opaque
	a.Meta = Meta
	return a, nil
}
