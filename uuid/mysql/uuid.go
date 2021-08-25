package mysql

import (
	"database/sql/driver"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"strings"
)

type MyUuid uuid.UUID

func NewMyUuid() MyUuid {
	return MyUuid(uuid.NewV4())
}

func (u *MyUuid) SetFromString(id string) error {
	uuidId, err := uuid.FromString(id)
	*u = MyUuid(uuidId)
	return err
}

func (u *MyUuid) SetFromByte(id []byte) error {
	uuidId, err := uuid.FromBytes(id)
	*u = MyUuid(uuidId)
	return err
}

func (u *MyUuid) String() string {
	uuidId := uuid.UUID(*u)
	return uuidId.String()
}

func (u *MyUuid) Bytes() []byte {
	uuidId := uuid.UUID(*u)
	return uuidId.Bytes()
}

func (u *MyUuid) IsEmpty() bool {
	uid := uuid.UUID(*u)
	check := uid.String()
	return check == "" || check == "00000000-0000-0000-0000-000000000000"
}

// Value implements the driver.Valuer interface.
func (u MyUuid) Value() (driver.Value, error) {
	uid := uuid.UUID(u)
	return uid.Bytes(), nil
}

// Scan implements the sql.Scanner interface.
// A 16-byte slice is handled by UnmarshalBinary, while
// a longer byte slice or a string is handled by UnmarshalText.
func (u *MyUuid) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		uid, _ := uuid.FromBytes(src)
		*u = MyUuid(uid)
		if len(src) == uuid.Size {
			return uid.UnmarshalBinary(src)
		}
		return uid.UnmarshalText(src)

	case string:
		uid, _ := uuid.FromString(src)
		*u = MyUuid(uid)
		return uid.UnmarshalText([]byte(src))
	default:
		fmt.Println("Undefined source type")
		return nil
	}

	return fmt.Errorf("uuid: cannot convert %T to UUID", src)
}

// MarshalJSON string json conversion from variable of uuid type
func (u MyUuid) MarshalJSON() ([]byte, error) {
	uid := uuid.UUID(u)
	return []byte(`"` + uid.String() + `"`), nil
}

// UnmarshalJSON variable of uuid type from string json conversion
func (u *MyUuid) UnmarshalJSON(b []byte) error {
	uid, err := uuid.FromString(strings.Trim(string(b), "\""))
	if err != nil {
		return err
	}
	*u = MyUuid(uid)

	return nil
}
