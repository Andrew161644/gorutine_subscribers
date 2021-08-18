package local_resolver

import (
	"log"
	ss "proj/handlers"
	"proj/resolvers"
)

type LocalResolver struct {
	resolvers.AbstractResolver
}

func (resolver *LocalResolver) GoBroadCastEvent() {
	log.Println("examples listening started")
	subscribers := map[string]ss.ISubscriber{}
	unsubscribe := make(chan string)
	for {
		select {
		case id := <-unsubscribe:
			delete(subscribers, id)
		case s := <-resolver.Subscribers:
			subscribers[s.GetId()] = s
			s.Handle()
		case e := <-resolver.Commands:
			for id, s := range subscribers {
				go func(id string, s ss.ISubscriber) {
					select {
					case <-s.GetStop():
						unsubscribe <- id
						return
					default:
					}
					select {
					case <-s.GetStop():
						unsubscribe <- id
					case s.GetEvents() <- e:
					}
				}(id, s)
			}
		}
	}
}

func NewLocalResolver() *LocalResolver {
	return &LocalResolver{
		*resolvers.NewResolver(),
	}
}
