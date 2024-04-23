package domain_test

import (
	"testing"

	"github.com/tkhrk1010/go-samples/actor-model/cluster/cluster-user-account/shared/domain"
)

func TestManager_Init(t *testing.T) {
	m := &domain.Manager{}
	m.Init()

	if m.AccountMap == nil {
		t.Error("AccountMap should be initialized")
	}
}

func TestManager_RegisterAccount(t *testing.T) {
	m := &domain.Manager{AccountMap: make(map[string]bool)}

	// アカウントを登録
	account := "test_account"
	m.RegisterAccount(account)

	// 登録したアカウントが存在するか確認
	if !m.AccountMap[account] {
		t.Errorf("Account '%s' should be registered", account)
	}
}

func TestManager_DeregisterAccount(t *testing.T) {
	m := &domain.Manager{AccountMap: make(map[string]bool)}

	// アカウントを登録
	account := "test_account"
	m.RegisterAccount(account)

	// アカウントを登録解除
	m.DeregisterAccount(account)

	// 登録解除したアカウントが存在しないか確認
	if m.AccountMap[account] {
		t.Errorf("Account '%s' should be deregistered", account)
	}
}