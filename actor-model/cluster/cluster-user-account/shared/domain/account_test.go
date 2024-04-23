package domain_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
)

func TestAccount_Add(t *testing.T) {
	// テストケースを定義
	cases := []struct {
		name          string
		initialCount  int64
		addAmount     int64
		expectedCount int64
	}{
		{"正の値を加算", 100, 50, 150},
		{"負の値を加算", 100, -30, 70},
		{"ゼロを加算", 100, 0, 100},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			acc := &domain.Account{
				Name:  "Test Account",
				Count: c.initialCount,
			}

			acc.Add(c.addAmount)

			if acc.Count != c.expectedCount {
				t.Errorf("Expected count: %d, but got: %d", c.expectedCount, acc.Count)
			}
		})
	}
}

func TestAccount_Remove(t *testing.T) {
	// テストケースを定義
	cases := []struct {
		name          string
		initialCount  int64
		removeAmount  int64
		expectedCount int64
	}{
		{"正の値を減算", 100, 30, 70},
		{"負の値を減算", 100, -50, 150},
		{"ゼロを減算", 100, 0, 100},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			acc := &domain.Account{
				Name:  "Test Account",
				Count: c.initialCount,
			}

			acc.Remove(c.removeAmount)

			if acc.Count != c.expectedCount {
				t.Errorf("Expected count: %d, but got: %d", c.expectedCount, acc.Count)
			}
		})
	}
}