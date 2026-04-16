package tag_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/tag"
)

func TestPrintNoTags(t *testing.T) {
	r := tag.New()
	var buf bytes.Buffer
	p := tag.NewPrinter(&buf)
	p.Print(r)
	if !strings.Contains(buf.String(), "no tags") {
		t.Fatalf("expected empty message, got: %s", buf.String())
	}
}

func TestPrintHeaderPresent(t *testing.T) {
	r := tag.New()
	r.Set(80, []string{"http"})
	var buf bytes.Buffer
	p := tag.NewPrinter(&buf)
	p.Print(r)
	if !strings.Contains(buf.String(), "PORT") {
		t.Fatalf("expected header in output: %s", buf.String())
	}
}

func TestPrintTagsFormatted(t *testing.T) {
	r := tag.New()
	r.Set(22, []string{"ssh"})
	r.Set(443, []string{"https", "tls"})
	var buf bytes.Buffer
	p := tag.NewPrinter(&buf)
	p.Print(r)
	out := buf.String()
	if !strings.Contains(out, "ssh") {
		t.Errorf("expected ssh in output")
	}
	if !strings.Contains(out, "https, tls") {
		t.Errorf("expected https, tls in output")
	}
}

func TestPrintSortedByPort(t *testing.T) {
	r := tag.New()
	r.Set(9000, []string{"app"})
	r.Set(80, []string{"http"})
	var buf bytes.Buffer
	p := tag.NewPrinter(&buf)
	p.Print(r)
	out := buf.String()
	if strings.Index(out, "80") > strings.Index(out, "9000") {
		t.Error("expected port 80 before 9000")
	}
}
