package classify_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/classify"
	"github.com/user/portwatch/internal/scanner"
)

func makeResult(port int) scanner.Result {
	return scanner.Result{Port: port, Open: true, Timestamp: time.Now()}
}

func TestClassifyPrivilegedPortIsWarning(t *testing.T) {
	c := classify.New(nil)
	l := c.Classify(makeResult(80))
	if l != classify.LevelWarning {
		t.Fatalf("expected warning, got %s", l)
	}
}

func TestClassifyHighPortIsInfo(t *testing.T) {
	c := classify.New(nil)
	l := c.Classify(makeResult(8080))
	if l != classify.LevelInfo {
		t.Fatalf("expected info, got %s", l)
	}
}

func TestClassifyCriticalPortOverridesDefault(t *testing.T) {
	c := classify.New([]int{8080})
	l := c.Classify(makeResult(8080))
	if l != classify.LevelCritical {
		t.Fatalf("expected critical, got %s", l)
	}
}

func TestClassifyCriticalPrivilegedPort(t *testing.T) {
	c := classify.New([]int{22})
	l := c.Classify(makeResult(22))
	if l != classify.LevelCritical {
		t.Fatalf("expected critical, got %s", l)
	}
}

func TestClassifyAllReturnsMapForAllResults(t *testing.T) {
	c := classify.New([]int{443})
	results := []scanner.Result{makeResult(443), makeResult(8080), makeResult(22)}
	m := c.ClassifyAll(results)
	if m[443] != classify.LevelCritical {
		t.Errorf("443 should be critical")
	}
	if m[8080] != classify.LevelInfo {
		t.Errorf("8080 should be info")
	}
	if m[22] != classify.LevelWarning {
		t.Errorf("22 should be warning")
	}
}

func TestLevelString(t *testing.T) {
	cases := map[classify.Level]string{
		classify.LevelInfo:     "info",
		classify.LevelWarning:  "warning",
		classify.LevelCritical: "critical",
	}
	for l, want := range cases {
		if l.String() != want {
			t.Errorf("got %s, want %s", l.String(), want)
		}
	}
}
