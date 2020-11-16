package main

import (
	"fmt"
	"log"
	"os/exec"
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

type Segment struct {
	Duration int
	Foreword string
}

func (s Segment) play() {
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

func (t Tabata) play() {
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

func main() {
	Segment{Foreword: "Bemelegítés", Duration: 10}.play()
	Segment{Foreword: "Feladatok", Duration: 40}.play()
	//Tabata{Foreword: "Izometria", Count: 5, RestTime: 10, WorkTime: 20}.play()
	Segment{Foreword: "Nyújtás", Duration: 10}.play()
	instruct("Vége")
}
