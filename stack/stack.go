package stack

import (
	"bufio"
	"bytes"
	"regexp"
)

type Element struct {
	Func   string `json:"func,omitempty"`
	Source string `json:"source,omitempty"`
}

type Goroutine struct {
	Name     string    `json:"name,omitempty"`
	State    string    `json:"state,omitempty"`
	Elements []Element `json:"elements,omitempty"`
}

type Stack struct {
	Goroutines []Goroutine `json:"goroutines,omitempty"`
	Raw        []byte      `json:"-"`
}

var (
	goroutineRegexp = regexp.MustCompile(`^(.+?)\s\[(.+?)\]:$`)
	funcRegexp      = regexp.MustCompile(`^(\S.+)$`)
	sourceRegexp    = regexp.MustCompile(`^\s+(.+)$`)
)

const (
	goroutineNameGroup  = 1
	goroutineStateGroup = 2
	funcGroup           = 1
	sourceGroup         = 1
)

func Parse(stack []byte) (Stack, error) {
	res := Stack{
		Raw: stack,
	}

	currentGoroutine := -1 // Index of the current goroutine
	currentElement := -1   // Index of the current element

	s := bufio.NewScanner(bytes.NewReader(stack))

	for s.Scan() {
		line := s.Text()
		if groups := goroutineRegexp.FindStringSubmatch(line); groups != nil {
			// New goroutine
			res.Goroutines = append(res.Goroutines, Goroutine{
				Name:  groups[goroutineNameGroup],
				State: groups[goroutineStateGroup],
			})
			currentGoroutine = len(res.Goroutines) - 1
			currentElement = -1
		} else if currentGoroutine >= 0 {
			// We have at least 1 goroutine

			if groups := funcRegexp.FindStringSubmatch(line); groups != nil {
				res.Goroutines[currentGoroutine].Elements = append(res.Goroutines[currentGoroutine].Elements, Element{
					Func: groups[funcGroup],
				})
				currentElement = len(res.Goroutines[currentGoroutine].Elements) - 1
			} else if currentElement >= 0 {
				// We have at least 1 element

				if groups := sourceRegexp.FindStringSubmatch(line); groups != nil {
					res.Goroutines[currentGoroutine].Elements[currentElement].Source = groups[sourceGroup]
				}
			}
		}
	}

	return res, s.Err()
}
