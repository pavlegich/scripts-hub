package command

import (
	"context"
	"reflect"
	"testing"
)

func TestNewCommandService(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx  context.Context
		repo Repository
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
