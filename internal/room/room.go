package room

type Room struct {
	table *interfaces.ITable
}

func NewRoom(tableName string) *Room {
	if tableName == "" {
		return nil
	}

}
