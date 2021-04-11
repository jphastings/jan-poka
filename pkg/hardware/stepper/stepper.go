package stepper

import (
	"fmt"
	"github.com/jphastings/jan-poka/pkg/math"
	"io/ioutil"
	"os"
	"path"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/host/rpi"
	"strconv"
	"time"
)

var stateDir string = "/run/jan-poka/"
var Motors = map[string]math.Degrees{
	"28BYJ-48": 45 / 256.0,
}

var stepOff = stepperPins{gpio.Low, gpio.Low, gpio.Low, gpio.Low}
var stepSeq = []stepperPins{
	{gpio.High, gpio.Low, gpio.Low, gpio.High},
	{gpio.High, gpio.Low, gpio.Low, gpio.Low},
	{gpio.High, gpio.High, gpio.Low, gpio.Low},
	{gpio.Low, gpio.High, gpio.Low, gpio.Low},
	{gpio.Low, gpio.High, gpio.High, gpio.Low},
	{gpio.Low, gpio.Low, gpio.High, gpio.Low},
	{gpio.Low, gpio.Low, gpio.High, gpio.High},
	{gpio.Low, gpio.Low, gpio.Low, gpio.High},
}
var seqLen = len(stepSeq)

type stepperPins struct {
	p0 gpio.Level
	p1 gpio.Level
	p2 gpio.Level
	p3 gpio.Level
}

type Stepper struct {
	Name      string
	stateFile string

	CurrentAngle math.Degrees
	anglePerStep math.Degrees
	currentStep  int

	p0 gpio.PinOut
	p1 gpio.PinOut
	p2 gpio.PinOut
	p3 gpio.PinOut

	stepSpacing time.Duration
}

func SetStateDir(path string) error {
	err := os.Mkdir(path, os.ModeDir)
	if err != nil {
		return err
	}
	stateDir = path
	return nil
}

func Pi2Quad(anglePerStep math.Degrees) []*Stepper {
	return []*Stepper{
		New("M0", rpi.P1_11, rpi.P1_12, rpi.P1_13, rpi.P1_15, 0, anglePerStep),
		New("M1", rpi.P1_16, rpi.P1_18, rpi.P1_22, rpi.P1_7, 0, anglePerStep),
		New("M2", rpi.P1_33, rpi.P1_32, rpi.P1_31, rpi.P1_29, 0, anglePerStep),
		New("M3", rpi.P1_38, rpi.P1_37, rpi.P1_36, rpi.P1_35, 0, anglePerStep),
	}
}

func New(name string, p0, p1, p2, p3 gpio.PinOut, startAngle, anglePerStep math.Degrees) *Stepper {
	s := &Stepper{
		Name:      name,
		stateFile: path.Join(stateDir, name),

		CurrentAngle: startAngle,
		anglePerStep: anglePerStep,

		currentStep: 0,
		p0:          p0,
		p1:          p1,
		p2:          p2,
		p3:          p3,

		stepSpacing: 1 * time.Millisecond,
	}
	_ = s.readFromState()
	return s
}

func (s *Stepper) SetSpeed(rpm float64) {
	usPerStep := (s.anglePerStep * 60000000) / (math.Degrees(rpm) * 360)
	s.stepSpacing = time.Duration(usPerStep) * time.Microsecond
}

func (s *Stepper) Off() error {
	err := s.applyStep(stepOff)
	if err != nil {
		return err
	}
	return s.writeToState()
}

func (s *Stepper) SetAngle(angle math.Degrees) error {
	angleChange := math.ModDeg(angle - s.CurrentAngle)
	if angleChange > 180 {
		angleChange = angleChange - 360
	}

	steps := int(angleChange / s.anglePerStep)
	return s.Step(steps)
}

// +ve is clockwise
func (s *Stepper) Step(steps int) error {
	stepUnit := 2 // TODO: Microstepping?
	anglePerStep := s.anglePerStep
	if steps < 0 {
		anglePerStep *= -1
		stepUnit = -2
		steps *= -1
	}

	for i := 0; i < steps; i++ {
		s.currentStep = (s.currentStep + stepUnit + seqLen) % seqLen
		if err := s.applyStep(stepSeq[s.currentStep]); err != nil {
			return err
		}
		s.CurrentAngle += anglePerStep
		<-time.NewTimer(s.stepSpacing).C
	}
	return nil
}

func (s *Stepper) applyStep(seq stepperPins) error {
	if err := s.p0.Out(seq.p0); err != nil {
		return err
	}
	if err := s.p1.Out(seq.p1); err != nil {
		return err
	}
	if err := s.p2.Out(seq.p2); err != nil {
		return err
	}
	if err := s.p3.Out(seq.p3); err != nil {
		return err
	}
	return nil
}

func (s *Stepper) readFromState() error {
	data, err := ioutil.ReadFile(s.stateFile)
	if err != nil {
		return err
	}
	degs, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	s.CurrentAngle = math.Degrees(degs)
	return nil
}

func (s *Stepper) writeToState() error {
	return ioutil.WriteFile(s.stateFile, []byte(fmt.Sprintf("%f", s.CurrentAngle)), 0644)
}