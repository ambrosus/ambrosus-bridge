package external_logger

type ExternalLogger interface {
	LogError(msg string) error
}
