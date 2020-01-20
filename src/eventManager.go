package main

type event struct {
	name string
	data interface{}
}

type subscriber interface {
	notify(ev *event)
}

type eventManager struct {
	subscribers map[string][]*subscriber
}

func newEventManager() *eventManager {
	return &eventManager{subscribers: make(map[string][]*subscriber)}
}

func (em *eventManager) subscribe(eventName string, sub *subscriber) {
	em.subscribers[eventName] = append(em.subscribers[eventName], sub)
}

func (em *eventManager) unsubscribe(eventName string, sub *subscriber) {
	if _, ok := em.subscribers[eventName]; !ok {
		return
	}

	loc := -1
	for i, currentSub := range em.subscribers[eventName] {
		if currentSub == sub {
			loc = i
			break
		}
	}

	if loc == -1 {
		return
	}

	em.subscribers[eventName] = append(em.subscribers[eventName][:loc], em.subscribers[eventName][loc+1:]...)
}

func (em *eventManager) notifySubscribers(ev *event) {
	if _, ok := em.subscribers[ev.name]; !ok {
		return
	}

	for _, sub := range em.subscribers[ev.name] {
		(*sub).notify(ev)
	}
}
