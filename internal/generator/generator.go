package generator

import (
	"math"
	"math/rand"
	"practice/storage"
	"time"
)

func RandomTimestamp() time.Time { // функция генерирует рандомное начальное время
	randomTime := rand.Int63n(time.Now().Unix()-94608000) + 94608000
	randomNow := time.Unix(randomTime, 0)
	return randomNow
}

func Generate(num int, c chan *storage.Car) {
	prevTimestamp := RandomTimestamp()
	//log.Println("time is", prevTimestamp)
	prevCoords := storage.Coordinates{-90 + rand.Float64()*180, -180 + rand.Float64()*360} //res[i] = min + rand.Float64() * (max - min)
	for {
		t := rand.Intn(60)
		newTimestamp := prevTimestamp.Add(time.Duration(t) * time.Second)
		newSpeed := rand.Intn(120)
		//log.Println(".............", t, newSpeed, newTimestamp)
		var s float64
		s = (float64(newSpeed) * float64(t)) / 3600 / 111
		var newCoords storage.Coordinates
		x := -float64(s) + rand.Float64()*(float64(s)*2)                                   // рандомное изменение широты
		y := math.Sqrt(math.Pow(float64(s), 2) - math.Pow(x, 2))                           // изменение долготы
		newCoords = storage.Coordinates{prevCoords.Latitude + x, prevCoords.Longitude + y} // это заглушка, пересчитать все

		car := &storage.Car{0, num, newSpeed, newCoords, newTimestamp}
		c <- car
		prevTimestamp = newTimestamp
		prevCoords = newCoords
		//log.Println("coords", newCoords)
		//time.Sleep(1 * time.Second)
	}
}
