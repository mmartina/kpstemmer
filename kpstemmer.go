package kpstemmer

var (
	g_v    Bitmask = []byte{17, 65, 16, 1}  // 'aeiouy'
	g_v_WX Bitmask = []byte{17, 65, 208, 1} // v + 'wx'
	g_AOU  Bitmask = []byte{1, 64, 16}      // 'aou'
	g_AIOU Bitmask = []byte{1, 65, 16}      // 'aiou'
)

const (
	rune_a rune = 'a'
	rune_u rune = 'u'
	rune_y rune = 'y'
)

type KpStemmer struct {
	snowballProgram
	b_GE_removed bool
	i_p1         int
	i_p2         int
}

// define C  as test (not 'ij' non-v)
func (s *KpStemmer) r_C() bool {
	r := s.limit - s.cursor
	if !s.eq_s_b("ij") && s.out_grouping_b(g_v, rune_a, rune_y) {
		s.cursor = s.limit - r
		return true
	}
	return false
}

// define Lose_infix as
func (s *KpStemmer) r_Lose_infix() bool {
	// next
	{
		if s.cursor >= s.limit {
			return false
		}
		s.cursor++
	}
	// gopast (['ge']) test hop 3 (goto v goto non-v)
	for ; ; s.cursor++ {
		s.bra = s.cursor
		if s.eq_s("ge") {
			s.ket = s.cursor
			break
		}
		if s.cursor >= s.limit {
			return false
		}
	}
	if s.cursor+3 <= s.limit {
		for ; ; s.cursor++ {
			c := s.cursor
			if s.in_grouping(g_v, rune_a, rune_y) {
				s.cursor = c
				break
			}
			if s.cursor >= s.limit {
				return false
			}
		}
		for ; ; s.cursor++ {
			c := s.cursor
			if s.out_grouping(g_v, rune_a, rune_y) {
				s.cursor = c
				break
			}
			if s.cursor >= s.limit {
				return false
			}
		}
		s.b_GE_removed = true // set GE_removed
		s.slice_del()         // delete
		return true
	}
	return false
}

// define Lose_prefix as
func (s *KpStemmer) r_Lose_prefix() bool {
	// ['ge'] test hop 3 (goto v goto non-v)
	s.bra = s.cursor
	if s.eq_s("ge") {
		s.ket = s.cursor
		if s.cursor+3 <= s.limit {
			for ; ; s.cursor++ {
				c := s.cursor
				if s.in_grouping(g_v, rune_a, rune_y) {
					s.cursor = c
					break
				}
				if s.cursor >= s.limit {
					return false
				}
			}
			for ; ; s.cursor++ {
				c := s.cursor
				if s.out_grouping(g_v, rune_a, rune_y) {
					s.cursor = c
					break
				}
				if s.cursor >= s.limit {
					return false
				}
			}
			s.b_GE_removed = true // set GE_removed
			s.slice_del()         // delete
			return true
		}
	}
	return false
}

// define R1 as (setmark x $x >= p1)
func (s *KpStemmer) r_R1() bool {
	return s.cursor >= s.i_p1
}

// define R2 as (setmark x $x >= p2)
func (s *KpStemmer) r_R2() bool {
	return s.cursor >= s.i_p2
}

var a_step_1 = []*Among{
	NewAmong("nde", -1, 7),
	NewAmong("en", -1, 6),
	NewAmong("s", -1, 2),
	NewAmong("'s", 2, 1),
	NewAmong("es", 2, 4),
	NewAmong("ies", 4, 3),
	NewAmong("aus", 2, 5),
}

