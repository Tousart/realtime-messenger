package infrastructure

type Consumer interface {
	ConsumeMessages()
}
