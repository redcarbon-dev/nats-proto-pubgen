// Code generated by pubgen. DO NOT EDIT.

package pubgen

type Publisher interface {
	Publish(topic string, msg []byte) error
}

type PubCli struct {
	p Publisher
}

func (p *PubCli) publish(topic string, msg []byte) error {
	return p.p.Publish(topic, msg)
}

func NewPublisher(p Publisher) *PubCli {
	return &PubCli{
		p: p,
	}
}
