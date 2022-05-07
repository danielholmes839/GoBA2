package realtime

type ID string

func (id ID) ID() string {
	return string(id)
}