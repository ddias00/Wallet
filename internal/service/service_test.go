package service

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/wallet/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Mock Repository
type mockRepo struct {
	users          map[string]*model.User
	wallets        map[uuid.UUID]*model.Wallet
	transactions   []model.Transaction
	refreshTokens  map[string]*model.RefreshToken
	findPixKeyFunc func(key string) (*model.PixKey, error)
}

func newMockRepo() *mockRepo {
	m := &mockRepo{
		users:         make(map[string]*model.User),
		wallets:       make(map[uuid.UUID]*model.Wallet),
		transactions:  []model.Transaction{},
		refreshTokens: make(map[string]*model.RefreshToken),
	}
	// Default implementation
	m.findPixKeyFunc = func(key string) (*model.PixKey, error) {
		if key == "valid@pix.com" {
			return &model.PixKey{
				ID:     uuid.New(),
				UserID: uuid.New(),
				Key:    key,
				Status: "active",
			}, nil
		}
		return nil, gorm.ErrRecordNotFound
	}
	return m
}

func (m *mockRepo) CreateUser(user *model.User) error {
	user.ID = uuid.New()
	m.users[user.Email] = user
	return nil
}

func (m *mockRepo) FindUserByEmail(email string) (*model.User, error) {
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) FindUserByCPF(cpf string) (*model.User, error) {
	for _, user := range m.users {
		if user.CPF == cpf {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) FindUserByID(id uuid.UUID) (*model.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) CreateWallet(wallet *model.Wallet) error {
	wallet.ID = uuid.New()
	m.wallets[wallet.UserID] = wallet
	return nil
}

func (m *mockRepo) FindWalletByUserID(userID uuid.UUID) (*model.Wallet, error) {
	if wallet, ok := m.wallets[userID]; ok {
		return wallet, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) UpdateWalletBalance(tx *gorm.DB, walletID uuid.UUID, amount float64, isCredit bool) error {
	return nil
}

func (m *mockRepo) CreateTransaction(tx *gorm.DB, transaction *model.Transaction) error {
	return nil
}

func (m *mockRepo) FindTransactionsByUserID(userID uuid.UUID, limit, offset int) ([]model.Transaction, error) {
	return m.transactions, nil
}

func (m *mockRepo) UpdateTransactionStatus(tx *gorm.DB, id uuid.UUID, status string) error {
	return nil
}

func (m *mockRepo) FindPixKeyByKey(key string) (*model.PixKey, error) {
	return m.findPixKeyFunc(key)
}

func (m *mockRepo) CreateBoleto(boleto *model.Boleto) error {
	return nil
}

func (m *mockRepo) CreateRefreshToken(token *model.RefreshToken) error {
	m.refreshTokens[token.Token] = token
	return nil
}

func (m *mockRepo) FindRefreshToken(token string) (*model.RefreshToken, error) {
	if rt, ok := m.refreshTokens[token]; ok {
		return rt, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockRepo) DeleteRefreshTokensByUserID(userID uuid.UUID) error {
	return nil
}

func (m *mockRepo) Transaction(fn func(*gorm.DB) error) error {
	return fn(nil)
}

func TestServiceRegisterWithMock(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	tests := []struct {
		name    string
		cpf     string
		email   string
		phone   string
		wantErr bool
	}{
		{"valid", "12345678909", "test@example.com", "11987654321", false},
		{"invalid_cpf", "123", "test2@example.com", "11987654321", true},
		{"invalid_email", "12345678910", "invalid", "11987654321", true},
		{"invalid_phone", "12345678911", "test3@example.com", "123", true},
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

func TestServiceLoginWithMock(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	// Create test user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	testUser := &model.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	repo.users[testUser.Email] = testUser

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{"valid", "test@example.com", "password123", false},
		{"wrong_password", "test@example.com", "wrong", true},
		{"wrong_email", "wrong@example.com", "password123", true},
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

func TestServiceGetBalance(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	userID := uuid.New()
	wallet := &model.Wallet{
		ID:      uuid.New(),
		UserID:  userID,
		Balance: 100.50,
	}
	repo.wallets[userID] = wallet

	result, err := svc.GetBalance(userID)
	if err != nil {
		t.Errorf("GetBalance() error = %v", err)
	}
	if result.Balance != 100.50 {
		t.Errorf("GetBalance() balance = %v, want 100.50", result.Balance)
	}
}

func TestServiceGetUserByID(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	testUser := &model.User{
		ID:    uuid.New(),
		Name:  "Test User",
		Email: "test@example.com",
	}
	repo.users[testUser.Email] = testUser

	result, err := svc.GetUserByID(testUser.ID)
	if err != nil {
		t.Errorf("GetUserByID() error = %v", err)
	}
	if result.Name != "Test User" {
		t.Errorf("GetUserByID() name = %v, want Test User", result.Name)
	}
}

func TestValidators(t *testing.T) {
	tests := []struct {
		name  string
		input string
		fn    func(string) bool
		want  bool
	}{
		{"valid_cpf", "12345678909", validateCPF, true},
		{"invalid_cpf", "123", validateCPF, false},
		{"valid_email", "test@example.com", validateEmail, true},
		{"invalid_email", "invalid", validateEmail, false},
		{"valid_phone", "11987654321", validatePhone, true},
		{"invalid_phone", "123", validatePhone, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(tt.input); got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestServiceCreditWallet(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	userID := uuid.New()
	wallet := &model.Wallet{ID: uuid.New(), UserID: userID, Balance: 100}
	repo.wallets[userID] = wallet

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"valid", 50.0, false},
		{"zero_amount", 0, true},
		{"negative_amount", -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreditWallet(userID, tt.amount, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("CreditWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceDebitWallet(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	userID := uuid.New()
	wallet := &model.Wallet{ID: uuid.New(), UserID: userID, Balance: 100}
	repo.wallets[userID] = wallet

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"valid", 50.0, false},
		{"insufficient_balance", 200, true},
		{"zero_amount", 0, true},
		{"negative_amount", -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.DebitWallet(userID, tt.amount, "test")
			if (err != nil) != tt.wantErr {
				t.Errorf("DebitWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceTransfer(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	fromUserID := uuid.New()
	toUserID := uuid.New()
	repo.wallets[fromUserID] = &model.Wallet{ID: uuid.New(), UserID: fromUserID, Balance: 100}
	repo.wallets[toUserID] = &model.Wallet{ID: uuid.New(), UserID: toUserID, Balance: 50}

	tests := []struct {
		name    string
		from    uuid.UUID
		to      uuid.UUID
		amount  float64
		wantErr bool
	}{
		{"valid", fromUserID, toUserID, 50.0, false},
		{"insufficient_balance", fromUserID, toUserID, 200, true},
		{"same_user", fromUserID, fromUserID, 10, true},
		{"zero_amount", fromUserID, toUserID, 0, true},
		{"negative_amount", fromUserID, toUserID, -10, true},
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

func TestServiceGenerateBoleto(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	userID := uuid.New()

	tests := []struct {
		name    string
		amount  float64
		wantErr bool
	}{
		{"valid", 100.0, false},
		{"zero_amount", 0, true},
		{"negative_amount", -10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boleto, err := svc.GenerateBoleto(userID, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateBoleto() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && boleto == nil {
				t.Error("GenerateBoleto() returned nil boleto")
			}
		})
	}
}

func TestServiceGetTransactions(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	userID := uuid.New()
	transactions, err := svc.GetTransactions(userID, 10, 0)
	if err != nil {
		t.Errorf("GetTransactions() error = %v", err)
	}
	if transactions == nil {
		t.Error("GetTransactions() returned nil")
	}
}

func TestServiceRegisterDuplicates(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	// Register first user
	svc.Register("Test", "12345678909", "test@example.com", "11987654321", "password123")

	// Try duplicate email
	_, err := svc.Register("Test2", "12345678910", "test@example.com", "11987654322", "password123")
	if err == nil {
		t.Error("Register() should fail with duplicate email")
	}

	// Try duplicate CPF
	_, err = svc.Register("Test3", "12345678909", "test2@example.com", "11987654323", "password123")
	if err == nil {
		t.Error("Register() should fail with duplicate CPF")
	}
}

func TestServiceLoginBlocked(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	blockedUser := &model.User{
		ID:        uuid.New(),
		Email:     "blocked@example.com",
		Password:  string(hashedPassword),
		IsBlocked: true,
	}
	repo.users[blockedUser.Email] = blockedUser

	_, _, err := svc.Login("blocked@example.com", "password123")
	if err == nil {
		t.Error("Login() should fail for blocked user")
	}
}

func TestServiceCreatePixPayment(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	fromUserID := uuid.New()
	toUserID := uuid.New()
	repo.wallets[fromUserID] = &model.Wallet{ID: uuid.New(), UserID: fromUserID, Balance: 100}
	repo.wallets[toUserID] = &model.Wallet{ID: uuid.New(), UserID: toUserID, Balance: 50}

	// Configure mock to return toUserID for valid PIX key
	repo.findPixKeyFunc = func(key string) (*model.PixKey, error) {
		if key == "valid@pix.com" {
			return &model.PixKey{
				ID:     uuid.New(),
				UserID: toUserID,
				Key:    key,
				Status: "active",
			}, nil
		}
		return nil, gorm.ErrRecordNotFound
	}

	tests := []struct {
		name    string
		pixKey  string
		amount  float64
		wantErr bool
	}{
		{"valid", "valid@pix.com", 50.0, false},
		{"invalid_pix_key", "invalid@pix.com", 50.0, true},
		{"zero_amount", "valid@pix.com", 0, true},
		{"negative_amount", "valid@pix.com", -10, true},
		{"insufficient_balance", "valid@pix.com", 200, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.CreatePixPayment(fromUserID, tt.pixKey, tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePixPayment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceRefreshToken(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	// Create user and login to get valid tokens
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	testUser := &model.User{
		ID:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	repo.users[testUser.Email] = testUser

	_, refreshToken, _ := svc.Login("test@example.com", "password123")

	// Test refresh with valid token
	newAccessToken, newRefreshToken, err := svc.RefreshToken(refreshToken)
	if err != nil {
		t.Errorf("RefreshToken() error = %v", err)
	}
	if newAccessToken == "" || newRefreshToken == "" {
		t.Error("RefreshToken() returned empty tokens")
	}

	// Test refresh with invalid token
	_, _, err = svc.RefreshToken("invalid-token")
	if err == nil {
		t.Error("RefreshToken() should fail with invalid token")
	}
}

func TestServiceWalletNotFound(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	userID := uuid.New()

	// Test CreditWallet with non-existent wallet
	err := svc.CreditWallet(userID, 50, "test")
	if err == nil {
		t.Error("CreditWallet() should fail with non-existent wallet")
	}

	// Test DebitWallet with non-existent wallet
	err = svc.DebitWallet(userID, 50, "test")
	if err == nil {
		t.Error("DebitWallet() should fail with non-existent wallet")
	}

	// Test Transfer with non-existent sender wallet
	toUserID := uuid.New()
	repo.wallets[toUserID] = &model.Wallet{ID: uuid.New(), UserID: toUserID, Balance: 50}
	err = svc.Transfer(userID, toUserID, 50, "test")
	if err == nil {
		t.Error("Transfer() should fail with non-existent sender wallet")
	}

	// Test Transfer with non-existent recipient wallet
	repo.wallets[userID] = &model.Wallet{ID: uuid.New(), UserID: userID, Balance: 100}
	delete(repo.wallets, toUserID)
	err = svc.Transfer(userID, toUserID, 50, "test")
	if err == nil {
		t.Error("Transfer() should fail with non-existent recipient wallet")
	}
}

func TestServicePixSameUser(t *testing.T) {
	repo := newMockRepo()
	svc := &Service{
		repo:       repo,
		jwtSecret:  "test-secret",
		accessExp:  15 * time.Minute,
		refreshExp: 7 * 24 * time.Hour,
	}

	userID := uuid.New()
	repo.wallets[userID] = &model.Wallet{ID: uuid.New(), UserID: userID, Balance: 100}

	// Mock PIX key pointing to same user
	repo.findPixKeyFunc = func(key string) (*model.PixKey, error) {
		return &model.PixKey{
			ID:     uuid.New(),
			UserID: userID,
			Key:    key,
			Status: "active",
		}, nil
	}

	err := svc.CreatePixPayment(userID, "same@pix.com", 50)
	if err == nil {
		t.Error("CreatePixPayment() should fail when transferring to yourself")
	}
}

func TestGenerateBarcode(t *testing.T) {
	barcode := generateBarcode()
	if len(barcode) != 47 {
		t.Errorf("generateBarcode() length = %d, want 47", len(barcode))
	}
}
