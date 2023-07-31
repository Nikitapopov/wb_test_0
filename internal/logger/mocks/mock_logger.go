package mocks

type MockLogger interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
}

type MockLogging struct{}

func NewMockLogger() MockLogger {
	return MockLogging{}
}

func (l MockLogging) Error(args ...interface{}) {
}

func (MockLogging) Debug(args ...interface{}) {
}

func (MockLogging) Debugf(format string, args ...interface{}) {
}

func (MockLogging) Errorf(format string, args ...interface{}) {
}

func (MockLogging) Info(args ...interface{}) {
}

func (MockLogging) Infof(format string, args ...interface{}) {
}

func (MockLogging) Warn(args ...interface{}) {
}

func (MockLogging) Warnf(format string, args ...interface{}) {
}
