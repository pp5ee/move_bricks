package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// Config holds application configuration
type Config struct {
	HuobiAccessKey string
	HuobiSecretKey string
	GateAccessKey  string
	GateSecretKey  string
	SMTPHost       string
	SMTPPort       int
	SMTPUsername   string
	SMTPPassword   string
	EmailFrom      string
	EmailTo        string

	// Arbitrage config
	HuobiSymbol    string
	GateSymbol     string
	MinProfit      float64
	TradeAmount    float64
	CheckInterval  int // milliseconds
	WSEnabled      bool
	PriceEnabled   bool
	ArbitrageEnabled bool
}

func loadConfig() *Config {
	checkInterval := 500 // default 500ms
	if iv := os.Getenv("ARBITRAGE_CHECK_INTERVAL"); iv != "" {
		if v, err := strconv.Atoi(iv); err == nil {
			checkInterval = v
		}
	}

	minProfit := 0.001 // default 0.001 USDT
	if mp := os.Getenv("MIN_PROFIT"); mp != "" {
		if v, err := strconv.ParseFloat(mp, 64); err == nil {
			minProfit = v
		}
	}

	tradeAmount := 2.0 // default 2 NAS
	if ta := os.Getenv("TRADE_AMOUNT"); ta != "" {
		if v, err := strconv.ParseFloat(ta, 64); err == nil {
			tradeAmount = v
		}
	}

	return &Config{
		HuobiAccessKey:    os.Getenv("HUOBI_ACCESS_KEY"),
		HuobiSecretKey:    os.Getenv("HUOBI_SECRET_KEY"),
		GateAccessKey:     os.Getenv("GATE_ACCESS_KEY"),
		GateSecretKey:     os.Getenv("GATE_SECRET_KEY"),
		SMTPHost:          os.Getenv("SMTP_HOST"),
		SMTPPort:          587,
		SMTPUsername:      os.Getenv("SMTP_USERNAME"),
		SMTPPassword:      os.Getenv("SMTP_PASSWORD"),
		EmailFrom:         os.Getenv("EMAIL_FROM"),
		EmailTo:           os.Getenv("EMAIL_TO"),
		HuobiSymbol:       os.Getenv("HUOBI_SYMBOL"),
		GateSymbol:        os.Getenv("GATE_SYMBOL"),
		MinProfit:         minProfit,
		TradeAmount:       tradeAmount,
		CheckInterval:     checkInterval,
		WSEnabled:         os.Getenv("WS_ENABLED") == "true",
		PriceEnabled:      os.Getenv("PRICE_NOTICE_ENABLED") == "true",
		ArbitrageEnabled:  os.Getenv("ARBITRAGE_ENABLED") == "true",
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("[move-bricks] ")

	log.Println("=== Go Trading Bot Starting ===")

	config := loadConfig()

	// Create Huobi client
	huobi := NewHuobiClient()

	// Create Gate client
	gate := NewGateClient()

	// Create email sender
	emailConfig := DefaultEmailConfig()
	emailSender := NewEmailSender(emailConfig)

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	// Check and display balances
	log.Println("Checking account balances...")
	trader := &ArbitrageTrader{
		huobiClient: huobi,
		gateClient:  gate,
	}
	hbNAS, hbUSDT, gateNAS, gateUSDT, err := trader.CheckBalance()
	if err != nil {
		log.Printf("Failed to fetch balances: %v", err)
	} else {
		log.Printf("Huobi - NAS: %.4f, USDT: %.4f", hbNAS, hbUSDT)
		log.Printf("Gate   - NAS: %.4f, USDT: %.4f", gateNAS, gateUSDT)

		// Check if balances are sufficient
		if hbNAS < 1 {
			log.Println("WARNING: Huobi NAS balance insufficient!")
			emailSender.Send("火币NAS余额不足!", fmt.Sprintf("火币NAS余额: %.4f", hbNAS))
		}
		if gateNAS < 1 {
			log.Println("WARNING: Gate NAS balance insufficient!")
			emailSender.Send("Gate NAS余额不足!", fmt.Sprintf("Gate NAS余额: %.4f", gateNAS))
		}
	}

	// Start WebSocket listener if enabled
	if config.WSEnabled {
		wsClient := NewHuobiWSClient(func(symbol string, kline KLine) {
			log.Printf("K-line update: %s - O: %.8f H: %.8f L: %.8f C: %.8f",
				symbol, kline.Open, kline.High, kline.Low, kline.Close)
		})
		if err := wsClient.Start(); err != nil {
			log.Printf("WebSocket error: %v", err)
		} else {
			// Subscribe to symbols
			if config.HuobiSymbol != "" {
				if err := wsClient.Subscribe(config.HuobiSymbol, "1min"); err != nil {
					log.Printf("Subscribe error for %s: %v", config.HuobiSymbol, err)
				}
			}
			go wsClient.WaitForSignal()
		}
	}

	// Start price monitoring if enabled
	if config.PriceEnabled && config.HuobiSymbol != "" {
		noticeConfig := &PriceNoticeConfig{
			Symbols:       []string{config.HuobiSymbol},
			HighThreshold: 0,
			LowThreshold:  0,
			CheckInterval: config.CheckInterval,
			Enabled:       true,
		}
		priceChecker := NewPriceChecker(noticeConfig, emailSender, huobi)
		if err := priceChecker.Start(ctx); err != nil {
			log.Printf("Price monitor error: %v", err)
		}
	}

	// Start arbitrage trading if enabled
	if config.ArbitrageEnabled && config.HuobiSymbol != "" && config.GateSymbol != "" {
		arbitrageConfig := &ArbitrageConfig{
			HuobiSymbol:   config.HuobiSymbol,
			GateSymbol:    config.GateSymbol,
			MinProfit:     config.MinProfit,
			TradeAmount:   config.TradeAmount,
			CheckInterval: config.CheckInterval,
		}
		arbitrageTrader := NewArbitrageTrader(arbitrageConfig, huobi, gate, emailSender)
		if err := arbitrageTrader.Start(ctx); err != nil {
			log.Printf("Arbitrage error: %v", err)
		}
	}

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Exiting")
}