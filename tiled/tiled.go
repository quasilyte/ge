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
	// TODO: handle multi-row tilesets.
	allTiles := make([]Tile, 0, len(tileset.Tiles))
	id := 0
	for _, t := range tileset.Tiles {
		if t.ID != id {
			// Insert the missing (implicit) tiles.
			idSeq := id
			for idSeq < t.ID {
				allTiles = append(allTiles, Tile{
					ID:          idSeq,
					Probability: 1,
				})
				idSeq++
			}
		}
		allTiles = append(allTiles, t)
		id++
	}
	for id < tileset.NumTiles {
		allTiles = append(allTiles, Tile{
			ID:          id,
			Probability: 1,
		})
		id++
	}
	tileset.Tiles = allTiles
	return &tileset, nil
}

type Tile struct {
	ID int `json:"id"`

	Probability float64 `json:"probability"`
}
