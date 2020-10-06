package sample

import (
	"reflect"
	"strconv"
	"testing"
)

func TestStream(t *testing.T) {
	stream := SampleStreamOf()
	stream.AddAll(
		Sample{
			Str: "1",
			Int: 1,
		},
		Sample{
			Str: "2",
			Int: 2,
		},
		Sample{
			Str: "3",
			Int: 3,
		},
		Sample{
			Str: "4",
			Int: 4,
		},
		Sample{
			Str: "5",
			Int: 5,
		},
	)
	cloned1 := stream.Clone()
	cloned2 := stream.Clone()
	cloned3 := stream.Clone()

	stream.ForEach(func(arg Sample, index int) {
		if arg.Int != index+1 {
			t.Fatal("Unexpect Value: AddAll ", index)
		}
	})
	if stream.Add(Sample{Str: "6", Int: 6}); stream.Len() != 6 {
		t.Fatal("Unexpect Value stream length.", stream)
	}
	if !stream.AllMatch(func(_ Sample, _ int) bool { return true }) ||
		stream.AllMatch(func(_ Sample, _ int) bool { return false }) {
		t.Fatal("Unexpect Value stream AllMatch.", stream)
	}
	if stream.Concat([]Sample{{Str: "7", Int: 7}, {Str: "8", Int: 8}}); stream.Get(6).Str != "7" || stream.Get(6).Int != 7 {
		t.Fatal("Unexpect Value stream Concat.", stream)
	}
	if cloned1.Delete(0); cloned1.Len() == 6 || stream.Len() == cloned1.Len() || cloned1.Get(0).Str != stream.Get(1).Str {
		t.Fatal("Unexpect Value stream Delete.", cloned1)
	}
	if cloned1.DeleteRange(0, 2); cloned1.Get(0).Int != 5 {
		t.Fatal("Unexpect Value stream DeleteRange.", cloned1)
	}
	if cloned1 = stream.Copy(); !reflect.DeepEqual(cloned1.Last(), stream.Last()) {
		t.Fatal("Unexpect Value stream Copy.", cloned1)
	}
	if cloned1.Filter(func(arg Sample, _ int) bool { return arg.Int%2 == 0 }); cloned1.Len() != 4 || cloned1.Get(0).Str != "2" {
		t.Fatal("Unexpect Value stream Filter.", cloned1)
	}
	if val := stream.Find(func(arg Sample, _ int) bool { return arg.Str == "5" }); cloned1.Find(func(arg Sample, _ int) bool { return arg.Str == "5" }) != nil || val.Str != "5" {
		t.Fatal("Unexpect Value stream Find.", cloned1)
	}
	if index := stream.FindIndex(func(arg Sample, _ int) bool { return arg.Int == 8 }); index != stream.Len()-1 {
		t.Fatal("Unexpect Value stream FindIndex.", stream)
	}
	if !reflect.DeepEqual(stream.First(), stream.Get(0)) {
		t.Fatal("Unexpect Value stream First.", stream)
	}
	if m := stream.GroupBy(func(arg Sample, _ int) string { return strconv.Itoa(arg.Int % 4) }); len(m["1"]) != 2 || m["1"][0].Int != 1 || m["2"][0].Int != 2 || m["3"][0].Int != 3 {
		t.Fatal("Unexpect Value stream GroupBy.", m)
	}
	if v := stream.GroupByValues(func(arg Sample, _ int) string { return strconv.Itoa(arg.Int % 4) }); len(v) != 4 {
		t.Fatal("Unexpect Value stream GroupByValues.", v)
	}
	if tmp := SampleStreamOf(); stream.IsEmpty() || !tmp.IsEmpty() {
		t.Fatal("Unexpect Value stream IsEmpty.", tmp)
	}
	if tmp := SampleStreamOf(); !stream.IsPreset() || tmp.IsPreset() {
		t.Fatal("Unexpect Value stream IsPreset.", tmp)
	}
	if cloned2.Map(func(arg Sample, index int) Sample { return Sample{Str: "test", Int: index} }); cloned2.First().Str != "test" && cloned2.First().Int != 0 && cloned2.Last().Str != "test" && cloned2.Last().Int != 4 {
		t.Fatal("Unexpect Value stream Map.", cloned2)
	}
	if cloned2.Add(Sample{Str: "last", Int: 999}); cloned2.Last().Str != "last" || cloned2.Last().Int != 999 {
		t.Fatal("Unexpect Value stream make slice.", cloned2)
	}
	if !cloned2.NoneMatch(func(_ Sample, _ int) bool { return false }) || cloned1.NoneMatch(func(_ Sample, _ int) bool { return true }) {
		t.Fatal("Unexpect Value stream NoneMatch.", cloned2)
	}
	if stream.Get(8888) != nil || stream.Get(0) == nil || stream.Get(stream.Len()-1) == nil || stream.Get(-1) != nil {
		t.Fatal("Unexpect Value stream Get.", stream)
	}
	if cloned3.ReduceInit(func(result, current Sample, index int) Sample { current.Int += result.Int; return current }, Sample{Int: 0}); cloned3.Last().Int != 15 {
		t.Fatal("Unexpect Value stream ReduceInit.", cloned3)
	}
	if cloned3.Replace(func(arg Sample, index int) Sample { return Sample{Str: "test", Int: 5} }); cloned2.First().Str != "test" && cloned2.First().Int != 5 && cloned2.Last().Str != "test" && cloned2.Last().Int != 5 {
		t.Fatal("Unexpect Value stream Map.", cloned2)
	}
	if cloned3.Reduce(func(result, current Sample, index int) Sample { current.Int += result.Int; return current }); cloned3.Last().Int != 25 {
		t.Fatal("Unexpect Value stream Reduce.", cloned3)
	}
	if cloned2.Reverse(); cloned2.First().Int != 999 || cloned2.Get(1).Int != 4 || cloned2.Get(2).Int != 3 || cloned2.Get(3).Int != 2 || cloned2.Get(4).Int != 1 || cloned2.Last().Int != 0 {
		t.Fatal("Unexpect Value stream Reverse.", cloned3)
	}
	if cloned2 = stream.Copy(); !stream.Equals(cloned2.Val()) || stream.Equals(cloned3.Val()) {
		t.Fatal("Unexpect Value stream Copy.", stream)
	}
	if index := stream.IndexOf(*stream.Get(3)); index != 3 {
		t.Fatal("Unexpect Value stream IndexOf.", stream)
	}
	if index := stream.IndexOf(*stream.Get(2)); index != 2 {
		t.Fatal("Unexpect Value stream IndexOf.", stream)
	}
	if index := stream.IndexOf(*stream.Get(5)); index != 5 {
		t.Fatal("Unexpect Value stream IndexOf.", stream)
	}
	if v := stream.Max(func(arg Sample, index int) float64 { return float64(arg.Int) }); v.Int != 8 {
		t.Fatal("Unexpect Value stream Max.", stream)
	}
	if v := stream.Min(func(arg Sample, index int) float64 { return float64(arg.Int) }); v.Int != 1 {
		t.Fatal("Unexpect Value stream Max.", stream)
	}
	if v := cloned2.Peek(func(arg *Sample, index int) { arg.Int = -1 }); v.Get(0).Int != -1 {
		t.Fatal("Unexpect Value stream Peek.", cloned2)
	}
	if v := cloned2.ForEach(func(arg Sample, index int) { arg.Int = 2 }); v.Get(0).Int == 2 {
		t.Fatal("Unexpect Value stream ForEach.", cloned2)
	}
	count := 0
	if stream.While(func(arg Sample, index int) bool { arg.Int = -1; count++; return index != 3 }); count != 4 && stream.Get(0).Int == -1 {
		t.Fatal("Unexpect Value stream While.", stream)
	}
	if v := stream.Get(0); v.Str != "1" {
		t.Fatal("Unexpect Value stream Get.", stream)
	} else {
		v.Str = ""
		(*v).Str = ""
		if stream.Get(0).Str != "1" {
			t.Fatal("Unexpect Value stream Get.", stream)
		}
	}
	if stream.Set(stream.Len()-1, Sample{Str: "last", Int: 0}); stream.Last().Str != "last" {
		t.Fatal("Unexpect Value stream Set.", stream)
	}
	if tmp := stream.Copy().Skip(2).Limit(2); tmp.Get(0).Int != 3 || tmp.Len() != 2 {
		t.Fatal("Unexpect Value stream Limit and Skip.", stream)
	}
}
