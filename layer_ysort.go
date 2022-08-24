package ge

import (
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type YSortLayer struct {
	Visible  bool
	disposed bool

	nodes ysortNodeContainer
}

func NewYSortLayer() *YSortLayer {
	return &YSortLayer{
		Visible: true,
		nodes: ysortNodeContainer{
			list: make([]ysortNode, 0, 16),
		},
	}
}

func (l *YSortLayer) AddGraphics(g SceneGraphics, pos Pos) {
	l.nodes.list = append(l.nodes.list, ysortNode{g: g, pos: pos, y: 0})
}

func (l *YSortLayer) IsDisposed() bool {
	return l.disposed
}

func (l *YSortLayer) Dispose() {
	l.disposed = true
}

func (l *YSortLayer) Draw(screen *ebiten.Image) {
	if !l.Visible {
		return
	}

	// First, remove all disposed objects from the list.
	list := l.nodes.list[:0]
	for _, n := range l.nodes.list {
		if n.g.IsDisposed() {
			continue
		}
		list = append(list, n)
	}

	// Now collect all y values for the sorting.
	// While at it, check whether we need to actually re-sort the nodes slice.
	needSort := false
	for i, n := range list {
		newY := n.pos.Base.Y + n.pos.Offset.Y
		if n.y != newY {
			needSort = true
			list[i].y = newY
		}
	}

	// Do the actual sorting, if needed.
	if needSort {
		// Use sort.Sort instead of sort.Slice as its measurable faster for
		// smaller slices (it's not slower for the big ones too).
		// TODO: could use generic sorting function here to make it ~10 times faster.
		// TODO: consider using several Y spans so we can avoid sorting the spans that haven't changed.
		// TODO: when viewports/cameras are a thing, we may want to avoid sorting the objects that are not visible.
		sort.Sort(&l.nodes)
	}

	// Do the actual rendering now.
	// The slice should be in the correct order in respect to the Y coordinates.
	for _, n := range list {
		n.g.Draw(screen)
	}

	l.nodes.list = list
}

type ysortNode struct {
	g   SceneGraphics
	pos Pos
	y   float64
}

type ysortNodeContainer struct {
	list []ysortNode
}

func (c *ysortNodeContainer) Len() int {
	return len(c.list)
}

func (c *ysortNodeContainer) Less(i, j int) bool {
	return c.list[i].y < c.list[j].y
}

func (c *ysortNodeContainer) Swap(i, j int) {
	c.list[i], c.list[j] = c.list[j], c.list[i]
}
