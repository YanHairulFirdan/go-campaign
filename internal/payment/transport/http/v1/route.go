package v1

import "github.com/gofiber/fiber/v2"

func RegisterRoute(fiberApp fiber.Router, h *paymentHandler) {
	r := fiberApp.Group("/payments")
	r.Get("/:transaction", h.GetDetailPayment)
}
