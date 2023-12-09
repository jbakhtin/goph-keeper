package appservices

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/output/logger/v1"

	"github.com/jbakhtin/goph-keeper/internal/server/implements/adapters/input/config/v1/drivers"
	configInterface "github.com/jbakhtin/goph-keeper/internal/server/interfaces/ports/input/config/v1"
)

func TestNewPasswordAppService(t *testing.T) {
	type args struct {
		cfg configInterface.Interface
		lgr logger.Interface
	}
	tests := []struct {
		name    string
		args    args
		want    *PasswordAppService
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPasswordAppService(tt.args.cfg, tt.args.lgr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPasswordAppService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPasswordAppService() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordAppService_CheckPassword(t *testing.T) {
	cfg, err := drivers.NewFormENV()
	if err != nil {
		log.Fatal(err)
	}

	h := hmac.New(sha256.New, []byte(cfg.GetAppKey()))
	h.Write([]byte("password"))
	dst := h.Sum(nil)
	hashedPassword := fmt.Sprintf("%x", dst)

	type fields struct {
		cfg configInterface.Interface
	}
	type args struct {
		password string
		need     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Unsuccessful",
			fields: fields{
				cfg: cfg,
			},
			args: args{
				password: "password_1",
				need:     hashedPassword,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Successful",
			fields: fields{
				cfg: cfg,
			},
			args: args{
				password: "password",
				need:     hashedPassword,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PasswordAppService{
				cfg: tt.fields.cfg,
			}
			got, err := p.CheckPassword(tt.args.password, tt.args.need)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswordAppService_HashPassword(t *testing.T) {
	type fields struct {
		cfg configInterface.Interface
	}
	type args struct {
		rawPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PasswordAppService{
				cfg: tt.fields.cfg,
			}
			got, err := p.HashPassword(tt.args.rawPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HashPassword() got = %v, want %v", got, tt.want)
			}
		})
	}
}
