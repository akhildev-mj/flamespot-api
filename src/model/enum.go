package model

type Type string
type Status string
type Sort string

const (
	TypeAll      Type = "ALL"
	TypeTakeaway Type = "TAKEAWAY"
	TypeDineIn   Type = "DINE_IN"
)

const (
	StatusAll     Status = "ALL"
	StatusOrdered Status = "ORDERED"
	StatusBilled  Status = "BILLED"
	StatusSaved   Status = "SAVED"
	StatusDeleted Status = "DELETED"
)

const (
	SortRecent Sort = "RECENT"
	SortName   Sort = "NAME"
	SortPrice  Sort = "PRICE"
)

func (e Type) IsValid() bool {
	switch e {
	case TypeTakeaway, TypeDineIn:
		return true
	}
	return false
}

func (e Type) IsValidFilter() bool {
	if e == TypeAll {
		return true
	}
	return e.IsValid()
}

func (e Status) IsValid() bool {
	switch e {
	case StatusOrdered, StatusBilled, StatusSaved, StatusDeleted:
		return true
	}
	return false
}

func (e Status) IsValidFilter() bool {
	if e == StatusAll {
		return true
	}
	return e.IsValid()
}

func (s Sort) IsValidForCategory() bool {
	switch s {
	case SortRecent, SortName:
		return true
	}
	return false
}

func (s Sort) IsValidForMenu() bool {
	switch s {
	case SortRecent, SortName, SortPrice:
		return true
	}
	return false
}

func (s Sort) IsValidForOrder() bool {
	switch s {
	case SortRecent, SortPrice:
		return true
	}
	return false
}

func FilterOrderTypes() []string {
	return []string{string(TypeAll), string(TypeTakeaway), string(TypeDineIn)}
}

func FilterOrderStatuses() []string {
	return []string{string(StatusAll), string(StatusOrdered), string(StatusBilled), string(StatusSaved), string(StatusDeleted)}
}

func CategorySorts() []string {
	return []string{string(SortRecent), string(SortName)}
}

func MenuSorts() []string {
	return []string{string(SortRecent), string(SortName), string(SortPrice)}
}

func OrderSorts() []string {
	return []string{string(SortRecent), string(SortPrice)}
}
