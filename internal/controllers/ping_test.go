package controllers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/NodeFactoryIo/vedran/internal/auth"
	"github.com/NodeFactoryIo/vedran/internal/models"
	mocks "github.com/NodeFactoryIo/vedran/mocks/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestApiController_PingHandler(t *testing.T) {
	tests := []struct {
		name string
		onPingSaveReturn interface{}
		expectedStatus int
	}{
		{name: "Valid ping request", onPingSaveReturn: nil, expectedStatus: http.StatusOK},
		{name: "Valid ping request", onPingSaveReturn: errors.New("DB error"), expectedStatus: http.StatusInternalServerError},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			timestamp := time.Now()
			// create mock controller
			nodeRepoMock := mocks.NodeRepository{}
			pingRepoMock := mocks.PingRepository{}
			metricsRepoMock := mocks.MetricsRepository{}
			pingRepoMock.On("Save", &models.Ping{
				NodeId:    "1",
				Timestamp: timestamp,
			}).Return(test.onPingSaveReturn)
			apiController := NewApiController(false, &nodeRepoMock, &pingRepoMock, &metricsRepoMock)
			handler := http.HandlerFunc(apiController.PingHandler)

			// create test request and populate context
			req, _ := http.NewRequest("POST", "/api/v1/node", bytes.NewReader(nil))
			c := &auth.RequestContext{
				NodeId:    "1",
				Timestamp: timestamp,
			}
			ctx := context.WithValue(req.Context(), auth.RequestContextKey, c)
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			// invoke test request
			handler.ServeHTTP(rr, req)

			// asserts
			assert.Equal(t, test.expectedStatus, rr.Code, fmt.Sprintf("Response status code should be %d", http.StatusOK))
			assert.True(t, pingRepoMock.AssertNumberOfCalls(t, "Save", 1))
		})
	}
}
