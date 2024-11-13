package config

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"wallet-service/pkg/logger"
)

type Cfg[T any] struct {
	mu    sync.Mutex
	vp    *viper.Viper
	value T
}

func NewCfg[T any]() *Cfg[T] {
	return &Cfg[T]{
		vp:    viper.New(),
		value: newT[T](),
	}
}

func newT[T any]() (t T) {
	return
}

func (c *Cfg[T]) StartWatch(path string) error {
	dir, name, ext := filepath.Dir(path), strings.ReplaceAll(filepath.Base(path), filepath.Ext(path), ""), strings.Trim(filepath.Ext(path), ".")

	var vp = c.vp
	vp.SetConfigName(name)
	vp.SetConfigType(ext)
	vp.AddConfigPath(dir)
	err := vp.ReadInConfig()
	if err != nil {
		return err
	}

	var update = func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		err := vp.Unmarshal(&c.value)
		if err != nil {
			logger.Error(context.Background(), "unmarshal to update config failed", zap.Error(err))
			return
		}
		file := vp.ConfigFileUsed()
		content, err := os.ReadFile(file)
		if err != nil {
			logger.Error(context.Background(), "read config file failed", zap.Error(err))
			return
		}
		logger.Info(context.Background(), "update config",
			zap.String("content", string(content)),
			zap.Any("obj", c.value),
		)
	}

	vp.OnConfigChange(func(e fsnotify.Event) {
		update()
	})
	vp.WatchConfig()

	update()

	return nil
}

func (c *Cfg[T]) Get() T {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func MustUseFileConfig[T any](path string) *Cfg[T] {
	var c = NewCfg[T]()
	err := c.StartWatch(path)
	if err != nil {
		logger.Panic(context.Background(), "start watch config failed", zap.Error(err))
	}
	return c
}

func GetConfig() *ServerConfig {
	var v = fileCfg.Get()
	return &v
}

var fileCfg = NewCfg[ServerConfig]()

// Init init config
func Init(configPath ...string) {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	pathConfigPath := filepath.Join(path, DefaultPathConfigPath)
	if len(configPath) > 0 {
		pathConfigPath = configPath[0]
	}
	fileCfg = MustUseFileConfig[ServerConfig](pathConfigPath)
}

const DefaultPathConfigPath = "config/config.yml"
