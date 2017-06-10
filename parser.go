package configbuilder

type Parser interface {
	GetTypeString() string
	GenerateConfig(meta *configMeta) error
}
