package external_logger

type ExternalLogger interface {
	LogError(err error) error
}
