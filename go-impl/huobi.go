package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

// HuobiClient handles Huobi API interactions
type HuobiClient struct {
	accessKey string
	secretKey string
}

// NewHuobiClient creates a new Huobi client from environment variables
func NewHuobiClient() *HuobiClient {
	return &HuobiClient{
		accessKey: os.Getenv("HUOBI_ACCESS_KEY"),
		secretKey: os.Getenv("HUOBI_SECRET_KEY"),
	}
}

// OrderBook represents the order book data
type OrderBook struct {
	Bids [][]interface{} `json:"bids"`
	Asks [][]interface{} `json:"asks"`
}

// OrderBookResponse represents the API response for order book
type OrderBookResponse struct {
	Status string    `json:"status"`
	Data   OrderBook `json:"data"`
}

// Balance represents balance information
type Balance struct {
	Currency string `json:"currency"`
	Type     string `json:"type"`
	Balance  string `json:"balance"`
}

// BalanceData represents balance data from API
type BalanceData struct {
	State   string    `json:"state"`
	List    []Balance `json:"list"`
}

// BalanceResponse represents the API response for balance
type BalanceResponse struct {
	Status string       `json:"status"`
	Data   *BalanceData `json:"data"`
}

// Order represents order information
type Order struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Symbol    string `json:"symbol"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	State     string `json:"state"`
	CreatedAt int64  `json:"created-at"`
}

// OrderResponse represents the API response for order
type OrderResponse struct {
	Status string `json:"status"`
	Data   Order  `json:"data"`
}

// sign generates HMAC-SHA256 signature for Huobi API
func (c *HuobiClient) sign(params map[string]string) string {
	// Sort parameters by key
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build query string
	var paramStrs []string
	for _, k := range keys {
		paramStrs = append(paramStrs, k+"="+params[k])
	}
	queryString := strings.Join(paramStrs, "&")

	// Create signature
	signString := fmt.Sprintf("POST\napi.huobipro.com\n/v1/order/orders\n%s", queryString)
	h := hmac.New(sha256.New, []byte(c.secretKey))
	h.Write([]byte(signString))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// FetchOrderBook retrieves order book for a trading pair
func (c *HuobiClient) FetchOrderBook(symbol string, limit int) (*OrderBook, error) {
	url := fmt.Sprintf("https://api.huobipro.com/market/depth?symbol=%s&type=step0", strings.ToLower(symbol))

	if limit > 0 {
		url += fmt.Sprintf("&depth=%d", limit)
	}

	params := map[string]string{
		"symbol": strings.ToLower(symbol),
		"type":   "step0",
	}

	// For public endpoints, no signature needed
	resp, err := fasthttp.Get(nil, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Release()

	var result OrderBookResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	if result.Status != "ok" {
		return nil, fmt.Errorf("API error: %s", result.Status)
	}

	return &result.Data, nil
}

// FetchBalance retrieves account balance
func (c *HuobiClient) FetchBalance() (*BalanceData, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	params := map[string]string{
		"AccessKeyId": c.accessKey,
		"SignatureMethod": "HmacSHA256",
		"SignatureVersion": "2",
		"Timestamp": timestamp,
	}

	params["Signature"] = c.sign(params)

	url := "https://api.huobipro.com/v1/account/accounts?" + buildQuery(params)

	req := fasthttp.NewRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")

	resp, err := fasthttp.Do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Release()

	var result BalanceResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	if result.Status != "ok" {
		return nil, fmt.Errorf("API error: %s", result.Status)
	}

	return result.Data, nil
}

// CreateMarketSellOrder creates a market sell order
func (c *HuobiClient) CreateMarketSellOrder(symbol string, amount float64) (*Order, error) {
	return c.createOrder(symbol, "sell-market", "", amount)
}

// CreateMarketBuyOrder creates a market buy order
func (c *HuobiClient) CreateMarketBuyOrder(symbol string, amount float64) (*Order, error) {
	return c.createOrder(symbol, "buy-market", "", amount)
}

// CreateLimitBuyOrder creates a limit buy order
func (c *HuobiClient) CreateLimitBuyOrder(symbol string, amount float64, price float64) (*Order, error) {
	return c.createOrder(symbol, "buy-limit", price, amount)
}

// CreateLimitSellOrder creates a limit sell order
func (c *HuobiClient) CreateLimitSellOrder(symbol string, amount float64, price float64) (*Order, error) {
	return c.createOrder(symbol, "sell-limit", price, amount)
}

// createOrder creates a generic order
func (c *HuobiClient) createOrder(symbol, orderType string, price string, amount float64) (*Order, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	// Convert symbol to huobi format (nasusdt)
	symbolLower := strings.ToLower(symbol)

	params := map[string]string{
		"AccessKeyId":     c.accessKey,
		"SignatureMethod": "HmacSHA256",
		"SignatureVersion": "2",
		"Timestamp":       timestamp,
		"symbol":          symbolLower,
		"type":            orderType,
		"amount":          strconv.FormatFloat(amount, 'f', -1, 64),
	}

	if price != "" {
		params["price"] = price
	}

	params["Signature"] = c.sign(params)

	url := "https://api.huobipro.com/v1/order/orders?" + buildQuery(params)

	req := fasthttp.NewRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")

	resp, err := fasthttp.Do(req, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Release()

	var result OrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	if result.Status != "ok" {
		return nil, fmt.Errorf("API error: %s", result.Status)
	}

	return &result.Data, nil
}

// buildQuery builds query string from params
func buildQuery(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var query []string
	for _, k := range keys {
		query = append(query, k+"="+urlEncode(params[k]))
	}
	return strings.Join(query, "&")
}

// urlEncode URL encodes a string
func urlEncode(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(s, "+", "%20"),
		"*", "%2A",
	)
}