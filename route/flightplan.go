package route

// linked list
type Command struct {
	Description string
	Action string
	Offset int
	Next *Command
}

type Plan struct {
	Name string
	Initial *Command
	Tail *Command
}


// a flight plan needs to have a start and an end - start must be a takeoff at T0, finish must be a landing at TN, and there can be any number of
// steps T1 .. TN-1 in between.  So in data structure terms it's a linked list of commands, and the head must be takeoff and the tail must be land

func NewPlan(name string) * Plan {
	to := Command{Description:"Take off", Offset: 0, Action: "take_off" }
	p := Plan{Name: name, Initial: &to }

	return &p
}

func append(plan *Plan, cmd *Command) {
	
}

// validate the flightplan contract so the tello doesn't get lost :P
//func validate(plan * Plan) {
//
//}