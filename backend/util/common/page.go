package common

import (
	commonpb "github.com/XHXHXHX/medical_marketing_proto/gen/go/proto/common"
	"github.com/pkg/errors"
)

const (
	defaultPageSize = 10
	maxPageSize = 100
)

type Page commonpb.Page

// 获取分页信息, 如果 p 为 nil,返回一个默认的分页信息.
// CurrentPage 从 1 开始.
// PageSize 必须大于 0, 小于 maxPageSize.
func GetPageInfo(p *commonpb.Page) (*Page, error) {
	if p == nil {
		return &Page{
			CurrentPage: 1,
			PageSize:    defaultPageSize,
		}, nil
	}
	if p.CurrentPage == 0 {
		return nil, errors.New("current_page must greater than zero")
	}
	if p.PageSize == 0 || p.PageSize > maxPageSize {
		return nil, errors.New("page_size must in [1, 100]")
	}
	return &Page{
		CurrentPage: p.CurrentPage,
		PageSize:    p.PageSize,
	}, nil
}


func (p *Page) AsRaw() *commonpb.Page {
	if p == nil {
		return nil
	}
	return &commonpb.Page{
		CurrentPage: p.CurrentPage,
		PageSize:    p.PageSize,
	}
}

// 获取分页起始位置.
func (p *Page) ToFrom() int {
	if p == nil {
		return 0
	}
	return int(p.PageSize * (p.CurrentPage - 1))
}

// 获取每页容量.
func (p *Page) ToLimit() int {
	if p == nil {
		return 0
	}
	return int(p.PageSize)
}

// func Page2Pagination(page *sharing.Page, total int64) *sharing.Pagination {
// panic if page is nil. page 必须是合法的, 必须保证 pageSize 不为 0.
func (p *Page) Page2Pagination(total int64) *commonpb.Pagination {
	pageCount := total / p.PageSize
	if total%p.PageSize > 0 {
		pageCount += 1
	}
	return &commonpb.Pagination{
		Total:       int32(total),
		CurrentPage: int32(p.CurrentPage),
		PageSize:    int32(p.PageSize),
		PageCount:   int32(pageCount),
	}
}
