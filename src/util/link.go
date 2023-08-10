package util

type Link interface {
	GetTarget() *File
	RemoveTarget()
	IsSymLink() bool
}

// Symlinks point to the path to the target file, not the file itself
type SymLink struct {
	name   string
	target *File
}

// Hard links point to an underlying file (not a directory)
type HardLink struct {
	name   string
	target *File
}

func NewSymLink(name string, target *File) *SymLink {
	return &SymLink{
		name:   name,
		target: target,
	}
}

func NewHardLink(name string, target *File) *HardLink {
	return &HardLink{
		name:   name,
		target: target,
	}
}

func (hl *HardLink) GetTarget() *File {
	return hl.target
}

func (sl *SymLink) GetTarget() *File {
	return sl.target
}

func (hl *HardLink) RemoveTarget() {
	hl.target = nil
}

func (sl *SymLink) RemoveTarget() {
	sl.target = nil
}

func (hl *HardLink) IsSymLink() bool {
	return false
}

func (sl *SymLink) IsSymLink() bool {
	return true
}
