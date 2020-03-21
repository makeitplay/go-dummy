package coach

import (
	"fmt"
	"github.com/makeitplay/arena"
	"github.com/makeitplay/arena/physics"
	"github.com/makeitplay/arena/units"
	"github.com/makeitplay/client-player-go"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
)

type DataSnapshot interface {
	Save(lavel string) error
}

type DataSaver interface {
	SaveSample(state client.GameInfo) (DataSnapshot, error)
	//SetPlayerPerspective()
	//SetFileNameFormatter(NameFormatter)
}

const imageWidth = 400
const imageLength = 200
const ballSize = 2
const playerSize = 4
const scaleDiff = 0.02

func NewImageBasedDataSever(dataDir string, ourTeamPlace arena.TeamPlace) DataSaver {
	return &imageDataSaver{dir: dataDir, ourTeamPlace: ourTeamPlace}
}

type imageDataSaver struct {
	dir          string
	ourTeamPlace arena.TeamPlace
	//labels            []string
	//playerPerspective *client.Player
}

type imageDataSnapshot struct {
	img draw.Image
	dir string
}

type Options struct {
	ShowPerspective bool
}

func (i *imageDataSnapshot) Save(label string) error {
	labelPath := filepath.Join(i.dir, label)
	_ = os.MkdirAll(labelPath, os.ModePerm)

	f, err := os.Create(fmt.Sprintf("%s/%d.jpeg", labelPath, time.Now().UnixNano()))
	if err != nil {
		return fmt.Errorf("could not create the image file: %s", err)
	}
	err = jpeg.Encode(f, i.img, &jpeg.Options{Quality: 100})
	if err != nil {
		return fmt.Errorf("could not encode the JPEG file: %s", err)
	}
	return nil
}

//opts Options
func (s *imageDataSaver) SaveSample(state client.GameInfo) (DataSnapshot, error) {
	upLeft := image.Point{}
	lowRight := image.Point{X: imageWidth, Y: imageLength}
	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	for _, p := range state.HomeTeam.Players {
		s.drawPlayer(img, p)
	}
	for _, p := range state.AwayTeam.Players {
		s.drawPlayer(img, p)
	}

	s.drawBall(img, state.Ball)

	return &imageDataSnapshot{img: img, dir: s.dir}, nil

}

func (s *imageDataSaver) drawPlayer(img draw.Image, player *client.Player) {
	playerColor := color.RGBA{R: 255, A: 255}
	if player.TeamPlace == s.ourTeamPlace {
		playerColor = color.RGBA{G: 255, A: 255}
	}
	halfBody := playerSize / 2
	coords := player.Coords
	if s.ourTeamPlace == arena.AwayTeam {
		coords = MirrorCoordsToAway(coords)
	}
	coords = fixCoord(coords)

	x := coords.PosX - halfBody
	y := coords.PosY - halfBody
	redRect := image.Rect(x, y, x+playerSize, y+playerSize) //  geometry of 2nd rectangle

	draw.Draw(img, redRect, &image.Uniform{C: playerColor}, image.ZP, draw.Src)
}

func (s *imageDataSaver) drawBall(img draw.Image, ball client.Ball) {
	halfBody := ballSize / 2

	coords := ball.Coords
	if s.ourTeamPlace == arena.AwayTeam {
		coords = MirrorCoordsToAway(coords)
	}
	coords = fixCoord(coords)

	x := coords.PosX - halfBody
	y := coords.PosY - halfBody
	redRect := image.Rect(x, y, x+ballSize, y+ballSize) //  geometry of 2nd rectangle

	draw.Draw(img, redRect, &image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.ZP, draw.Src)
}

func fixCoord(gamePos physics.Point) physics.Point {
	return physics.Point{
		PosX: int(float64(gamePos.PosX) * scaleDiff),
		PosY: imageLength - int(float64(gamePos.PosY)*scaleDiff),
	}
}

func MirrorCoordsToAway(coords physics.Point) physics.Point {
	return physics.Point{
		PosX: units.FieldWidth - coords.PosX,
		PosY: units.FieldHeight - coords.PosY,
	}
}