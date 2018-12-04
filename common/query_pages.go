package common

import (
	"fmt"
	"github.com/jucardi/go-db"
	"github.com/jucardi/go-db/pages"
)

// Page adds to the query the information required to fetch the requested page of objects.
func (q *AbstractQuery) Page(page ...*pages.Page) dbx.IQuery {
	return PageHandler(q.Q, page...)
}

// WrapPage attempts to obtain the items in the requested page and wraps the result in *pages.Paginated
func (q *AbstractQuery) WrapPage(result interface{}, page ...*pages.Page) (*pages.Paginated, error) {
	return WrapPageHandler(q.Q, result, page...)
}

func PageHandler(q dbx.IQuery, page ...*pages.Page) dbx.IQuery {
	if len(page) < 1 || page[0] == nil {
		return q
	}

	p := page[0]

	if len(p.Sort) > 0 {
		q.Sort(p.Sort...)
	}
	return q.Skip((p.Page - 1) * p.Size).Limit(p.Size)
}

func WrapPageHandler(q dbx.IQuery, result interface{}, page ...*pages.Page) (*pages.Paginated, error) {
	if len(page) < 1 || page[0] == nil {
		return wrap(q, result, nil)
	}

	n, err := q.Count()

	if err != nil {
		return nil, fmt.Errorf("unable to obtain a count of elements, %v", err)
	}

	p := page[0]
	q.Page(page...)

	return wrap(q, result, p, n)
}

func wrap(q dbx.IQuery, result interface{}, p *pages.Page, n ...int) (*pages.Paginated, error) {
	if err := q.All(result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal page, %v", err)
	}

	return pages.CreatePaginated(p, result, n...)
}
