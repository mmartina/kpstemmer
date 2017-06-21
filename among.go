package kpstemmer

type Among struct {
	s_size      int
	s           string
	substring_i int
	result      int
}

func NewAmong(s string, substring_i int, result int) (rcvr *Among) {
	rcvr = &Among{}
	rcvr.s_size = len(s)
	rcvr.s = s
	rcvr.substring_i = substring_i
	rcvr.result = result
	return
}
