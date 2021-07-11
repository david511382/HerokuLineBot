package bootstrap

import (
	"heroku-line-bot/util"
	"testing"
)

func TestDb_ScanUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		db      *Db
		args    args
		wantDb  *Db
		wantErr bool
	}{
		{
			"full",
			&Db{},
			args{
				url: "protocol://user:password@server:80/database?key=value",
			},
			&Db{
				Server: Server{
					Host: "server",
					Port: 80,
				},
				User:     "user",
				Password: "password",
				Database: "database",
				Param:    "key=value",
				Protocol: "protocol",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.db.ScanUrl(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("Db.ScanUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if ok, msg := util.Comp(tt.db, tt.wantDb); !ok {
				t.Fatal(msg)
			}
		})
	}
}
