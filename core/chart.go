package core

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"io"
	"os"
)

type lineChartIns struct {
	x     []string
	data  []float32
	title string
}

func (l *lineChartIns) lineData() []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(l.data); i++ {
		items = append(items, opts.LineData{Value: l.data[i]})
	}
	return items
}

func (l *lineChartIns) lineBase() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: l.title, Subtitle: "净值变化"}),
	)

	line.SetXAxis(l.x).
		AddSeries("Category A", l.lineData())
	return line
}

func (l *lineChartIns) ioWrite() {
	page := components.NewPage()
	page.AddCharts(
		l.lineBase(),
	)
	f, err := os.Create("html/line.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
