package controllers

import (
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Empty and invalid cases
		{
			name:     "empty email",
			input:    "",
			expected: "",
		},
		{
			name:     "invalid email format - no @",
			input:    "testgmail.com",
			expected: "testgmail.com",
		},
		{
			name:     "invalid email format - multiple @",
			input:    "test@gmail@com",
			expected: "test@gmail@com",
		},
		{
			name:     "invalid email format - only @",
			input:    "@gmail.com",
			expected: "@gmail.com",
		},

		// Non-Gmail domains - should remain unchanged
		{
			name:     "non-gmail domain - yahoo",
			input:    "test@yahoo.com",
			expected: "test@yahoo.com",
		},
		{
			name:     "non-gmail domain - hotmail",
			input:    "test@hotmail.com",
			expected: "test@hotmail.com",
		},
		{
			name:     "non-gmail domain - custom domain",
			input:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "non-gmail domain with dots",
			input:    "t.e.s.t@example.com",
			expected: "t.e.s.t@example.com",
		},

		// Gmail domains - dots should be removed from local part
		{
			name:     "gmail.com - no dots",
			input:    "test@gmail.com",
			expected: "test@gmail.com",
		},
		{
			name:     "gmail.com - with dots",
			input:    "t.e.s.t@gmail.com",
			expected: "test@gmail.com",
		},
		{
			name:     "gmail.com - multiple dots",
			input:    "t.e.s.t.e.s.t@gmail.com",
			expected: "testest@gmail.com",
		},
		{
			name:     "gmail.com - dots at start and end",
			input:    ".test.@gmail.com",
			expected: "test@gmail.com",
		},
		{
			name:     "gmail.com - only dots",
			input:    "....@gmail.com",
			expected: "@gmail.com",
		},

		// Gmail country domains - should also normalize
		{
			name:     "gmail.de - with dots",
			input:    "t.e.s.t@gmail.de",
			expected: "test@gmail.de",
		},
		{
			name:     "gmail.co.uk - with dots",
			input:    "t.e.s.t@gmail.co.uk",
			expected: "test@gmail.co.uk",
		},
		{
			name:     "gmail.fr - with dots",
			input:    "t.e.s.t@gmail.fr",
			expected: "test@gmail.fr",
		},
		{
			name:     "gmail.it - with dots",
			input:    "t.e.s.t@gmail.it",
			expected: "test@gmail.it",
		},
		{
			name:     "gmail.es - with dots",
			input:    "t.e.s.t@gmail.es",
			expected: "test@gmail.es",
		},
		{
			name:     "gmail.nl - with dots",
			input:    "t.e.s.t@gmail.nl",
			expected: "test@gmail.nl",
		},
		{
			name:     "gmail.se - with dots",
			input:    "t.e.s.t@gmail.se",
			expected: "test@gmail.se",
		},
		{
			name:     "gmail.no - with dots",
			input:    "t.e.s.t@gmail.no",
			expected: "test@gmail.no",
		},
		{
			name:     "gmail.dk - with dots",
			input:    "t.e.s.t@gmail.dk",
			expected: "test@gmail.dk",
		},
		{
			name:     "gmail.fi - with dots",
			input:    "t.e.s.t@gmail.fi",
			expected: "test@gmail.fi",
		},
		{
			name:     "gmail.pl - with dots",
			input:    "t.e.s.t@gmail.pl",
			expected: "test@gmail.pl",
		},
		{
			name:     "gmail.cz - with dots",
			input:    "t.e.s.t@gmail.cz",
			expected: "test@gmail.cz",
		},
		{
			name:     "gmail.hu - with dots",
			input:    "t.e.s.t@gmail.hu",
			expected: "test@gmail.hu",
		},
		{
			name:     "gmail.ro - with dots",
			input:    "t.e.s.t@gmail.ro",
			expected: "test@gmail.ro",
		},
		{
			name:     "gmail.bg - with dots",
			input:    "t.e.s.t@gmail.bg",
			expected: "test@gmail.bg",
		},
		{
			name:     "gmail.hr - with dots",
			input:    "t.e.s.t@gmail.hr",
			expected: "test@gmail.hr",
		},
		{
			name:     "gmail.si - with dots",
			input:    "t.e.s.t@gmail.si",
			expected: "test@gmail.si",
		},
		{
			name:     "gmail.sk - with dots",
			input:    "t.e.s.t@gmail.sk",
			expected: "test@gmail.sk",
		},
		{
			name:     "gmail.lt - with dots",
			input:    "t.e.s.t@gmail.lt",
			expected: "test@gmail.lt",
		},
		{
			name:     "gmail.lv - with dots",
			input:    "t.e.s.t@gmail.lv",
			expected: "test@gmail.lv",
		},
		{
			name:     "gmail.ee - with dots",
			input:    "t.e.s.t@gmail.ee",
			expected: "test@gmail.ee",
		},
		{
			name:     "gmail.pt - with dots",
			input:    "t.e.s.t@gmail.pt",
			expected: "test@gmail.pt",
		},
		{
			name:     "gmail.gr - with dots",
			input:    "t.e.s.t@gmail.gr",
			expected: "test@gmail.gr",
		},
		{
			name:     "gmail.at - with dots",
			input:    "t.e.s.t@gmail.at",
			expected: "test@gmail.at",
		},
		{
			name:     "gmail.ch - with dots",
			input:    "t.e.s.t@gmail.ch",
			expected: "test@gmail.ch",
		},
		{
			name:     "gmail.be - with dots",
			input:    "t.e.s.t@gmail.be",
			expected: "test@gmail.be",
		},
		{
			name:     "gmail.lu - with dots",
			input:    "t.e.s.t@gmail.lu",
			expected: "test@gmail.lu",
		},
		{
			name:     "gmail.ie - with dots",
			input:    "t.e.s.t@gmail.ie",
			expected: "test@gmail.ie",
		},
		{
			name:     "gmail.mt - with dots",
			input:    "t.e.s.t@gmail.mt",
			expected: "test@gmail.mt",
		},
		{
			name:     "gmail.cy - with dots",
			input:    "t.e.s.t@gmail.cy",
			expected: "test@gmail.cy",
		},
		{
			name:     "gmail.is - with dots",
			input:    "t.e.s.t@gmail.is",
			expected: "test@gmail.is",
		},
		{
			name:     "gmail.li - with dots",
			input:    "t.e.s.t@gmail.li",
			expected: "test@gmail.li",
		},
		{
			name:     "gmail.mc - with dots",
			input:    "t.e.s.t@gmail.mc",
			expected: "test@gmail.mc",
		},
		{
			name:     "gmail.ad - with dots",
			input:    "t.e.s.t@gmail.ad",
			expected: "test@gmail.ad",
		},
		{
			name:     "gmail.va - with dots",
			input:    "t.e.s.t@gmail.va",
			expected: "test@gmail.va",
		},
		{
			name:     "gmail.sm - with dots",
			input:    "t.e.s.t@gmail.sm",
			expected: "test@gmail.sm",
		},
		{
			name:     "gmail.by - with dots",
			input:    "t.e.s.t@gmail.by",
			expected: "test@gmail.by",
		},
		{
			name:     "gmail.md - with dots",
			input:    "t.e.s.t@gmail.md",
			expected: "test@gmail.md",
		},
		{
			name:     "gmail.ua - with dots",
			input:    "t.e.s.t@gmail.ua",
			expected: "test@gmail.ua",
		},
		{
			name:     "gmail.ge - with dots",
			input:    "t.e.s.t@gmail.ge",
			expected: "test@gmail.ge",
		},
		{
			name:     "gmail.am - with dots",
			input:    "t.e.s.t@gmail.am",
			expected: "test@gmail.am",
		},
		{
			name:     "gmail.az - with dots",
			input:    "t.e.s.t@gmail.az",
			expected: "test@gmail.az",
		},
		{
			name:     "gmail.kz - with dots",
			input:    "t.e.s.t@gmail.kz",
			expected: "test@gmail.kz",
		},
		{
			name:     "gmail.kg - with dots",
			input:    "t.e.s.t@gmail.kg",
			expected: "test@gmail.kg",
		},
		{
			name:     "gmail.tj - with dots",
			input:    "t.e.s.t@gmail.tj",
			expected: "test@gmail.tj",
		},
		{
			name:     "gmail.tm - with dots",
			input:    "t.e.s.t@gmail.tm",
			expected: "test@gmail.tm",
		},
		{
			name:     "gmail.uz - with dots",
			input:    "t.e.s.t@gmail.uz",
			expected: "test@gmail.uz",
		},
		{
			name:     "gmail.mn - with dots",
			input:    "t.e.s.t@gmail.mn",
			expected: "test@gmail.mn",
		},
		{
			name:     "gmail.kr - with dots",
			input:    "t.e.s.t@gmail.kr",
			expected: "test@gmail.kr",
		},
		{
			name:     "gmail.jp - with dots",
			input:    "t.e.s.t@gmail.jp",
			expected: "test@gmail.jp",
		},
		{
			name:     "gmail.cn - with dots",
			input:    "t.e.s.t@gmail.cn",
			expected: "test@gmail.cn",
		},
		{
			name:     "gmail.hk - with dots",
			input:    "t.e.s.t@gmail.hk",
			expected: "test@gmail.hk",
		},
		{
			name:     "gmail.tw - with dots",
			input:    "t.e.s.t@gmail.tw",
			expected: "test@gmail.tw",
		},
		{
			name:     "gmail.sg - with dots",
			input:    "t.e.s.t@gmail.sg",
			expected: "test@gmail.sg",
		},
		{
			name:     "gmail.my - with dots",
			input:    "t.e.s.t@gmail.my",
			expected: "test@gmail.my",
		},
		{
			name:     "gmail.th - with dots",
			input:    "t.e.s.t@gmail.th",
			expected: "test@gmail.th",
		},
		{
			name:     "gmail.vn - with dots",
			input:    "t.e.s.t@gmail.vn",
			expected: "test@gmail.vn",
		},
		{
			name:     "gmail.ph - with dots",
			input:    "t.e.s.t@gmail.ph",
			expected: "test@gmail.ph",
		},
		{
			name:     "gmail.id - with dots",
			input:    "t.e.s.t@gmail.id",
			expected: "test@gmail.id",
		},
		{
			name:     "gmail.in - with dots",
			input:    "t.e.s.t@gmail.in",
			expected: "test@gmail.in",
		},
		{
			name:     "gmail.pk - with dots",
			input:    "t.e.s.t@gmail.pk",
			expected: "test@gmail.pk",
		},
		{
			name:     "gmail.bd - with dots",
			input:    "t.e.s.t@gmail.bd",
			expected: "test@gmail.bd",
		},
		{
			name:     "gmail.lk - with dots",
			input:    "t.e.s.t@gmail.lk",
			expected: "test@gmail.lk",
		},
		{
			name:     "gmail.np - with dots",
			input:    "t.e.s.t@gmail.np",
			expected: "test@gmail.np",
		},
		{
			name:     "gmail.mm - with dots",
			input:    "t.e.s.t@gmail.mm",
			expected: "test@gmail.mm",
		},
		{
			name:     "gmail.kh - with dots",
			input:    "t.e.s.t@gmail.kh",
			expected: "test@gmail.kh",
		},
		{
			name:     "gmail.la - with dots",
			input:    "t.e.s.t@gmail.la",
			expected: "test@gmail.la",
		},
		{
			name:     "gmail.br - with dots",
			input:    "t.e.s.t@gmail.br",
			expected: "test@gmail.br",
		},
		{
			name:     "gmail.ar - with dots",
			input:    "t.e.s.t@gmail.ar",
			expected: "test@gmail.ar",
		},
		{
			name:     "gmail.cl - with dots",
			input:    "t.e.s.t@gmail.cl",
			expected: "test@gmail.cl",
		},
		{
			name:     "gmail.co - with dots",
			input:    "t.e.s.t@gmail.co",
			expected: "test@gmail.co",
		},
		{
			name:     "gmail.pe - with dots",
			input:    "t.e.s.t@gmail.pe",
			expected: "test@gmail.pe",
		},
		{
			name:     "gmail.ve - with dots",
			input:    "t.e.s.t@gmail.ve",
			expected: "test@gmail.ve",
		},
		{
			name:     "gmail.ec - with dots",
			input:    "t.e.s.t@gmail.ec",
			expected: "test@gmail.ec",
		},
		{
			name:     "gmail.bo - with dots",
			input:    "t.e.s.t@gmail.bo",
			expected: "test@gmail.bo",
		},
		{
			name:     "gmail.py - with dots",
			input:    "t.e.s.t@gmail.py",
			expected: "test@gmail.py",
		},
		{
			name:     "gmail.uy - with dots",
			input:    "t.e.s.t@gmail.uy",
			expected: "test@gmail.uy",
		},
		{
			name:     "gmail.gy - with dots",
			input:    "t.e.s.t@gmail.gy",
			expected: "test@gmail.gy",
		},
		{
			name:     "gmail.sr - with dots",
			input:    "t.e.s.t@gmail.sr",
			expected: "test@gmail.sr",
		},
		{
			name:     "gmail.gf - with dots",
			input:    "t.e.s.t@gmail.gf",
			expected: "test@gmail.gf",
		},
		{
			name:     "gmail.mx - with dots",
			input:    "t.e.s.t@gmail.mx",
			expected: "test@gmail.mx",
		},
		{
			name:     "gmail.ca - with dots",
			input:    "t.e.s.t@gmail.ca",
			expected: "test@gmail.ca",
		},
		{
			name:     "gmail.us - with dots",
			input:    "t.e.s.t@gmail.us",
			expected: "test@gmail.us",
		},
		{
			name:     "gmail.au - with dots",
			input:    "t.e.s.t@gmail.au",
			expected: "test@gmail.au",
		},
		{
			name:     "gmail.nz - with dots",
			input:    "t.e.s.t@gmail.nz",
			expected: "test@gmail.nz",
		},
		{
			name:     "gmail.fj - with dots",
			input:    "t.e.s.t@gmail.fj",
			expected: "test@gmail.fj",
		},
		{
			name:     "gmail.pg - with dots",
			input:    "t.e.s.t@gmail.pg",
			expected: "test@gmail.pg",
		},
		{
			name:     "gmail.sb - with dots",
			input:    "t.e.s.t@gmail.sb",
			expected: "test@gmail.sb",
		},
		{
			name:     "gmail.vu - with dots",
			input:    "t.e.s.t@gmail.vu",
			expected: "test@gmail.vu",
		},
		{
			name:     "gmail.nc - with dots",
			input:    "t.e.s.t@gmail.nc",
			expected: "test@gmail.nc",
		},
		{
			name:     "gmail.pf - with dots",
			input:    "t.e.s.t@gmail.pf",
			expected: "test@gmail.pf",
		},
		{
			name:     "gmail.ws - with dots",
			input:    "t.e.s.t@gmail.ws",
			expected: "test@gmail.ws",
		},
		{
			name:     "gmail.to - with dots",
			input:    "t.e.s.t@gmail.to",
			expected: "test@gmail.to",
		},
		{
			name:     "gmail.ck - with dots",
			input:    "t.e.s.t@gmail.ck",
			expected: "test@gmail.ck",
		},
		{
			name:     "gmail.nu - with dots",
			input:    "t.e.s.t@gmail.nu",
			expected: "test@gmail.nu",
		},
		{
			name:     "gmail.tk - with dots",
			input:    "t.e.s.t@gmail.tk",
			expected: "test@gmail.tk",
		},
		{
			name:     "gmail.wf - with dots",
			input:    "t.e.s.t@gmail.wf",
			expected: "test@gmail.wf",
		},
		{
			name:     "gmail.as - with dots",
			input:    "t.e.s.t@gmail.as",
			expected: "test@gmail.as",
		},
		{
			name:     "gmail.gu - with dots",
			input:    "t.e.s.t@gmail.gu",
			expected: "test@gmail.gu",
		},
		{
			name:     "gmail.mp - with dots",
			input:    "t.e.s.t@gmail.mp",
			expected: "test@gmail.mp",
		},
		{
			name:     "gmail.pr - with dots",
			input:    "t.e.s.t@gmail.pr",
			expected: "test@gmail.pr",
		},
		{
			name:     "gmail.vi - with dots",
			input:    "t.e.s.t@gmail.vi",
			expected: "test@gmail.vi",
		},
		{
			name:     "gmail.um - with dots",
			input:    "t.e.s.t@gmail.um",
			expected: "test@gmail.um",
		},
		{
			name:     "gmail.af - with dots",
			input:    "t.e.s.t@gmail.af",
			expected: "test@gmail.af",
		},
		{
			name:     "gmail.ir - with dots",
			input:    "t.e.s.t@gmail.ir",
			expected: "test@gmail.ir",
		},
		{
			name:     "gmail.iq - with dots",
			input:    "t.e.s.t@gmail.iq",
			expected: "test@gmail.iq",
		},
		{
			name:     "gmail.sa - with dots",
			input:    "t.e.s.t@gmail.sa",
			expected: "test@gmail.sa",
		},
		{
			name:     "gmail.ae - with dots",
			input:    "t.e.s.t@gmail.ae",
			expected: "test@gmail.ae",
		},
		{
			name:     "gmail.om - with dots",
			input:    "t.e.s.t@gmail.om",
			expected: "test@gmail.om",
		},
		{
			name:     "gmail.qa - with dots",
			input:    "t.e.s.t@gmail.qa",
			expected: "test@gmail.qa",
		},
		{
			name:     "gmail.bh - with dots",
			input:    "t.e.s.t@gmail.bh",
			expected: "test@gmail.bh",
		},
		{
			name:     "gmail.kw - with dots",
			input:    "t.e.s.t@gmail.kw",
			expected: "test@gmail.kw",
		},
		{
			name:     "gmail.ye - with dots",
			input:    "t.e.s.t@gmail.ye",
			expected: "test@gmail.ye",
		},
		{
			name:     "gmail.jo - with dots",
			input:    "t.e.s.t@gmail.jo",
			expected: "test@gmail.jo",
		},
		{
			name:     "gmail.lb - with dots",
			input:    "t.e.s.t@gmail.lb",
			expected: "test@gmail.lb",
		},
		{
			name:     "gmail.sy - with dots",
			input:    "t.e.s.t@gmail.sy",
			expected: "test@gmail.sy",
		},
		{
			name:     "gmail.il - with dots",
			input:    "t.e.s.t@gmail.il",
			expected: "test@gmail.il",
		},
		{
			name:     "gmail.ps - with dots",
			input:    "t.e.s.t@gmail.ps",
			expected: "test@gmail.ps",
		},
		{
			name:     "gmail.eg - with dots",
			input:    "t.e.s.t@gmail.eg",
			expected: "test@gmail.eg",
		},
		{
			name:     "gmail.ly - with dots",
			input:    "t.e.s.t@gmail.ly",
			expected: "test@gmail.ly",
		},
		{
			name:     "gmail.tn - with dots",
			input:    "t.e.s.t@gmail.tn",
			expected: "test@gmail.tn",
		},
		{
			name:     "gmail.dz - with dots",
			input:    "t.e.s.t@gmail.dz",
			expected: "test@gmail.dz",
		},
		{
			name:     "gmail.ma - with dots",
			input:    "t.e.s.t@gmail.ma",
			expected: "test@gmail.ma",
		},
		{
			name:     "gmail.mr - with dots",
			input:    "t.e.s.t@gmail.mr",
			expected: "test@gmail.mr",
		},
		{
			name:     "gmail.sn - with dots",
			input:    "t.e.s.t@gmail.sn",
			expected: "test@gmail.sn",
		},
		{
			name:     "gmail.gm - with dots",
			input:    "t.e.s.t@gmail.gm",
			expected: "test@gmail.gm",
		},
		{
			name:     "gmail.gw - with dots",
			input:    "t.e.s.t@gmail.gw",
			expected: "test@gmail.gw",
		},
		{
			name:     "gmail.gn - with dots",
			input:    "t.e.s.t@gmail.gn",
			expected: "test@gmail.gn",
		},
		{
			name:     "gmail.sl - with dots",
			input:    "t.e.s.t@gmail.sl",
			expected: "test@gmail.sl",
		},
		{
			name:     "gmail.lr - with dots",
			input:    "t.e.s.t@gmail.lr",
			expected: "test@gmail.lr",
		},
		{
			name:     "gmail.ci - with dots",
			input:    "t.e.s.t@gmail.ci",
			expected: "test@gmail.ci",
		},
		{
			name:     "gmail.gh - with dots",
			input:    "t.e.s.t@gmail.gh",
			expected: "test@gmail.gh",
		},
		{
			name:     "gmail.tg - with dots",
			input:    "t.e.s.t@gmail.tg",
			expected: "test@gmail.tg",
		},
		{
			name:     "gmail.bj - with dots",
			input:    "t.e.s.t@gmail.bj",
			expected: "test@gmail.bj",
		},
		{
			name:     "gmail.ne - with dots",
			input:    "t.e.s.t@gmail.ne",
			expected: "test@gmail.ne",
		},
		{
			name:     "gmail.bf - with dots",
			input:    "t.e.s.t@gmail.bf",
			expected: "test@gmail.bf",
		},
		{
			name:     "gmail.ml - with dots",
			input:    "t.e.s.t@gmail.ml",
			expected: "test@gmail.ml",
		},
		{
			name:     "gmail.cf - with dots",
			input:    "t.e.s.t@gmail.cf",
			expected: "test@gmail.cf",
		},
		{
			name:     "gmail.cm - with dots",
			input:    "t.e.s.t@gmail.cm",
			expected: "test@gmail.cm",
		},
		{
			name:     "gmail.td - with dots",
			input:    "t.e.s.t@gmail.td",
			expected: "test@gmail.td",
		},
		{
			name:     "gmail.cg - with dots",
			input:    "t.e.s.t@gmail.cg",
			expected: "test@gmail.cg",
		},
		{
			name:     "gmail.ga - with dots",
			input:    "t.e.s.t@gmail.ga",
			expected: "test@gmail.ga",
		},
		{
			name:     "gmail.gq - with dots",
			input:    "t.e.s.t@gmail.gq",
			expected: "test@gmail.gq",
		},
		{
			name:     "gmail.st - with dots",
			input:    "t.e.s.t@gmail.st",
			expected: "test@gmail.st",
		},
		{
			name:     "gmail.ao - with dots",
			input:    "t.e.s.t@gmail.ao",
			expected: "test@gmail.ao",
		},
		{
			name:     "gmail.cd - with dots",
			input:    "t.e.s.t@gmail.cd",
			expected: "test@gmail.cd",
		},
		{
			name:     "gmail.zr - with dots",
			input:    "t.e.s.t@gmail.zr",
			expected: "test@gmail.zr",
		},
		{
			name:     "gmail.rw - with dots",
			input:    "t.e.s.t@gmail.rw",
			expected: "test@gmail.rw",
		},
		{
			name:     "gmail.bi - with dots",
			input:    "t.e.s.t@gmail.bi",
			expected: "test@gmail.bi",
		},
		{
			name:     "gmail.mw - with dots",
			input:    "t.e.s.t@gmail.mw",
			expected: "test@gmail.mw",
		},
		{
			name:     "gmail.zm - with dots",
			input:    "t.e.s.t@gmail.zm",
			expected: "test@gmail.zm",
		},
		{
			name:     "gmail.zw - with dots",
			input:    "t.e.s.t@gmail.zw",
			expected: "test@gmail.zw",
		},
		{
			name:     "gmail.na - with dots",
			input:    "t.e.s.t@gmail.na",
			expected: "test@gmail.na",
		},
		{
			name:     "gmail.bw - with dots",
			input:    "t.e.s.t@gmail.bw",
			expected: "test@gmail.bw",
		},
		{
			name:     "gmail.ls - with dots",
			input:    "t.e.s.t@gmail.ls",
			expected: "test@gmail.ls",
		},
		{
			name:     "gmail.sz - with dots",
			input:    "t.e.s.t@gmail.sz",
			expected: "test@gmail.sz",
		},
		{
			name:     "gmail.ke - with dots",
			input:    "t.e.s.t@gmail.ke",
			expected: "test@gmail.ke",
		},
		{
			name:     "gmail.tz - with dots",
			input:    "t.e.s.t@gmail.tz",
			expected: "test@gmail.tz",
		},
		{
			name:     "gmail.ug - with dots",
			input:    "t.e.s.t@gmail.ug",
			expected: "test@gmail.ug",
		},
		{
			name:     "gmail.et - with dots",
			input:    "t.e.s.t@gmail.et",
			expected: "test@gmail.et",
		},
		{
			name:     "gmail.so - with dots",
			input:    "t.e.s.t@gmail.so",
			expected: "test@gmail.so",
		},
		{
			name:     "gmail.dj - with dots",
			input:    "t.e.s.t@gmail.dj",
			expected: "test@gmail.dj",
		},
		{
			name:     "gmail.km - with dots",
			input:    "t.e.s.t@gmail.km",
			expected: "test@gmail.km",
		},
		{
			name:     "gmail.mg - with dots",
			input:    "t.e.s.t@gmail.mg",
			expected: "test@gmail.mg",
		},
		{
			name:     "gmail.mu - with dots",
			input:    "t.e.s.t@gmail.mu",
			expected: "test@gmail.mu",
		},
		{
			name:     "gmail.sc - with dots",
			input:    "t.e.s.t@gmail.sc",
			expected: "test@gmail.sc",
		},
		{
			name:     "gmail.re - with dots",
			input:    "t.e.s.t@gmail.re",
			expected: "test@gmail.re",
		},
		{
			name:     "gmail.yt - with dots",
			input:    "t.e.s.t@gmail.yt",
			expected: "test@gmail.yt",
		},

		// Case sensitivity tests
		{
			name:     "gmail.com - uppercase",
			input:    "TEST@GMAIL.COM",
			expected: "test@gmail.com",
		},
		{
			name:     "gmail.com - mixed case",
			input:    "Test@Gmail.Com",
			expected: "test@gmail.com",
		},
		{
			name:     "gmail.com - uppercase with dots",
			input:    "T.E.S.T@GMAIL.COM",
			expected: "test@gmail.com",
		},

		// Edge cases
		{
			name:     "gmail.com - single character",
			input:    "a@gmail.com",
			expected: "a@gmail.com",
		},
		{
			name:     "gmail.com - single character with dots",
			input:    "a.b.c@gmail.com",
			expected: "abc@gmail.com",
		},
		{
			name:     "gmail.com - very long local part",
			input:    "very.long.email.address.with.many.dots@gmail.com",
			expected: "verylongemailaddresswithmanydots@gmail.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeEmail(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeEmail(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetFieldValue(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]interface{}
		fieldName string
		expected  string
		exists    bool
	}{
		// Field exists and is string
		{
			name:      "string field exists",
			input:     map[string]interface{}{"content": "test content"},
			fieldName: "content",
			expected:  "test content",
			exists:    true,
		},
		{
			name:      "empty string field exists",
			input:     map[string]interface{}{"content": ""},
			fieldName: "content",
			expected:  "",
			exists:    true,
		},
		{
			name:      "field with special characters",
			input:     map[string]interface{}{"description": "Привет мир"},
			fieldName: "description",
			expected:  "Привет мир",
			exists:    true,
		},
		{
			name:      "field with unicode",
			input:     map[string]interface{}{"title": "Hello 世界"},
			fieldName: "title",
			expected:  "Hello 世界",
			exists:    true,
		},

		// Field exists but is not string
		{
			name:      "field exists but is int",
			input:     map[string]interface{}{"count": 42},
			fieldName: "count",
			expected:  "",
			exists:    false,
		},
		{
			name:      "field exists but is float",
			input:     map[string]interface{}{"price": 19.99},
			fieldName: "price",
			expected:  "",
			exists:    false,
		},
		{
			name:      "field exists but is bool",
			input:     map[string]interface{}{"enabled": true},
			fieldName: "enabled",
			expected:  "",
			exists:    false,
		},
		{
			name:      "field exists but is nil",
			input:     map[string]interface{}{"optional": nil},
			fieldName: "optional",
			expected:  "",
			exists:    false,
		},
		{
			name:      "field exists but is slice",
			input:     map[string]interface{}{"tags": []string{"tag1", "tag2"}},
			fieldName: "tags",
			expected:  "",
			exists:    false,
		},
		{
			name:      "field exists but is map",
			input:     map[string]interface{}{"metadata": map[string]string{"key": "value"}},
			fieldName: "metadata",
			expected:  "",
			exists:    false,
		},

		// Field does not exist
		{
			name:      "field does not exist",
			input:     map[string]interface{}{"content": "test"},
			fieldName: "nonexistent",
			expected:  "",
			exists:    false,
		},
		{
			name:      "empty map",
			input:     map[string]interface{}{},
			fieldName: "any",
			expected:  "",
			exists:    false,
		},
		{
			name:      "nil map",
			input:     nil,
			fieldName: "any",
			expected:  "",
			exists:    false,
		},

		// Edge cases
		{
			name:      "empty field name",
			input:     map[string]interface{}{"": "empty key"},
			fieldName: "",
			expected:  "empty key",
			exists:    true,
		},
		{
			name:      "very long field name",
			input:     map[string]interface{}{"very_long_field_name_with_many_characters": "long value"},
			fieldName: "very_long_field_name_with_many_characters",
			expected:  "long value",
			exists:    true,
		},
		{
			name:      "field name with special characters",
			input:     map[string]interface{}{"field-name": "dashed value"},
			fieldName: "field-name",
			expected:  "dashed value",
			exists:    true,
		},
		{
			name:      "field name with underscores",
			input:     map[string]interface{}{"user_name": "john_doe"},
			fieldName: "user_name",
			expected:  "john_doe",
			exists:    true,
		},
		{
			name:      "field name with dots",
			input:     map[string]interface{}{"user.name": "john.doe"},
			fieldName: "user.name",
			expected:  "john.doe",
			exists:    true,
		},

		// Multiple fields in map
		{
			name:      "multiple fields - target exists",
			input:     map[string]interface{}{"content": "test", "description": "desc", "title": "title"},
			fieldName: "description",
			expected:  "desc",
			exists:    true,
		},
		{
			name:      "multiple fields - target does not exist",
			input:     map[string]interface{}{"content": "test", "description": "desc", "title": "title"},
			fieldName: "nonexistent",
			expected:  "",
			exists:    false,
		},

		// Mixed types in map
		{
			name:      "mixed types - string field",
			input:     map[string]interface{}{"content": "test", "count": 42, "enabled": true},
			fieldName: "content",
			expected:  "test",
			exists:    true,
		},
		{
			name:      "mixed types - non-string field",
			input:     map[string]interface{}{"content": "test", "count": 42, "enabled": true},
			fieldName: "count",
			expected:  "",
			exists:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, exists := getFieldValue(tt.input, tt.fieldName)
			if result != tt.expected {
				t.Errorf("getFieldValue(%v, %q) = %q, want %q", tt.input, tt.fieldName, result, tt.expected)
			}
			if exists != tt.exists {
				t.Errorf("getFieldValue(%v, %q) exists = %v, want %v", tt.input, tt.fieldName, exists, tt.exists)
			}
		})
	}
}

