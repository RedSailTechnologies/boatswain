package template

import (
	"github.com/redsailtechnologies/boatswain/pkg/ddd"
)

var entityName = "Template"

// Created is the event for when a new template is created
type Created struct {
	Timestamp int64
	UUID      string
	Name      string
	Type      Type
	YAML      string
}

// EventType marks this as an event
func (e Created) EventType() string {
	return entityName + "Created"
}

// Destroyed is the event for when a template is destroyed
type Destroyed struct {
	Timestamp int64
}

// EventType marks this as an event
func (e Destroyed) EventType() string {
	return entityName + "Destroyed"
}

// Updated is the event for when a template is updated
type Updated struct {
	Timestamp int64
	UUID      string
	Name      string
	Type      Type
	YAML      string
}

// EventType marks this as an event
func (e Updated) EventType() string {
	return entityName + "Updated"
}

// Template represents a template for a delivery object
type Template struct {
	events    []ddd.Event
	version   int
	destroyed bool

	uuid  string
	name  string
	ttype Type
	yaml  string
}

// Replay recreates the template from a series of events
func Replay(events []ddd.Event) *Template {
	t := &Template{}
	for _, event := range events {
		t.on(event)
	}
	return t
}

// Create handles create commands
func Create(uuid, name, yaml string, ttype Type, timestamp int64) (*Template, error) {
	if uuid == "" {
		return nil, ddd.IDError{}
	}
	err := validateFields(name, yaml, ttype)
	if err != nil {
		return nil, err
	}

	t := &Template{}
	t.on(&Created{
		Timestamp: timestamp,
		UUID:      uuid,
		Name:      name,
		Type:      ttype,
		YAML:      yaml,
	})
	return t, nil
}

// UUID returns this template's identifier
func (t *Template) UUID() string {
	return t.uuid
}

// Name returns this template's name
func (t *Template) Name() string {
	return t.uuid
}

// Type returns this template's type
func (t *Template) Type() Type {
	return t.ttype
}

// YAML returns this template's yaml
func (t *Template) YAML() string {
	return t.yaml
}

func (t *Template) on(event ddd.Event) {
	t.events = append(t.events, event)
	t.version++
	switch e := event.(type) {
	case *Created:
		t.uuid = e.UUID
		t.name = e.Name
		t.ttype = e.Type
		t.yaml = e.YAML
	case *Destroyed:
		t.destroyed = true
	case *Updated:
		t.name = e.Name
		t.ttype = e.Type
		t.yaml = e.YAML
	}
}

func validateFields(name, yaml string, t Type) error {
	if name == "" {
		return ddd.RequiredArgumentError{Arg: "Name"}
	}
	if yaml == "" {
		return ddd.RequiredArgumentError{Arg: "YAML"}
	}
	// TODO AdamP - check the template is valid here based on its type?
	return nil
}
