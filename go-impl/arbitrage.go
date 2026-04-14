package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

// ArbitrageConfig holds arbitrage trading configuration
type ArbitrageConfig struct {
	HuobiSymbol string // e.g., "nasusdt"
	GateSymbol  string // e.g., "NAS_USDT"
	MinProfit   float64
	TradeAmount float64
	CheckInterval int // milliseconds
}

// ArbitrageTrader handles the arbitrage trading logic
type ArbitrageTrader struct {
	config       *ArbitrageConfig
	huobiClient  *HuobiClient
	gateClient   *GateClient
	emailSender  *EmailSender
	counter      int
	times        int
	mu           sync.Mutex
	ctx          context.Context
}

// NewArbitrageTrader creates a new arbitrage trader
func NewArbitrageTrader(config *ArbitrageConfig, huobi *HuobiClient, gate *GateClient, email *EmailSender) *ArbitrageTrader {
	return &ArbitrageTrader{
		config:      config,
		huobiClient: huobi,
		gateClient:  gate,
		emailSender: email,
		counter:     0,
		times:       0,
	}
}

// Start starts the arbitrage trading loop
func (t *ArbitrageTrader) Start(ctx context.Context) error {
	log.Printf("[Arbitrage] Starting trader for %s/%s", t.config.HuobiSymbol, t.config.GateSymbol)

	ticker := time.NewTicker(time.Duration(t.config.CheckInterval) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			t.executeArbitrage()
		}
	}
}

func (t *ArbitrageTrader) executeArbitrage() {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.times++

	// Fetch order books from both exchanges
	hbOrderBook, err := t.huobiClient.FetchOrderBook(t.config.HuobiSymbol, 5)
	if err != nil {
		log.Printf("[Arbitrage] Failed to fetch Huobi orderbook: %v", err)
		return
	}

	gateOrderBook, err := t.gateClient.FetchOrderBook(t.config.GateSymbol, 5)
	if err != nil {
		log.Printf("[Arbitrage] Failed to fetch Gate orderbook: %v", err)
		return
	}

	// Get best prices
	// Huobi: bid=买价(买一), ask=卖价(卖一)
	// Gate: bid=买价, ask=卖价
	var hbBidPrice, hbAskPrice, gateBidPrice, gateAskPrice float64

	if len(hbOrderBook.Bids) > 0 {
		hbBidPrice, _ = strconv.ParseFloat(fmt.Sprintf("%v", hbOrderBook.Bids[0][0]), 64)
	}
	if len(hbOrderBook.Asks) > 0 {
		hbAskPrice, _ = strconv.ParseFloat(fmt.Sprintf("%v", hbOrderBook.Asks[0][0]), 64)
	}
	if len(gateOrderBook.Bids) > 0 {
		gateBidPrice, _ = strconv.ParseFloat(fmt.Sprintf("%v", gateOrderBook.Bids[0][0]), 64)
	}
	if len(gateOrderBook.Asks) > 0 {
		gateAskPrice, _ = strconv.ParseFloat(fmt.Sprintf("%v", gateOrderBook.Asks[0][0]), 64)
	}

	// Get amounts
	hbBidAmount := 0.0
	hbAskAmount := 0.0
	gateBidAmount := 0.0
	gateAskAmount := 0.0

	if len(hbOrderBook.Bids) > 0 {
		hbBidAmount, _ = strconv.ParseFloat(fmt.Sprintf("%v", hbOrderBook.Bids[0][1]), 64)
	}
	if len(hbOrderBook.Asks) > 0 {
		hbAskAmount, _ = strconv.ParseFloat(fmt.Sprintf("%v", hbOrderBook.Asks[0][1]), 64)
	}
	if len(gateOrderBook.Bids) > 0 {
		gateBidAmount, _ = strconv.ParseFloat(fmt.Sprintf("%v", gateOrderBook.Bids[0][1]), 64)
	}
	if len(gateOrderBook.Asks) > 0 {
		gateAskAmount, _ = strconv.ParseFloat(fmt.Sprintf("%v", gateOrderBook.Asks[0][1]), 64)
	}

	log.Printf("[Arbitrage] Huobi: bid=%.8f(%.4f) ask=%.8f(%.4f), Gate: bid=%.8f(%.4f) ask=%.8f(%.4f)",
		hbBidPrice, hbBidAmount, hbAskPrice, hbAskAmount,
		gateBidPrice, gateBidAmount, gateAskPrice, gateAskAmount)

	// Strategy 1: Buy on Gate, Sell on Huobi
	// Gate bid price < Huobi ask price -> buy on Gate, sell on Huobi
	profit1 := hbAskPrice - gateBidPrice
	if profit1 > t.config.MinProfit && gateBidAmount > t.config.TradeAmount && hbAskAmount > t.config.TradeAmount {
		log.Printf("[Arbitrage] Opportunity 1: Buy on Gate @ %.8f, Sell on Huobi @ %.8f, Profit: %.8f",
			gateBidPrice, hbAskPrice, profit1)
		t.executeTrade(gateBidPrice, hbAskPrice, "gate_to_huobi", profit1)
		return
	}

	// Strategy 2: Buy on Huobi, Sell on Gate
	// Huobi bid price < Gate ask price -> buy on Huobi, sell on Gate
	profit2 := gateAskPrice - hbBidPrice
	if profit2 > t.config.MinProfit && hbBidAmount > t.config.TradeAmount && gateAskAmount > t.config.TradeAmount {
		log.Printf("[Arbitrage] Opportunity 2: Buy on Huobi @ %.8f, Sell on Gate @ %.8f, Profit: %.8f",
			hbBidPrice, gateAskPrice, profit2)
		t.executeTrade(hbBidPrice, gateAskPrice, "huobi_to_gate", profit2)
		return
	}
}

