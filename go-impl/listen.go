package main

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketConfig holds WebSocket connection configuration
type WebSocketConfig struct {
	URL           string
	ReconnectDelay time.Duration
	MaxRetries    int
}

// KLine represents a K-line/candlestick data point
type KLine struct {
	ID     int64   `json:"id"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Amount float64 `json:"amount"`
	Vol    float64 `json:"vol"`
	Count  int     `json:"count"`
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Channel string          `json:"ch"`
	Data    json.RawMessage `json:"tick"`
}

// WSPing represents a ping message
type WSPing struct {
	Ping int64 `json:"ping"`
}

// WSSubscribe represents a subscription request
type WSSubscribe struct {
	Sub string `json:"sub"`
	ID  string `json:"id"`
}

// HuobiWSClient handles Huobi WebSocket connections
type HuobiWSClient struct {
	conn       *websocket.Conn
	config     WebSocketConfig
	mu         sync.Mutex
	connected  bool
	stopChan   chan struct{}
	doneChan   chan struct{}
	subs       map[string]bool
	callback   func(symbol string, kline KLine)
}

// NewHuobiWSClient creates a new Huobi WebSocket client
func NewHuobiWSClient(callback func(symbol string, kline KLine)) *HuobiWSClient {
	return &HuobiWSClient{
		config: WebSocketConfig{
			URL:           "wss://api.huobipro.com/ws",
			ReconnectDelay: 5 * time.Second,
			MaxRetries:    10,
		},
		stopChan: make(chan struct{}),
		doneChan: make(chan struct{}),
		subs:     make(map[string]bool),
		callback: callback,
	}
}

// Connect establishes WebSocket connection with retry logic
func (c *HuobiWSClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.connected {
		return nil
	}

	var err error
	var conn *websocket.Conn

	for retries := 0; retries < c.config.MaxRetries; retries++ {
		conn, _, err = websocket.DefaultDialer.Dial(c.config.URL, nil)
		if err == nil {
			c.conn = conn
			c.connected = true
			log.Println("WebSocket connected")
			return nil
		}
		log.Printf("Connection failed (attempt %d/%d): %v", retries+1, c.config.MaxRetries, err)
		time.Sleep(c.config.ReconnectDelay)
	}

	return fmt.Errorf("failed to connect after %d retries: %w", c.config.MaxRetries, err)
}

// Disconnect closes the WebSocket connection
func (c *HuobiWSClient) Disconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return nil
	}

	close(c.stopChan)
	<-c.doneChan

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Subscribe subscribes to a channel
func (c *HuobiWSClient) Subscribe(symbol string, interval string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	channel := fmt.Sprintf("market.%s.kline.%s", strings.ToLower(symbol), interval)
	c.subs[channel] = true

	subMsg := WSSubscribe{
		Sub: channel,
		ID:  "id" + fmt.Sprintf("%d", time.Now().Unix()),
	}

	data, err := json.Marshal(subMsg)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// Start starts the WebSocket message handler
func (c *HuobiWSClient) Start() error {
	if err := c.Connect(); err != nil {
		return err
	}

	go c.readLoop()
	go c.pingLoop()

	return nil
}

// readLoop handles incoming WebSocket messages
func (c *HuobiWSClient) readLoop() {
	defer func() {
		close(c.doneChan)
		c.mu.Lock()
		c.connected = false
		c.mu.Unlock()
	}()

	for {
		select {
		case <-c.stopChan:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}

			c.handleMessage(message)
		}
	}
}

// pingLoop sends periodic ping messages
func (c *HuobiWSClient) pingLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return
		case <-ticker.C:
			pingMsg := fmt.Sprintf(`{"ping":%d}`, time.Now().UnixMilli())
			c.mu.Lock()
			if c.connected && c.conn != nil {
				c.conn.WriteMessage(websocket.TextMessage, []byte(pingMsg))
			}
			c.mu.Unlock()
		}
	}
}

// handleMessage processes incoming WebSocket messages
func (c *HuobiWSClient) handleMessage(message []byte) {
	// Try to decompress gzip
	reader, err := gzip.NewReader(strings.NewReader(string(message)))
	if err == nil {
		defer reader.Close()
		scanner := bufio.NewScanner(reader)
		if scanner.Scan() {
			message = scanner.Bytes()
		}
	}

	// Check for ping
	var ping WSPing
	if err := json.Unmarshal(message, &ping); err == nil {
		if ping.Ping > 0 {
			pong := fmt.Sprintf(`{"pong":%d}`, ping.Ping)
			c.mu.Lock()
			if c.connected && c.conn != nil {
				c.conn.WriteMessage(websocket.TextMessage, []byte(pong))
			}
			c.mu.Unlock()
			// Resubscribe to channels
			for ch := range c.subs {
				subMsg := WSSubscribe{Sub: ch, ID: "id" + fmt.Sprintf("%d", time.Now().Unix())}
				data, _ := json.Marshal(subMsg)
				c.conn.WriteMessage(websocket.TextMessage, data)
			}
		}
		return
	}

	// Parse regular message
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	if wsMsg.Channel == "" {
		return
	}

	// Extract symbol from channel (e.g., "market.nasusdt.kline.1min")
	parts := strings.Split(wsMsg.Channel, ".")
	if len(parts) >= 2 {
		symbol := parts[1]

		// Parse K-line data
		var kline KLine
		if err := json.Unmarshal(wsMsg.Data, &kline); err == nil {
			if c.callback != nil {
				c.callback(symbol, kline)
			}
		}
	}
}

// WaitForSignal waits for interrupt signal to shutdown
func (c *HuobiWSClient) WaitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
	c.Disconnect()
}