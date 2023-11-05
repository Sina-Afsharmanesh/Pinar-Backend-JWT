package httpLocal

import (
	"jwt/log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

var (
	HttpServer = createHttpServer()
)

func createHttpServer() *echo.Echo {
	server := echo.New()
	server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency:       true,
		LogProtocol:      false,
		LogRemoteIP:      true,
		LogHost:          false,
		LogMethod:        true,
		LogURI:           true,
		LogURIPath:       false,
		LogRoutePath:     true,
		LogRequestID:     false,
		LogReferer:       false,
		LogUserAgent:     true,
		LogStatus:        true,
		LogError:         true,
		LogContentLength: false,
		LogResponseSize:  false,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.ServLogger.Info("REQUESR", zap.String("IP", v.RemoteIP), zap.String("URI", v.URI), zap.String("route", v.RoutePath), zap.String("method", v.Method), zap.Int("status", v.Status), zap.String("user agent", v.UserAgent), zap.Duration("delay", v.Latency), zap.NamedError("Error", v.Error))

			return nil
		},
	}))
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: middleware.DefaultCORSConfig.AllowHeaders,
	}))
	return server
}
