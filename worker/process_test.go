package worker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	type payload struct {
		test string
	}
	type payloadOut struct {
		test string
	}
	type env struct {
		static string
	}
	procFunc := func(ei interface{}, pli interface{}) (interface{}, error) {
		e := ei.(*env)
		pl := pli.(*payload)
		plo := &payloadOut{test: pl.test + e.static}
		return plo, nil
	}
	e := &env{static: "_env"}
	proc := NewProcess(e, procFunc)
	proc.Run()
	pl := &payload{test: "test1"}
	taskID, err := proc.Send(pl)
	require.NoError(t, err)
	expectPlO := &payloadOut{
		test: pl.test + e.static,
	}
	for {
		time.Sleep(time.Millisecond)
		plo, err := proc.Check(taskID)
		require.NoError(t, err)
		if plo != nil {
			require.Equal(t, expectPlO, plo)
			break
		}
	}
}
