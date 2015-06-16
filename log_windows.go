package main

//this is a wrapper log formatter to augment the default logrus TextFormatter
//to terminate log lines on windows with \r\n instead of just \n

import (
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetFormatter(&carriageReturnFormater{new(log.TextFormatter)})
}

type carriageReturnFormater struct {
	originalFormatter log.Formatter
}

func (f *carriageReturnFormater) Format(entry *log.Entry) ([]byte, error) {
	bytes, err := f.originalFormatter.Format(entry)
	bytes[len(bytes)-1] = '\r'
	bytes = append(bytes, '\n')
	return bytes, err
}
