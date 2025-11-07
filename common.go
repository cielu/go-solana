package solana

// SlotRange the first slot --> lastSlot
type SlotRange struct {
	FirstSlot uint64 `json:"firstSlot,omitempty"`
	LastSlot  uint64 `json:"lastSlot,omitempty"`
}

// DataSlice Request a slice of the data range.
type DataSlice struct {
	Length uint64 `json:"length"`
	Offset uint64 `json:"offset"`
}
