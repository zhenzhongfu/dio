package main

type A struct {
	a int
}

func (aa *A) get() int {
	return aa.a
}

func (aa *A) get2() int {
	return aa.a
}

func main() {
	aa := A{1}
	for i := 1; i < 100; i++ {
		tmp := i
		go func() {
			if tmp < 30 {
				aa.get()
			} else {
				aa.get2()
			}
		}()
	}
}
