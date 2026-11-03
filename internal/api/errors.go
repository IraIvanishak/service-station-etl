package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

// errorHandler turns every error returned by a handler into
// {"error": "..."} with a matching status code.
func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var fe *fiber.Error
	var pqErr *pq.Error
	switch {
	case errors.As(err, &fe):
		code = fe.Code
	case errors.Is(err, sql.ErrNoRows):
		code = fiber.StatusNotFound
	case errors.As(err, &pqErr) && pqErr.Code == foreignKeyViolation:
		// Deleting a row others still point at is a conflict; inserting a
		// reference to a row that does not exist is a bad request.
		code = fiber.StatusBadRequest
		if strings.Contains(pqErr.Detail, "is still referenced") {
			code = fiber.StatusConflict
		}
	}
	if code >= fiber.StatusInternalServerError {
		log.Printf("%s %s: %v", c.Method(), c.Path(), err)
	}
	return c.Status(code).JSON(fiber.Map{"error": err.Error()})
}

const foreignKeyViolation = "23503"

// deleted turns "0 rows affected" from an :execrows delete into a 404.
func deleted(rows int64, err error) error {
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func badRequest(format string, args ...any) error {
	return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf(format, args...))
}

// idParam reads an integer route parameter such as /repairs/:id.
func idParam(c *fiber.Ctx, name string) (int32, error) {
	id, err := strconv.ParseInt(c.Params(name), 10, 32)
	if err != nil {
		return 0, badRequest("%s must be an integer", name)
	}
	return int32(id), nil
}

func parseBody(c *fiber.Ctx, v any) error {
	if err := c.BodyParser(v); err != nil {
		return badRequest("invalid request body: %v", err)
	}
	return nil
}
