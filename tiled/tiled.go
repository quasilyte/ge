package tiled

import "encoding/json"

// https://doc.mapeditor.org/en/latest/reference/json-map-format/

type Tileset struct {
	Type string `json:"type"`

	Version      string `json:"version"`
	TiledVersion string `json:"tiledversion"`

	Name string `json:"name"`

	Spacing float64 `json:"spacing"`
	Margin  float64 `json:"margin"`

	NumTiles   int     `json:"tilecount"`
	NumColumns int     `json:"columns"`
	TileWidth  float64 `json:"tilewidth"`
	TileHeight float64 `json:"tileheight"`

	Tiles []Tile `json:"tiles"`
}

func UnmarshalTileset(jsonData []byte) (*Tileset, error) {
	var tileset Tileset
	if err := json.Unmarshal(jsonData, &tileset); err != nil {
		return nil, err
	}

	always := 1.0
	if tileset.Tiles == nil && tileset.NumTiles != 0 {
		// All tiles are perfectly ordered.
		// All tiles have default values.
		// IDs go from 0 to N-1.
		tiles := make([]Tile, tileset.NumTiles)
		for i := range tiles {
			tiles[i].ID = i
			tiles[i].Probability = &always
		}
		tileset.Tiles = tiles
	} else {
		for i := range tileset.Tiles {
			if tileset.Tiles[i].Probability == nil {
				tileset.Tiles[i].Probability = &always
			}
		}
	}

	return &tileset, nil
}

func (tileset *Tileset) TileByID(id int) *Tile {
	for i := range tileset.Tiles {
		if tileset.Tiles[i].ID == id {
			return &tileset.Tiles[i]
		}
	}
	return nil
}

type Tile struct {
	ID int `json:"id"`

	Class string `json:"class"`

	Probability *float64 `json:"probability"`
}

type Map struct {
	Height int
	Width  int

	Tilesets []TilesetRef

	Layers []MapLayer
}

type MapLayer struct {
	Name    string   `json:"name"`
	Objects []Object `json:"objects"`
}

type Object struct {
	GID      int          `json:"gid"`
	X        int          `json:"x"`
	Y        int          `json:"y"`
	Rotation int          `json:"rotation"`
	Props    []ObjectProp `json:"properties"`
}

func (o *Object) GetProp(name string) *ObjectProp {
	for i := range o.Props {
		if o.Props[i].Name == name {
			return &o.Props[i]
		}
	}
	return nil
}

func (o *Object) GetIntProp(name string, defaultValue int) int {
	p := o.GetProp(name)
	if p == nil {
		return defaultValue
	}
	if p.Type != "int" {
		return defaultValue
	}
	return int(p.Value.(float64))
}

func (o *Object) GetStringProp(name string, defaultValue string) string {
	p := o.GetProp(name)
	if p == nil {
		return defaultValue
	}
	if p.Type != "string" {
		return defaultValue
	}
	return p.Value.(string)
}

type ObjectProp struct {
	Name  string `json:"name"`
	Type  string
	Value any
}

type TilesetRef struct {
	FirstGID int    `json:"firstgid"`
	Source   string `json:"source"`
}

func UnmarshalMap(jsonData []byte) (*Map, error) {
	var m Map
	if err := json.Unmarshal(jsonData, &m); err != nil {
		return nil, err
	}
	return &m, nil
}
