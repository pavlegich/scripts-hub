package command

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pavlegich/scripts-hub/internal/entities"
	errs "github.com/pavlegich/scripts-hub/internal/errors"
	"github.com/pavlegich/scripts-hub/internal/mocks"
	repo "github.com/pavlegich/scripts-hub/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestNewCommandService(t *testing.T) {
	ctx := context.Background()

	type args struct {
		repo repo.Repository
	}
	tests := []struct {
		name string
		args args
		want *CommandService
	}{
		{
			name: "ok",
			args: args{
				repo: nil,
			},
			want: &CommandService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCommandService(ctx, tt.args.repo)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommandService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandService_Create(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	s := NewCommandService(ctx, mockRepo)

	type expected struct {
		cmd *entities.Command
		err error
	}
	type args struct {
		command *entities.Command
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		wantErr  error
		want     int
	}{
		{
			name: "success",
			args: args{
				command: &entities.Command{
					Name:   "ok",
					Script: "pwd",
				},
			},
			expected: expected{
				cmd: &entities.Command{
					ID:     1,
					Name:   "ok",
					Script: "pwd",
				},
				err: nil,
			},
			wantErr: nil,
			want:    1,
		},
		{
			name: "command_already_exists",
			args: args{
				command: &entities.Command{
					Name:   "ok",
					Script: "pwd",
				},
			},
			expected: expected{
				cmd: nil,
				err: errs.ErrCmdAlreadyExists,
			},
			wantErr: errs.ErrCmdAlreadyExists,
			want:    -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().CreateCommand(gomock.Any(), gomock.Any()).
				Return(tt.expected.cmd, tt.expected.err).Times(1)

			got, err := s.Create(ctx, tt.args.command)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCommandService_List(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	s := NewCommandService(ctx, mockRepo)

	type expected struct {
		cmds []*entities.Command
		err  error
	}
	tests := []struct {
		name     string
		expected expected
		wantErr  error
		want     []*entities.Command
	}{
		{
			name: "success",
			expected: expected{
				cmds: []*entities.Command{
					{
						ID:     1,
						Name:   "ok",
						Script: "pwd",
					},
				},
				err: nil,
			},
			wantErr: nil,
			want: []*entities.Command{
				{
					ID:     1,
					Name:   "ok",
					Script: "pwd",
				},
			},
		},
		{
			name: "no_data_in_db",
			expected: expected{
				cmds: nil,
				err:  errs.ErrCmdNotFound,
			},
			wantErr: errs.ErrCmdNotFound,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetAllCommands(gomock.Any()).
				Return(tt.expected.cmds, tt.expected.err).Times(1)

			got, err := s.List(ctx)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCommandService_Unload(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	s := NewCommandService(ctx, mockRepo)

	type expected struct {
		cmd *entities.Command
		err error
	}
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		wantErr  error
		want     *entities.Command
	}{
		{
			name: "success",
			args: args{
				name: "ok",
			},
			expected: expected{
				cmd: &entities.Command{
					ID:     1,
					Name:   "ok",
					Script: "pwd",
				},
				err: nil,
			},
			wantErr: nil,
			want: &entities.Command{
				ID:     1,
				Name:   "ok",
				Script: "pwd",
			},
		},
		{
			name: "no_data_in_db",
			args: args{
				name: "nothing",
			},
			expected: expected{
				cmd: nil,
				err: errs.ErrCmdNotFound,
			},
			wantErr: errs.ErrCmdNotFound,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().GetCommandByName(gomock.Any(), gomock.Any()).
				Return(tt.expected.cmd, tt.expected.err).Times(1)

			got, err := s.Unload(ctx, tt.args.name)

			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestCommandService_Delete(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	s := NewCommandService(ctx, mockRepo)

	type expected struct {
		err error
	}
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				name: "ok",
			},
			expected: expected{
				err: nil,
			},
			wantErr: nil,
		},
		{
			name: "no_data_in_db",
			args: args{
				name: "nothing",
			},
			expected: expected{
				err: errs.ErrCmdNotFound,
			},
			wantErr: errs.ErrCmdNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().DeleteCommandByName(gomock.Any(), gomock.Any()).
				Return(tt.expected.err).Times(1)

			err := s.Delete(ctx, tt.args.name)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommandService_Update(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	s := NewCommandService(ctx, mockRepo)

	type expected struct {
		err error
	}
	type args struct {
		command *entities.Command
	}
	tests := []struct {
		name     string
		args     args
		expected expected
		wantErr  error
	}{
		{
			name: "success",
			args: args{
				command: &entities.Command{
					Name:   "ok",
					Script: "pwd",
					Output: "/path",
				},
			},
			expected: expected{
				err: nil,
			},
			wantErr: nil,
		},
		{
			name: "no_data_in_db",
			args: args{
				command: &entities.Command{
					Name:   "nothing",
					Script: "pwd",
					Output: "/nopath",
				},
			},
			expected: expected{
				err: errs.ErrCmdNotFound,
			},
			wantErr: errs.ErrCmdNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().AppendCommandOutputByName(gomock.Any(), gomock.Any()).
				Return(tt.expected.err).Times(1)

			err := s.AppendOutput(ctx, tt.args.command)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
