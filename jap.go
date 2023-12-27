package jap

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// JAPClient is a client for the JustAnotherPanel API.
type JAPClient struct {
	key      string
	endpoint string
}

// New creates a new JAPClient with the given API key.
func New(key string) JAPClient {
	return JAPClient{
		key:      key,
		endpoint: "https://justanotherpanel.com/api/v2",
	}
}

// Service represents the structure of each service in the API response.
type Service struct {
	Service  string `json:"service"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Category string `json:"category"`
	Rate     string `json:"rate"`
	Min      string `json:"min"`
	Max      string `json:"max"`
	Refill   bool   `json:"refill"`
	Cancel   bool   `json:"cancel"`
}

// ListServices retrieves the list of services from the API.
func (c *JAPClient) ListServices() ([]Service, error) {
	body := struct {
		Key    string `json:"key"`
		Action string `json:"action"`
	}{
		Key:    c.key,
		Action: "services",
	}
	bytes, err := c.post(body)
	if err != nil {
		return nil, err
	}

	var response []Service
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// AddOrder adds an order with the given parameters and returns the order ID as a string.
func (c *JAPClient) AddOrder(service, link string, quantity int, runs, interval *int) (string, error) {
	orderRequest := struct {
		Key      string `json:"key"`
		Action   string `json:"action"`
		Service  string `json:"service"`
		Link     string `json:"link"`
		Quantity int    `json:"quantity"`
		Runs     *int   `json:"runs,omitempty"`
		Interval *int   `json:"interval,omitempty"`
	}{
		Key:      c.key,
		Action:   "add",
		Service:  service,
		Link:     link,
		Quantity: quantity,
		Runs:     runs,
		Interval: interval,
	}

	bytes, err := c.post(orderRequest)
	if err != nil {
		return "", err
	}

	var response struct {
		OrderID int `json:"order"`
	}
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(response.OrderID), nil
}

// OrderStatusResponse represents the JSON structure of the response for the order status request.
type OrderStatusResponse struct {
	OrderStatus map[string]OrderStatus `json:"orderStatus"`
}

// GetOrderStatus checks the status of an order with the given order ID and returns the status.
func (c *JAPClient) GetOrderStatus(orderID string) (OrderStatusResponse, error) {
	body := struct {
		Key    string `json:"key"`
		Action string `json:"action"`
		Order  string `json:"order"`
	}{
		Key:    c.key,
		Action: "status",
		Order:  orderID,
	}
	bytes, err := c.post(body)
	if err != nil {
		return OrderStatusResponse{}, err
	}

	var response OrderStatusResponse
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return OrderStatusResponse{}, err
	}

	return response, nil
}

// GetUserBalance retrieves the user's balance from the API.
func (c *JAPClient) GetUserBalance() (UserBalanceResponse, error) {
	body := struct {
		Key    string `json:"key"`
		Action string `json:"action"`
	}{
		Key:    c.key,
		Action: "balance",
	}
	bytes, err := c.post(body)
	if err != nil {
		return UserBalanceResponse{}, err
	}

	var response UserBalanceResponse
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return UserBalanceResponse{}, err
	}

	return response, nil
}

// post is a helper method to perform POST requests for the JAPClient.
func (c *JAPClient) post(body interface{}) ([]byte, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// The response type will depend on the method calling post, so we return the raw JSON
	// and let the calling method handle unmarshalling.
	return responseBody, nil
}

// OrderStatus details for an order.
type OrderStatus struct {
	Charge     string `json:"charge,omitempty"`
	StartCount string `json:"start_count,omitempty"`
	Status     string `json:"status,omitempty"`
	Remains    string `json:"remains,omitempty"`
	Currency   string `json:"currency,omitempty"`
	Error      string `json:"error,omitempty"`
}

// UserBalanceResponse represents the JSON structure of the response for the user balance request.
type UserBalanceResponse struct {
	Balance  string `json:"balance"`
	Currency string `json:"currency"`
}

func (c *JAPClient) RedditUpvote(link string, quantity int) (string, error) {
	return c.AddOrder("6228", link, quantity, nil, nil)
}
