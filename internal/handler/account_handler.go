package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mitayush/visa-txn/internal/apperrors"
	"github.com/mitayush/visa-txn/internal/dto"
	"github.com/mitayush/visa-txn/internal/service"
)

var accountCreateErrHTTPStatus = map[error]int{
	apperrors.ErrAccountExists: http.StatusConflict,
}

var accountGetErrHTTPStatus = map[error]int{
	apperrors.ErrAccountNotFound:  http.StatusNotFound,
	apperrors.ErrInvalidAccountID: http.StatusBadRequest,
}

type AccountHandler struct {
	svc *service.AccountService
}

func NewAccountHandler(svc *service.AccountService) *AccountHandler {
	return &AccountHandler{svc: svc}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	acc := dto.CreateAccountRequestToModel(&req)

	account, err := h.svc.CreateAccount(r.Context(), acc)
	if err != nil {
		writeJSON(w, accountCreateErrHTTPStatus[err], dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, dto.CreateAccountResponse{
		AccountID:      account.AccountID,
		DocumentNumber: account.DocumentNumber,
		CreatedAt:      account.CreatedAt,
	})
}

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: "invalid account ID"})
		return
	}
	acc, err := h.svc.GetAccount(r.Context(), accountID)
	if err != nil {
		writeJSON(w, accountGetErrHTTPStatus[err], dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, dto.AccountToResponse(acc))
}
