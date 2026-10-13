package auth

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Sabbir185/gpi/entity"
)

type AuthRepository interface {
	CreateOTP(identity string, otp string, OtpType entity.OtpType, OtpChannel entity.OtpChannel, ttl time.Duration) error
	FindOTPByIdentityAndType(identity string, OtpType entity.OtpType) (*entity.Otp, error)
}

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepository{db: db}
}

func (s *authRepository) CreateOTP(identity string, otp string, OtpType entity.OtpType, OtpChannel entity.OtpChannel, ttl time.Duration) error {
	fmt.Println("identity", identity)
	fmt.Println("otp", otp)
	fmt.Println("OtpType", OtpType)
	fmt.Println("OtpChannel", OtpChannel)
	fmt.Println("ttl", ttl)
	return nil
}

func (s *authRepository) FindOTPByIdentityAndType(identity string, OtpType entity.OtpType) (*entity.Otp, error) {
	return nil, nil
}
