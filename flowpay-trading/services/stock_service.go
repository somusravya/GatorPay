package services

import (
	"math"
	"math/rand"
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

// StockInfo holds detailed stock metadata
type StockInfo struct {
	Symbol      string
	Name        string
	Sector      string
	MarketCap   string
	PERatio     float64
	DivYield    float64
	High52W     float64
	Low52W      float64
	BasePrice   float64
}

var stockDatabase = []StockInfo{
	// Technology
	{Symbol: "AAPL", Name: "Apple Inc.", Sector: "Technology", MarketCap: "2.89T", PERatio: 29.5, DivYield: 0.55, High52W: 199.62, Low52W: 143.90, BasePrice: 185.0},
	{Symbol: "MSFT", Name: "Microsoft Corp.", Sector: "Technology", MarketCap: "2.78T", PERatio: 35.2, DivYield: 0.74, High52W: 384.30, Low52W: 275.37, BasePrice: 370.0},
	{Symbol: "GOOGL", Name: "Alphabet Inc.", Sector: "Technology", MarketCap: "1.72T", PERatio: 24.8, DivYield: 0.0, High52W: 153.78, Low52W: 102.21, BasePrice: 141.0},
	{Symbol: "AMZN", Name: "Amazon.com Inc.", Sector: "Technology", MarketCap: "1.55T", PERatio: 60.3, DivYield: 0.0, High52W: 189.77, Low52W: 118.35, BasePrice: 178.0},
	{Symbol: "NVDA", Name: "NVIDIA Corp.", Sector: "Technology", MarketCap: "1.12T", PERatio: 62.1, DivYield: 0.04, High52W: 502.66, Low52W: 222.97, BasePrice: 455.0},
	{Symbol: "META", Name: "Meta Platforms Inc.", Sector: "Technology", MarketCap: "820B", PERatio: 27.5, DivYield: 0.0, High52W: 340.70, Low52W: 167.05, BasePrice: 310.0},
	// Finance
	{Symbol: "JPM", Name: "JPMorgan Chase & Co.", Sector: "Finance", MarketCap: "430B", PERatio: 10.8, DivYield: 2.65, High52W: 172.96, Low52W: 123.11, BasePrice: 155.0},
	{Symbol: "V", Name: "Visa Inc.", Sector: "Finance", MarketCap: "510B", PERatio: 30.2, DivYield: 0.76, High52W: 270.42, Low52W: 220.68, BasePrice: 258.0},
	{Symbol: "GS", Name: "Goldman Sachs Group", Sector: "Finance", MarketCap: "120B", PERatio: 14.2, DivYield: 2.8, High52W: 389.58, Low52W: 292.50, BasePrice: 355.0},
	// Healthcare
	{Symbol: "JNJ", Name: "Johnson & Johnson", Sector: "Healthcare", MarketCap: "390B", PERatio: 22.5, DivYield: 2.95, High52W: 175.97, Low52W: 143.13, BasePrice: 162.0},
	{Symbol: "UNH", Name: "UnitedHealth Group", Sector: "Healthcare", MarketCap: "480B", PERatio: 23.8, DivYield: 1.35, High52W: 558.10, Low52W: 436.38, BasePrice: 520.0},
	{Symbol: "PFE", Name: "Pfizer Inc.", Sector: "Healthcare", MarketCap: "160B", PERatio: 18.6, DivYield: 5.7, High52W: 37.80, Low52W: 25.20, BasePrice: 28.5},
	// Crypto / Digital
	{Symbol: "BTC", Name: "Bitcoin (Simulated)", Sector: "Crypto", MarketCap: "850B", PERatio: 0.0, DivYield: 0.0, High52W: 73750.0, Low52W: 24930.0, BasePrice: 62500.0},
	{Symbol: "ETH", Name: "Ethereum (Simulated)", Sector: "Crypto", MarketCap: "290B", PERatio: 0.0, DivYield: 0.0, High52W: 4090.0, Low52W: 1520.0, BasePrice: 3250.0},
	{Symbol: "SOL", Name: "Solana (Simulated)", Sector: "Crypto", MarketCap: "45B", PERatio: 0.0, DivYield: 0.0, High52W: 126.0, Low52W: 18.0, BasePrice: 105.0},
	// Consumer / Other
	{Symbol: "TSLA", Name: "Tesla Inc.", Sector: "Consumer", MarketCap: "680B", PERatio: 58.5, DivYield: 0.0, High52W: 299.29, Low52W: 152.37, BasePrice: 210.5},
	{Symbol: "DIS", Name: "Walt Disney Co.", Sector: "Consumer", MarketCap: "180B", PERatio: 65.2, DivYield: 0.0, High52W: 123.74, Low52W: 78.73, BasePrice: 95.0},
	{Symbol: "NKE", Name: "Nike Inc.", Sector: "Consumer", MarketCap: "155B", PERatio: 28.1, DivYield: 1.4, High52W: 131.31, Low52W: 88.66, BasePrice: 102.0},
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

func findStock(symbol string) *StockInfo {
	symbol = strings.ToUpper(symbol)
	for _, s := range stockDatabase {
		if s.Symbol == symbol {
			return &s
		}
	}
	return nil
}

func (s *StockService) GetQuote(symbol string) (map[string]interface{}, error) {
	symbol = strings.ToUpper(symbol)
	if data, ok := s.getFromCache("quote_" + symbol); ok {
		return data.(map[string]interface{}), nil
	}

	stock := findStock(symbol)
	basePrice := 150.25
	if stock != nil {
		basePrice = stock.BasePrice
	}

	// Simulate small random movement from base
	jitter := basePrice * 0.005 * (rand.Float64() - 0.5)
	price := basePrice + jitter
	changePercent := (jitter / basePrice) * 100

	mock := map[string]interface{}{
		"symbol":         symbol,
		"price":          math.Round(price*100) / 100,
		"change":         math.Round(jitter*100) / 100,
		"change_percent": math.Round(changePercent*100) / 100,
		"volume":         12500000 + rand.Intn(5000000),
		"timestamp":      time.Now(),
	}

	s.setToCache("quote_"+symbol, mock, 30*time.Second)
	return mock, nil
}

func (s *StockService) GetDetails(symbol string) (map[string]interface{}, error) {
	symbol = strings.ToUpper(symbol)
	if data, ok := s.getFromCache("details_" + symbol); ok {
		return data.(map[string]interface{}), nil
	}

	stock := findStock(symbol)
	mock := map[string]interface{}{
		"symbol":         symbol,
		"name":           "Mocked Company Inc.",
		"sector":         "Technology",
		"market_cap":     "N/A",
		"pe_ratio":       0.0,
		"dividend_yield": 0.0,
		"52_week_high":   200.0,
		"52_week_low":    100.0,
	}

	if stock != nil {
		mock["name"] = stock.Name
		mock["sector"] = stock.Sector
		mock["market_cap"] = stock.MarketCap
		mock["pe_ratio"] = stock.PERatio
		mock["dividend_yield"] = stock.DivYield
		mock["52_week_high"] = stock.High52W
		mock["52_week_low"] = stock.Low52W
	}

	s.setToCache("details_"+symbol, mock, 60*time.Minute)
	return mock, nil
}

func (s *StockService) GetChart(symbol string) ([]map[string]interface{}, error) {
	symbol = strings.ToUpper(symbol)
	if data, ok := s.getFromCache("chart_" + symbol); ok {
		return data.([]map[string]interface{}), nil
	}

	stock := findStock(symbol)
	basePrice := 140.0
	if stock != nil {
		basePrice = stock.BasePrice * 0.92 // Start chart ~8% lower for visual trend
	}

	var chart []map[string]interface{}
	now := time.Now()
	price := basePrice
	for i := 30; i >= 0; i-- {
		// Simulate realistic random walk
		move := price * 0.015 * (rand.Float64() - 0.45) // slight upward bias
		price += move
		if price < basePrice*0.85 {
			price = basePrice * 0.85
		}

		timePoint := now.AddDate(0, 0, -i).Format("2006-01-02")
		chart = append(chart, map[string]interface{}{
			"time":  timePoint,
			"value": math.Round(price*100) / 100,
		})
	}

	s.setToCache("chart_"+symbol, chart, 5*time.Minute)
	return chart, nil
}

func (s *StockService) Search(query string) ([]map[string]string, error) {
	query = strings.ToUpper(query)
	var results []map[string]string

	for _, stock := range stockDatabase {
		if strings.Contains(strings.ToUpper(stock.Symbol), query) ||
			strings.Contains(strings.ToUpper(stock.Name), query) ||
			strings.Contains(strings.ToUpper(stock.Sector), query) {
			results = append(results, map[string]string{
				"symbol": stock.Symbol,
				"name":   stock.Name,
				"sector": stock.Sector,
			})
		}
	}

	if len(results) == 0 {
		// Return all stocks if no match
		for _, stock := range stockDatabase {
			results = append(results, map[string]string{
				"symbol": stock.Symbol,
				"name":   stock.Name,
				"sector": stock.Sector,
			})
		}
	}

	return results, nil
}

func (s *StockService) MarketSummary() ([]map[string]interface{}, error) {
	if data, ok := s.getFromCache("market_summary"); ok {
		return data.([]map[string]interface{}), nil
	}

	mock := []map[string]interface{}{
		{"index": "S&P 500", "price": 4500.25 + rand.Float64()*20, "change_percent": 0.5 + rand.Float64()*0.3},
		{"index": "NASDAQ", "price": 14000.5 + rand.Float64()*50, "change_percent": 1.2 + rand.Float64()*0.5},
		{"index": "DOW", "price": 35000.0 + rand.Float64()*100, "change_percent": -0.2 + rand.Float64()*0.4},
		{"index": "BTC/USD", "price": 62500.0 + rand.Float64()*200, "change_percent": 2.1 + rand.Float64()*1.0},
	}
	s.setToCache("market_summary", mock, 30*time.Second)
	return mock, nil
}
