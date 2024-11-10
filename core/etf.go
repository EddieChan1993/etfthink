package core

import (
	"fmt"
	"math"
)

type etfDays struct {
	all      []*etfDaysPer
	pinCa    float64 //关键点变化幅度 6%
	turnCa   float64 //移转信号幅度 3%
	pin1     float64 //关键点1
	pin2     float64 //关键点2
	starIsUp bool
	oldPin   *etfDaysPer
}

type etfDaysPer struct {
	dateD string
	val   float64
}

func (e *etfDays) think() {
	isUp := e.starIsUp
	fmt.Println(e.isUpStr(isUp))
	e.oldPin = e.all[0]
	for _, per := range e.all {
		caRate := math.Abs(e.oldPin.val-per.val) / e.oldPin.val * 100
		if isUp {
			if e.pin1 != 0 {
				if per.val > e.pin1 && math.Abs(per.val-e.pin1)/e.pin1*100 > e.turnCa {
					isUp = true
					fmt.Print(fmt.Sprintf("%s %.4f -彻底突破关键点1，%s趋势形成\n", per.dateD, per.val, e.isUpStr(isUp)))
					e.pin1 = 0
					e.pin2 = 0
				} else {
					if e.oldPin.val > per.val {
						if caRate >= e.turnCa {
							fmt.Print(fmt.Sprintf("%s %.4f -反向=%.0f点 %.4f%%【逆转警告】\n", per.dateD, per.val, e.turnCa, caRate))
						}
						if caRate >= e.pinCa {
							isUp = false
							fmt.Print(fmt.Sprintf("%s %.4f -转向>=%.0f点-次级%s\n", per.dateD, e.pinCa, per.val, e.isUpTmpStr(isUp)))
							e.oldPin = per
						}
					} else {
						fmt.Print(fmt.Sprintf("%s %.4f -自然%s持续\n", per.dateD, per.val, e.isUpTmpStr(isUp)))
						e.oldPin = per
					}
				}
			} else {
				if e.oldPin.val > per.val {
					if caRate >= e.pinCa {
						isUp = false
						fmt.Print(fmt.Sprintf("%s %.4f -【关键点1】-自然%s阶段\n", e.oldPin.dateD, e.oldPin.val, e.isUpTmpStr(isUp)))
						e.oldPin = per
						e.pin1 = per.val
					}
				} else {
					e.oldPin = per
				}
			}
		} else {
			if e.pin2 != 0 {
				if per.val < e.pin2 && math.Abs(per.val-e.pin2)/e.pin1*100 > e.turnCa {
					isUp = false
					fmt.Print(fmt.Sprintf("%s %.4f -彻底突破关键点2，%s趋势形成\n", per.dateD, per.val, e.isUpStr(isUp)))
					e.pin1 = 0
					e.pin2 = 0
				} else {
					if per.val > e.oldPin.val {
						if caRate >= e.turnCa {
							fmt.Print(fmt.Sprintf("%s %.4f -反向=%.0f点 %.4f%%【逆转警告】\n", per.dateD, per.val, e.turnCa, caRate))
						}
						if caRate >= e.pinCa {
							isUp = true
							fmt.Print(fmt.Sprintf("%s %.4f -转向>=%.0f点-次级%s\n", per.dateD, e.pinCa, per.val, e.isUpTmpStr(isUp)))
							e.oldPin = per
						}
					} else {
						fmt.Print(fmt.Sprintf("%s %.4f -自然%s持续\n", per.dateD, per.val, e.isUpTmpStr(isUp)))
						e.oldPin = per
					}
				}
			} else {
				if e.oldPin.val < per.val {
					if caRate >= e.pinCa {
						isUp = true
						fmt.Print(fmt.Sprintf("%s %.4f -【关键点2】-自然%s阶段\n", e.oldPin.dateD, e.oldPin.val, e.isUpTmpStr(isUp)))
						e.oldPin = per
						e.pin2 = per.val
					}
				} else {
					e.oldPin = per
				}
			}
		}
	}
}

func (e *etfDays) isUpStr(isUp bool) string {
	if isUp {
		return "上升"
	} else {
		return "下降"
	}
}

func (e *etfDays) isUpTmpStr(isUp bool) string {
	if isUp {
		return "回升"
	} else {
		return "回撤"
	}
}
