package game

type field uint8

const (
	fEmpty field = iota
	fUnknown
	fHit
	fMissed
	fOccupied
)

type sea struct {
	width  int
	height int

	data []field
}

func newSea(width, height int) *sea {
	size := width * height
	data := make([]field, size)
	for i := 0; i < size; i++ {
		data[i] = fEmpty
	}
	return &sea{
		width:  width,
		height: height,
		data:   data,
	}
}

func (sea *sea) toIndices(index, shipWidth, shipHeight int) []int {
	size := shipWidth * shipHeight
	indices := make([]int, size)
	for i := 0; i < size; i++ {
		row := i / shipWidth
		col := i % shipWidth
		indices[i] = index + col + row*sea.width
	}
	return indices
}

func (sea *sea) canPut(index, shipWidth, shipHeight int) bool {
	if index/sea.width+shipHeight > sea.height {
		return false
	}
	if index%sea.width+shipWidth > sea.width {
		return false
	}
	indices := sea.toIndices(index, shipWidth, shipHeight)
	for _, i := range indices {
		if sea.data[i] != fEmpty {
			return false
		}
	}
	return true
}

func (sea *sea) put(index, shipWidth, shipHeight int) {
	indices := sea.toIndices(index, shipWidth, shipHeight)
	for _, i := range indices {
		sea.data[i] = fOccupied
	}
}

func (sea *sea) canAttack(index int) bool {
	if index < 0 || index >= len(sea.data) {
		return false
	}
	return sea.data[index] != fHit && sea.data[index] != fMissed
}

func (sea *sea) attack(index int) bool {
	if sea.data[index] == fOccupied {
		sea.data[index] = fHit
		return true
	} else {
		sea.data[index] = fMissed
		return false
	}
}

func (sea *sea) setData(index int, updatedType field) {
	sea.data[index] = updatedType
}
