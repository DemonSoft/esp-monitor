package domain

import (
	"fmt"
	"strings"
)

type DomainLog struct {
	kafka Kafka
}

type Kafka interface {
	LogMessage(str string)
}

var Log = DomainLog{}

func Create() DomainLog {
	return Log
}

func (d *DomainLog) Log(params ...any) {
	var b strings.Builder
	for _, val := range params {
		v := fmt.Sprintf("%v", val)
		b.WriteString(v)
	}

	result := b.String()
	fmt.Println(result)

	if d.kafka != nil {
		d.kafka.LogMessage(result)
	}
}

func (s *DomainLog) UpdateKafka(kafka Kafka) {
	s.kafka = kafka
	Log = *s
}
