package asyncawait

type Coroutine[A, I, O any] struct {
	f  func(*Coroutine[A, I, O], A)
	cI chan *WrapI[I]
	cO chan *WrapO[O]
}

type WrapI[I any] struct {
	Value I
	Error error
}

type WrapO[O any] struct {
	Value O
	Done  bool
}

func NewCoroutine[A, I, O any](f func(*Coroutine[A, I, O], A), a A) *Coroutine[A, I, O] {
	c := &Coroutine[A, I, O]{
		f:  f,
		cI: make(chan *WrapI[I]),
		cO: make(chan *WrapO[O]),
	}

	go func() {
		<-c.cI
		c.f(c, a)
		close(c.cI)
		c.cO <- &WrapO[O]{Done: true}
		close(c.cO)
	}()

	return c
}

func (c *Coroutine[A, I, O]) Invoke() (O, bool) {
	c.cI <- nil
	o := <-c.cO
	return o.Value, o.Done
}

func (c *Coroutine[A, I, O]) Resume(i I, e error) (O, bool) {
	c.cI <- &WrapI[I]{Value: i, Error: e}
	o := <-c.cO
	return o.Value, o.Done
}

func (c *Coroutine[A, I, O]) Yield(o O) (I, error) {
	c.cO <- &WrapO[O]{Value: o, Done: false}
	i := <-c.cI
	return i.Value, i.Error
}
