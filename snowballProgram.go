package kpstemmer

import "bytes"

type snowball struct {
	current        string
	cursor         int
	limit          int
	limit_backward int
	bra            int
	ket            int
}

type bitmask []byte

func (s *snowball) eq_s(str string) bool {
	s_size := len(str)
	if s.limit-s.cursor < s_size {
		return false
	}
	buf := s.current[s.cursor:]
	for i := 0; i != s_size; i++ {
		if buf[i] != str[i] {
			return false
		}
	}
	s.cursor += s_size
	return true
}

func (s *snowball) eq_s_b(str string) bool {
	s_size := len(str)
	if s.cursor-s.limit_backward < s_size {
		return false
	}
	buf := s.current[s.cursor-s_size:]
	for i := 0; i != s_size; i++ {
		if buf[i] != str[i] {
			return false
		}
	}
	s.cursor -= s_size
	return true
}

type among struct {
	s_size      int
	s           string
	substring_i int
	result      int
}

func newAmong(s string, substring_i int, result int) (rcvr *among) {
	rcvr = &among{}
	rcvr.s_size = len(s)
	rcvr.s = s
	rcvr.substring_i = substring_i
	rcvr.result = result
	return
}

func (s *snowball) find_among_b(v []*among) int {
	i := 0
	j := len(v)
	c := s.cursor
	lb := s.limit_backward
	common_i := 0
	common_j := 0
	first_key_inspected := false
	for {
		k := i + (j-i)>>1
		diff := 0
		var common int
		if common_i < common_j {
			common = common_i
		} else {
			common = common_j
		}
		w := v[k]
		buf := s.current
		for i2 := w.s_size - 1 - common; i2 >= 0; i2-- {
			if c-common == lb {
				diff = -1
				break
			}
			diff = int(buf[c-1-common]) - int(w.s[i2])
			if diff != 0 {
				break
			}
			common++
		}
		if diff < 0 {
			j = k
			common_j = common
		} else {
			i = k
			common_i = common
		}
		if j-i <= 1 {
			if i > 0 || j == i || first_key_inspected {
				break
			}
			first_key_inspected = true
		}
	}
	for {
		w := v[i]
		if common_i >= w.s_size {
			s.cursor = c - w.s_size
			return w.result
		}
		i = w.substring_i
		if i < 0 {
			return 0
		}
	}
	return -1
}

func (s *snowball) GetCurrent() string {
	result := string(s.current)
	s.current = ""
	return result
}

func (s *snowball) in_grouping(bitmask bitmask, min rune, max rune) bool {
	if s.cursor >= s.limit {
		return false
	}
	ch := rune(s.current[s.cursor])
	if ch > max || ch < min {
		return false
	}
	b := byte(ch - min)
	if bitmask[b>>3]&(0X1<<(b&0X7)) == 0 {
		return false
	}
	s.cursor++
	return true
}

func (s *snowball) in_grouping_b(bitmask bitmask, min rune, max rune) bool {
	if s.cursor <= s.limit_backward {
		return false
	}
	ch := rune(s.current[s.cursor-1])
	if ch > max || ch < min {
		return false
	}
	b := byte(ch - min)
	if bitmask[b>>3]&(0X1<<(b&0X7)) == 0 {
		return false
	}
	s.cursor--
	return true
}

func (s *snowball) insert(c_bra int, c_ket int, str string) {
	adjustment := s.replace_s(c_bra, c_ket, str)
	if c_bra <= s.bra {
		s.bra += adjustment
	}
	if c_bra <= s.ket {
		s.ket += adjustment
	}
}

func (s *snowball) out_grouping(bitmask bitmask, min rune, max rune) bool {
	if s.cursor >= s.limit {
		return false
	}
	ch := rune(s.current[s.cursor])
	if ch > max || ch < min {
		s.cursor++
		return true
	}
	b := byte(ch - min)
	if bitmask[b>>3]&(0X1<<(b&0X7)) == 0 {
		s.cursor++
		return true
	}
	return false
}

func (s *snowball) out_grouping_b(bitmask bitmask, min rune, max rune) bool {
	if s.cursor <= s.limit_backward {
		return false
	}
	ch := rune(s.current[s.cursor-1])
	if ch > max || ch < min {
		s.cursor--
		return true
	}
	b := byte(ch - min)
	if bitmask[b>>3]&(0X1<<(b&0X7)) == 0 {
		s.cursor--
		return true
	}
	return false
}

func (s *snowball) replace_s(c_bra int, c_ket int, str string) int {
	adjustment := len(str) - (c_ket - c_bra)
	s.current = replaceString(s.current, c_bra, c_ket, str)
	s.limit += adjustment
	if s.cursor >= c_ket {
		s.cursor += adjustment
	} else if s.cursor > c_bra {
		s.cursor = c_bra
	}
	return adjustment
}

func (s *snowball) SetCurrent(value string) {
	s.current = replaceString(s.current, 0, len(s.current), value)
	s.cursor = 0
	s.limit = len(s.current)
	s.limit_backward = 0
	s.bra = s.cursor
	s.ket = s.limit
}

func (s *snowball) slice_check() {
	if s.bra < 0 || s.bra > s.ket || s.ket > s.limit || s.limit > len(s.current) {
		// should not occur
	}
}

func (s *snowball) slice_del() {
	s.slice_from("")
}

func (s *snowball) slice_from(str string) {
	s.slice_check()
	s.replace_s(s.bra, s.ket, str)
}

func (s *snowball) slice_to(buf string) string {
	s.slice_check()
	return replaceString(buf, 0, len(buf), s.current[s.bra:s.ket])
}

func replaceString(src string, start int, end int, sub string) string {
	var buf bytes.Buffer
	buf.Grow(len(src) + len(sub))
	buf.WriteString(src[:start])
	buf.WriteString(sub)
	buf.WriteString(src[end:])
	return buf.String()
}
