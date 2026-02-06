package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfigDefault(t *testing.T) {
	expected := &Config{
		RabbitMQ: RabbitMQCfg{
			Addr:          DEFAULT_RABBITMQ_ADDR,
			MessagesQueue: DEFAULT_MESSAGES_QUEUE,
		},
		Redis: RedisCfg{
			Addr: DEFAULT_REDIS_ADDR,
		},
		Server: ServerCfg{
			Addr:     DEFAULT_SERVER_ADDR,
			NodeAddr: DEFAULT_NODE_ADDR,
		},
	}

	actual := LoadConfig()

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("config mismatch (-want +got):%s\n", diff)
	}
}
