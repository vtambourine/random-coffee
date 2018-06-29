package main

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
)

func main() {
	fmt.Println("Cron jobs")

	t := make(chan int)
	c := cron.New()
	c.AddFunc("@every 1s", func () {
		fmt.Print(".")
	})

	c.AddFunc("0 */1 10-14 * * FRI", func () {
		log.Println("\nit's time")
	})

	c.Start()

	<-t
}
