package pkg

import (
	"github.com/cloud-native-toolkit/atkmod"
	cloudevent "github.com/cloudevents/sdk-go/v2"
)

func GetEventData(eventStr string) (*atkmod.EventData, error) {
	event, err := atkmod.LoadEvent(eventStr)
	if err != nil {
		return nil, err
	}
	data, err := atkmod.LoadEventData(event)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func NewEventWithData(evntType string, src string, sub string, data *atkmod.EventData) cloudevent.Event {
	event := cloudevent.NewEvent()
	event.SetType(evntType)
	event.SetSource(src)
	event.SetSubject(sub)
	event.SetData(cloudevent.ApplicationJSON, data)
	return event
}
