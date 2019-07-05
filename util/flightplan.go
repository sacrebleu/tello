package util

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		l := scanner.Text()
		if !strings.HasPrefix(l, "#"){
			lines = append(lines, l)
		}
		//i++
	}
	//fmt.Println(lines)

	app.FlightPlan = lines
}