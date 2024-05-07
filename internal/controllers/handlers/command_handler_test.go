// Package handlers_test contains tests for the handlers package.
package handlers_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pavlegich/scripts-hub/internal/controllers/handlers"
	"github.com/pavlegich/scripts-hub/internal/entities"
	errs "github.com/pavlegich/scripts-hub/internal/errors"
	"github.com/pavlegich/scripts-hub/internal/infra/config"
	"github.com/pavlegich/scripts-hub/internal/mocks"
	"github.com/stretchr/testify/require"
)

func TestCommandHandler_HandleCommand(t *testing.T) {
	ctx := context.Background()

	cfg := &config.Config{
		Address: `localhost:8080`,
	}

	type args struct {
		method string
	}
	tests := []struct {
		name     string
		args     args
		wantCode int
	}{
		{
			name: "method_not_allowed",
			args: args{
				method: http.MethodPatch,
			},
			wantCode: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Controller
			ctrl := handlers.NewController(ctx, cfg)
			mh, err := ctrl.BuildRoute(ctx, nil, nil)
			require.NoError(t, err)

			// Form new request
			url := `http://` + cfg.Address + `/command`

			r := httptest.NewRequest(tt.args.method, url, nil)
			w := httptest.NewRecorder()

			mh.ServeHTTP(w, r)

			// Get response
			resp := w.Result()
			defer resp.Body.Close()

			// Check status code
			require.Equal(t, tt.wantCode, resp.StatusCode)
		})
	}
}

func TestCommandHandler_HandleCreateCommand(t *testing.T) {
	ctx := context.Background()

	// Initialize mock repository
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mocks.NewMockRepository(mockCtrl)

	cfg := &config.Config{
		Address:   `localhost:8080`,
		RateLimit: 1,
	}

	type expCreate struct {
		want bool
		cmd  *entities.Command
		err  error
	}
	type expAppend struct {
		want bool
		err  error
	}
	type expected struct {
		create expCreate
		append expAppend
	}
	type args struct {
		reqBody string
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
				reqBody: `{"name": "pwd", "script": "pwd"}`,
			},
			expected: expected{
				create: expCreate{
					want: true,
					cmd: &entities.Command{
						ID:     1,
						Name:   "pwd",
						Script: "pwd",
					},
					err: nil,
				},
				append: expAppend{
					want: true,
					err:  nil,
				},
			},
			wantCode: http.StatusCreated,
			wantBody: `{"command_id": 1}`,
		},
		{
			name: "incorrect_body",
			args: args{
				reqBody: `{{"name": "pwd", "script": "pwd"}`,
			},
			expected: expected{},
			wantCode: http.StatusBadRequest,
			wantBody: ``,
		},
		{
			name: "empty_name_and_script",
			args: args{
				reqBody: `{}`,
			},
			expected: expected{},
			wantCode: http.StatusBadRequest,
			wantBody: ``,
		},
		{
			name: "unknown_command",
			args: args{
				reqBody: `{"name": "unknown", "script": "unknown"}`,
			},
			expected: expected{},
			wantCode: http.StatusBadRequest,
			wantBody: ``,
		},
		{
			name: "command_already_exists",
			args: args{
				reqBody: `{"name": "exists", "script": "ls"}`,
			},
			expected: expected{
				create: expCreate{
					want: true,
					cmd:  nil,
					err:  errs.ErrCmdAlreadyExists,
				},
			},
			wantCode: http.StatusConflict,
			wantBody: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mocks expected response
			if tt.expected.create.want {
				mockRepo.EXPECT().CreateCommand(gomock.Any(), gomock.Any()).
					Return(tt.expected.create.cmd, tt.expected.create.err).Times(1)
			}
			if tt.expected.append.want {
				mockRepo.EXPECT().AppendCommandOutputByName(gomock.Any(), gomock.Any()).
					Return(tt.expected.append.err).Times(1)
			}

			// Controller
			ctrl := handlers.NewController(ctx, cfg)
			ch := make(chan entities.Command)
			mh, err := ctrl.BuildRoute(ctx, mockRepo, ch)
			require.NoError(t, err)

			// Form new request
			url := `http://` + cfg.Address + `/command`

			r := httptest.NewRequest(http.MethodPost, url, bytes.NewBufferString(tt.args.reqBody))
			w := httptest.NewRecorder()

			mh.ServeHTTP(w, r)
			time.Sleep(100 * time.Millisecond)

			// Get response
			resp := w.Result()
			gotBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Check status code
			require.Equal(t, tt.wantCode, resp.StatusCode)
			if !(tt.wantBody == ``) {
				require.JSONEq(t, tt.wantBody, string(gotBody))
			}
		})
	}
}

