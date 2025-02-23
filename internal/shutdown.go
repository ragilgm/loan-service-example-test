package internal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo"
	"github.com/typical-go/typical-go/pkg/errkit"
	"go.uber.org/dig"
)

// Shutdown infra
func Shutdown(p struct {
	dig.In
	Pg   *sql.DB
	Echo *echo.Echo
}) error {

	fmt.Printf("Shutdown at %s", time.Now().String())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	errs := errkit.Errors{
		p.Pg.Close(),
		p.Echo.Shutdown(ctx),
	}

	return errs.Unwrap()
}
