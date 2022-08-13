package store

import (
	"database/sql"

	"github.com/tim3-p/gophkeeper/internal/common"
	// sqlite sql package
	_ "github.com/mattn/go-sqlite3"
)

// StoreRecordWithID stores record with the ID specified
func (s *Store) StoreRecordWithID(id int64, user string, record common.Record) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, err := s.db.Exec(`INSERT INTO records
		(id, user_id, name, type, opaque, meta)
		VALUES(?, (SELECT id from users where user=?), ?, ?, ?, ?)`,
		id,
		user,
		record.Name,
		record.Type,
		record.Opaque,
		record.Meta,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetRecordID returns the ID of the record with the given type and name
func (s *Store) GetRecordID(user string, t common.RecordType, name string) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	row := s.db.QueryRow(
		`SELECT records.id
			FROM records JOIN users ON records.user_id = users.id
			WHERE users.user = ?
			AND records.type = ?
			AND records.name = ?`,
		user, t, name,
	)

	var id int64
	err := row.Scan(&id)
	if err == sql.ErrNoRows {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}

// StoreRecord stores Record data for given user
func (s *Store) StoreRecord(user string, record common.Record) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	res, err := s.db.Exec(`INSERT INTO records
		(user_id, name, type, opaque, meta)
		VALUES((SELECT id from users where user=?), ?, ?, ?, ?)`,
		user,
		record.Name,
		record.Type,
		record.Opaque,
		record.Meta,
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

// UpdateRecordByID updates an record by ID
func (s *Store) UpdateRecordByID(user string, id int64, record common.Record) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	res, err := s.db.Exec(`UPDATE records
		SET name = ?, type = ?, opaque = ?, meta = ?
		WHERE id in
		( SELECT records.id FROM records
			JOIN users ON records.user_id = users.id
			WHERE users.user = ? AND records.id = ?)`,
		record.Name,
		record.Type,
		record.Opaque,
		record.Meta,
		user,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrNotFound
	}

	return nil
}

// UpdateRecordByTypeName updates an record by type and name
func (s *Store) UpdateRecordByTypeName(user string, t common.RecordType, name string, record common.Record) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	res, err := s.db.Exec(`UPDATE records
		SET name = ?, type = ?, opaque = ?, meta = ?
		WHERE id in
		( SELECT records.id FROM records
			JOIN users ON records.user_id = users.id
			WHERE users.user = ?
			AND records.type = ?
			and records.name = ?)`,
		record.Name,
		record.Type,
		record.Opaque,
		record.Meta,
		user,
		t,
		name,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return ErrNotFound
	}

	return nil
}

// ListRecords returns list of stored records for the given user
// with name and type fields filled
func (s *Store) ListRecords(user string) (common.Records, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	records := make(common.Records)
	rows, err := s.db.Query(
		`SELECT records.id, records.type, records.name
			FROM records JOIN users ON records.user_id = users.id
			WHERE users.user = ?`,
		user,
	)
	if err != nil {
		return records, err
	}

	for rows.Next() {
		var id int64
		var record common.Record
		err = rows.Scan(&id, &record.Type, &record.Name)
		if err != nil {
			return records, err
		}
		records[id] = record
	}
	return records, nil
}

// ListRecordsByType returns list of stored records of the given type
// for the given user with name and type fields filled
func (s *Store) ListRecordsByType(user string, t common.RecordType) (common.Records, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	records := make(common.Records)
	rows, err := s.db.Query(
		`SELECT records.id, records.name
			FROM records JOIN users ON records.user_id = users.id
			WHERE users.user = ? AND records.type = ?`,
		user, t,
	)
	if err != nil {
		return records, err
	}

	for rows.Next() {
		var id int64
		var record common.Record
		record.Type = t
		err = rows.Scan(&id, &record.Name)
		if err != nil {
			return records, err
		}
		records[id] = record
	}
	return records, nil
}

// GetRecordByID returns stored record by ID
func (s *Store) GetRecordByID(user string, id int64) (common.Record, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var record common.Record

	row := s.db.QueryRow(
		`SELECT records.name, records.type, records.opaque, records.meta
			FROM records JOIN users ON records.user_id = users.id
			WHERE users.user = ? AND records.id = ?`,
		user, id,
	)

	err := row.Scan(&record.Name,
		&record.Type,
		&record.Opaque,
		&record.Meta,
	)
	if err == sql.ErrNoRows {
		return record, ErrNotFound
	}
	if err != nil {
		return record, err
	}
	return record, nil
}

// GetRecordByTypeName returns stored record by type and name
func (s *Store) GetRecordByTypeName(user string, t common.RecordType, name string) (common.Record, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var record common.Record

	row := s.db.QueryRow(
		`SELECT records.name, records.type, records.opaque, records.meta
			FROM records JOIN users ON records.user_id = users.id
			WHERE users.user = ?
			AND records.type = ?
			AND records.name = ?`,
		user, t, name,
	)

	err := row.Scan(&record.Name,
		&record.Type,
		&record.Opaque,
		&record.Meta,
	)
	if err == sql.ErrNoRows {
		return record, ErrNotFound
	}
	if err != nil {
		return record, err
	}
	return record, nil
}

// DeleteRecordByID deletes the specified record record by ID
func (s *Store) DeleteRecordByID(user string, id int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	res, err := s.db.Exec(
		`DELETE FROM records
		 	WHERE id in
			( SELECT records.id FROM records
				JOIN users ON records.user_id = users.id
				WHERE users.user = ? AND records.id = ?
			)`,
		user, id,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil
	}
	if rows != 1 {
		return ErrNotFound
	}
	return nil
}

// DeleteRecordByTypeName deletes the specified record record by type and name
func (s *Store) DeleteRecordByTypeName(user string, t common.RecordType, name string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	res, err := s.db.Exec(
		`DELETE FROM records
		 	WHERE id in
			( SELECT records.id FROM records
				JOIN users ON records.user_id = users.id
				WHERE users.user = ?
				AND records.type = ?
				AND records.name = ?
			)`,
		user, t, name,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return nil
	}
	if rows != 1 {
		return ErrNotFound
	}
	return nil
}
