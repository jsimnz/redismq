package redismq

import (
	"fmt"
	"github.com/adeven/redis"
	"time"
)

//Consumers are Watchers that have writing commands
type Consumer struct {
	Broker
}

func (self *Queue) AddConsumer(name string) (c *Consumer, err error) {
	c = &Consumer{Broker{Name: name, Queue: self}}
	//check uniqueness and start heartbeat
	added := self.redisClient.SAdd(self.WorkerKey(), name).Val()
	if added == 0 {
		val := self.redisClient.Get(c.HeartbeatName()).Val()
		if val == "ping" {
			return nil, fmt.Errorf("consumer with this name is already active!")
		}
	}
	c.StartHeartbeat()
	return c, nil
}

func (self *Broker) HeartbeatName() string {
	return self.WorkingName() + "::heartbeat"
}

func (self *Consumer) StartHeartbeat() {
	go func() {
		for {
			self.Queue.redisClient.SetEx(self.HeartbeatName(), 1, "ping")
			time.Sleep(500 * time.Millisecond)
		}
	}()
}

func (self *Consumer) parseRedisAnswer(answer *redis.StringReq) (*Package, error) {
	if answer.Err() != nil {
		return nil, answer.Err()
	}
	p, err := UnmarshalPackage(answer.Val(), self.Queue, self)
	if err != nil {
		return nil, err
	}
	return p, nil
}
