package level

import "github.com/BurntSushi/toml"

type Config struct {
	MinX, MinY, MaxX, MaxY float32
	TileWidth, TileHeight  float32
	ViewMaxRange           float32
	MapData                string
	TerrainType            string
}

func (c *Config) LoadFromFile(f string) {
	_, err := toml.DecodeFile(f, c)
	if err != nil {
		panic(err)
	}
}

func (c *Config) LoadFromData(data string) {
	_, err := toml.Decode(data, c)
	if err != nil {
		panic(err)
	}
}
