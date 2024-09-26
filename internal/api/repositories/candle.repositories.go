package repositories

import (
	"errors"

	"github.com/hoangtm1601/go-binance-crawler/internal/models"
	"gorm.io/gorm"
)

type CandleRepository struct {
	db *gorm.DB
}

func NewCandleRepository(db *gorm.DB) *CandleRepository {
	return &CandleRepository{db: db}
}

func (r *CandleRepository) Create(candle *models.Candle) error {
	return r.db.Create(candle).Error
}

func (r *CandleRepository) InsertMany(candles []*models.Candle) error {
	return r.db.Create(candles).Error
}

func (r *CandleRepository) GetByID(id uint) (*models.Candle, error) {
	var candle models.Candle
	err := r.db.First(&candle, id).Error
	return &candle, err
}

func (r *CandleRepository) GetBySymbol(symbol string) ([]models.Candle, error) {
	var candles []models.Candle
	err := r.db.Where("symbol = ?", symbol).Find(&candles).Error
	return candles, err
}

func (r *CandleRepository) GetLatestCandleByInterval(symbol string, interval models.CandleInterval) (*models.Candle, error) {
	var candle models.Candle
	err := r.db.Where("symbol = ? and interval = ?", symbol, interval).Order("start DESC").First(&candle).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &candle, err
}

func (r *CandleRepository) BulkInsertAndUpdate(candlesToInsert []*models.Candle, candlesToUpdate []*models.Candle) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Bulk insert new candles
		if len(candlesToInsert) > 0 {
			for i := 0; i < len(candlesToInsert); i += 1000 {
				end := i + 1000
				if end > len(candlesToInsert) {
					end = len(candlesToInsert)
				}
				if err := tx.Create(candlesToInsert[i:end]).Error; err != nil {
					return err
				}
			}
		}

		// Update existing candles
		for _, candle := range candlesToUpdate {
			if err := tx.Model(&models.Candle{}).
				Where("symbol = ? AND interval = ? AND start = ?", candle.Symbol, candle.Interval, candle.Start).
				Updates(map[string]interface{}{
					"last_end": candle.LastEnd,
					"hi":       candle.Hi,
					"lo":       candle.Lo,
					"cl":       candle.Cl,
					"bv":       candle.Bv,
					"qv":       candle.Qv,
					"cnt":      candle.Cnt,
					"tbv":      candle.Tbv,
					"tqv":      candle.Tqv,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
