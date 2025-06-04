/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Fiber服务器适配器实现 (Fiber server adapter implementation)
 */

package fiber

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/fiber/middleware"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// FiberServer Fiber服务器适配器 (Fiber server adapter)
// 实现server.WebFramework接口 (Implements server.WebFramework interface)
type FiberServer struct {
	config           *server.ServerConfig         // 服务器配置 (Server configuration)
	serviceContainer services.ServiceContainer    // 服务容器 (Service container)
	fiber            *fiber.App                   // Fiber实例 (Fiber instance)
	logger           services.Logger              // 日志服务 (Logger service)
}

// NewFiberServer 创建Fiber服务器适配器 (Create Fiber server adapter)
func NewFiberServer(config *server.ServerConfig, serviceContainer services.ServiceContainer) (*FiberServer, error) {
	// 创建Fiber配置 (Create Fiber configuration)
	fiberConfig := fiber.Config{
		Prefork:       false, // 不使用预分叉模式 (Don't use prefork mode)
		CaseSensitive: true,  // 路径大小写敏感 (Case sensitive paths)
		StrictRouting: false, // 不严格路由 (Not strict routing)
		ServerHeader:  "lmcc-go-sdk/fiber",
		AppName:       "lmcc-go-sdk Fiber Server",
		ReadTimeout:   config.ReadTimeout,
		WriteTimeout:  config.WriteTimeout,
		IdleTimeout:   config.IdleTimeout,
		BodyLimit:     int(config.MaxHeaderBytes),
	}

	// 创建Fiber应用 (Create Fiber app)
	app := fiber.New(fiberConfig)

	// 创建服务器实例 (Create server instance)
	fiberServer := &FiberServer{
		config:           config,
		serviceContainer: serviceContainer,
		fiber:            app,
		logger:           serviceContainer.GetLogger(),
	}

	// 设置中间件 (Setup middleware)
	fiberServer.setupMiddleware()

	// 记录服务器创建 (Log server creation)
	if fiberServer.logger != nil {
		fiberServer.logger.Infow("Fiber server adapter created",
			"host", config.Host,
			"port", config.Port,
			"mode", config.Mode,
			"framework", config.Framework,
		)
	}

	return fiberServer, nil
}

// Start 启动服务器 (Start the server)
func (s *FiberServer) Start(ctx context.Context) error {
	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	
	if s.logger != nil {
		s.logger.Infow("Starting Fiber server",
			"address", address,
			"read_timeout", s.config.ReadTimeout.Seconds(),
			"write_timeout", s.config.WriteTimeout.Seconds(),
		)
	}

	// 在goroutine中启动服务器 (Start server in goroutine)
	errChan := make(chan error, 1)
	go func() {
		if err := s.fiber.Listen(address); err != nil {
			errChan <- err
		}
	}()

	// 等待上下文取消或启动错误 (Wait for context cancellation or startup error)
	select {
	case <-ctx.Done():
		return s.Stop(context.Background())
	case err := <-errChan:
		return err
	}
}

// Stop 停止服务器 (Stop the server)
func (s *FiberServer) Stop(ctx context.Context) error {
	if s.fiber == nil {
		return nil
	}

	if s.logger != nil {
		s.logger.Infow("Stopping Fiber server",
			"address", fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		)
	}

	// 使用Fiber的Shutdown方法 (Use Fiber's Shutdown method)
	if err := s.fiber.Shutdown(); err != nil {
		if s.logger != nil {
			s.logger.Errorw("Failed to stop Fiber server", "error", err)
		}
		return err
	}

	if s.logger != nil {
		s.logger.Info("Fiber server stopped successfully")
	}
	return nil
}

// RegisterRoute 注册路由 (Register a route)
func (s *FiberServer) RegisterRoute(method, path string, handler server.Handler) error {
	// 包装处理器 (Wrap handler)
	fiberHandler := s.wrapHandler(handler)
	
	// 注册路由 (Register route)
	s.fiber.Add(method, path, fiberHandler)
	
	return nil
}

