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
	keepDays int32
}

type etfDaysPer struct {
	dateD string
	val   float64
}

func (e *etfDays) think() {
	fmt.Printf("//==================== %s 趋势====================//\n)", e.isUpStr(e.starIsUp))
	e.oldPin = e.all[0]
	caRate := 0.0
	for _, per := range e.all {
		caRate = math.Abs(e.oldPin.val-per.val) / e.oldPin.val * 100
		if e.starIsUp {
			e.upThink(per, caRate)
		} else {
			e.downThink(per, caRate)
		}
	}
}

func (e *etfDays) upThink(per *etfDaysPer, caRate float64) {
	e.keepDays++
	if e.pin1 != 0 {
		if per.val > e.pin1 && math.Abs(per.val-e.pin1)/e.pin1*100 > e.turnCa {
			fmt.Printf("%s %.4f -突破【关键点1】，%s趋势恢复\n", per.dateD, per.val, e.isUpStr(e.starIsUp))
			e.pin1 = 0
			e.pin2 = 0
		} else {
			if e.oldPin.val > per.val {
				fmt.Printf("%s %.4f -持续%s后 %d天 反向\n", per.dateD, per.val, e.isUpStr(e.starIsUp), e.keepDays)
				e.keepDays = 0
				if caRate >= e.turnCa {
					fmt.Printf("%s %.4f -反向波动 %.0f点 %.4f%%【逆转警告】\n", per.dateD, per.val, e.turnCa, caRate)
				}
				if caRate >= e.pinCa {
					e.starIsUp = false
					fmt.Printf("%s %.4f -转向>=%.0f点-次级%s\n", per.dateD, e.pinCa, per.val, e.isUpTmpStr(e.starIsUp))
					e.oldPin = per
					return
				}
			} else {
				e.oldPin = per
			}
		}
	} else {
		if e.oldPin.val > per.val {
			fmt.Printf("%s %.4f -持续%s后 %d天 反向\n", per.dateD, per.val, e.isUpStr(e.starIsUp), e.keepDays)
			e.keepDays = 0
			if caRate >= e.pinCa {
				e.starIsUp = false
				fmt.Printf("%s %.4f -【关键点1 - %s】开始进入-自然%s阶段\n", per.dateD, per.val, e.oldPin.dateD, e.isUpTmpStr(e.starIsUp))
				e.oldPin = per
				e.pin1 = per.val
				return
			}
		} else {
			e.oldPin = per
		}
	}
}

func (e *etfDays) downThink(per *etfDaysPer, caRate float64) {
	e.keepDays++
	if e.pin2 != 0 {
		if per.val < e.pin2 && math.Abs(per.val-e.pin2)/e.pin1*100 > e.turnCa {
			fmt.Print(fmt.Sprintf("%s %.4f -突破【关键点2】，%s趋势恢复\n", per.dateD, per.val, e.isUpStr(e.starIsUp)))
			e.pin1 = 0
			e.pin2 = 0
		} else {
			if per.val > e.oldPin.val {
				fmt.Printf("%s %.4f -持续%s后 %d天 反向\n", per.dateD, per.val, e.isUpStr(e.starIsUp), e.keepDays)
				e.keepDays = 0
				if caRate >= e.turnCa {
					fmt.Print(fmt.Sprintf("%s %.4f -反向波动 %.0f点 %.4f%%【逆转警告】\n", per.dateD, per.val, e.turnCa, caRate))
				}
				if caRate >= e.pinCa {
					e.starIsUp = true
					fmt.Print(fmt.Sprintf("%s %.4f -转向>=%.0f点-次级%s\n", per.dateD, e.pinCa, per.val, e.isUpTmpStr(e.starIsUp)))
					e.oldPin = per
					return
				}
			} else {
				e.oldPin = per
			}
		}
	} else {
		if e.oldPin.val < per.val {
			fmt.Printf("%s %.4f -持续%s后 %d天 反向\n", per.dateD, per.val, e.isUpStr(e.starIsUp), e.keepDays)
			e.keepDays = 0
			if caRate >= e.pinCa {
				e.starIsUp = true
				fmt.Printf("%s %.4f -【关键点2 - %s】开始进入-自然%s阶段\n", per.dateD, per.val, e.oldPin.dateD, e.isUpTmpStr(e.starIsUp))
				e.oldPin = per
				e.pin2 = per.val
				return
			}
		} else {
			e.oldPin = per
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
