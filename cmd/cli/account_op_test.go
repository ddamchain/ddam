package cli

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	defaultPassword = "123"
	defaultDuration = uint(120)
	exportKeyInfo   = "0xb0882314bec44ee5fd20a741cbb1167456b53728edd0b59f33129e561d39210e"
)

var (
	exKey  string
	exaddr string
)

func TestInitAccountManager(t *testing.T) {
	aop, err := initAccountManager("./keystore", false)
	if err != nil {
		t.Fatal("init account err:", err)
	}
	//NewAccount
	result := aop.NewAccount(defaultPassword)
	if result.Status != 0 {
		t.Fatalf("NewAccount err")
	}
	fmt.Println("NewAccount result:", result)

	//AccountList
	result = aop.AccountList()
	if result.Status != 0 {
		t.Fatalf("AccountList err")
	}
	fmt.Println("AccountList:：", result, reflect.TypeOf(result.Data))

	//unLock
	//accountList:=result.Data.([]string)
	exaddr = result.Data.([]string)[0]
	result = aop.UnLock(exaddr, defaultPassword, defaultDuration)
	if result.Status != 0 {
		t.Fatalf("unLock err")
	}
	fmt.Printf("UnLock addr :%s,result:%s", exaddr, result.Message)

	//AccountInfo
	result = aop.AccountInfo()
	if result.Status != 0 {
		t.Fatalf("AccountInfo err")
	}
	fmt.Println("AccountInfo :", result, reflect.TypeOf(result.Data))

	//ExprotKey
	result = aop.ExportKey()
	if result.Status != 0 {
		t.Fatalf("exprotKey err")
	}
	fmt.Println("exprotKey	:", result)
	exKey = result.Data.(string)

	//DeleteAccount
	result = aop.DeleteAccount()
	if result.Status != 0 {
		t.Fatalf("DeleteAccount err")
	}
	fmt.Println("DeleteAccount", result)

	//ImprotKey
	result = aop.NewAccountByImportKey(exKey, defaultPassword)
	if result.Status != 0 {
		t.Fatalf("ImprotKey err")
	}

	if result.Data != exaddr {
		t.Fatalf("Inconsistent private key export and export")
	}
	fmt.Println("ImprotKey :", result)

}
