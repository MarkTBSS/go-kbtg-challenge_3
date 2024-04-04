package wallet

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Storer interface {
	Wallets() ([]Wallet, error)
	WalletsByType(walletType string) ([]Wallet, error)
}

type Handler struct {
	store Storer
}

type Err struct {
	Message string `json:"message"`
}

func New(database Storer) *Handler {
	return &Handler{store: database}
}

// WalletHandler
//
// @Summary		Get all wallets
// @Description	Get all wallets
// @Tags			wallet
// @Accept			json
// @Produce		json
// @Success		200	{object}	Wallet
// @Router			/api/v1/wallets [get]
// @Failure		500	{object}	Err
func (handler *Handler) WalletsHandler(context echo.Context) error {
	wallets, err := handler.store.Wallets()
	if err != nil {
		return context.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return context.JSON(http.StatusOK, wallets)
}

func (handler *Handler) WalletsByTypeHandler(context echo.Context) error {
	walletType := context.QueryParam("wallet_type")
	wallets, err := handler.store.WalletsByType(walletType)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return context.JSON(http.StatusOK, wallets)
}
