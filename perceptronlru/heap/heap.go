package heap

type HeapItem struct {
	Priority float64
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
	if i >= heap.Size {
		println("error: i", i, "size", heap.Size, "j", j)
	}
	a := heap.elements[i]
	heap.elements[i] = heap.elements[j]
	heap.elements[j] = a
	heap.elements[i].Position = i
	heap.elements[j].Position = j
}

func (heap *Heap) Up(index int) {
	for {
		i := (index - 1) / 2 // parent
		if i == index || !(heap.elements[i].Priority < heap.elements[index].Priority) {
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
		p1 := heap.elements[j1].Priority
		j2 := j1 + 1
		if j2 < heap.Size {
			p2 := heap.elements[j2].Priority
			if !(p1 < p2) {
				j = j2
			}
		}
		if !(heap.elements[index].Priority < heap.elements[j].Priority) {
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
	if element.Position >= heap.Size {
		println("messing up the declaration of position")
	}
	if heap.elements[element.Position] != element {
		println("Not inserting where I want to")
	}
	heap.Up(element.Position)
}

func (heap *Heap) Insert(element interface{}, Priority float64) *HeapItem {
	item := &HeapItem{
		Value:    element,
		Priority: Priority,
		Position: -1,
	}
	heap.Push(item)
	return item
}

func (heap *Heap) Remove(index int) *HeapItem {
	e := heap.elements[index]
	n := heap.Size - 1
	heap.Swap(index, n)
	heap.elements = heap.elements[0:heap.Size]
	heap.Up(index)
	heap.Down(index)
	return e
}

func (heap *Heap) Reinsert(index int, priority float64) {
	if index < heap.Size {
		if heap.elements[index].Position != index {
			println("we have an element at the wrong position")
		}
	}
	item := heap.Remove(index)
	item.Priority = priority
	heap.Push(item)
}


func (heap *Heap) Head() *HeapItem {
	return heap.elements[0]
}

func (heap *Heap) Check(name string) {
	for i := 0; i < heap.Size; i++ {
		if heap.elements[i].Position != i {
			println("element at wrong position", i, heap.elements[i].Position, name)
		}
	}
}

func (heap *Heap) Index(i int) *HeapItem {
	return heap.elements[i]
}