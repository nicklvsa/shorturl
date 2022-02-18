package shared

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/nicklvsa/shorturl/shared/db"
	"github.com/nicklvsa/shorturl/shared/logger"
)

type MetricsConfig struct {
	CollectAllTime bool              `json:"collect_all_time"`
	Periods        map[string]string `json:"periods"`
}

type Config struct {
	DB            *redis.Client
	MetricsConfig *MetricsConfig
}

func (m MetricsConfig) GetMetricPeriods() (map[string]time.Duration, error) {
	data := make(map[string]time.Duration)

	for period, name := range m.Periods {
		if _, found := data[name]; !found {
			periodLen := len(period) - 1
			format := period[periodLen]

			value, err := strconv.Atoi(period[:periodLen])
			if err != nil {
				return nil, err
			}

			durFormat := time.Hour
			switch format {
			case 'd':
				durFormat = time.Hour * 24
			case 'h':
				durFormat = time.Hour
			case 'm':
				durFormat = time.Minute
			}

			data[name] = time.Duration(value) * durFormat
		}
	}

	if m.CollectAllTime {
		data["total_count"] = time.Duration(0)
	}

	logger.Infof("%+v", data)

	return data, nil
}

func InitConfig() (*Config, error) {
	// initialize our redis db
	db, err := db.InitRedis(context.Background())
	if err != nil {
		return nil, err
	}

	metricsBytes, err := os.ReadFile("metrics-config.json")
	if err != nil {
		return nil, err
	}

	var metricsCfg MetricsConfig
	if err := json.Unmarshal(metricsBytes, &metricsCfg); err != nil {
		return nil, err
	}

	return &Config{
		DB:            db,
		MetricsConfig: &metricsCfg,
	}, nil
}
