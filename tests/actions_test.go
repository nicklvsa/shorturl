package tests

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/nicklvsa/shorturl/actions"
	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/tests/utils"
)

func TestCreateURLMapping(t *testing.T) {
	client, mock := redismock.NewClientMock()

	cfg := &shared.Config{
		DB:            client,
		MetricsConfig: &shared.MetricsConfig{},
	}

	act := actions.NewActions(cfg)
	act.Ctx = context.TODO()

	mock.Regexp().ExpectSet("(short::*)\\w+", "123::http://hello.world", 0).SetVal("1")

	_, err := act.CreateURLMapping("http://hello.world", "123", nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteShortURL(t *testing.T) {
	client, mock := redismock.NewClientMock()

	cfg := &shared.Config{
		DB:            client,
		MetricsConfig: &shared.MetricsConfig{},
	}

	act := actions.NewActions(cfg)
	act.Ctx = context.TODO()

	mock.ExpectGet("short::cool123").SetVal("abc::abc")
	mock.ExpectDel("short::cool123").SetVal(1)

	if err := act.DeleteShortURL("cool123", "abc"); err != nil {
		t.Fatal(err)
	}
}

func TestIncrShortURLCount(t *testing.T) {
	client, mock := redismock.NewClientMock()

	cfg := &shared.Config{
		DB: client,
		MetricsConfig: &shared.MetricsConfig{
			Periods: map[string]string{
				"12h": "helloworld",
			},
		},
	}

	act := actions.NewActions(cfg)
	act.Ctx = context.TODO()

	mock.ExpectIncrBy("helloworld::someshortid", 1).SetVal(1)

	if err := act.IncrShortURLCount("someshortid"); err != nil {
		t.Fatal(err)
	}
}

func TestGetShortURLMetrics(t *testing.T) {
	client, mock := redismock.NewClientMock()

	cfg := &shared.Config{
		DB: client,
		MetricsConfig: &shared.MetricsConfig{
			Periods: map[string]string{
				"12h": "helloworld",
			},
		},
	}

	act := actions.NewActions(cfg)
	act.Ctx = context.TODO()

	mock.ExpectGet("short::someshortid").SetVal("abc::abc")

	data, err := act.GetShortURLMetrics("someshortid", "abc")
	if err != nil {
		t.Fatal(err)
	}

	utils.AssertInt(data["helloworld"], 0, t)
}

func TestGetLongURL(t *testing.T) {
	client, mock := redismock.NewClientMock()

	cfg := &shared.Config{
		DB:            client,
		MetricsConfig: &shared.MetricsConfig{},
	}

	act := actions.NewActions(cfg)
	act.Ctx = context.TODO()

	mock.ExpectGet("short::someshortid").SetVal("abc::example.com")

	longURL, err := act.GetLongURL("someshortid")
	if err != nil {
		t.Fatal(err)
	}

	utils.AssertStr(longURL, "example.com", t)
}
