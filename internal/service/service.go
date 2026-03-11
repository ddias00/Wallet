package service

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/wallet/internal/model"
	"github.com/wallet/internal/repository"
	"github.com/wallet/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo       repository.RepositoryInterface
	jwtSecret  string
	accessExp  time.Duration
	refreshExp time.Duration
}

func New(repo repository.RepositoryInterface, jwtSecret string, accessExp, refreshExp time.Duration) *Service {
	return &Service{
		repo:       repo,
		jwtSecret:  jwtSecret,
		accessExp:  accessExp,
		refreshExp: refreshExp,
	}
}

// Auth
func (s *Service) Register(name, cpf, email, phone, password string) (*model.User, error) {
	if !validateCPF(cpf) {
		return nil, errors.New("invalid CPF")
	}
	if !validateEmail(email) {
		return nil, errors.New("invalid email")
	}
	if !validatePhone(phone) {
		return nil, errors.New("invalid phone")
	}

	if _, err := s.repo.FindUserByEmail(email); err == nil {
		return nil, errors.New("email already exists")
	}
	if _, err := s.repo.FindUserByCPF(cpf); err == nil {
		return nil, errors.New("CPF already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     name,
		CPF:      cpf,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	wallet := &model.Wallet{UserID: user.ID, Balance: 0}
	s.repo.CreateWallet(wallet)

	return user, nil
}

func (s *Service) Login(email, password string) (string, string, error) {
	user, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if user.IsBlocked {
		return "", "", errors.New("account is blocked")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := auth.GenerateTokenPair(user.ID, s.jwtSecret, s.accessExp, s.refreshExp)
	if err != nil {
		return "", "", err
	}

	s.repo.DeleteRefreshTokensByUserID(user.ID)
	s.repo.CreateRefreshToken(&model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.refreshExp),
	})

	return accessToken, refreshToken, nil
}

func (s *Service) RefreshToken(refreshToken string) (string, string, error) {
	storedToken, err := s.repo.FindRefreshToken(refreshToken)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	claims, err := auth.ValidateToken(refreshToken, s.jwtSecret)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	accessToken, newRefreshToken, err := auth.GenerateTokenPair(claims.UserID, s.jwtSecret, s.accessExp, s.refreshExp)
	if err != nil {
		return "", "", err
	}

	s.repo.DeleteRefreshTokensByUserID(storedToken.UserID)
	s.repo.CreateRefreshToken(&model.RefreshToken{
		UserID:    storedToken.UserID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(s.refreshExp),
	})

	return accessToken, newRefreshToken, nil
}

// Wallet
func (s *Service) GetBalance(userID uuid.UUID) (*model.Wallet, error) {
	return s.repo.FindWalletByUserID(userID)
}

func (s *Service) GetUserByID(userID uuid.UUID) (*model.User, error) {
	return s.repo.FindUserByID(userID)
}

func (s *Service) CreditWallet(userID uuid.UUID, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	wallet, err := s.repo.FindWalletByUserID(userID)
	if err != nil {
		return errors.New("wallet not found")
	}

	return s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.UpdateWalletBalance(tx, wallet.ID, amount, true); err != nil {
			return err
		}

		transaction := &model.Transaction{
			ToUserID:    &userID,
			Amount:      amount,
			Type:        "credit",
			Status:      "confirmed",
			Description: description,
		}
		return s.repo.CreateTransaction(tx, transaction)
	})
}

func (s *Service) DebitWallet(userID uuid.UUID, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	wallet, err := s.repo.FindWalletByUserID(userID)
	if err != nil {
		return errors.New("wallet not found")
	}

	if wallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	return s.repo.Transaction(func(tx *gorm.DB) error {
		if err := s.repo.UpdateWalletBalance(tx, wallet.ID, amount, false); err != nil {
			return err
		}

		transaction := &model.Transaction{
			FromUserID:  &userID,
			Amount:      amount,
			Type:        "debit",
			Status:      "confirmed",
			Description: description,
		}
		return s.repo.CreateTransaction(tx, transaction)
	})
}

