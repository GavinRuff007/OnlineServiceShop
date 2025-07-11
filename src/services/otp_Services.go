package services

import (
	"RestGoTest/src/cache"
	"RestGoTest/src/config"
	"RestGoTest/src/constant"
	"RestGoTest/src/helper/service_errors"
	"RestGoTest/src/pkg/logging"
	"RestGoTest/src/util"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type OtpUsecase struct {
	logger      logging.Logger
	cfg         *config.Config
	redisClient *redis.Client
}

type otpDto struct {
	Value string
	Used  bool
}

func NewOtpUsecase(cfg *config.Config) *OtpUsecase {
	logger := logging.NewLogger(cfg)
	redis := cache.GetRedis()
	return &OtpUsecase{logger: logger, cfg: cfg, redisClient: redis}
}

func (u *OtpUsecase) SendOtp(mobileNumber string) error {
	otp := util.GenerateOtp()
	err := u.SetOtp(mobileNumber, otp)
	if err != nil {
		return err
	}
	return nil
}

func (u *OtpUsecase) SetOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constant.RedisOtpDefaultKey, mobileNumber)
	val := &otpDto{
		Value: otp,
		Used:  false,
	}

	res, err := cache.Get[otpDto](u.redisClient, key)
	if err == nil && !res.Used {
		return &service_errors.ServiceError{EndUserMessage: constant.OptExists}
	} else if err == nil && res.Used {
		return &service_errors.ServiceError{EndUserMessage: constant.OtpUsed}
	}
	err = cache.Set(u.redisClient, key, val, u.cfg.Otp.ExpireTime*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (u *OtpUsecase) ValidateOtp(mobileNumber string, otp string) error {
	key := fmt.Sprintf("%s:%s", constant.RedisOtpDefaultKey, mobileNumber)
	res, err := cache.Get[otpDto](u.redisClient, key)
	if err != nil {
		return err
	} else if res.Used {
		return &service_errors.ServiceError{EndUserMessage: constant.OtpUsed}
	} else if !res.Used && res.Value != otp {
		return &service_errors.ServiceError{EndUserMessage: constant.OtpNotValid}
	} else if !res.Used && res.Value == otp {
		res.Used = true
		err = cache.Set(u.redisClient, key, res, u.cfg.Otp.ExpireTime*time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}
