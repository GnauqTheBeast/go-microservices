package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Kafka KafkaConfig
	HTTP  HTTPConfig
}

type KafkaConfig struct {
	Brokers       []string
	ConsumerGroup string
	Topics        struct {
		UserEvents    string
		Notifications string
	}
	Consumer struct {
		InitialOffset string
		Strategy      string
	}
	Producer struct {
		ReturnSuccesses bool
		ReturnErrors    bool
	}
}

type HTTPConfig struct {
	Port string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// Kafka defaults
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.consumerGroup", "notification-group")
	viper.SetDefault("kafka.topics.userEvents", "user-events")
	viper.SetDefault("kafka.topics.notifications", "notifications")
	viper.SetDefault("kafka.consumer.initialOffset", "oldest")
	viper.SetDefault("kafka.consumer.strategy", "roundrobin")
	viper.SetDefault("kafka.producer.returnSuccesses", true)
	viper.SetDefault("kafka.producer.returnErrors", true)

	// HTTP defaults
	viper.SetDefault("http.port", ":8080")
}
