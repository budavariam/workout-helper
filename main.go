package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

func countdown() {
	for i := 3; i > 0; i-- {
		countdown := exec.Command("say", fmt.Sprintf("%d", i))
		err := countdown.Run()
		log.Printf("Failed to say number: %d. Err: %v", i, err)
	}
}

type Segment struct {
	Duration int
	Foreword string
}

func segment(segment Segment) {
	foreword := exec.Command("say", segment.Foreword)
	err := foreword.Run()
	log.Printf("Error during segment foreword. Err: %v", err)
	timer1 := time.NewTimer(time.Duration(segment.Duration) * time.Second)
	<-timer1.C
	countdown()
}

func main() {
	segment(Segment{Foreword: "Bemelegítés", Duration: 10})
	segment(Segment{Foreword: "Feladatok", Duration: 40})
	segment(Segment{Foreword: "Nyújtás", Duration: 10})
}
