package action

import (
	"errors"
	"log"

	"github.com/tim3-p/gophkeeper/cmd/client/internal/config"
	"github.com/tim3-p/gophkeeper/internal/client"
	"github.com/tim3-p/gophkeeper/internal/common"
)

const minPasswordLen = 5

func checkPasswordRequirements(password string) bool {
	if len(password) < minPasswordLen {
		return false
	}
	return true
}

func actUser(subop config.OpSubtype, user common.User) error {
	clnt := client.NewClient(config.Cfg.ServerAddr,
		config.Cfg.UserName,
		config.Cfg.Password,
		config.Cfg.CacheFile,
		config.Cfg.HTTPSInsecure,
	)
	switch subop {
	case config.OpSubtypeUserRegister:
		id, err := clnt.RegisterUser(config.Cfg.FullName)
		if err != nil {
			return err
		}
		log.Printf("user is registered with id %d", id)
	case config.OpSubtypeUserVerify:
		err := clnt.VerifyUser()
		if err != nil {
			return err
		}
		log.Printf("user is verified")
	case config.OpSubtypeUserPasswordChange:
		ok := checkPasswordRequirements(config.Op.User.Password)
		if !ok {
			return errors.New("password does not match requirements")
		}
		err := clnt.ChangePassword(config.Op.User)
		if err != nil {
			return err
		}
		log.Printf("password is changed")
	}
	return nil
}

func actCache(subop config.OpSubtype) error {
	clnt := client.NewClient(config.Cfg.ServerAddr,
		config.Cfg.UserName,
		config.Cfg.Password,
		config.Cfg.CacheFile,
		config.Cfg.HTTPSInsecure,
	)
	switch subop {
	case config.OpSubtypeCacheClean:
		err := clnt.CleanCache()
		if err != nil {
			return err
		}
		log.Println("cache is cleaned")
	case config.OpSubtypeCacheSync:
		if err := clnt.SyncCacheByType(common.AccountRecord); err != nil {
			return err
		}
		if err := clnt.SyncCacheByType(common.NoteRecord); err != nil {
			return err
		}
		if err := clnt.SyncCacheByType(common.CardRecord); err != nil {
			return err
		}
		if err := clnt.SyncCacheByType(common.BinaryRecord); err != nil {
			return err
		}
		log.Println("cache is synchronized")
	}
	return nil
}

// ChooseAct performs client actions
func ChooseAct() error {
	switch config.Op.Op {
	case config.OpTypeUser:
		return actUser(config.Op.Subop, config.Op.User)
	case config.OpTypeCache:
		return actCache(config.Op.Subop)
	case config.OpTypeAccount:
		return actRecord(config.Op.Subop, config.Op.Account)
	case config.OpTypeNote:
		return actRecord(config.Op.Subop, config.Op.Note)
	case config.OpTypeCard:
		return actRecord(config.Op.Subop, config.Op.Card)
	case config.OpTypeBinary:
		return actRecord(config.Op.Subop, config.Op.Binary)
	default:
		return errors.New("unknown operation type")
	}
}
