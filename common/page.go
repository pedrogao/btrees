package common

type (
	// Page of disk
	Page struct {
		Id    uint32
		Flags uint32
		Data  uintptr
	}

	PageProvider interface {
		Fetch() *Page
		Write(page *Page)
		Flush() *Page
		Delete() *Page
	}
)