// define Step_1 as
func (s *KpStemmer) r_Step_1() bool {
	// [among ( (])
	s.ket = s.cursor
	among_var := s.find_among_b(a_step_1)
	if among_var > 0 {
		s.bra = s.cursor
		switch among_var {
		case 1:
			// '{'}s' (delete)
			s.slice_del()
			return true
		case 2:
			// 's'    (R1 not ('t' R1) C delete)
			if s.r_R1() {
				r := s.limit - s.cursor
				if !s.eq_s_b("t") || !s.r_R1() {
					s.cursor = s.limit - r
					if s.r_C() {
						s.slice_del()
						return true
					}
				}
			}
		case 3:
			// 'ies'  (R1 <-'ie')
			if s.r_R1() {
				s.slice_from("ie")
				return true
			}
		case 4:
			// 'es'
			//	      (('ar' R1 C ] delete lengthen_V) or
			//	       ('er' R1 C ] delete) or
			//	       (R1 C <-'e'))
			r := s.limit - s.cursor
			if s.eq_s_b("ar") && s.r_R1() && s.r_C() {
				s.bra = s.cursor
				s.slice_del()
				s.r_lengthen_V()
				return true
			}
			s.cursor = s.limit - r
			if s.eq_s_b("er") && s.r_R1() && s.r_C() {
				s.bra = s.cursor
				s.slice_del()
				return true
			}
			s.cursor = s.limit - r
			if s.r_R1() && s.r_C() {
				s.slice_from("e")
				return true
			}
		case 5:
			// 'aus'  (R1 V <-'au')
			if s.r_R1() && s.r_V() {
				s.slice_from("au")
				return true
			}
		case 6:
			// 'en'   (('hed' R1 ] <-'heid') or
			//	       ('nd' delete) or
			//	       ('d' R1 C ] delete) or
			//	       ('i' or 'j' V delete) or
			//	       (R1 C delete lengthen_V))
			r := s.limit - s.cursor
			if s.eq_s_b("hed") && s.r_R1() {
				s.bra = s.cursor
				s.slice_from("heid")
				return true
			}
			s.cursor = s.limit - r
			if s.eq_s_b("nd") {
				s.slice_del()
				return true
			}
			s.cursor = s.limit - r
			if s.eq_s_b("d") && s.r_R1() && s.r_C() {
				s.bra = s.cursor
				s.slice_del()
				return true
			}
			s.cursor = s.limit - r
			if (s.eq_s_b("i") || s.eq_s_b("j")) && s.r_V() {
				s.slice_del()
				return true
			}
			s.cursor = s.limit - r
			if s.r_R1() && s.r_C() {
				s.slice_del()
				s.r_lengthen_V()
				return true

			}
		case 7:
			// 'nde'  (<-'nd')
			s.slice_from("nd")
			return true
		}
	}
	// )
	return false
}

var a_step_1c = []*Among{
	NewAmong("d", -1, 1),
	NewAmong("t", -1, 2),
}

func (s *KpStemmer) r_Step_1c() {
	// [among ( (] R1 C)
	s.ket = s.cursor
	among_var := s.find_among_b(a_step_1c)
	if among_var > 0 {
		s.bra = s.cursor
		if s.r_R1() && s.r_C() {
			switch among_var {
			case 1:
				// 'd' (not ('n' R1) delete)
				r := s.limit - s.cursor
				if !s.eq_s_b("n") || !s.r_R1() {
					s.cursor = s.limit - r
					s.slice_del()
					return
				}
			case 2:
				// 't' (not ('h' R1) delete)
				r := s.limit - s.cursor
				if !s.eq_s_b("h") || !s.r_R1() {
					s.cursor = s.limit - r
					s.slice_del()
					return
				}
			}
		}
	}
	// )
}

var a_step_2 = []*Among{
	NewAmong("de", -1, 5),
	NewAmong("ge", -1, 2),
	NewAmong("ische", -1, 4),
	NewAmong("je", -1, 1),
	NewAmong("lijke", -1, 3),
	NewAmong("le", -1, 9),
	NewAmong("ene", -1, 10),
	NewAmong("re", -1, 8),
	NewAmong("se", -1, 7),
	NewAmong("te", -1, 6),
	NewAmong("ieve", -1, 11),
}

