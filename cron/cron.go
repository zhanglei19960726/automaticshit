package cron

import (
	"automaticshit/automaticshit"
	"automaticshit/common/config"
	"automaticshit/common/context"
	"automaticshit/notic"

	cron "github.com/robfig/cron/v3"
)

type Cron struct {
	spec    string
	cro     *cron.Cron
	entryID cron.EntryID
	exeFunc func()
}

func NewCron(ctx context.IContext, spec string, notic notic.INotic, autoMatci automaticshit.IAutoMaticShit) (c *Cron, err error) {
	c = &Cron{
		spec: spec,
		cro:  cron.New(cron.WithLogger(newCronLog(ctx)), cron.WithSeconds()),
		exeFunc: func() {
			if err := notic.NoticShit(ctx, autoMatci); err != nil {
				ctx.Error("NotciShit error", err.Error())
			}
		},
	}
	c.entryID, err = c.cro.AddFunc(spec, c.exeFunc)
	if err != nil {
		return
	}
	c.cro.Start()
	return
}

func (c *Cron) ReloadConfig(ctx context.IContext) {
	ctx.Debug("reloadConfig")
	cfg := config.GetConfig()
	if cfg.CronConfig.Space == c.spec {
		return
	}
	ctx.Debug("remove func")
	c.cro.Remove(c.entryID)
	id, err := c.cro.AddFunc(cfg.CronConfig.Space, c.exeFunc)
	if err != nil {
		ctx.Error("AddFunc error", err.Error())
		return
	}
	c.entryID = id
}
