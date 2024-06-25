package utils

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxKey int

var ctxContentKey ctxKey

func addValuesToContext(
	ctx context.Context,
	addValues func(map[string]interface{}) map[string]interface{},
) context.Context {
	var values map[string]interface{}
	values, ok := ctx.Value(ctxContentKey).(map[string]interface{})
	if !ok {
		values = make(map[string]interface{})
	}
	values = addValues(values)
	return context.WithValue(ctx, ctxContentKey, values)
}

func WithValue(ctx context.Context, key string, value interface{}) context.Context {
	return addValuesToContext(
		ctx,
		func(vals map[string]interface{}) map[string]interface{} {
			vals[key] = value
			return vals
		},
	)
}

func WithValues(ctx context.Context, newValues map[string]interface{}) context.Context {
	return addValuesToContext(
		ctx,
		func(values map[string]interface{}) map[string]interface{} {
			for key, value := range newValues {
				values[key] = value
			}
			return values
		},
	)
}

func HasRole(roles []string, requestedRole string) bool {
	for _, role := range roles {
		if role == requestedRole {
			return true
		}
	}
	return false
}

func GetJWTData(ctx context.Context) (int64, []string, bool) {
	var values map[string]interface{}
	values, ok := ctx.Value(ctxContentKey).(map[string]interface{})
	if !ok {
		return 0, nil, false
	}
	return values["userId"].(int64), values["userRoles"].([]string), true
}

type ContextLoggerHook struct{}

func (clh *ContextLoggerHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.PanicLevel,
		logrus.FatalLevel,
	}
}

func (clh *ContextLoggerHook) Fire(entry *logrus.Entry) error {
	if entry.Context != nil {
		ctx := entry.Context
		var values map[string]interface{}
		values, ok := ctx.Value(ctxContentKey).(map[string]interface{})
		if ok {
			entry.Data["app_context"] = values
		}
	}
	return nil
}
