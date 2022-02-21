package worker

import (
	"fmt"

	"github.com/cocoonspace/fsm"
)

type State int

type Event int

const (
	Offline State = iota
	Idle
	PickingUp
	Delivering
	DroppingOff
)

func (s State) String() string {
	return [...]string{"Offline", "Idle", "PickingUp", "Delivering", "DroppingOff"}[s]
}

const (
	Ready Event = iota
	SignOff
	PickUp
	Deliver
	DropOff
)

func (s Event) String() string {
	return [...]string{"Ready", "SignOff", "PickUp", "Deliver", "DropOff"}[s]
}

type StateController struct {
	fcm *fsm.FSM
}

func (c *StateController) Ready() error {
	res := c.fcm.Event(fsm.Event(Ready))
	return c.checkState(res, Ready, Idle)
}

func (c *StateController) SignOff() error {
	res := c.fcm.Event(fsm.Event(SignOff))
	return c.checkState(res, SignOff, Offline)
}

func (c *StateController) PickUp() error {
	res := c.fcm.Event(fsm.Event(PickUp))
	return c.checkState(res, PickUp, PickingUp)
}

func (c *StateController) Deliver() error {
	res := c.fcm.Event(fsm.Event(Deliver))
	return c.checkState(res, Deliver, Delivering)
}

func (c *StateController) DropOff() error {
	res := c.fcm.Event(fsm.Event(DropOff))
	return c.checkState(res, DropOff, DroppingOff)
}

func (c *StateController) checkState(res bool, event Event, expected State) error {
	current := State(c.fcm.Current())
	if !res {
		return fmt.Errorf("state transition failed for event=%v current=%v", event, current)
	}
	if current != expected {
		return fmt.Errorf("state check failed for event=%v current=%v expected=%v", event, current, expected)
	}
	return nil
}

func NewStateController(initialState State) *StateController {
	f := fsm.New(fsm.State(initialState))
	fmt.Printf("initial state=%v\n", State(f.Current()))
	f.Transition(fsm.On(fsm.Event(Ready)), fsm.Src(fsm.State(Offline)), fsm.Dst(fsm.State(Idle)))
	f.Transition(fsm.On(fsm.Event(SignOff)), fsm.Src(fsm.State(Idle)), fsm.Dst(fsm.State(Offline)))
	f.Transition(fsm.On(fsm.Event(PickUp)), fsm.Src(fsm.State(Idle)), fsm.Dst(fsm.State(PickingUp)))
	f.Transition(fsm.On(fsm.Event(Deliver)), fsm.Src(fsm.State(PickingUp)), fsm.Dst(fsm.State(Delivering)))
	f.Transition(fsm.On(fsm.Event(DropOff)), fsm.Src(fsm.State(Delivering)), fsm.Dst(fsm.State(DroppingOff)))
	f.Transition(fsm.On(fsm.Event(Ready)), fsm.Src(fsm.State(DroppingOff)), fsm.Dst(fsm.State(Idle)))

	f.Enter(func(state fsm.State) {
		fmt.Printf("enter state=%v\n", State(state))
	})
	f.Exit(func(state fsm.State) {
		fmt.Printf("exit state=%v\n", State(state))
	})

	return &StateController{fcm: f}
}
