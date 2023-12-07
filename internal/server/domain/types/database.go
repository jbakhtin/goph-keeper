package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/go-faster/errors"
)

type ID uint64
type TimeStamp time.Time
type FingerPrint map[string]string
type RefreshToken string

func (id *ID) Scan(value any) error {
	switch vt := value.(type) {
	case int64:
		*id = ID(vt)
	default:
		return errors.New(fmt.Sprintf("can not convert %x to %x type", vt, "ID"))
	}

	return nil
}

func (id *ID) UnmarshalJSON(data []byte) (err error) {
	value, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("Can not convert this []byte to int")
	}

	*id = ID(value)

	return
}

//func (id ID) MarshalJSON() (b []byte, err error) {
//	return json.Marshal(id)
//}

func (ts *TimeStamp) Scan(value any) error {
	switch vt := value.(type) {
	case time.Time:
		*ts = TimeStamp(vt)
	default:
		return errors.New(fmt.Sprintf("can not convert %x to %x type", vt, "TimeStamp"))
	}

	return nil
}

func (ts TimeStamp) Value() (driver.Value, error) {
	stringValue := time.Time(ts)
	return stringValue, nil
}

func (fp *FingerPrint) Scan(value any) error {
	switch vt := value.(type) {
	case []byte:
		return json.Unmarshal(vt, &fp)
	default:
		return errors.New(fmt.Sprintf("can not convert %s to %s type", reflect.TypeOf(vt), "FingerPrint"))
	}
}
