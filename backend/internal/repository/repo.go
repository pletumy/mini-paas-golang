package repository

type Page struct {
	Limit  int
	Offset int
}

type Sort struct {
	Field string
	Desc  bool
}

type ListStruct[T any] struct {
	Items []T
	Total int64
}

func (p Page) Sanitize(maxLimit int) Page {
	if p.Limit <= 0 || p.Limit > maxLimit {
		p.Limit = maxLimit
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
	return p
}

type ListResult[T any] struct {
	Items []T
	Total int64
}
