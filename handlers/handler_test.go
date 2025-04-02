package handlers

import (
	"awesomeProject/models"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockService struct {
	processTransactionFunc func(ctx context.Context, id string, userid uint64, amount uint64, state string, source string) error
	getBalanceFunc         func(ctx context.Context, id uint64) (uint64, error)
}

func (m *mockService) ProcessTransaction(ctx context.Context, id string, userid uint64, amount uint64, state string, source string) error {
	return m.processTransactionFunc(ctx, id, userid, amount, state, source)
}

func (m *mockService) GetBalance(ctx context.Context, id uint64) (uint64, error) {
	return m.getBalanceFunc(ctx, id)
}

func (m *mockService) UpdateBalance(ctx context.Context, userId uint64, amount uint64, state string) error {
	return nil
}

func (m *mockService) ValidateBalance(ctx context.Context, userId uint64, amount uint64, state string) error {
	return nil
}

func TestGetBalance(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockService    mockService
		expectedStatus int
		expectedBody   *models.UserBalance
	}{
		{
			name:   "successful balance retrieval",
			userID: "1",
			mockService: mockService{
				getBalanceFunc: func(ctx context.Context, id uint64) (uint64, error) {
					return 1000, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: &models.UserBalance{
				ID:      1,
				Balance: "10.00",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&tt.mockService)
			req := httptest.NewRequest(http.MethodGet, "/user/"+tt.userID, nil)
			w := httptest.NewRecorder()

			handler.GetBalance(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != nil {
				var got models.UserBalance
				err := json.NewDecoder(w.Body).Decode(&got)
				if err != nil {
					t.Fatalf("Failed to decode response body: %v", err)
				}

				if got.ID != tt.expectedBody.ID || got.Balance != tt.expectedBody.Balance {
					t.Errorf("Expected body %v, got %v", tt.expectedBody, got)
				}
			}
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    models.TransactionRequest
		sourceType     string
		mockService    mockService
		expectedStatus int
	}{
		{
			name:   "successful transaction",
			userID: "1",
			requestBody: models.TransactionRequest{
				State:         "win",
				Amount:        "10.00",
				TransactionId: "test-1",
			},
			sourceType: "game",
			mockService: mockService{
				processTransactionFunc: func(ctx context.Context, id string, userid uint64, amount uint64, state string, source string) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&tt.mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/user/"+tt.userID+"/transaction", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Source-Type", tt.sourceType)
			w := httptest.NewRecorder()

			handler.CreateTransaction(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
