package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ketul1009/stockscreener-backend/db"
	"github.com/ketul1009/stockscreener-backend/pkg/logger"
	utils "github.com/ketul1009/stockscreener-backend/utils"
	"go.uber.org/zap"
)

type WatchlistService struct {
	DB   *db.Queries
	Pool *pgxpool.Pool
}

type Watchlist struct {
	ID        int64                    `json:"id"`
	Name      string                   `json:"name"`
	UserID    string                   `json:"user_id"`
	StockList []map[string]interface{} `json:"stock_list"`
}

func (s *WatchlistService) CreateWatchlist(ctx context.Context, watchlist *Watchlist) (*Watchlist, error) {

	uniqueStocks := utils.NewSet[string]()
	stockList := []map[string]interface{}{}
	for _, stock := range watchlist.StockList {
		if !uniqueStocks.Has(stock["symbol"].(string)) {
			stockList = append(stockList, stock)
			uniqueStocks.Add(stock["symbol"].(string))
		}
	}
	stockListJSON, err := json.Marshal(stockList)
	if err != nil {
		return nil, err
	}

	dbWatchlist, err := s.DB.CreateWatchlist(ctx, db.CreateWatchlistParams{
		Name:      watchlist.Name,
		UserID:    pgtype.UUID{Bytes: uuid.MustParse(watchlist.UserID), Valid: true},
		Stocks:    stockListJSON,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})

	if err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"unique_watchlist_name\" (SQLSTATE 23505)" {
			return nil, errors.New("watchlist name already exists")
		}
		logger.Error("Failed to create watchlist", zap.Error(err))
		return nil, err
	}

	return &Watchlist{
		ID:        int64(dbWatchlist.ID),
		Name:      dbWatchlist.Name,
		UserID:    dbWatchlist.UserID.String(),
		StockList: watchlist.StockList,
	}, nil
}

func (s *WatchlistService) GetWatchlist(ctx context.Context, id int32) (*Watchlist, error) {
	watchlist, err := s.DB.GetWatchlist(ctx, id)
	if err != nil {
		return nil, err
	}

	var stockList []map[string]interface{}
	err = json.Unmarshal(watchlist.Stocks, &stockList)
	if err != nil {
		return nil, err
	}

	return &Watchlist{
		ID:        int64(watchlist.ID),
		Name:      watchlist.Name,
		UserID:    watchlist.UserID.String(),
		StockList: stockList,
	}, nil
}

func (s *WatchlistService) GetAllWatchlists(ctx context.Context, userID string) ([]*Watchlist, error) {
	watchlists, err := s.DB.GetAllWatchlists(ctx, pgtype.UUID{Bytes: uuid.MustParse(userID), Valid: true})
	if err != nil {
		return nil, err
	}

	var watchlistList []*Watchlist
	for _, watchlist := range watchlists {
		var stockList []map[string]interface{}
		err = json.Unmarshal(watchlist.Stocks, &stockList)
		if err != nil {
			return nil, err
		}
		watchlistList = append(watchlistList, &Watchlist{
			ID:        int64(watchlist.ID),
			Name:      watchlist.Name,
			UserID:    watchlist.UserID.String(),
			StockList: stockList,
		})
	}

	return watchlistList, nil
}

func (s *WatchlistService) UpdateWatchlist(ctx context.Context, id int32, watchlist *Watchlist) (*Watchlist, error) {
	uniqueStocks := utils.NewSet[string]()
	stockList := []map[string]interface{}{}
	for _, stock := range watchlist.StockList {
		if !uniqueStocks.Has(stock["symbol"].(string)) {
			stockList = append(stockList, stock)
			uniqueStocks.Add(stock["symbol"].(string))
		}
	}
	stockListJSON, err := json.Marshal(stockList)
	if err != nil {
		return nil, err
	}

	updatedWatchlist, err := s.DB.UpdateWatchlist(ctx, db.UpdateWatchlistParams{
		ID:        id,
		UserID:    pgtype.UUID{Bytes: uuid.MustParse(watchlist.UserID), Valid: true},
		Name:      watchlist.Name,
		Stocks:    stockListJSON,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &Watchlist{
		ID:        int64(updatedWatchlist.ID),
		Name:      updatedWatchlist.Name,
		UserID:    updatedWatchlist.UserID.String(),
		StockList: watchlist.StockList,
	}, nil
}

func (s *WatchlistService) DeleteWatchlist(ctx context.Context, id int32) error {
	err := s.DB.DeleteWatchlist(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
