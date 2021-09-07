package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/output/mqtt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

const (
	stepAmount = 0.1
)

var (
	curPos = common.WallPos{}
	pub    *mqtt.Config
)

// TODO: Use a subset of env.Config?

type MQTTEnv struct {
	MQTTBroker   string `default:"mqtt.local:1883"`
	MQTTUsername string `default:"jan-poka"`
	MQTTPassword string `default:""`
	MQTTTopic    string `default:"home/geo/target"`

	TCPTimeout time.Duration `default:"1s"`
}

func main() {
	log.SetFlags(0)

	var environment MQTTEnv
	err := envconfig.Process("jp", &environment)
	if err != nil {
		panic(err)
	}

	pub, err = mqtt.New(
		environment.MQTTBroker,
		environment.MQTTUsername,
		environment.MQTTPassword,
		environment.MQTTTopic,
		environment.TCPTimeout,
	)
	if err != nil {
		panic(err)
	}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	log.Println("Move the mapper to your home point.")

	if err := initialReset(); err != nil {
		log.Fatalf("Could not complete zeroing! %v", err)
	}

	log.Println("Use Q and S to shorten and lengthen the left wire")
	log.Println("Use W and A to shorten and lengthen the right wire")

	stepMove("Press ESC to quit or Enter when you're done")

	pub.Publish(mqtt.Message{Reset: true})

	log.Println("All set! Your mapper is now reset and at its home position.")
}

func initialReset() error {
	log.Println("Press R to move to the home position from here, H to move to the existing home position.")
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			return fmt.Errorf("you quit the program")
		}

		switch char {
		case 'r', 'R':
			pub.Publish(mqtt.Message{Reset: true})
			return nil
		case 'h', 'H':
			curPos.LengthLeft = 0
			curPos.LengthRight = 0
			publishPosition()
			return nil
		}
	}
}

func stepMove(explain string) {
	log.Println(explain)
	if err := move(); err != nil {
		log.Fatalf("Could not complete zeroing! %v", err)
	}
}

func move() error {
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		switch key {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			log.Println()
			return fmt.Errorf("you quit the program")
		case keyboard.KeyEnter:
			log.Println()
			return nil
		}

		switch char {
		case 'q', 'Q':
			curPos.LengthLeft -= stepAmount
		case 's', 'S':
			curPos.LengthLeft += stepAmount
		case 'w', 'W':
			curPos.LengthRight -= stepAmount
		case 'a', 'A':
			curPos.LengthRight += stepAmount
		}

		publishPosition()
		fmt.Printf("\rLeft: %10.1f Right: %10.1f", float32(curPos.LengthLeft), float32(curPos.LengthRight))
	}
}

func publishPosition() {
	left := float32(curPos.LengthLeft)
	right := float32(curPos.LengthRight)

	pub.Publish(mqtt.Message{
		CalculatedMapperLeft:  left,
		CalculatedMapperRight: right,
		CalculatedMapperLengths: map[int]mqtt.Lengths{
			// TODO: What if more than 1?
			0: {Left: left, Right: right},
		},
	})
}
