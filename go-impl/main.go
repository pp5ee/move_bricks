package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// Config holds application configuration
type Config struct {
	HuobiAccessKey string
	HuobiSecretKey string
	SMTPHost       string
	SMTPPort       int
	SMTPUsername   string
	SMTPPassword   string
	EmailFrom      string
	EmailTo        string
	Symbols        []string
	PriceHigh      float64
	PriceLow       float64
	PriceInterval  int
	WSEnabled      bool
	PriceEnabled   bool
}

func loadConfig() *Config {
	symbolsStr := os.Getenv("SYMBOLS")
	symbols := []string{}
	if symbolsStr != "" {
		symbols = []string{symbolsStr} // For single symbol, split by comma if needed
	}

	priceInterval := 60
	if iv := os.Getenv("PRICE_INTERVAL"); iv != "" {
		fmt.Sscanf(iv, "%d", &priceInterval)
	}

	return &Config{
		HuobiAccessKey: os.Getenv("HUOBI_ACCESS_KEY"),
		HuobiSecretKey: os.Getenv("HUOBI_SECRET_KEY"),
		SMTPHost:       os.Getenv("SMTP_HOST"),
		SMTPPort:       587,
		SMTPUsername:   os.Getenv("SMTP_USERNAME"),
		SMTPPassword:   os.Getenv("SMTP_PASSWORD"),
		EmailFrom:      os.Getenv("EMAIL_FROM"),
		EmailTo:        os.Getenv("EMAIL_TO"),
		Symbols:        symbols,
		PriceHigh:      0,
		PriceLow:       0,
		PriceInterval:  priceInterval,
		WSEnabled:      os.Getenv("WS_ENABLED") == "true",
		PriceEnabled:   os.Getenv("PRICE_NOTICE_ENABLED") == "true",
	}
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("[move-bricks] ")

	config := loadConfig()

	// Create Huobi client
	huobi := NewHuobiClient()

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
			for _, symbol := range config.Symbols {
				if err := wsClient.Subscribe(symbol, "1min"); err != nil {
					log.Printf("Subscribe error for %s: %v", symbol, err)
				}
			}
			go wsClient.WaitForSignal()
		}
	}

	// Start price monitoring if enabled
	if config.PriceEnabled && len(config.Symbols) > 0 {
		noticeConfig := &PriceNoticeConfig{
			Symbols:       config.Symbols,
			HighThreshold: config.PriceHigh,
			LowThreshold:  config.PriceLow,
			CheckInterval: config.PriceInterval,
			Enabled:       true,
		}
		priceChecker := NewPriceChecker(noticeConfig, emailSender, huobi)
		if err := priceChecker.Start(ctx); err != nil {
			log.Printf("Price monitor error: %v", err)
		}
	}

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Exiting")
}