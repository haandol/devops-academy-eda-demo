package consumer

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/haandol/devops-academy-eda-demo/pkg/config"
	"github.com/haandol/devops-academy-eda-demo/pkg/port/primaryport/consumerport"
	"github.com/haandol/devops-academy-eda-demo/pkg/util"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kzap"
)

type KafkaConsumer struct {
	client        *kgo.Client
	topic         string
	messageExpiry time.Duration
	handler       consumerport.HandlerFunc
	batchSize     int
}

const ConsumerTimeout = 30 * time.Second

func NewKafkaConsumer(cfg *config.Kafka, groupID, topic string) *KafkaConsumer {
	opts := buildConsumerOpts(cfg.Seeds, groupID, topic)
	if strings.Contains(cfg.Seeds[0], "9094") {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		log.Panic(err)
	}

	return &KafkaConsumer{
		client:        client,
		topic:         topic,
		messageExpiry: time.Duration(cfg.MessageExpirySec) * time.Second,
		batchSize:     cfg.BatchSize,
		handler:       nil,
	}
}

func buildConsumerOpts(seeds []string, group, topic string) []kgo.Opt {
	return []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.ConsumerGroup(group),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(),
		kgo.Balancers(kgo.CooperativeStickyBalancer()), // explicit default rebalancer
		kgo.FetchMaxWait(1 * time.Second),
		kgo.FetchMaxBytes(70 * 1024 * 1024), // 70MB
		kgo.AllowAutoTopicCreation(),        // TODO: only for the dev
		kgo.WithLogger(kzap.New(
			util.GetLogger().With("package", "consumer").Desugar(),
			kzap.Level(kgo.LogLevelWarn),
		)),
	}
}

func (c *KafkaConsumer) RegisterHandler(h consumerport.HandlerFunc) error {
	logger := util.GetLogger().With(
		"module", "KafkaConsumer",
		"func", "RegisterHandler",
	)
	logger.Info("Registering handler...")

	if c.handler != nil {
		logger.Error("handler already registered")
		return errors.New("handler already registered")
	}

	c.handler = h
	logger.Info("registered handler")

	return nil
}

// Consume - consume messages from Kafka and dispatch to handlers
func (c *KafkaConsumer) Consume(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "KafkaConsumer",
		"func", "Consume",
		"topic", c.topic,
	)
	logger.Infow("Consuming Topic", "topic", c.topic)

	// check initialized
	if c.handler == nil {
		return errors.New("handler not registered")
	}

	for {
		logger.Info("Polling...")
		ctx := context.Background()

		fetches := c.client.PollRecords(ctx, c.batchSize)
		if fetches.IsClientClosed() {
			return errors.New("kafka client closed")
		}
		if errs := fetches.Errors(); len(errs) > 0 {
			return errs[0].Err
		}

		if err := c.handleFetchesInOrder(ctx, &fetches); err != nil {
			return err
		}

		if err := c.client.CommitUncommittedOffsets(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (c *KafkaConsumer) handleFetchesInOrder(ctx context.Context, fetches *kgo.Fetches) error {
	logger := util.GetLogger().With(
		"module", "KafkaConsumer",
		"func", "handleFetchesInOrder",
		"topic", c.topic,
	)

	var errs []error
	fetches.EachRecord(func(record *kgo.Record) {
		key := string(record.Key)
		logger.Infow("Message received", "key", key)

		message := &consumerport.Message{
			Topic:     record.Topic,
			Key:       key,
			Value:     record.Value,
			Timestamp: record.Timestamp,
		}
		if c.messageExpiry > 0 && time.Since(record.Timestamp) > c.messageExpiry {
			logger.Warnw("message expired", "expirySec", c.messageExpiry, "key", key)
			return
		}

		if err := c.handler(ctx, message); err != nil {
			errs = append(errs, err)
		}
	})
	if len(errs) > 0 {
		return fmt.Errorf("%v", errs)
	}

	return nil
}

func (c *KafkaConsumer) Close(ctx context.Context) error {
	logger := util.GetLogger().With(
		"module", "KafkaConsumer",
		"func", "Close",
		"topic", c.topic,
	)
	logger.Info("Closing...")

	c.client.Close()
	return nil
}
