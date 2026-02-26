package bootstrap

import (
	"context"
	"crypto/tls"
	"go-admin/internal/config"
	"go-admin/internal/wirex"
	"go-admin/pkg/errors"
	"go-admin/pkg/logging"
	"go-admin/pkg/middleware"
	"go-admin/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func startHTTPServer(ctx context.Context, injector *wirex.Injector) (func(), error) {
	if config.C.IsDebug() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()

	e.GET("/health", func(c *gin.Context) {
		util.ResOK(c)
	})

	e.NoMethod(func(c *gin.Context) {
		util.ResError(c, errors.MethodNotAllowed("", "Method Not Allowed"))
	})
	e.NoRoute(func(c *gin.Context) {
		util.ResError(c, errors.NotFound("", "Not Found"))
	})

	// Register middlewares
	useHTTPMiddlewares(e)

	// Register routers
	if err := injector.M.RegisterRouters(ctx, e); err != nil {
		return nil, err
	}

	addr := config.C.General.HTTP.Addr

	logging.Context(ctx).Info("HTTP server is listening", zap.String("addr", addr))

	srv := &http.Server{
		Addr:              addr,
		Handler:           e,
		ReadTimeout:       time.Second * time.Duration(config.C.General.HTTP.ReadTimeout),
		ReadHeaderTimeout: time.Second * time.Duration(config.C.General.HTTP.ReadTimeout),
		WriteTimeout:      time.Second * time.Duration(config.C.General.HTTP.WriteTimeout),
		IdleTimeout:       time.Second * time.Duration(config.C.General.HTTP.IdleTimeout),
	}

	go func() {
		var err error
		//如果配置了证书和私钥就用https否则就用http
		if config.C.General.HTTP.CertFile != "" && config.C.General.HTTP.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(config.C.General.HTTP.CertFile, config.C.General.HTTP.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			logging.Context(ctx).Error("Failed to listen http server", zap.Error(err))
		}
	}()
	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.C.General.HTTP.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logging.Context(ctx).Error("Failed to shutdown http server", zap.Error(err))
		}
	}, nil
}

func useHTTPMiddlewares(e *gin.Engine) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Enable:                 config.C.Middleware.CORS.Enable,
		AllowAllOrigins:        config.C.Middleware.CORS.AllowAllOrigins,
		AllowOrigins:           config.C.Middleware.CORS.AllowOrigins,
		AllowMethods:           config.C.Middleware.CORS.AllowMethods,
		AllowHeaders:           config.C.Middleware.CORS.AllowHeaders,
		AllowCredentials:       config.C.Middleware.CORS.AllowCredentials,
		ExposeHeaders:          config.C.Middleware.CORS.ExposeHeaders,
		MaxAge:                 config.C.Middleware.CORS.MaxAge,
		AllowWildcard:          config.C.Middleware.CORS.AllowWildcard,
		AllowBrowserExtensions: config.C.Middleware.CORS.AllowBrowserExtensions,
		AllowWebSockets:        config.C.Middleware.CORS.AllowWebSockets,
		AllowFiles:             config.C.Middleware.CORS.AllowFiles,
	}))
}
