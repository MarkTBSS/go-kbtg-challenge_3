package wallet

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

type StubStore struct {
	wallets []Wallet
	err     error
}

func (s StubStore) Wallets() ([]Wallet, error) {
	return s.wallets, s.err
}

func (s StubStore) GetByType(walletType string) ([]Wallet, error) {
	// Implement GetByType method for StubStore
	// For testing purposes, you can return a list of wallets filtered by walletType
	var filteredWallets []Wallet
	for _, w := range s.wallets {
		if w.WalletType == walletType {
			filteredWallets = append(filteredWallets, w)
		}
	}
	return filteredWallets, nil
}

func TestWallets(t *testing.T) {
	t.Run("given unable to get wallets should return 500 and error message", func(t *testing.T) {
		echoInstance := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/wallets", nil)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		// Create a recorder to record the response
		record := httptest.NewRecorder()
		// Create a context from the request and recorder
		context := echoInstance.NewContext(request, record)
		stubError := StubStore{err: echo.ErrInternalServerError}
		s := New(stubError)
		// Call the WalletHandler function
		s.WalletHandler(context)
		// Check the status code of the response
		if record.Code != http.StatusInternalServerError {
			t.Errorf("expected status code %d but got %d", http.StatusInternalServerError, record.Code)
		}
	})

	t.Run("given user able to getting wallet should return list of wallets", func(t *testing.T) {
		echoInstance := echo.New()
		request := httptest.NewRequest(http.MethodGet, "/api/v1/wallets", nil)
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		// Create a recorder to record the response
		record := httptest.NewRecorder()
		// Create a context from the request and recorder
		context := echoInstance.NewContext(request, record)
		timeMock := time.Now().UTC()
		stubStore := StubStore{
			wallets: []Wallet{
				{
					ID:         1,
					UserID:     1,
					UserName:   "John Doe",
					WalletName: "John Savings",
					WalletType: "Savings",
					Balance:    1000.00,
					CreatedAt:  timeMock,
				},
				{
					ID:         2,
					UserID:     2,
					UserName:   "Jane Doe",
					WalletName: "Jane Savings",
					WalletType: "Savings",
					Balance:    2000.00,
					CreatedAt:  timeMock,
				},
			},
		}
		s := New(stubStore)
		// Call the WalletHandler function
		err := s.WalletHandler(context)
		// Check if an error occurred
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		wantUserName := "John Doe"
		want := []Wallet{
			{
				ID:         1,
				UserID:     1,
				UserName:   wantUserName,
				WalletName: "John Savings",
				WalletType: "Savings",
				Balance:    1000.00,
				CreatedAt:  timeMock,
			},
			{
				ID:         2,
				UserID:     2,
				UserName:   "Jane Doe",
				WalletName: "Jane Savings",
				WalletType: "Savings",
				Balance:    2000.00,
				CreatedAt:  timeMock,
			},
		}
		gotJson := record.Body.Bytes()
		var got []Wallet
		if err := json.Unmarshal(gotJson, &got); err != nil {
			t.Errorf("unable to unmarshal json: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v but got %v", want, got)
		}
	})
}
