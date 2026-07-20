package tmuxify

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type ProjectConfig struct {
	Roots []string `toml:"roots"`
	Ignore []string `toml:"ignore"`
	MaxDepth int `toml:"max_depth"`
}

var defaultProjectConfig = ProjectConfig{
	Roots: []string{},
	Ignore: []string{".git", "node_modules", ".cache", ".bun", ".cargo"},
	MaxDepth: 4,
}

func GetAppConfig() ProjectConfig{
	homeDir , err := os.UserHomeDir();
	if err != nil {
		fmt.Printf("Tmuxify was unable identify user home dir")
		return defaultProjectConfig
	}

	data , err := os.ReadFile(filepath.Join(homeDir , ".tmuxify-conf.toml"))
	if err != nil {
		return defaultProjectConfig;
	}
	config := ProjectConfig{}
	err = toml.Unmarshal(data , &config);
	if err != nil{
		return defaultProjectConfig;
	}
	return config
}
