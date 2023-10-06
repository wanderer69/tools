package tree

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type Payload struct {
	symbol string
	id     string
}

func TestAddString(t *testing.T) {
	tree := NewTree()
	str := "abcdefg"
	str = "алия"
	p := &Payload{
		symbol: str,
		id:     uuid.NewString(),
	}
	err := tree.AddString(str, p)
	require.NoError(t, err)

	pos, payload, ok := tree.CheckString(str)
	fmt.Printf("pos %v, ok %v\r\n", pos, ok)
	require.True(t, ok)
	require.Equal(t, 0, pos)
	require.Equal(t, p, payload)

	require.True(t, false)
}

func TestDictionary(t *testing.T) {
	tree := NewTree()

	lst := []string{
		"аббат",
		"абак",
		"аббревиатура",
		"аудиенция",
	}
	dia := []*DictionaryItem{}
	for i := range lst {
		di := DictionaryItem{
			Item: lst[i],
			Id:   uuid.NewString(),
		}
		dia = append(dia, &di)
	}
	err := tree.AddToDictionary(dia)
	require.NoError(t, err)

	pos, payload, ok := tree.CheckString(lst[0])
	fmt.Printf("pos %v, ok %v pl %#v\r\n", pos, ok, payload.(*DictionaryItem))
	require.True(t, ok)
	require.Equal(t, 0, pos)
	//require.Equal(t, p, payload)

	str := "абак аббат аббревиатура аудиенция"
	bstr := []byte(str)
	for i := range bstr {
		//fmt.Printf("str[%v] %x\r\n", i, bstr[i])
		tree.CheckBySymbol(bstr[i], i)
	}
	for i := range tree.ResultItems {
		if tree.ResultItems[i].Payload != nil {
			payload := tree.ResultItems[i].Payload.(*DictionaryItem)
			fmt.Printf("pos %v pl %#v\r\n", tree.ResultItems[i].BeginPos, payload)
		} else {
			fmt.Printf("pos %v\r\n", tree.ResultItems[i].BeginPos)
		}
	}
	require.True(t, false)
}
