package v1

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go-campaign.com/internal/payment/service"
)

type paymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(svc *service.PaymentService) *paymentHandler {
	return &paymentHandler{
		svc: svc,
	}
}

func (h *paymentHandler) GetDetailPayment(ctx *fiber.Ctx) error {
	if ctx.Params("transaction") == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "transaction id is required",
		})
	}

	transactionID, err := uuid.Parse(ctx.Params("transaction"))

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to parse id: %v", err),
		})
	}

	payment, err := h.svc.GetDetailPayment(ctx.Context(), transactionID)

	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": payment,
	})
}
