package common

type Page struct {
	CurrentPage int
	PageSize int
}

type Pagination struct {
	Page
	PageCount int
	Total int
}

func (p *Page) ToPagination(total int) *Pagination {
	pageCount := total / p.PageSize
	if total%p.PageSize > 0 {
		pageCount += 1
	}
	return &Pagination{
		Page:        *p,
		PageCount: pageCount,
		Total:       total,
	}
}
