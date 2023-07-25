package schema

type Queue struct {
	Status string `bson:"status,omitempty" json:"status,omitempty" redis:"status"`
	Player string `bson:"player,omitempty" json:"player,omitempty" redis:"player"`
}

func (queue Queue) Database() string {
	return ""
}

func (queue Queue) Collection() string {
	return ""
}

func (queue Queue) Key() string {
	return queue.Status
}
