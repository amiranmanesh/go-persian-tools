package digit

import "testing"

// sink keeps the compiler from optimizing the benchmarked calls away.
var (
	sink    string
	sinkInt int64
)

const (
	benchMixed  = "شماره حساب ۱۲۳۴۵۶۷۸۹۰ و مبلغ 1234567 ریال برای علي رضا"
	benchAmount = "1234567890"
)

func BenchmarkToPersianDigits(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = ToPersianDigits(benchMixed)
	}
}

func BenchmarkToEnglishDigits(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = ToEnglishDigits(benchMixed)
	}
}

func BenchmarkOnlyNumbers(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = OnlyNumbers(benchMixed)
	}
}

func BenchmarkCurrency(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = Currency(benchAmount)
	}
}

func BenchmarkAddCommas(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = AddCommas(14555478854)
	}
}

func BenchmarkRemoveCommas(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = RemoveCommas("14,555,478,854")
	}
}

func BenchmarkToWords(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sink = ToWords(9223372036854775807)
	}
}
