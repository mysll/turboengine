package tower

import (
	"fmt"
	"testing"
	"turboengine/gameplay/object"
)

type aoi struct {
	t *testing.T
}

func (a *aoi) OnEnterAOI(watcher, target object.ObjectId) {
	a.t.Logf("[%d] :%d enter \n", watcher, target)
}

func (a *aoi) OnLeaveAOI(watcher, target object.ObjectId) {
	a.t.Logf("[%d] :%d leave \n", watcher, target)
}

func drawAoi(toi *TowerAOI) {
	fmt.Printf(" |")
	for i := 0; i <= toi.max.X; i++ {
		fmt.Printf("%d\t|", i)
	}
	fmt.Println()
	for i := 0; i <= toi.max.Y; i++ {
		fmt.Printf("%d|", i)
		for j := 0; j <= toi.max.X; j++ {
			for id := range toi.towers[i][j].Ids {
				fmt.Printf("%d,", id)
			}
			fmt.Printf("\t|")
		}
		fmt.Println()
	}
}
func TestNewTowerAOI(t *testing.T) {
	toi := NewTowerAOI(500, 500, 50, 50, 10, &aoi{t})
	toi.Enter(1, object.Vec3{75, 0, 75}, 2)
	toi.Enter(2, object.Vec3{125, 0, 500}, 2)
	toi.Enter(3, object.Vec3{200, 0, 150}, 2)
	drawAoi(toi)
	toi.Move(1, object.Vec3{75, 0, 75}, object.Vec3{25, 0, 125}, 2)
	drawAoi(toi)
	toi.Move(2, object.Vec3{125, 0, 500}, object.Vec3{75, 0, 75}, 2)
	drawAoi(toi)
	toi.Level(1, object.Vec3{25, 0, 125}, 2)
	drawAoi(toi)
}