// define Step_2 as
func (s *KpStemmer) r_Step_2() bool {
	// [among ( (])
	s.ket = s.cursor
	among_var := s.find_among_b(a_step_2)
	if among_var > 0 {
		s.bra = s.cursor
		switch among_var {
		case 1:
			// 'je'   (('{'}t' ] delete) or
			//         ('et'   ] R1 C delete) or
			//         ('rnt'  ] <-'rn') or
			//         ('t'    ] R1 VX delete) or
			//         ('ink'  ] <-'ing') or
			//         ('mp'   ] <-'m') or
			//         ('{'}'  ] R1 delete) or
			// 	       (] R1 C delete))
			r := s.limit - s.cursor
			if s.eq_s_b("'t") {
				s.bra = s.cursor
				s.slice_del()
				return true
			}
			s.cursor = s.limit - r
			if s.eq_s_b("et") {
				s.bra = s.cursor
				if s.r_R1() && s.r_C() {
					s.slice_del()
					return true
				}
			}
			s.cursor = s.limit - r
			if s.eq_s_b("rnt") {
				s.bra = s.cursor
				s.slice_from("rn")
				return true
			}
			s.cursor = s.limit - r
			if s.eq_s_b("t") {
				s.bra = s.cursor
				if s.r_R1() && s.r_VX() {
					s.slice_del()
					return true
				}
			}
			s.cursor = s.limit - r
			if s.eq_s_b("ink") {
				s.bra = s.cursor
				s.slice_from("ing")
				return true
			}
			s.cursor = s.limit - r
			if s.eq_s_b("mp") {
				s.bra = s.cursor
				s.slice_from("m")
				return true
			}
			s.cursor = s.limit - r
			if s.eq_s_b("'") {
				s.bra = s.cursor
				if s.r_R1() {
					s.slice_del()
					return true
				}
			}
			s.cursor = s.limit - r
			{
				s.bra = s.cursor
				if s.r_R1() && s.r_C() {
					s.slice_del()
					return true
				}
			}
		case 2:
			// 'ge'   (R1 <-'g')
			if s.r_R1() {
				s.slice_from("g")
				return true
			}
		case 3:
			// 'lijke'(R1 <-'lijk')
			if s.r_R1() {
				s.slice_from("lijk")
				return true
			}
		case 4:
			// 'ische'(R1 <-'isch')
			if s.r_R1() {
				s.slice_from("isch")
				return true
			}
		case 5:
			// 'de'   (R1 C delete)
			if s.r_R1() && s.r_C() {
				s.slice_del()
				return true
			}
		case 6:
			// 'te'   (R1 <-'t')
			if s.r_R1() {
				s.slice_from("t")
				return true
			}
		case 7:
			// 'se'   (R1 <-'s')
			if s.r_R1() {
				s.slice_from("s")
				return true
			}
		case 8:
			// 're'   (R1 <-'r')
			if s.r_R1() {
				s.slice_from("r")
				return true
			}
		case 9:
			// 'le'   (R1 delete attach 'l' lengthen_V)
			if s.r_R1() {
				s.slice_del()
				s.insert(s.cursor, s.cursor, "l")
				s.r_lengthen_V()
				return true
			}
		case 10:
			// 'ene'  (R1 C delete attach 'en' lengthen_V)
			if s.r_R1() && s.r_C() {
				s.slice_del()
				s.insert(s.cursor, s.cursor, "en")
				s.r_lengthen_V()
				return true
			}
		case 11:
			// 'ieve' (R1 C <-'ief')
			if s.r_R1() && s.r_C() {
				s.slice_from("ief")
				return true
			}
		}
	}
	// )
	return false
}

var a_step_3 = []*Among{
	NewAmong("heid", -1, 3),
	NewAmong("fie", -1, 7),
	NewAmong("gie", -1, 8),
	NewAmong("atie", -1, 1),
	NewAmong("isme", -1, 5),
	NewAmong("ing", -1, 5),
	NewAmong("arij", -1, 6),
	NewAmong("erij", -1, 5),
	NewAmong("sel", -1, 3),
	NewAmong("rder", -1, 4),
	NewAmong("ster", -1, 3),
	NewAmong("iteit", -1, 2),
	NewAmong("dst", -1, 10),
	NewAmong("tst", -1, 9),
}