// Transaction
func (s *Service) Transfer(fromUserID, toUserID uuid.UUID, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if fromUserID == toUserID {
		return errors.New("cannot transfer to yourself")
	}

	fromWallet, err := s.repo.FindWalletByUserID(fromUserID)
	if err != nil {
		return errors.New("sender wallet not found")
	}
	if fromWallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	toWallet, err := s.repo.FindWalletByUserID(toUserID)
	if err != nil {
		return errors.New("recipient wallet not found")
	}

	return s.repo.Transaction(func(tx *gorm.DB) error {
		transaction := &model.Transaction{
			FromUserID:  &fromUserID,
			ToUserID:    &toUserID,
			Amount:      amount,
			Type:        "transfer",
			Status:      "pending",
			Description: description,
		}

		if err := s.repo.CreateTransaction(tx, transaction); err != nil {
			return err
		}

		if err := s.repo.UpdateWalletBalance(tx, fromWallet.ID, amount, false); err != nil {
			return errors.New("insufficient balance")
		}

		if err := s.repo.UpdateWalletBalance(tx, toWallet.ID, amount, true); err != nil {
			return err
		}

		return s.repo.UpdateTransactionStatus(tx, transaction.ID, "confirmed")
	})
}

func (s *Service) GetTransactions(userID uuid.UUID, limit, offset int) ([]model.Transaction, error) {
	return s.repo.FindTransactionsByUserID(userID, limit, offset)
}

// PIX
func (s *Service) CreatePixPayment(fromUserID uuid.UUID, pixKey string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	recipientPixKey, err := s.repo.FindPixKeyByKey(pixKey)
	if err != nil {
		return errors.New("PIX key not found")
	}

	toUserID := recipientPixKey.UserID
	if fromUserID == toUserID {
		return errors.New("cannot transfer to yourself")
	}

	fromWallet, err := s.repo.FindWalletByUserID(fromUserID)
	if err != nil {
		return errors.New("sender wallet not found")
	}
	if fromWallet.Balance < amount {
		return errors.New("insufficient balance")
	}

	toWallet, err := s.repo.FindWalletByUserID(toUserID)
	if err != nil {
		return errors.New("recipient wallet not found")
	}

	return s.repo.Transaction(func(tx *gorm.DB) error {
		transaction := &model.Transaction{
			FromUserID:  &fromUserID,
			ToUserID:    &toUserID,
			Amount:      amount,
			Type:        "transfer",
			Status:      "pending",
			Description: fmt.Sprintf("PIX transfer to %s", pixKey),
		}

		if err := s.repo.CreateTransaction(tx, transaction); err != nil {
			return err
		}

		if err := s.repo.UpdateWalletBalance(tx, fromWallet.ID, amount, false); err != nil {
			return errors.New("insufficient balance")
		}

		if err := s.repo.UpdateWalletBalance(tx, toWallet.ID, amount, true); err != nil {
			return err
		}

		return s.repo.UpdateTransactionStatus(tx, transaction.ID, "confirmed")
	})
}

// Boleto
func (s *Service) GenerateBoleto(userID uuid.UUID, amount float64) (*model.Boleto, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be positive")
	}

	boleto := &model.Boleto{
		UserID:  userID,
		Amount:  amount,
		Code:    fmt.Sprintf("BOL%d%s", time.Now().Unix(), uuid.New().String()[:8]),
		Barcode: generateBarcode(),
		Status:  "pending",
		DueDate: time.Now().Add(3 * 24 * time.Hour),
	}

	if err := s.repo.CreateBoleto(boleto); err != nil {
		return nil, err
	}

	return boleto, nil
}

// Validators
func validateCPF(cpf string) bool {
	cpf = regexp.MustCompile(`[^0-9]`).ReplaceAllString(cpf, "")
	return len(cpf) == 11
}

func validateEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

func validatePhone(phone string) bool {
	phone = regexp.MustCompile(`[^0-9]`).ReplaceAllString(phone, "")
	return len(phone) >= 10 && len(phone) <= 11
}

func generateBarcode() string {
	return fmt.Sprintf("%d%s%s", time.Now().Unix(), uuid.New().String()[:20], "00000000000000000")[:47]
}
