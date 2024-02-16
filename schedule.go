package tflux

import (
	"fmt"
	"math/rand"
	"time"
	"strings"
)

type TimeScale uint

const (
	MinuteScale TimeScale = iota
	HourScale
	DayScale
	WeekScale
	MonthScale
)

func checkTimeScaleValue(value TimeScale) error {
	switch value {
	case MinuteScale, HourScale, DayScale, WeekScale, MonthScale:
		return nil
	default:
		return fmt.Errorf("")
	}
}

// ---------------- RANDOM SCHEDULE ----------------------

type RandomSchedule struct {
	startTime time.Time
	nextEvent time.Time
	precEvent time.Time
	timeScale TimeScale
	bounds    struct {
		lower uint64
		upper uint64
	}
}

func NewRandomSchedule(
	startTime time.Time,
	tScale TimeScale,
	lBound, uBound int) (*RandomSchedule, error) {

	if uBound < lBound {
		return nil, fmt.Errorf(
			"the lower bound of a random time interval" +
				"must be less than the upper bound")
	}

	err := checkTimeScaleValue(tScale)
	if err != nil {
		return nil, err
	}

	rs := RandomSchedule{
		startTime: startTime,
		nextEvent: startTime,
		timeScale: tScale,
		bounds: struct {
			lower uint64
			upper uint64
		}{uint64(lBound), uint64(uBound)},
	}
	return &rs, nil
}

func (rs *RandomSchedule) setNextEvent() {
	r := rand.Uint64()
	u := rs.bounds.upper
	l := rs.bounds.lower
	delta := int(r%(u-l+1) + l)

	rs.precEvent = rs.nextEvent
	switch rs.timeScale {
	case MinuteScale:
		rs.nextEvent = rs.precEvent.Add(time.Minute * time.Duration(delta))
	case HourScale:
		rs.nextEvent = rs.precEvent.Add(time.Hour * time.Duration(delta))
	case DayScale:
		rs.nextEvent = rs.precEvent.AddDate(0, 0, delta)
	case WeekScale:
		rs.nextEvent = rs.precEvent.AddDate(0, 0, delta*7)
	case MonthScale:
		rs.nextEvent = rs.precEvent.AddDate(0, delta, 0)
	}
}

func (rs *RandomSchedule) GetNextEvent() time.Time {
	defer rs.setNextEvent()
	return rs.nextEvent
}

// -------------------- PATTERN SCHEDULE -------------------------

type PatternSchedule struct {
	startTime time.Time
	nextEvent time.Time
	precEvent time.Time
	timeScale TimeScale
	patternQueue []uint
	patternConst []uint
}

func NewPatternSchedule(startTime time.Time, tScale TimeScale, pattern []uint) (*PatternSchedule, error) {
	err := checkTimeScaleValue(tScale)
	if err != nil {
		return nil, err
	}

	ps := PatternSchedule{
		startTime: startTime,
		nextEvent: startTime,
		timeScale: tScale,
		patternQueue: make([]uint, 0, len(pattern)),
		patternConst: make([]uint, len(pattern)),
	}
	copy(ps.patternConst, pattern)
	return &ps, nil
}

func (ps *PatternSchedule) getNextDelta() int {
	if len(ps.patternQueue) > 0 {
		delta := ps.patternQueue[0]
		ps.patternQueue = ps.patternQueue[1:]
		return int(delta)
	}
	ps.patternQueue = append(ps.patternQueue, ps.patternConst...)
	return ps.getNextDelta()
}

func (ps *PatternSchedule) setNextEvent() {
	delta := ps.getNextDelta()
	ps.precEvent = ps.nextEvent
	switch ps.timeScale {
	case MinuteScale:
		ps.nextEvent = ps.precEvent.Add(time.Minute * time.Duration(delta))
	case HourScale:
		ps.nextEvent = ps.precEvent.Add(time.Hour * time.Duration(delta))
	case DayScale:
		ps.nextEvent = ps.precEvent.AddDate(0, 0, delta)
	case WeekScale:
		ps.nextEvent = ps.precEvent.AddDate(0, 0, delta*7)
	case MonthScale:
		ps.nextEvent = ps.precEvent.AddDate(0, delta, 0)
	}
}

func (ps *PatternSchedule) GetNextEvent() time.Time {
	defer ps.setNextEvent()
	return ps.nextEvent
}

// -------------------- CRONTAB SCHEDULE -------------------------

type CronSchedule struct {
	startTime time.Time
	nextEvent time.Time
	precEvent time.Time
	timeScale TimeScale
	patternQueue []uint
	patternConst []uint
}

func parseCronExpr(expr string) ([]Schedule, error) {
	cronAllowedSymbols := [][]string{
		[]string{"*", ",", "-", "/"},
		[]string{"*", ",", "-", "/"},
		[]string{"*", ",", "-", "/"},
		[]string{"*", ",", "-", "/"},
		[]string{"*", ",", "-", "/"},
	}
	exprArr := strings.Split(expr, " ")
	if len(exprArr) != 5 {
		return nil, fmt.Errorf("the cron expression must have 5 fields")
	}
	
	for index, exprField := range exprArr {
		if strings.Split(exprField, "*") == 
	}
}

func NewCronSchedule(startTime time.Time, cronExpression string) (*CronSchedule, error) {
	shceduleList, err := parseCronExpr(cronExpression)
	if err != nil {
		return nil, err
	}

	ps := PatternSchedule{
		startTime: startTime,
		nextEvent: startTime,
		timeScale: tScale,
		patternQueue: make([]uint, 0, len(pattern)),
		patternConst: make([]uint, len(pattern)),
	}
	copy(ps.patternConst, pattern)
	return &ps, nil
}

