package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"subscription-service/internal/model"
	"subscription-service/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscription создает новую подписку
// @Summary Создать подписку
// @Description Создает новую запись о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.service.CreateSubscription(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subscription)
}

// GetSubscription получает подписку по ID
// @Summary Получить подписку
// @Description Возвращает подписку по её ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	id := c.Param("id")

	subscription, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// UpdateSubscription обновляет подписку
// @Summary Обновить подписку
// @Description Обновляет данные подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param input body model.UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateSubscription(c.Request.Context(), id, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription updated successfully"})
}

// DeleteSubscription удаляет подписку
// @Summary Удалить подписку
// @Description Удаляет подписку по её ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteSubscription(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "subscription deleted successfully"})
}

// ListSubscriptions возвращает список подписок
// @Summary Список подписок
// @Description Возвращает список подписок с возможностью фильтрации
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	var userID *string
	var serviceName *string

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		userID = &userIDStr
	}

	if serviceNameStr := c.Query("service_name"); serviceNameStr != "" {
		serviceName = &serviceNameStr
	}

	subscriptions, err := h.service.ListSubscriptions(c.Request.Context(), userID, serviceName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}

// GetSummary возвращает суммарную стоимость подписок за период
// @Summary Сумма подписок за период
// @Description Возвращает суммарную стоимость всех подписок за выбранный период
// @Tags summary
// @Accept json
// @Produce json
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Param start_period query string true "Начало периода (MM-YYYY)"
// @Param end_period query string true "Конец периода (MM-YYYY)"
// @Success 200 {object} model.SubscriptionSummary
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /summary [get]
func (h *SubscriptionHandler) GetSummary(c *gin.Context) {
	var req model.SummaryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.ServiceName != nil && *req.ServiceName == "" {
		req.ServiceName = nil
	}

	summary, err := h.service.GetSummary(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
