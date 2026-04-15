package middleware

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type QueryNormalization map[string]int

func PaginationQueryNormalizer(q QueryNormalization) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		args := ctx.Request().URI().QueryArgs()

		for k, minValue := range q {
			pageValue := strings.TrimSpace(string(args.Peek(k)))

			if pageValue != "" {
				if num, err := strconv.Atoi(pageValue); err != nil || num < minValue {
					args.Set(k, strconv.Itoa(minValue))
				}
			}
		}

		return ctx.Next()
	}
}
