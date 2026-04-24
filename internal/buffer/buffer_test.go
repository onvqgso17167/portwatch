package buffer

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestAddAndLen(t *testing.T) {
	b := New(5)
	if b.Len() != 0 {
		t.Fatalf("expected 0, got %d", b.Len())
	}
	b.Add(epoch, []string{"80", "443"})
	if b.Len() != 1 {
		t.Fatalf("expected 1, got %d", b.Len())
	}
}

func TestEvictsOldestWhenFull(t *testing.T) {
	b := New(3)
	b.Add(epoch, []string{"80"})
	b.Add(epoch.Add(time.Second), []string{"443"})
	b.Add(epoch.Add(2*time.Second), []string{"8080"})
	b.Add(epoch.Add(3*time.Second), []string{"9090"})

	if b.Len() != 3 {
		t.Fatalf("expected 3, got %d", b.Len())
	}
	all := b.All()
	if all[0].Results[0] != "443" {
		t.Errorf("expected oldest to be 443, got %s", all[0].Results[0])
	}
	if all[2].Results[0] != "9090" {
		t.Errorf("expected newest to be 9090, got %s", all[2].Results[0])
	}
}

func TestAllReturnsCopy(t *testing.T) {
	b := New(5)
	b.Add(epoch, []string{"22"})
	all := b.All()
	all[0].Results[0] = "mutated"

	original := b.All()
	if original[0].Results[0] == "mutated" {
		t.Error("All() should return an independent copy")
	}
}

func TestReset(t *testing.T) {
	b := New(5)
	b.Add(epoch, []string{"80"})
	b.Add(epoch, []string{"443"})
	b.Reset()
	if b.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", b.Len())
	}
}

func TestCapMinimumIsOne(t *testing.T) {
	b := New(0)
	b.Add(epoch, []string{"80"})
	b.Add(epoch, []string{"443"})
	if b.Len() != 1 {
		t.Fatalf("expected 1, got %d", b.Len())
	}
	all := b.All()
	if all[0].Results[0] != "443" {
		t.Errorf("expected 443, got %s", all[0].Results[0])
	}
}

func TestConcurrentAdd(t *testing.T) {
	b := New(100)
	done := make(chan struct{})
	for i := 0; i < 50; i++ {
		go func() {
			b.Add(epoch, []string{"80"})
			done <- struct{}{}
		}()
	}
	for i := 0; i < 50; i++ {
		<-done
	}
	if b.Len() > 100 {
		t.Errorf("buffer exceeded capacity: %d", b.Len())
	}
}
