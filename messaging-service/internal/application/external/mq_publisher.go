package external

type Event any

type MqPublisher interface {
	Publish(Event) error
}
