package core

import "log"

func (c *Core) Debugf(template string, args ...any) {
	if c.logger != nil {
		c.logger.Sugar().Debugf(template, args...)
	} else {
		log.Printf(template, args...)
	}
}

func (c *Core) Infof(template string, args ...any) {
	if c.logger != nil {
		c.logger.Sugar().Infof(template, args...)
	} else {
		log.Printf(template, args...)
	}
}

func (c *Core) Errorf(template string, args ...any) {
	if c.logger != nil {
		c.logger.Sugar().Errorf(template, args...)
	} else {
		log.Printf(template, args...)
	}
}

func (c *Core) Debug(args ...any) {
	if c.logger != nil {
		c.logger.Sugar().Debug(args)
	} else {
		log.Print(args...)
	}
}

func (c *Core) Info(args ...any) {
	if c.logger != nil {
		c.logger.Sugar().Info(args)
	} else {
		log.Print(args...)
	}
}

func (c *Core) Error(args ...any) {
	if c.logger != nil {
		c.logger.Sugar().Error(args)
	} else {
		log.Print(args...)
	}
}
