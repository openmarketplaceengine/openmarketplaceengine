package itinerary

import "fmt"

type StepType int

const (
	goToLocation StepType = iota
	pickUpPassengers
	dropOffPassengers
)

type step struct {
	prev     *step
	next     *step
	id       string
	stepType StepType
}

type flow struct {
	head *step
	tail *step
}

func newFlow() *flow {
	return &flow{}
}

func (f *flow) addStep(stepID string, stepType StepType) {
	s := &step{
		next:     f.head,
		id:       stepID,
		stepType: stepType,
	}
	if f.head != nil {
		f.head.prev = s
	}
	f.head = s

	l := f.head
	for l.next != nil {
		l = l.next
	}
	f.tail = l
}

func (f *flow) addStepBefore(beforeID string, stepID string, stepType StepType) error {
	s := &step{
		id:       stepID,
		stepType: stepType,
	}

	if f.head == nil {
		return fmt.Errorf("flow has no steps")
	}
	start := f.head
	for start != nil {
		if start.id == beforeID {
			s.prev = start
			s.next = start.next
			start.next = s
		}
		start = start.next
	}
	f.tail = start
	return nil
}

func (f *flow) removeStep(stepID string) {
	step := f.head
	for step != nil {
		if step.id == stepID {
			step.prev.next = step.next
			step.next.prev = step.prev
		}
		step = step.next
	}

	f.tail = step
}

func (f *flow) dump() {
	step := f.head
	for step != nil {
		fmt.Printf("%+v ->", step.id)
		step = step.next
	}
	fmt.Println()
}