func TestCommandHandler_HandleGetCommand(t *testing.T) {
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
	type query struct {
		want   bool
		key    string
		values []string
	}
	type args struct {
		query query
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
				query: query{
					want:   true,
					key:    "name",
					values: []string{"pwd"},
				},
			},
			expected: expected{
				cmd: &entities.Command{
					ID:     1,
					Name:   "pwd",
					Script: "pwd",
					Output: "/path",
				},
				err: nil,
			},
			wantCode: http.StatusOK,
			wantBody: `{"id": 1, "name": "pwd", "script": "pwd", "output": "/path"}`,
		},
		{
			name: "incorrect_query_key",
			args: args{
				query: query{
					want:   true,
					key:    "incorrect",
					values: []string{"pwd"},
				},
			},
			expected: expected{},
			wantCode: http.StatusBadRequest,
			wantBody: ``,
		},
		{
			name:     "no_query",
			args:     args{},
			expected: expected{},
			wantCode: http.StatusBadRequest,
			wantBody: ``,
		},
		{
			name: "incorrect_query_values_count",
			args: args{
				query: query{
					want:   true,
					key:    "name",
					values: []string{"pwd", "second"},
				},
			},
			expected: expected{},
			wantCode: http.StatusBadRequest,
			wantBody: ``,
		},
		{
			name: "command_not_found",
			args: args{
				query: query{
					want:   true,
					key:    "name",
					values: []string{"unknown"},
				},
			},
			expected: expected{
				cmd: nil,
				err: errs.ErrCmdNotFound,
			},
			wantCode: http.StatusNotFound,
			wantBody: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock expected response
			if !reflect.DeepEqual(tt.expected, expected{}) {
				mockRepo.EXPECT().GetCommandByName(gomock.Any(), gomock.Any()).
					Return(tt.expected.cmd, tt.expected.err).Times(1)
			}

			// Controller
			ctrl := handlers.NewController(ctx, cfg)
			mh, err := ctrl.BuildRoute(ctx, mockRepo, nil)
			require.NoError(t, err)

			// Form new request
			url := `http://` + cfg.Address + `/command`

			r := httptest.NewRequest(http.MethodGet, url, nil)
			if tt.args.query.want {
				q := r.URL.Query()
				for _, v := range tt.args.query.values {
					q.Add(tt.args.query.key, v)
				}
				r.URL.RawQuery = q.Encode()
			}
			w := httptest.NewRecorder()

			mh.ServeHTTP(w, r)

			// Get response
			resp := w.Result()
			gotBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Check status code
			require.Equal(t, tt.wantCode, resp.StatusCode)
			if !(tt.wantBody == ``) {
				require.JSONEq(t, tt.wantBody, string(gotBody))
			}
		})
	}
}

func TestCommandHandler_HandleCommands(t *testing.T) {
	ctx := context.Background()

	// Initialize mock repository
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mocks.NewMockRepository(mockCtrl)

	cfg := &config.Config{
		Address: `localhost:8080`,
	}

	type expected struct {
		cmds []*entities.Command
		err  error
	}
	type args struct {
		method string
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
				method: http.MethodGet,
			},
			expected: expected{
				cmds: []*entities.Command{
					{
						ID:     1,
						Name:   "pwd",
						Script: "pwd",
						Output: "/path",
					},
					{
						ID:     2,
						Name:   "ls",
						Script: "ls",
						Output: "Documents",
					},
				},
				err: nil,
			},
			wantCode: http.StatusOK,
			wantBody: `[{"id": 1, "name": "pwd", "script": "pwd", "output": "/path"}, 
			{"id": 2, "name": "ls", "script": "ls", "output": "Documents"}]`,
		},
		{
			name: "incorrect_method",
			args: args{
				method: http.MethodPost,
			},
			expected: expected{},
			wantCode: http.StatusMethodNotAllowed,
			wantBody: ``,
		},
		{
			name: "commands_not_found",
			args: args{
				method: http.MethodGet,
			},
			expected: expected{
				err: errs.ErrCmdNotFound,
			},
			wantCode: http.StatusNotFound,
			wantBody: ``,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock expected response
			if !reflect.DeepEqual(tt.expected, expected{}) {
				mockRepo.EXPECT().GetAllCommands(gomock.Any()).
					Return(tt.expected.cmds, tt.expected.err).Times(1)
			}

			// Controller
			ctrl := handlers.NewController(ctx, cfg)
			mh, err := ctrl.BuildRoute(ctx, mockRepo, nil)
			require.NoError(t, err)

			// Form new request
			url := `http://` + cfg.Address + `/commands`

			r := httptest.NewRequest(tt.args.method, url, nil)
			w := httptest.NewRecorder()

			mh.ServeHTTP(w, r)

			// Get response
			resp := w.Result()
			gotBody, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			defer resp.Body.Close()

			// Check status code
			require.Equal(t, tt.wantCode, resp.StatusCode)
			if !(tt.wantBody == ``) {
				require.JSONEq(t, tt.wantBody, string(gotBody))
			}
		})
	}
}

