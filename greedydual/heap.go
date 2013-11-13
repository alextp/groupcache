package greedydual

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
	a := heap.elements[i]
	heap.elements[i] = heap.elements[j]
	heap.elements[j] = a
	heap.elements[i].Position = i
	heap.elements[j].Position = j
}

func less(i, j *HeapItem) bool {
	pi := i.Priority
	pj := j.Priority
	return pi <= pj
}

func (heap *Heap) Up(index int) {
	for {
		parent := (index - 1) / 2 // parent
		if parent == index || less(heap.elements[parent], heap.elements[index]) {
			break
		}
		heap.Swap(parent, index)
		index = parent
	}
}

func (heap *Heap) Down(index int) {
	for {
		left := 2*index + 1
		if left >= heap.Size || left < 0 {
			break
		}
		child := left
		right := left + 1
		if right < heap.Size {
			if !less(heap.elements[left], heap.elements[right]) {
				child = right
			}
		}
		if !less(heap.elements[child], heap.elements[index]) {
			break
		}
		heap.Swap(index, child)
		index = child
	}
}

func (heap *Heap) Push(element *HeapItem) {
	element.Position = heap.Size
	heap.Size += 1
	heap.elements = append(heap.elements, element)
	heap.Up(element.Position)
}

func (heap *Heap) Insert(element interface{}, Priority float64) *HeapItem {
	item := &HeapItem{
		Value:    element,
		Priority: Priority,
		Position: -1,
	}
	heap.Push(item)
	heap.Check("insert")
	return item
}

func (heap *Heap) Remove(index int) *HeapItem {
	n := heap.Size - 1
	heap.Swap(index, n)
	e := heap.elements[n]
	heap.Size -= 1
	heap.elements = heap.elements[0:heap.Size]
	heap.Down(index)
	heap.Check("remove")
	return e
}

func (heap *Heap) Reinsert(index int, Priority float64) {
	heap.Check("pre-reinsert")
	e := heap.elements[index]
	for {
		if e.Position == 0 {
			break
		}
		heap.Swap(e.Position, (e.Position-1)/2)
	}
	heap.Remove(0)
	e.Priority = Priority
	heap.Push(e)
	heap.Check("reinsert")
}

func (heap *Heap) Pop() *HeapItem {
	return heap.Remove(0)
}

func (heap *Heap) Head() *HeapItem {
	return heap.elements[0]
}

func (heap *Heap) Check(name string) {
	if 1 == 2 {
		for i := 0; i < heap.Size; i++ {
			parent := (i-1)/2
			if heap.elements[parent].Priority > heap.elements[i].Priority {
				println("error", name, parent, i, heap.elements[parent].Priority, heap.elements[i].Priority)
				println(heap.elements[23423423423])
			}
		}
	}
}