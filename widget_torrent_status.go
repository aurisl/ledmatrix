package main

import (
	"github.com/dustin/go-humanize"
	"github.com/fogleman/gg"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"
	"time"
)

type TorrentStatus struct {
	ctx                *gg.Context
	close              chan bool
	widget             *Widget
	reloadTorrent      bool
	torrentInformation string
	torrentProgress    string
	tick               int64
	clientActive       bool
	borderShader       *BorderShared
	percentage         float64
	config             WidgetTorrentStatusConfig
}

func DrawTorrentStatus(toolkit *rgbmatrix.ToolKit, close chan bool, widget *Widget, config WidgetTorrentStatusConfig) {

	animation := &TorrentStatus{
		ctx:                gg.NewContext(32, 32),
		close:              close,
		widget:             widget,
		reloadTorrent:      true,
		clientActive:       false,
		torrentInformation: `N/A`,
		torrentProgress:    "",
		tick:               0,
		borderShader:       NewBorderShader(),
		percentage:         0,
		config:             config,
	}

	toolkit.PlayAnimation(animation)
}

func (animation *TorrentStatus) Next() (image.Image, <-chan time.Time, error) {

	animation.tick++

	initializeTorrentCanvas(animation)

	if animation.reloadTorrent == true || animation.tick%30 == 0 {

		UTorrentClient, err := NewUTorrentClient(
			animation.config.TorrentWebApiUrl,
			animation.config.Username,
			animation.config.Password)

		if err != nil {
			animation.widget.name = ""
			return nil, nil, io.EOF
		}

		torrentList, err := UTorrentClient.getList()
		if err != nil {
			animation.widget.name = ""
			return nil, nil, io.EOF
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

	isScrollingCompleted := DrawScrollingText(animation.torrentInformation, 12, animation.ctx)
	if isScrollingCompleted == true {
		animation.reloadTorrent = true
	}

	drawLineTorrent(animation)
	animation.borderShader.DrawBorderShader(animation.ctx)

	percentageColor := getBlendedColor(uint8(animation.percentage))
	DrawText(animation.torrentProgress, 8, 30, animation.ctx, percentageColor)

	select {
	case <-animation.close:
		return nil, nil, io.EOF
	default:
		return animation.ctx.Image(), time.After(time.Millisecond * 70), nil
	}
}

func getBlendedColor(percentage uint8) color.RGBA {

	var red uint8 = 0
	var green uint8 = 0
	var blue uint8 = 0

	if percentage > 50 {
		red = uint8(1 - float64(2*(percentage-50))/100*255)
		green = 255
	} else {
		red = 255
		green = uint8(float64(2*percentage) / 100 * 255)
	}

	return color.RGBA{R: red, G: green, B: blue, A: 255}
}

func drawLineTorrent(animation *TorrentStatus) {

	grad := gg.NewLinearGradient(1, 1, 32, 1)
	grad.AddColorStop(0, color.RGBA{0, 255, 0, 255})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 255})
	grad.AddColorStop(0.5, color.RGBA{255, 0, 0, 255})
	animation.ctx.SetStrokeStyle(grad)

	animation.ctx.DrawLine(1, 17, 32, 17)
	animation.ctx.SetLineWidth(1)
	animation.ctx.Stroke()
}

func drawRedScreen(animation *TorrentStatus) {

	if animation.tick%10 == 0 {
		grad := gg.NewRadialGradient(15, 15, 0, 15, 15, 32)

		grad.AddColorStop(1, color.RGBA{255, 0, 0, 255})

		animation.ctx.SetFillStyle(grad)
		animation.ctx.DrawRectangle(0, 0, 32, 32)
		animation.ctx.Fill()

		return
	}

	animation.ctx.SetRGB(0, 0, 0)
	animation.ctx.Clear()
	animation.ctx.SetColor(color.RGBA{255, 255, 255, 255})
	animation.ctx.Fill()
}

func initializeTorrentCanvas(animation *TorrentStatus) {

	animation.ctx.SetRGB(0, 0, 0)
	animation.ctx.Clear()
	animation.ctx.SetColor(color.RGBA{255, 255, 255, 255})
}
