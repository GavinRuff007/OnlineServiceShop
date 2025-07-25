package dependency

import (
	"RestGoTest/src/config"
	contractRepository "RestGoTest/src/domain/repository"
	infraRepository "RestGoTest/src/infra/repository"
)

func GetUserRepository(cfg *config.Config) contractRepository.UserRepository {
	return infraRepository.NewUserRepository(cfg)
}

func GetOrderRepository(cfg *config.Config) contractRepository.OrderRepository {
	return infraRepository.NewOrderRepository(cfg)
}
