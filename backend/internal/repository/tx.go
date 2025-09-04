package repository

import (
	"context"

	"gorm.io/gorm"
)

type txKeyType string

const txKey txKeyType = "txDB"

func withTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func getDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok && tx != nil {
		return tx
	}
	return db
}

type TxMangager interface {
	WithinTransaction(ctx context.Context, fn func(ctxTx context.Context) error) error
}

type gormTxMangager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) TxMangager {
	return &gormTxMangager{db: db}
}

func (m *gormTxMangager) WithinTransaction(ctx context.Context, fn func(ctxTx context.Context) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		ctxTx := withTx(ctx, tx)
		return fn(ctxTx)
	})
}
