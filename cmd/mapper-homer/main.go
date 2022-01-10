package main

import (
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/jphastings/jan-poka/pkg/common"
	"github.com/jphastings/jan-poka/pkg/output/mapper"
	"github.com/jphastings/jan-poka/pkg/output/mqtt"
	"github.com/kelseyhightower/envconfig"
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
	MQTTPort int `default:"1883"`

	Persistence string `default:"~/.jan-poka"`
}

func check(err error) {
	if err != nil {
		log.Fatalf("Could not complete zeroing! %v", err)
	}
}

func main() {
	log.SetFlags(0)

	var environment MQTTEnv
	err := envconfig.Process("jp", &environment)
	check(err)
	if strings.HasPrefix(environment.Persistence, "~/") {
		usr, err := user.Current()
		check(err)
		environment.Persistence = filepath.Join(usr.HomeDir, environment.Persistence[2:])
	}

	pub, err = mqtt.New(environment.MQTTPort, environment.Persistence)
	check(err)

	check(keyboard.Open())
	defer keyboard.Close()

	check(phaseInit())

	log.Println("Use Q and S to shorten and lengthen the left wire")
	log.Println("Use W and A to shorten and lengthen the right wire")

	check(stepMove("Move the mapper to your home point."))

	resetPosition()

	check(phaseAddMaps(environment.Persistence))
}

func phaseInit() error {
	log.Println("Press R to find the home position from here, H to move to the previous home position and start from there.")
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
			curPos.Left = 0
			curPos.Right = 0
			publishPosition()
			return nil
		}
	}
}

func phaseAddMaps(configRoot string) error {
	m, err := mapper.New(configRoot)
	if err != nil {
		return err
	}

	var mapper mapper.State

	// TODO: Deal with more than one mapper
	if len(m.Mappers) >= 1 {
		mapper = m.Mappers[0]
	}

	for mapID := range mapper.Maps {
		// TODO: How to georeference?
		_ = mapID
	}
	return fmt.Errorf("not implemented")
}

func stepMove(explain string) error {
	log.Println(explain)
	return move()
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
			curPos.Left -= stepAmount
		case 's', 'S':
			curPos.Left += stepAmount
		case 'w', 'W':
			curPos.Right -= stepAmount
		case 'a', 'A':
			curPos.Right += stepAmount
		}

		publishPosition()
		fmt.Printf("\rLeft: %10.1f Right: %10.1f", float32(curPos.Left), float32(curPos.Right))
	}
}

func publishPosition() {
	left := float32(curPos.Left)
	right := float32(curPos.Right)

	pub.Publish(mqtt.Message{
		CalculatedMapperLeft:  left,
		CalculatedMapperRight: right,
		CalculatedMapperLengths: map[int]mqtt.Lengths{
			// TODO: What if more than 1?
			0: {Left: left, Right: right},
		},
	})
}

func resetPosition() {
	pub.Publish(mqtt.Message{Reset: true})
}
