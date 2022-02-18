package tests

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/tests/utils"
)

var testConfig = `
	{
		"collect_all_time": false,
		"periods": {
			"1m": "past_minute",
			"24h": "past_day", 
			"7d": "past_week"
		}
	}
`

func TestConfig(t *testing.T) {
	var metricsCfg shared.MetricsConfig
	if err := json.Unmarshal([]byte(testConfig), &metricsCfg); err != nil {
		t.Fatal(err)
	}

	utils.AssertFalse(metricsCfg.CollectAllTime, t)
	utils.AssertStr(metricsCfg.Periods["7d"], "past_week", t)
	utils.AssertStr(metricsCfg.Periods["24h"], "past_day", t)
	utils.AssertStr(metricsCfg.Periods["1m"], "past_minute", t)
}

func TestGetPeriods(t *testing.T) {
	var metricsCfg shared.MetricsConfig
	if err := json.Unmarshal([]byte(testConfig), &metricsCfg); err != nil {
		t.Fatal(err)
	}

	durs, err := metricsCfg.GetMetricPeriods()
	if err != nil {
		t.Fatal(err)
	}

	utils.AssertDuration(durs["past_minute"], 1*time.Minute, t)
	utils.AssertDuration(durs["past_week"], 168*time.Hour, t)
	utils.AssertDuration(durs["past_day"], 24*time.Hour, t)
}
