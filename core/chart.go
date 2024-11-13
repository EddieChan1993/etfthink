package core

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
)

type lineChart struct {
	x      []string
	y      []float32
	title  string
	points map[string]float32
	pointY []opts.EffectScatterData
}

func (l *lineChart) setPoints(points map[string]float32) {
	pointY := make([]opts.EffectScatterData, len(l.x))
	for i, x := range l.x {
		if val, had := points[x]; had {
			pointY[i] = opts.EffectScatterData{
				Name:  x,
				Value: val,
			}
		}

	}
	l.pointY = pointY
}

func (l *lineChart) lineData() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(l.y); i++ {
		items = append(items, opts.LineData{Value: l.y[i]})
	}
	return items
}

func (l *lineChart) pointBasic() *charts.EffectScatter {
	scatter := charts.NewEffectScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "关键点"}),
	)

	scatter.SetXAxis(l.points).AddSeries("", l.pointY)

	return scatter
}

func (l *lineChart) descLine() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: l.title, Subtitle: "净值变化"}),
	)
	line.SetXAxis(l.x).
		AddSeries("单位净值", l.lineData()).
		Overlap(l.pointBasic())
	//line.Overlap(scatterBase())
	return line
}

func (l *lineChart) ioWrite() {
	page := components.NewPage()
	page.AddCharts(
		l.descLine(),
	)
	f, err := os.Create("html/line.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
