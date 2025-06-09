/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Echo服务器适配器实现 (Echo server adapter implementation)
 */

package echo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	echoMiddleware "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/echo/middleware"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// EchoServer Echo服务器适配器 (Echo server adapter)
// 实现server.WebFramework接口 (Implements server.WebFramework interface)
type EchoServer struct {
	config           *server.ServerConfig         // 服务器配置 (Server configuration)
	serviceContainer services.ServiceContainer    // 服务容器 (Service container)
	echo             *echo.Echo                   // Echo实例 (Echo instance)
	httpServer       *http.Server                 // HTTP服务器 (HTTP server)
	logger           services.Logger              // 日志服务 (Logger service)
}

// NewEchoServer 创建Echo服务器适配器 (Create Echo server adapter)
func NewEchoServer(config *server.ServerConfig, serviceContainer services.ServiceContainer) (*EchoServer, error) {
	// 创建Echo实例 (Create Echo instance)
	e := echo.New()
	
	// 配置Echo实例 (Configure Echo instance)
	e.HideBanner = true
	e.HidePort = true
	
	// 根据模式设置日志级别 (Set log level based on mode)
	if config.IsDebugMode() {
		e.Debug = true
	}

	// 创建服务器实例 (Create server instance)
	server := &EchoServer{
		config:           config,
		serviceContainer: serviceContainer,
		echo:             e,
		logger:           serviceContainer.GetLogger(),
	}

	// 设置基本中间件 (Set basic middleware)
	if err := server.setupMiddleware(); err != nil {
		return nil, serviceContainer.GetErrorHandler().Wrap(err, "failed to setup middleware")
	}

	// 记录服务器创建日志 (Log server creation)
	server.logger.Infow("Echo server adapter created",
		"host", config.Host,
		"port", config.Port,
		"mode", config.Mode,
		"framework", "echo",
	)

	return server, nil
}

// Start 启动服务器 (Start the server)
func (s *EchoServer) Start(ctx context.Context) error {
	// 创建HTTP服务器 (Create HTTP server)
	s.httpServer = &http.Server{
		Addr:           fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		Handler:        s.echo,
		ReadTimeout:    s.config.ReadTimeout,
		WriteTimeout:   s.config.WriteTimeout,
		IdleTimeout:    s.config.IdleTimeout,
		MaxHeaderBytes: s.config.MaxHeaderBytes,
	}

	s.logger.Infow("Starting Echo server",
		"address", s.httpServer.Addr,
		"read_timeout", s.config.ReadTimeout,
		"write_timeout", s.config.WriteTimeout,
	)

	// 启动服务器 (Start server)
	if s.config.TLS.Enabled {
		return s.httpServer.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
	} else {
		return s.httpServer.ListenAndServe()
	}
}

// Stop 停止服务器 (Stop the server)
func (s *EchoServer) Stop(ctx context.Context) error {
	if s.httpServer == nil {
		return nil
	}

	s.logger.Infow("Stopping Echo server", "address", s.httpServer.Addr)

	// 设置关闭超时 (Set shutdown timeout)
	if s.config.GracefulShutdown.Enabled {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.config.GracefulShutdown.Timeout)
		defer cancel()
	}

	// 优雅关闭服务器 (Gracefully shutdown server)
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return s.serviceContainer.GetErrorHandler().Wrap(err, "failed to shutdown Echo server")
	}

	s.logger.Infow("Echo server stopped successfully")
	return nil
}

// RegisterRoute 注册路由 (Register a route)
func (s *EchoServer) RegisterRoute(method, path string, handler server.Handler) error {
	// 将统一处理器转换为Echo处理器 (Convert unified handler to Echo handler)
	echoHandler := s.wrapHandler(handler)
	
	// 注册路由 (Register route)
	s.echo.Add(method, path, echoHandler)
	
	s.logger.Debugw("Route registered",
		"method", method,
		"path", path,
		"framework", "echo",
	)
	
	return nil
}

// RegisterMiddleware 注册全局中间件 (Register global middleware)
func (s *EchoServer) RegisterMiddleware(middleware server.Middleware) error {
	// 将统一中间件转换为Echo中间件 (Convert unified middleware to Echo middleware)
	echoMiddleware := s.wrapMiddleware(middleware)
	
	// 注册中间件 (Register middleware)
	s.echo.Use(echoMiddleware)
	
	s.logger.Debugw("Global middleware registered", "framework", "echo")
	
	return nil
}

