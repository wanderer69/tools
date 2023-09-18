package list

import (
	"sync"
)

type ListItem struct {
	Head interface{}
	Prev *ListItem
	Next *ListItem
	ID   string
}

type List struct {
	Root   *ListItem
	End    *ListItem
	Mutex  *sync.Mutex
	Used   chan bool
	DictID map[string]*ListItem
	Size   int
}

func NewList() *List {
	lq := List{}
	lq.Root = nil
	lq.End = nil
	lq.Used = make(chan bool)
	lq.DictID = make(map[string]*ListItem)
	lq.Mutex = &sync.Mutex{}
	lq.Size = 0
	return &lq
}

func (lq *List) Push(id string, data interface{}) {
	li := ListItem{data, nil, nil, id}
	lq.Mutex.Lock()
	if lq.Root == nil {
		lq.Root = &li
		lq.End = &li
	} else {
		li.Prev = lq.End
		lq.End.Next = &li
		lq.End = &li
	}
	lq.DictID[id] = &li
	lq.Size = lq.Size + 1
	lq.Mutex.Unlock()
}

func (lq *List) Pop() (data interface{}) {
	lq.Mutex.Lock()
	if lq.Root == nil {
	} else {
		li := lq.End
		if li.Prev == nil {
			lq.End = nil
			lq.Root = nil
		} else {
			lq.End = li.Prev
			lq.End.Next = nil
		}
		data = li.Head
		lq.Size = lq.Size - 1
	}
	lq.Mutex.Unlock()
	return
}

func (lq *List) First() (data interface{}) {
	lq.Mutex.Lock()
	if lq.Root == nil {
	} else {
		li := lq.Root
		if li.Next == nil {
			lq.End = nil
			lq.Root = nil
		} else {
			lq.Root = li.Next
			lq.Root.Prev = nil
		}
		data = li.Head
		lq.Size = lq.Size - 1
	}
	lq.Mutex.Unlock()
	return
}

func (lq *List) Find(fn func(head interface{}) bool) (data interface{}) {
	lq.Mutex.Lock()
	if lq.Root == nil {
	} else {
		li := lq.Root
		for {
			if li == nil {
				break
			} else {
				if fn != nil {
					if fn(li.Head) {
						data = li.Head
						break
					}
				}
				li = li.Next
			}
		}
	}
	lq.Mutex.Unlock()
	return
}

func (lq *List) Get(id string) (data interface{}) {
	lq.Mutex.Lock()
	if lq.Root == nil {
	} else {
		if li, ok := lq.DictID[id]; ok {
			if li == lq.Root {
				li := lq.Root
				if li.Next == nil {
					lq.End = nil
					lq.Root = nil
				} else {
					lq.Root = li.Next
					lq.Root.Prev = nil
				}
			} else {
				if li == lq.End {
					li := lq.End
					if li.Prev == nil {
						lq.End = nil
						lq.Root = nil
					} else {
						lq.End = li.Prev
						lq.End.Next = nil
					}
				} else {
					li.Next.Prev = li.Prev
					li.Prev.Next = li.Next
				}
			}
			lq.Size = lq.Size - 1
			data = li.Head
		}
	}
	lq.Mutex.Unlock()
	return
}
