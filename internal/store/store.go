package store

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"os"
	"sync"

	"github.com/tim3-p/gophkeeper/internal/common"
	// sqlite sql package
	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultDBFile = "secret_storage.db"
)

// ErrNotFound is to indicate the absence of the record
var ErrNotFound = errors.New("Record not found")

// ErrAlreadyExists is to indicate the record already exist
var ErrAlreadyExists = errors.New("Entity already exists")

// Store is the secret storage
type Store struct {
	db     *sql.DB
	dbFile string
	mutex  sync.Mutex
}

// CloseDB closes database
func (s *Store) CloseDB() error {
	if s == nil {
		return errors.New("store is nill")
	}
	err := s.db.Close()
	if err != nil {
		return err
	}
	return nil
}

// DropStore removes the storage completely. Use for tests
func DropStore(dbFile string) error {
	err := os.Remove(dbFile)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// NewStore initializes new storage or opens existing one
func NewStore(storeName string) (*Store, error) {
	secretStore := &Store{}
	var err error

	secretStore.db, err = sql.Open("sqlite3", storeName)
	if err != nil {
		return secretStore, err
	}
	err = secretStore.db.Ping()
	if err != nil {
		return secretStore, err
	}

	_, err = secretStore.db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		user TEXT NOT NULL UNIQUE CHECK (length(user) >= 3),
		full_name TEXT,
		password_hash TEXT
	)`)
	if err != nil {
		return secretStore, err
	}

	_, err = secretStore.db.Exec(`CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY,
		user_id INTEGER NOT NULL,
		name TEXT NOT NULL CHECK (length(name) >= 1),
		type TEXT NOT NULL,
		opaque TEXT,
		meta TEXT,
		UNIQUE(user_id,name,type),
		FOREIGN KEY (user_id)
		  REFERENCES users (id)
		    ON DELETE CASCADE
		    ON UPDATE NO ACTION
	)`)
	if err != nil {
		return secretStore, err
	}

	return secretStore, nil
}

func (s *Store) isUserExists(userName string) (bool, error) {
	row := s.db.QueryRow(
		`SELECT count(*) FROM users WHERE user = ?`,
		userName,
	)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	if count == 1 {
		return true, nil
	}
	return false, nil
}

// CheckUserAuth checks the user password match
func (s *Store) CheckUserAuth(userName string, userPass string) (bool, error) {

	row := s.db.QueryRow(
		`SELECT password_hash FROM users WHERE user = ?`,
		userName,
	)

	var dbPasswordHash string
	err := row.Scan(&dbPasswordHash)
	if err == sql.ErrNoRows {
		return false, ErrNotFound
	}
	if err != nil {
		return false, err
	}
	hash := sha256.Sum256([]byte(userPass))
	passwordHash := hex.EncodeToString(hash[:])

	if dbPasswordHash == passwordHash {
		return true, nil
	}

	return false, nil
}

// AddUser creates user account
func (s *Store) AddUser(user common.User) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hash := sha256.Sum256([]byte(user.Password))
	passwordHash := hex.EncodeToString(hash[:])

	userExists, err := s.isUserExists(user.Name)
	if err != nil {
		return 0, err
	}

	if userExists {
		return 0, ErrAlreadyExists
	}

	res, err := s.db.Exec(`INSERT INTO users
		(user, full_name, password_hash)
		VALUES(?, ?, ?)`,
		user.Name,
		user.FullName,
		passwordHash,
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// ChangeUserPassword takes username and new passwords and replace password
func (s *Store) ChangeUserPassword(user, newPass string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	hash := sha256.Sum256([]byte(newPass))
	passwordHash := hex.EncodeToString(hash[:])

	_, err := s.db.Exec(`UPDATE users
		SET password_hash = ?
		where user = ?`,
		passwordHash, user,
	)
	if err != nil {
		return err
	}
	return nil
}
