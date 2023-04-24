package ge

import (
	"github.com/quasilyte/ge/physics"
)

const zindexMax = 6

type RootScene struct {
	context *Context

	controller   SceneController
	objects      []SceneObject
	addedObjects []SceneObject

	delayedFuncs []delayedFunc

	collisionEngine physics.CollisionEngine

	graphics [zindexMax][]SceneGraphics

	subSceneArray [zindexMax]Scene
}

type SimulationRunner struct {
	root *RootScene
}

func NewSimulatedScene(ctx *Context, controller SceneController) (*SimulationRunner, *Scene) {
	root := newRootScene()
	root.context = ctx
	root.controller = controller
	scene := &root.subSceneArray[1]
	return &SimulationRunner{root: root}, scene
}

func newRootScene() *RootScene {
	root := &RootScene{
		objects:      make([]SceneObject, 0, 32),
		addedObjects: make([]SceneObject, 0, 8),
		graphics: [zindexMax][]SceneGraphics{
			make([]SceneGraphics, 0, 16),
			make([]SceneGraphics, 0, 24),
			make([]SceneGraphics, 0, 8),
		},
	}
	for i := range root.subSceneArray {
		root.subSceneArray[i].zindex = uint8(i)
		root.subSceneArray[i].root = root
	}
	return root
}

func (scene *RootScene) addObject(o SceneObject, zindex uint) {
	if zindex < zindexMax {
		scene.addedObjects = append(scene.addedObjects, o)
		o.Init(&scene.subSceneArray[zindex])
		return
	}
	panic("z index overflow")
}

func (scene *RootScene) update(delta float64) {
	if len(scene.delayedFuncs) != 0 {
		funcs := scene.delayedFuncs[:0]
		for _, fn := range scene.delayedFuncs {
			fn.delay -= delta
			if fn.delay <= 0 {
				fn.action()
			} else {
				funcs = append(funcs, fn)
			}
		}
		scene.delayedFuncs = funcs
	}

	scene.collisionEngine.CalculateFrame()

	scene.controller.Update(delta)

	liveObjects := scene.objects[:0]
	for _, o := range scene.objects {
		if o.IsDisposed() {
			continue
		}
		o.Update(delta)
		liveObjects = append(liveObjects, o)
	}
	scene.objects = liveObjects
	scene.objects = append(scene.objects, scene.addedObjects...)
	scene.addedObjects = scene.addedObjects[:0]
}

type delayedFunc struct {
	delay  float64
	action func()
}
