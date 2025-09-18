package streamdeck

import (
	"fmt"
)

const DialMax = 100

type EventKind int

const (
	EventUnknown = iota
	EventKeyPressed
	EventKeyReleased
	EventDialPressed
	EventDialReleased
	EventDialTurn
)

func (ev EventKind) String() string {
	switch ev {
	case EventKeyPressed:
		return "key-pressed"
	case EventKeyReleased:
		return "key-released"
	case EventDialPressed:
		return "dial-pressed"
	case EventDialReleased:
		return "dial-released"
	case EventDialTurn:
		return "dial-turn"
	default:
		return "unknown"
	}
}

type Event struct {
	Kind  EventKind
	Which int
}

func (e Event) String() string {
	return fmt.Sprintf("%s:%d", e.Kind.String(), e.Which)
}

type State struct {
	Keys     []bool
	DialPush []bool
	DialPos  []int // 0 -> 100
}

func (s *State) Update(c *Config, b []byte) ([]Event, error) {
	if b[0] != 1 {
		return nil, fmt.Errorf("why isn't it starting with 1, %v", b)
	}

	// see https://github.com/dh1tw/streamdeck/pull/9#discussion_r2187628307
	if c != nil && c.ConvertKey {
		return s.updateKeyPressOriginal(b[1:])
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
		return nil, fmt.Errorf("unknown event type %d", b[1])
	}
}

func applyBools(in []bool, data []byte) ([]int, []bool) {
	changed := []int{}
	for i, b := range data {
		if len(in) <= i {
			in = append(in, false)
		}

		prev := in[i]

		if b == 0 {
			in[i] = false
		} else {
			in[i] = true
		}

		if prev != in[i] {
			changed = append(changed, i)
		}
	}
	return changed, in
}

func (s *State) updateKeyPressOriginal(data []byte) ([]Event, error) {
	changedKeys := []int{}
	updateEvents := []Event{}

	if len(data) < 15 {
		return nil, fmt.Errorf("wrong amount of data for updateKeyPressOriginal %d", len(data))
	}
	nd := make([]byte, 15)
	for x := 0; x < 3; x++ {
		for y := 0; y < 5; y++ {
			nd[(x*5)+y] = data[(x*5)+4-y]
		}
	}

	changedKeys, s.Keys = applyBools(s.Keys, nd)
	for _, changedKey := range changedKeys {
		if s.Keys[changedKey] {
			updateEvents = append(updateEvents, Event{EventKeyPressed, changedKey})
		} else {
			updateEvents = append(updateEvents, Event{EventKeyReleased, changedKey})
		}
	}
	return updateEvents, nil
}

func (s *State) updateKeyPress(data []byte) ([]Event, error) {
	changedKeys := []int{}
	updateEvents := []Event{}

	changedKeys, s.Keys = applyBools(s.Keys, data)
	for _, changedKey := range changedKeys {
		if s.Keys[changedKey] {
			updateEvents = append(updateEvents, Event{EventKeyPressed, changedKey})
		} else {
			updateEvents = append(updateEvents, Event{EventKeyReleased, changedKey})
		}
	}
	return updateEvents, nil
}

func (s *State) updateDialPush(data []byte) ([]Event, error) {
	var changedDialPushs []int
	changedDialPushs, s.DialPush = applyBools(s.DialPush, data)

	updateEvents := []Event{}
	for _, changedKey := range changedDialPushs {
		if s.DialPush[changedKey] {
			updateEvents = append(updateEvents, Event{EventDialPressed, changedKey})
		} else {
			updateEvents = append(updateEvents, Event{EventDialReleased, changedKey})
		}
	}
	return updateEvents, nil
}

func (s *State) updateDialTurn(data []byte) ([]Event, error) {
	var changedDialTurns int
	updateEvents := []Event{}
	changedDialTurns, s.DialPos = applyDelta(s.DialPos, data)

	if changedDialTurns >= 0 {
		updateEvents = append(updateEvents, Event{EventDialTurn, changedDialTurns})
		return updateEvents, nil
	}
	return nil, nil
}

func applyDelta(in []int, data []byte) (int, []int) {
	changed := -1

	for i, d := range data {
		if d != 0 {
			changed = i
		}
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
	return changed, in
}
