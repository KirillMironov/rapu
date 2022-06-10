package echox

import (
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type grpcError interface {
	GRPCStatus() *status.Status
	Error() string
}

func GRPCErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	var httpError *echo.HTTPError

	switch err := err.(type) {
	case *echo.HTTPError:
		httpError = err
		if err.Internal != nil {
			if internalErr, ok := err.Internal.(*echo.HTTPError); ok {
				httpError = internalErr
			}
		}
	case grpcError:
		st, ok := status.FromError(err)
		if !ok {
			httpError = echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			break
		}

		switch st.Code() {
		case codes.InvalidArgument:
			httpError = echo.NewHTTPError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		case codes.NotFound:
			httpError = echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
		case codes.AlreadyExists:
			httpError = echo.NewHTTPError(http.StatusConflict, http.StatusText(http.StatusConflict))
		case codes.Unauthenticated:
			httpError = echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		default:
			c.Logger().Error(err)
			httpError = echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		}
	default:
		httpError = echo.NewHTTPError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	code := httpError.Code
	message := httpError.Message
	if m, ok := httpError.Message.(string); ok {
		message = echo.Map{"message": m}
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(httpError.Code)
	} else {
		err = c.JSON(code, message)
	}
	if err != nil {
		c.Logger().Error(err)
	}
}
