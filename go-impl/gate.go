package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

// GateClient handles Gate.io API interactions
type GateClient struct {
	accessKey string
	secretKey string
}

// NewGateClient creates a new Gate.io client from environment variables
func NewGateClient() *GateClient {
	return &GateClient{
		accessKey: os.Getenv("GATE_ACCESS_KEY"),
		secretKey: os.Getenv("GATE_SECRET_KEY"),
	}
}

// GateOrderBook represents Gate.io order book data
type GateOrderBook struct {
	Bids [][]interface{} `json:"bids"`
	Asks [][]interface{} `json:"asks"`
}

// GateOrderBookResponse represents the API response for order book
type GateOrderBookResponse struct {
	Method string        `json:"method"`
	Params GateOrderBook `json:"params"`
	ID     int           `json:"id"`
}

// GateBalance represents balance information
type GateBalance struct {
	Available map[string]string `json:"available"`
	Locked    map[string]string `json:"locked"`
}

// GateBalanceResponse represents the API response for balance
type GateBalanceResponse struct {
	Method string       `json:"method"`
	Params []interface{} `json:"params"`
	ID     int          `json:"id"`
}

// GateOrder represents order information
type GateOrder struct {
	ID        string `json:"id"`
	Currency  string `json:"currency"`
	Type      string `json:"type"`
	Side      string `json:"side"`
	Amount    string `json:"amount"`
	Price     string `json:"price"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"create_time"`
}

// GateOrderResponse represents the API response for order
type GateOrderResponse struct {
	ID     int        `json:"id"`
	Result bool       `json:"result"`
	Order  *GateOrder `json:"order"`
}

// sign generates HMAC-SHA512 signature for Gate.io API
func (c *GateClient) sign(params map[string]string, queryString string) string {
	signString := queryString
	h := hmac.New(sha512.New, []byte(c.secretKey))
	h.Write([]byte(signString))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// FetchOrderBook retrieves order book for a trading pair
func (c *GateClient) FetchOrderBook(symbol string, limit int) (*GateOrderBook, error) {
	url := fmt.Sprintf("https://api.gateio.ws/api2/1/orderBook/%s", symbol)

	if limit > 0 {
		url += fmt.Sprintf("?limit=%d", limit)
	}

	resp, err := fasthttp.Get(nil, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Release()

	var result GateOrderBook
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// FetchBalance retrieves account balance
func (c *GateClient) FetchBalance() (*GateBalance, error) {
	timestamp := time.Now().Unix()

	params := map[string]string{
		"method": "balance_info",
		"api_key": c.accessKey,
	}

	// Build query string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var paramStrs []string
	for _, k := range keys {
		paramStrs = append(paramStrs, k+"="+params[k])
	}
	queryString := "?" + strings.Join(paramStrs, "&") + fmt.Sprintf("&sign=%s&timestamp=%d", c.sign(params, queryString), timestamp)

	url := "https://api.gateio.ws/api2/1/spot/balance_info" + queryString

	req := fasthttp.NewRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")

	resp, err := fasthttp.Do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Release()

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	available := make(map[string]string)
	if avail, ok := result["available"].(map[string]interface{}); ok {
		for k, v := range avail {
			available[k] = fmt.Sprintf("%v", v)
		}
	}

	locked := make(map[string]string)
	if lock, ok := result["locked"].(map[string]interface{}); ok {
		for k, v := range lock {
			locked[k] = fmt.Sprintf("%v", v)
		}
	}

	return &GateBalance{
		Available: available,
		Locked:    locked,
	}, nil
}

// CreateMarketBuyOrder creates a market buy order
func (c *GateClient) CreateMarketBuyOrder(symbol string, amount float64) (*GateOrder, error) {
	return c.createOrder(symbol, "buy", "market", amount, 0)
}

// CreateMarketSellOrder creates a market sell order
func (c *GateClient) CreateMarketSellOrder(symbol string, amount float64) (*GateOrder, error) {
	return c.createOrder(symbol, "sell", "market", amount, 0)
}

// CreateLimitBuyOrder creates a limit buy order
func (c *GateClient) CreateLimitBuyOrder(symbol string, amount float64, price float64) (*GateOrder, error) {
	return c.createOrder(symbol, "buy", "limit", amount, price)
}

// CreateLimitSellOrder creates a limit sell order
func (c *GateClient) CreateLimitSellOrder(symbol string, amount float64, price float64) (*GateOrder, error) {
	return c.createOrder(symbol, "sell", "limit", amount, price)
}

// createOrder creates an order on Gate.io
func (c *GateClient) createOrder(symbol, side, orderType string, amount float64, price float64) (*GateOrder, error) {
	timestamp := time.Now().Unix()

	params := map[string]string{
		"currency":  symbol,
		"type":      orderType,
		"side":      side,
		"amount":    fmt.Sprintf("%f", amount),
		"api_key":   c.accessKey,
	}

	if price > 0 {
		params["price"] = fmt.Sprintf("%f", price)
	}

	// Build query string
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var paramStrs []string
	for _, k := range keys {
		paramStrs = append(paramStrs, k+"="+params[k])
	}
	queryString := strings.Join(paramStrs, "&")
	queryString += fmt.Sprintf("&sign=%s&timestamp=%d", c.sign(params, queryString), timestamp)

	url := "https://api.gateio.ws/api2/1/spot/orders"

	req := fasthttp.NewRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBody(queryString)

	resp, err := fasthttp.Do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Release()

	var result GateOrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	if !result.Result {
		return nil, fmt.Errorf("order creation failed: %s", string(resp.Body()))
	}

	return result.Order, nil
}