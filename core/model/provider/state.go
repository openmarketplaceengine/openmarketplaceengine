package provider

import (
	"fmt"

	"github.com/cocoonspace/fsm"
)

type State int

type Event int

const (
	Offline State = iota
	Online
	PickingUp
	Delivering
)

func (s State) String() string {
	return [...]string{"Offline", "Online", "PickingUp", "Delivering"}[s]
}

const (
	GoOnline Event = iota
	GoOffline
	PickUp
	Deliver
	CompleteDelivery
)

func (s Event) String() string {
	return [...]string{"GoOnline", "GoOffline", "PickUp", "Deliver", "CompleteDelivery"}[s]
}

type StateController struct {
	fcm *fsm.FSM
}

func (c *StateController) GoOnline() error {
	res := c.fcm.Event(fsm.Event(GoOnline))
	return c.checkState(res, GoOnline, Online)
}

func (c *StateController) GoOffline() error {
	res := c.fcm.Event(fsm.Event(GoOffline))
	return c.checkState(res, GoOffline, Offline)
}

func (c *StateController) PickUp() error {
	res := c.fcm.Event(fsm.Event(PickUp))
	return c.checkState(res, PickUp, PickingUp)
}

func (c *StateController) Deliver() error {
	res := c.fcm.Event(fsm.Event(Deliver))
	return c.checkState(res, Deliver, Delivering)
}

func (c *StateController) CompleteDelivery() error {
	res := c.fcm.Event(fsm.Event(CompleteDelivery))
	return c.checkState(res, CompleteDelivery, Online)
}

func (c StateController) checkState(res bool, event Event, expected State) error {
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
	f.Transition(fsm.On(fsm.Event(GoOnline)), fsm.Src(fsm.State(Offline)), fsm.Dst(fsm.State(Online)))
	f.Transition(fsm.On(fsm.Event(GoOffline)), fsm.Src(fsm.State(Online)), fsm.Dst(fsm.State(Offline)))
	f.Transition(fsm.On(fsm.Event(PickUp)), fsm.Src(fsm.State(Online)), fsm.Dst(fsm.State(PickingUp)))
	f.Transition(fsm.On(fsm.Event(Deliver)), fsm.Src(fsm.State(PickingUp)), fsm.Dst(fsm.State(Delivering)))
	f.Transition(fsm.On(fsm.Event(CompleteDelivery)), fsm.Src(fsm.State(Delivering)), fsm.Dst(fsm.State(Online)))

	f.EnterState(fsm.State(Offline), func() {
		fmt.Println("going offline")
	})
	f.Enter(func(state fsm.State) {
		fmt.Printf("enter state=%v\n", State(state))
	})
	f.Exit(func(state fsm.State) {
		fmt.Printf("exit state=%v\n", State(state))
	})

	return &StateController{fcm: f}
}