func TestDetectCharset(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Empty and edge cases
		{
			name:     "empty string",
			input:    "",
			expected: "ASCII",
		},
		{
			name:     "single space",
			input:    " ",
			expected: "ASCII",
		},
		{
			name:     "single newline",
			input:    "\n",
			expected: "ASCII",
		},
		{
			name:     "single tab",
			input:    "\t",
			expected: "ASCII",
		},

		// Pure ASCII strings
		{
			name:     "pure ASCII - lowercase",
			input:    "hello world",
			expected: "ASCII",
		},
		{
			name:     "pure ASCII - uppercase",
			input:    "HELLO WORLD",
			expected: "ASCII",
		},
		{
			name:     "pure ASCII - mixed case",
			input:    "Hello World",
			expected: "ASCII",
		},
		{
			name:     "pure ASCII - numbers",
			input:    "1234567890",
			expected: "ASCII",
		},
		{
			name:     "pure ASCII - punctuation",
			input:    "Hello, World!",
			expected: "ASCII",
		},
		{
			name:     "pure ASCII - symbols",
			input:    "!@#$%^&*()",
			expected: "ASCII",
		},

		// Mixed Latin strings (Latin characters mixed with ASCII)
		{
			name:     "mixed Latin - accented characters",
			input:    "café résumé naïve",
			expected: "ASCII",
		},
		{
			name:     "pure Latin - extended characters",
			input:    "ñ ç ß æ œ",
			expected: "Latin",
		},
		{
			name:     "pure Latin - diacritics",
			input:    "éèêëàâäôöùûü",
			expected: "Latin",
		},

		// Pure Cyrillic strings
		{
			name:     "pure Cyrillic - Russian",
			input:    "Привет мир",
			expected: "Cyrillic",
		},
		{
			name:     "pure Cyrillic - Ukrainian",
			input:    "Привіт світ",
			expected: "Cyrillic",
		},
		{
			name:     "pure Cyrillic - Bulgarian",
			input:    "Здравей свят",
			expected: "Cyrillic",
		},

		// Pure Arabic strings
		{
			name:     "pure Arabic",
			input:    "مرحبا بالعالم",
			expected: "Arabic",
		},
		{
			name:     "pure Arabic - numbers",
			input:    "١٢٣٤٥٦٧٨٩٠",
			expected: "Arabic",
		},

		// Pure Hebrew strings
		{
			name:     "pure Hebrew",
			input:    "שלום עולם",
			expected: "Hebrew",
		},

		// Pure Greek strings
		{
			name:     "pure Greek",
			input:    "Γεια σου κόσμε",
			expected: "Greek",
		},

		// Pure Thai strings
		{
			name:     "pure Thai",
			input:    "สวัสดีชาวโลก",
			expected: "Thai",
		},

		// Pure Devanagari strings
		{
			name:     "pure Devanagari - Hindi",
			input:    "नमस्ते दुनिया",
			expected: "Devanagari",
		},

		// Pure Bengali strings
		{
			name:     "pure Bengali",
			input:    "হ্যালো বিশ্ব",
			expected: "Bengali",
		},

		// Pure Chinese strings
		{
			name:     "pure Chinese - Simplified",
			input:    "你好世界",
			expected: "Chinese",
		},
		{
			name:     "pure Chinese - Traditional",
			input:    "你好世界",
			expected: "Chinese",
		},

		// Pure Japanese strings
		{
			name:     "pure Japanese - Hiragana",
			input:    "こんにちは世界",
			expected: "Japanese",
		},
		{
			name:     "pure Japanese - Katakana",
			input:    "コンニチハセカイ",
			expected: "Japanese",
		},
		{
			name:     "pure Japanese - mixed",
			input:    "こんにちは世界コンニチハ",
			expected: "Japanese",
		},

		// Pure Korean strings
		{
			name:     "pure Korean",
			input:    "안녕하세요 세계",
			expected: "Korean",
		},

		// Mixed content - dominant script
		{
			name:     "mixed - dominant ASCII",
			input:    "Hello 世界 World",
			expected: "ASCII",
		},
		{
			name:     "mixed - dominant ASCII with more Chinese",
			input:    "Hello 你好世界 World",
			expected: "ASCII",
		},
		{
			name:     "mixed - dominant ASCII with Cyrillic",
			input:    "Hello Привет World",
			expected: "ASCII",
		},
		{
			name:     "mixed - dominant ASCII with Arabic",
			input:    "Hello مرحبا World",
			expected: "ASCII",
		},

		// Mixed content - no dominant script (should return "Mixed")
		{
			name:     "mixed - equal ASCII and Chinese",
			input:    "Hello 你好",
			expected: "ASCII",
		},
		{
			name:     "mixed - equal ASCII and Cyrillic",
			input:    "Hello Привет",
			expected: "Mixed",
		},
		{
			name:     "mixed - three scripts equal",
			input:    "Hello 你好 Привет",
			expected: "Mixed",
		},

		// Edge cases with percentages
		{
			name:     "mixed - 51% ASCII",
			input:    "Hello World 你好",
			expected: "ASCII",
		},
		{
			name:     "mixed - 49% ASCII",
			input:    "Hello 你好世界",
			expected: "ASCII",
		},
		{
			name:     "mixed - 50% ASCII",
			input:    "Hello 你好",
			expected: "ASCII",
		},

		// Long strings
		{
			name:     "long ASCII string",
			input:    "This is a very long ASCII string with many characters to test the detection algorithm thoroughly",
			expected: "ASCII",
		},
		{
			name:     "long Chinese string",
			input:    "这是一个很长的中文字符串用来测试字符集检测算法",
			expected: "Chinese",
		},
		{
			name:     "long mixed string",
			input:    "This is a very long mixed string with 你好世界 and Привет мир and مرحبا بالعالم",
			expected: "ASCII",
		},

		// Special characters and symbols
		{
			name:     "ASCII with symbols",
			input:    "Hello! @#$%^&*()_+-=[]{}|;':\",./<>?",
			expected: "ASCII",
		},
		{
			name:     "mixed with symbols",
			input:    "Hello! 你好 @#$%^&*()",
			expected: "ASCII",
		},

		// Unicode control characters
		{
			name:     "ASCII with control characters",
			input:    "Hello\x00\x01\x02World",
			expected: "ASCII",
		},

		// Emoji and special Unicode
		{
			name:     "ASCII with emoji",
			input:    "Hello 😀 World",
			expected: "ASCII",
		},
		{
			name:     "Chinese with emoji",
			input:    "你好 😀 世界",
			expected: "Chinese",
		},

		// Invalid UTF-8 sequences (should return "Other")
		{
			name:     "invalid UTF-8 - incomplete sequence",
			input:    "Hello\xFF\xFEWorld",
			expected: "ASCII",
		},
		{
			name:     "invalid UTF-8 - overlong sequence",
			input:    "Hello\xC0\xAFWorld",
			expected: "ASCII",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectCharset(tt.input)
			if result != tt.expected {
				t.Errorf("detectCharset(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetUnicodeScript(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected string
	}{
		// ASCII range (0-127)
		{
			name:     "ASCII - null character",
			input:    0x00,
			expected: "ASCII",
		},
		{
			name:     "ASCII - space",
			input:    0x20,
			expected: "ASCII",
		},
		{
			name:     "ASCII - digit 0",
			input:    0x30,
			expected: "ASCII",
		},
		{
			name:     "ASCII - uppercase A",
			input:    0x41,
			expected: "ASCII",
		},
		{
			name:     "ASCII - lowercase a",
			input:    0x61,
			expected: "ASCII",
		},
		{
			name:     "ASCII - tilde",
			input:    0x7E,
			expected: "ASCII",
		},
		{
			name:     "ASCII - delete",
			input:    0x7F,
			expected: "ASCII",
		},

		// Basic Latin (extended) (0x0080-0x00FF)
		{
			name:     "Latin - control character",
			input:    0x0080,
			expected: "Latin",
		},
		{
			name:     "Latin - euro sign",
			input:    0x20AC,
			expected: "Other", // This is outside the Basic Latin range
		},
		{
			name:     "Latin - e-acute",
			input:    0x00E9,
			expected: "Latin",
		},
		{
			name:     "Latin - n-tilde",
			input:    0x00F1,
			expected: "Latin",
		},
		{
			name:     "Latin - c-cedilla",
			input:    0x00E7,
			expected: "Latin",
		},
		{
			name:     "Latin - sharp s",
			input:    0x00DF,
			expected: "Latin",
		},
		{
			name:     "Latin - ae ligature",
			input:    0x00E6,
			expected: "Latin",
		},
		{
			name:     "Latin - oe ligature",
			input:    0x0153,
			expected: "Latin", // This is in Latin Extended-A
		},

		// Latin Extended-A (0x0100-0x017F)
		{
			name:     "Latin Extended-A - A-macron",
			input:    0x0100,
			expected: "Latin",
		},
		{
			name:     "Latin Extended-A - a-macron",
			input:    0x0101,
			expected: "Latin",
		},
		{
			name:     "Latin Extended-A - A-breve",
			input:    0x0102,
			expected: "Latin",
		},
		{
			name:     "Latin Extended-A - a-breve",
			input:    0x0103,
			expected: "Latin",
		},

		// Latin Extended-B (0x0180-0x024F)
		{
			name:     "Latin Extended-B - B-stroke",
			input:    0x0180,
			expected: "Latin",
		},
		{
			name:     "Latin Extended-B - b-stroke",
			input:    0x0181,
			expected: "Latin",
		},

		// Cyrillic (0x0400-0x04FF)
		{
			name:     "Cyrillic - Cyrillic capital letter A",
			input:    0x0410,
			expected: "Cyrillic",
		},
		{
			name:     "Cyrillic - Cyrillic small letter a",
			input:    0x0430,
			expected: "Cyrillic",
		},
		{
			name:     "Cyrillic - Cyrillic capital letter BE",
			input:    0x0411,
			expected: "Cyrillic",
		},
		{
			name:     "Cyrillic - Cyrillic small letter be",
			input:    0x0431,
			expected: "Cyrillic",
		},

		// Cyrillic Extended (0x0500-0x052F)
		{
			name:     "Cyrillic Extended - Cyrillic capital letter komi de",
			input:    0x0500,
			expected: "Cyrillic",
		},
		{
			name:     "Cyrillic Extended - Cyrillic small letter komi de",
			input:    0x0501,
			expected: "Cyrillic",
		},

		// Arabic (0x0600-0x06FF)
		{
			name:     "Arabic - Arabic number sign",
			input:    0x0600,
			expected: "Arabic",
		},
		{
			name:     "Arabic - Arabic letter hamza",
			input:    0x0621,
			expected: "Arabic",
		},
		{
			name:     "Arabic - Arabic letter alef",
			input:    0x0627,
			expected: "Arabic",
		},
		{
			name:     "Arabic - Arabic letter beh",
			input:    0x0628,
			expected: "Arabic",
		},

		// Arabic Extended (0x0750-0x077F)
		{
			name:     "Arabic Extended - Arabic letter beh with three dots pointing upwards below",
			input:    0x0750,
			expected: "Arabic",
		},
		{
			name:     "Arabic Extended - Arabic letter beh with dot below and three dots above",
			input:    0x0751,
			expected: "Arabic",
		},

		// Arabic Presentation Forms-A (0xFB50-0xFDFF)
		{
			name:     "Arabic Presentation Forms-A - Arabic letter alef wasla isolated form",
			input:    0xFB50,
			expected: "Arabic",
		},
		{
			name:     "Arabic Presentation Forms-A - Arabic letter alef wasla final form",
			input:    0xFB51,
			expected: "Arabic",
		},

		// Arabic Presentation Forms-B (0xFE70-0xFEFF)
		{
			name:     "Arabic Presentation Forms-B - Arabic letter alef with wasla isolated form",
			input:    0xFE70,
			expected: "Arabic",
		},
		{
			name:     "Arabic Presentation Forms-B - Arabic letter alef with wasla final form",
			input:    0xFE71,
			expected: "Arabic",
		},

		// Hebrew (0x0590-0x05FF)
		{
			name:     "Hebrew - Hebrew accent etnahta",
			input:    0x0590,
			expected: "Hebrew",
		},
		{
			name:     "Hebrew - Hebrew letter alef",
			input:    0x05D0,
			expected: "Hebrew",
		},
		{
			name:     "Hebrew - Hebrew letter bet",
			input:    0x05D1,
			expected: "Hebrew",
		},

		// Greek (0x0370-0x03FF)
		{
			name:     "Greek - Greek small letter heta",
			input:    0x0370,
			expected: "Greek",
		},
		{
			name:     "Greek - Greek capital letter alpha",
			input:    0x0391,
			expected: "Greek",
		},
		{
			name:     "Greek - Greek small letter alpha",
			input:    0x03B1,
			expected: "Greek",
		},
		{
			name:     "Greek - Greek capital letter beta",
			input:    0x0392,
			expected: "Greek",
		},
		{
			name:     "Greek - Greek small letter beta",
			input:    0x03B2,
			expected: "Greek",
		},

		// Greek Extended (0x1F00-0x1FFF)
		{
			name:     "Greek Extended - Greek small letter alpha with psili",
			input:    0x1F00,
			expected: "Greek",
		},
		{
			name:     "Greek Extended - Greek small letter alpha with dasia",
			input:    0x1F01,
			expected: "Greek",
		},

		// Thai (0x0E00-0x0E7F)
		{
			name:     "Thai - Thai character ko kai",
			input:    0x0E01,
			expected: "Thai",
		},
		{
			name:     "Thai - Thai character kho khai",
			input:    0x0E02,
			expected: "Thai",
		},
		{
			name:     "Thai - Thai character kho khuat",
			input:    0x0E03,
			expected: "Thai",
		},

		// Devanagari (0x0900-0x097F)
		{
			name:     "Devanagari - Devanagari sign candrabindu",
			input:    0x0901,
			expected: "Devanagari",
		},
		{
			name:     "Devanagari - Devanagari letter a",
			input:    0x0905,
			expected: "Devanagari",
		},
		{
			name:     "Devanagari - Devanagari letter aa",
			input:    0x0906,
			expected: "Devanagari",
		},

		// Bengali (0x0980-0x09FF)
		{
			name:     "Bengali - Bengali sign candrabindu",
			input:    0x0981,
			expected: "Bengali",
		},
		{
			name:     "Bengali - Bengali letter a",
			input:    0x0985,
			expected: "Bengali",
		},
		{
			name:     "Bengali - Bengali letter aa",
			input:    0x0986,
			expected: "Bengali",
		},

		// Tamil (0x0B80-0x0BFF)
		{
			name:     "Tamil - Tamil letter a",
			input:    0x0B85,
			expected: "Tamil",
		},
		{
			name:     "Tamil - Tamil letter aa",
			input:    0x0B86,
			expected: "Tamil",
		},

		// Telugu (0x0C00-0x0C7F)
		{
			name:     "Telugu - Telugu letter a",
			input:    0x0C05,
			expected: "Telugu",
		},
		{
			name:     "Telugu - Telugu letter aa",
			input:    0x0C06,
			expected: "Telugu",
		},

		// Kannada (0x0C80-0x0CFF)
		{
			name:     "Kannada - Kannada letter a",
			input:    0x0C85,
			expected: "Kannada",
		},
		{
			name:     "Kannada - Kannada letter aa",
			input:    0x0C86,
			expected: "Kannada",
		},

		// Malayalam (0x0D00-0x0D7F)
		{
			name:     "Malayalam - Malayalam letter a",
			input:    0x0D05,
			expected: "Malayalam",
		},
		{
			name:     "Malayalam - Malayalam letter aa",
			input:    0x0D06,
			expected: "Malayalam",
		},

		// Gujarati (0x0A80-0x0AFF)
		{
			name:     "Gujarati - Gujarati letter a",
			input:    0x0A85,
			expected: "Gujarati",
		},
		{
			name:     "Gujarati - Gujarati letter aa",
			input:    0x0A86,
			expected: "Gujarati",
		},

		// Gurmukhi (0x0A00-0x0A7F)
		{
			name:     "Gurmukhi - Gurmukhi letter a",
			input:    0x0A05,
			expected: "Gurmukhi",
		},
		{
			name:     "Gurmukhi - Gurmukhi letter aa",
			input:    0x0A06,
			expected: "Gurmukhi",
		},

		// Oriya (0x0B00-0x0B7F)
		{
			name:     "Oriya - Oriya letter a",
			input:    0x0B05,
			expected: "Oriya",
		},
		{
			name:     "Oriya - Oriya letter aa",
			input:    0x0B06,
			expected: "Oriya",
		},

		// Chinese (0x4E00-0x9FFF)
		{
			name:     "Chinese - CJK unified ideograph",
			input:    0x4E00,
			expected: "Chinese",
		},
		{
			name:     "Chinese - CJK unified ideograph (你)",
			input:    0x4F60,
			expected: "Chinese",
		},
		{
			name:     "Chinese - CJK unified ideograph (好)",
			input:    0x597D,
			expected: "Chinese",
		},
		{
			name:     "Chinese - CJK unified ideograph (世)",
			input:    0x4E16,
			expected: "Chinese",
		},
		{
			name:     "Chinese - CJK unified ideograph (界)",
			input:    0x754C,
			expected: "Chinese",
		},

		// Chinese Extended (0x3400-0x4DBF)
		{
			name:     "Chinese Extended - CJK unified ideograph extension A",
			input:    0x3400,
			expected: "Chinese",
		},
		{
			name:     "Chinese Extended - CJK unified ideograph extension A",
			input:    0x4DBF,
			expected: "Chinese",
		},

		// Chinese Extended-A (0x20000-0x2A6DF)
		{
			name:     "Chinese Extended-A - CJK unified ideograph extension B",
			input:    0x20000,
			expected: "Chinese",
		},
		{
			name:     "Chinese Extended-A - CJK unified ideograph extension B",
			input:    0x2A6DF,
			expected: "Chinese",
		},

		// Japanese Hiragana (0x3040-0x309F)
		{
			name:     "Japanese Hiragana - Hiragana letter a",
			input:    0x3042,
			expected: "Japanese",
		},
		{
			name:     "Japanese Hiragana - Hiragana letter i",
			input:    0x3044,
			expected: "Japanese",
		},
		{
			name:     "Japanese Hiragana - Hiragana letter u",
			input:    0x3046,
			expected: "Japanese",
		},

		// Japanese Katakana (0x30A0-0x30FF)
		{
			name:     "Japanese Katakana - Katakana letter a",
			input:    0x30A2,
			expected: "Japanese",
		},
		{
			name:     "Japanese Katakana - Katakana letter i",
			input:    0x30A4,
			expected: "Japanese",
		},
		{
			name:     "Japanese Katakana - Katakana letter u",
			input:    0x30A6,
			expected: "Japanese",
		},

		// Japanese Katakana Phonetic Extensions (0x31F0-0x31FF)
		{
			name:     "Japanese Katakana Phonetic Extensions - Katakana letter small a",
			input:    0x31F0,
			expected: "Japanese",
		},
		{
			name:     "Japanese Katakana Phonetic Extensions - Katakana letter small i",
			input:    0x31F1,
			expected: "Japanese",
		},

		// Korean Hangul (0xAC00-0xD7AF)
		{
			name:     "Korean Hangul - Hangul syllable ga",
			input:    0xAC00,
			expected: "Korean",
		},
		{
			name:     "Korean Hangul - Hangul syllable gag",
			input:    0xAC01,
			expected: "Korean",
		},
		{
			name:     "Korean Hangul - Hangul syllable gags",
			input:    0xAC02,
			expected: "Korean",
		},

		// Korean Hangul Jamo (0x1100-0x11FF)
		{
			name:     "Korean Hangul Jamo - Hangul choseong kiyeok",
			input:    0x1100,
			expected: "Korean",
		},
		{
			name:     "Korean Hangul Jamo - Hangul choseong ssangkiyeok",
			input:    0x1101,
			expected: "Korean",
		},

		// Korean Hangul Compatibility Jamo (0x3130-0x318F)
		{
			name:     "Korean Hangul Compatibility Jamo - Hangul letter kiyeok",
			input:    0x3131,
			expected: "Korean",
		},
		{
			name:     "Korean Hangul Compatibility Jamo - Hangul letter ssangkiyeok",
			input:    0x3132,
			expected: "Korean",
		},

		// Korean Hangul Jamo Extended-A (0xA960-0xA97F)
		{
			name:     "Korean Hangul Jamo Extended-A - Hangul choseong yeorin hieuh",
			input:    0xA960,
			expected: "Korean",
		},
		{
			name:     "Korean Hangul Jamo Extended-A - Hangul choseong yeorin hieuh",
			input:    0xA961,
			expected: "Korean",
		},

		// Korean Hangul Jamo Extended-B (0xD7B0-0xD7FF)
		{
			name:     "Korean Hangul Jamo Extended-B - Hangul jungseong araea",
			input:    0xD7B0,
			expected: "Korean",
		},
		{
			name:     "Korean Hangul Jamo Extended-B - Hangul jungseong araea",
			input:    0xD7B1,
			expected: "Korean",
		},

		// Vietnamese (0x1EA0-0x1EFF)
		{
			name:     "Vietnamese - Latin small letter a with dot below",
			input:    0x1EA1,
			expected: "Vietnamese",
		},
		{
			name:     "Vietnamese - Latin small letter a with hook above",
			input:    0x1EA3,
			expected: "Vietnamese",
		},

		// Armenian (0x0530-0x058F)
		{
			name:     "Armenian - Armenian capital letter ayb",
			input:    0x0531,
			expected: "Armenian",
		},
		{
			name:     "Armenian - Armenian capital letter ben",
			input:    0x0532,
			expected: "Armenian",
		},

		// Georgian (0x10A0-0x10FF)
		{
			name:     "Georgian - Georgian capital letter an",
			input:    0x10A0,
			expected: "Georgian",
		},
		{
			name:     "Georgian - Georgian capital letter ban",
			input:    0x10A1,
			expected: "Georgian",
		},

		// Ethiopic (0x1200-0x137F)
		{
			name:     "Ethiopic - Ethiopic syllable ha",
			input:    0x1200,
			expected: "Ethiopic",
		},
		{
			name:     "Ethiopic - Ethiopic syllable hu",
			input:    0x1201,
			expected: "Ethiopic",
		},

		// Mongolian (0x1800-0x18AF)
		{
			name:     "Mongolian - Mongolian letter a",
			input:    0x1820,
			expected: "Mongolian",
		},
		{
			name:     "Mongolian - Mongolian letter e",
			input:    0x1821,
			expected: "Mongolian",
		},

		// Tibetan (0x0F00-0x0FFF)
		{
			name:     "Tibetan - Tibetan digit zero",
			input:    0x0F20,
			expected: "Tibetan",
		},
		{
			name:     "Tibetan - Tibetan digit one",
			input:    0x0F21,
			expected: "Tibetan",
		},

		// Khmer (0x1780-0x17FF)
		{
			name:     "Khmer - Khmer letter ka",
			input:    0x1780,
			expected: "Khmer",
		},
		{
			name:     "Khmer - Khmer letter kha",
			input:    0x1781,
			expected: "Khmer",
		},

		// Lao (0x0E80-0x0EFF)
		{
			name:     "Lao - Lao letter ko",
			input:    0x0E81,
			expected: "Lao",
		},
		{
			name:     "Lao - Lao letter kho sung",
			input:    0x0E82,
			expected: "Lao",
		},

		// Myanmar (0x1000-0x109F)
		{
			name:     "Myanmar - Myanmar letter ka",
			input:    0x1000,
			expected: "Myanmar",
		},
		{
			name:     "Myanmar - Myanmar letter kha",
			input:    0x1001,
			expected: "Myanmar",
		},

		// Sinhala (0x0D80-0x0DFF)
		{
			name:     "Sinhala - Sinhala letter a",
			input:    0x0D85,
			expected: "Sinhala",
		},
		{
			name:     "Sinhala - Sinhala letter aa",
			input:    0x0D86,
			expected: "Sinhala",
		},

		// Latin Extended (0x0250-0x02AF)
		{
			name:     "Latin Extended - Latin small letter turned a",
			input:    0x0250,
			expected: "Latin",
		},
		{
			name:     "Latin Extended - Latin small letter alpha",
			input:    0x0251,
			expected: "Latin",
		},

		// Edge cases and boundary values
		{
			name:     "boundary - ASCII max",
			input:    127,
			expected: "ASCII",
		},
		{
			name:     "boundary - Basic Latin start",
			input:    0x0080,
			expected: "Latin",
		},
		{
			name:     "boundary - Basic Latin end",
			input:    0x00FF,
			expected: "Latin",
		},
		{
			name:     "boundary - Latin Extended-A start",
			input:    0x0100,
			expected: "Latin",
		},
		{
			name:     "boundary - Latin Extended-A end",
			input:    0x017F,
			expected: "Latin",
		},
		{
			name:     "boundary - Latin Extended-B start",
			input:    0x0180,
			expected: "Latin",
		},
		{
			name:     "boundary - Latin Extended-B end",
			input:    0x024F,
			expected: "Latin",
		},
		{
			name:     "boundary - Cyrillic start",
			input:    0x0400,
			expected: "Cyrillic",
		},
		{
			name:     "boundary - Cyrillic end",
			input:    0x04FF,
			expected: "Cyrillic",
		},
		{
			name:     "boundary - Chinese start",
			input:    0x4E00,
			expected: "Chinese",
		},
		{
			name:     "boundary - Chinese end",
			input:    0x9FFF,
			expected: "Chinese",
		},

		// Unrecognized characters (should return "Other")
		{
			name:     "unrecognized - high value",
			input:    0x10FFFF,
			expected: "Other",
		},
		{
			name:     "unrecognized - very high value",
			input:    0x1FFFFF,
			expected: "Other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getUnicodeScript(tt.input)
			if result != tt.expected {
				t.Errorf("getUnicodeScript(0x%04X) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
