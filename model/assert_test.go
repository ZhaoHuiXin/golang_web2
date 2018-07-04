package model

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func BenchmarkAssertAssert(b *testing.B) {
	count := 0
	var delta interface{}
	delta = 1
	for n := 0; n < b.N; n++ {
		count = count + delta.(int)
	}
}

func BenchmarkAssertDirect(b *testing.B) {
	count := 0
	delta := 1
	for n := 0; n < b.N; n++ {
		count = count + delta
	}
}

type T int

func (t T) Value() int {
	return int(t)
}

type I interface{ Value() int }

func BenchmarkAssertMethod(b *testing.B) {
	count := 0
	var delta I = T(1)
	for n := 0; n < b.N; n++ {
		count = count + delta.Value()
	}
}

func BenchmarkAssertSumMethod(b *testing.B) {
	var sum = func(x, y interface{}) int {
		return x.(int) + y.(int)
	}

	count := 0
	delta := 1
	for n := 0; n < b.N; n++ {
		count = sum(count, delta)
	}
}

func BenchmarkAssertSwitchType(b *testing.B) {
	count := 0
	var delta interface{} = 1
	for n := 0; n < b.N; n++ {
		switch v := delta.(type) {
		case int, int32, int64:
			count = count + v.(int)
		default:
		}
	}
}

type A struct {
	Name, N1, N2, N3, N4, N5 string
}

type B struct {
	Name, N1, N2, N3, N4, N5 string
}

// dst should be a pointer to struct, src should be a struct
func CopyReflect(dst interface{}, src interface{}) (err error) {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		err = errors.New("dst isn't a pointer to struct")
		return
	}
	dstElem := dstValue.Elem()
	if dstElem.Kind() != reflect.Struct {
		err = errors.New("pointer doesn't point to struct")
		return
	}

	srcValue := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Struct {
		err = errors.New("src isn't struct")
		return
	}

	for i := 0; i < srcType.NumField(); i++ {
		sf := srcType.Field(i)
		sv := srcValue.FieldByName(sf.Name)
		// make sure the value which in dst is valid and can set
		if dv := dstElem.FieldByName(sf.Name); dv.IsValid() && dv.CanSet() {
			dv.Set(sv)
		}
	}
	return
}

func CopyA2B(a *A, b *B) {
	b.Name = a.Name
	b.N1 = a.N1
	b.N2 = a.N2
	b.N3 = a.N3
	b.N4 = a.N4
	b.N5 = a.N5
}

func CopyA2Map(a *A) map[string]interface{} {
	b := make(map[string]interface{})
	b["Name"] = a.Name
	b["N1"] = a.N1
	b["N2"] = a.N2
	b["N3"] = a.N3
	b["N4"] = a.N4
	b["N5"] = a.N5
	return b
}

func BenchmarkCopyReflect(b *testing.B) {
	for n := 0; n < b.N; n++ {
		a := A{Name: "a", N1: "n1", N2: "n2", N3: "n3", N4: "n4", N5: "n5"}
		d := B{}
		CopyReflect(&d, a)
	}
}

func BenchmarkCopyA2B(b *testing.B) {
	for n := 0; n < b.N; n++ {
		a := A{Name: "a", N1: "n1", N2: "n2", N3: "n3", N4: "n4", N5: "n5"}
		d := B{}
		CopyA2B(&a, &d)
	}
}

func BenchmarkCopyMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		a := A{Name: "a", N1: "n1", N2: "n2", N3: "n3", N4: "n4", N5: "n5"}
		d := CopyA2Map(&a)
		_ = d
	}
}

func BenchmarkCopyA2BByJSON(b *testing.B) {
	for n := 0; n < b.N; n++ {
		a := A{Name: "a", N1: "n1", N2: "n2", N3: "n3", N4: "n4", N5: "n5"}
		buf, _ := json.Marshal(a)
		d := B{}
		json.Unmarshal(buf, &d)
	}
}
