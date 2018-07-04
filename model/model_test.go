package model

import (
	"testing"
)

func Test_Model(t *testing.T) {
	role := Role{}
	role.table = "role"

	role.List()
	t.Log("test role model")
}
