package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wallet/internal/dto"
	"github.com/wallet/internal/middleware"
	"github.com/wallet/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with CPF, email, phone and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	user, err := h.service.Register(req.Name, req.CPF, req.Email, req.Phone, req.Password)
	if err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	c.JSON(201, dto.Response{Success: true, Message: "User registered successfully", Data: user})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(401, dto.Response{Success: false, Error: err.Error()})
		return
	}

	c.JSON(200, dto.Response{
		Success: true,
		Message: "Login successful",
		Data: dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh Token Request"
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, dto.Response{Success: false, Error: err.Error()})
		return
	}

	c.JSON(200, dto.Response{
		Success: true,
		Message: "Token refreshed successfully",
		Data: dto.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current authenticated user profile
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/auth/me [get]
func (h *Handler) GetProfile(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		c.JSON(404, dto.Response{Success: false, Error: "user not found"})
		return
	}

	c.JSON(200, dto.Response{Success: true, Message: "Profile retrieved", Data: user})
}

// GetBalance godoc
// @Summary Get wallet balance
// @Description Get current wallet balance for authenticated user
// @Tags wallet
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Router /api/v1/wallet/balance [get]
func (h *Handler) GetBalance(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	wallet, err := h.service.GetBalance(userID)
	if err != nil {
		c.JSON(404, dto.Response{Success: false, Error: "wallet not found"})
		return
	}

	c.JSON(200, dto.Response{
		Success: true,
		Message: "Wallet balance retrieved",
		Data: gin.H{
			"balance":   wallet.Balance,
			"wallet_id": wallet.ID,
		},
	})
}

// CreditWallet godoc
// @Summary Credit wallet
// @Description Add funds to wallet
// @Tags wallet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.WalletOperationRequest true "Credit Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/v1/wallet/credit [post]
func (h *Handler) CreditWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.WalletOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	if err := h.service.CreditWallet(userID, req.Amount, req.Description); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	c.JSON(200, dto.Response{Success: true, Message: "Wallet credited successfully"})
}

// DebitWallet godoc
// @Summary Debit wallet
// @Description Remove funds from wallet
// @Tags wallet
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.WalletOperationRequest true "Debit Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /api/v1/wallet/debit [post]
func (h *Handler) DebitWallet(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.WalletOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	if err := h.service.DebitWallet(userID, req.Amount, req.Description); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	c.JSON(200, dto.Response{Success: true, Message: "Wallet debited successfully"})
}

// Transfer godoc
// @Summary Transfer money
// @Description Transfer money to another user
// @Tags transactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.TransferRequest true "Transfer Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/transactions/transfer [post]
func (h *Handler) Transfer(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	toUserID, err := uuid.Parse(req.ToUserID)
	if err != nil {
		c.JSON(400, dto.Response{Success: false, Error: "invalid recipient user ID"})
		return
	}

	if err := h.service.Transfer(userID, toUserID, req.Amount, req.Description); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	c.JSON(200, dto.Response{Success: true, Message: "Transfer completed successfully"})
}

// GetTransactions godoc
// @Summary Get transaction history
// @Description Get transaction history with pagination
// @Tags transactions
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/transactions [get]
func (h *Handler) GetTransactions(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	transactions, err := h.service.GetTransactions(userID, limit, offset)
	if err != nil {
		c.JSON(500, dto.Response{Success: false, Error: "failed to retrieve transactions"})
		return
	}

	c.JSON(200, dto.Response{Success: true, Message: "Transactions retrieved", Data: transactions})
}

// CreatePixPayment godoc
// @Summary Create PIX payment
// @Description Create a PIX payment using PIX key
// @Tags payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.PixPaymentRequest true "PIX Payment Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/pix [post]
func (h *Handler) CreatePixPayment(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.PixPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	if err := h.service.CreatePixPayment(userID, req.PixKey, req.Amount); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	c.JSON(200, dto.Response{Success: true, Message: "PIX payment completed successfully"})
}

// GenerateBoleto godoc
// @Summary Generate boleto
// @Description Generate a boleto for payment
// @Tags payments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.BoletoRequest true "Boleto Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /api/v1/boleto [post]
func (h *Handler) GenerateBoleto(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.BoletoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	boleto, err := h.service.GenerateBoleto(userID, req.Amount)
	if err != nil {
		c.JSON(400, dto.Response{Success: false, Error: err.Error()})
		return
	}

	c.JSON(200, dto.Response{Success: true, Message: "Boleto generated successfully", Data: boleto})
}

// Health godoc
// @Summary Health check
// @Description Check if API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
