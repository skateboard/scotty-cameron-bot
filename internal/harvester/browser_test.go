package harvester

import (
	"fmt"
	"testing"
	"time"
)

func TestHarvester(t *testing.T) {
	harv := New("harvester", "https://www.scottycameron.com/")
	go func() {
		err := harv.Initialize()
		if err != nil {
			t.Error(err)
		}
	}()

	time.Sleep(10 * time.Second)

	solv, err := harv.Harvest("51829642-2cda-4b09-896c-594f89d700cc")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(solv)

	harv.Destroy()
}
