package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (s *Storage) Reservation(ctx context.Context, UserId int64, ServiceId int64, OrderId int64, Price decimal.Decimal, description *string) error {
	logger := s.Logger.With(zap.Int64("userID", UserId), zap.Int64("ServiceID", ServiceId), zap.Int64("OrderID", OrderId))
	logger.Debug("reservation of funds")

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

	id, _, err := s.Transfer(ctx, UserId, reserveAccountID, Price, description, asNestedTo(tx))
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

	firstInsertExec := `INSERT INTO deferred_expenses (account_id, service_id, order_id, operation, price, tx_id)
			VALUES ($1, $2, $3, $4, $5, $6);`

	_, err = tx.Exec(
		ctx,
		firstInsertExec,
		UserId,
		ServiceId,
		OrderId,
		ExpensesTypeReservation,
		Price,
		id,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			logger.Error("adding unique order error", zap.Error(err))
			return ErrOrderId
		}
		logger.Error("failed to insert record", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	return err
}
