package storage

import (
	"context"
	"errors"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s *Storage) Revenue(ctx context.Context, UserId int64, ServiceId int64, OrderId int64, Sum decimal.Decimal, description *string) error {
	logger := s.Logger.With(zap.Int64("userID", UserId), zap.Int64("ServiceID", ServiceId), zap.Int64("OrderID", OrderId))
	logger.Debug("reservation of funds")

	var amount decimal.Decimal
	var exist bool

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				logger.Error("error rolls back the transaction", zap.Error(err))
			}
		}
	}()

	firstSelectQuery := `SELECT price FROM deferred_expenses WHERE account_id = $3 AND service_id = $4 AND order_id = $1 AND operation = $2;`

	err = tx.QueryRow(
		ctx,
		firstSelectQuery,
		OrderId,
		ExpensesTypeReservation,
		UserId,
		ServiceId,
	).Scan(&amount)

	if err != nil {
		if amount.Equal(decimal.NewFromInt(0)) {
			logger.Error("order exists error", zap.Error(ErrReserveExist))
			return ErrReserveExist
		}
		logger.Error("query error", zap.Error(err))
		return err
	}

	if Sum.GreaterThan(amount) {
		logger.Error("revenue recognition error", zap.Error(ErrRevenue))
		return ErrRevenue
	}

	id, _, err := s.Transfer(ctx, reserveAccountID, cacheBookAccountID, Sum, description, asNestedTo(tx))
	if err != nil {
		switch {
		case errors.Is(err, ErrSerialization):
			logger.Warn("transaction isolation level error", zap.Error(err))
			return ErrSerialization
		case errors.Is(err, ErrTransfer):
			logger.Error("insufficient funds on the sender's account", zap.Error(ErrTransfer))
			return ErrTransfer
		case errors.Is(err, ErrUserAvailability):
			logger.Error("error returning user balance with specified id: user does not exist", zap.Error(err))
			return ErrUserAvailability
		default:
			logger.Error("error updating balance", zap.Error(err))
			return err
		}
	}

	secondSelectQuery := `SELECT EXISTS (SELECT 1 FROM deferred_expenses df WHERE (account_id = $4 AND service_id = $5 AND order_id = $2 AND operation = $3)
							AND (EXISTS(SELECT 1 FROM deferred_expenses WHERE order_id = $2 AND operation = $1)
							OR EXISTS(SELECT 1 FROM consolidated_report cr WHERE cr.order_id = df.order_id)));`

	err = tx.QueryRow(
		ctx,
		secondSelectQuery,
		ExpensesTypeUnreservation,
		OrderId,
		ExpensesTypeReservation,
		UserId,
		ServiceId,
	).Scan(&exist)
	if err != nil {
		logger.Error("failed to Query", zap.Error(err))
		return err
	}

	if exist {
		s.Logger.Error("", zap.Error(ErrRecordExist))
		return ErrRecordExist
	}

	firstInsertExec := `INSERT INTO consolidated_report (account_id, service_id, order_id, sum, tx_id)
							VALUES ($1, $2, $3, $4, $5);`

	_, err = tx.Exec(
		ctx,
		firstInsertExec,
		UserId,
		ServiceId,
		OrderId,
		Sum,
		id,
	)
	if err != nil {
		logger.Error("failed to Insert", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	return err
}
