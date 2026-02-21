package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLoadConfigDefault(t *testing.T) {
	expected := &Config{
		Redis: RedisCfg{
			Addr: DEFAULT_REDIS_ADDR,
		},
		Server: ServerCfg{
			Addr:     DEFAULT_SERVER_ADDR,
			NodeAddr: DEFAULT_NODE_ADDR,
		},
		PostgreSQL: PostgreSQLCfg{
			Addr: DEFAULT_POSTGRESQL_ADDR,
		},
	}

	actual := LoadConfig()

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("config mismatch (-want +got):%s\n", diff)
	}
}
