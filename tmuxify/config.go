package tmuxify

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Session struct {
		Name string `toml:"name"`
		Main uint32 `toml:"main"`
	} `toml:"session"`

	Window []struct{
		Name string `toml:"name"`
		Cmds []string `toml:"cmds"`
	} `toml:"window"`
};

func GetConfig() Config{
	data ,err := os.ReadFile("./.tmuxify.toml");
	if err != nil {
		fmt.Print(err);
		os.Exit(-1);
	}

	config := Config{}
	err = toml.Unmarshal(data , &config);
	if err != nil {
		fmt.Print(err)
		os.Exit(-1);
	}

	return config
}
