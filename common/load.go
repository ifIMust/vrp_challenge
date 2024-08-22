package common

type Load struct {
	Index   int
	Pickup  *Location
	Dropoff *Location
}

// Compute distance from Pickup to Dropoff
func (l *Load) Cost() float64 {
	return l.Pickup.Distance(l.Dropoff)
}

// Compute distance from Dropoff to the depot/origin
func (l *Load) HomeCost() float64 {
	return l.Dropoff.Distance(HomeLocation)
}

// LoadMap is a helper type for organizing Loads
type LoadMap map[int]*Load

func (lm LoadMap) Duplicate() LoadMap {
	dup := make(map[int]*Load)
	for k, v := range lm {
		dup[k] = v
	}
	return dup
}

// AsMap is a helper function for dealing with collections of Loads
func AsMap(loads []*Load) LoadMap {
	m := make(map[int]*Load)
	for _, l := range loads {
		m[l.Index] = l
	}
	return m
}

// LoadSorter for sorting a collection of loads based on proximity to a reference location.
// It supports the use of the Go standard library sort.Sort operation.
// This is meant as a single use sorter for ever-changing reference locations.
// A minheap would be preferred for repeated pop operations.
type LoadSorterEntry struct {
	Load     *Load
	Distance float64
}

type LoadSorter struct {
	LoadEntries []*LoadSorterEntry
	Reference   *Location
}

func NewLoadSorter(loads LoadMap, reference *Location) *LoadSorter {
	l := LoadSorter{}
	l.LoadEntries = make([]*LoadSorterEntry, 0)
	l.Reference = reference

	for _, load := range loads {
		l.AddEntry(load)
	}
	return &l
}

func (l *LoadSorter) AddEntry(load *Load) {
	e := LoadSorterEntry{load, l.Reference.Distance(load.Pickup)}
	l.LoadEntries = append(l.LoadEntries, &e)
}

func (l *LoadSorter) Len() int {
	return len(l.LoadEntries)
}

func (l *LoadSorter) Swap(i, j int) {
	l.LoadEntries[i], l.LoadEntries[j] = l.LoadEntries[j], l.LoadEntries[i]
}

func (l *LoadSorter) Less(i, j int) bool {
	return l.LoadEntries[i].Distance < l.LoadEntries[j].Distance
}

// Pop returns the Load that is nearest to the reference
// This is only true if sort.Sort has been called on the LoadSorter
func (l *LoadSorter) Pop() *Load {
	if len(l.LoadEntries) == 0 {
		return nil
	}
	loadEntry := l.LoadEntries[0]
	l.LoadEntries = l.LoadEntries[1:]
	return loadEntry.Load
}