// define Step_3 as
func (s *KpStemmer) r_Step_3() bool {
	// [among ( (])
	s.ket = s.cursor
	among_var := s.find_among_b(a_step_3)
	if among_var > 0 {
		s.bra = s.cursor
		switch among_var {
		case 1:
			// 'atie'  (R1 <-'eer')
			if s.r_R1() {
				s.slice_from("eer")
				return true
			}
		case 2:
			// 'iteit' (R1 delete lengthen_V)
			if s.r_R1() {
				s.slice_del()
				s.r_lengthen_V()
				return true
			}
		case 3:
			// 'heid'
			// 'sel'
			// 'ster'  (R1 delete)
			if s.r_R1() {
				s.slice_del()
				return true
			}
		case 4:
			// 'rder'  (<-'r')
			s.slice_from("r")
			return true
		case 5:
			// 'ing'
			// 'isme'
			// 'erij'  (R1 delete lengthen_V)
			if s.r_R1() {
				s.slice_del()
				s.r_lengthen_V()
				return true
			}
		case 6:
			// 'arij'  (R1 C <-'aar')
			if s.r_R1() && s.r_C() {
				s.slice_from("aar")
				return true
			}
		case 7:
			// 'fie'   (R2 delete attach 'f' lengthen_V)
			if s.r_R2() {
				s.slice_del()
				s.insert(s.cursor, s.cursor, "f")
				s.r_lengthen_V()
				return true
			}
		case 8:
			// 'gie'   (R2 delete attach 'g' lengthen_V)
			if s.r_R2() {
				s.slice_del()
				s.insert(s.cursor, s.cursor, "g")
				s.r_lengthen_V()
				return true
			}
		case 9:
			// 'tst'   (R1 C <-'t')
			if s.r_R1() && s.r_C() {
				s.slice_from("t")
				return true
			}
		case 10:
			// 'dst'   (R1 C <-'d')
			if s.r_R1() && s.r_C() {
				s.slice_from("d")
				return true
			}
		}
	}
	// )
	return false
}

var (
	a_step_4_1 = []*Among{
		NewAmong("end", -1, 10),
		NewAmong("atief", -1, 2),
		NewAmong("erig", -1, 10),
		NewAmong("achtig", -1, 9),
		NewAmong("ioneel", -1, 1),
		NewAmong("baar", -1, 3),
		NewAmong("laar", -1, 5),
		NewAmong("naar", -1, 4),
		NewAmong("raar", -1, 6),
		NewAmong("eriger", -1, 10),
		NewAmong("achtiger", -1, 9),
		NewAmong("lijker", -1, 8),
		NewAmong("tant", -1, 7),
		NewAmong("erigst", -1, 10),
		NewAmong("achtigst", -1, 9),
		NewAmong("lijkst", -1, 8),
	}
	a_step_4_2 = []*Among{
		NewAmong("ig", -1, 1),
		NewAmong("iger", -1, 1),
		NewAmong("igst", -1, 1),
	}
)

// define Step_4 as
func (s *KpStemmer) r_Step_4() bool {
	r := s.limit - s.cursor
	// (   [among ( (])
	{
		s.ket = s.cursor
		among_var := s.find_among_b(a_step_4_1)
		if among_var > 0 {
			s.bra = s.cursor
			switch among_var {
			case 1:
				// 'ioneel'  (R1 <-'ie')
				if s.r_R1() {
					s.slice_from("ie")
					return true
				}
			case 2:
				// 'atief'   (R1 <-'eer')
				if s.r_R1() {
					s.slice_from("eer")
					return true
				}
			case 3:
				// 'baar'    (R1 delete)
				if s.r_R1() {
					s.slice_del()
					return true
				}
			case 4:
				// 'naar'    (R1 V <-'n')
				if s.r_R1() && s.r_V() {
					s.slice_from("n")
					return true
				}
			case 5:
				// 'laar'    (R1 V <-'l')
				if s.r_R1() && s.r_V() {
					s.slice_from("l")
					return true
				}
			case 6:
				// 'raar'    (R1 V <-'r')
				if s.r_R1() && s.r_V() {
					s.slice_from("r")
					return true
				}
			case 7:
				// 'tant'    (R1 <-'teer')
				if s.r_R1() {
					s.slice_from("teer")
					return true
				}
			case 8:
				// 'lijker'
				// 'lijkst'  (R1 <-'lijk')
				if s.r_R1() {
					s.slice_from("lijk")
					return true
				}
			case 9:
				// 'achtig'
				// 'achtiger'
				// 'achtigst'(R1 delete)
				if s.r_R1() {
					s.slice_del()
					return true
				}
			case 10:
				// 'eriger'
				// 'erigst'
				// 'erig'
				// 'end'     (R1 C delete lengthen_V)
				if s.r_R1() && s.r_C() {
					s.slice_del()
					s.r_lengthen_V()
					return true
				}
			}
		}
	}
	// )
	s.cursor = s.limit - r
	// or
	// (   [among ( (])
	{
		s.ket = s.cursor
		among_var := s.find_among_b(a_step_4_2)
		if among_var > 0 {
			s.bra = s.cursor
			switch among_var {
			case 1:
				// 'iger'
				// 'igst'
				// 'ig'      (R1 C delete lengthen_V)
				if s.r_R1() && s.r_C() {
					s.slice_del()
					s.r_lengthen_V()
					return true
				}
			}
		}
	}
	// )
	return false
}

