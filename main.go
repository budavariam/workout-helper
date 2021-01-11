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

const (
	text_start       = "Start"
	text_rest        = "Pihenő"
	text_warmup      = "Bemelegítés"
	text_stretch     = "Nyújtás"
	text_tabata      = "Izometria"
	text_tasks       = "Feladatok"
	text_nth_task    = "%d. feladatsor"
	text_last_minute = "1 perc"
)

func instruct(text string) {
	// Get the list of availlable voices with: 'say -v ?'
	readText := exec.Command("say", "-v", "Mariska", text)
	err := readText.Run()
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

func countSeconds(seconds int) {
	if seconds < 60 {
		workTimer := time.NewTimer(time.Duration(seconds-3) * time.Second)
		<-workTimer.C
		countdown()
	} else {
		timeUntilNotice := time.Duration(seconds-60) * time.Second
		timeAfterNotice := time.Duration(60-3) * time.Second
		workTimer := time.NewTimer(timeUntilNotice)
		<-workTimer.C
		instruct(text_last_minute)
		workTimer = time.NewTimer(timeAfterNotice)
		<-workTimer.C
		countdown()
	}
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
	countSeconds(s.Duration * 60)
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
		instruct(text_start)
		countSeconds(t.WorkTime)
		instruct(text_rest)
		countSeconds(t.RestTime)
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
					Foreword: fmt.Sprintf(text_nth_task, i+1),
					Duration: duration,
				},
			)
		} else if sectionType == 't' {
			result = append(result,
				Tabata{
					Foreword: fmt.Sprintf(text_tabata),
					RestTime: 10,
					WorkTime: 20,
					Count:    8,
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

	middleTasks := []Section{Segment{Foreword: text_tasks, Duration: 40}}
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
		wp.sections = append(wp.sections, Segment{Foreword: text_warmup, Duration: 10})
	}
	wp.sections = append(wp.sections, middleTasks...)
	if !*skipStretch {
		wp.sections = append(wp.sections, Segment{Foreword: text_stretch, Duration: 10})
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
