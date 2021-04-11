package main

import (
	"github.com/jphastings/jan-poka/pkg/pointer/mqtt"
	"github.com/jphastings/jan-poka/pkg/tracker"
)

func init() {
	configurables = append(configurables, configurable{
		"Publishing to MQTT",
		func() bool { return environment.MQTTBroker != "" && environment.MQTTTopic != "" },
		configureMQTT,
	})
}

func configureMQTT() (tracker.OnTracked, error) {
	pub, err := mqtt.New(
		environment.MQTTBroker,
		environment.MQTTUsername,
		environment.MQTTPassword,
		environment.MQTTTopic,
		environment.TCPTimeout,
	)
	if err != nil {
		return nil, err
	}

	return pub.TrackerCallback, nil
}
