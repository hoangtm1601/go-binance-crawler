package services

import (
	"context"
	"log"
	"time"

	"github.com/adshao/go-binance/v2"
	"github.com/hoangtm1601/go-binance-crawler/internal/initializers"
	"github.com/hoangtm1601/go-binance-crawler/internal/models"
	"github.com/hoangtm1601/go-binance-crawler/utils"
)

type CrawlersService struct {
	binanceClient *binance.Client
	logger        *log.Logger
	candleService *CandleService
	config        initializers.Config
}

func NewCrawlersService(binanceClient *binance.Client, logger *log.Logger, candleService *CandleService, config *initializers.Config) *CrawlersService {
	return &CrawlersService{
		binanceClient: binanceClient,
		logger:        logger,
		candleService: candleService,
		config:        *config,
	}
}

func (s *CrawlersService) Crawl() {
	startTime := s.getCheckPoint()
	for {
		readableStartTime := time.Unix(0, startTime*int64(time.Millisecond)).UTC().Format(time.RFC3339)
		startTimer := time.Now()
		s.logger.Printf("CrawlService::crawl() | Start crawling for date time: %s", readableStartTime)

		queryEndTime := time.Now().Add(-1 * time.Minute).Truncate(time.Minute).Add(5 * time.Second)

		candles, err := s.binanceClient.NewKlinesService().Symbol(s.config.CrawlerSymbol).Limit(s.config.BinanceCandleLimit).Interval("1m").StartTime(startTime).EndTime(queryEndTime.UnixNano() / int64(time.Millisecond)).Do(context.Background())
		if err != nil {
			s.logger.Printf("CrawlService::crawl() | BinanceClient error: %s", err.Error())
			time.Sleep(time.Duration(s.config.CrawlerSleepTimeOnceFailed) * time.Millisecond)
			continue
		}

		if len(candles) == 0 {
			s.delayOnCatchUp()
			continue
		}

		// Convert binance.Kline to models.Candle
		var candleModels []*models.Candle
		for _, kline := range candles {
			candleModels = append(candleModels, &models.Candle{
				Symbol:   s.config.CrawlerSymbol,
				Interval: s.config.CrawlerDefaultInterval,
				Start:    kline.OpenTime,
				End:      kline.CloseTime,
				Op:       utils.StringToFloat64(kline.Open),
				Hi:       utils.StringToFloat64(kline.High),
				Lo:       utils.StringToFloat64(kline.Low),
				Cl:       utils.StringToFloat64(kline.Close),
				Bv:       utils.StringToFloat64(kline.Volume),
				Qv:       utils.StringToFloat64(kline.QuoteAssetVolume),
				Cnt:      kline.TradeNum,
				Tbv:      utils.StringToFloat64(kline.TakerBuyBaseAssetVolume),
				Tqv:      utils.StringToFloat64(kline.TakerBuyQuoteAssetVolume),
			})
		}

		s.candleService.ProcessCandles(s.config.CrawlerSymbol, candleModels)

		lastCandleEndTime := candles[len(candles)-1].CloseTime
		startTime = lastCandleEndTime + 1
		endTimer := time.Since(startTimer)
		s.logger.Printf("CrawlService::crawl() | Done crawling for date time: %s | Time consumed: %.2f sec", time.Unix(0, lastCandleEndTime*int64(time.Millisecond)).UTC().Format(time.RFC3339), endTimer.Seconds())

		if time.Now().UnixNano()/int64(time.Millisecond) > lastCandleEndTime+60000 {
			time.Sleep(time.Duration(s.config.CrawlerSleepTimeBetweenPastCrawling) * time.Millisecond)
			continue
		}

		s.delayOnCatchUp(startTime)
	}
}

func (s *CrawlersService) delayOnCatchUp(nextCandle ...int64) {
	sleepTimeToNextWindow := s.getNextTimeWindowDelayTime(nextCandle...)
	s.logger.Printf("CrawlService::crawl() | Catch up now ~ lets take a rest in %.2f s", float64(sleepTimeToNextWindow)/1000)
	time.Sleep(time.Duration(sleepTimeToNextWindow) * time.Millisecond)
}

func (s *CrawlersService) getNextTimeWindowDelayTime(nextCandle ...int64) int64 {
	now := time.Now()
	nextWindow := now
	if len(nextCandle) > 0 {
		nextWindow = time.Unix(0, nextCandle[0]*int64(time.Millisecond))
	}
	nextWindow = nextWindow.Add(time.Minute).Truncate(time.Minute).Add(5 * time.Second)
	return nextWindow.UnixNano()/int64(time.Millisecond) - now.UnixNano()/int64(time.Millisecond)
}

func (s *CrawlersService) getCheckPoint() int64 {
	latestCandle, err := s.candleService.GetLatestCandleByInterval(s.config.CrawlerSymbol, s.config.CrawlerDefaultInterval)
	if latestCandle != nil {
		return latestCandle.End
	}
	if err != nil {
		return int64(s.config.CrawlerInitialCrawlingTime)
	}
	return int64(s.config.CrawlerInitialCrawlingTime)
}
