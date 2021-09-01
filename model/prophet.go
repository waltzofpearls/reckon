package model

import (
	"context"
	"runtime"
	"time"

	"github.com/DataDog/go-python3"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
)

type Prophet struct {
	logger *zap.Logger
}

func NewProphet(lg *zap.Logger) Prophet {
	return Prophet{
		logger: lg,
	}
}

func (p Prophet) Train(ctx context.Context, module *python3.PyObject, data prom.Metric, duration time.Duration) Forecasts {
	logger := p.logger.With(zap.String("metric_name", data.Name), zap.Any("metric_labels", data.Labels))
	logger.Info("train model with data", zap.Int("length", len(data.Values)), zap.String("duration", duration.String()))

	runtime.LockOSThread()
	s := python3.PyGILState_Ensure()
	defer python3.PyGILState_Release(s)

	argData := python3.PyList_New(len(data.Values)) // return value: new reference, gets stolen later
	for i, value := range data.Values {
		samplePair := python3.PyList_New(2)
		sampleTime := python3.PyLong_FromLongLong(int64(value.Time.Unix()))
		sampleValue := python3.PyFloat_FromDouble(value.Value)
		retTime := python3.PyList_SetItem(samplePair, 0, sampleTime)
		retValue := python3.PyList_SetItem(samplePair, 1, sampleValue)
		retSample := python3.PyList_SetItem(argData, i, samplePair)
		if retTime != 0 || retValue != 0 || retSample != 0 {
			if python3.PyErr_Occurred() != nil {
				python3.PyErr_Print()
			}
			sampleTime.DecRef()
			sampleValue.DecRef()
			samplePair.DecRef()
			argData.DecRef()
			logger.Error("error setting sample")
			return nil
		}
		sampleTime = nil
		sampleValue = nil
		samplePair = nil
	}
	argDuration := python3.PyFloat_FromDouble(duration.Minutes())

	args := python3.PyTuple_New(2) // return value: new reference
	if args == nil {
		argData.DecRef()
		argDuration.DecRef()
		logger.Error("error creating args tuple")
		return nil
	}
	defer args.DecRef()
	if ret := python3.PyTuple_SetItem(args, 0, argData); ret != 0 { // steals ref to argData
		if python3.PyErr_Occurred() != nil {
			python3.PyErr_Print()
		}
		argData.DecRef()
		argDuration.DecRef()
		argData = nil
		logger.Error("error setting args tuple argData")
		return nil
	}
	argData = nil
	if ret := python3.PyTuple_SetItem(args, 1, argDuration); ret != 0 {
		if python3.PyErr_Occurred() != nil {
			python3.PyErr_Print()
		}
		argDuration.DecRef()
		logger.Error("error setting args tuple argDuration")
		return nil
	}
	argDuration = nil

	dict := python3.PyModule_GetDict(module) // return value: borrowed
	if !(dict != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		logger.Error("could not get dict for module")
		return nil
	}
	train := python3.PyDict_GetItemString(dict, "train")
	if !(train != nil && python3.PyCallable_Check(train)) { // return value: borrowed
		logger.Error("could not find function train()")
		return nil
	}
	returned := train.CallObject(args)
	if !(returned != nil && python3.PyErr_Occurred() == nil) { // return value: new reference
		python3.PyErr_Print()
		logger.Error("error calling function detect")
		return nil
	}
	defer returned.DecRef()

	forecasts, err := toForecastsList(returned)
	if err != nil {
		logger.Error("error converting python dict to go map", zap.Error(err))
		return nil
	}

	return forecasts
}

func toForecastsList(pyDict *python3.PyObject) (Forecasts, error) {
	var (
		pos       int
		key, val  *python3.PyObject
		forecasts Forecasts
	)
	for python3.PyDict_Next(pyDict, &pos, &key, &val) {
		timestamp := int64(python3.PyFloat_AsDouble(key))
		if timestamp == -1 || python3.PyErr_Occurred() != nil {
			continue
		}
		var (
			subPos         int
			subKey, subVal *python3.PyObject
		)
		values := make(map[Column]float64)
		for python3.PyDict_Next(val, &subPos, &subKey, &subVal) {
			column := python3.PyUnicode_AsUTF8(subKey)
			if column == "" || python3.PyErr_Occurred() != nil {
				continue
			}
			value := python3.PyFloat_AsDouble(subVal)
			if value == -1 || python3.PyErr_Occurred() != nil {
				continue
			}
			values[Column(column)] = value
		}
		subKey = nil
		subVal = nil
		forecasts = append(forecasts, Forecast{
			Timestamp: timestamp,
			Values:    values,
		})
	}
	key = nil
	val = nil

	return forecasts, nil
}
