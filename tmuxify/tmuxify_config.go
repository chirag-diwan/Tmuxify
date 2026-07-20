package tmuxify

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type session_t struct{
	Name string `toml:"name"`
	Main int `toml:"main"`
}

type window_t struct{
	Name string `toml:"name"`
	Cmds []string `toml:"cmds"`
}

type Config struct {
	Session session_t`toml:"session"`
	Window []window_t `toml:"window"`
};

func GetDefaultConfig(root_dir string)Config{
	return Config{
		Session: session_t{
			Name: filepath.Base(root_dir),
			Main: 1,
		},
		Window: []window_t{
			{ Name: "main", },
		},
	}
}

func GetConfig(root_dir string) Config{
	data ,err := os.ReadFile(filepath.Join(root_dir , ".tmuxify.toml"));
	if err != nil {
		fmt.Print(err);
		return GetDefaultConfig(root_dir)
	}

	config := Config{}
	err = toml.Unmarshal(data , &config);

	if err != nil {
		fmt.Print(err)
		return GetDefaultConfig(root_dir)
	}

	return config
}
