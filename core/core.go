package core

import (
	"fmt"
	"github.com/xuri/excelize/v2"
)

func Run(path string, isUp bool) {
	etfIns := &etfDays{
		all:      nil,
		pinCa:    4.0, //关键点和次级转折幅度定义
		turnCa:   3.0, //逆转警告幅度定义
		pin1:     &etfDaysPer{},
		pin2:     &etfDaysPer{},
		starIsUp: isUp,
		lastPin:  nil,
		keepDays: 0,
	}
	etfIns.all = initEtfData(path)
	etfIns.think()
}

func initEtfData(path string) []*etfDaysPer {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Println(err)
		return nil
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
		return nil
	}
	res := make([]*etfDaysPer, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		res = append(res, &etfDaysPer{
			dateD: rows[i][0],
			val:   ToFloat64(rows[i][1]),
		})
	}
	return res
}
