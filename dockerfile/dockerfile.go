package dockerfile

type Stage struct {
	Args         map[string]struct{}
	Dependencies map[string]struct{}
	Name         string
}

type Dockerfile struct {
	Stages   []Stage
	MetaArgs map[string]string
}
