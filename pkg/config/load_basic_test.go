/*
 * Author: Martin <lmccc.dev@gmail.com>
 * Co-Author: AI Assistant
 * Description: This code was collaboratively developed by Martin and AI Assistant.
 * Contains basic configuration loading tests (YAML, TOML).
 */

package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_YAML_Valid(t *testing.T) {
	yamlContent := `
server:
  host: "127.0.0.1"
  port: 9999
  mode: "debug"
  readTimeout: 10s
  writeTimeout: 15s
  gracefulShutdownTimeout: 20s
log:
  level: "debug"
  format: "json"
  output: "file"
  filename: "test.log"
  maxSize: 200
  maxBackups: 10
  maxAge: 14
  compress: true
database:
  type: "postgres"
  host: "db.example.com"
  port: 5432
  user: "testuser"
  password: "testpassword"
  dbName: "testdb"
  maxIdleConns: 5
  maxOpenConns: 50
  connMaxLifetime: 30m
tracing:
  enabled: true
  provider: "jaeger"
  endpoint: "http://jaeger:14268/api/traces"
  samplerType: "probabilistic"
  samplerParam: 0.5
metrics:
  enabled: true
  provider: "prometheus"
  port: 9191
  path: "/testmetrics"
customFeature:
  apiKey: "test-api-key"
  rateLimit: 100
  enabled: true
`
	configFile, cleanup := createTempConfigFile(t, yamlContent, "yaml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "yaml"))
	require.NoError(t, err, "LoadConfig should not return an error for valid YAML")

	// Assert Server Config
	require.NotNil(t, loadedCfg.Server, "Server config should be loaded")
	assert.Equal(t, "127.0.0.1", loadedCfg.Server.Host)
	assert.Equal(t, 9999, loadedCfg.Server.Port)
	assert.Equal(t, "debug", loadedCfg.Server.Mode)
	assert.Equal(t, 10*time.Second, loadedCfg.Server.ReadTimeout)
	assert.Equal(t, 15*time.Second, loadedCfg.Server.WriteTimeout)
	assert.Equal(t, 20*time.Second, loadedCfg.Server.GracefulShutdownTimeout)

	// Assert Log Config
	require.NotNil(t, loadedCfg.Log, "Log config should be loaded")
	assert.Equal(t, "debug", loadedCfg.Log.Level)
	assert.Equal(t, "json", loadedCfg.Log.Format)
	assert.Equal(t, "file", loadedCfg.Log.Output)
	assert.Equal(t, "test.log", loadedCfg.Log.Filename)
	assert.Equal(t, 200, loadedCfg.Log.MaxSize)
	assert.Equal(t, 10, loadedCfg.Log.MaxBackups)
	assert.Equal(t, 14, loadedCfg.Log.MaxAge)
	assert.True(t, loadedCfg.Log.Compress)

	// Assert Database Config
	require.NotNil(t, loadedCfg.Database, "Database config should be loaded")
	assert.Equal(t, "postgres", loadedCfg.Database.Type)
	assert.Equal(t, "db.example.com", loadedCfg.Database.Host)
	assert.Equal(t, 5432, loadedCfg.Database.Port)
	assert.Equal(t, "testuser", loadedCfg.Database.User)
	assert.Equal(t, "testpassword", loadedCfg.Database.Password)
	assert.Equal(t, "testdb", loadedCfg.Database.DBName)
	assert.Equal(t, 5, loadedCfg.Database.MaxIdleConns)
	assert.Equal(t, 50, loadedCfg.Database.MaxOpenConns)
	assert.Equal(t, 30*time.Minute, loadedCfg.Database.ConnMaxLifetime)

	// Assert Tracing Config
	require.NotNil(t, loadedCfg.Tracing, "Tracing config should be loaded")
	assert.True(t, loadedCfg.Tracing.Enabled)
	assert.Equal(t, "jaeger", loadedCfg.Tracing.Provider)
	assert.Equal(t, "http://jaeger:14268/api/traces", loadedCfg.Tracing.Endpoint)
	assert.Equal(t, "probabilistic", loadedCfg.Tracing.SamplerType)
	assert.Equal(t, 0.5, loadedCfg.Tracing.SamplerParam)

	// Assert Metrics Config
	require.NotNil(t, loadedCfg.Metrics, "Metrics config should be loaded")
	assert.True(t, loadedCfg.Metrics.Enabled)
	assert.Equal(t, "prometheus", loadedCfg.Metrics.Provider)
	assert.Equal(t, 9191, loadedCfg.Metrics.Port)
	assert.Equal(t, "/testmetrics", loadedCfg.Metrics.Path)

	// Assert Custom Feature Config
	require.NotNil(t, loadedCfg.CustomFeature, "CustomFeature config should be loaded")
	assert.Equal(t, "test-api-key", loadedCfg.CustomFeature.APIKey)
	assert.Equal(t, 100, loadedCfg.CustomFeature.RateLimit)
	assert.True(t, loadedCfg.CustomFeature.Enabled)
}

func TestLoadConfig_TOML_Valid(t *testing.T) {
	tomlContent := `
[server]
host = "192.168.1.1"
port = 8888
mode = "test"
readTimeout = "8s"
writeTimeout = "12s"
gracefulShutdownTimeout = "18s"

[log]
level = "warn"
format = "text"
output = "stderr"

[database]
type = "mysql"
host = "mysql.internal"
port = 3307
user = "tomluser"
password = "tomlpass"
dbName = "tomldb"
maxIdleConns = 20
maxOpenConns = 150
connMaxLifetime = "2h"

[tracing]
enabled = true
provider = "zipkin"
endpoint = "http://zipkin:9411/api/v2/spans"
samplerType = "const"
samplerParam = 1.0

[metrics]
enabled = true
provider = "prometheus"
port = 9292
path = "/tomlmetrics"

[customFeature]
apiKey = "toml-api-key"
rateLimit = 50
enabled = false
`
	configFile, cleanup := createTempConfigFile(t, tomlContent, "toml")
	defer cleanup()

	var loadedCfg testAppConfig
	initializeTestConfig(&loadedCfg) // Initialize before loading

	err := LoadConfig(&loadedCfg, WithConfigFile(configFile, "toml"))
	require.NoError(t, err, "LoadConfig should not return an error for valid TOML")

	// Assert Server Config
	require.NotNil(t, loadedCfg.Server, "Server config should be loaded")
	assert.Equal(t, "192.168.1.1", loadedCfg.Server.Host)
	assert.Equal(t, 8888, loadedCfg.Server.Port)
	assert.Equal(t, "test", loadedCfg.Server.Mode)
	assert.Equal(t, 8*time.Second, loadedCfg.Server.ReadTimeout) // Viper reads TOML duration strings
	assert.Equal(t, 12*time.Second, loadedCfg.Server.WriteTimeout)
	assert.Equal(t, 18*time.Second, loadedCfg.Server.GracefulShutdownTimeout)

	// Assert Log Config (Only specified fields)
	require.NotNil(t, loadedCfg.Log, "Log config should be loaded")
	assert.Equal(t, "warn", loadedCfg.Log.Level)
	assert.Equal(t, "text", loadedCfg.Log.Format)
	assert.Equal(t, "stderr", loadedCfg.Log.Output)
	// Check defaults for unspecified fields
	assert.Equal(t, "app.log", loadedCfg.Log.Filename) // Default
	assert.Equal(t, 100, loadedCfg.Log.MaxSize)        // Default

	// Assert Database Config
	require.NotNil(t, loadedCfg.Database, "Database config should be loaded")
	assert.Equal(t, "mysql", loadedCfg.Database.Type)
	assert.Equal(t, "mysql.internal", loadedCfg.Database.Host)
	assert.Equal(t, 3307, loadedCfg.Database.Port)
	assert.Equal(t, "tomluser", loadedCfg.Database.User)
	assert.Equal(t, "tomlpass", loadedCfg.Database.Password)
	assert.Equal(t, "tomldb", loadedCfg.Database.DBName)
	assert.Equal(t, 20, loadedCfg.Database.MaxIdleConns)
	assert.Equal(t, 150, loadedCfg.Database.MaxOpenConns)
	assert.Equal(t, 2*time.Hour, loadedCfg.Database.ConnMaxLifetime) // Viper reads TOML duration strings

	// Assert Tracing Config
	require.NotNil(t, loadedCfg.Tracing, "Tracing config should be loaded")
	assert.True(t, loadedCfg.Tracing.Enabled)
	assert.Equal(t, "zipkin", loadedCfg.Tracing.Provider)
	assert.Equal(t, "http://zipkin:9411/api/v2/spans", loadedCfg.Tracing.Endpoint)
	assert.Equal(t, "const", loadedCfg.Tracing.SamplerType)
	assert.Equal(t, 1.0, loadedCfg.Tracing.SamplerParam)

	// Assert Metrics Config
	require.NotNil(t, loadedCfg.Metrics, "Metrics config should be loaded")
	assert.True(t, loadedCfg.Metrics.Enabled)
	assert.Equal(t, "prometheus", loadedCfg.Metrics.Provider)
	assert.Equal(t, 9292, loadedCfg.Metrics.Port)
	assert.Equal(t, "/tomlmetrics", loadedCfg.Metrics.Path)

	// Assert Custom Feature Config
	require.NotNil(t, loadedCfg.CustomFeature, "CustomFeature config should be loaded")
	assert.Equal(t, "toml-api-key", loadedCfg.CustomFeature.APIKey)
	assert.Equal(t, 50, loadedCfg.CustomFeature.RateLimit)
	assert.False(t, loadedCfg.CustomFeature.Enabled)
} 