package logx

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xuzq3/glib/logx/internal"
	"strings"
)

type logrusJsonFormatter struct{}

func (f *logrusJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buf := internal.BufferPool.Get()
	defer internal.BufferPool.Put(buf)

	buf.WriteByte('{')
	buf.WriteString(fmt.Sprintf(`"%s":"%s"`, levelKey, entry.Level.String()))

	buf.WriteByte(',')
	buf.WriteString(fmt.Sprintf(`"%s":"%s"`, timeKey, entry.Time.Format("2006-01-02T15:04:05.000")))

	_, fileVal := logrusCallerPrettyfier(entry.Caller)
	if fileVal != "" {
		buf.WriteByte(',')
		buf.WriteString(fmt.Sprintf(`"%s":"%s"`, callerKey, fileVal))
	}

	buf.WriteByte(',')
	message := entry.Message
	if message != "" {
		message = strings.TrimSuffix(message, "\n")
		message = strings.ReplaceAll(message, "\"", "\\\"")
	}
	buf.WriteString(fmt.Sprintf(`"%s":"%s"`, msgKey, message))

	if len(entry.Data) > 0 {
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

		//buf.WriteByte(',')
		//buf.WriteString(fmt.Sprintf(`"%s":`, dataKey))
		//buf.Write(b)

		buf.WriteByte(',')
		buf.Write(b[1 : len(b)-1])
	}

	//buf.WriteByte(',')
	//buf.WriteString(`"data":`)
	//encoder := json.NewEncoder(buf)
	//if err := encoder.Encode(data); err != nil {
	//	return nil, fmt.Errorf("failed to marshal fields to JSON, %v", err)
	//}
	//buf.Truncate(buf.Len() - 1)

	buf.WriteByte('}')
	buf.WriteByte('\n')

	return buf.Bytes(), nil
}
