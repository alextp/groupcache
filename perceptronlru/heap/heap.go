package heap

type HeapItem struct {
	priority float64
	Position int
	Value    interface{}
}

type ItemArray []*HeapItem

type Heap struct {
	Size     int
	elements ItemArray
}

func NewHeap() *Heap {
	return &Heap{
		Size:     0,
		elements: ItemArray{},
	}
}

func (heap *Heap) Swap(i, j int) {
	a := heap.elements[i]
	heap.elements[i] = heap.elements[j]
	heap.elements[j] = a
	heap.elements[i].Position = i
	heap.elements[j].Position = j
}

func (heap *Heap) Up(index int) {
	for {
		i := (index - 1) / 2 // parent
		if i == index || !(heap.elements[index].priority < heap.elements[i].priority) {
			break
		}
		heap.Swap(i, index)
		index = i
	}
}

func (heap *Heap) Down(index int) {
	for {
		j1 := 2*index + 1
		if j1 >= heap.Size || j1 < 0 {
			break
		}
		j := j1
		p1 := heap.elements[j1].priority
		j2 := j1 + 1
		if j2 < heap.Size {
			p2 := heap.elements[j2].priority
			if !(p1 < p2) {
				j = j2
			}
		}
		if !(heap.elements[j].priority < heap.elements[index].priority) {
			break
		}
		heap.Swap(index, j)
		index = j
	}
}

func (heap *Heap) Push(element *HeapItem) {
	element.Position = heap.Size
	heap.Size += 1
	heap.elements = append(heap.elements, element)
	heap.Up(element.Position)
}

func (heap *Heap) Insert(element interface{}, priority float64) *HeapItem {
	item := &HeapItem{
		Value:    element,
		priority: priority,
		Position: -1,
	}
	heap.Push(item)
	return item
}

func (heap *Heap) Remove(index int) *HeapItem {
	n := heap.Size - 1
	heap.Swap(index, n)
	heap.Down(index)
	e := heap.elements[n]
	heap.Size -= 1
	heap.elements = heap.elements[0:heap.Size]
	return e
}

func (heap *Heap) Reinsert(index int, priority float64) {
	item := heap.Remove(index)
	item.priority = priority
	heap.Push(item)
}

func (heap *Heap) Pop() *HeapItem {
	return heap.Remove(0)
}