// Group 创建路由组 (Create a route group)
func (s *EchoServer) Group(prefix string, middlewares ...server.Middleware) server.RouteGroup {
	// 创建Echo路由组 (Create Echo route group)
	echoGroup := s.echo.Group(prefix)
	
	// 添加中间件到路由组 (Add middleware to route group)
	for _, mw := range middlewares {
		echoMiddleware := s.wrapMiddleware(mw)
		echoGroup.Use(echoMiddleware)
	}
	
	// 创建路由组适配器 (Create route group adapter)
	return NewEchoRouteGroupWithPrefix(echoGroup, s, prefix)
}

// GetNativeEngine 获取原生Echo实例 (Get native Echo instance)
func (s *EchoServer) GetNativeEngine() interface{} {
	return s.echo
}

// GetConfig 获取服务器配置 (Get server configuration)
func (s *EchoServer) GetConfig() *server.ServerConfig {
	return s.config
}

// GetEchoEngine 获取Echo引擎实例 (Get Echo engine instance)
// 提供类型安全的访问方法 (Provides type-safe access method)
func (s *EchoServer) GetEchoEngine() *echo.Echo {
	return s.echo
}

// setupMiddleware 设置基本中间件 (Set up basic middleware)
func (s *EchoServer) setupMiddleware() error {
	// 导入统一中间件包 (Import unified middleware package)
	// 这里我们需要动态导入，因为Go不支持在函数内部导入
	// 所以我们直接使用完整的包路径
	
	// 恢复中间件 (Recovery middleware) - 使用统一实现
	if s.config.Middleware.Recovery.Enabled {
		recoveryConfig := &echoMiddleware.RecoveryConfig{
			Enabled:            s.config.Middleware.Recovery.Enabled,
			PrintStack:         s.config.Middleware.Recovery.PrintStack,
			StackSize:          4 << 10, // 4KB
			DisableStackAll:    false,
			DisableColorOutput: false,
		}
		recoveryMiddleware := echoMiddleware.NewRecoveryMiddleware(recoveryConfig, s.serviceContainer)
		s.echo.Use(recoveryMiddleware.Handler())
	}

	// 日志中间件 (Logger middleware) - 使用统一实现
	if s.config.Middleware.Logger.Enabled {
		loggerConfig := &echoMiddleware.LoggerConfig{
			Enabled:      s.config.Middleware.Logger.Enabled,
			Format:       "${time_rfc3339} ${status} ${method} ${uri} ${latency_human} ${bytes_in}/${bytes_out}\n",
			SkipPaths:    []string{"/health", "/metrics"},
			EnableColors: false,
			TimeFormat:   "2006-01-02T15:04:05Z07:00",
		}
		loggerMiddleware := echoMiddleware.NewLoggerMiddleware(loggerConfig, s.serviceContainer)
		s.echo.Use(loggerMiddleware.Handler())
	}

	// 请求ID中间件 (Request ID middleware) - 保持原生实现
	s.echo.Use(middleware.RequestID())

	// 安全中间件 (Security middleware) - 保持原生实现
	s.echo.Use(middleware.Secure())

	// CORS中间件 (CORS middleware) - 使用统一实现
	if s.config.CORS.Enabled {
		corsConfig := &echoMiddleware.CORSConfig{
			Enabled:          s.config.CORS.Enabled,
			AllowOrigins:     s.config.CORS.AllowOrigins,
			AllowMethods:     s.config.CORS.AllowMethods,
			AllowHeaders:     s.config.CORS.AllowHeaders,
			ExposeHeaders:    s.config.CORS.ExposeHeaders,
			AllowCredentials: s.config.CORS.AllowCredentials,
			MaxAge:           int(s.config.CORS.MaxAge.Seconds()),
		}
		corsMiddleware := echoMiddleware.NewCORSMiddleware(corsConfig, s.serviceContainer)
		s.echo.Use(corsMiddleware.Handler())
	}

	return nil
}

// wrapHandler 包装统一处理器为Echo处理器 (Wrap unified handler to Echo handler)
func (s *EchoServer) wrapHandler(handler server.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 创建统一上下文 (Create unified context)
		ctx := NewEchoContext(c, s.serviceContainer)
		
		// 调用统一处理器 (Call unified handler)
		return handler.Handle(ctx)
	}
}

// wrapMiddleware 包装统一中间件为Echo中间件 (Wrap unified middleware to Echo middleware)
func (s *EchoServer) wrapMiddleware(middleware server.Middleware) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 创建统一上下文 (Create unified context)
			ctx := NewEchoContext(c, s.serviceContainer)
			
			// 定义下一个处理器 (Define next handler)
			nextHandler := func() error {
				return next(c)
			}
			
			// 调用统一中间件 (Call unified middleware)
			return middleware.Process(ctx, nextHandler)
		}
	}
} 