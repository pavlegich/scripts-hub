package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/pavlegich/scripts-hub/internal/infra/config"
)

func TestNewController(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
		want *Controller
	}{
		{
			name: "ok",
			args: args{
				cfg: cfg,
			},
			want: &Controller{
				cfg: cfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewController(ctx, tt.args.cfg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}
