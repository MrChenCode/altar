package test

import (
	"altar/application/config"
	"altar/application/library/util"
	"github.com/stretchr/testify/require"
	"testing"
	"unsafe"
)

type boolFuncStruct struct {
	id   int
	name string
}

type boolOtherType string

func TestUtilInt(t *testing.T) {
	require.Equal(t, util.Int(int(100)), int(100))
	require.Equal(t, util.Int("abc"), int(0))
	require.Equal(t, util.Int("22a3bc33"), int(0))
	require.Equal(t, util.Int(""), int(0))
	require.Equal(t, util.Int("\n"), int(0))
	require.Equal(t, util.Int(float64(5.6444)), int(0))
	require.Equal(t, util.Int(float64(5)), int(5))
	require.Equal(t, util.Int([]int{123}), int(0))
	require.Equal(t, util.Int(map[string]string{"abc": "ef"}), int(0))
	require.Equal(t, util.Int(&boolFuncStruct{}), int(0))
}

func TestUtilBool(t *testing.T) {
	require.True(t, util.Bool(true))
	require.False(t, util.Bool(false))

	require.False(t, util.Bool(int(0)))
	require.True(t, util.Bool(int(1)))
	require.False(t, util.Bool(int8(0)))
	require.True(t, util.Bool(int8(100)))
	require.False(t, util.Bool(int16(0)))
	require.True(t, util.Bool(int16(12)))
	require.False(t, util.Bool(int32(0)))
	require.True(t, util.Bool(int32(5)))
	require.False(t, util.Bool(int64(0)))
	require.True(t, util.Bool(int64(-1)))

	require.False(t, util.Bool(uint(0)))
	require.True(t, util.Bool(uint(1)))
	require.False(t, util.Bool(uint8(0)))
	require.True(t, util.Bool(uint8(100)))
	require.False(t, util.Bool(uint16(0)))
	require.True(t, util.Bool(uint16(12)))
	require.False(t, util.Bool(uint32(0)))
	require.True(t, util.Bool(uint32(5)))
	require.False(t, util.Bool(uint64(0)))
	require.True(t, util.Bool(uint64(10)))

	a := "abc"
	ptr := uintptr(unsafe.Pointer(&a))
	require.True(t, util.Bool(ptr))
	require.False(t, util.Bool(uintptr(0)))

	require.False(t, util.Bool(float32(0.0)))
	require.False(t, util.Bool(float32(0)))
	require.True(t, util.Bool(float32(3.1415926)))
	require.False(t, util.Bool(float64(0.000)))
	require.False(t, util.Bool(float64(0)))
	require.True(t, util.Bool(float64(1.232)))

	require.True(t, util.Bool(complex64(complex(1, 2.3))))
	require.False(t, util.Bool(complex64(complex(0, 0))))
	require.True(t, util.Bool(complex64(complex(0, 1.2))))
	require.True(t, util.Bool(complex64(complex(1, 0.0))))
	require.True(t, util.Bool(complex128(complex(1, 2.3))))
	require.False(t, util.Bool(complex128(complex(0, 0))))
	require.True(t, util.Bool(complex128(complex(0, 1.3))))
	require.True(t, util.Bool(complex128(complex(1, 0.0))))

	require.True(t, util.Bool("a"))
	require.True(t, util.Bool(" a "))
	require.False(t, util.Bool(""))
	require.False(t, util.Bool(" "))
	require.False(t, util.Bool("\n\t"))

	var ch chan int
	require.False(t, util.Bool(ch))
	ch = make(chan int)
	require.True(t, util.Bool(ch))
	var fs func()
	require.False(t, util.Bool(fs))
	fs = func() {}
	require.True(t, util.Bool(fs))
	var m map[string]int
	require.False(t, util.Bool(m))
	m = make(map[string]int)
	require.False(t, util.Bool(m))
	m["a"] = 1
	require.True(t, util.Bool(m))
	delete(m, "a")
	require.False(t, util.Bool(m))
	var up unsafe.Pointer
	require.False(t, util.Bool(up))
	var upt = "abc"
	up = unsafe.Pointer(&upt)
	require.True(t, util.Bool(up))
	var ss []int
	require.False(t, util.Bool(ss))
	ss = make([]int, 0, 1)
	require.False(t, util.Bool(ss))
	ss = append(ss, 1)
	require.True(t, util.Bool(ss))

	var intptr *int
	require.False(t, util.Bool(intptr))
	inttmp := 100
	intptr = &inttmp
	require.True(t, util.Bool(&intptr))
	inttmp2 := 0
	intptr2 := &inttmp2
	require.False(t, util.Bool(&intptr2))

	var intera interface{}
	var interb interface{}
	interb = &intera
	require.False(t, util.Bool(&interb))
	for _, ix := range []interface{}{0, "", false, "\t\n", nil, []int{}, map[string]int{}} {
		intera = ix
		require.False(t, util.Bool(&interb))
	}
	for _, ix := range []interface{}{1, "0", true, "\ta\n", []int{1}} {
		intera = ix
		require.True(t, util.Bool(&interb))
	}
	var structa boolFuncStruct
	require.True(t, util.Bool(&structa))
	require.True(t, util.Bool(structa))
	require.True(t, util.Bool(&boolFuncStruct{}))

	var structb *boolFuncStruct
	require.False(t, util.Bool(structb))
	require.False(t, util.Bool(&structb))

	var other boolOtherType
	require.False(t, util.Bool(other))
	require.False(t, util.Bool(&other))
	other = "a"
	require.True(t, util.Bool(other))
	require.True(t, util.Bool(&other))
}

func TestUtilSubstr(t *testing.T) {
	str := "0123456789"

	require.Equal(t, util.Substr(str, 0, 5), "01234")
	require.Equal(t, util.Substr(str, 3, 5), "34567")
	require.Equal(t, util.Substr(str, 7, 15), "789")
	require.Equal(t, util.Substr(str, 12, 5), "")
	require.Equal(t, util.Substr(str, 2, -2), "")
	require.Equal(t, util.Substr(str, 4, 0), "456789")
	require.Equal(t, util.Substr(str, -3, 0), "789")
	require.Equal(t, util.Substr(str, -4, 10), "6789")
}

func TestUtilBookType(t *testing.T) {
	require.Equal(t, util.GetBookType("012"), config.BookTypeDefault)
	require.Equal(t, util.GetBookType("0123"), config.BookTypeDefault)
	require.Equal(t, util.GetBookType("32535"), config.BookTypeDefault)
	require.Equal(t, util.GetBookType(234699), config.BookTypeComic)
	require.Equal(t, util.GetBookType("34546698"), config.BookTypeAudio)
	require.Equal(t, util.GetBookType("12345657"), config.BookTypeDefault)
	require.Equal(t, util.GetBookType(234564643699), config.BookTypeComic)
	require.Equal(t, util.GetBookType("34546899"), config.BookTypeEpub)
	require.Equal(t, util.GetBookType(2435465900), config.BookTypeEpub)
}

func TestUtilFbkImg(t *testing.T) {
	require.Equal(t, util.FbkImg("abc"), config.CdnDomainRoute+"abc")
	require.Equal(t, util.FbkImg("comic/head.jpg"), config.CdnDomainRoute+"comic/head.jpg")
	require.Equal(t, util.FbkImg(""), "")
	require.Equal(t, util.FbkImg("https://abc.com/s.jpg"), "https://abc.com/s.jpg")
	require.Equal(t, util.FbkImg("http://a.com/b.png"), "http://a.com/b.png")
}
