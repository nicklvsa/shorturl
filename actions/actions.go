package actions

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/shared/errs"
	"github.com/nicklvsa/shorturl/shared/logger"
	"github.com/rs/xid"
)

type Actions struct {
	Config *shared.Config
	Ctx    context.Context
}

func NewActions(cfg *shared.Config) *Actions {
	return &Actions{
		Config: cfg,
		Ctx:    context.Background(),
	}
}

func randomizeString(str string) string {
	bytes := make([]byte, len(str))
	for i := range bytes {
		bytes[i] = str[rand.Intn(len(str))]
	}

	return string(bytes)
}

func generateShortURL(longURL, employeeID string) string {
	return randomizeString(xid.New().String())
}

func (a Actions) canAccess(key, employeeID string) bool {
	storedVal, err := a.Config.DB.Get(a.Ctx, key).Result()
	if err != nil {
		logger.Errorf(err.Error())
		return false
	}

	val := shared.ShortenDBVal(employeeID, strings.Split(storedVal, "::")[1])

	return storedVal == val
}

func (a Actions) keyExists(key string) bool {
	count, err := a.Config.DB.Exists(a.Ctx, key).Result()
	if err != nil {
		return false
	}

	return count > 0
}

func (a Actions) setMetrics(shortID string) error {
	periods, err := a.Config.MetricsConfig.GetMetricPeriods()
	if err != nil {
		return err
	}

	for name, period := range periods {
		key := fmt.Sprintf("%s::%s", name, shortID)

		if !a.keyExists(key) {
			if err := a.setMetric(name, shortID, period); err != nil {
				return err
			}
		}
	}

	return nil
}

func (a Actions) setMetric(name, shortID string, period time.Duration) error {
	key := fmt.Sprintf("%s::%s", name, shortID)
	logger.Infof("Setting ttl for %s to %s", key, period.String())

	if err := a.Config.DB.Set(a.Ctx, key, 0, period).Err(); err != nil {
		return err
	}

	return nil
}

func (a Actions) CreateURLMapping(longURL, employeeID string, expirationMins *int) (string, error) {
	shortID := generateShortURL(longURL, employeeID)

	// ensure a shorturl cannot be duplicated
	for {
		if val, _ := a.GetLongURL(shortID); val == "" {
			break
		}

		logger.Warnf("found colliding short id, regenerating...")
		shortID = generateShortURL(longURL, employeeID)
	}

	dur := 0 * time.Minute
	if expirationMins != nil {
		dur = time.Duration(*expirationMins) * time.Minute
	}

	key := shared.ShortenDBKey(shortID)
	val := shared.ShortenDBVal(employeeID, longURL)

	if err := a.Config.DB.Set(a.Ctx, key, val, dur).Err(); err != nil {
		logger.Errorf("unable to set url mapping! Error: %s", err.Error())
		return "", err
	}

	if err := a.setMetrics(shortID); err != nil {
		logger.Errorf("unable to set metrics! Error: %s", err.Error())
		return "", err
	}

	return shortID, nil
}

func (a Actions) DeleteShortURL(shortID, employeeID string) error {
	key := shared.ShortenDBKey(shortID)

	if !a.canAccess(key, employeeID) {
		return errors.New("unauthorized")
	}

	if err := a.Config.DB.Del(a.Ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (a Actions) IncrShortURLCount(shortID string) error {
	periods, err := a.Config.MetricsConfig.GetMetricPeriods()
	if err != nil {
		return err
	}

	for name := range periods {
		key := fmt.Sprintf("%s::%s", name, shortID)

		if err := a.Config.DB.IncrBy(a.Ctx, key, 1).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (a Actions) GetShortURLMetrics(shortID, employeeID string) (map[string]int, error) {
	shortKey := shared.ShortenDBKey(shortID)

	periods, err := a.Config.MetricsConfig.GetMetricPeriods()
	if err != nil {
		return nil, err
	}

	if !a.canAccess(shortKey, employeeID) {
		return nil, errs.UnauthorizedAPIError.Err()
	}

	data := make(map[string]int)

	for name, period := range periods {
		key := fmt.Sprintf("%s::%s", name, shortID)

		if a.keyExists(key) {
			logger.Infof("key %s exists", name)

			keyTotal, err := a.Config.DB.Get(a.Ctx, key).Result()
			if err != nil {
				return nil, err
			}

			keyNum, err := strconv.Atoi(keyTotal)
			if err != nil {
				return nil, err
			}

			data[name] = keyNum
		} else {
			logger.Infof("key %s does not exist", name)
			a.setMetric(name, shortID, period)
			data[name] = 0
		}
	}

	return data, nil
}

func (a Actions) GetLongURL(shortID string) (string, error) {
	key := shared.ShortenDBKey(shortID)

	result, err := a.Config.DB.Get(a.Ctx, key).Result()
	if err != nil {
		return "", err
	}

	return strings.Split(result, "::")[1], nil
}
