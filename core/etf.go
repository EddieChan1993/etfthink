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
	lastPin  *etfDaysPer
	keepDays int32
}

type etfDaysPer struct {
	dateD string
	val   float64
}

func (e *etfDays) think() {
	fmt.Printf("//==================== %s 趋势====================//\n", e.isUpStr(e.starIsUp))
	e.lastPin = e.all[0]
	caRate := 0.0
	for _, per := range e.all {
		caRate = math.Abs(e.lastPin.val-per.val) / e.lastPin.val * 100
		if e.starIsUp {
			e.upThink(per, caRate)
		} else {
			e.downThink(per, caRate)
		}
	}
	fmt.Printf("//==================== End ====================//\n")
}

func (e *etfDays) upThink(per *etfDaysPer, caRate float64) {
	if e.pin1 != 0 {
		if per.val > e.pin1 && math.Abs(per.val-e.pin1)/e.pin1*100 > e.turnCa {
			e.keepDays++
			fmt.Printf(" %s 突破【关键点1】，%s趋势恢复\n", e.log(per), e.isUpStr(e.starIsUp))
			e.pin1 = 0
			e.pin2 = 0
		} else {
			if e.lastPin.val > per.val {
				e.logKeepDays(per)
				e.keepDays = 0
				if caRate >= e.turnCa {
					fmt.Print(fmt.Sprintf(" %s 于[%s]%s波动 %.0f点 %.2f【逆转警告】\n", e.log(per), e.log(e.lastPin), e.isUpStr(!e.starIsUp), e.turnCa, caRate))
				}
				if caRate >= e.pinCa {
					e.starIsUp = false
					fmt.Print(fmt.Sprintf(" %s 于[%s]转向 %.2f点->次级%s\n", e.log(per), e.log(e.lastPin), caRate, e.isUpTmpStr(e.starIsUp)))
					e.lastPin = per
					return
				}
			} else {
				e.keepDays++
				e.lastPin = per
			}
		}
	} else {
		if e.lastPin.val > per.val {
			e.logKeepDays(per)
			e.keepDays = 0
			if caRate >= e.pinCa {
				e.starIsUp = false
				fmt.Printf(" %s【关键点1 %s】开始进入->自然%s阶段 \n", e.log(per), e.lastPin.dateD, e.isUpTmpStr(e.starIsUp))
				e.pin1 = e.lastPin.val
				e.lastPin = per
				return
			}
		} else {
			e.keepDays++
			e.lastPin = per
		}
	}
}

func (e *etfDays) downThink(per *etfDaysPer, caRate float64) {
	if e.pin2 != 0 {
		if per.val < e.pin2 && math.Abs(per.val-e.pin2)/e.pin1*100 > e.turnCa {
			e.keepDays++
			fmt.Print(fmt.Sprintf(" %s 突破【关键点2】，%s趋势恢复\n", e.log(per), e.isUpStr(e.starIsUp)))
			e.pin1 = 0
			e.pin2 = 0
		} else {
			if per.val > e.lastPin.val {
				e.logKeepDays(per)
				e.keepDays = 0
				if caRate >= e.turnCa {
					fmt.Print(fmt.Sprintf(" %s 于[%s]%s波动 %.0f点 %.2f【逆转警告】\n", e.log(per), e.log(e.lastPin), e.isUpStr(!e.starIsUp), e.turnCa, caRate))
				}
				if caRate >= e.pinCa {
					e.starIsUp = true
					fmt.Print(fmt.Sprintf(" %s 于[%s]转向 %.2f点->次级%s\n", e.log(per), e.log(e.lastPin), caRate, e.isUpTmpStr(e.starIsUp)))
					e.lastPin = per
					return
				}
			} else {
				e.keepDays++
				e.lastPin = per
			}
		}
	} else {
		if e.lastPin.val < per.val {
			e.logKeepDays(per)
			e.keepDays = 0
			if caRate >= e.pinCa {
				e.starIsUp = true
				fmt.Printf(" %s【关键点2 %s】开始进入->自然%s阶段\n", e.log(per), e.lastPin.dateD, e.isUpTmpStr(e.starIsUp))
				e.pin2 = e.lastPin.val
				e.lastPin = per
				return
			}
		} else {
			e.keepDays++
			e.lastPin = per
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

func (e *etfDays) log(per *etfDaysPer) string {
	return fmt.Sprintf("%s %.4f", per.dateD, per.val)
}

func (e *etfDays) logKeepDays(per *etfDaysPer) {
	if e.keepDays <= 0 {
		return
	}
	fmt.Printf(" %s 持续%s %d天 反向日\n", e.log(per), e.isUpStr(e.starIsUp), e.keepDays)
}
