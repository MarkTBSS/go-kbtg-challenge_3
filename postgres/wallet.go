package postgres

import (
	"time"

	"github.com/MarkTBSS/go-kbtg-challenge_3/wallet"
)

type Wallet struct {
	ID         int       `postgres:"id"`
	UserID     int       `postgres:"user_id"`
	UserName   string    `postgres:"user_name"`
	WalletName string    `postgres:"wallet_name"`
	WalletType string    `postgres:"wallet_type"`
	Balance    float64   `postgres:"balance"`
	CreatedAt  time.Time `postgres:"created_at"`
}

func (postgres *Postgres) Wallets(walletType string) ([]wallet.Wallet, error) {
	var query string
	if walletType == "" {
		query = "SELECT * FROM user_wallet"
	} else {
		query = "SELECT * FROM user_wallet WHERE wallet_type = $1"
	}

	rows, err := postgres.Database.Query(query, walletType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []wallet.Wallet
	for rows.Next() {
		var w wallet.Wallet // Assuming wallet.Wallet is the type of your wallet struct
		err := rows.Scan(
			&w.ID,
			&w.UserID,
			&w.UserName,
			&w.WalletName,
			&w.WalletType,
			&w.Balance,
			&w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}
	return wallets, nil
}
