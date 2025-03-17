package orderservice

import (
	"context"
	"encoding/json"
	"github.com/physicist2018/gopher-mart-single/internal/entities"
	orderrepo "github.com/physicist2018/gopher-mart-single/internal/interfaces/repositories/order"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

// OrderService предоставляет методы для работы с заказами, включая проверку их статуса в системе начисления баллов.
type OrderService struct {
	ordersRepository orderrepo.Repository // Репозиторий для работы с заказами
	baseURL          string               // Базовый URL системы начисления баллов
	logger           *zerolog.Logger      // Логгер для записи событий
}

// OrderResponse представляет структуру ответа от системы начисления баллов.
type OrderResponse struct {
	OrderID string  `json:"order"`   // ID заказа
	Status  string  `json:"status"`  // Статус заказа (например, "PROCESSED", "INVALID")
	Accrual float64 `json:"accrual"` // Начисленные баллы
}

// NewOrderService создает новый экземпляр OrderService.
// ordersRepository - репозиторий для работы с заказами.
// baseURL - базовый URL системы начисления баллов.
// logger - логгер для записи событий.
// Возвращает указатель на OrderService.
func NewOrderService(ordersRepository orderrepo.Repository, baseURL string, logger *zerolog.Logger) *OrderService {
	return &OrderService{
		ordersRepository: ordersRepository,
		baseURL:          baseURL,
		logger:           logger,
	}
}

// CheckOrdersForAccrual проверяет статус заказов, которые еще не обработаны (не имеют статуса INVALID или PROCESSED).
// Возвращает ошибку, если не удалось получить список заказов.
func (os *OrderService) CheckOrdersForAccrual() error {
	// Получаем список заказов, которые еще не обработаны
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	orders, err := os.ordersRepository.GetPendingOrders(ctx)
	cancel()
	if err != nil {
		return err
	}

	// Параллельно проверяем статус каждого заказа в системе начисления баллов
	for _, order := range orders {
		go os.checkOrderStatus(order)
	}

	return nil
}

// checkOrderStatus проверяет статус заказа в системе начисления баллов и обновляет его в репозитории.
// order - заказ, статус которого нужно проверить.
func (os *OrderService) checkOrderStatus(order entities.Order) {
	url := os.baseURL + "/api/orders/" + order.ID
	for {
		// Выполняем HTTP-запрос к системе начисления баллов
		resp, err := http.Get(url)
		if err != nil {
			os.logger.Error().Err(err).Msgf("Failed to fetch order %s", order.ID)
			return
		}
		defer resp.Body.Close()

		// Обрабатываем случай, когда превышен лимит запросов (429 Too Many Requests)
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			retryDuration, err := time.ParseDuration(retryAfter + "s") // Интерпретируем как секунды
			if err != nil {
				os.logger.Error().Err(err).Msgf("Invalid Retry-After header for order %s", order.ID)
				return
			}
			os.logger.Warn().Msgf("Rate limit exceeded for order %s, retrying after %s", order.ID, retryDuration)
			time.Sleep(retryDuration)
			continue // Повторяем запрос после указанного интервала
		}

		// Обрабатываем случай, когда заказ не зарегистрирован в системе (204 No Content)
		if resp.StatusCode == http.StatusNoContent {
			log.Info().Msgf("заказ %s не зарегистрирован в системе расчёта.", order.ID)
			return // Заканчиваем проверку для этого заказа
		}

		// Обрабатываем внутреннюю ошибку сервера (500 Internal Server Error)
		if resp.StatusCode == http.StatusInternalServerError {
			log.Error().Msgf("внутренняя ошибка сервера при обработке заказа %s", order.ID)
			return
		}

		// Декодируем JSON-ответ
		var orderResp OrderResponse
		if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
			log.Error().Err(err).Msgf("Failed to decode response for order %s", order.ID)
			return
		}

		// Обрабатываем статус заказа
		switch orderResp.Status {
		case "PROCESSED":
			ctxUpdate, ctxUpdateCancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = os.ordersRepository.UpdateOrderStatus(ctxUpdate, orderResp.OrderID, orderResp.Status, orderResp.Accrual)
			if err != nil {
				os.logger.Error().Err(err).Msgf("Failed to update order %s", order.ID)
			} else {
				os.logger.Info().Msgf("Order %s updated: Status %s, Accrual %f", order.ID, orderResp.Status, orderResp.Accrual)
			}
			ctxUpdateCancel()
			//
			ctxCreate, ctxCreateCancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = os.ordersRepository.CreateAccrual(ctxCreate, order.ID, orderResp.OrderID, orderResp.Accrual)
			ctxCreateCancel()
		case "INVALID":
			// Обновляем статус заказа в репозитории
			ctxUpdate, ctxUpdateCancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = os.ordersRepository.UpdateOrderStatus(ctxUpdate, orderResp.OrderID, orderResp.Status, orderResp.Accrual)
			ctxUpdateCancel()
			if err != nil {
				os.logger.Error().Err(err).Msgf("Failed to update order %s", order.ID)
			} else {
				os.logger.Info().Msgf("Order %s updated: Status %s, Accrual %f", order.ID, orderResp.Status, orderResp.Accrual)
			}

		case "REGISTERED", "PROCESSING":
			// Заказ все еще в процессе обработки

			os.logger.Info().Msgf("Order %s still in progress: Status %s", order.ID, orderResp.Status)
			ctxUpdate, ctxUpdateCancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = os.ordersRepository.UpdateOrderStatus(ctxUpdate, orderResp.OrderID, "PROCESSING", 0)
			ctxUpdateCancel()
		}

		break
	}
}
