package window

import (
	"fmt"
	"time"

	"github.com/mad01/k8s-node-terminator/pkg/annotations"
	"k8s.io/api/core/v1"
)

// ParseTest test input time window
// expected format is h:mm AM/PM
//		example: 2:01 AM
func ParseTest(input string) error {
	_, err := time.Parse("3:04 PM", input)
	if err != nil {
		return fmt.Errorf("failed to parse input \"%v\" time %v", input, err.Error())
	}
	return nil
}

func getTimeWindow(hour, minute, second int) (*time.Time, error) {
	current := time.Now()

	hourStr := fmt.Sprintf("%d", hour)
	minuteStr := fmt.Sprintf("%d", minute)
	secondStr := fmt.Sprintf("%d", second)

	if hour <= 9 {
		hourStr = fmt.Sprintf("0%d", hour)
	}

	if minute <= 9 {
		minuteStr = fmt.Sprintf("0%d", minute)
	}

	if second <= 9 {
		secondStr = fmt.Sprintf("0%d", second)
	}

	timeString := fmt.Sprintf("%d-%d-%d %v:%v:%v",
		current.Year(),
		current.Month(),
		current.Day(),
		hourStr,
		minuteStr,
		secondStr,
	)

	w, err := time.Parse("2006-01-02 15:04:05", timeString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timeString%v %v", timeString, err.Error())
	}
	return &w, nil

}

// GetMaintenanceWindowFromAnnotations doc
func GetMaintenanceWindowFromAnnotations(node *v1.Node) (*MaintenanceWindow, error) {
	var from, to string
	var window *MaintenanceWindow

	a := node.GetAnnotations()
	if _, ok := a[annotations.NodeAnnotationFromWindow]; ok {
		from = a[annotations.NodeAnnotationFromWindow]
	}
	if _, ok := a[annotations.NodeAnnotationToWindow]; ok {
		to = a[annotations.NodeAnnotationToWindow]
	}

	if from != "" && to != "" {
		w, err := NewMaintenanceWindow(from, to)
		if err != nil {
			return nil, fmt.Errorf("failed to create new MaintenanceWindow %v", err.Error())
		}
		window = w
	}

	return window, nil
}

// NewMaintenanceWindow into
func NewMaintenanceWindow(from, to string) (*MaintenanceWindow, error) {
	fromTime, err := time.Parse("3:04 PM", from)
	if err != nil {
		return nil, fmt.Errorf("failed to parse from time %v", err.Error())
	}
	toTime, err := time.Parse("3:04 PM", to)
	if err != nil {
		return nil, fmt.Errorf("failed to parse to time %v", err.Error())
	}

	start, err := getTimeWindow(fromTime.Hour(), fromTime.Minute(), fromTime.Second())
	if err != nil {
		return nil, fmt.Errorf("failed to get time window %v", err.Error())
	}

	stop, err := getTimeWindow(toTime.Hour(), toTime.Minute(), toTime.Second())
	if err != nil {
		return nil, fmt.Errorf("failed to get time window %v", err.Error())
	}

	m := MaintenanceWindow{
		fromString: from,
		toString:   to,
		from:       start,
		to:         stop,
	}
	return &m, nil

}

// MaintenanceWindow info
type MaintenanceWindow struct {
	fromString string
	toString   string
	from       *time.Time
	to         *time.Time
}

// FromString returns the input string hh:mm AM/PM
func (m *MaintenanceWindow) FromString() string {
	return m.fromString
}

// ToString returns the input string hh:mm AM/PM
func (m *MaintenanceWindow) ToString() string {
	return m.toString
}

// From returns from window time
func (m *MaintenanceWindow) From() *time.Time {
	return m.from
}

// To returns to window time
func (m *MaintenanceWindow) To() *time.Time {
	return m.to
}

// InMaintenanceWindow info
func (m *MaintenanceWindow) InMaintenanceWindow() bool {
	currentTimeUnix := time.Now().Unix()
	from := m.from.Unix()
	to := m.to.Unix()

	inWindow := func() bool {
		if (currentTimeUnix >= from) && (currentTimeUnix <= to) {
			return true
		}

		return false
	}

	if inWindow() == true {
		return true
	}

	return false
}
