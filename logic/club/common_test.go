package club

import (
	"heroku-line-bot/util"
	"testing"
)

func Test_calculateActivityPay(t *testing.T) {
	type args struct {
		people      int
		ballConsume util.Float
		courtFee    util.Float
		clubSubsidy util.Float
	}
	tests := []struct {
		name              string
		args              args
		wantActivityFee   util.Float
		wantClubMemberFee int
		wantGuestFee      int
	}{
		{
			"overflow",
			args{
				people:      14,
				ballConsume: util.NewFloat(16),
				courtFee:    util.NewFloat(1920),
				clubSubsidy: util.NewFloat(0),
			},
			util.NewFloat(2416),
			175,
			175,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotActivityFee, gotClubMemberFee, gotGuestFee := calculateActivityPay(tt.args.people, tt.args.ballConsume, tt.args.courtFee, tt.args.clubSubsidy)
			if ok, msg := util.Comp(gotActivityFee.Value(), tt.wantActivityFee.Value()); !ok {
				t.Errorf(msg)
				return
			}
			if ok, msg := util.Comp(gotClubMemberFee, tt.wantClubMemberFee); !ok {
				t.Errorf(msg)
				return
			}
			if ok, msg := util.Comp(gotGuestFee, tt.wantGuestFee); !ok {
				t.Errorf(msg)
				return
			}
		})
	}
}
