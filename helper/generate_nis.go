package helper

import (
	"fmt"
	"sync"
	"time"
)

var (
	mu             sync.Mutex
	sequenceNumber = 1
	lastYear       = time.Now().Year()
)

func GenerateNis() string {
	mu.Lock()
	defer mu.Unlock()

	currentYear := time.Now().Year()
	// Reset sequenceNumber jika tahun telah berganti
	if currentYear != lastYear {
		sequenceNumber = 1
		lastYear = currentYear
	}

	// Mendapatkan dua digit terakhir dari tahun saat ini
	yearSuffix := currentYear % 100

	// Menghasilkan kode dengan format "MYYXXX"
	code := fmt.Sprintf("M%d%03d", yearSuffix, sequenceNumber)
	sequenceNumber++

	return code
}
