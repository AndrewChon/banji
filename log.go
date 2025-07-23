package banji

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

type LogKind struct {
	Kind        string
	Description string
}

type Severity int8

const (
	DebugSeverity Severity = iota - 1
	InfoSeverity
	WarnSeverity
	ErrorSeverity
	PanicSeverity
	CriticalSeverity
)

type Field struct {
	K string
	V any
}
type LogProperties []Field

type Log struct {
	id          uuid.UUID
	kind        string
	description string
	postmark    time.Time
	severity    Severity
	properties  LogProperties
}

func NewLog(k LogKind, s Severity) *Log {
	return &Log{
		id:          uuid.New(),
		kind:        k.Kind,
		description: k.Description,
		postmark:    time.Now(),
		severity:    s,
		properties:  LogProperties{},
	}
}

func NewLogWithProperties(k LogKind, s Severity, fields []Field) *Log {
	l := NewLog(k, s)
	for _, field := range fields {
		l.properties = append(l.properties, field)
	}
	return l
}

func (l *Log) ID() uuid.UUID {
	return l.id
}

func (l *Log) Kind() string {
	return l.kind
}

func (l *Log) Description() string {
	return l.description
}

func (l *Log) Postmark() time.Time {
	return l.postmark
}

func (l *Log) Severity() Severity {
	return l.severity
}

func (l *Log) Properties() LogProperties {
	return l.properties
}

func (l *Log) With(k string, v any) *Log {
	l.properties = append(l.properties, Field{k, v})
	return l
}

func (l *Log) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"id":          l.id.String(),
		"kind":        l.kind,
		"description": l.description,
		"postmark":    l.postmark,
		"severity":    l.severity,
	}

	for _, f := range l.properties {
		m[strings.ToLower(f.K)] = f.V
	}

	return json.Marshal(m)
}
