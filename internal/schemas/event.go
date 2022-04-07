package schemas

type Event struct {
	Service string
	Action  string
	Data    []byte
}
