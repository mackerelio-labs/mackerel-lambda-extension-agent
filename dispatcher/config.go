package dispatcher

type DispatcherConfig struct {
	MinBatchSize int64 `env:"DISPATCH_MIN_BATCH_SIZE" envDefault:"1"`
}
