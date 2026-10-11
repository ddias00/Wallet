package unit

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wallet/internal/model"
	"github.com/wallet/internal/repository"
	"github.com/wallet/internal/service"
	"github.com/wallet/pkg/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.Wallet{}, &model.Transaction{}, &model.PixKey{}, &model.Boleto{}, &model.RefreshToken{}); err != nil {
		t.Fatal(err)
	}
	return db
}

// Service Tests
func TestRepositoryRejectsInsufficientDebit(t *testing.T) {
	db := setupDB(t)
	repo := repository.New(db)
	wallet := &model.Wallet{UserID: uuid.New(), Balance: 5}
	if err := repo.CreateWallet(wallet); err != nil {
		t.Fatal(err)
	}
	if err := repo.UpdateWalletBalance(nil, wallet.ID, 10, false); err == nil {
		t.Fatal("debit must fail when the conditional update affects no rows")
	}
	stored, err := repo.FindWalletByUserID(wallet.UserID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Balance != 5 {
		t.Fatalf("rejected debit changed balance to %v", stored.Balance)
	}
}

func TestServiceRegister(t *testing.T) {
	db := setupDB(t)
	repo := repository.New(db)
	svc := service.New(repo, "secret", 0, 0)

	tests := []struct {
		name    string
		cpf     string
		email   string
		phone   string
		wantErr bool
	}{
		{"valid", "12345678909", "test@example.com", "11987654321", false},
		{"invalid cpf", "123", "test2@example.com", "11987654321", true},
		{"invalid email", "12345678910", "invalid", "11987654321", true},
		{"invalid phone", "12345678911", "test3@example.com", "123", true},
		{"duplicate email", "12345678912", "test@example.com", "11987654321", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Register("Test", tt.cpf, tt.email, tt.phone, "password123")
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceLogin(t *testing.T) {
	db := setupDB(t)
	repo := repository.New(db)
	svc := service.New(repo, "secret", 0, 0)

	svc.Register("Test", "12345678909", "test@example.com", "11987654321", "password123")

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{"valid", "test@example.com", "password123", false},
		{"wrong password", "test@example.com", "wrong", true},
		{"wrong email", "wrong@example.com", "password123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := svc.Login(tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceTransfer(t *testing.T) {
	db := setupDB(t)
	repo := repository.New(db)
	svc := service.New(repo, "secret", 0, 0)

	user1, _ := svc.Register("User1", "12345678909", "user1@example.com", "11987654321", "pass")
	user2, _ := svc.Register("User2", "98765432100", "user2@example.com", "11987654322", "pass")

	wallet1, _ := repo.FindWalletByUserID(user1.ID)
	repo.UpdateWalletBalance(nil, wallet1.ID, 100, true)

	tests := []struct {
		name    string
		from    uuid.UUID
		to      uuid.UUID
		amount  float64
		wantErr bool
	}{
		{"valid", user1.ID, user2.ID, 50, false},
		{"insufficient balance", user1.ID, user2.ID, 200, true},
		{"negative amount", user1.ID, user2.ID, -10, true},
		{"self transfer", user1.ID, user1.ID, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Transfer(tt.from, tt.to, tt.amount, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("Transfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceBoleto(t *testing.T) {
	db := setupDB(t)
	repo := repository.New(db)
	svc := service.New(repo, "secret", 0, 0)

	user, _ := svc.Register("Test", "12345678909", "test@example.com", "11987654321", "pass")

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"valid", 100, false},
		{"zero", 0, true},
		{"negative", -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GenerateBoleto(user.ID, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateBoleto() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Repository Tests
func TestRepositoryUser(t *testing.T) {
	db := setupDB(t)
	repo := repository.New(db)

	user := &model.User{Name: "Test", CPF: "12345678909", Email: "test@example.com", Phone: "11987654321", Password: "hashed"}

	if err := repo.CreateUser(user); err != nil {
		t.Errorf("CreateUser() error = %v", err)
	}

	found, err := repo.FindUserByEmail("test@example.com")
	if err != nil || found.Email != user.Email {
		t.Errorf("FindUserByEmail() failed")
	}
}

func TestRepositoryWallet(t *testing.T) {
	db := setupDB(t)
	repo := repository.New(db)

	user := &model.User{Name: "Test", CPF: "12345678909", Email: "test@example.com", Phone: "11987654321", Password: "hashed"}
	repo.CreateUser(user)

	wallet := &model.Wallet{UserID: user.ID, Balance: 100}
	repo.CreateWallet(wallet)

	// Credit
	repo.UpdateWalletBalance(nil, wallet.ID, 50, true)
	found, _ := repo.FindWalletByUserID(user.ID)
	if found.Balance != 150 {
		t.Errorf("Balance = %f, want 150", found.Balance)
	}

	// Debit
	repo.UpdateWalletBalance(nil, wallet.ID, 50, false)
	found, _ = repo.FindWalletByUserID(user.ID)
	if found.Balance != 100 {
		t.Errorf("Balance = %f, want 100", found.Balance)
	}
}

func TestRepositoryTransaction(t *testing.T) {
	db := setupDB(t)
	repo := repository.New(db)

	user := &model.User{Name: "Test", CPF: "12345678909", Email: "test@example.com", Phone: "11987654321", Password: "hashed"}
	repo.CreateUser(user)

	transaction := &model.Transaction{FromUserID: &user.ID, ToUserID: &user.ID, Amount: 100, Type: "transfer", Status: "pending"}
	repo.CreateTransaction(nil, transaction)

	transactions, _ := repo.FindTransactionsByUserID(user.ID, 10, 0)
	if len(transactions) != 1 {
		t.Errorf("Transactions count = %d, want 1", len(transactions))
	}
}

// Auth Tests
func TestAuthJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	accessToken, refreshToken, err := auth.GenerateTokenPair(userID, secret, 15*time.Minute, 7*24*time.Hour)
	if err != nil || accessToken == "" || refreshToken == "" {
		t.Error("GenerateTokenPair() failed")
	}

	claims, err := auth.ValidateToken(accessToken, secret)
	if err != nil || claims.UserID != userID {
		t.Error("ValidateToken() failed")
	}

	_, err = auth.ValidateToken("invalid", secret)
	if err == nil {
		t.Error("ValidateToken() should fail for invalid token")
	}
}
