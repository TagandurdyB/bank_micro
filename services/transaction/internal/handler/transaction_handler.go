package handler

import (
	proto "bank_micro/proto/gen"
	"bank_micro/services/transaction/internal/repository"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type TransactionHandler struct {
	proto.UnimplementedTransactionServiceServer
	repo *repository.TransactionRepository
}

func NewTransactionHandler(repo *repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{repo: repo}
}

func (h *TransactionHandler) ReadAll(ctx context.Context, req *proto.ReadAllRequest) (*proto.ReadAllResponse, error) {
	txs, err := h.repo.GetAll(req.AccountId, req.ToAccountId)
	if err != nil {
		return nil, err
	}

	var protoTxs []*proto.TransactionResponse
	for _, tx := range txs {
		toAcc := ""
		if tx.ToAccountID != nil {
			toAcc = tx.ToAccountID.String()
		}

		protoTxs = append(protoTxs, &proto.TransactionResponse{
			Id:              tx.ID.String(),
			AccountId:       tx.AccountID.String(),
			ToAccountId:     toAcc,
			Amount:          tx.Amount,
			TransactionType: tx.TransactionType,
			CreatedAt:       timestamppb.New(tx.CreatedAt),
		})
	}

	return &proto.ReadAllResponse{Transactions: protoTxs}, nil
}
