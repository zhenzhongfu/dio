package pkg2

//
type SA struct {
	num int
}

func (sa *SA) Get() int {
	return sa.num
}

func (sa *SA) Set(n int) {
	sa.num = n
}

type ISA interface {
	Get() int
	Set(int)
}