func (t *ArbitrageTrader) executeTrade(buyPrice, sellPrice float64, direction string, profit float64) {
	amount := t.config.TradeAmount
	baseSymbol := t.config.HuobiSymbol[:len(t.config.HuobiSymbol)-4] // remove "usdt"
	quoteSymbol := "usdt"

	var buyErr, sellErr error

	if direction == "gate_to_huobi" {
		// Buy on Gate, Sell on Huobi
		_, buyErr = t.gateClient.CreateLimitBuyOrder(t.config.GateSymbol, amount, buyPrice)
		if buyErr == nil {
			_, sellErr = t.huobiClient.CreateLimitSellOrder(t.config.HuobiSymbol, amount, sellPrice)
		}
	} else {
		// Buy on Huobi, Sell on Gate
		_, buyErr = t.huobiClient.CreateLimitBuyOrder(t.config.HuobiSymbol, amount, buyPrice)
		if buyErr == nil {
			_, sellErr = t.gateClient.CreateLimitSellOrder(t.config.GateSymbol, amount, sellPrice)
		}
	}

	t.counter++

	if buyErr != nil || sellErr != nil {
		log.Printf("[Arbitrage] Trade failed: buy_err=%v, sell_err=%v", buyErr, sellErr)
		subject := fmt.Sprintf("[搬砖] 交易失败! 方向: %s", direction)
		content := fmt.Sprintf("买入价: %.8f, 卖出价: %.8f, 利润: %.8f\n错误: %v / %v",
			buyPrice, sellPrice, profit, buyErr, sellErr)
		t.emailSender.Send(subject, content)
	} else {
		log.Printf("[Arbitrage] Trade executed successfully! Direction: %s, Profit: %.8f", direction, profit)
		subject := fmt.Sprintf("[搬砖] 交易成功! 方向: %s", direction)
		content := fmt.Sprintf("买入价: %.8f, 卖出价: %.8f, 数量: %.4f, 利润: %.8f",
			buyPrice, sellPrice, amount, profit)
		t.emailSender.Send(subject, content)
	}

	log.Printf("[Arbitrage] %s -- 已搬砖次数: %d -- 检测到可搬砖次数: %d",
		time.Now().Format("2006-01-02 15:04:05"), t.counter, t.times)
}

// CheckBalance checks account balances
func (t *ArbitrageTrader) CheckBalance() (hbNAS, hbUSDT, gateNAS, gateUSDT float64, err error) {
	// Get Huobi balance
	hbBalance, err := t.huobiClient.FetchBalance()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	for _, b := range hbBalance.List {
		if b.Currency == "nas" && b.Type == "trade" {
			hbNAS, _ = strconv.ParseFloat(b.Balance, 64)
		}
		if b.Currency == "usdt" && b.Type == "trade" {
			hbUSDT, _ = strconv.ParseFloat(b.Balance, 64)
		}
	}

	// Get Gate balance
	gateBalance, err := t.gateClient.FetchBalance()
	if err != nil {
		return 0, 0, 0, 0, err
	}

	if nas, ok := gateBalance.Available["NAS"]; ok {
		gateNAS, _ = strconv.ParseFloat(nas, 64)
	}
	if usdt, ok := gateBalance.Available["USDT"]; ok {
		gateUSDT, _ = strconv.ParseFloat(usdt, 64)
	}

	return hbNAS, hbUSDT, gateNAS, gateUSDT, nil
}

// round rounds a float to specified decimal places
func round(x float64, decimals int) float64 {
	p := math.Pow10(decimals)
	return math.Round(x*p) / p
}