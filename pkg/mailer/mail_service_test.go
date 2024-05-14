package mailer

import (
	"context"
	"errors"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/mailer/mock"
	"go.uber.org/mock/gomock"
)

func init() {
	if _, err := os.Stat("./../../.env"); !os.IsNotExist(err) {
		err := godotenv.Load(os.ExpandEnv("./../../.env"))
		if err != nil {
			log.Fatalf("Error getting env %v\n", err)
		}
	}
}

func TestNewMailService(t *testing.T) {
	tests := []struct {
		name string
		want *mailService
	}{
		{
			name: "implemented",
			want: NewMailService(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMailService(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMailService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mailer_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	mailer := mock.NewMockMailerInterface(ctrl)

	type args struct {
		ctx          context.Context
		recipient    string
		templateFile string
		data         interface{}
	}
	tests := []struct {
		name    string
		m       *mailService
		args    args
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "error template file",
			m: &mailService{
				sender: NewMailService().sender,
				mailer: mailer,
			},
			args: args{
				ctx:          context.Background(),
				recipient:    "allendragneel@gmail.com",
				templateFile: "notifmail.tmpl",
				data: map[string]interface{}{
					"Sighter": "nanda",
					"TigerID": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "error template file with no subject",
			m: &mailService{
				sender: NewMailService().sender,
				mailer: mailer,
			},
			args: args{
				ctx:          context.Background(),
				recipient:    "allendragneel@gmail.com",
				templateFile: "test_nosubject.tmpl",
				data: map[string]interface{}{
					"Sighter": "nanda",
					"TigerID": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "error template file with no html body",
			m: &mailService{
				sender: NewMailService().sender,
				mailer: mailer,
			},
			args: args{
				ctx:          context.Background(),
				recipient:    "allendragneel@gmail.com",
				templateFile: "test_nohtmlbody.tmpl",
				data: map[string]interface{}{
					"Sighter": "nanda",
					"TigerID": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "error template file with no plain body",
			m: &mailService{
				sender: NewMailService().sender,
				mailer: mailer,
			},
			args: args{
				ctx:          context.Background(),
				recipient:    "allendragneel@gmail.com",
				templateFile: "test_noplainbody.tmpl",
				data: map[string]interface{}{
					"Sighter": "nanda",
					"TigerID": "2",
				},
			},
			wantErr: true,
		},
		{
			name: "error dial and send mailer helper",
			m: &mailService{
				sender: NewMailService().sender,
				mailer: mailer,
			},
			args: args{
				ctx:          context.Background(),
				recipient:    "allendragneel@gmail.com",
				templateFile: "notif_mail.tmpl",
				data: map[string]interface{}{
					"Sighter": "nanda",
					"TigerID": "2",
				},
			},
			wantErr: true,
			mocks: []*gomock.Call{
				mailer.EXPECT().DialAndSend(gomock.Any()).Return(errors.New("any error")),
			},
		},
		{
			name: "success",
			m: &mailService{
				sender: NewMailService().sender,
				mailer: mailer,
			},
			args: args{
				ctx:          context.Background(),
				recipient:    "allendragneel@gmail.com",
				templateFile: "notif_mail.tmpl",
				data: map[string]interface{}{
					"Sighter": "nanda",
					"TigerID": "2",
				},
			},
			wantErr: false,
			mocks: []*gomock.Call{
				mailer.EXPECT().DialAndSend(gomock.Any()).Return(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Send(tt.args.ctx, tt.args.recipient, tt.args.templateFile, tt.args.data); err != nil {
				if !tt.wantErr {
					t.Errorf("mailer.Send() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
