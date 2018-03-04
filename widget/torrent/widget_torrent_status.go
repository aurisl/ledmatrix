package torrent

import (
	"github.com/dustin/go-humanize"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"
	"time"
	"github.com/aurisl/ledmatrix/command"
	"github.com/aurisl/ledmatrix/config"
	"github.com/aurisl/ledmatrix/draw"
	"github.com/aurisl/ledmatrix/matrix"
)

type animation struct {
	ctx                *gg.Context
	close              chan bool
	widgetCommand      *command.WidgetCommand
	reloadTorrent      bool
	torrentInformation string
	torrentProgress    string
	tick               int64
	clientActive       bool
	borderShader       *draw.BorderShared
	percentage         float64
	config             config.WidgetTorrentStatus
}

func Draw(toolkit *matrix.Toolkit,
	close chan bool,
	widgetCommand *command.WidgetCommand,
	widgetConfig config.WidgetTorrentStatus) {

	animation := &animation{
		ctx:                toolkit.Ctx,
		close:              close,
		widgetCommand:      widgetCommand,
		reloadTorrent:      true,
		clientActive:       false,
		torrentInformation: `N/A`,
		torrentProgress:    "",
		tick:               0,
		borderShader:       draw.NewBorderShader(),
		percentage:         0,
		config:             widgetConfig,
	}

	toolkit.MatrixToolkit.PlayAnimation(animation)
}

func (animation *animation) Next() (image.Image, <-chan time.Time, error) {

	animation.tick++

	draw.ClearCanvas(animation.ctx)

	if animation.reloadTorrent == true || animation.tick%30 == 0 {
		if updateTorrentInformation(animation) == false {
			animation.widgetCommand.Name = ""
			return nil, nil, io.EOF
		}
	}

	if animation.clientActive == true && animation.torrentProgress == "" {
		drawRedScreen(animation)

		select {
		case <-animation.close:
			return nil, nil, io.EOF
		default:
			return animation.ctx.Image(), time.After(time.Millisecond * 50), nil
		}
	}

	isScrollingCompleted := draw.TextScrolling(animation.torrentInformation, 12, animation.ctx)
	if isScrollingCompleted == true {
		animation.reloadTorrent = true
	}

	draw.GradientLine(animation.ctx)
	animation.borderShader.DrawBorderShader(animation.ctx)

	percentageColor := draw.BlendingPercentageColor(uint8(animation.percentage))
	draw.Text(animation.torrentProgress, 8, 30, animation.ctx, percentageColor)

	select {
	case <-animation.close:
		return nil, nil, io.EOF
	default:
		return animation.ctx.Image(), time.After(time.Millisecond * 70), nil
	}
}

func updateTorrentInformation(animation *animation) bool {

	UTorrentClient, err := NewUTorrentClient(
		animation.config.TorrentWebApiUrl,
		animation.config.Username,
		animation.config.Password)

	if err != nil {
		animation.widgetCommand.Name = ""
		return false
	}

	torrentList, err := UTorrentClient.getList()
	if err != nil {
		animation.widgetCommand.Name = ""
		return false
	}

	animation.torrentInformation = ""
	animation.torrentProgress = ""
	for _, element := range torrentList.Torrents {
		if element.Remaining == 0 {
			continue
		}

		name := strings.Replace(element.Name, ".", " ", -1)
		animation.torrentInformation = name + " : " + humanize.Bytes(element.Remaining)
		animation.percentage = element.Progress / 10
		animation.torrentProgress = strconv.FormatFloat(animation.percentage, 'f', 0, 64) + "%"

		break
	}

	animation.clientActive = true
	animation.reloadTorrent = false

	return true
}

func drawRedScreen(animation *animation) {

	if animation.tick % 10 == 0 {
		animation.ctx.SetColor(color.RGBA{255, 0, 0, 255})
	} else {
		animation.ctx.SetColor(color.RGBA{0, 0, 0, 255})
	}
	animation.ctx.Clear()
}
