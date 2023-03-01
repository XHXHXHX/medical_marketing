package conf

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/imdario/mergo"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// Builder 用于构建 Configuration
type Builder struct {
	configFile  string // 目前只支持 yaml
	overrideENV bool
}

// FromYaml 从 yaml 文件构建
func (b *Builder) FromYaml(file string) *Builder {
	b.configFile = file
	return b
}

// WithENV 使用环境变量覆盖值
func (b *Builder) WithENV() *Builder {
	b.overrideENV = true
	return b
}

// Build 构建 Configuration
// 构建顺序如下:
// 1. 从配置文件 .yaml
// 2. 使用环境变量覆盖对应变量 (Optioanl)
func (b *Builder) Build(cfg interface{}) error {
	if cfg == nil {
		return errors.New("params config is nil")
	}
	if reflect.ValueOf(cfg).Type().Kind() != reflect.Ptr {
		return errors.New("params config should be a pointer")
	}

	unmarshal := func(filepath string, dst interface{}) error {
		content, err := ioutil.ReadFile(filepath)
		if err != nil {
			return errors.Wrapf(err, "read config file: %s err", filepath)
		}

		if err := yaml.Unmarshal(content, dst); err != nil {
			return errors.Wrap(err, "yaml unmarshal err")
		}
		return nil
	}

	// 1. 从 .yaml 中构建
	// 1.1 default.yaml
	if defaultFile, err := b.fullConfigPath(`default.yaml`); err == nil {
		if err := unmarshal(defaultFile, cfg); err != nil {
			return err
		}
	}

	// 1.2 load from config file (i.e. dev.yaml/staging.yaml...)
	// TODO 暂时要求必须存在 配置文件
	if len(b.configFile) == 0 {
		return errors.New("config file is not specified")
	}
	configFile, err := b.fullConfigPath(b.configFile)

	if err != nil {
		return err
	}
	tmpCfg := reflect.New(reflect.ValueOf(cfg).Elem().Type()).Interface()
	if err := unmarshal(configFile, tmpCfg); err != nil {
		return err
	}

	if err := mergo.Merge(cfg, tmpCfg, mergo.WithOverride); err != nil {
		return errors.Wrap(err, "merge configuration err")
	}

	// 2. 使用环境变量覆盖
	if b.overrideENV {
		if err := envconfig.Process("", cfg); err != nil {
			return errors.Wrap(err, "override ENV err")
		}
	}
	return nil
}

// TODO 增加方法可以从更多的 path 搜索,现在默认从 ./config 和 . 中寻找
func (b *Builder) fullConfigPath(filename string) (string, error) {
	for _, folder := range []string{"./config", ".", "../../config", "../config"} {
		fullpath := path.Join(folder, filename)
		if _, err := os.Stat(fullpath); err == nil {
			return fullpath, nil
		}
	}
	return "", errors.Errorf("file: %s not found", filename)
}

