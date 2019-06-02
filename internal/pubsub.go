package internal

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const DefaultSub = "enrichment-worker-test"
const DefaultTopic = "churn-enrichment-test"

func CreateTopicIfNotExists(c *pubsub.Client, topic string) (*pubsub.Topic, error) {
	ctx := context.Background()
	t := c.Topic(topic)
	ok, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		return t, nil
	}
	t, err = c.CreateTopic(ctx, topic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create non existent topic")
	}
	return t, nil
}

func CreateSub(client *pubsub.Client, name string, topic *pubsub.Topic) (*pubsub.Subscription, error) {
	ctx := context.Background()
	sub := client.Subscription(name)
	ok, err := sub.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if ok {
		log.Infof("Using existing sub %v", sub)
		return sub, nil
	}
	sub, err = client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic: topic,
	})
	if err != nil {
		return sub, err
	}
	log.Infof("Subscripted created %v", sub)
	return sub, nil
}
