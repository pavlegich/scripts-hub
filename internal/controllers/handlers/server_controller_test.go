package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/pavlegich/scripts-hub/internal/infra/config"
	"github.com/pavlegich/scripts-hub/internal/repository"
)

func TestNewController(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	type args struct {
		ctx  context.Context
		repo repository.Repository
		cfg  *config.Config
	}
	tests := []struct {
		name string
		args args
		want *Controller
	}{
		{
			name: "ok",
			args: args{
				ctx:  ctx,
				repo: nil,
				cfg:  cfg,
			},
			want: &Controller{
				repo: nil,
				cfg:  cfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewController(tt.args.ctx, tt.args.repo, tt.args.cfg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}
