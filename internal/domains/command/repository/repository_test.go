// Package repository contains repository object
// and methods for interaction between service and storage.
package repository

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
)

func TestNewCommandRepository(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx context.Context
		db  *sql.DB
	}
	tests := []struct {
		name string
		args args
		want *Repository
	}{
		{
			name: "ok",
			args: args{
				ctx: ctx,
				db:  nil,
			},
			want: &Repository{
				db: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCommandRepository(tt.args.ctx, tt.args.db)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDataRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}
