package domain

type PathOptions struct {
	Description string
	ID          string
	Tags        []string
	Deprecated  bool
}

func (p *PathOptions) Defaults() {

}
