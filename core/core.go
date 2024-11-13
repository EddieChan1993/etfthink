package core

import (
	"fmt"
	"github.com/EddieChan1993/gcore/utils/cast"
	"github.com/xuri/excelize/v2"
	"path/filepath"
)

func Run(path string, isUp bool) {
	etfIns := &etfDays{
		all:          nil,
		pinCa:        4.0, //关键点和次级转折幅度定义
		turnCa:       3.0, //逆转警告幅度定义
		pin1:         &etfDaysPer{},
		pin2:         &etfDaysPer{},
		starIsUp:     isUp,
		lastPin:      nil,
		keepDays:     0,
		keepTurnDays: 0,
		points:       make(map[string]float32),
	}
	etfInsData, lineChartIns := initEtfData(path)
	etfIns.all = etfInsData
	points := etfIns.think()
	lineChartIns.setPoints(points)
	lineChartIns.ioWrite()
}

func initEtfData(path string) ([]*etfDaysPer, *lineChart) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Get all the rows in the Sheet1.
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	res := make([]*etfDaysPer, 0, len(rows))
	x := make([]string, 0, len(rows))
	y := make([]float32, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		res = append(res, &etfDaysPer{
			dateD: rows[i][0],
			val:   cast.ToFloat64(rows[i][1]),
		})
		x = append(x, rows[i][0])
		y = append(y, cast.ToFloat32(rows[i][1]))
	}
	filename := filepath.Base(path)
	return res, &lineChart{
		x:     x,
		y:     y,
		title: filename,
	}
}
