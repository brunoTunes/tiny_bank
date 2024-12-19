package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestUser_validate(t *testing.T) {
	type fields struct {
		ID        string
		Name      string
		DeletedAt *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    *User
		wantErr error
	}{
		{
			name: "valid",
			fields: fields{
				ID:   "id",
				Name: "name",
			},
			want: &User{
				ID:   "id",
				Name: "name",
			},
		},
		{
			name: "invalid name, want invalidEmptyNameError",
			fields: fields{
				ID:   "id",
				Name: "",
			},
			wantErr: invalidEmptyNameError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:        tt.fields.ID,
				Name:      tt.fields.Name,
				DeletedAt: tt.fields.DeletedAt,
			}
			got, err := u.validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("validate() (-want +got):\n%s", diff)
			}
		})
	}
}
