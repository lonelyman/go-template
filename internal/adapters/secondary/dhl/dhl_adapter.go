package dhl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DHLAdapter handles communication with DHL API
type DHLAdapter struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewDHLAdapter creates a new DHL adapter
func NewDHLAdapter(baseURL, apiKey string) *DHLAdapter {
	return &DHLAdapter{
		baseURL: baseURL,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// TrackingRequest represents a tracking request
type TrackingRequest struct {
	TrackingNumber string `json:"trackingNumber"`
}

// TrackingResponse represents a tracking response
type TrackingResponse struct {
	TrackingNumber string `json:"trackingNumber"`
	Status         string `json:"status"`
	Location       string `json:"location"`
	EstimatedDate  string `json:"estimatedDate"`
}

// TrackShipment tracks a shipment using DHL API
func (d *DHLAdapter) TrackShipment(trackingNumber string) (*TrackingResponse, error) {
	req := TrackingRequest{
		TrackingNumber: trackingNumber,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", d.baseURL+"/track", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+d.apiKey)

	resp, err := d.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DHL API returned status: %d", resp.StatusCode)
	}

	var trackingResp TrackingResponse
	if err := json.NewDecoder(resp.Body).Decode(&trackingResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &trackingResp, nil
}
