package list

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestList(t *testing.T) {
	lst := NewList()
	item := lst.Find(func(head interface{}) bool {
		return false
	})
	require.Nil(t, item)
	datas := []string{"data1", "data2", "data3"}
	ids := []string{RandStringBytes(64), RandStringBytes(64), RandStringBytes(64)}

	lst.Push(ids[0], datas[0])
	lst.Push(ids[1], datas[1])
	lst.Push(ids[2], datas[2])
	require.Equal(t, 3, lst.Size)

	item = lst.Find(func(head interface{}) bool {
		return head.(string) == datas[2]
	})
	require.Equal(t, datas[2], item)
	d := lst.Get(ids[0])
	require.Equal(t, datas[0], d)
	d = lst.Get(ids[1])
	require.Equal(t, datas[1], d)
	d = lst.Get(ids[2])
	require.Equal(t, datas[2], d)

	lst.Push(ids[0], datas[0])
	lst.Push(ids[1], datas[1])
	lst.Push(ids[2], datas[2])

	d = lst.First()
	require.Equal(t, datas[0], d)
	d = lst.First()
	require.Equal(t, datas[1], d)
	d = lst.First()
	require.Equal(t, datas[2], d)
	require.Equal(t, 0, lst.Size)
}
