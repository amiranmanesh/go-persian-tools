package text

import "testing"

// The reference word lists behind TestFinglishAccuracy.
//
// Finglish has no single standard — u/oo, e/eh and a/aa are all in daily use —
// so each entry lists every spelling a Persian speaker would accept.
//
// The list is split in two on purpose. finglishReference is what the rules were
// developed against; finglishHeldOut was written afterwards and never consulted
// while tuning them, so it measures whether the rules generalize instead of
// merely fitting the words they were built from.
type finglishCase struct {
	word string
	ok   []string
}

// finglishReference is the tuning list.
func finglishReference() []finglishCase {
	return []finglishCase{
		{"سلام", []string{"salam"}},
		{"خوب", []string{"khub", "khoob"}},
		{"بد", []string{"bad"}},
		{"مرد", []string{"mard"}},
		{"زن", []string{"zan"}},
		{"آب", []string{"ab", "aab"}},
		{"نان", []string{"nan", "naan"}},
		{"کتاب", []string{"ketab", "ketaab"}},
		{"تهران", []string{"tehran", "tehraan"}},
		{"ایران", []string{"iran", "iraan"}},
		{"دنیا", []string{"donya", "donia", "donyaa"}},
		{"پدر", []string{"pedar"}},
		{"مادر", []string{"madar", "maadar"}},
		{"برادر", []string{"baradar", "baraadar"}},
		{"خواهر", []string{"khahar", "khaahar"}},
		{"پسر", []string{"pesar"}},
		{"دختر", []string{"dokhtar"}},
		{"دوست", []string{"dust", "doost"}},
		{"خانه", []string{"khane", "khaneh", "khaaneh"}},
		{"مدرسه", []string{"madrese", "madreseh"}},
		{"دانشگاه", []string{"daneshgah", "daaneshgaah", "daneshgaah"}},
		{"شهر", []string{"shahr"}},
		{"خیابان", []string{"khiaban", "khiabaan", "khiyaban"}},
		{"کوه", []string{"kuh", "kooh"}},
		{"دریا", []string{"darya", "daryaa", "daria"}},
		{"زمین", []string{"zamin"}},
		{"آسمان", []string{"aseman", "asemaan", "aasemaan"}},
		{"روز", []string{"ruz", "rooz"}},
		{"شب", []string{"shab"}},
		{"سال", []string{"sal", "saal"}},
		{"ماه", []string{"mah", "maah"}},
		{"امروز", []string{"emruz", "emrooz"}},
		{"فردا", []string{"farda", "fardaa"}},
		{"صبح", []string{"sobh"}},
		{"غذا", []string{"ghaza", "ghazaa"}},
		{"چای", []string{"chay", "chaay", "chai"}},
		{"شیر", []string{"shir"}},
		{"گل", []string{"gol"}},
		{"درخت", []string{"derakht"}},
		{"ماشین", []string{"mashin", "maashin"}},
		{"در", []string{"dar"}},
		{"کار", []string{"kar", "kaar"}},
		{"دست", []string{"dast"}},
		{"سر", []string{"sar"}},
		{"چشم", []string{"cheshm"}},
		{"بزرگ", []string{"bozorg"}},
		{"گرم", []string{"garm"}},
		{"سرد", []string{"sard"}},
		{"زیاد", []string{"ziad", "ziaad", "ziyad"}},
		{"قشنگ", []string{"ghashang"}},
		{"سخت", []string{"sakht"}},
		{"رفتن", []string{"raftan"}},
		{"آمدن", []string{"amadan", "aamadan"}},
		{"خوردن", []string{"khordan"}},
		{"دیدن", []string{"didan"}},
		{"گفتن", []string{"goftan"}},
		{"امیر", []string{"amir"}},
		{"علی", []string{"ali"}},
		{"رضا", []string{"reza", "rezaa"}},
		{"مریم", []string{"maryam"}},
		{"شیراز", []string{"shiraz", "shiraaz"}},
		{"اصفهان", []string{"esfahan", "esfahaan", "isfahan"}},
		{"مشهد", []string{"mashhad"}},
		{"کرمان", []string{"kerman", "kermaan"}},
		{"زنجان", []string{"zanjan", "zanjaan"}},
		{"یزد", []string{"yazd"}},
	}
}

