/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: Gin服务器适配器实现
 */

package gin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server"
	ginMiddleware "github.com/lmcc-dev/lmcc-go-sdk/pkg/server/plugins/gin/middleware"
	"github.com/lmcc-dev/lmcc-go-sdk/pkg/server/services"
)

// GinServer Gin服务器适配器 (Gin server adapter)
// 将Gin引擎适配到统一的WebFramework接口 (Adapts Gin engine to unified WebFramework interface)
type GinServer struct {
	engine     *gin.Engine
	config     *server.ServerConfig
	httpServer *http.Server
	routes     map[string]*GinRouteGroup
	services   services.ServiceContainer
}

// NewGinServer 创建Gin服务器适配器 (Create Gin server adapter)
// 保留向后兼容性 (Maintain backward compatibility)
func NewGinServer(config *server.ServerConfig) *GinServer {
	return NewGinServerWithServices(config, nil)
}

// NewGinServerWithServices 创建带服务容器的Gin服务器适配器 (Create Gin server adapter with service container)
func NewGinServerWithServices(config *server.ServerConfig, serviceContainer services.ServiceContainer) *GinServer {
	// 如果没有提供服务容器，创建默认的 (If no service container provided, create default one)
	if serviceContainer == nil {
		serviceContainer = services.NewServiceContainerWithDefaults()
	}
	
	// 设置Gin模式 (Set Gin mode)
	switch config.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	
	// 创建Gin引擎 (Create Gin engine)
	engine := gin.New()
	
	// 应用Gin特定配置 (Apply Gin-specific configuration)
	if ginConfig, ok := config.Plugins["gin"].(map[string]interface{}); ok {
		applyGinConfig(engine, ginConfig)
	}
	
	// 创建HTTP服务器 (Create HTTP server)
	httpServer := &http.Server{
		Addr:           config.GetAddress(),
		Handler:        engine,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		IdleTimeout:    config.IdleTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	
	ginServer := &GinServer{
		engine:     engine,
		config:     config,
		httpServer: httpServer,
		routes:     make(map[string]*GinRouteGroup),
		services:   serviceContainer,
	}
	
	// 设置中间件 (Setup middleware)
	ginServer.setupMiddleware()
	
	return ginServer
}

// Start 启动服务器 (Start server)
func (s *GinServer) Start(ctx context.Context) error {
	logger := s.services.GetLogger()
	
	// 启动HTTP服务器 (Start HTTP server)
	go func() {
		var err error
		if s.config.TLS.Enabled {
			if s.config.TLS.CertFile != "" && s.config.TLS.KeyFile != "" {
				logger.Infof("Starting HTTPS server on %s", s.config.GetAddress())
				err = s.httpServer.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
			} else {
				err = fmt.Errorf("TLS enabled but cert file or key file not provided")
			}
		} else {
			logger.Infof("Starting HTTP server on %s", s.config.GetAddress())
			err = s.httpServer.ListenAndServe()
		}
		
		if err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start server: %v", err)
		}
	}()
	
	return nil
}

// Stop 停止服务器 (Stop server)
func (s *GinServer) Stop(ctx context.Context) error {
	logger := s.services.GetLogger()
	logger.Info("Stopping server...")
	
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		logger.Errorf("Error stopping server: %v", err)
	} else {
		logger.Info("Server stopped successfully")
	}
	
	return err
}

// RegisterRoute 注册路由 (Register route)
func (s *GinServer) RegisterRoute(method, path string, handler server.Handler) error {
	// 将统一Handler适配为Gin Handler (Adapt unified Handler to Gin Handler)
	ginHandler := s.adaptHandler(handler)
	
	// 注册路由 (Register route)
	s.engine.Handle(method, path, ginHandler)
	
	return nil
}

// RegisterMiddleware 注册全局中间件 (Register global middleware)
func (s *GinServer) RegisterMiddleware(middleware server.Middleware) error {
	// 将统一Middleware适配为Gin Middleware (Adapt unified Middleware to Gin Middleware)
	ginMiddleware := s.adaptMiddleware(middleware)
	
	// 注册中间件 (Register middleware)
	s.engine.Use(ginMiddleware)
	
	return nil
}

// Group 创建路由组 (Create route group)
func (s *GinServer) Group(prefix string, middlewares ...server.Middleware) server.RouteGroup {
	// 创建Gin路由组 (Create Gin route group)
	ginGroup := s.engine.Group(prefix)
	
	// 应用中间件 (Apply middlewares)
	for _, middleware := range middlewares {
		ginMiddleware := s.adaptMiddleware(middleware)
		ginGroup.Use(ginMiddleware)
	}
	
	// 创建路由组适配器 (Create route group adapter)
	routeGroup := NewGinRouteGroupWithServices(ginGroup, prefix, s, s.services)
	
	// 存储路由组 (Store route group)
	s.routes[prefix] = routeGroup
	
	return routeGroup
}

// GetNativeEngine 获取原生Gin引擎 (Get native Gin engine)
func (s *GinServer) GetNativeEngine() interface{} {
	return s.engine
}

