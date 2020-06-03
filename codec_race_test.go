// +build race

package ep

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/advanderveer/ep/epcoding"
)

func TestConcurrentCodecHandling(t *testing.T) {
	for i, c := range []struct {
		fn      interface{}
		body    string
		expCode int
		expBody string
	}{
		{func() {}, ``, 200, ``},
		{func(s string) string { return s }, `"foo"`, 200, `"foo"` + "\n"},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			srv := httptest.NewServer(New(
				RequestDecoding(epcoding.JSON{}),
				ResponseEncoding(epcoding.JSON{}),
			).Handle(c.fn))

			var wg sync.WaitGroup
			defer wg.Wait()

			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					resp, _ := http.Post(srv.URL, "application/json", strings.NewReader(c.body))
					if resp.StatusCode != c.expCode {
						t.Fatalf("expected code: %d, got: %d", c.expCode, resp.StatusCode)
					}

					body, _ := ioutil.ReadAll(resp.Body)
					if string(body) != c.expBody {
						t.Fatalf("expected body: '%s', got: '%s'", c.expBody, string(body))
					}
				}()
			}
		})
	}
}
