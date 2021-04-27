package tower

import (
	"fmt"
	"testing"
	. "turboengine/common/datatype"
)

type aoi struct {
	t *testing.T
}

func (a *aoi) OnEnterAOI(watcher, target ObjectId) {
	fmt.Printf("[%d] :%d enter \n", watcher, target)
}

func (a *aoi) OnLeaveAOI(watcher, target ObjectId) {
	fmt.Printf("[%d] :%d leave \n", watcher, target)
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
	toi := NewTowerAOI(-250, -250, 250, 250, 50, 50, 200, &aoi{t})
	fmt.Println("info:", toi.max)
	toi.Enter(1, Vec3{0, 0, 0}, 100)
	toi.Enter(2, Vec3{125, 0, 250}, 100)
	toi.Enter(3, Vec3{200, 0, 150}, 100)
	drawAoi(toi)
	toi.Move(1, Vec3{0, 0, 0}, Vec3{100, 0, 125}, 100)
	drawAoi(toi)
	toi.Move(2, Vec3{125, 0, 250}, Vec3{75, 0, 75}, 100)
	drawAoi(toi)
	toi.Level(1, Vec3{100, 0, 125}, 100)
	drawAoi(toi)
}