// GetConfig 获取服务器配置 (Get server configuration)
func (s *GinServer) GetConfig() *server.ServerConfig {
	return s.config
}

// GetServices 获取服务容器 (Get service container)
func (s *GinServer) GetServices() services.ServiceContainer {
	return s.services
}

// GetGinEngine 获取Gin引擎 (Get Gin engine)
// 类型安全的方法获取Gin引擎 (Type-safe method to get Gin engine)
func (s *GinServer) GetGinEngine() *gin.Engine {
	return s.engine
}

// GetHTTPServer 获取HTTP服务器 (Get HTTP server)
func (s *GinServer) GetHTTPServer() *http.Server {
	return s.httpServer
}

// adaptHandler 将统一Handler适配为Gin Handler (Adapt unified Handler to Gin Handler)
func (s *GinServer) adaptHandler(handler server.Handler) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		// 创建统一上下文 (Create unified context)
		ctx := NewGinContext(ginCtx)
		
		// 调用处理器 (Call handler)
		if err := handler.Handle(ctx); err != nil {
			// 使用服务容器的错误处理器 (Use service container's error handler)
			errorHandler := s.services.GetErrorHandler()
			logger := s.services.GetLogger()
			
			// 记录错误 (Log error)
			logger.Errorf("Handler error: %v", err)
			
			// 处理错误 (Handle error)
			ginCtx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
				"message": errorHandler.GetStackTrace(err),
			})
			ginCtx.Abort()
		}
	}
}

// adaptMiddleware 将统一Middleware适配为Gin Middleware (Adapt unified Middleware to Gin Middleware)
func (s *GinServer) adaptMiddleware(middleware server.Middleware) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		// 创建统一上下文 (Create unified context)
		ctx := NewGinContext(ginCtx)
		
		// 调用中间件 (Call middleware)
		err := middleware.Process(ctx, func() error {
			ginCtx.Next()
			return nil
		})
		
		if err != nil {
			// 使用服务容器的错误处理器 (Use service container's error handler)
			errorHandler := s.services.GetErrorHandler()
			logger := s.services.GetLogger()
			
			// 记录错误 (Log error)
			logger.Errorf("Middleware error: %v", err)
			
			// 处理错误 (Handle error)
			ginCtx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
				"message": errorHandler.GetStackTrace(err),
			})
			ginCtx.Abort()
		}
	}
}

// setupMiddleware 设置中间件 (Setup middleware)
func (s *GinServer) setupMiddleware() {
	// 设置恢复中间件 (Setup recovery middleware) - 使用原生实现
	if s.config.Middleware.Recovery.Enabled {
		if s.config.Middleware.Recovery.PrintStack {
			s.engine.Use(gin.RecoveryWithWriter(gin.DefaultWriter))
		} else {
			s.engine.Use(gin.Recovery())
		}
	}

	// 设置日志中间件 (Setup logger middleware) - 使用原生实现
	if s.config.Middleware.Logger.Enabled {
		s.engine.Use(gin.Logger())
	}

	// 设置CORS中间件 (Setup CORS middleware) - 使用统一实现
	if s.config.CORS.Enabled {
		corsMiddleware := ginMiddleware.NewCORSMiddleware(&s.config.CORS)
		s.engine.Use(s.adaptMiddleware(corsMiddleware))
	}
}



// applyGinConfig 应用Gin特定配置 (Apply Gin-specific configuration)
func applyGinConfig(engine *gin.Engine, config map[string]interface{}) {
	// 设置信任的代理 (Set trusted proxies)
	if trustedProxies, ok := config["trusted_proxies"].([]string); ok && len(trustedProxies) > 0 {
		engine.SetTrustedProxies(trustedProxies)
	}
	
	// 设置是否重定向尾部斜杠 (Set redirect trailing slash)
	if redirectTrailingSlash, ok := config["redirect_trailing_slash"].(bool); ok {
		engine.RedirectTrailingSlash = redirectTrailingSlash
	}
	
	// 设置是否重定向固定路径 (Set redirect fixed path)
	if redirectFixedPath, ok := config["redirect_fixed_path"].(bool); ok {
		engine.RedirectFixedPath = redirectFixedPath
	}
	
	// 设置是否处理方法不允许 (Set handle method not allowed)
	if handleMethodNotAllowed, ok := config["handle_method_not_allowed"].(bool); ok {
		engine.HandleMethodNotAllowed = handleMethodNotAllowed
	}
	
	// 设置最大多部分内存 (Set max multipart memory)
	if maxMultipartMemory, ok := config["max_multipart_memory"].(int64); ok {
		engine.MaxMultipartMemory = maxMultipartMemory
	}
	
	// 设置HTML模板 (Set HTML templates)
	if templatePattern, ok := config["template_pattern"].(string); ok && templatePattern != "" {
		engine.LoadHTMLGlob(templatePattern)
	}
	
	// 设置静态文件路径 (Set static file path)
	if staticConfig, ok := config["static"].(map[string]interface{}); ok {
		if relativePath, ok := staticConfig["relative_path"].(string); ok {
			if root, ok := staticConfig["root"].(string); ok {
				engine.Static(relativePath, root)
			}
		}
	}
}