package action

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/tim3-p/gophkeeper/cmd/client/internal/config"
	"github.com/tim3-p/gophkeeper/internal/client"
	"github.com/tim3-p/gophkeeper/internal/common"
	"github.com/tim3-p/gophkeeper/internal/crypt"
)

const defaultFileMode = 0600

func getRecord(clnt *client.Client) (common.Record, error) {
	var eRecord common.Record
	var err error
	if config.Op.RecordID != 0 {
		eRecord, err = clnt.GetRecordByID(config.Op.RecordID)
	} else {
		eRecord, err = clnt.GetRecordByTypeName(
			config.Op.RecordType,
			config.Op.RecordName,
		)
	}
	if err != nil {
		return eRecord, err
	}
	record, err := crypt.DecryptRecord(*config.Key, eRecord)
	if err != nil {
		return record, err
	}
	return record, nil
}

func mergeRecord(newRecord, oldRecord common.Record) common.Record {
	if !config.Op.RecordChange.Name {
		newRecord.Name = oldRecord.Name
	}
	if !config.Op.RecordChange.Opaque {
		newRecord.Opaque = oldRecord.Opaque
	}
	if !config.Op.RecordChange.Meta {
		newRecord.Meta = oldRecord.Meta
	}
	return newRecord
}

func actRecord(subop config.OpSubtype, subrecord common.Opaque) error {
	clnt := client.NewClient(config.Cfg.ServerAddr,
		config.Cfg.UserName,
		config.Cfg.Password,
		config.Cfg.CacheFile,
		config.Cfg.HTTPSInsecure,
	)
	switch subop {
	case config.OpSubtypeRecordStore:
		record := common.Record{
			Name: config.Op.RecordName,
			Type: config.Op.RecordType,
			Meta: config.Op.RecordMeta,
		}

		err := subrecord.Check()
		if err != nil {
			return err
		}

		opaque, err := subrecord.Pack()
		if err != nil {
			return err
		}
		record.Opaque = string(opaque)

		eRecord, err := crypt.EncryptRecord(*config.Key, record)
		if err != nil {
			return err
		}
		id, err := clnt.StoreRecord(eRecord)
		if err != nil {
			return err
		}
		fmt.Printf("record stored with id %d\n", id)
	case config.OpSubtypeRecordGet:
		record, err := getRecord(clnt)
		if err != nil {
			return err
		}

		fmt.Println(record)

		if config.Op.RecordType == common.BinaryRecord {
			err = writeDecodeFile(config.Op.FileName, record.Opaque)
			if err != nil {
				return err
			}
			fmt.Printf("  File %s is written\n", config.Op.FileName)
		}
	case config.OpSubtypeRecordList:
		records, err := clnt.ListRecordsByType(config.Op.RecordType)
		if err != nil {
			return err
		}
		fmt.Println(records)
	case config.OpSubtypeRecordUpdate:
		record := common.Record{
			Name: config.Op.RecordName,
			Type: config.Op.RecordType,
			Meta: config.Op.RecordMeta,
		}

		// check opaque fields for validity if their change is requested
		if config.Op.RecordChange.Opaque {
			err := subrecord.Check()
			if err != nil {
				return err
			}
		}

		// pack provided opaque fields
		opaque, err := subrecord.Pack()
		if err != nil {
			return err
		}
		record.Opaque = string(opaque)

		// use server-side data for unchanged name, meta or opaque
		if !config.Op.RecordChange.Name || !config.Op.RecordChange.Opaque || !config.Op.RecordChange.Meta {
			serverRecord, err := getRecord(clnt)
			if err != nil {
				return err
			}
			record = mergeRecord(record, serverRecord)
		}

		eRecord, err := crypt.EncryptRecord(*config.Key, record)
		if err != nil {
			return err
		}

		if config.Op.RecordID != 0 {
			err = clnt.UpdateRecordByID(config.Op.RecordID, eRecord)
		} else {
			err = clnt.UpdateRecordByTypeName(config.Op.RecordType,
				config.Op.RecordName,
				eRecord,
			)
		}
		if err != nil {
			return err
		}
		fmt.Println("record updated")
	case config.OpSubtypeRecordDelete:
		if config.Op.RecordID != 0 {
			err := clnt.DeleteRecordByID(config.Op.RecordID)
			if err != nil {
				return err
			}
			fmt.Printf("Record %d deleted\n", config.Op.RecordID)
		} else {
			err := clnt.DeleteRecordByTypeName(
				config.Op.RecordType,
				config.Op.RecordName,
			)
			if err != nil {
				return err
			}
			fmt.Printf("Record %s deleted\n", config.Op.RecordName)
		}
	}
	return nil
}

func writeDecodeFile(file, str string) error {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, data, defaultFileMode)
	if err != nil {
		return err
	}
	return nil
}
