package banji

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Level int8

const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

type LogField struct {
	K string `json:"k"`
	V any    `json:"v"`
}

type LogData []LogField

type Log struct {
	id          uuid.UUID
	postmark    time.Time
	kind        string
	description string
	level       Level
	data        LogData
}

func NewLog(kind, description string, level Level) *Log {
	return &Log{
		id:          uuid.New(),
		postmark:    time.Now(),
		kind:        kind,
		description: description,
		level:       level,
		data:        LogData{},
	}
}

func (l *Log) ID() uuid.UUID {
	return l.id
}

func (l *Log) Postmark() time.Time {
	return l.postmark
}

func (l *Log) Kind() string {
	return l.kind
}

func (l *Log) Description() string {
	return l.description
}

func (l *Log) Level() Level {
	return l.level
}

func (l *Log) Data() LogData {
	return l.data
}

func (l *Log) Value(key string) any {
	for _, field := range l.data {
		if field.K == key {
			return field.V
		}
	}

	return nil
}

func (l *Log) Add(k string, v any) *Log {
	l.data = append(l.data, LogField{
		K: k,
		V: v.(string),
	})
	return l
}

func (l *Log) With(x ...any) error {
	var logData LogData
	length := len(x)

	if length%2 != 0 || length == 0 {
		return errors.New("arguments must be in pairs")
	}

	for i := 0; i < length; i += 2 {
		k, ok := x[i].(string)
		if !ok {
			return errors.New("key must be a string")
		}

		v := x[i+1]
		logData = append(logData, LogField{
			K: k,
			V: v,
		})
	}

	l.data = logData
	return nil
}

func (ld LogData) MarshalJSON() ([]byte, error) {
	dataMap := make(map[string]any, len(ld))
	for _, field := range ld {
		dataMap[field.K] = field.V
	}

	return json.Marshal(dataMap)
}
