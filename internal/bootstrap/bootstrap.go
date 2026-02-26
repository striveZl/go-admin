package bootstrap

import (
	"context"
	"fmt"
	"go-admin/internal/config"
	"go-admin/internal/wirex"
	"go-admin/pkg/logging"
	"go-admin/pkg/util"
	"os"
	"strings"

	"go.uber.org/zap"
)

type RunConfig struct {
	WorkDir   string //工作目录
	Configs   string //目录或文件
	StaticDir string //静态文件目录
}

// Run 函数初始化并启动一个服务，包括配置和日志记录，并处理

// 退出时的清理工作

func Run(ctx context.Context, runCfg RunConfig) error {
	defer func() {
		if err := zap.L().Sync(); err != nil {
			fmt.Printf("zap日志记录同步器失败:%s \n", err.Error())
		}
	}()

	//加载配置
	workDir := runCfg.WorkDir
	staticDir := runCfg.StaticDir

	configNames := parseConfigNames(runCfg.Configs)
	if len(configNames) == 0 {
		return fmt.Errorf("no config names provided")
	}

	config.MustLoad(workDir, configNames...)
	config.C.General.WorkDir = workDir

	config.C.Middleware.Static.Dir = staticDir

	config.C.Print()

	if err := logging.InitWithConfig(ctx, &config.C.Logger); err != nil {
		return err
	}

	ctx = logging.NewTag(ctx, logging.TagKeyMain)

	logging.Context(ctx).Info("服务器启动...",
		zap.String("version", config.C.General.Version),
		zap.Int("pid", os.Getpid()),
		zap.String("workdir", workDir),
		zap.String("config", runCfg.Configs),
		zap.String("static", staticDir),
	)

	injector, cleanInjectorFn, err := wirex.BuildInjector(ctx)
	if err != nil {
		return err
	}

	// Initialize global prometheus metrics.
	// prom.Init()

	return util.Run(ctx, func(ctx context.Context) (func(), error) {
		httpServerCleanFn, err := startHTTPServer(ctx, injector)
		if err != nil {
			return cleanInjectorFn, err
		}

		return func() {
			httpServerCleanFn()

			if err := injector.M.Release(ctx); err != nil {
				logging.Context(ctx).Error("failed to release injector", zap.Error(err))
			}

			cleanInjectorFn()
		}, nil
	})

}

func parseConfigNames(raw string) []string {
	items := strings.Split(raw, ",")
	names := make([]string, 0, len(items))
	for _, item := range items {
		name := strings.TrimSpace(item)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}
