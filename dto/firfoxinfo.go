package dto

import "github.com/shopspring/decimal"

type FirFoxInfo struct {
	Value decimal.Decimal `json:"value"`
	Nonce int             `json:"nonce"`
}
