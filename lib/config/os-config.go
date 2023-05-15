package config

type operatingSystem struct {
	name                string
	code                string
	version             string
	kernel              string
	shell               string
	init                string
	containerized       bool
	containerManagement string
}

func (os *operatingSystem) SetName(name string) {
	os.name = name
}

func (os *operatingSystem) GetName() string {
	return os.name
}

func (os *operatingSystem) SetCode(code string) {
	os.code = code
}

func (os *operatingSystem) GetCode() string {
	return os.code
}

func (os *operatingSystem) SetVersion(version string) {
	os.version = version
}

func (os *operatingSystem) GetVersion() string {
	return os.version
}

func (os *operatingSystem) SetKernel(kernel string) {
	os.kernel = kernel
}

func (os *operatingSystem) GetKernel() string {
	return os.kernel
}

func (os *operatingSystem) SetShell(shell string) {
	os.shell = shell
}

func (os *operatingSystem) GetShell() string {
	return os.shell
}

func (os *operatingSystem) SetInit(init string) {
	os.init = init
}

func (os *operatingSystem) GetInit() string {
	return os.init
}
