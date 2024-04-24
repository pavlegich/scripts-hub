package handlers

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/pavlegich/scripts-hub/internal/infra/config"
)

func TestNewController(t *testing.T) {
	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	type args struct {
		ctx context.Context
		db  *sql.DB
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
				ctx: ctx,
				db:  nil,
				cfg: cfg,
			},
			want: &Controller{
				db:  nil,
				cfg: cfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewController(tt.args.ctx, tt.args.db, tt.args.cfg)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}
