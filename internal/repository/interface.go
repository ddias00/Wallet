package repository

import (
	"github.com/google/uuid"
	"github.com/wallet/internal/model"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	CreateUser(user *model.User) error
	FindUserByEmail(email string) (*model.User, error)
	FindUserByCPF(cpf string) (*model.User, error)
	FindUserByID(userID uuid.UUID) (*model.User, error)
	CreateWallet(wallet *model.Wallet) error
	FindWalletByUserID(userID uuid.UUID) (*model.Wallet, error)
	UpdateWalletBalance(tx *gorm.DB, walletID uuid.UUID, amount float64, isCredit bool) error
	CreateTransaction(tx *gorm.DB, transaction *model.Transaction) error
	FindTransactionsByUserID(userID uuid.UUID, limit, offset int) ([]model.Transaction, error)
	UpdateTransactionStatus(tx *gorm.DB, id uuid.UUID, status string) error
	FindPixKeyByKey(key string) (*model.PixKey, error)
	CreateBoleto(boleto *model.Boleto) error
	CreateRefreshToken(token *model.RefreshToken) error
	FindRefreshToken(token string) (*model.RefreshToken, error)
	DeleteRefreshTokensByUserID(userID uuid.UUID) error
	Transaction(fn func(*gorm.DB) error) error
}
