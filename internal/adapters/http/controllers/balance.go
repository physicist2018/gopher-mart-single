package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/balance"
	"github.com/physicist2018/gopher-mart-single/pkg/middlewares"
)

type BalanceController struct {
	balanceUseCase balance.UseCase
}

func NewBalanceController(balanceUseCase balance.UseCase) *BalanceController {
	return &BalanceController{balanceUseCase: balanceUseCase}
}

func (b *BalanceController) Balance(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middlewares.UserIDKey{}).(int)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	balance, err := b.balanceUseCase.Execute(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(balance)

}
