package uuid

import "game_server/stack"

type UUID struct {
	available_uuids *stack.Stack
	highest_index   int
}

func (ui *UUID) Next() int {

	if ui.available_uuids.Len() > 0 {
		return ui.available_uuids.Pop().(int)
	} else {
		ui.highest_index++
		return ui.highest_index
	}
}

func (ui *UUID) Free(id int) {
	ui.available_uuids.Push(id)
}

func New() *UUID {
	return &UUID{stack.New(), -1}
}
