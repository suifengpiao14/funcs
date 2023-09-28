package funcs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetIp(t *testing.T) {
	ip, err := GetIp()
	require.NoError(t, err)
	fmt.Println(ip)
}