// finglishHeldOut is the validation list.
func finglishHeldOut() []finglishCase {
	return []finglishCase{
		{"باران", []string{"baran", "baraan"}},
		{"برف", []string{"barf"}},
		{"آتش", []string{"atash", "aatash"}},
		{"سنگ", []string{"sang"}},
		{"چوب", []string{"chub", "choob"}},
		{"نامه", []string{"name", "nameh", "naameh"}},
		{"پنجره", []string{"panjere", "panjareh", "panjereh"}},
		{"دیوار", []string{"divar", "divaar"}},
		{"میز", []string{"miz"}},
		{"صندلی", []string{"sandali"}},
		{"کلاس", []string{"kelas", "kelaas", "klas"}},
		{"استاد", []string{"ostad", "ostaad", "ustad"}},
		{"کتابخانه", []string{"ketabkhane", "ketabkhaneh", "ketaabkhaaneh"}},
		{"بیمارستان", []string{"bimarestan", "bimaarestaan", "bimarestaan"}},
		{"فروشگاه", []string{"forushgah", "forooshgah", "forushgaah"}},
		{"تلفن", []string{"telefon", "telfon"}},
		{"دوچرخه", []string{"docharkhe", "docharkheh", "doocharkheh"}},
		{"هواپیما", []string{"havapeyma", "havaapeymaa", "havapayma"}},
		{"قطار", []string{"ghatar", "ghataar"}},
		{"بازار", []string{"bazar", "baazaar", "bazaar"}},
		{"نمک", []string{"namak"}},
		{"شکر", []string{"shekar", "shakar"}},
		{"برنج", []string{"berenj", "beranj"}},
		{"سیب", []string{"sib"}},
		{"پرتقال", []string{"porteghal", "porteghaal"}},
		{"هندوانه", []string{"hendevane", "hendevaneh", "hendvaneh"}},
		{"گوشت", []string{"gusht", "goosht"}},
		{"نمکدان", []string{"namakdan", "namakdaan"}},
		{"آشپزخانه", []string{"ashpazkhane", "ashpazkhaneh", "aashpazkhaaneh"}},
		{"خورشید", []string{"khorshid"}},
		{"ستاره", []string{"setare", "setareh", "sotareh"}},
		{"ابر", []string{"abr"}},
		{"باد", []string{"bad", "baad"}},
		{"رنگ", []string{"rang"}},
		{"سبز", []string{"sabz"}},
		{"سفید", []string{"sefid", "safid"}},
		{"سیاه", []string{"siah", "siyah", "siaah"}},
		{"زرد", []string{"zard"}},
		{"بلند", []string{"boland", "baland"}},
		{"کوتاه", []string{"kutah", "kootah", "kutaah"}},
		{"تند", []string{"tond", "tand"}},
		{"آرام", []string{"aram", "aaraam", "araam"}},
		{"خوشحال", []string{"khoshhal", "khoshhaal"}},
		{"ناراحت", []string{"narahat", "naaraahat", "narahaat"}},
		{"دانستن", []string{"danestan", "daanestan"}},
		{"نوشتن", []string{"neveshtan", "nveshtan"}},
		{"خواندن", []string{"khandan", "khaandan"}},
		{"شنیدن", []string{"shenidan", "shanidan"}},
		{"نشستن", []string{"neshastan", "nashastan"}},
		{"فروختن", []string{"forukhtan", "forookhtan"}},
		{"تبریز", []string{"tabriz"}},
		{"اهواز", []string{"ahvaz", "ahvaaz"}},
		{"رشت", []string{"rasht"}},
		{"همدان", []string{"hamedan", "hamadan", "hamedaan"}},
		{"قزوین", []string{"ghazvin", "gazvin"}},
		{"بندر", []string{"bandar"}},
		{"کاشان", []string{"kashan", "kaashaan"}},
		{"سمنان", []string{"semnan", "semnaan"}},
		{"گرگان", []string{"gorgan", "gorgaan"}},
		{"ساری", []string{"sari", "saari"}},
	}
}

// accuracy reports how many words transliterate to an accepted spelling.
func accuracy(t *testing.T, name string, cases []finglishCase, floor float64) {
	t.Helper()

	hit := 0
	for _, c := range cases {
		got := Finglish(c.word)
		for _, want := range c.ok {
			if got == want {
				hit++
				break
			}
		}
	}

	rate := 100 * float64(hit) / float64(len(cases))
	t.Logf("%s: %d/%d = %.1f%%", name, hit, len(cases), rate)
	if rate < floor {
		t.Errorf("%s accuracy fell to %.1f%% (%d/%d), below the %.1f%% this implementation reached",
			name, rate, hit, len(cases), floor)
	}
}

// TestFinglishAccuracy measures the transliterator against the reference word
// lists and fails if it ever does worse than it does today.
//
// The floors are the rates the current rules achieve. They are deliberately
// exact rather than generous: the point of the test is to catch a change that
// quietly makes the output worse, which a loose threshold would hide.
//
// For scale, the previous state-machine implementation scored 47.0% on the
// tuning list and 38.3% on the held-out one.
func TestFinglishAccuracy(t *testing.T) {
	t.Parallel()

	accuracy(t, "tuning", finglishReference(), 65.1)
	accuracy(t, "held-out", finglishHeldOut(), 53.3)
}
