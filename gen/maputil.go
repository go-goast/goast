package gen

type K interface{}
type V interface{}
type Map map[K]V

func (mp Map) Copy() Map {
	var newMap Map = make(map[K]V)
	for k, v := range mp {
		newMap[k] = v
	}
	return newMap
}

func (mp *Map) Init() {
	var newMap Map = make(map[K]V)
	mp = &newMap
}
