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
	//heap.Check("swap")
}

func (heap *Heap) Up(index int) {
	//heap.Check("preup")
	for {
		i := (index - 1) / 2 // parent
		if i == index || (heap.elements[i].Priority < heap.elements[index].Priority) {
			break
		}
		heap.Swap(i, index)
		index = i
	}
	heap.Check("up")
}

func (heap *Heap) Down(index int) {
	for {
		left := 2*index + 1
		if left >= heap.Size || left < 0 {
			break
		}
		j := left
		p1 := heap.elements[left].Priority
		right := left + 1
		if right < heap.Size {
			p2 := heap.elements[right].Priority
			if !(p1 < p2) {
				j = right
			}
		}
		if (heap.elements[index].Priority < heap.elements[j].Priority) {
			break
		}
		heap.Swap(index, j)
		index = j
	}
	heap.Check("down")
}

func (heap *Heap) Push(element *HeapItem) {
	heap.Check("prepush")
	element.Position = heap.Size
	heap.Size += 1
	heap.elements = append(heap.elements, nil)
	heap.elements[element.Position] = element
	heap.Check("push1")
	heap.Up(element.Position)
	heap.Check("push")
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
	if index != 0 {
		println(index)
		println(heap.elements[22432423423])
	}
	e := heap.elements[index]
	n := heap.Size - 1
	heap.Swap(index, n)
	heap.elements = heap.elements[0:heap.Size]
	heap.Size -= 1
	heap.Down(index)
	heap.Check("remove")
	return e
}

func (heap *Heap) Reinsert(index int, priority float64) {
	if index < heap.Size {
		if heap.elements[index].Position != index {
			println("we have an element at the wrong position")
			println(heap.elements[2323423423])
		}
	}
	e := heap.elements[index]
	for {
		if e.Position == 0 {
			break
		}
		heap.Swap(e.Position, (e.Position-1)/2)
	}
	//heap.Check("pre remove")
	heap.Remove(0)
	//heap.Check("post remove")
	for i := 0; i < heap.Size; i++ {
		if heap.elements[i] == e {
			println("this shouldn't happen, we're not removing", i, e.Position, heap.Size)
			println(heap.elements[234234324])
		}
	}
	e.Priority = priority
	heap.Push(e)
	heap.Check("reinsert")
}


func (heap *Heap) Head() *HeapItem {
	if heap.elements[0].Position != 0 {
		println("wrong head position", heap.elements[0].Position)
		println(heap.elements[3424324234234])
	}
	return heap.elements[0]
}

func (heap *Heap) Check(name string) {
	if 1 == 1 {
		return
	}
	for i := 0; i < heap.Size; i++ {
		if heap.elements[i].Position != i {
			println("element at wrong position", i, heap.elements[i].Position, name)
			println(heap.elements[34234234234])
		}
		left := 2*i+1
		if left < heap.Size {
			if heap.elements[left].Priority < heap.elements[i].Priority {
				println("violating heap property after", name, i, "left", left, heap.Size)
				println(heap.elements[left*1000])
			}
			right := left + 1
			if right < heap.Size {
				if heap.elements[right].Priority < heap.elements[i].Priority {
					println("violating heap property after", name, i, "right", right)
					println(heap.elements[left*1000])
				}
			}
		}
	}
}

func (heap *Heap) Index(i int) *HeapItem {
	return heap.elements[i]
}