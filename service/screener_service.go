package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ketul1009/stockscreener-backend/db"
	engine "github.com/ketul1009/stockscreener-backend/stock-engine"
)

type ScreenerService struct {
	DB   *db.Queries
	Pool *pgxpool.Pool
}

type Screener struct {
	ID            int64                    `json:"id"`
	Name          string                   `json:"name"`
	StockUniverse string                   `json:"stock_universe"`
	UserID        string                   `json:"user_id"`
	Rules         []map[string]interface{} `json:"rules"`
}

func (s *ScreenerService) CreateScreener(ctx context.Context, screener *Screener) (*Screener, error) {
	rulesJSON, err := json.Marshal(screener.Rules)
	if err != nil {
		return nil, err
	}

	dbScreener, err := s.DB.CreateScreener(ctx, db.CreateScreenerParams{
		UserID: pgtype.UUID{Bytes: uuid.MustParse(screener.UserID), Valid: true},
		Name:   screener.Name,
		Rules:  rulesJSON,
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
		ID:     int64(dbScreener.ID),
		UserID: dbScreener.UserID.String(),
		Name:   dbScreener.Name,
		Rules:  rules,
	}, nil
}

func (s *ScreenerService) GetScreeners(ctx context.Context, username string) ([]Screener, error) {
	screeners, err := s.DB.GetScreeners(ctx, pgtype.UUID{Bytes: uuid.MustParse(username), Valid: true})
	if err != nil {
		return nil, err
	}

	var screenerList []Screener
	for _, screener := range screeners {
		var rules []map[string]interface{}
		if err := json.Unmarshal(screener.Rules, &rules); err != nil {
			return nil, err
		}

		screenerList = append(screenerList, Screener{
			ID:     int64(screener.ID),
			UserID: screener.UserID.String(),
			Name:   screener.Name,
			Rules:  rules,
		})
	}

	return screenerList, nil
}

func (s *ScreenerService) GetScreener(ctx context.Context, id int64) (*Screener, error) {
	screener, err := s.DB.GetScreener(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	var rules []map[string]interface{}
	if err := json.Unmarshal(screener.Rules, &rules); err != nil {
		return nil, err
	}

	return &Screener{
		ID:            int64(screener.ID),
		UserID:        screener.UserID.String(),
		StockUniverse: screener.StockUniverse,
		Name:          screener.Name,
		Rules:         rules,
	}, nil
}

func (s *ScreenerService) UpdateScreener(ctx context.Context, id int64, screener *Screener) (*Screener, error) {
	rulesJSON, err := json.Marshal(screener.Rules)
	if err != nil {
		return nil, err
	}

	dbScreener, err := s.DB.UpdateScreener(ctx, db.UpdateScreenerParams{
		ID:            int32(id),
		Name:          screener.Name,
		Rules:         rulesJSON,
		StockUniverse: screener.StockUniverse,
	})
	if err != nil {
		return nil, err
	}

	var rules []map[string]interface{}
	if err := json.Unmarshal(dbScreener.Rules, &rules); err != nil {
		return nil, err
	}

	return &Screener{
		ID:            int64(dbScreener.ID),
		UserID:        dbScreener.UserID.String(),
		StockUniverse: dbScreener.StockUniverse,
		Name:          dbScreener.Name,
		Rules:         rules,
	}, nil
}

func (s *ScreenerService) DeleteScreener(ctx context.Context, id int64) error {

	err := s.DB.DeleteScreener(ctx, int32(id))
	if err != nil {
		return err
	}

	return nil
}

func (s *ScreenerService) GetStockUniverse(ctx context.Context) ([]engine.Stock, error) {
	query := `SELECT symbol, indicators.indicators 
	FROM stocks s
	JOIN (
		SELECT indicators, stock_id
		FROM indicators_data 
		WHERE (
			DATE = '2025-05-14'
		)
	) AS indicators ON s.id = indicators.stock_id`
	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var stocks []engine.Stock
	for rows.Next() {
		var symbol string
		var indicatorsJSON []byte
		if err := rows.Scan(&symbol, &indicatorsJSON); err != nil {
			return nil, err
		}
		var indicators map[string]interface{}
		if err := json.Unmarshal(indicatorsJSON, &indicators); err != nil {
			return nil, err
		}
		stocks = append(stocks, engine.Stock{Symbol: symbol, Indicators: indicators})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()
	return stocks, nil
}
