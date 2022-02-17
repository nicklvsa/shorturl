package actions

import (
	"context"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/nicklvsa/shorturl/shared"
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

	storedVal, err := a.Config.DB.Get(a.Ctx, key).Result()
	if err != nil {
		return err
	}

	val := shared.ShortenDBVal(employeeID, strings.Split(storedVal, "::")[1])

	if storedVal != val {
		return errors.New("unauthorized")
	}

	if err := a.Config.DB.Del(a.Ctx, key).Err(); err != nil {
		return err
	}

	return nil
}

func (a Actions) IncrShortURLCount(shortID string) error {
	totalKey := shared.TotalCountDBKey(shortID)
	weekKey := shared.WeekCountDBKey(shortID)
	dayKey := shared.DayCountDBKey(shortID)

	if err := a.Config.DB.Incr(a.Ctx, totalKey).Err(); err != nil {
		return err
	}

	if err := a.Config.DB.Incr(a.Ctx, weekKey).Err(); err != nil {
		return err
	}

	if err := a.Config.DB.Incr(a.Ctx, dayKey).Err(); err != nil {
		return err
	}

	return nil
}

func (a Actions) GetLongURL(shortID string) (string, error) {
	key := shared.ShortenDBKey(shortID)

	result, err := a.Config.DB.Get(a.Ctx, key).Result()
	if err != nil {
		return "", err
	}

	return strings.Split(result, "::")[1], nil
}
