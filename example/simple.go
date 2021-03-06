package main

import (
	"fmt"
	"github.com/adeven/goenv"
	"github.com/adeven/redismq"
)

func main() {
	goenv := goenv.DefaultGoenv()
	testQueue := redismq.NewQueue(goenv, "clicks")
	for i := 0; i < 10; i++ {
		testQueue.Put("testpayload")
	}
	consumer, err := testQueue.AddConsumer("testconsumer")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		p, err := consumer.Get()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(p.CreatedAt)
		err = p.Ack()
		if err != nil {
			fmt.Println(err)
		}
	}
}
