package repository

import (
	"time"

	"github.com/google/uuid"
	"github.com/wallet/internal/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// User
func (r *Repository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *Repository) FindUserByCPF(cpf string) (*model.User, error) {
	var user model.User
	err := r.db.Where("cpf = ?", cpf).First(&user).Error
	return &user, err
}

func (r *Repository) FindUserByID(userID uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	return &user, err
}

// Wallet
func (r *Repository) CreateWallet(wallet *model.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *Repository) FindWalletByUserID(userID uuid.UUID) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.Where("user_id = ?", userID).First(&wallet).Error
	return &wallet, err
}

func (r *Repository) UpdateWalletBalance(tx *gorm.DB, walletID uuid.UUID, amount float64, isCredit bool) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	if isCredit {
		return db.Model(&model.Wallet{}).Where("id = ?", walletID).UpdateColumn("balance", gorm.Expr("balance + ?", amount)).Error
	}
	return db.Model(&model.Wallet{}).Where("id = ? AND balance >= ?", walletID, amount).UpdateColumn("balance", gorm.Expr("balance - ?", amount)).Error
}

// Transaction
func (r *Repository) CreateTransaction(tx *gorm.DB, transaction *model.Transaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(transaction).Error
}

func (r *Repository) FindTransactionsByUserID(userID uuid.UUID, limit, offset int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("from_user_id = ? OR to_user_id = ?", userID, userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error
	return transactions, err
}

func (r *Repository) UpdateTransactionStatus(tx *gorm.DB, id uuid.UUID, status string) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Model(&model.Transaction{}).Where("id = ?", id).Update("status", status).Error
}

// PixKey
func (r *Repository) FindPixKeyByKey(key string) (*model.PixKey, error) {
	var pixKey model.PixKey
	err := r.db.Where("key = ? AND status = ?", key, "active").First(&pixKey).Error
	return &pixKey, err
}

// Boleto
func (r *Repository) CreateBoleto(boleto *model.Boleto) error {
	return r.db.Create(boleto).Error
}

// RefreshToken
func (r *Repository) CreateRefreshToken(token *model.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *Repository) FindRefreshToken(token string) (*model.RefreshToken, error) {
	var refreshToken model.RefreshToken
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).First(&refreshToken).Error
	return &refreshToken, err
}

func (r *Repository) DeleteRefreshTokensByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}

// Transaction helper
func (r *Repository) Transaction(fn func(*gorm.DB) error) error {
	return r.db.Transaction(fn)
}
