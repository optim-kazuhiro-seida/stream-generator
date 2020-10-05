package main

import (
	"bytes"
	"flag"
	"go/build"
	"go/format"
	"go/types"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/loader"
)

type Variable struct {
	TypeName    string
	PackageName string
}

func scanPkg(dir string) *loader.PackageInfo {
	p, err := build.ImportDir(dir, build.FindOnly)
	if err != nil {
		panic(err)
	}
	conf := loader.Config{TypeChecker: types.Config{FakeImportC: true}}
	conf.Import(p.ImportPath)
	program, err := conf.Load()
	if err != nil {
		panic(err)
	}
	return program.Package(p.ImportPath)
}
func main() {
	var (
		typeName = flag.String("type", "", "struct name, -type={{Struct}}")
		dir, err = filepath.Abs(".")
		buf      bytes.Buffer
	)
	flag.Parse()
	if typeName == nil {
		log.Fatalf("Nothing type name")
	}
	if err != nil {
		log.Fatalf("Missing the directory: %s", err)
	}
	if err := templateGenerator.Execute(&buf, Variable{
		TypeName:    *typeName,
		PackageName: scanPkg(dir).Pkg.Name(),
	}); err != nil {
		log.Fatalf("Cannot generate: %s", err)
	}
	if src, err := format.Source(buf.Bytes()); err != nil {
		log.Fatalf("Missing format")
	} else if err := ioutil.WriteFile(filepath.Join(dir, strings.ToLower(*typeName)+"_stream.go"), src, 0644); err != nil {
		log.Fatalf("Writing a file: %s", err)
	}

}

