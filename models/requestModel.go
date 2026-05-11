package models

type TxHashRequest struct {
	TxHash string `json:"tx_hash" binding:"required"`
}
