package api

type Pagination struct {
	Page int `query:"page" json:"page"`
	Size int `query:"size" json:"size"`
}

func (p *Pagination) FormatPageAndSize(minSize, maxSize, defaultSize int) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size == 0 {
		p.Size = defaultSize
	}
	if p.Size < minSize {
		p.Size = maxSize
	}
	if p.Size > maxSize {
		p.Size = maxSize
	}
}

func (p *Pagination) CalcOffset() (int, int) {
	offset := (p.Page - 1) * p.Size
	if offset < 0 {
		offset = 0
	}
	return offset, p.Size
}
