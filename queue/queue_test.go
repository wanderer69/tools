package queue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := NewQueue()
	type payload struct {
		test string
	}
	tests := []string{"1", "2", "3"}
	for i := range tests {
		qi := &QueueItem{Value: &payload{test: tests[i]}}
		require.NoError(t, q.Push(qi))
	}
	n := 0
	qi, err := q.Get()
	require.NoError(t, err)
	for {
		if qi != nil {
			require.Equal(t, tests[n], qi.Value.(*payload).test)
		}
		n += 1
		fmt.Printf("n %v len(tests) %v qi %#v\r\n", n, len(tests), qi)
		qi, err = q.Next(qi)
		fmt.Printf("qi %#v err %v\r\n", qi, err)

		if n > len(tests) {
			require.Error(t, err)
			break
		}
		require.NoError(t, err)
	}
	for i := range tests {
		v, err := q.Get()
		require.NoError(t, err)
		require.Equal(t, tests[i], v.Value.(*payload).test)
		v, err = q.Pop()
		require.NoError(t, err)
		require.Equal(t, tests[i], v.Value.(*payload).test)
	}
	_, err = q.Get()
	require.Error(t, err)
	_, err = q.Pop()
	require.Error(t, err)
}
