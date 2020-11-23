package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func instruct(text string) {
	countdown := exec.Command("say", text)
	err := countdown.Run()
	if err != nil {
		log.Printf("Failed to say: %s. Err: %v", text, err)
	} else {
		log.Printf("%s", text)
	}
}

func countdown() {
	for i := 3; i > 0; i-- {
		instruct(fmt.Sprintf("%d", i))
	}
	time.Sleep(1 * time.Second)
}

type Section interface {
	start()
	print()
}

type Segment struct {
	Duration int
	Foreword string
}

func (s Segment) start() {
	instruct(s.Foreword)
	timer1 := time.NewTimer(time.Duration(s.Duration) * time.Minute)
	<-timer1.C
	countdown()
}

func (s Segment) print() {
	fmt.Printf("%s (%d min)\n", s.Foreword, s.Duration)
}

type Tabata struct {
	Foreword string
	Count    int
	RestTime int
	WorkTime int
}

func (t Tabata) start() {
	instruct(t.Foreword)
	time.Sleep(10 * time.Second)
	for i := 0; i < t.Count; i++ {
		instruct("Start")
		workTimer := time.NewTimer(time.Duration(t.WorkTime)*time.Second - 3)
		<-workTimer.C
		countdown()
		instruct("Pihenő")
		restTimer := time.NewTimer(time.Duration(t.RestTime)*time.Second - 3)
		<-restTimer.C
		countdown()
	}
}

func (t Tabata) print() {
	fmt.Printf("%s %d x ( %ds work + %ds rest )\n", t.Foreword, t.Count, t.WorkTime, t.RestTime)
}

type WorkoutPlan struct {
	sections []Section
}

func (wp *WorkoutPlan) validateParameters(text string) bool {
	matched, err := regexp.Match(`^((s\d+|t) )+\s*$`, []byte(text+" "))
	if err != nil {
		log.Fatalf("Failed to validate parameters. Err: %v\n", err)
	}
	return matched
}
func (wp *WorkoutPlan) parseParameters(params string) []Section {
	parts := strings.Split(params, " ")
	result := []Section{}
	for i, s := range parts {
		sectionType := s[0]
		if sectionType == 's' {
			duration, err := strconv.Atoi(s[1:])
			if err != nil {
				log.Printf("Failed to parse duration. Err: %v\n", err)
				duration = 10
			}
			result = append(result,
				Segment{
					Foreword: fmt.Sprintf("%d. feladatsor", i),
					Duration: duration,
				},
			)
		} else if sectionType == 't' {
			result = append(result,
				Tabata{
					Foreword: fmt.Sprintf("Izometria"),
					RestTime: 10,
					WorkTime: 20,
					Count:    10,
				},
			)
		}
	}
	return result
}
func (wp *WorkoutPlan) init() {
	skipWarmup := flag.Bool("skip-warmup", false, "Do not start the workout with warmup")
	skipStretch := flag.Bool("skip-stretch", false, "Do not end the workout with warmup")
	flag.Parse()

	middleTasks := []Section{Segment{Foreword: "Feladatok", Duration: 40}}
	args := flag.Args()
	if len(args) > 0 {
		params := args[0]
		if !wp.validateParameters(params) {
			log.Fatalf(`Invalid parameters provided. '%v'.

Tabata('t') and Segments ('s') sections can be defined as the first string parameter. These sections must be separated by a space character.

Tabata:
  not configurable. 
Segments:
  must provide the duration in minutes without any unit or space between. 
			
e.g: 't s10 s20':
  0. Warmup
  1. A Tabata session
  2. 10 minutes for one task
  3. 20 minutes for another task
  Bonus: Stretching
`,
				os.Args[1])
		}
		middleTasks = wp.parseParameters(params)
	}

	wp.sections = []Section{}
	if !*skipWarmup {
		wp.sections = append(wp.sections, Segment{Foreword: "Bemelegítés", Duration: 10})
	}
	wp.sections = append(wp.sections, middleTasks...)
	if !*skipStretch {
		wp.sections = append(wp.sections, Segment{Foreword: "Nyújtás", Duration: 10})
	}

	fmt.Printf("Workut Plan\n\n")
	for _, s := range wp.sections {
		s.print()
	}
	fmt.Printf("\n////////////\n\n")
}

func (wp *WorkoutPlan) start() {
	for _, s := range wp.sections {
		s.start()
	}
}

func main() {
	workout := WorkoutPlan{}
	workout.init()
	workout.start()
	instruct("Vége")
}