var a_step_6 = []*Among{
	NewAmong("bb", -1, 1),
	NewAmong("cc", -1, 2),
	NewAmong("dd", -1, 3),
	NewAmong("ff", -1, 4),
	NewAmong("gg", -1, 5),
	NewAmong("hh", -1, 6),
	NewAmong("jj", -1, 7),
	NewAmong("kk", -1, 8),
	NewAmong("ll", -1, 9),
	NewAmong("mm", -1, 10),
	NewAmong("nn", -1, 11),
	NewAmong("pp", -1, 12),
	NewAmong("qq", -1, 13),
	NewAmong("rr", -1, 14),
	NewAmong("ss", -1, 15),
	NewAmong("tt", -1, 16),
	NewAmong("v", -1, 21),
	NewAmong("vv", 16, 17),
	NewAmong("ww", -1, 18),
	NewAmong("xx", -1, 19),
	NewAmong("z", -1, 22),
	NewAmong("zz", 20, 20),
}

// define Step_6 as
func (s *KpStemmer) r_Step_6() {
	// [among ( (])
	s.ket = s.cursor
	among_var := s.find_among_b(a_step_6)
	if among_var > 0 {
		s.bra = s.cursor
		switch among_var {
		case 1:
			// 'bb'   (<-'b')
			s.slice_from("b")
		case 2:
			// 'cc'   (<-'c')
			s.slice_from("c")
		case 3:
			// 'dd'   (<-'d')
			s.slice_from("d")
		case 4:
			// 'ff'   (<-'f')
			s.slice_from("f")
		case 5:
			// 'gg'   (<-'g')
			s.slice_from("g")
		case 6:
			// 'hh'   (<-'h')
			s.slice_from("h")
		case 7:
			// 'jj'   (<-'j')
			s.slice_from("j")
		case 8:
			// 'kk'   (<-'k')
			s.slice_from("k")
		case 9:
			// 'll'   (<-'l')
			s.slice_from("l")
		case 10:
			// 'mm'   (<-'m')
			s.slice_from("m")
		case 11:
			// 'nn'   (<-'n')
			s.slice_from("n")
		case 12:
			// 'pp'   (<-'p')
			s.slice_from("p")
		case 13:
			// 'qq'   (<-'q')
			s.slice_from("q")
		case 14:
			// 'rr'   (<-'r')
			s.slice_from("r")
		case 15:
			// 'ss'   (<-'s')
			s.slice_from("s")
		case 16:
			// 'tt'   (<-'t')
			s.slice_from("t")
		case 17:
			// 'vv'   (<-'v')
			s.slice_from("v")
		case 18:
			// 'ww'   (<-'w')
			s.slice_from("w")
		case 19:
			// 'xx'   (<-'x')
			s.slice_from("x")
		case 20:
			// 'zz'   (<-'z')
			s.slice_from("z")
		case 21:
			// 'v'    (<-'f')
			s.slice_from("f")
		case 22:
			// 'z'    (<-'s')
			s.slice_from("s")
		}
	}
	// )
}

var a_step_7 = []*Among{
	NewAmong("ft", -1, 2),
	NewAmong("kt", -1, 1),
	NewAmong("pt", -1, 3),
}

