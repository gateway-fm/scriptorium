package clog

import "context"

type fieldMapType struct{}

var fieldMap = fieldMapType{}

func (l *CustomLogger) AddKeysValuesToCtx(ctx context.Context, kv map[string]interface{}) context.Context {
	fields := ctx.Value(fieldMap)

	if fields == nil {
		return context.WithValue(ctx, fieldMap, kv)
	}

	l.mu.Lock()

	for k, v := range kv {
		if v != nil {

		}
		fields.(map[string]interface{})[k] = v
	}

	l.mu.Unlock()

	return context.WithValue(ctx, fieldMap, fields)
}

func (l *CustomLogger) fieldsFromCtx(ctx context.Context) map[string]interface{} {
	if fm := ctx.Value(fieldMap); fm != nil {
		return fm.(map[string]interface{})
	}

	return nil
}
