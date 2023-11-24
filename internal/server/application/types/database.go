package types

import (
	"encoding/json"
	"fmt"
	"github.com/go-faster/errors"
	"reflect"
	"strconv"
	"time"
)

type Id uint64
type TimeStamp time.Time
type FingerPrint map[string]string

func (id *Id) Scan(value any) error {
	switch vt := value.(type) {
	case int64:
		*id = Id(vt)
	default:
		return errors.New(fmt.Sprintf("can not convert %x to %x type", vt, "Id"))
	}

	return nil
}

func (id *Id) UnmarshalJSON(data []byte) (err error) {
	value, err := strconv.Atoi(string(data))
	if err != nil {
		fmt.Println("Can not convert this []byte to int")
	}

	*id = Id(value)

	return
}

//func (id Id) MarshalJSON() (b []byte, err error) {
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

func (fp *FingerPrint) Scan(value any) error {
	switch vt := value.(type) {
	case []byte:
		return json.Unmarshal(vt, &fp)
	default:
		return errors.New(fmt.Sprintf("can not convert %s to %s type", reflect.TypeOf(vt), "FingerPrint"))
	}
}