// define Step_7 as
func (s *KpStemmer) r_Step_7() bool {
	// [among ( (])
	s.ket = s.cursor
	among_var := s.find_among_b(a_step_7)
	if among_var > 0 {
		s.bra = s.cursor
		switch among_var {
		case 1:
			// 'kt'   (<-'k')
			s.slice_from("k")
		case 2:
			// 'ft'   (<-'f')
			s.slice_from("f")
		case 3:
			// 'pt'   (<-'p')
			s.slice_from("p")
		}
		return true
	}
	// )
	return false
}

// define V  as test (v or 'ij')
func (s *KpStemmer) r_V() bool {
	r := s.limit - s.cursor
	if s.in_grouping_b(g_v, rune_a, rune_y) || s.eq_s_b("ij") {
		s.cursor = s.limit - r
		return true
	}
	return false
}

// define VX as test (next v or 'ij')
func (s *KpStemmer) r_VX() bool {
	r := s.limit - s.cursor
	if s.cursor > s.limit_backward {
		s.cursor--
		if s.in_grouping_b(g_v, rune_a, rune_y) || s.eq_s_b("ij") {
			s.cursor = s.limit - r
			return true
		}
	}
	return false
}

// define lengthen_V as do (
func (s *KpStemmer) r_lengthen_V() {
	r_1 := s.limit - s.cursor
	// non-v_WX [
	if s.out_grouping_b(g_v_WX, rune_a, rune_y) {
		s.ket = s.cursor
		ok := false
		{
			r_2 := s.limit - s.cursor
			// (AOU] test (non-v or atlimit)) or
			if s.in_grouping_b(g_AOU, rune_a, rune_u) {
				s.bra = s.cursor
				ok = s.out_grouping_b(g_v, rune_a, rune_y) || s.cursor <= s.limit_backward
			}
			s.cursor = s.limit - r_2
		}
		if !ok {
			// ('e'] test (non-v or atlimit
			//             not AIOU
			//             not (next AIOU non-v)))
			if s.eq_s_b("e") {
				s.bra = s.cursor
				r_3 := s.limit - s.cursor
				if (s.out_grouping_b(g_v, rune_a, rune_y) || s.cursor <= s.limit_backward) && !s.in_grouping_b(g_AIOU, rune_a, rune_u) {
					ok = s.cursor <= s.limit_backward
					if !ok {
						s.cursor--
						ok = !s.in_grouping_b(g_AIOU, rune_a, rune_u) || !s.out_grouping_b(g_v, rune_a, rune_y)
					}
				}
				s.cursor = s.limit - r_3
			}
		}
		// ->ch insert ch
		if ok {
			s_ch := s.slice_to("")
			c := s.cursor
			s.insert(s.cursor, s.cursor, string(s_ch))
			s.cursor = c
		}
	}
	s.cursor = s.limit - r_1
}

// define measure as
func (s *KpStemmer) r_measure() {
	c := s.cursor
	// do (
	{
		s.cursor = s.limit // tolimit
		s.i_p1 = s.cursor  // setmark p1
		s.i_p2 = s.cursor  // setmark p2
	}
	// )
	s.cursor = c
	ok := true
	// do(
	// repeat non-v  atleast 1 ('ij' or v)  non-v  setmark p1
	if ok {
		for s.out_grouping(g_v, rune_a, rune_y) {
		}
		t := 0
		for s.eq_s("ij") || s.in_grouping(g_v, rune_a, rune_y) {
			t++
		}
		ok = t >= 1 && s.out_grouping(g_v, rune_a, rune_y)
		if ok {
			s.i_p1 = s.cursor
		}
	}
	// repeat non-v  atleast 1 ('ij' or v)  non-v  setmark p2
	if ok {
		for s.out_grouping(g_v, rune_a, rune_y) {
		}
		t := 0
		for s.eq_s("ij") || s.in_grouping(g_v, rune_a, rune_y) {
			t++
		}
		ok = t >= 1 && s.out_grouping(g_v, rune_a, rune_y)
		if ok {
			s.i_p2 = s.cursor
		}
	}
	// )
	s.cursor = c
}

