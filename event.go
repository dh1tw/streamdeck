package streamdeck

import (
	"fmt"
)

const DialMax = 100

type State struct {
	Keys     []bool
	DialPush []bool
	DialPos  []int // 0 -> 100
}

func (s *State) Update(b []byte) error {
	if b[0] != 1 {
		return fmt.Errorf("why isn't it starting with 1, %v", b)
	}

	switch b[1] {
	case 0:
		return s.updateKeyPress(b[4:])
	case 3:
		if b[4] == 0 {
			return s.updateDialPush(b[5:])
		}
		return s.updateDialTurn(b[5:])
	default:
		return fmt.Errorf("unknown event type %d", b[1])
	}
}

func applyBools(in []bool, data []byte) []bool {
	for i, b := range data {
		if len(in) <= i {
			in = append(in, false)
		}

		if b == 0 {
			in[i] = false
		} else {
			in[i] = true
		}
	}
	return in
}

func (s *State) updateKeyPress(data []byte) error {
	s.Keys = applyBools(s.Keys, data)
	return nil
}

func (s *State) updateDialPush(data []byte) error {
	s.DialPush = applyBools(s.DialPush, data)
	return nil

}

func (s *State) updateDialTurn(data []byte) error {
	s.DialPos = applyDelta(s.DialPos, data)
	return nil
}

func applyDelta(in []int, data []byte) []int {

	for i, d := range data {
		if len(in) <= i {
			in = append(in, 50)
		}

		output := in[i]

		if d < 0x80 {
			output = min(DialMax, output+int(d))
		} else {
			output = max(0, output-(256-int(d)))
		}

		in[i] = output
	}
	return in
}
