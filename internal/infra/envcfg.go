package infra

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

func LoadDatabaseCfg() (*DatabaseCfg, error) {
	var cfg DatabaseCfg
	prefix := "PG"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}

	return &cfg, nil
}

func LoadKafkaCfg() (*KafkaCfg, error) {
	var cfg KafkaCfg
	prefix := "KAFKA"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}

	return &cfg, nil
}

func LoadSMTPConfig() (*SMTPCfg, error) {
	var cfg SMTPCfg
	prefix := "SMTP"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}

	return &cfg, nil
}

func LoadEchoCfg() (*EchoCfg, error) {
	var cfg EchoCfg
	prefix := "APP"
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("%s: %w", prefix, err)
	}
	return &cfg, nil
}
