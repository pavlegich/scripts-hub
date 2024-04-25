package command

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pavlegich/scripts-hub/internal/entities"
	"github.com/pavlegich/scripts-hub/internal/mocks"
	repo "github.com/pavlegich/scripts-hub/internal/repository"
)

func TestNewCommandService(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
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
				ctx:  ctx,
				repo: nil,
			},
			want: &CommandService{
				repo: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCommandService(tt.args.ctx, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandService_Create(t *testing.T) {
	ctx := context.Background()
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	s := NewCommandService(ctx, mockRepo)
	_ = s

	type args struct {
		command *entities.Command
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// got, err := s.Create(tt.args.ctx, tt.args.command)
			// if (err != nil) != tt.wantErr {
			// 	t.Errorf("CommandService.Create() error = %v, wantErr %v", err, tt.wantErr)
			// 	return
			// }
			// if got != tt.want {
			// 	t.Errorf("CommandService.Create() = %v, want %v", got, tt.want)
			// }
		})
	}
}
