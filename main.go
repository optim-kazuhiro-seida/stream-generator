package main

import (
	"flag"
	"go/build"
	"go/types"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

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
	)
	flag.Parse()
	if typeName == nil || *typeName == "" {
		log.Fatalf("Nothing type name")
	} else {
		log.Println("Generating Struct name: ", *typeName)
	}
	if err != nil {
		log.Fatalf("Missing the directory: %s", err)
	}

	result := strings.Replace(template, "{{.TypeName}}", *typeName, -1)
	log.Println("Scan Type Name.")

	result = strings.Replace(result, "{{.PackageName}}", scanPkg(dir).Pkg.Name(), -1)
	log.Println("Scan Package Name.")

	if err := ioutil.WriteFile(filepath.Join(dir, "stream_"+strings.ToLower(*typeName)+".go"), []byte(result), 0644); err != nil {
		log.Fatalf("Writing a file: %s", err)
	} else {
		log.Println("Generated " + filepath.Join(dir, "stream_"+strings.ToLower(*typeName)+".go"))
	}
}

var template = `
package {{.PackageName}}

import (
	"math"
	"reflect"
	"sort"
)
type {{.TypeName}}Stream []{{.TypeName}}
func {{.TypeName}}StreamOf(arg ...{{.TypeName}}) {{.TypeName}}Stream {
	return arg
}
func {{.TypeName}}StreamFrom(arg []{{.TypeName}}) {{.TypeName}}Stream {
	return arg
}
func Create{{.TypeName}}Stream(arg ...{{.TypeName}}) *{{.TypeName}}Stream {
    tmp := {{.TypeName}}StreamOf(arg...)
    return &tmp
}
func Generate{{.TypeName}}Stream(arg []{{.TypeName}}) *{{.TypeName}}Stream {
    tmp := {{.TypeName}}StreamFrom(arg)
    return &tmp
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
func (self *{{.TypeName}}Stream) AllMatch(fn func({{.TypeName}}, int) bool) bool {
	for i, v := range *self {
		if !fn(v, i) {
			return false
		}
	}
	return true
}
func (self *{{.TypeName}}Stream) AnyMatch(fn func({{.TypeName}}, int) bool) bool {
	for i, v := range *self {
		if fn(v, i) {
			return true
		}
	}
	return false
}
func (self *{{.TypeName}}Stream) Clone() *{{.TypeName}}Stream {
	temp := make([]{{.TypeName}}, self.Len())
	copy(temp, *self)
	return (*{{.TypeName}}Stream)(&temp)
}
func (self *{{.TypeName}}Stream) Copy() *{{.TypeName}}Stream {
	return self.Clone()
}
func (self *{{.TypeName}}Stream) Concat(arg []{{.TypeName}}) *{{.TypeName}}Stream {
	return self.AddAll(arg...)
}
func (self *{{.TypeName}}Stream) Contains(arg {{.TypeName}}) bool {
	return self.FindIndex(func(_arg {{.TypeName}}, index int) bool { return reflect.DeepEqual(_arg, arg) }) != -1
}
func (self *{{.TypeName}}Stream) Clean() *{{.TypeName}}Stream {
    return Create{{.TypeName}}Stream()
}
func (self *{{.TypeName}}Stream) Delete(index int) *{{.TypeName}}Stream {
	return self.DeleteRange(index, index)
}
func (self *{{.TypeName}}Stream) DeleteRange(startIndex, endIndex int) *{{.TypeName}}Stream {
	*self = append((*self)[:startIndex], (*self)[endIndex+1:]...)
	return self
}
func (self *{{.TypeName}}Stream) Distinct() *{{.TypeName}}Stream {
	stack := {{.TypeName}}StreamOf()
	return self.Filter(func(arg {{.TypeName}}, _ int) bool {
		if !stack.Contains(arg) {
			stack.Add(arg)
			return true
		}
		return false
	})
}
func (self *{{.TypeName}}Stream) Equals(arr []{{.TypeName}}) bool {
	if (*self == nil) != (arr == nil) || len(*self) != len(arr) {
		return false
	}
	for i := range *self {
		if !reflect.DeepEqual((*self)[i], arr[i]) {
			return false
		}
	}
	return true
}
func (self *{{.TypeName}}Stream) Filter(fn func({{.TypeName}}, int) bool) *{{.TypeName}}Stream {
	_array := {{.TypeName}}StreamOf()
	self.ForEach(func(v {{.TypeName}}, i int) {
		if fn(v, i) {
			_array.Add(v)
		}
	})
	*self = _array
	return self
}
func (self *{{.TypeName}}Stream) Find(fn func({{.TypeName}}, int) bool) *{{.TypeName}} {
	i := self.FindIndex(fn)
	if -1 != i {
		return &(*self)[i]
	}
	return nil
}
func (self *{{.TypeName}}Stream) FindIndex(fn func({{.TypeName}}, int) bool) int {
	for i, v := range self.Val() {
		if fn(v, i) {
			return i
		}
	}
	return -1
}
func (self *{{.TypeName}}Stream) First() *{{.TypeName}} {
	return self.Get(0)
}
func (self *{{.TypeName}}Stream) ForEach(fn func({{.TypeName}}, int)) *{{.TypeName}}Stream {
	for i, v := range self.Val() {
		fn(v, i)
	}
	return self
}
func (self *{{.TypeName}}Stream) ForEachRight(fn func({{.TypeName}}, int)) *{{.TypeName}}Stream {
	for i := self.Len() - 1; i >= 0; i-- {
		fn(*self.Get(i), i)
	}
	return self
}
func (self *{{.TypeName}}Stream) GroupBy(fn func({{.TypeName}}, int) string) map[string][]{{.TypeName}} {
    m := map[string][]{{.TypeName}}{}
    for i, v := range self.Val() {
        key := fn(v, i)
        m[key] = append(m[key], v)
    }
    return m
}
func (self *{{.TypeName}}Stream) GroupByValues(fn func({{.TypeName}}, int) string) [][]{{.TypeName}} {
	tmp := [][]{{.TypeName}}{}
	m := self.GroupBy(fn)
	for _, v := range m {
		tmp = append(tmp, v)
	}
	return tmp
}
func (self *{{.TypeName}}Stream) IndexOf(arg {{.TypeName}}) int {
	for index, _arg := range *self {
		if reflect.DeepEqual(_arg, arg) {
			return index
		}
	}
	return -1
}
func (self *{{.TypeName}}Stream) IsEmpty() bool {
	return self.Len() == 0
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
func (self *{{.TypeName}}Stream) Limit(limit int) *{{.TypeName}}Stream {
	self.Slice(0, limit)
	return self
}
func (self *{{.TypeName}}Stream) Map(fn func({{.TypeName}}, int) {{.TypeName}}) *{{.TypeName}}Stream {
	return self.ForEach(func(v {{.TypeName}}, i int) { self.Set(i, fn(v, i)) })
}
func (self *{{.TypeName}}Stream) MapAny(fn func({{.TypeName}}, int) interface{}) []interface{} {
	_array := make([]interface{}, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Map2Int(fn func({{.TypeName}}, int) int) []int {
	_array := make([]int, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Map2Int32(fn func({{.TypeName}}, int) int32) []int32 {
	_array := make([]int32, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Map2Int64(fn func({{.TypeName}}, int) int64) []int64 {
	_array := make([]int64, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Map2Float32(fn func({{.TypeName}}, int) float32) []float32 {
	_array := make([]float32, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Map2Float64(fn func({{.TypeName}}, int) float64) []float64 {
	_array := make([]float64, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Map2Bool(fn func({{.TypeName}}, int) bool) []bool {
	_array := make([]bool, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Map2Bytes(fn func({{.TypeName}}, int) []byte) [][]byte {
	_array := make([][]byte, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Map2String(fn func({{.TypeName}}, int) string) []string {
	_array := make([]string, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}
func (self *{{.TypeName}}Stream) Max(fn func({{.TypeName}}, int) float64) *{{.TypeName}} {
	f := self.Get(0)
	if f == nil {
		return nil
	}
	m := fn(*f, 0)
	index := 0
	for i := 1; i < self.Len(); i++ {
		v := fn(*self.Get(i), i)
		m = math.Max(m, v)
		if m == v {
			index = i
		}
	}
	return self.Get(index)
}
func (self *{{.TypeName}}Stream) Min(fn func({{.TypeName}}, int) float64) *{{.TypeName}} {
	f := self.Get(0)
	if f == nil {
		return nil
	}
	m := fn(*f, 0)
	index := 0
	for i := 1; i < self.Len(); i++ {
		v := fn(*self.Get(i), i)
		m = math.Min(m, v)
		if m == v {
			index = i
		}
	}
	return self.Get(index)
}
func (self *{{.TypeName}}Stream) NoneMatch(fn func({{.TypeName}}, int) bool) bool {
	return !self.AnyMatch(fn)
}
func (self *{{.TypeName}}Stream) Get(index int) *{{.TypeName}} {
	if self.Len() > index && index >= 0 {
		tmp := (*self)[index]
        return &tmp
	}
	return nil
}
func (self *{{.TypeName}}Stream) Peek(fn func(*{{.TypeName}}, int)) *{{.TypeName}}Stream {
    for i, v := range *self {
        fn(&v, i)
        self.Set(i, v)
    }
    return self
}
func (self *{{.TypeName}}Stream) Reduce(fn func({{.TypeName}}, {{.TypeName}}, int) {{.TypeName}}) *{{.TypeName}}Stream {
	return self.ReduceInit(fn, {{.TypeName}}{})
}
func (self *{{.TypeName}}Stream) ReduceInit(fn func({{.TypeName}}, {{.TypeName}}, int) {{.TypeName}}, initialValue {{.TypeName}}) *{{.TypeName}}Stream {
	result :={{.TypeName}}StreamOf()
	self.ForEach(func(v {{.TypeName}}, i int) {
		if i == 0 {
			result.Add(fn(initialValue, v, i))
		} else {
			result.Add(fn(result[i-1], v, i))
		}
	})
	*self = result
	return self
}
func (self *{{.TypeName}}Stream) ReduceInterface(fn func(interface{}, {{.TypeName}}, int) interface{}) []interface{} {
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
func (self *{{.TypeName}}Stream) ReduceString(fn func(string, {{.TypeName}}, int) string) []string {
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
func (self *{{.TypeName}}Stream) ReduceInt(fn func(int, {{.TypeName}}, int) int) []int {
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
func (self *{{.TypeName}}Stream) ReduceInt32(fn func(int32, {{.TypeName}}, int) int32) []int32 {
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
func (self *{{.TypeName}}Stream) ReduceInt64(fn func(int64, {{.TypeName}}, int) int64) []int64 {
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
func (self *{{.TypeName}}Stream) ReduceFloat32(fn func(float32, {{.TypeName}}, int) float32) []float32 {
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
func (self *{{.TypeName}}Stream) ReduceFloat64(fn func(float64, {{.TypeName}}, int) float64) []float64 {
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
func (self *{{.TypeName}}Stream) ReduceBool(fn func(bool, {{.TypeName}}, int) bool) []bool {
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
	for i, j := 0, self.Len()-1; i < j; i, j = i+1, j-1 {
		(*self)[i], (*self)[j] = (*self)[j], (*self)[i]
	}
	return self
}
func (self *{{.TypeName}}Stream) Replace(fn func({{.TypeName}}, int) {{.TypeName}}) *{{.TypeName}}Stream {
	return self.Map(fn)
}
func (self *{{.TypeName}}Stream) Set(index int, val {{.TypeName}}) *{{.TypeName}}Stream {
    if len(*self) > index {
        (*self)[index] = val
    }
    return self
}
func (self *{{.TypeName}}Stream) Skip(skip int) *{{.TypeName}}Stream {
	self.Slice(skip, self.Len()-skip)
	return self
}
func (self *{{.TypeName}}Stream) SkippingEach(fn func({{.TypeName}}, int) int) *{{.TypeName}}Stream {
	for i := 0; i < self.Len(); i++ {
		skip := fn(*self.Get(i), i)
		i += skip
	}
	return self
}
func (self *{{.TypeName}}Stream) Slice(startIndex, n int) *{{.TypeName}}Stream {
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
func (self *{{.TypeName}}Stream) Sort(fn func(i, j int) bool) *{{.TypeName}}Stream {
	sort.Slice(*self, fn)
	return self
}
func (self *{{.TypeName}}Stream) SortStable(fn func(i, j int) bool) *{{.TypeName}}Stream {
	sort.SliceStable(*self, fn)
	return self
}
func (self *{{.TypeName}}Stream) ToList() []{{.TypeName}} {
	return self.Val()
}
func (self *{{.TypeName}}Stream) Unique() *{{.TypeName}}Stream {
	return self.Distinct()
}
func (self *{{.TypeName}}Stream) Val() []{{.TypeName}} {
	if self == nil {
		return []{{.TypeName}}{}
	}
	return *self.Copy()
}
func (self *{{.TypeName}}Stream) While(fn func({{.TypeName}}, int) bool) *{{.TypeName}}Stream {
    for i, v := range self.Val() {
        if !fn(v, i) {
            break
        }
    }
    return self
}
`
