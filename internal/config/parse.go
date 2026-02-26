package config

import (
	"encoding/json"
	"github.com/creasty/defaults"
	"go-admin/pkg/errors"
	"go-admin/pkg/toml"
	"os"
	"path/filepath"
	"sync"
)

var (
	once sync.Once
	C    = new(Config)
)

func MustLoad(dir string, names ...string) {
	once.Do(func() {
		if err := Load(dir, names...); err != nil {
			panic(err)
		}
	})
}

// 从目录中加载各种格式的配置文件，并将它们解析为

// 结构体。
func Load(dir string, names ...string) error {
	// Set default values
	if err := defaults.Set(C); err != nil {
		return err
	}

	supportExts := map[string]struct{}{
		".json": {},
		".toml": {},
	}
	parseFile := func(name string) error {
		ext := filepath.Ext(name)
		if _, ok := supportExts[ext]; !ok {
			return nil
		}

		buf, err := os.ReadFile(name)
		if err != nil {
			return errors.Wrapf(err, "failed to read config file %s", name)
		}

		switch ext {
		case ".json":
			err = json.Unmarshal(buf, C)
		case ".toml":
			err = toml.Unmarshal(buf, C)
		}
		return errors.Wrapf(err, "failed to unmarshal config %s", name)
	}

	for _, name := range names {
		fullname := filepath.Join(dir, name)
		info, err := os.Stat(fullname)
		if err != nil {
			return errors.Wrapf(err, "failed to get config file %s", name)
		}

		if info.IsDir() {
			err := filepath.WalkDir(fullname, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				} else if d.IsDir() {
					return nil
				}
				return parseFile(path)
			})
			if err != nil {
				return errors.Wrapf(err, "failed to walk config dir %s", name)
			}
			continue
		}
		if err := parseFile(fullname); err != nil {
			return err
		}
	}

	return nil
}
