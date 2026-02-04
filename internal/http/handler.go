package api

import (
	"fmt"
	"strconv"
	"subscriptions/internal/models"
	"subscriptions/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	storage service.SubscriptionService
	logger  *zap.Logger
}

func NewHandler(subscriptions service.SubscriptionService, logger *zap.Logger) *Handler {
	return &Handler{
		storage: subscriptions,
		logger:  logger,
	}
}

func (h *Handler) InitRoutes(ginRouter *gin.RouterGroup) {
	router := ginRouter.Group("/subscriptions")
	{
		router.POST("", h.CreateSubscription)
		router.GET("/:id", h.GetSubscription)
		router.PUT("/:id", h.UpdateSubscription)
		router.DELETE("/:id", h.DeleteSubscription)
		router.GET("", h.ListSubscriptions)
		router.POST("/calculate", h.CalculateTotalCost)
	}
}

// CreateSubscription Добавить информацию о подписке
// @Summary Создать новую подписку
// @Description Создает новую запись о подписке пользователя
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body models.CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(c *gin.Context) {
	var req models.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	subscription := &models.Subscription{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	err := h.storage.CreateSubscription(c.Request.Context(), subscription)
	if err != nil {
		h.logger.Error("Failed to create subscription", zap.Error(err))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Subscription created")
	c.JSON(201, subscription)
}

// GetSubscription получает подписку по ID
// @Summary Получить подписку
// @Description Возвращает информацию о подписке по её идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path string true "ID подписки"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format (must be integer)"})
		return
	}

	subscription, err := h.storage.GetSubscriptionByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to get subscription", zap.Error(err))
		c.JSON(404, gin.H{"error": "Subscription not found"})
		return
	}

	h.logger.Info(fmt.Sprintf("Get subscription with ID = %d", id))
	c.JSON(200, subscription)
}

// UpdateSubscription обновляет подписку
// @Summary Обновить подписку
// @Description Обновляет информацию о подписке
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "ID подписки"
// @Param input body models.UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid ID format (must be integer)", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid ID format (must be integer)"})
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	subscription, err := h.storage.UpdateSubscription(c.Request.Context(), uint(id), &req)
	if err != nil {
		h.logger.Error("Failed to update subscription", zap.Error(err))
		c.JSON(404, gin.H{"error": "Subscription not found"})
		return
	}

	h.logger.Info(fmt.Sprintf("Subscription with ID = %d updated successfully", id))
	c.JSON(200, subscription)
}

// DeleteSubscription удаляет подписку
// @Summary Удалить подписку по идентификатору
// @Description Удаляет запись о подписке
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Error("Invalid ID format (must be integer)", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid ID format (must be integer)"})
		return
	}

	if id <= 0 {
		h.logger.Error("Invalid ID format (must be greater than 0)", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid ID format (must be greater than 0)"})
		return
	}

	err = h.storage.DeleteSubscription(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("Failed to delete subscription", zap.Error(err))
		c.JSON(404, gin.H{"error": "Subscription not found"})
		return
	}

	h.logger.Info(fmt.Sprintf("Subscription with ID = %d deleted successfully", id))
	// ответ успешный, но без тела
	c.Status(204)
}

// ListSubscriptions получает список всех подписок
// @Summary Список всех подписок
// @Description Возвращает список всех подписок (с учетом смещения и количества записей)
// @Tags subscriptions
// @Produce json
// @Param limit query int false "Количество записей" default(100)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *Handler) ListSubscriptions(c *gin.Context) {
	limit := c.DefaultQuery("limit", "100")
	offset := c.DefaultQuery("offset", "0")

	subscriptions, err := h.storage.ListSubscriptions(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list subscriptions", zap.Error(err))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	h.logger.Info("Get list of all subscriptions")
	c.JSON(200, subscriptions)
}

// CalculateTotalCost рассчитывает суммарную стоимость подписок
// @Summary Рассчитать суммарную стоимость
// @Description Рассчитывает суммарную стоимость всех подписок за период с фильтрацией по пользователю, датам и имени подписки
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body models.TotalCostRequest true "Фильтры для расчета (user_id, start_date, end_date, service_name)"
// @Success 200 {object} models.TotalCostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/calculate [post]
func (h *Handler) CalculateTotalCost(c *gin.Context) {
	var req models.TotalCostRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	totalCost, err := h.storage.CalculateSubscriptionsTotalCost(
		c.Request.Context(),
		req.StartDate,
		req.EndDate,
		req.UserID,
		req.ServiceName,
	)

	if err != nil {
		h.logger.Error("Failed to calculate total cost", zap.Error(err))
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	response := models.TotalCostResponse{
		TotalCost: totalCost,
	}
	response.Period.StartDate = req.StartDate
	response.Period.EndDate = req.EndDate
	response.Filters.UserID = req.UserID
	response.Filters.ServiceName = req.ServiceName

	c.JSON(200, response)
}
