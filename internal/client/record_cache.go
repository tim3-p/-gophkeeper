package client

import (
	"github.com/tim3-p/gophkeeper/internal/common"
	"github.com/tim3-p/gophkeeper/internal/store"
)

func (c *Client) cacheDeleteRecordByID(id int64) error {
	if c.CacheFile == "" {
		return nil
	}
	err := c.Store.DeleteRecordByID(c.UserName, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheRecordWithID(storeID int64, record common.Record) error {
	if c.CacheFile == "" {
		return nil
	}

	_, err := c.Store.AddUser(common.User{
		Name: c.UserName,
	})
	if err != nil && err != store.ErrAlreadyExists {
		return err
	}

	err = c.Store.DeleteRecordByID(c.UserName, storeID)
	if err != nil && err != store.ErrNotFound {
		return err
	}

	err = c.Store.StoreRecordWithID(storeID, c.UserName, record)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) cacheGetRecordByID(id int64) (common.Record, error) {
	if c.CacheFile == "" {
		return common.Record{}, nil
	}
	return c.Store.GetRecordByID(c.UserName, id)
}

func (c *Client) cacheGetRecordID(t common.RecordType,
	name string,
) (int64, error) {
	if c.CacheFile == "" {
		return 0, nil
	}
	return c.Store.GetRecordID(c.UserName, t, name)
}

func (c *Client) cacheGetRecordByTypeName(t common.RecordType,
	name string,
) (common.Record, error) {
	if c.CacheFile == "" {
		return common.Record{}, nil
	}
	return c.Store.GetRecordByTypeName(c.UserName, t, name)
}

func (c *Client) cacheListRecords() (common.Records, error) {
	records := make(common.Records)
	if c.CacheFile == "" {
		return records, nil
	}
	return c.Store.ListRecords(c.UserName)
}

func (c *Client) cacheListRecordsByType(
	t common.RecordType,
) (common.Records, error) {
	records := make(common.Records)
	if c.CacheFile == "" {
		return records, nil
	}
	return c.Store.ListRecordsByType(c.UserName, t)
}
