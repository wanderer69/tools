package tree

import (
	"fmt"
)

type SymbolItem struct {
	treeItem *TreeItem
	payload  interface{}
}

type TreeItem struct {
	symbolItems [256]*SymbolItem
}

type Tree struct {
	Root       *TreeItem
	StackItems []*StackItem

	ResultItems []*ResultItem
}

func NewTree() *Tree {
	return &Tree{}
}

func (t *Tree) AddString(str string, payload interface{}) error {
	currentTI := t.Root
	isRoot := true
	var prevSI *SymbolItem
	var currentSI *SymbolItem
	// var parentLetter byte
	bstr := []byte(str)
	for i := range bstr {
		b := bstr[i]
		if currentTI == nil {
			currentTI = &TreeItem{}
			if isRoot {
				t.Root = currentTI
				isRoot = false
			} else {
				if prevSI != nil {
					prevSI.treeItem = currentTI
				}
			}
		} else {
			isRoot = false
		}
		currentSI = currentTI.symbolItems[b]
		if currentSI == nil {
			currentSI = &SymbolItem{}
			currentTI.symbolItems[b] = currentSI
		}
		prevSI = currentSI
		currentTI = currentSI.treeItem
	}
	if currentSI.payload != nil {
		return fmt.Errorf("duplicate %s", str)
	}
	currentSI.payload = payload
	return nil
}

func (t *Tree) CheckString(str string) (int, interface{}, bool) {
	currentTI := t.Root
	// var parentLetter byte
	var currentSI *SymbolItem
	bstr := []byte(str)
	for i := range bstr {
		b := bstr[i]
		//fmt.Printf("bstr %x\r\n", b)
		if currentTI == nil {
			return i, nil, false
		}
		currentSI = currentTI.symbolItems[b]
		if currentSI == nil {
			return i, nil, false
		}
		currentTI = currentSI.treeItem
	}
	return 0, currentSI.payload, true
}

type DictionaryItem struct {
	Item string
	Id   string
}

func (t *Tree) AddToDictionary(lst []*DictionaryItem) error {
	for i := range lst {
		err := t.AddString(lst[i].Item, lst[i])
		if err != nil {
			return err
		}
	}
	return nil
}

type StackItem struct {
	pos       int
	currentSI *SymbolItem
	currentTI *TreeItem
}

type ResultItem struct {
	Payload  interface{}
	BeginPos int
}

func remove(slice []*StackItem, i int) []*StackItem {
	//fmt.Printf("len(slice) %v, i %v\r\n", len(slice), i)
	if len(slice) == 0 {
		return slice
	}
	if len(slice) == 1 {
		return slice[:0]
	}
	copy(slice[i:], slice[i+1:])
	return slice[:len(slice)-1]
}

func (t *Tree) CheckBySymbol(s byte, pos int) {
	items := []int{}
	if len(t.StackItems) > 0 {
		for i := range t.StackItems {
			if t.StackItems[i].currentTI == nil {
				// delete item
				items = append(items, i)
			} else {
				t.StackItems[i].currentSI = t.StackItems[i].currentTI.symbolItems[s]
				if t.StackItems[i].currentSI == nil {
					// delete item
					items = append(items, i)
				} else {
					t.StackItems[i].currentTI = t.StackItems[i].currentSI.treeItem
				}
			}
			if t.StackItems[i].currentSI != nil {
				if t.StackItems[i].currentSI.payload != nil {
					ri := ResultItem{
						BeginPos: t.StackItems[i].pos,
						Payload:  t.StackItems[i].currentSI.payload,
					}
					t.ResultItems = append(t.ResultItems, &ri)
					items = append(items, i)
				}
			}
		}
	}

	for i := range items {
		t.StackItems = remove(t.StackItems, items[i]-i)
	}

	currentTI := t.Root
	if currentTI == nil {
		return
	}
	currentSI := currentTI.symbolItems[s]
	if currentSI == nil {
		return
	}
	currentTI = currentSI.treeItem

	if currentSI.payload != nil {
		ri := ResultItem{
			BeginPos: pos,
			Payload:  currentSI.payload,
		}
		t.ResultItems = append(t.ResultItems, &ri)
	}

	si := StackItem{
		pos:       pos,
		currentSI: currentSI,
		currentTI: currentTI,
	}
	t.StackItems = append(t.StackItems, &si)
}
