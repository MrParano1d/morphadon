package core

type DefaultSystem[C Context] struct {
	ctx C
}

var _ System[*TodoContext] = (*DefaultSystem[*TodoContext])(nil)

func NewDefaultSystem[C Context]() *DefaultSystem[C] {
	return &DefaultSystem[C]{}
}

func (s *DefaultSystem[C]) Init(app App[C]) error {
	return nil
}

func (s *DefaultSystem[C]) Context() C {
	return s.ctx
}

func (s *DefaultSystem[C]) SetContext(ctx C) {
	s.ctx = ctx
}

func (s *DefaultSystem[C]) Actions() []Action[C] {
	return make([]Action[C], 0)
}
