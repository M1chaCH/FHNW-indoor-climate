package obs

import (
	"fmt"
	"sensor_hub_backend/logs"
)

type Observable[T any] struct {
	Name     string
	TypeName string
	channels map[int]chan<- T
	index    int
}

func NewObservable[T any](name string) *Observable[T] {
	typeName := fmt.Sprintf("%T", (*T)(nil))
	o := Observable[T]{
		Name:     name,
		TypeName: typeName,
		channels: make(map[int]chan<- T),
		index:    0,
	}
	logs.LogInfo("'%v' created", o)
	return &o
}

// Emit loops over every subscribed channel and emits the value
// (this action is blocking, if a channel has no listener anymore, then we are stuck here...)
func (o *Observable[T]) Emit(value T) {
	logs.LogDebug("'%v' emitting value", o)
	for _, channel := range o.channels {
		channel <- value
	}
}

func (o *Observable[T]) NewChannel() (<-chan T, int) {
	channel := make(chan T)
	return channel, o.Subscribe(channel)
}

func (o *Observable[T]) Subscribe(channel chan<- T) int {
	i := o.index
	o.channels[i] = channel
	o.index++
	logs.LogDebug("'%v' New subscription (%d)", o, i)
	return i
}

// Unsubscribe removes the channel from the list of subscribers
// IT IS CRUCIAL that this method is called when the channel is no longer needed,
// otherwise the Emit method will block forever
func (o *Observable[T]) Unsubscribe(index int) {
	delete(o.channels, index)
	logs.LogDebug("'%v' Unsubscribed (%d)", o, index)
}

func (o *Observable[T]) String() string {
	return fmt.Sprintf("Observable[%s] (%s)", o.Name, o.TypeName)
}
