package unique

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	newbase64 "github.com/wanderer69/tools/new_base64"
)

var (
	_unique map[string]string
	_r      *rand.Rand
)

func init() {
	InitUniqueValue()
}

func InitUniqueValue() {
	_r = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	_unique = make(map[string]string)
}

func UniqueValue(len_n int) string {
	if false {
		var str string
		for {
			var bytes_array []byte
			for i := 0; i < len_n; i++ {
				bytes := _r.Intn(35)
				if bytes > 9 {
					bytes = bytes + 7
				}
				bytes_array = append(bytes_array, byte(bytes+16*3))
			}
			str = string(bytes_array)
			if _, ok := _unique[str]; !ok {
				_unique[str] = str
				break
			}
		}
		return str
	}
	return UniqueValueUUIDNew()
}

// 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000111111111111111111111111111111
// 0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999111111111122222222223333333333
// 0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789
// 01234567012345670123456701234567012345670123456701234567012345670123456701234567012345670123456701234567012345670123456701234567
// 012345012345012345012345012345012345012345012345012345012345012345012345012345012345012345012345012345012345012345012345012345012345
// 000000111111222222333333444444555555666666777777888888999999000000111111222222333333444444555555666666777777888888999999000000111111
// !                       !                       !                       !                       !                       !       !
func UniqueValueUUID() string {
	data := uuid.Must(uuid.NewRandom())
	//16 байт - 16*8
	result := base64.StdEncoding.EncodeToString(data[:])
	return result
}

func UniqueValueUUIDNew() string {
	data := uuid.Must(uuid.NewRandom())
	//16 байт - 16*8
	result := newbase64.BytesEncode(data[:]) // base64.StdEncoding.EncodeToString(data[:])

	bn, err := newbase64.BytesDecode(result)
	if err != nil {
		panic(err)
	}
	for i := range bn {
		if data[i] != bn[i] {
			panic(fmt.Errorf("not equeal %v %v", data[i], bn[i]))
		}
	}
	return string(result)
}
