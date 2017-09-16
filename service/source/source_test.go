package source_test

import (
	"github.com/koko990/logcollector/service/source"
	"testing"
)

func TestKuberResource_FetchResource(t *testing.T) {
	var s = source.NewSourceFactory("http://10.110.18.26:8080")
	s.FetchResource()
}
