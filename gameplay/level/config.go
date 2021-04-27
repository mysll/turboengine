package level

import "github.com/BurntSushi/toml"

type Config struct {
	MinX, MinY, MaxX, MaxY float32
	TileWidth, TileHeight  float32
	ViewMaxRange           float32
}

func (c *Config) LoadFromFile(f string) {
	_, err := toml.DecodeFile(f, c)
	if err != nil {
		panic(err)
	}
}
