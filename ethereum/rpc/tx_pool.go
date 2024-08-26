package rpc

func (b *BackendImpl) PendingTransactionsCount() (int, error) {
	res, err := b.clientCtx.Client.UnconfirmedTxs(b.ctx, nil)
	if err != nil {
		return 0, err
	}
	return res.Total, nil
}
