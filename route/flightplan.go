package route

import (
	"fmt"
	"strings"
	"log"
)

// linked list
type Command struct {
	Description string
	Action string
	Offset int
	next *Command
}

type Plan struct {
	Name string
	Head *Command
	Tail *Command
}


// a flight plan needs to have a start and an end - start must be a takeoff at T0, finish must be a landing at TN, and there can be any number of
// steps T1 .. TN-1 in between.  So in data structure terms it's a linked list of commands, and the head must be takeoff and the tail must be land

func NewPlan(name string) * Plan {
	to := Command{Description:"Take off", Offset: 0, Action: "take_off" }
	fi := Command{Description:"Land", Offset: 300, Action: "land"}
	to.next = &fi

	// initialise the plan with a start and an end
	p := Plan{Name: name, Head: &to, Tail: &fi}

	return &p
}

// append a command to a flight plan - will always be appended just before the Tail
func Append(plan *Plan, cmd *Command) {

	var curr = plan.Head
	var prior = curr

	for curr.next != nil {
		prior = curr
		curr = curr.next
	}

	prior.next = cmd
	cmd.next = curr
}

func Tabulate(plan *Plan) {
	var output [] string

	log.Println("Flightplan: ",plan.Name)
	for curr := plan.Head ; curr != nil; {
		output = append(output, fmt.Sprintf("%3d %20s\t%s\n", curr.Offset, curr.Action, curr.Description))
		curr = curr.next
	}

	fmt.Print(strings.Join(output, "\n"))
}
