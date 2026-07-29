package gin_support

import (
	"GameManager/adapters/in/utils"
	"github.com/knadh/koanf/parsers/yaml" // yaml 解析器
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"go.uber.org/zap"
	"sync"
)

type Config struct {
	Server struct {
		Port int    `koanf:"port"`
		Host string `koanf:"host"`
	} `koanf:"server"`
}

var config *Config
var once = sync.Once{}

func initConfig() {
	log := utils.GetLogger()
	k := koanf.New(".")

	// 加载默认值（使用 confmap.Provider）
	defaults := map[string]interface{}{
		"server.port": 10086,
		"server.host": "127.0.0.1",
	}

	// delim = "." 表示 map 的 key 是扁平的（dotted notation）
	if err := k.Load(confmap.Provider(defaults, "."), nil); err != nil {
		log.Fatal("加载默认值失败", zap.Error(err))
	}

	// 再加载配置文件，没有就用默认值
	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		log.Error("配置文件加载失败（可能不存在），使用默认值", zap.Error(err))
	}

	var cfg = Config{}
	// 反序列化到结构体
	if err := k.Unmarshal("", &cfg); err != nil {
		log.Fatal("结构体反序列化失败", zap.Error(err))
	}
	config = &cfg
}

func GetConfig() *Config {
	if config == nil {
		once.Do(initConfig)
	}
	return config
}
