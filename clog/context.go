package clog

import "context"

type fieldMapType struct{}

var fieldMap = fieldMapType{}

func (l *CustomLogger) AddKeysValuesToCtx(ctx context.Context, kv map[string]interface{}) context.Context {
	fields := ctx.Value(fieldMap)

	if fields == nil {
		return context.WithValue(ctx, fieldMap, kv)
	}

	for k, v := range kv {
		fields.(map[string]interface{})[k] = v
	}

	return context.WithValue(ctx, fieldMap, fields)
}

func (l *CustomLogger) fieldsFromCtx(ctx context.Context) map[string]interface{} {
	return ctx.Value(fieldMap).(map[string]interface{})
}
