package coach

import (
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/client-player-go"
	"testing"
)

func TestImageDataSaver_SaveSample(t *testing.T) {
	saver := NewImageBasedDataSever("./", arena.HomeTeam)

	onlyMe := &client.Player{
		TeamPlace: arena.HomeTeam,
	}

	ball := client.Ball{}
	ball.Coords = physics.Point{
		PosX: arena.FieldCenter.PosX + 500,
		PosY: arena.FieldCenter.PosY + 500,
	}
	onlyMe.Coords = arena.FieldCenter
	_, err := saver.SaveSample(client.GameInfo{
		Ball: ball,
		AwayTeam: client.Team{
			Players: []*client.Player{},
		},
		HomeTeam: client.Team{
			Players: []*client.Player{onlyMe},
		},
	})

	if err != nil {
		t.Fatalf("retuened error: %s", err)
	}
}
