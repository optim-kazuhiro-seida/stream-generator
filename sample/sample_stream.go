package sample

type SampleStream []Sample

func SampleStreamOf(arg ...Sample) SampleStream {
	return arg
}
func SampleStreamFrom(arg []Sample) SampleStream {
	return arg
}

func (self *SampleStream) Add(arg Sample) *SampleStream {
	return self.AddAll(arg)

}

func (self *SampleStream) AddAll(arg ...Sample) *SampleStream {
	*self = append(*self, arg...)
	return self
}

func (self *SampleStream) AddSafe(arg *Sample) *SampleStream {
	if arg != nil {
		self.Add(*arg)
	}
	return self

}
func (self *SampleStream) AllMatch(fn func(arg Sample, index int) bool) bool {
	for i, v := range *self {
		if !fn(v, i) {
			return false
		}
	}
	return true
}

func (self *SampleStream) AnyMatch(fn func(arg Sample, index int) bool) bool {
	for i, v := range *self {
		if fn(v, i) {
			return true
		}
	}
	return false
}
func (self *SampleStream) Clone() *SampleStream {
	temp := make([]Sample, self.Len())
	copy(temp, *self)
	return (*SampleStream)(&temp)
}

func (self *SampleStream) Copy() *SampleStream {
	return self.Clone()
}

func (self *SampleStream) Concat(arg []Sample) *SampleStream {
	return self.AddAll(arg...)
}

func (self *SampleStream) Delete(index int) *SampleStream {
	return self.DeleteRange(index, index)
}

func (self *SampleStream) DeleteRange(startIndex int, endIndex int) *SampleStream {
	if startIndex < 0 || endIndex+1 > self.Len() {
		return self
	}
	*self = append((*self)[:startIndex], (*self)[endIndex+1:]...)
	return self
}

func (self *SampleStream) Filter(fn func(arg Sample, index int) bool) *SampleStream {
	_array := SampleStreamOf()
	self.ForEach(func(v Sample, i int) {
		if fn(v, i) {
			_array.Add(v)
		}
	})
	*self = _array
	return self
}

func (self *SampleStream) Find(fn func(arg Sample, index int) bool) *Sample {
	i := self.FindIndex(fn)
	if -1 != i {
		return &(*self)[i]
	}
	return nil
}

func (self *SampleStream) FindIndex(fn func(arg Sample, index int) bool) int {
	for i, v := range *self {
		if fn(v, i) {
			return i
		}
	}
	return -1
}

func (self *SampleStream) First() *Sample {
	return self.Get(0)
}

func (self *SampleStream) ForEach(fn func(arg Sample, index int)) *SampleStream {
	for i, v := range *self {
		fn(v, i)
	}
	return self
}
func (self *SampleStream) ForEachRight(fn func(arg Sample, index int)) *SampleStream {
	for i := len(*self) - 1; i >= 0; i-- {
		fn((*self)[i], i)
	}
	return self
}
func (self *SampleStream) GroupBy(fn func(arg Sample, index int) string) map[string][]Sample {
	m := map[string][]Sample{}
	for i, v := range *self {
		key := fn(v, i)
		m[key] = append(m[key], v)
	}
	return m
}
func (self *SampleStream) GroupByValues(fn func(arg Sample, index int) string) [][]Sample {
	tmp := [][]Sample{}
	m := self.GroupBy(fn)
	for _, v := range m {
		tmp = append(tmp, v)
	}
	return tmp
}
func (self *SampleStream) IsEmpty() bool {
	return self.Len() == 0
}
func (self *SampleStream) IsPreset() bool {
	return !self.IsEmpty()
}
func (self *SampleStream) Last() *Sample {
	return self.Get(self.Len() - 1)
}

func (self *SampleStream) Len() int {
	return len(*self)
}

func (self *SampleStream) Map(fn func(arg Sample, index int) Sample) *SampleStream {
	for i, v := range *self {
		(*self).Set(i, fn(v, i))
	}
	return self
}

