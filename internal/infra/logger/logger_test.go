// Package logger contains objects and methods for logging the events.
package logger

import (
	"context"
	"testing"
)

func TestInit(t *testing.T) {
	ctx := context.Background()

	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				level: "Info",
			},
			wantErr: false,
		},
		{
			name: "wrong_level",
			args: args{
				level: "Wrong",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Init(ctx, tt.args.level)

			if (err != nil) != tt.wantErr {
				t.Errorf("Init() = %v, want %v", err != nil, tt.wantErr)
			}
		})
	}
}
