package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// PriceNoticeConfig holds price alert configuration
type PriceNoticeConfig struct {
	Symbols        []string
	HighThreshold  float64
	LowThreshold   float64
	CheckInterval  int // seconds
	Enabled        bool
}

// DefaultPriceNoticeConfig creates default config from environment
func DefaultPriceNoticeConfig() *PriceNoticeConfig {
	symbolsStr := os.Getenv("SYMBOLS")
	symbols := []string{}
	if symbolsStr != "" {
		symbols = strings.Split(symbolsStr, ",")
	}

	return &PriceNoticeConfig{
		Symbols:       symbols,
		HighThreshold: 0,
		LowThreshold:  0,
		CheckInterval: 60,
		Enabled:       os.Getenv("PRICE_NOTICE_ENABLED") == "true",
	}
}

// Validate validates the price notice config
func (c *PriceNoticeConfig) Validate() error {
	if len(c.Symbols) == 0 {
		return fmt.Errorf("at least one symbol is required")
	}
	return nil
}

// PriceChecker checks price thresholds
type PriceChecker struct {
	config       *PriceNoticeConfig
	emailSender  *EmailSender
	huobiClient  *HuobiClient
	lastNotified map[string]int64 // timestamp of last notification
}

// NewPriceChecker creates a new price checker
func NewPriceChecker(config *PriceNoticeConfig, emailSender *EmailSender, huobiClient *HuobiClient) *PriceChecker {
	return &PriceChecker{
		config:       config,
		emailSender:  emailSender,
		huobiClient:  huobiClient,
		lastNotified: make(map[string]int64),
	}
}

// CheckPrices checks all configured symbols for threshold breaches
func (p *PriceChecker) CheckPrices() error {
	for _, symbol := range p.config.Symbols {
		orderBook, err := p.huobiClient.FetchOrderBook(symbol, 1)
		if err != nil {
			return fmt.Errorf("failed to fetch order book for %s: %w", symbol, err)
		}

		if len(orderBook.Asks) > 0 {
			price := getFloat(orderBook.Asks[0])
			p.checkThreshold(symbol, price)
		}
	}
	return nil
}

// checkThreshold checks a single price against thresholds
func (p *PriceChecker) checkThreshold(symbol string, price float64) {
	highEnabled := p.config.HighThreshold > 0
	lowEnabled := p.config.LowThreshold > 0

	if highEnabled && price >= p.config.HighThreshold {
		p.sendAlert(symbol, price, p.config.HighThreshold, "HIGH")
	}

	if lowEnabled && price <= p.config.LowThreshold {
		p.sendAlert(symbol, price, p.config.LowThreshold, "LOW")
	}
}

// sendAlert sends a price alert
func (p *PriceChecker) sendAlert(symbol string, price, threshold float64, alertType string) {
	// Rate limit: only send one alert per minute per symbol
	lastNotified := p.lastNotified[symbol]
	now := nowTimestamp()
	if now-lastNotified < 60 {
		return
	}

	p.lastNotified[symbol] = now

	if p.emailSender != nil {
		p.emailSender.SendAlert(symbol, price, threshold, alertType)
	}
}

// Start starts the price checking loop
func (p *PriceChecker) Start(ctx context.Context) error {
	if err := p.config.Validate(); err != nil {
		return err
	}

	ticker := newTicker(p.config.CheckInterval)
	defer ticker.stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.c:
			if err := p.CheckPrices(); err != nil {
				fmt.Fprintf(os.Stderr, "Price check error: %v\n", err)
			}
		}
	}
}

// Ticker is a simple interval ticker
type Ticker struct {
	c    chan struct{}
	stop chan struct{}
}

func newTicker(interval int) *Ticker {
	t := &Ticker{
		c:    make(chan struct{}),
		stop: make(chan struct{}),
	}
	go func() {
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case t.c <- struct{}{}:
				default:
				}
			case <-t.stop:
				return
			}
		}
	}()
	return t
}

func (t *Ticker) stop() {
	close(t.stop)
}

func nowTimestamp() int64 {
	return time.Now().Unix()
}

func getFloat(v []interface{}) float64 {
	if len(v) == 0 {
		return 0
	}
	switch val := v[0].(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}
	return 0
}