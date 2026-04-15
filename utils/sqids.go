package utils

import (
	"errors"

	"github.com/sqids/sqids-go"
)

var uidGenerator11 *sqids.Sqids

func init() {
	uidGenerator11, _ = sqids.New(sqids.Options{
		Alphabet:  "0123456789",
		MinLength: 11,
		Blocklist: make([]string, 0),
	})
}

func EncodeUId(uid int64) (string, error) {
	return uidGenerator11.Encode([]uint64{uint64(uid)})
}
func DecodeUId(uid string) (int64, error) {
	resArr := uidGenerator11.Decode(uid)
	if len(resArr) != 1 {
		return 0, errors.New("invalid uid")
	}
	number := resArr[0]
	if number < 0 {
		return 0, errors.New("invalid uid")
	}
	id, err := EncodeUId(int64(number))
	if err != nil {
		return 0, err
	}
	if uid != id {
		return 0, errors.New("invalid uid")
	}
	return int64(number), nil
}
