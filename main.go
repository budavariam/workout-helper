package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func instruct(text string) {
	log.Printf("Instruct START: %v", time.Now())
	countdown := exec.Command("say", text)
	err := countdown.Run()
	if err != nil {
		log.Printf("Failed to say: %s. Err: %v", text, err)
	} else {
		log.Printf("%s", text)
	}
	log.Printf("Instruct END: %v", time.Now())
}

func countdown() {
	for i := 3; i > 0; i-- {
		instruct(fmt.Sprintf("%d", i))
	}
	time.Sleep(1 * time.Second)
}

type Section interface {
	start()
}

type Segment struct {
	Duration int
	Foreword string
}

func (s Segment) start() {
	instruct(s.Foreword)
	timer1 := time.NewTimer(time.Duration(s.Duration) * time.Second)
	<-timer1.C
	countdown()
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

func getMiddlePart() []Section {
	result := []Section{Segment{Foreword: "Feladatok", Duration: 40}}
	if len(os.Args) > 0 {
		params := os.Args[1]
		parts := strings.Split(params, " ")
		result = []Section{}
		for i, s := range parts {
			sectionType := s[0]
			if sectionType == 's' {
				duration, err := strconv.Atoi(s[1:])
				if err != nil {
					log.Printf("Failed to parse duration. Err: %v", err)
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
						Count:    5,
					},
				)
			}
		}
	}
	return result
}

func main() {
	middlePart := getMiddlePart()
	Segment{Foreword: "Bemelegítés", Duration: 10}.start()
	for _, s := range middlePart {
		s.start()
	}
	Segment{Foreword: "Nyújtás", Duration: 10}.start()
	instruct("Vége")
}
