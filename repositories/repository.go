package repository

import (
	// model "github.com/tientp-floware/mgodb-stream/models"
	setting "github.com/tientp-floware/mgodb-stream/repositories/setting"
)

type (
	// Service instance
	Service struct {
		Setting Setting
	}
	// Setting list function can use
	Setting interface {
	}
)

// New creates new user application service
func New() *Service {
	return &Service{Setting: setting.NewSetting()}
}
