package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

//Plug defines the Plug/Relais Device
type Plug struct {
	livePin        *gpio.DirectPinDriver
	neutralPin     *gpio.DirectPinDriver
	servo          *gpio.ServoDriver
	rawServo       *gpio.DirectPinDriver
	firmataAdaptor *firmata.FirmataAdaptor
}

//NewPlug returns a new Plug structure
func NewPlug(serial string, livePinNumber string, neutralPinNumber string, servoPinNumber string) (*Plug, error) {
	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", serial)
	livePin := gpio.NewDirectPinDriver(firmataAdaptor, "pin", livePinNumber)
	neutralPin := gpio.NewDirectPinDriver(firmataAdaptor, "pin", neutralPinNumber)
	servo := gpio.NewServoDriver(firmataAdaptor, "servo", servoPinNumber)
	rawServo := gpio.NewDirectPinDriver(firmataAdaptor, "pin", servoPinNumber)
	if err := firmataAdaptor.Connect(); err != nil {
		return nil, err[0]
	}
	servo.Move(40)
	gobot.After(1*time.Second, func() {
		rawServo.DigitalRead()
	})
	return &Plug{livePin, neutralPin, servo, rawServo, firmataAdaptor}, nil
}

func (p *Plug) set(level byte) {
	p.livePin.DigitalWrite(level)
	p.neutralPin.DigitalWrite(level)
}

//On turns on the connected device
func (p *Plug) On() {
	p.set(byte(0))
}

//Off turns off the connected device
func (p *Plug) Off() {
	p.set(byte(1))
}

func (p *Plug) Input() {
	p.servo.Move(uint8(0))
	gobot.After(1*time.Second, func() {
		p.servo.Move(40)
		gobot.After(1*time.Second, func() {
			p.rawServo.DigitalRead()
		})
	})
}

//Disconnect from the arduino
func (p *Plug) Disconnect() error {
	return p.firmataAdaptor.Disconnect()
}
