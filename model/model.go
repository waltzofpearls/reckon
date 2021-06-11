package model

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/DataDog/go-python3"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
)

type Trainer interface {
	Train(context.Context, *python3.PyObject, prom.Metric, time.Duration) Forecasts
}

const TypeProphet = "Prophet"

func New(lg *zap.Logger, name string) (Trainer, error) {
	switch name {
	case TypeProphet:
		return NewProphet(lg), nil
	}
	return nil, fmt.Errorf("%s is not a valid model type", name)
}

type Column string

const (
	Yhat      Column = "yhat"
	YhatUpper Column = "yhat_upper"
	YhatLower Column = "yhat_lower"
	Original  Column = "original"
)

type Forecast struct {
	Timestamp int64
	Values    map[Column]float64
}

func (f Forecast) IsZero() bool {
	return f.Timestamp == 0 && len(f.Values) == 0
}

type Forecasts []Forecast

func (f Forecasts) Nearest(timestamp int64) Forecast {
	if len(f) == 0 {
		return Forecast{}
	}
	current := f[0]
	for _, forecast := range f {
		if math.Abs(float64(timestamp-forecast.Timestamp)) < math.Abs(float64(timestamp-current.Timestamp)) {
			current = forecast
		}
	}
	return current
}

func InitPythonModule(lg *zap.Logger) (module *python3.PyObject, cleanup func(), err error) {
	python3.Py_Initialize()

	if !python3.Py_IsInitialized() {
		lg.Error("error initializing python interpreter")
		return nil, nil, errors.New("error initializing python interpreter")
	}

	path, err := filepath.Abs(filepath.Dir(os.Args[0]) + "/model")
	if err != nil {
		lg.Error("error finding absolute /path/to/reckon/model", zap.Error(err))
		return nil, nil, fmt.Errorf("error finding absolute /path/to/reckon/model: %w", err)
	}

	returned := python3.PyRun_SimpleString("import sys\nsys.path.append(\"" + path + "\")")
	if returned != 0 {
		lg.Error("error appending to python sys.path", zap.String("path", path))
		return nil, nil, fmt.Errorf("error appending to python sys.path: %s", path)
	}

	imported := python3.PyImport_ImportModule("prophet_model") // return value: new ref
	if !(imported != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		lg.Error("failed to import module prophet_model")
		return nil, nil, errors.New("failed to import module prophet_model")
	}

	module = python3.PyImport_AddModule("prophet_model") // return value: borrowed ref (from imported)
	if !(module != nil && python3.PyErr_Occurred() == nil) {
		python3.PyErr_Print()
		lg.Error("failed to add module prophet_model")
		return nil, nil, errors.New("failed to add module prophet_model")
	}
	return module, func() {
		defer python3.Py_Finalize()
		defer imported.DecRef()
	}, nil
}
