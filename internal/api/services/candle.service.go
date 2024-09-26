package services

import (
	"log"

	"github.com/hoangtm1601/go-binance-crawler/internal/api/repositories"
	"github.com/hoangtm1601/go-binance-crawler/internal/models"
	"github.com/hoangtm1601/go-binance-crawler/utils"
)

type CandleService struct {
	repo *repositories.CandleRepository
}

func NewCandleService(repo *repositories.CandleRepository) *CandleService {
	return &CandleService{repo: repo}
}

func (s *CandleService) CreateCandle(candle *models.Candle) error {
	return s.repo.Create(candle)
}

func (s *CandleService) CreateCandles(candles []*models.Candle) error {
	return s.repo.InsertMany(candles)
}

func (s *CandleService) GetCandleByID(id uint) (*models.Candle, error) {
	return s.repo.GetByID(id)
}

func (s *CandleService) GetCandlesBySymbol(symbol string) ([]models.Candle, error) {
	return s.repo.GetBySymbol(symbol)
}

func (s *CandleService) GetLatestCandleByInterval(symbol string, interval models.CandleInterval) (*models.Candle, error) {
	return s.repo.GetLatestCandleByInterval(symbol, interval)
}

func (s *CandleService) ProcessCandles(symbol string, candles []*models.Candle) {
	candleInsertToDatabase := make(map[models.CandleInterval][]*models.Candle)
	candleUpdateToDatabase := make(map[models.CandleInterval]map[int64]bool)

	for _, candle := range candles {
		for _, interval := range []models.CandleInterval{models.OneMin, models.FiveMin, models.FifteenMin, models.ThirtyMin, models.SixtyMin, models.TwoFortyMin, models.SevenTwentyMin, models.FourteenFortyMin} {
			if interval != models.OneMin && len(candleInsertToDatabase[interval]) == 0 {
				latestCandle, err := s.repo.GetLatestCandleByInterval(symbol, interval)

				if err == nil && latestCandle != nil && latestCandle.LastEnd < latestCandle.End {
					candleInsertToDatabase[interval] = append([]*models.Candle{latestCandle}, candleInsertToDatabase[interval]...)
					if candleUpdateToDatabase[interval] == nil {
						candleUpdateToDatabase[interval] = make(map[int64]bool)
					}
					candleUpdateToDatabase[interval][latestCandle.Start] = true
				}
			}

			startDate := candle.Start
			endDate := candle.Start + (int64(utils.GetMinute(interval))*60000 - 1)

			if interval == models.OneMin && endDate != candle.End {
				log.Printf("Warning with candle data, interval mismatch! \n %v", candle)
			}

			if len(candleInsertToDatabase[interval]) == 0 {
				newCandle := &models.Candle{
					Symbol:   candle.Symbol,
					Interval: interval,
					Start:    startDate,
					End:      endDate,
					LastEnd:  candle.End,
					Op:       candle.Op,
					Hi:       candle.Hi,
					Lo:       candle.Lo,
					Cl:       candle.Cl,
					Bv:       candle.Bv,
					Qv:       candle.Qv,
					Cnt:      candle.Cnt,
					Tbv:      candle.Tbv,
					Tqv:      candle.Tqv,
				}
				candleInsertToDatabase[interval] = append([]*models.Candle{newCandle}, candleInsertToDatabase[interval]...)
				// s.cacheService.LeftPushOne(interval, newCandle)
				// s.indicatorService.CalcCurrentIndicators(interval)
				continue
			}

			lastKnown := &candleInsertToDatabase[interval][0] // Use a pointer

			if (*lastKnown).Start >= (candle.Start) {
				continue
			}

			if (*lastKnown).LastEnd >= (candle.End) {
				continue
			}

			if (*lastKnown).End < (candle.Start) {
				// Handle missing previous candle
			}

			if (*lastKnown).End >= (candle.End) {
				if (*lastKnown).LastEnd < (candle.End) {
					(*lastKnown).LastEnd = candle.End
					(*lastKnown).Hi = max((*lastKnown).Hi, candle.Hi)
					(*lastKnown).Lo = min((*lastKnown).Lo, candle.Lo)
					(*lastKnown).Cl = candle.Cl
					(*lastKnown).Bv += candle.Bv
					(*lastKnown).Qv += candle.Qv
					(*lastKnown).Cnt += candle.Cnt
					(*lastKnown).Tbv += candle.Tbv
					(*lastKnown).Tqv += candle.Tqv
					// s.cacheService.Update(interval, lastKnown)
				}
			} else {
				if startDate != candle.Start {
					log.Printf("Warning, startDate mismatch! startDate: %v - data.start: %v %v", startDate, candle.Start, interval)
				}
				newCandle := &models.Candle{
					Symbol:   candle.Symbol,
					Interval: interval,
					Start:    startDate,
					End:      endDate,
					LastEnd:  candle.End,
					Op:       candle.Op,
					Hi:       candle.Hi,
					Lo:       candle.Lo,
					Cl:       candle.Cl,
					Bv:       candle.Bv,
					Qv:       candle.Qv,
					Cnt:      candle.Cnt,
					Tbv:      candle.Tbv,
					Tqv:      candle.Tqv,
				}
				// s.cacheService.LeftPushOne(interval, newCandle)
				candleInsertToDatabase[interval] = append([]*models.Candle{newCandle}, candleInsertToDatabase[interval]...)
			}
			// s.indicatorService.CalcCurrentIndicators(interval)
		}
	}

	// Prepare candlesToUpdate and candlesToInsert
	var candlesToInsert []*models.Candle
	var candlesToUpdate []*models.Candle

	for interval, candles := range candleInsertToDatabase {
		for _, candle := range candles {
			if _, exists := candleUpdateToDatabase[interval]; exists {
				candlesToUpdate = append(candlesToUpdate, candle)
				continue
			}
			candlesToInsert = append(candlesToInsert, candle)
		}
	}

	// Use BulkInsertAndUpdate from repositories
	err := s.repo.BulkInsertAndUpdate(candlesToInsert, candlesToUpdate)
	if err != nil {
		log.Fatalf("CandleService::ProcessCandles() Error in BulkInsertAndUpdate: %v", err)
	}
}
