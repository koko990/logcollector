package calculate

type Unit int

const (
	_      Unit = iota
	minute
	hour
	day
)

type CalcParam struct {
	CurrentVal       int
	PrevVal          int
	timeUnit         Unit
	CurrentTimestamp int
	PrevTimestamp    int
}

//formula of moving average
//aver(a1~2)=(a1+a2)/2
//aver(a1~3)=(a1+a2+a3)/3
//movAver(a1~2)=(n1*a1+n2*a2)/(n) n1=1 n2=1 n=2
//movAver(a1~3)=(n1*movAver(a1~2)+n2*a3)/(n) = (a1+a2+a3)/3  n1=2 n2=3-2  n=3
func calcAverMoving(val int64, prev int64, w int64, sum int64) (res int64) {
	return (w*prev + val*(sum-w)) / sum
}

//
func (c CalcParam) getTimeAver() int {
	switch c.timeUnit {
	case minute:
		return int(calcAverMoving(int64(c.CurrentVal), int64(c.PrevVal),
			int64(c.PrevTimestamp%(60)), int64(c.CurrentTimestamp)))
	case hour:
		return int(calcAverMoving(int64(c.CurrentVal), int64(c.PrevVal),
			int64(c.PrevTimestamp%(60*60)), int64(c.CurrentTimestamp)))
	case day:
		return int(calcAverMoving(int64(c.CurrentVal), int64(c.PrevVal),
			int64(c.PrevTimestamp%(60*60*24)), int64(c.CurrentTimestamp)))
	}
	return 0
}

//
func (c CalcParam) GetMinuteAver() int {
	c.timeUnit = minute
	return c.getTimeAver()
}

//
func (c CalcParam) GetHourAver() int {
	c.timeUnit = hour
	return c.getTimeAver()
}

//
func (c CalcParam) GetDayAver() int {
	c.timeUnit = day
	return c.getTimeAver()
}
