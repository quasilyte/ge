package ge

import "testing"

func TestIfaceImpl(t *testing.T) {
	_ = []SceneGraphicsLayer{
		(*SimpleLayer)(nil),
		(*YSortLayer)(nil),
		(*MultiLayer)(nil),
	}
}