func TestCommandHandler_HandleDeleteCommand(t *testing.T) {
	ctx := context.Background()

	// Initialize mock repository
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockRepo := mocks.NewMockRepository(mockCtrl)

	cfg := &config.Config{
		Address:   `localhost:8080`,
		RateLimit: 1,
	}

	type query struct {
		want   bool
		key    string
		values []string
	}
	type args struct {
		query   query
		reqBody string
	}
	type expCreate struct {
		want bool
		cmd  *entities.Command
		err  error
	}
	type expAppend struct {
		want bool
		err  error
	}
	type expDelete struct {
		want bool
		err  error
	}
	type expected struct {
		create expCreate
		append expAppend
		delete expDelete
	}
	type wantCode struct {
		create int
		delete int
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		wantCode wantCode
	}{
		{
			name: "success",
			args: args{
				query: query{
					want:   true,
					key:    "name",
					values: []string{"pwd"},
				},
				reqBody: `{"name": "pwd", "script": "pwd"}`,
			},
			expected: expected{
				create: expCreate{
					want: true,
					cmd: &entities.Command{
						ID:     1,
						Name:   "pwd",
						Script: "pwd",
					},
					err: nil,
				},
				append: expAppend{
					want: true,
					err:  nil,
				},
				delete: expDelete{
					want: true,
					err:  nil,
				},
			},
			wantCode: wantCode{
				create: http.StatusCreated,
				delete: http.StatusNoContent,
			},
		},
		{
			name: "incorrect_query_values_count",
			args: args{
				query: query{
					want:   true,
					key:    "name",
					values: []string{"pwd", "second"},
				},
			},
			expected: expected{},
			wantCode: wantCode{
				delete: http.StatusBadRequest,
			},
		},
		{
			name:     "no_query",
			args:     args{},
			expected: expected{},
			wantCode: wantCode{
				delete: http.StatusBadRequest,
			},
		},
		{
			name: "incorrect_query_key",
			args: args{
				query: query{
					want:   true,
					key:    "unknown",
					values: []string{"pwd"},
				},
			},
			expected: expected{},
			wantCode: wantCode{
				delete: http.StatusBadRequest,
			},
		},
		{
			name: "command_not_found",
			args: args{
				query: query{
					want:   true,
					key:    "name",
					values: []string{"uknown"},
				},
			},
			expected: expected{
				delete: expDelete{
					want: true,
					err:  errs.ErrCmdNotFound,
				},
			},
			wantCode: wantCode{
				delete: http.StatusNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mocks expected response
			if tt.expected.create.want {
				mockRepo.EXPECT().CreateCommand(gomock.Any(), gomock.Any()).
					Return(tt.expected.create.cmd, tt.expected.create.err).Times(1)
			}
			if tt.expected.append.want {
				mockRepo.EXPECT().AppendCommandOutputByName(gomock.Any(), gomock.Any()).
					Return(tt.expected.append.err).Times(1)
			}
			if tt.expected.delete.want {
				mockRepo.EXPECT().DeleteCommandByName(gomock.Any(), gomock.Any()).
					Return(tt.expected.delete.err).Times(1)
			}

			// Controller
			ctrl := handlers.NewController(ctx, cfg)
			ch := make(chan entities.Command)
			mh, err := ctrl.BuildRoute(ctx, mockRepo, ch)
			require.NoError(t, err)

			// CREATE COMMAND

			// Form new request
			url := `http://` + cfg.Address + `/command`

			if tt.wantCode.create != 0 {
				r := httptest.NewRequest(http.MethodPost, url, bytes.NewBufferString(tt.args.reqBody))
				w := httptest.NewRecorder()

				mh.ServeHTTP(w, r)
				time.Sleep(100 * time.Millisecond)

				// Get response
				resp := w.Result()
				defer resp.Body.Close()

				// Check status code
				require.Equal(t, tt.wantCode.create, resp.StatusCode)
			}

			// DELETE COMMAND

			// Form new request
			r := httptest.NewRequest(http.MethodDelete, url, nil)
			if tt.args.query.want {
				q := r.URL.Query()
				for _, v := range tt.args.query.values {
					q.Add(tt.args.query.key, v)
				}
				r.URL.RawQuery = q.Encode()
			}
			w := httptest.NewRecorder()

			mh.ServeHTTP(w, r)

			// Get response
			resp := w.Result()
			defer resp.Body.Close()

			// Check status code
			require.Equal(t, tt.wantCode.delete, resp.StatusCode)
		})
	}
}