func (s *KpStemmer) Stem() {
	b_Y_found := false // unset Y_found
	b_stemmed := false // unset stemmed
	// do ( ['y'] <-'Y' set Y_found )
	{
		c := s.cursor
		s.bra = s.cursor
		if s.eq_s("y") {
			s.ket = s.cursor
			s.slice_from("Y")
			b_Y_found = true
		}
		s.cursor = c
	}
	// do repeat(goto (v  ['y'])<-'Y' set Y_found )
	{
		c_1 := s.cursor
		for ; s.cursor <= s.limit; s.cursor++ {
			c_2 := s.cursor
			if s.in_grouping(g_v, rune_a, rune_y) {
				s.bra = s.cursor
				if s.eq_s("y") {
					s.ket = s.cursor
					s.cursor = c_2
					s.slice_from("Y")
					b_Y_found = true
					continue
				}
			}
			s.cursor = c_2
		}
		s.cursor = c_1
	}
	// measure
	s.r_measure()
	// backwards
	{
		s.limit_backward = s.cursor
		s.cursor = s.limit
		// do (Step_1 set stemmed )
		{
			r := s.limit - s.cursor
			if s.r_Step_1() {
				b_stemmed = true
			}
			s.cursor = s.limit - r
		}
		// do (Step_2 set stemmed )
		{
			r := s.limit - s.cursor
			if s.r_Step_2() {
				b_stemmed = true
			}
			s.cursor = s.limit - r
		}
		// do (Step_3 set stemmed )
		{
			r := s.limit - s.cursor
			if s.r_Step_3() {
				b_stemmed = true
			}
			s.cursor = s.limit - r
		}
		// do (Step_4 set stemmed )
		{
			r := s.limit - s.cursor
			if s.r_Step_4() {
				b_stemmed = true
			}
			s.cursor = s.limit - r
		}
		s.cursor = s.limit_backward
	}
	s.b_GE_removed = false // // unset GE_removed
	// do (Lose_prefix and measure)
	{
		c := s.cursor
		if s.r_Lose_prefix() {
			s.cursor = c
			s.r_measure()
		}
		s.cursor = c
	}
	// backwards
	{
		s.limit_backward = s.cursor
		s.cursor = s.limit
		// do (GE_removed Step_1c)
		{
			if s.b_GE_removed {
				s.r_Step_1c()
			}
		}
		s.cursor = s.limit_backward
	}
	s.b_GE_removed = false // // unset GE_removed
	// do (Lose_infix and measure)
	{
		c := s.cursor
		if s.r_Lose_infix() {
			s.cursor = c
			s.r_measure()
		}
		s.cursor = c
	}
	// backwards
	{
		s.limit_backward = s.cursor
		s.cursor = s.limit
		// do (GE_removed Step_1c)
		{
			if s.b_GE_removed {
				s.r_Step_1c()
			}
		}
		s.cursor = s.limit_backward
	}
	// backwards
	{
		s.limit_backward = s.cursor
		s.cursor = s.limit
		// do (Step_7 set stemmed )
		{
			r := s.limit - s.cursor
			if s.r_Step_7() {
				b_stemmed = true
			}
			s.cursor = s.limit - r
		}
		// do (stemmed or GE_removed Step_6)
		{
			r := s.limit - s.cursor
			if b_stemmed || s.b_GE_removed {
				s.r_Step_6()
			}
			s.cursor = s.limit - r
		}
		s.cursor = s.limit_backward
	}
	// do(Y_found  repeat(goto (['Y']) <-'y'))
	{
		if b_Y_found {
			c_1 := s.cursor
			for ; s.cursor <= s.limit; s.cursor++ {
				c_2 := s.cursor
				s.bra = s.cursor
				if s.eq_s("Y") {
					s.ket = s.cursor
					s.cursor = c_2
					s.slice_from("y")
					continue
				}
				s.cursor = c_2
			}
			s.cursor = c_1
		}
	}
}

func Stem(value string) string {
	s := KpStemmer{}
	s.SetCurrent(value)
	s.Stem()
	return s.GetCurrent()
}
