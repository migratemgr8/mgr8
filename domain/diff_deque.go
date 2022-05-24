package domain

type DiffDeque []Diff

func NewDiffDeque() DiffDeque{
	return []Diff{}
}

func (d *DiffDeque) Add(diff Diff) {
	*d = append(*d, diff)
}

func (d *DiffDeque) Extend(diffDeque DiffDeque) {
	*d = append(*d, diffDeque...)
}

func (d DiffDeque) GetUpStatements(deparser Deparser) []string {
	var statements []string
	for _, diff := range d {
		statements = append(statements, diff.Up(deparser))
	}
	return statements
}

func (d DiffDeque) GetDownStatements(deparser Deparser) []string {
	var statements []string
	for _, diff := range d {
		statements = append(statements, diff.Down(deparser))
	}
	return statements
}
