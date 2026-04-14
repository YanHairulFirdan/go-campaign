package app

import (
	"github.com/gofiber/fiber/v2"
)

type Bootable = func(router fiber.Router, deps *Dependencies)