func (self *SampleStream) MapAny(fn func(arg Sample, index int) interface{}) []interface{} {
	_array := make([]interface{}, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) Map2Int(fn func(arg Sample, index int) int) []int {
	_array := make([]int, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) Map2Int32(fn func(arg Sample, index int) int32) []int32 {
	_array := make([]int32, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) Map2Int64(fn func(arg Sample, index int) int64) []int64 {
	_array := make([]int64, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) Map2Float32(fn func(arg Sample, index int) float32) []float32 {
	_array := make([]float32, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) Map2Float64(fn func(arg Sample, index int) float64) []float64 {
	_array := make([]float64, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) Map2Bool(fn func(arg Sample, index int) bool) []bool {
	_array := make([]bool, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) Map2Bytes(fn func(arg Sample, index int) []byte) [][]byte {
	_array := make([][]byte, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) Map2String(fn func(arg Sample, index int) string) []string {
	_array := make([]string, 0, len(*self))
	for i, v := range *self {
		_array = append(_array, fn(v, i))
	}
	return _array
}

func (self *SampleStream) NoneMatch(fn func(arg Sample, index int) bool) bool {
	return !self.AnyMatch(fn)
}

func (self *SampleStream) Get(index int) *Sample {
	if self.Len() > index && index >= 0 {
		return &(*self)[index]
	}
	return nil
}
func (self *SampleStream) Reduce(fn func(result, current Sample, index int) Sample) *SampleStream {
	return self.ReduceInit(fn, Sample{})
}
func (self *SampleStream) ReduceInit(fn func(result, current Sample, index int) Sample, initialValue Sample) *SampleStream {
	result := []Sample{}
	self.ForEach(func(v Sample, i int) {
		if i == 0 {
			result = append(result, fn(initialValue, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	})
	*self = result
	return self
}

func (self *SampleStream) ReduceInterface(fn func(result interface{}, current Sample, index int) interface{}) []interface{} {
	result := []interface{}{}
	for i, v := range *self {
		if i == 0 {
			result = append(result, fn(Sample{}, v, i))
		} else {
			result = append(result, fn(result[i-1], v, i))
		}
	}
	return result
}
func (self *SampleStream) ReduceString(fn func(result string, current Sample, index int) string) []string {
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
func (self *SampleStream) ReduceInt(fn func(result int, current Sample, index int) int) []int {
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
func (self *SampleStream) ReduceInt32(fn func(result int32, current Sample, index int) int32) []int32 {
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
func (self *SampleStream) ReduceInt64(fn func(result int64, current Sample, index int) int64) []int64 {
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
func (self *SampleStream) ReduceFloat32(fn func(result float32, current Sample, index int) float32) []float32 {
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
func (self *SampleStream) ReduceFloat64(fn func(result float64, current Sample, index int) float64) []float64 {
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
func (self *SampleStream) ReduceBool(fn func(result bool, current Sample, index int) bool) []bool {
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
func (self *SampleStream) Reverse() *SampleStream {
	for i, j := 0, len(*self)-1; i < j; i, j = i+1, j-1 {
		(*self)[i], (*self)[j] = (*self)[j], (*self)[i]
	}
	return self
}

func (self *SampleStream) Replace(fn func(arg Sample, index int) Sample) *SampleStream {
	return self.Map(fn)
}

func (self *SampleStream) Set(index int, val Sample) {
	if len(*self) > index && index >= 0 {
		(*self)[index] = val
	}
}

func (self *SampleStream) Slice(startIndex int, n int) *SampleStream {
	last := startIndex + n
	if len(*self)-1 < startIndex {
		*self = []Sample{}
	} else if len(*self) < last {
		*self = (*self)[startIndex:len(*self)]
	} else {
		*self = (*self)[startIndex:last]
	}
	return self
}

func (self *SampleStream) ToList() []Sample {
	return self.Val()
}

func (self *SampleStream) Val() []Sample {
	if self == nil {
		return []Sample{}
	}
	return *self
}