package entropy

import (
	"testing"
)

const (
	high         = "XZsVLgrqzy1wMabsl8TO9SuiKmOhWsz6qbBo6u8WhMDiLysEAG"
	expectedHigh = 4.963856189774723
	low          = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
)

func TestLowEntropy(t *testing.T) {
	t.Logf("Low entropy: %f", StringEntropy(low))
}

func TestHighEntropy(t *testing.T) {
	eh := StringEntropy(high)
	t.Logf("High entropy: %f", eh)
	if eh != expectedHigh {
		t.Logf("Entropy %f != %f", eh, expectedHigh)
		t.Fail()
	}
}
