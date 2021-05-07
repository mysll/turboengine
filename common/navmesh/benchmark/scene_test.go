package benchmark

import (
	"math/rand"
	"testing"
	"time"
	"turboengine/common/navmesh"
	"turboengine/common/navmesh/tests"
)

const path1 = "../tests/meshes/scene1.obj.tile.bin"
const path2 = "../tests/meshes/scene1.obj.tilecache.bin"

var scn1 *navmesh.StaticScene
var scn2 *navmesh.DynamicScene
var scn3 *navmesh.DynamicScene

func init() {
	rand.Seed(time.Now().UTC().UnixNano())

	scn1 = navmesh.NewStaticScene()
	tests.InitScene(scn1.Scene, path1)

	scn2 = navmesh.NewDynamicScene(navmesh.HEIGHT_MODE_1)
	tests.InitScene(scn2.Scene, path1)

	scn3 = navmesh.NewDynamicScene(navmesh.HEIGHT_MODE_2)
	tests.InitScene(scn3.Scene, path2)
}

func Benchmark_Scene1(t *testing.B) {
	t.N = 10000
	for i := 0; i < t.N; i++ {
		scn1.Simulation(0.025)
	}
}

func Benchmark_Scene2(t *testing.B) {
	t.N = 10000
	for i := 0; i < t.N; i++ {
		scn2.Simulation(0.025)
	}
}

func Benchmark_Scene3(t *testing.B) {
	t.N = 10000
	for i := 0; i < t.N; i++ {
		scn3.Simulation(0.025)
	}
}
