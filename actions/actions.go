package actions

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/nicklvsa/shorturl/shared"
	"github.com/nicklvsa/shorturl/shared/logger"
	"github.com/rs/xid"
)

type MetricsView struct {
	TotalCount    int `json:"total_count"`
	PastDayCount  int `json:"past_day_count"`
	PastWeekCount int `json:"past_week_count"`
}

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

func (a Actions) setMetrics(shortID string) error {
	totalCountKey := shared.TotalCountDBKey(shortID)
	weekKey := shared.WeekCountDBKey(shortID)
	dayKey := shared.DayCountDBKey(shortID)

	count, err := a.Config.DB.Exists(a.Ctx, totalCountKey).Result()
	if err != nil {
		return err
	}

	if count <= 0 {
		if err := a.Config.DB.Set(a.Ctx, totalCountKey, 0, 0).Err(); err != nil {
			return err
		}
	}

	count, err = a.Config.DB.Exists(a.Ctx, dayKey).Result()
	if err != nil {
		return err
	}

	if count <= 0 {
		if err := a.Config.DB.Set(a.Ctx, dayKey, 0, 24*time.Hour).Err(); err != nil {
			return err
		}
	}

	count, err = a.Config.DB.Exists(a.Ctx, weekKey).Result()
	if err != nil {
		return err
	}

	if count <= 0 {
		if err := a.Config.DB.Set(a.Ctx, weekKey, 0, 168*time.Hour).Err(); err != nil {
			return err
		}
	}

	return nil
}

func (a Actions) CreateURLMapping(longURL, employeeID string, expirationMins *int) (string, error) {
	shortID := generateShortURL(longURL, employeeID)

	for {
		if val, _ := a.GetLongURL(shortID); val == "" {
			break
		}

		shortID = generateShortURL(longURL, employeeID)
	}

	dur := 0 * time.Minute
	if expirationMins != nil {
		dur = time.Duration(*expirationMins) * time.Minute
	}

	key := shared.ShortenDBKey(shortID)
	val := shared.ShortenDBVal(employeeID, longURL)

	if err := a.Config.DB.Set(a.Ctx, key, val, dur).Err(); err != nil {
		return "", err
	}

	if err := a.setMetrics(shortID); err != nil {
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
	logger.Infof("trying to incr")

	totalKey := shared.TotalCountDBKey(shortID)
	weekKey := shared.WeekCountDBKey(shortID)
	dayKey := shared.DayCountDBKey(shortID)

	if err := a.Config.DB.IncrBy(a.Ctx, totalKey, 1).Err(); err != nil {
		return err
	}

	if err := a.Config.DB.IncrBy(a.Ctx, weekKey, 1).Err(); err != nil {
		return err
	}

	if err := a.Config.DB.IncrBy(a.Ctx, dayKey, 1).Err(); err != nil {
		return err
	}

	return nil
}

func (a Actions) GetShortURLMetrics(shortID, employeeID string) (*MetricsView, error) {
	key := shared.ShortenDBKey(shortID)
	totalKey := shared.TotalCountDBKey(shortID)
	weekKey := shared.WeekCountDBKey(shortID)
	dayKey := shared.DayCountDBKey(shortID)

	if !a.canAccess(key, employeeID) {
		return nil, errors.New("unauthorized")
	}

	totalRes, err := a.Config.DB.Get(a.Ctx, totalKey).Result()
	if err != nil {
		return nil, err
	}

	totalNum, err := strconv.Atoi(totalRes)
	if err != nil {
		return nil, err
	}

	dayRes, err := a.Config.DB.Get(a.Ctx, dayKey).Result()
	if err != nil {
		return nil, err
	}

	dayNum, err := strconv.Atoi(dayRes)
	if err != nil {
		return nil, err
	}

	weekRes, err := a.Config.DB.Get(a.Ctx, weekKey).Result()
	if err != nil {
		return nil, err
	}

	weekNum, err := strconv.Atoi(weekRes)
	if err != nil {
		return nil, err
	}

	return &MetricsView{
		TotalCount:    totalNum,
		PastDayCount:  dayNum,
		PastWeekCount: weekNum,
	}, nil
}

func (a Actions) GetLongURL(shortID string) (string, error) {
	key := shared.ShortenDBKey(shortID)

	result, err := a.Config.DB.Get(a.Ctx, key).Result()
	if err != nil {
		return "", err
	}

	return strings.Split(result, "::")[1], nil
}
