// Package crypto provides cryptographic utilities including UUID generation and parsing.
package crypto

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
)

// UUID represents a UUID v4
type UUID struct {
	uuid.UUID
}

// NewUUID generates a new UUID v4
func NewUUID() (*UUID, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	return &UUID{UUID: u}, nil
}

// ParseUUID parses a string into a UUID
func ParseUUID(s string) (*UUID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return nil, err
	}
	return &UUID{UUID: u}, nil
}

// IsNil returns true if the UUID is nil
func (u *UUID) IsNil() bool {
	return u.UUID == uuid.Nil
}

// String returns the string representation of the UUID
func (u *UUID) String() string {
	return u.UUID.String()
}

// Scan implements the sql.Scanner interface for UUID
func (u *UUID) Scan(src interface{}) error {
	if src == nil {
		*u = UUID{UUID: uuid.Nil}
		return nil
	}

	switch s := src.(type) {
	case []byte:
		parsed, err := uuid.ParseBytes(s)
		if err != nil {
			return err
		}
		*u = UUID{UUID: parsed}
		return nil
	case string:
		parsed, err := uuid.Parse(s)
		if err != nil {
			return err
		}
		*u = UUID{UUID: parsed}
		return nil
	default:
		return fmt.Errorf("cannot scan %T into UUID", src)
	}
}

// Value implements the driver.Valuer interface for UUID
func (u UUID) Value() (driver.Value, error) {
	if u.IsNil() {
		return nil, nil
	}
	return u.String(), nil
}

// UUIDScan implements the sql.Scanner interface for UUID (legacy support)
func UUIDScan(u *UUID, src interface{}) error {
	return u.Scan(src)
}

// UUIDValue implements the driver.Valuer interface for UUID (legacy support)
func UUIDValue(u UUID) (driver.Value, error) {
	return u.Value()
}
