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
		return "key-push"
	case EventKeyReleased:
		return "key-unpush"
	case EventDialPressed:
		return "dial-push"
	case EventDialReleased:
		return "dial-unpush"
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

func (s *State) Update(c *Config, b []byte) (Event, error) {
	if b[0] != 1 {
		return Event{}, fmt.Errorf("why isn't it starting with 1, %v", b)
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
		return Event{}, fmt.Errorf("unknown event type %d", b[1])
	}
}

func applyBools(in []bool, data []byte) (int, []bool) {
	changed := -1
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
			changed = i
		}
	}
	return changed, in
}

func (s *State) updateKeyPress(data []byte) (Event, error) {
	var changed int
	changed, s.Keys = applyBools(s.Keys, data)
	if changed >= 0 {
		if s.Keys[changed] {
			return Event{EventKeyPressed, changed}, nil
		}
		return Event{EventKeyReleased, changed}, nil
	}
	return Event{EventUnknown, changed}, nil
}

func (s *State) updateDialPush(data []byte) (Event, error) {
	var changed int
	changed, s.DialPush = applyBools(s.DialPush, data)
	if changed >= 0 {
		if s.DialPush[changed] {
			return Event{EventDialPressed, changed}, nil
		}
		return Event{EventDialReleased, changed}, nil
	}
	return Event{EventUnknown, changed}, nil
}

func (s *State) updateDialTurn(data []byte) (Event, error) {
	var changed int
	changed, s.DialPos = applyDelta(s.DialPos, data)
	if changed >= 0 {
		return Event{EventDialTurn, changed}, nil
	}
	return Event{EventUnknown, changed}, nil
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
