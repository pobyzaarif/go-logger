package main

import (
	"time"

	"github.com/pobyzaarif/goutil/logger"
)

var log = logger.NewLog("ab")

func main() {
	log.Info("hello")
	// rand.Seed(time.Now().UnixNano())
	// timeDelay := rand.Intn(10)
	// fmt.Printf("-- delay %v seconds --\n", timeDelay)
	time.Sleep(time.Duration(4) * time.Second)

	mappp := map[string]interface{}{"any": "any"}
	log.InfoWithData("test map", mappp)

	// process something here
	log.TimerStart()
	// timeDelay = rand.Intn(5)
	// fmt.Printf("-- delay %v seconds --\n", timeDelay)
	time.Sleep(time.Duration(2) * time.Second)
	log.InfoWithData("test map", mappp)

	// process something here
	// timeDelay = rand.Intn(3)
	// fmt.Printf("-- delay %v seconds --\n", timeDelay)
	time.Sleep(time.Duration(1) * time.Second)
	log.InfoWithData("test map", mappp)
}
