package config

import (
	"music-sharing/user-microservice/internal/infrastructure"
	"sync"

	"gorm.io/gorm"
)

var (
	Container *DIContainer
	initOnce  sync.Once
)

type (
	DIContainer struct {
		CommmandBus *infrastructure.CommandBus
		QueryBus    *infrastructure.QueryBus
		Database    *gorm.DB
	}
)

func Initialize() error {

	initOnce.Do(func() {
		Container = &DIContainer{}

		Container.CommmandBus = infrastructure.GetCommandBus()
		Container.QueryBus = infrastructure.GetQueryBus()
		Container.Database = infrastructure.GetDB()
	})

	return nil
}
