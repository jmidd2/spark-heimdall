package device

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Device represents any connectable device (PC, server, etc.)
type Device struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IPAddress   string `json:"ip_address"`
	Protocol    string `json:"protocol"` // "vnc", "rdp", etc.
	Port        int    `json:"port"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	FullScreen  bool   `json:"full_screen"`
	Description string `json:"description,omitempty"`
	Screen      string `json:"screen,omitempty"`
}

type Devices []Device

func (m *Store) GetAll() Devices {
	return m.Devices
}

func (m *Store) Get(id string) (Device, bool) {
	for _, device := range m.Devices {
		if device.ID == id {
			return device, true
		}
	}

	return Device{}, false
}

func (m *Store) Add(device Device) error {
	// Generate ID if not provided
	if m.highestDeviceId == 0 {
		if err := m.findHighestDeviceId(); err != nil {
			return err
		}
	} else {
		m.highestDeviceId += 1
	}

	if device.ID == "" {
		device.ID = fmt.Sprintf("pc%d", m.highestDeviceId)
	}

	err := m.Devices.ValidateNew(device)
	if err != nil {
		return err
	}

	m.Devices = append(m.Devices, device)

	return nil
}

func (m *Store) Update(device Device) error {
	for i, existingPC := range m.Devices {
		if existingPC.ID == device.ID {
			m.Devices[i] = device
			return nil
		}
	}

	return fmt.Errorf("PC with ID %s not found", device.ID)
}

func (m *Store) Delete(id string) error {
	newDevices := make([]Device, 0, len(m.Devices))
	found := false

	for _, d := range m.Devices {
		if d.ID == id {
			found = true
		} else {
			newDevices = append(newDevices, d)
		}
	}

	if !found {
		return fmt.Errorf("PC with ID %s not found", id)
	}

	m.Devices = newDevices

	return nil
}

func (d Devices) ValidateNew(device Device) error {
	// Check for duplicate ID
	for _, existingPC := range d {
		if existingPC.ID == device.ID {
			return fmt.Errorf("PC with ID %s already exists", device.ID)
		}
	}

	return nil
}

func (m *Store) findHighestDeviceId() error {
	highestDeviceId := 1
	for _, pc := range m.Devices {
		part := strings.Split(pc.ID, "pc")
		if len(part) != 2 {
			log.Printf("invalid PC ID: %s", pc.ID)
			continue
		}

		test, err := strconv.Atoi(part[1])
		if err != nil {
			return err
		}

		if test > highestDeviceId {
			highestDeviceId = test
		}
	}
	m.highestDeviceId = highestDeviceId + 1
	log.Printf("Set highest device id to: %d", m.highestDeviceId)
	return nil
}

// DevicesStore defines operations for managing devices
type DevicesStore interface {
	GetAll() Devices
	Get(id string) (Device, bool)
	Add(device Device) error
	Update(device Device) error
	Delete(id string) error
}

type Store struct {
	Devices         Devices `json:"devices"`
	highestDeviceId int
}

// Ensure Manager implements DeviceManager
var _ DevicesStore = (*Store)(nil)
