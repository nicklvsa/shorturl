package actions

import (
	"context"
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

	if err := a.Config.DB.Set(a.Ctx, key, val, dur); err != nil {
		return "", err.Err()
	}

	return shortID, nil
}

func (a Actions) GetLongURL(shortID string) (string, error) {
	key := shared.ShortenDBKey(shortID)

	result, err := a.Config.DB.Get(a.Ctx, key).Result()
	if err != nil {
		return "", err
	}

	return strings.Split(result, "::")[1], nil
}
