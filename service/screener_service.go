package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ketul1009/stockscreener-backend/db"
)

type ScreenerService struct {
	DB *db.Queries
}

type Screener struct {
	ID       int64
	Name     string
	Username string
	Rules    []map[string]interface{}
}

func (s *ScreenerService) CreateScreener(ctx context.Context, screener *Screener) (*Screener, error) {
	rulesJSON, err := json.Marshal(screener.Rules)
	if err != nil {
		return nil, err
	}

	dbScreener, err := s.DB.CreateScreener(ctx, db.CreateScreenerParams{
		Username: screener.Username,
		Name:     screener.Name,
		Rules:    rulesJSON,
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var rules []map[string]interface{}
	if err := json.Unmarshal(dbScreener.Rules, &rules); err != nil {
		return nil, err
	}

	return &Screener{
		ID:       int64(dbScreener.ID),
		Username: dbScreener.Username,
		Name:     dbScreener.Name,
		Rules:    rules,
	}, nil
}
