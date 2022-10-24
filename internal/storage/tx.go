package storage

import "github.com/jackc/pgx/v4"

type txOptions struct {
	runAsChild bool
	parentTx   pgx.Tx
}

func defaultTxOptions() *txOptions {
	return &txOptions{
		runAsChild: false,
		parentTx:   nil,
	}
}

func buildOptions(options ...TxOption) *txOptions {
	resultOptions := defaultTxOptions()
	for _, o := range options {
		o.apply(resultOptions)
	}
	return resultOptions
}

type TxOption interface {
	apply(options *txOptions)
}

type txOptionFunc func(options *txOptions)

func (f txOptionFunc) apply(opts *txOptions) { f(opts) }

func asNestedTo(parentTx pgx.Tx) TxOption {
	return txOptionFunc(func(opts *txOptions) {
		opts.runAsChild = true
		opts.parentTx = parentTx
	})
}
