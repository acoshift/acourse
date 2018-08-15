package paginate

// Paginate type
type Paginate struct {
	page    int64
	perPage int64
	items   int64
}

// New creates paginate from page, per page, and items
func New(page, perPage, items int64) *Paginate {
	if items < 0 {
		items = 0
	}
	if perPage <= 0 {
		perPage = 1
	}
	if page <= 0 {
		page = 1
	} else if m := maxPage(items, perPage); page > m {
		page = max(m, 1)
	}
	return &Paginate{
		page:    page,
		perPage: perPage,
		items:   items,
	}
}

// FromLimitOffset creates new paginate from limit, offset and count
func FromLimitOffset(limit, offset, count int64) *Paginate {
	if count < 0 {
		count = 0
	}
	if limit <= 0 {
		limit = 1
	}
	if offset < 0 {
		offset = 0
	} else if offset > count {
		offset = count
	}
	return &Paginate{
		page:    offset/limit + 1,
		perPage: limit,
		items:   count,
	}
}

// Page returns page
func (p *Paginate) Page() int64 {
	return p.page
}

// PerPage returns per page
func (p *Paginate) PerPage() int64 {
	return p.perPage
}

// Items returns items
func (p *Paginate) Items() int64 {
	return p.items
}

// Count is the alias for Items
func (p *Paginate) Count() int64 {
	return p.items
}

// Limit returns per page
func (p *Paginate) Limit() int64 {
	return p.perPage
}

// Offset returns offset for current page
func (p *Paginate) Offset() int64 {
	return (p.page - 1) * p.perPage
}

// LimitOffset returns limit and offet
func (p *Paginate) LimitOffset() (limit, offset int64) {
	return p.Limit(), p.Offset()
}

func maxPage(items, perPage int64) int64 {
	m := items % perPage
	if m > 0 {
		m = 1
	}
	return max(items/perPage+m, 1)
}

// MaxPage returns max page
func (p *Paginate) MaxPage() int64 {
	return maxPage(p.items, p.perPage)
}

// CanPrev returns is current page can go prev
func (p *Paginate) CanPrev() bool {
	return p.page > 1
}

// CanNext returns is current page can go next
func (p *Paginate) CanNext() bool {
	return p.page < p.MaxPage()
}

// Prev returns prev page
func (p *Paginate) Prev() int64 {
	return max(p.page-1, 1)
}

// Next returns next page
func (p *Paginate) Next() int64 {
	return min(p.page+1, p.MaxPage())
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

// Pages returns page numbers for paginate
//
// around is the number of page around the current page
// ex. if current page is 10 and around is 3
// the result is 0 7 8 9 10 11 12 13 0
//
// edge is the number of page at the edge
// ex. if current page is 10, max page is 20 and edge is 2
// the result is 1 2 0 10 0 19 20
//
// then if current page is 10, max page is 20,
// around is 3, and edge is 2
// the result is
// 1 2 0 7 8 9 10 11 12 13 0 19 20
func (p *Paginate) Pages(around, edge int64) []int64 {
	xs := make([]int64, 0)
	maxPage := p.MaxPage()

	var current int64 = 1
	var m int64

	if edge == p.page-around-2 {
		m = p.page - around - 1
	} else {
		m = min(edge, p.page-around-1)
	}
	for ; current <= m; current++ {
		xs = append(xs, current)
	}

	if current < p.page-around {
		xs = append(xs, 0)
	}

	current = max(current, p.page-around)

	if p.page+around+1 == maxPage-edge {
		m = p.page + around + 1
	} else {
		m = min(p.page+around, maxPage)
	}
	for ; current <= m; current++ {
		xs = append(xs, current)
	}

	if current < maxPage-edge {
		xs = append(xs, 0)
	}

	current = max(current, maxPage-edge+1)
	for ; current <= maxPage; current++ {
		xs = append(xs, current)
	}
	return xs
}

// MovablePaginate is the paginate for movable list
type MovablePaginate struct {
	page    int64
	perPage int64
	pages   int64
	cnt     int64
}

// NewMovable creates new movable paginate
func NewMovable(page, perPage, pages int64) *MovablePaginate {
	if perPage <= 0 {
		perPage = 1
	}
	if page <= 0 {
		page = 1
	}
	if pages <= 0 {
		pages = 1
	}
	return &MovablePaginate{
		page:    page,
		perPage: perPage,
		pages:   pages,
		cnt:     (last(page, pages) - page + 1) * perPage,
	}
}

// Page returns page
func (p *MovablePaginate) Page() int64 {
	return p.page
}

// PerPage returns per page
func (p *MovablePaginate) PerPage() int64 {
	return p.perPage
}

// Count return count
func (p *MovablePaginate) Count() int64 {
	return p.cnt
}

// SetCount sets count
// where cnt is the count of next page to show
//
// `select count(*) from (select * from table offset {CountOffset} limit {CountLimit}) as t;`
func (p *MovablePaginate) SetCount(cnt int64) *MovablePaginate {
	p.cnt = cnt
	return p
}

// CountLimit returns limit for count
func (p *MovablePaginate) CountLimit() int64 {
	return (last(p.page, p.pages)-p.page+1)*p.perPage + 1
}

// CountOffset returns offset for count
func (p *MovablePaginate) CountOffset() int64 {
	return (p.page - 1) * p.perPage
}

// Counting runs set count from counter function
func (p *MovablePaginate) Counting(counter func(limit, offset int64) int64) {
	p.SetCount(counter(p.CountLimit(), p.CountOffset()))
}

// Limit returns per page
func (p *MovablePaginate) Limit() int64 {
	return p.perPage
}

// Offset returns offset for current page
func (p *MovablePaginate) Offset() int64 {
	return (p.page - 1) * p.perPage
}

// LimitOffset returns limit and offet
func (p *MovablePaginate) LimitOffset() (limit, offset int64) {
	return p.Limit(), p.Offset()
}

// MaxPage returns max page
func (p *MovablePaginate) MaxPage() int64 {
	if p.cnt <= p.perPage {
		return p.page
	}
	maxPage := p.page + p.cnt/p.perPage - 1
	if p.cnt%p.perPage > 0 {
		maxPage++
	}
	return maxPage
}

// CanPrev returns is current page can go prev
func (p *MovablePaginate) CanPrev() bool {
	return p.page > 1
}

// CanNext returns is current page can go next
func (p *MovablePaginate) CanNext() bool {
	return p.page < p.MaxPage()
}

// Prev returns prev page
func (p *MovablePaginate) Prev() int64 {
	return max(p.page-1, 1)
}

// Next returns next page
func (p *MovablePaginate) Next() int64 {
	return min(p.page+1, p.MaxPage())
}

// First returns first page in paginate
func (p *MovablePaginate) First() int64 {
	return max(p.page-p.pages/2, 1)
}

// Last returns last page in paginate without calculate count
func last(page, pages int64) int64 {
	last := page + pages/2
	if pages%2 == 0 {
		last = last - min(page-pages/2, 1)
	} else {
		last = last - min(page-pages/2-1, 0)
	}
	return last
}

// Pages returns page number for paginate without last page
func (p *MovablePaginate) Pages() []int64 {
	xs := make([]int64, 0, p.pages)

	for i := p.First(); i <= p.MaxPage(); i++ {
		xs = append(xs, i)
	}

	return xs
}
