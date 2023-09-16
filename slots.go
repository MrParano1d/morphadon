package morphadon

type Slots map[string]any

func (slots Slots) Has(slotName string) bool {
	_, ok := slots[slotName]
	return ok
}

