package generator

import (
	"fmt"
	"math/rand"
	"time"
)

type Coordinates struct {
	width     float64
	longitude float64
}

type telematic struct {
	timestamp time.Time
	speed     int
	coords    *Coordinates
}

func randomTimestamp() time.Time {
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
	randomNow := time.Unix(randomTime, 0)
	return randomNow
}

func generate(c chan *telematic) {
	//var tList []*telematic
	prevTimestamp := randomTimestamp()
	prevCoords := &Coordinates{-90 + rand.Float64()*180, -180 + rand.Float64()*360} //res[i] = min + rand.Float64() * (max - min)
	for i := 0; i < 10; i++ {
		newTimestamp := prevTimestamp.Add(time.Duration(rand.Intn(60)) * time.Second)
		newSpeed := rand.Intn(120)
		//fmt.Println(newTimestamp.Unix(), newTimestamp.Unix()/3600)
		//r := (newTimestamp.Unix()/3600 - prevTimestamp.Unix()/3600) * int64(newSpeed)
		//fmt.Println(prevTimestamp, newTimestamp, r)
		var newCoords *Coordinates
		//x := -float64(r) + rand.Float64()*(float64(r)*2)
		//y := math.Pow(float64(r), 2) - math.Pow(x, 2)
		//fmt.Println(r, x, y)
		newCoords = &Coordinates{prevCoords.width, prevCoords.longitude} // это заглушка, пересчитать все
		//if r < 1 {
		//	newCoords = &coordinates{prevCoords.width + (-90 + rand.Float64()*180), -180 + rand.Float64()*360}
		//} else {
		//	newCoords = &coordinates{-90 + rand.Float64()*180, prevCoords.longitude + (-180 + rand.Float64()*360)}
		//}
		//fmt.Printf("prevCoords: %v, prevTime: %v, newCoords: %v, newTime: %v\n", prevCoords, prevTimestamp, newCoords, newTimestamp)
		t := &telematic{
			timestamp: newTimestamp,
			speed:     newSpeed,
			coords:    newCoords,
		}
		c <- t
		//tList = append(tList, t)
		prevTimestamp = newTimestamp
		prevCoords = newCoords
	}
	close(c)
	//return tList
}

func Generator() {
	kafkaCh := make(chan *telematic)
	fmt.Println("created")
	go generate(kafkaCh)
	for {
		val, ok := <-kafkaCh
		if ok == false {
			fmt.Println(val, ok, "<-- loop broke!")
			break // exit break loop
		} else {
			fmt.Println(val, ok)
		}
	}
	//for i := range kafkaCh {
	//	fmt.Println(i)
	//}
	//return nil
}
