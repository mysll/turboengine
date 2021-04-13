package entity

import (
	"fmt"
	"testing"
	"turboengine/gameplay/object"
)

func TestNewPlayer(t *testing.T) {
	player := NewPlayer()
	player.NameChange(OnNameChange)
	player.SetSilent(true)
	player.SetName("test")
	player.SetSilent(false)
	player.SetName("test1")
	player.SetPos(object.V3(1, 2, 3))
	fmt.Println(player.Pos().X(), player.Pos().Y(), player.Pos().Z())
	fmt.Println(player.Dirty(), player.PublicDirty(), player.PrivateDirty())
}

func OnNameChange(self interface{}, index int, val interface{}) {
	fmt.Printf("old:%v, new:%v\n", val, self.(*Player).Name())
}
