package external_logger

type ExternalLogger interface {
	LogError(prefix, msg string) error
	LogWarning(prefix, msg string) error
}
