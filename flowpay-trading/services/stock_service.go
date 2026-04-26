package services

import (
	"strings"
	"sync"
	"time"
)

type StockService struct {
	cache      map[string]CacheItem
	cacheMutex sync.RWMutex
}

type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

func NewStockService() *StockService {
	return &StockService{
		cache: make(map[string]CacheItem),
	}
}

func (s *StockService) getFromCache(key string) (interface{}, bool) {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	item, exists := s.cache[key]
	if exists && time.Now().Before(item.ExpiresAt) {
		return item.Data, true
	}
	return nil, false
}

func (s *StockService) setToCache(key string, data interface{}, ttl time.Duration) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	s.cache[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Mocked stock data generators
func (s *StockService) GetQuote(symbol string) (map[string]interface{}, error) {
	symbol = strings.ToUpper(symbol)
	if data, ok := s.getFromCache("quote_" + symbol); ok {
		return data.(map[string]interface{}), nil
	}

	// Mock response mimicking Yahoo Finance or Twelve Data
	mock := map[string]interface{}{
		"symbol": symbol,
		"price":  150.25,
		"change": 1.25,
		"change_percent": 0.84,
		"volume": 12500000,
		"timestamp": time.Now(),
	}
	
	if symbol == "AAPL" { mock["price"] = 185.0 }
	if symbol == "TSLA" { mock["price"] = 210.5 }
	
	s.setToCache("quote_"+symbol, mock, 1*time.Minute)
	return mock, nil
}

func (s *StockService) GetDetails(symbol string) (map[string]interface{}, error) {
	symbol = strings.ToUpper(symbol)
	if data, ok := s.getFromCache("details_" + symbol); ok {
		return data.(map[string]interface{}), nil
	}

	mock := map[string]interface{}{
		"symbol": symbol,
		"name": "Mocked Company Inc.",
		"sector": "Technology",
		"market_cap": "1.5T",
		"pe_ratio": 25.5,
		"dividend_yield": 1.5,
		"52_week_high": 200.0,
		"52_week_low": 100.0,
	}

	s.setToCache("details_"+symbol, mock, 60*time.Minute)
	return mock, nil
}

func (s *StockService) GetChart(symbol string) ([]map[string]interface{}, error) {
	symbol = strings.ToUpper(symbol)
	if data, ok := s.getFromCache("chart_" + symbol); ok {
		return data.([]map[string]interface{}), nil
	}

	// Mocking time series chart data for TradingView lightweight-charts
	var chart []map[string]interface{}
	basePrice := 140.0
	now := time.Now()
	for i := 30; i >= 0; i-- {
		timePoint := now.AddDate(0, 0, -i).Format("2006-01-02")
		chart = append(chart, map[string]interface{}{
			"time": timePoint,
			"value": basePrice + float64(i)*0.5,
		})
	}

	s.setToCache("chart_"+symbol, chart, 5*time.Minute)
	return chart, nil
}

func (s *StockService) Search(query string) ([]map[string]string, error) {
	mock := []map[string]string{
		{"symbol": "AAPL", "name": "Apple Inc."},
		{"symbol": "TSLA", "name": "Tesla Inc."},
		{"symbol": "GOOGL", "name": "Alphabet Inc."},
		{"symbol": "AMZN", "name": "Amazon.com Inc."},
		{"symbol": "MSFT", "name": "Microsoft Corp."},
	}
	return mock, nil
}

func (s *StockService) MarketSummary() ([]map[string]interface{}, error) {
	if data, ok := s.getFromCache("market_summary"); ok {
		return data.([]map[string]interface{}), nil
	}

	mock := []map[string]interface{}{
		{"index": "S&P 500", "price": 4500.25, "change_percent": 0.5},
		{"index": "NASDAQ", "price": 14000.5, "change_percent": 1.2},
		{"index": "DOW", "price": 35000.0, "change_percent": -0.2},
	}
	s.setToCache("market_summary", mock, 5*time.Minute)
	return mock, nil
}
