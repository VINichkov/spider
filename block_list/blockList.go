package block_list

import (
	"strconv"
	"sync"
)

type BlockList struct {
	list map[string]bool
	sync.Mutex
}

func NewBlockList() *BlockList {
	return &BlockList{
		make(map[string]bool),
		sync.Mutex{},
	}
}

func (list *BlockList) Push(code int) {
	str := strconv.Itoa(code)
	list.Lock()
	list.list[str] = true
	list.Unlock()
}


func (list *BlockList) Include(code int) bool {
	str := strconv.Itoa(code)
	list.Lock()
	result := list.list[str]
	list.Unlock()
	return result
}
