package text

import "testing"

// sink keeps the compiler from optimizing the benchmarked calls away.
var sink string

const benchMixed = "شماره حساب ۱۲۳۴۵۶۷۸۹۰ و مبلغ 1234567 ریال برای علي رضا"

func BenchmarkNormalize(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = Normalize(benchMixed)
	}
}

func BenchmarkFixArabic(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = FixArabic(benchMixed)
	}
}

func BenchmarkSwitchToPersianKey(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = SwitchToPersianKey("sghl o,fd")
	}
}

func BenchmarkFinglish(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = Finglish("سلام خوبی")
	}
}

func BenchmarkOnlyPersianAlpha(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = OnlyPersianAlpha(benchMixed)
	}
}
