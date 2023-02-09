package auth

import (
	"fmt"
	"strings"
	"testing"
)

func TestPCKE(t *testing.T) {
	t.Skip()
	result, _ := generateRandomString(128)
	fmt.Println(string(result))
	randomString := "21d42492841a8eb0a0d9841ef1112579f17126b3508d0b018a9fa498"
	expeted := "tCclVDIuqLhiOb9vpn2YBsa_zYJ77s1hh0Hcq7CIKdo"
	actual := sha256UrlEncode([]byte(randomString))
	if strings.Compare(expeted, actual) == 0 {
		t.Errorf("actual: %s, expected: %s", actual, expeted)
	}
}
