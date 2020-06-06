module github.com/stiletto/damb

go 1.14

require (
	github.com/Microsoft/go-winio v0.4.15-0.20200113171025-3fe6c5262873 // indirect
	github.com/OpenPeeDeeP/depguard v1.0.1
	github.com/Sirupsen/logrus v1.6.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/logrusorgru/aurora v0.0.0-20200102142835-e9ef32dff381
	github.com/moby/buildkit v0.6.4
	github.com/openshift/imagebuilder v1.1.4
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.4.1
	github.com/stretchr/testify v1.6.1 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c
)

replace github.com/Sirupsen/logrus v1.6.0 => github.com/sirupsen/logrus v1.6.0

replace github.com/containerd/containerd v1.3.0-0.20190507210959-7c1e88399ec0 => github.com/containerd/containerd v1.3.0-beta.2.0.20190823190603-4a2f61c4f2b4

replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
