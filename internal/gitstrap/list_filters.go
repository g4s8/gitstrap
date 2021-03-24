package gitstrap

// ListFilter for list results
type ListFilter interface {
	check(*listResult) bool
}

// LfNop - list filter does nothing
var LfNop ListFilter = &lfNop{}

type lfNop struct{}

func (f *lfNop) check(r *listResult) bool {
	return true
}

// LfForks - list filter by fork criteria
func LfForks(origin ListFilter, fork bool) ListFilter {
	return &lfFork{origin, fork}
}

type lfFork struct {
	origin ListFilter
	fork   bool
}

func (f *lfFork) check(r *listResult) bool {
	return f.origin.check(r) && r.fork == f.fork
}

// LfStarsCriteria - criteria of repository stars for filtering
type LfStarsCriteria func(int) bool

// LfStarsGt - list filter stars criteria: greater than `val`
func LfStarsGt(val int) LfStarsCriteria {
	return func(x int) bool {
		return x > val
	}
}

// LfStarsLt - list filter stars criteria: less than `val`
func LfStarsLt(val int) LfStarsCriteria {
	return func(x int) bool {
		return x < val
	}
}

// LfStars - list filter by stars count
func LfStars(origin ListFilter, criteria LfStarsCriteria) ListFilter {
	return &lfStars{origin, criteria}
}

type lfStars struct {
	origin   ListFilter
	criteria LfStarsCriteria
}

func (f *lfStars) check(i *listResult) bool {
	return f.origin.check(i) && f.criteria(i.stars)
}
