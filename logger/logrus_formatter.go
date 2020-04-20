package logger

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/xuzq3/glib/logger/internal"
)

type logrusJsonFormatter struct{}

func (f *logrusJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buf := internal.BufferPool.Get()
	defer internal.BufferPool.Put(buf)

	if entry.Data["type"] == "biz" {
		buf.WriteString("[datacentercollector]")
	}

	buf.WriteByte('{')
	buf.WriteString(fmt.Sprintf(`"level":"%s"`, entry.Level.String()))

	buf.WriteByte(',')
	buf.WriteString(fmt.Sprintf(`"time":"%s"`, entry.Time.Format("2006-01-02T15:04:05.000")))

	_, fileVal := logrusCallerPrettyfier(entry.Caller)
	if fileVal != "" {
		buf.WriteByte(',')
		buf.WriteString(fmt.Sprintf(`"caller":"%s"`, fileVal))
	}

	buf.WriteByte(',')
	buf.WriteString(fmt.Sprintf(`"msg":"%s"`, entry.Message))

	data := make(Fields, len(entry.Data))
	for k, v := range entry.Data {
		switch v := v.(type) {
		case error:
			data[k] = v.Error()
		default:
			data[k] = v
		}
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	}
	buf.WriteByte(',')
	buf.WriteString(fmt.Sprintf(`"data":`))
	buf.Write(b)

	//buf.WriteByte(',')
	//buf.WriteString(`"data":`)
	//encoder := json.NewEncoder(buf)
	//if err := encoder.Encode(data); err != nil {
	//	return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	//}
	//buf.Truncate(buf.Len() - 1)

	buf.WriteByte('}')
	if entry.Data["type"] == "biz" {
		buf.WriteString("[datacentercollector]")
	}
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}
