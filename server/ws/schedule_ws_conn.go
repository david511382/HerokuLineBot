package ws

import (
	"heroku-line-bot/global"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type IWsSender interface {
	Send(
		messageType int, p []byte,
	) error
	Close() error
}

type IScheduleWsConnJob interface {
	Run(IWsSender)
	Listen(IWsSender, *WsConnReadMessage)
}

type cronJob struct {
	wsSender IWsSender
	job      IScheduleWsConnJob
	spec     string
}

func newCronJob(spec string, wsSender IWsSender, job IScheduleWsConnJob) *cronJob {
	return &cronJob{
		job:      job,
		wsSender: wsSender,
		spec:     spec,
	}
}

func (j *cronJob) Run() {
	j.job.Run(j.wsSender)
}

type ScheduleWsConn struct {
	conn *WsConn
	cr   *cron.Cron
	jobs []*cronJob
}

func NewScheduleWsConn(c *gin.Context) (r *ScheduleWsConn, resultErr error) {
	conn, err := NewWsConn(c)
	if err != nil {
		resultErr = err
		return
	}

	r = &ScheduleWsConn{
		conn: conn,
		cr:   cron.NewWithLocation(global.Location),
		jobs: make([]*cronJob, 0),
	}
	conn.SetCloseListener(r.cr.Stop)
	r.conn.SetMessageListener(func(wcrm *WsConnReadMessage) {
		for _, cronJob := range r.jobs {
			cronJob.job.Listen(r.conn, wcrm)
		}
	})
	return
}

// spec: "0 */29 * * * *"
// Seconds
// Minutes
// Hours
// Day-of-Month
// Month
// Day-of-Week
// Year (optional field)
func (w *ScheduleWsConn) AddJob(spec string, job IScheduleWsConnJob) {
	cronJob := newCronJob(spec, w.conn, job)
	w.jobs = append(w.jobs, cronJob)
	return
}

func (w *ScheduleWsConn) Serve() (resultErr error) {
	for _, cronJob := range w.jobs {
		if err := w.cr.AddJob(cronJob.spec, cronJob); err != nil {
			resultErr = err
			return
		}
	}

	w.conn.Serve()

	for _, cronJob := range w.jobs {
		cronJob.Run()
	}

	w.cr.Start()

	return
}

func (w *ScheduleWsConn) SetListenHeartBeatTimeout(timeout time.Duration) {
	w.conn.SetListenHeartBeatTimeout(timeout)
}
