package rmq_test

import (
	"context"
	rmqClient "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/bang-go/crab/core/mq/rmq"
	"log"
	"testing"
	"time"
)

var (
	Topic        = ""
	GroupName    = ""
	Endpoint     = ""
	NameSpace    = ""
	AccessKey    = ""
	AccessSecret = ""
)

func TestConsumer(t *testing.T) {
	consumer, err := rmq.NewSimpleConsumer(&rmq.Config{
		NameSpace:     NameSpace,
		Endpoint:      Endpoint,
		ConsumerGroup: GroupName,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    AccessKey,
			AccessSecret: AccessSecret,
		},
	},
		rmqClient.WithAwaitDuration(10*time.Second),
		rmqClient.WithSubscriptionExpressions(map[string]*rmqClient.FilterExpression{
			Topic: rmqClient.SUB_ALL,
		}))
	if err != nil {
		log.Fatal(err)
	}
	_ = consumer.Start()
	for {
		mvs, err := consumer.Receive(context.Background(), 16, 20*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		for _, mv := range mvs {
			consumer.Ack(context.Background(), mv)
			log.Println(mv)
		}
	}
}
