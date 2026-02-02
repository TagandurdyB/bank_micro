package handler

import (
	proto "bank_micro/proto/gen"
	"bank_micro/services/account/internal/model"
	"bank_micro/services/account/internal/service"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type AccountHandler struct {
	proto.UnimplementedAccountServiceServer
	service *service.AccountService
}

func NewAccountHandler(s *service.AccountService) *AccountHandler {
	return &AccountHandler{service: s}
}

// 1. CreateAccount
func (h *AccountHandler) CreateAccount(ctx context.Context, req *proto.CreateAccountRequest) (*proto.AccountResponse, error) {
	acc, err := h.service.Create(req.Currency, req.InitialBalance)
	if err != nil {
		return nil, err
	}
	return h.mapToProto(acc), nil
}

// 2. GetAccount
func (h *AccountHandler) GetAccount(ctx context.Context, req *proto.GetAccountRequest) (*proto.AccountResponse, error) {
	acc, err := h.service.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return h.mapToProto(acc), nil
}

// 3. GetAllAccounts
func (h *AccountHandler) GetAllAccounts(ctx context.Context, req *proto.Empty) (*proto.GetAllAccountsResponse, error) {
	accounts, err := h.service.GetAll()
	if err != nil {
		return nil, err
	}

	var protoAccounts []*proto.AccountResponse
	for _, acc := range accounts {
		protoAccounts = append(protoAccounts, h.mapToProto(&acc))
	}

	return &proto.GetAllAccountsResponse{Accounts: protoAccounts}, nil
}

// 4. UpdateAccount
func (h *AccountHandler) UpdateAccount(ctx context.Context, req *proto.UpdateAccountRequest) (*proto.AccountResponse, error) {
	acc, err := h.service.Update(req.Id, req.Balance, req.IsLocked)
	if err != nil {
		return nil, err
	}
	return h.mapToProto(acc), nil
}

// 5. DeleteAccount
func (h *AccountHandler) DeleteAccount(ctx context.Context, req *proto.DeleteAccountRequest) (*proto.DeleteAccountResponse, error) {
	err := h.service.Delete(req.Id)
	if err != nil {
		return &proto.DeleteAccountResponse{Success: false}, err
	}
	return &proto.DeleteAccountResponse{Success: true}, nil
}

// 6. Deposit (RabbitMQ)
func (h *AccountHandler) Deposit(ctx context.Context, req *proto.DepositRequest) (*proto.DepositResponse, error) {
	msg, err := h.service.ProcessDeposit(req.Id, req.Amount)
	if err != nil {
		return nil, err
	}
	return &proto.DepositResponse{Message: msg}, nil
}

// Model -> Proto
func (h *AccountHandler) mapToProto(acc *model.Account) *proto.AccountResponse {
	return &proto.AccountResponse{
		Id:        acc.ID.String(),
		Balance:   acc.Balance,
		Currency:  acc.Currency,
		IsLocked:  acc.IsLocked,
		CreatedAt: timestamppb.New(acc.CreatedAt),
	}
}
