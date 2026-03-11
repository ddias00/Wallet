package dto

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	CPF      string `json:"cpf" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TransferRequest struct {
	ToUserID    string  `json:"to_user_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

type PixPaymentRequest struct {
	PixKey string  `json:"pix_key" binding:"required"`
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type BoletoRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

type WalletOperationRequest struct {
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Description string  `json:"description"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