var templateGenerator = template.Must(template.New("").Parse(`
package {{.PackageName}}
type {{.TypeName}}Stream []{{.TypeName}}

func {{.TypeName}}StreamOf(arg ...{{.TypeName}}) {{.TypeName}}Stream {
	return arg
}
func {{.TypeName}}StreamFrom(arg []{{.TypeName}}) {{.TypeName}}Stream {
	return arg
}

func (self *{{.TypeName}}Stream) Add(arg {{.TypeName}}) *{{.TypeName}}Stream {
	return self.AddAll(arg)

}

func (self *{{.TypeName}}Stream) AddAll(arg ...{{.TypeName}}) *{{.TypeName}}Stream {
	*self = append(*self, arg...)
	return self
}

func (self *{{.TypeName}}Stream) AddSafe(arg *{{.TypeName}}) *{{.TypeName}}Stream {
    if arg != nil {
        self.Add(*arg)
    }
	return self

}
func (self *{{.TypeName}}Stream) AllMatch(fn func(arg {{.TypeName}}, index int) bool) bool {
	for i, v := range *self {
		if !fn(v, i) {
			return false
		}
	}
	return true
}

func (self *{{.TypeName}}Stream) AnyMatch(fn func(arg {{.TypeName}}, index int) bool) bool {
	for i, v := range *self {
		if fn(v, i) {
			return true
		}
	}
	return false
}
func (self *{{.TypeName}}Stream) Clone() *{{.TypeName}}Stream {
	temp := {{.TypeName}}StreamOf()
	temp = *self
	return &temp
}

func (self *{{.TypeName}}Stream) Copy() *{{.TypeName}}Stream {
	return self.Clone()
}

func (self *{{.TypeName}}Stream) Concat(arg []{{.TypeName}}) *{{.TypeName}}Stream {
	return self.AddAll(arg...)
}

func (self *{{.TypeName}}Stream) Delete(index int) *{{.TypeName}}Stream {
	if len(*self) > index+1 {
		*self = append((*self)[:index], (*self)[index+1:]...)
	} else {
		*self = append((*self)[:index])
	}
	return self
}

func (self *{{.TypeName}}Stream) DeleteRange(startIndex int, endIndex int) *{{.TypeName}}Stream {
	*self = append((*self)[:startIndex], (*self)[endIndex+1:]...)
	return self
}

func (self *{{.TypeName}}Stream) Filter(fn func(arg {{.TypeName}}, index int) bool) *{{.TypeName}}Stream {
	_array := []{{.TypeName}}{}
	for i, v := range *self {
		if fn(v, i) {
			_array = append(_array, v)
		}
	}
	*self = _array
	return self
}

func (self *{{.TypeName}}Stream) Find(fn func(arg {{.TypeName}}, index int) bool) *{{.TypeName}} {
	i := self.FindIndex(fn)
	if -1 != i {
		return &(*self)[i]
	}
	return nil
}

func (self *{{.TypeName}}Stream) FindIndex(fn func(arg {{.TypeName}}, index int) bool) int {
	for i, v := range *self {
		if fn(v, i) {
			return i
		}
	}
	return -1
}

func (self *{{.TypeName}}Stream) First() *{{.TypeName}} {
	return self.Get(0)
}

func (self *{{.TypeName}}Stream) ForEach(fn func(arg {{.TypeName}}, index int)) *{{.TypeName}}Stream {
	for i, v := range *self {
		fn(v, i)
	}
	return self
}
func (self *{{.TypeName}}Stream) ForEachRight(fn func(arg {{.TypeName}}, index int)) *{{.TypeName}}Stream {
	for i := len(*self) - 1; i >= 0; i-- {
		fn((*self)[i], i)
	}
	return self
}
func (self *{{.TypeName}}Stream) GroupBy(fn func(arg {{.TypeName}}, index int) string) map[string][]{{.TypeName}} {
    m := map[string][]{{.TypeName}}{}
    for i, v := range *self {
        key := fn(v, i)
        m[key] = append(m[key], v)
    }
    return m
}
func (self *{{.TypeName}}Stream) GroupByValues(fn func(arg {{.TypeName}}, index int) string) [][]{{.TypeName}} {
	tmp := [][]{{.TypeName}}{}
	m := self.GroupBy(fn)
	for _, v := range m {
		tmp = append(tmp, v)
	}
	return tmp
}
func (self *{{.TypeName}}Stream) IsEmpty() bool {
	if self.Len() == 0 {
		return true
	}
	return false
}
func (self *{{.TypeName}}Stream) IsPreset() bool {
	return !self.IsEmpty()
}
func (self *{{.TypeName}}Stream) Last() *{{.TypeName}} {
	return self.Get(self.Len() - 1)
}

func (self *{{.TypeName}}Stream) Len() int {
    if self == nil {
		return 0
	}
	return len(*self)
}

func (self *{{.TypeName}}Stream) Map(fn func(arg {{.TypeName}}, index int) {{.TypeName}}) *{{.TypeName}}Stream {
	_array := make([]{{.TypeName}}, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	*self = _array
	return self
}

func (self *{{.TypeName}}Stream) MapAny(fn func(arg {{.TypeName}}, index int) interface{}) []interface{} {
	_array := make([]interface{}, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) Map2Int(fn func(arg {{.TypeName}}, index int) int) []int {
	_array := make([]int, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) Map2Int32(fn func(arg {{.TypeName}}, index int) int32) []int32 {
	_array := make([]int32, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) Map2Int64(fn func(arg {{.TypeName}}, index int) int64) []int64 {
	_array := make([]int64, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) Map2Float32(fn func(arg {{.TypeName}}, index int) float32) []float32 {
	_array := make([]float32, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) Map2Float64(fn func(arg {{.TypeName}}, index int) float64) []float64 {
	_array := make([]float64, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) Map2Bool(fn func(arg {{.TypeName}}, index int) bool) []bool {
	_array := make([]bool, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) Map2Bytes(fn func(arg {{.TypeName}}, index int) []byte) [][]byte {
	_array := make([][]byte, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) Map2String(fn func(arg {{.TypeName}}, index int) string) []string {
	_array := make([]string, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *{{.TypeName}}Stream) NoneMatch(fn func(arg {{.TypeName}}, index int) bool) bool {
	return !self.AnyMatch(fn)
}

func (self *{{.TypeName}}Stream) Get(index int) *{{.TypeName}} {
	if self.Len() > index && index >= 0 {
		return &(*self)[index]
	}
	return nil
}
func (self *{{.TypeName}}Stream) ReduceInit(fn func(result, current {{.TypeName}}, index int) {{.TypeName}}, initialValue {{.TypeName}}) *{{.TypeName}}Stream {
	result := []{{.TypeName}}{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn(initialValue, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	*self = result
	return self
}
func (self *{{.TypeName}}Stream) Reduce(fn func(result, current {{.TypeName}}, index int) {{.TypeName}}) *{{.TypeName}}Stream {
	result := []{{.TypeName}}{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn({{.TypeName}}{}, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	*self = result
	return self
}
func (self *{{.TypeName}}Stream) ReduceInterface(fn func(result interface{}, current {{.TypeName}}, index int) interface{}) []interface{} {
	result := []interface{}{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn({{.TypeName}}{}, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *{{.TypeName}}Stream) ReduceString(fn func(result string, current {{.TypeName}}, index int) string) []string {
	result := []string{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn("", v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *{{.TypeName}}Stream) ReduceInt(fn func(result int, current {{.TypeName}}, index int) int) []int {
	result := []int{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn(0, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *{{.TypeName}}Stream) ReduceInt32(fn func(result int32, current {{.TypeName}}, index int) int32) []int32 {
	result := []int32{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn(0, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *{{.TypeName}}Stream) ReduceInt64(fn func(result int64, current {{.TypeName}}, index int) int64) []int64 {
	result := []int64{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn(0, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *{{.TypeName}}Stream) ReduceFloat32(fn func(result float32, current {{.TypeName}}, index int) float32) []float32 {
	result := []float32{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn(0.0, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *{{.TypeName}}Stream) ReduceFloat64(fn func(result float64, current {{.TypeName}}, index int) float64) []float64 {
	result := []float64{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn(0.0, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *{{.TypeName}}Stream) ReduceBool(fn func(result bool, current {{.TypeName}}, index int) bool) []bool {
	result := []bool{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn(false, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *{{.TypeName}}Stream) Reverse() *{{.TypeName}}Stream {
	_array := make([]{{.TypeName}}, 0, len(*self))
	for i := len(*self) - 1; i >= 0; i-- {
		_array = append(_array, (*self)[i])
	}
	*self = _array
	return self
}

func (self *{{.TypeName}}Stream) Replace(fn func(arg {{.TypeName}}, index int) {{.TypeName}}) *{{.TypeName}}Stream {
	for i, v := range *self {
		(*self)[i] = fn(v, i)
	}
	return self
}

func (self *{{.TypeName}}Stream) Set(index int, val {{.TypeName}}) {
	if len(*self) > index {
		(*self)[index] = val
	}
}

func (self *{{.TypeName}}Stream) Slice(startIndex int, n int) *{{.TypeName}}Stream {
    last := startIndex+n
    if len(*self)-1 < startIndex {
        *self = []{{.TypeName}}{}
    } else if len(*self) < last {
        *self = (*self)[startIndex:len(*self)]
    } else {
        *self = (*self)[startIndex:last]
    }
	return self
}

func (self *{{.TypeName}}Stream) ToList() []{{.TypeName}} {
	return self.Val()
}

func (self *{{.TypeName}}Stream) Val() []{{.TypeName}} {
    if self == nil {R
        return []{{.TypeName}}{}
    }
	return *self
}

`))