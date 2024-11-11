package core

import (
	"fmt"
	"math"
)

type etfDays struct {
	all          []*etfDaysPer
	pinCa        float64     //关键点变化幅度 6%
	turnCa       float64     //逆转信号幅度 3%
	pin1         *etfDaysPer //关键点1
	pin2         *etfDaysPer //关键点2
	starIsUp     bool
	lastPin      *etfDaysPer
	keepDays     int
	keepTurnDays int
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
	if e.pin1.val != 0 {
		if per.val > e.pin1.val {
			e.logKeepDays(per, e.starIsUp)
			if math.Abs(per.val-e.pin1.val)/e.pin1.val*100 > e.turnCa {
				fmt.Printf(" %s 彻底突破【%s】%s趋势形成\n", e.log(per), e.pin(1), e.isUpStr(e.starIsUp))
				e.pin1 = &etfDaysPer{}
				e.pin2 = &etfDaysPer{}
			} else {
				fmt.Printf(" %s 突破【%s】%s\n", e.log(per), e.pin(1), e.isUpStr(e.starIsUp))
			}
		} else {
			if e.lastPin.val > per.val {
				e.logKeepTurnDays(per, !e.starIsUp)
				if caRate >= e.pinCa {
					e.starIsUp = false
					fmt.Print(fmt.Sprintf(" %s 于[%s]%s>= %.0f点(%.2f)->次级%s\n", e.log(per), e.log(e.lastPin), e.isUpStr(!e.starIsUp), e.pinCa, caRate, e.isUpTmpStr(e.starIsUp)))
					e.lastPin = per
				} else if caRate >= e.turnCa {
					fmt.Print(fmt.Sprintf(" %s 于[%s]%s>= %.0f点(%.2f)【逆转警告】\n", e.log(per), e.log(e.lastPin), e.isUpStr(!e.starIsUp), e.turnCa, caRate))
				}
			} else {
				e.logKeepDays(per, e.starIsUp)
				e.lastPin = per
			}
		}
	} else {
		if e.lastPin.val > per.val {
			e.logKeepTurnDays(per, !e.starIsUp)
			if caRate >= e.pinCa {
				e.starIsUp = false
				e.pin1 = e.lastPin
				e.lastPin = per
				e.keepTurnDays, e.keepDays = e.keepDays, e.keepTurnDays
				fmt.Printf(" %s 出现【%s】->自然%s阶段 \n", e.log(per), e.pin(1), e.isUpTmpStr(e.starIsUp))
				return
			}
		} else {
			e.logKeepDays(per, e.starIsUp)
			e.lastPin = per
		}
	}
}

func (e *etfDays) downThink(per *etfDaysPer, caRate float64) {
	if e.pin2.val != 0 {
		if per.val < e.pin2.val {
			e.logKeepDays(per, e.starIsUp)
			if math.Abs(per.val-e.pin2.val)/e.pin1.val*100 > e.turnCa {
				fmt.Print(fmt.Sprintf(" %s 突破【%s】%s趋势形成\n", e.log(per), e.pin(2), e.isUpStr(e.starIsUp)))
				e.pin1 = &etfDaysPer{}
				e.pin2 = &etfDaysPer{}
			} else {
				fmt.Printf(" %s 突破【%s】%s\n", e.log(per), e.pin(2), e.isUpStr(e.starIsUp))
			}
		} else {
			if per.val > e.lastPin.val {
				e.logKeepTurnDays(per, !e.starIsUp)
				if caRate >= e.pinCa {
					e.starIsUp = true
					fmt.Print(fmt.Sprintf(" %s 于[%s]%s>= %.0f点(%.2f)->次级%s\n", e.log(per), e.log(e.lastPin), e.isUpStr(!e.starIsUp), e.pinCa, caRate, e.isUpTmpStr(e.starIsUp)))
					e.lastPin = per
				} else if caRate >= e.turnCa {
					fmt.Print(fmt.Sprintf(" %s 于[%s]%s>= %.0f点(%.2f)【逆转警告】\n", e.log(per), e.log(e.lastPin), e.isUpStr(!e.starIsUp), e.turnCa, caRate))
				}
			} else {
				e.logKeepDays(per, e.starIsUp)
				e.lastPin = per
			}
		}
	} else {
		if e.lastPin.val < per.val {
			e.logKeepTurnDays(per, !e.starIsUp)
			if caRate >= e.pinCa {
				e.starIsUp = true
				e.pin2 = e.lastPin
				e.lastPin = per
				fmt.Printf(" %s 出现【%s】->自然%s阶段\n", e.log(per), e.pin(2), e.isUpTmpStr(e.starIsUp))
				e.keepTurnDays, e.keepDays = e.keepDays, e.keepTurnDays
				return
			}
		} else {
			e.logKeepDays(per, e.starIsUp)
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

// logKeepDays 确实同向增长
func (e *etfDays) logKeepDays(per *etfDaysPer, isUp bool) {
	//if e.isTurn {
	//	e.isTurn = false
	//	e.keepDays = 0
	//}
	e.keepTurnDays = 0
	e.keepDays++
	fmt.Printf(" %s 于[%s]持续[%s] %d天\n", e.log(per), e.log(e.lastPin), e.isUpStr(isUp), e.keepDays)
}

// logKeepTurnDays 确实反响增长
func (e *etfDays) logKeepTurnDays(per *etfDaysPer, isUp bool) {
	e.keepTurnDays++
	e.keepDays = 0
	fmt.Printf(" %s 于[%s]持续-%s- %d天\n", e.log(per), e.log(e.lastPin), e.isUpStr(isUp), e.keepTurnDays)
}

func (e *etfDays) pin(pinK int) string {
	if pinK == 1 {
		return fmt.Sprintf("关键点1 %s %.4f", e.pin1.dateD, e.pin1.val)
	} else {
		return fmt.Sprintf("关键点2 %s %.4f", e.pin2.dateD, e.pin2.val)
	}
}
