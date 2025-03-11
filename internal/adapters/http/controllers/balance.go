package controllers

import "github.com/physicist2018/gopher-mart-single/internal/interfaces/usecases/balance"

type BalanceController struct {
	balanceUseCase balance.BalanceUseCase
}
