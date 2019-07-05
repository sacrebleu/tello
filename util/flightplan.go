package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)


type Plan struct {
	Index int
	Steps []string
	decoratedSteps [] string
}

func (plan * Plan) Render() [] string {
	return plan.decoratedSteps
}

// set the index to the next step of the plan
func (plan * Plan) Next() {
	if len(plan.Steps) < 1 {
		return
	}
	plan.Index = (plan.Index + 1) % len(plan.Steps)

	copy(plan.decoratedSteps, plan.Steps) // copy the original
	plan.decoratedSteps[plan.Index] = fmt.Sprintf("[%s](fg:green)", plan.Steps[plan.Index])
}

// readLines reads a whole file into memory
// and returns a slice of its lines.
func (app * Application) ReadPlan(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to read %s", path))
		return
	}
	defer file.Close()

	var lines []string
	var dlines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		l := scanner.Text()
		if !strings.HasPrefix(l, "#"){
			lines = append(lines, l)
			dlines = append(dlines,l)
		}
	}

	//lines[0] = fmt.Sprintf("[%s](fg:green)", lines[0])
	app.FlightPlan = Plan{Index: 0, Steps: lines, decoratedSteps: dlines}
}