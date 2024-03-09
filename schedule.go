package tflux

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"
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
	startTime    time.Time
	nextEvent    time.Time
	precEvent    time.Time
	timeScale    TimeScale
	patternQueue []uint
	patternConst []uint
}

func NewPatternSchedule(startTime time.Time, tScale TimeScale, pattern []uint) (*PatternSchedule, error) {
	err := checkTimeScaleValue(tScale)
	if err != nil {
		return nil, err
	}

	ps := PatternSchedule{
		startTime:    startTime,
		nextEvent:    startTime,
		timeScale:    tScale,
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

var cronExprFieldNames map[int]string = map[int]string{
	0: "minute",
	1: "hour",
	2: "monthDay",
	3: "month",
	4: "weekDay",
} 

type CronSchedule struct {
	startTime     time.Time
	scheduleTable map[string]struct{
		units []int
		next int
	}
}

func validate(unit, position int) error {
	doCheck := func(lb, ub int, unitName string) error {
		if unit < lb || unit > ub {
			return fmt.Errorf("invalid %s, must between %d and %d inclusive", unitName, lb, ub)
		}
		return nil
	}
	switch position {
	case 0: // minute
		return doCheck(0, 59, "minutes")
	case 1: // hour
		return doCheck(0, 23, "hours")
	case 2: // day of month
		return doCheck(1, 31, "day of month")
	case 3: // month
		return doCheck(1, 12, "month")
	case 4: // days of week
		return doCheck(1, 7, "day of week")
	}
	return nil
}

func parseField(position int, field string) ([]int, error) {
	dedupe := func(list []int) []int {
		encountered := make(map[int]bool)
		result := make([]int, 0)
		for _, i := range list {
			if encountered[i] {
				continue
			}
			encountered[i] = true
			result = append(result, i)
		}
		return result
	}
	sequence := func(from, to int) []int {
		result := make([]int, 0)
		for i := from; i <= to; i++ {
			result = append(result, i)
		}
		return result
	}	
	result := make([]int, 0)
	elements := strings.Split(field, ",")
	for _, element := range elements {
		if element == "*" {
			switch position {
			case 0:
				return sequence(0, 59), nil
			case 1:
				return sequence(0, 23), nil
			case 2:
				return sequence(1, 31), nil
			case 3:
				return sequence(1, 12), nil
			case 4:
				return sequence(1, 7), nil
			}
		}
		bounds := strings.Split(element, "-")
		if len(bounds) > 1 {
			lb, err := strconv.Atoi(bounds[0])
			if err != nil {
				return []int{}, err
			}
			ub, err := strconv.Atoi(bounds[1])
			if err != nil {
				return []int{}, err
			}
			err = validate(lb, position)
			if err != nil {
				return []int{}, err
			}
			err = validate(ub, position)
			if err != nil {
				return []int{}, err
			}
			if lb > ub {
				return []int{}, fmt.Errorf(
					"lower bound is larger than upper bound in range %d-%d", lb, ub,
				)
			}
			result = append(result, sequence(lb, ub)...)
		} else {
			unit, err := strconv.Atoi(element)
			if err != nil {
				return []int{}, err
			}
			err = validate(unit, position)
			if err != nil {
				return []int{}, err
			}
			result = append(result, unit)
		}
	}
	result = dedupe(result)
	slices.Sort(result)
	return result, nil
}

//startTime time.Time
func NewCronSchedule(cronExpression string) (*CronSchedule, error) {
	crs := CronSchedule{startTime: time.Now()}
	err := crs.parseCronExpr(cronExpression)
	if err != nil {
		return nil, err
	}
	return &crs, nil
}

func (crs *CronSchedule) String() string {
	s := ""
	for fieldName, times := range crs.scheduleTable {
		s += fmt.Sprintf("%s: %v\n", fieldName, times.units)
	}
	return s
} 

func (crs *CronSchedule) parseCronExpr(expr string) error {
	exprArr := strings.Split(expr, " ")
	if len(exprArr) != 5 {
		return fmt.Errorf("a cron expression has five string fields, ")
	}
	crs.scheduleTable = make(map[string]struct{units []int; next int})
	for position, field := range exprArr {
		fieldName := cronExprFieldNames[position]
		scheduleUnits, err := parseField(position, field)
		if err != nil {
			return err
		}
		crs.scheduleTable[fieldName] = struct{units []int; next int}{
			units: scheduleUnits, 
			next: 0,
		}
	}
	return nil
}