// RegisterMiddleware 注册全局中间件 (Register global middleware)
func (s *FiberServer) RegisterMiddleware(middleware server.Middleware) error {
	// 包装中间件 (Wrap middleware)
	fiberMiddleware := s.wrapMiddleware(middleware)
	
	// 注册中间件 (Register middleware)
	s.fiber.Use(fiberMiddleware)
	
	return nil
}

// Group 创建路由组 (Create a route group)
func (s *FiberServer) Group(prefix string, middlewares ...server.Middleware) server.RouteGroup {
	// 创建Fiber路由组 (Create Fiber route group)
	fiberGroup := s.fiber.Group(prefix)
	
	// 添加中间件到路由组 (Add middleware to route group)
	for _, mw := range middlewares {
		fiberMiddleware := s.wrapMiddleware(mw)
		fiberGroup.Use(fiberMiddleware)
	}
	
	// 创建路由组适配器 (Create route group adapter)
	return NewFiberRouteGroupWithPrefix(fiberGroup, prefix, s.serviceContainer)
}

// GetNativeEngine 获取原生框架实例 (Get native framework instance)
func (s *FiberServer) GetNativeEngine() interface{} {
	return s.fiber
}

// GetConfig 获取服务器配置 (Get server configuration)
func (s *FiberServer) GetConfig() *server.ServerConfig {
	return s.config
}

// GetFiberApp 获取Fiber应用实例 (Get Fiber app instance)
func (s *FiberServer) GetFiberApp() *fiber.App {
	return s.fiber
}

// wrapHandler 包装处理器 (Wrap handler)
func (s *FiberServer) wrapHandler(handler server.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 创建上下文适配器 (Create context adapter)
		ctx := NewFiberContext(c, s.serviceContainer)
		
		// 调用处理器 (Call handler)
		return handler.Handle(ctx)
	}
}

// wrapMiddleware 包装中间件 (Wrap middleware)
func (s *FiberServer) wrapMiddleware(middleware server.Middleware) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 创建上下文适配器 (Create context adapter)
		ctx := NewFiberContext(c, s.serviceContainer)
		
		// 调用中间件 (Call middleware)
		return middleware.Process(ctx, func() error {
			return c.Next()
		})
	}
}

// setupMiddleware 设置中间件 (Setup middleware)
func (s *FiberServer) setupMiddleware() {
	// 设置恢复中间件 (Setup recovery middleware)
	if s.config.Middleware.Recovery.Enabled {
		recoveryMiddleware := middleware.NewRecoveryMiddleware(&middleware.RecoveryConfig{
			PrintStack: s.config.Middleware.Recovery.PrintStack,
		}, s.serviceContainer)
		s.fiber.Use(recoveryMiddleware.Handler())
	}

	// 设置日志中间件 (Setup logger middleware)
	if s.config.Middleware.Logger.Enabled {
		loggerMiddleware := middleware.NewLoggerMiddleware(&middleware.LoggerConfig{
			Format:    "[${time}] ${status} - ${method} ${path} ${latency}\n",
			TimeZone:  "Local",
			TimeFormat: "15:04:05",
		}, s.serviceContainer)
		s.fiber.Use(loggerMiddleware.Handler())
	}

	// 设置CORS中间件 (Setup CORS middleware)
	if s.config.CORS.Enabled {
		corsConfig := &middleware.CORSConfig{
			AllowOrigins:     s.config.CORS.AllowOrigins,
			AllowMethods:     s.config.CORS.AllowMethods,
			AllowHeaders:     s.config.CORS.AllowHeaders,
			AllowCredentials: s.config.CORS.AllowCredentials,
			ExposeHeaders:    s.config.CORS.ExposeHeaders,
			MaxAge:           int(s.config.CORS.MaxAge.Seconds()),
		}
		corsMiddleware := middleware.NewCORSMiddleware(corsConfig, s.serviceContainer)
		s.fiber.Use(corsMiddleware.Handler())
	}
} 