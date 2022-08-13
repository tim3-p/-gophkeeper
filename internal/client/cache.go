package client

import "github.com/tim3-p/gophkeeper/internal/common"

// SyncCacheByType tries to cache all records of the given type
func (c *Client) SyncCacheByType(t common.RecordType) error {
	records, err := c.ListRecordsByType(t)
	if err != nil {
		return err
	}

	for id := range records {
		_, err := c.GetRecordByID(id)
		if err != nil {
			return err
		}
	}
	return nil
}

// CleanCache erases all cached records
func (c *Client) CleanCache() error {
	records, err := c.cacheListRecords()
	if err != nil {
		return err
	}

	for id := range records {
		err := c.cacheDeleteRecordByID(id)
		if err != nil {
			return err
		}
	}
	return nil
}
