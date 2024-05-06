// Package handlers_test contains tests for the handlers package.
package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pavlegich/scripts-hub/internal/controllers/handlers"
	"github.com/pavlegich/scripts-hub/internal/entities"
	"github.com/pavlegich/scripts-hub/internal/infra/config"
	"github.com/pavlegich/scripts-hub/internal/mocks"
	"github.com/stretchr/testify/require"
)

func TestCommandHandler_HandleCreateCommand(t *testing.T) {
	ctx := context.Background()

	// Initialize mock repository
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mocks.NewMockRepository(mockCtrl)

	cfg := &config.Config{
		Address: `localhost:8080`,
	}

	type expected struct {
		cmd *entities.Command
		err error
	}
	type args struct {
		req *entities.Command
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		wantCode int
		wantBody string
	}{
		{
			name: "success",
			args: args{
				req: &entities.Command{
					Name:   "pwd",
					Script: "pwd",
				},
			},
			expected: expected{
				cmd: &entities.Command{
					ID:     1,
					Name:   "pwd",
					Script: "pwd",
				},
				err: nil,
			},
			wantCode: http.StatusCreated,
			wantBody: `{"command_id": 1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock expected response
			mockRepo.EXPECT().CreateCommand(gomock.Any(), gomock.Any()).
				Return(tt.expected.cmd, tt.expected.err).Times(1)

			// Controller
			ctrl := handlers.NewController(ctx, mockRepo, cfg)
			mh, err := ctrl.BuildRoute(ctx)
			require.NoError(t, err)

			// Form new request
			url := `http://` + cfg.Address + `/command`

			reqBody, err := json.Marshal(tt.args.req)
			require.NoError(t, err)
			r := httptest.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()

			mh.ServeHTTP(w, r)

			// Get response
			resp := w.Result()
			defer resp.Body.Close()

			gotBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			// Check status code
			require.Equal(t, tt.wantCode, resp.StatusCode)
			require.JSONEq(t, tt.wantBody, string(gotBody))
		})
	}
}
