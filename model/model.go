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
	"github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/exports"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
)

type Trainer interface {
	Train(context.Context, *python3.PyObject, prom.Metric, time.Duration) (Forecasts, error)
}

const (
	TypeProphet = "Prophet"
	TypeTangram = "Tangram"
)

func New(lg *zap.Logger, name string) (Trainer, error) {
	switch name {
	case TypeProphet:
		return NewProphet(lg), nil
	case TypeTangram:
		return NewTangram(lg), nil
	}
	return nil, fmt.Errorf("%s is not a valid model type", name)
}

type Column string

const (
	Yhat      Column = "yhat"
	YhatUpper Column = "yhat_upper"
	YhatLower Column = "yhat_lower"
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

	state := python3.PyEval_SaveThread()

	return module, func() {
		imported.DecRef()
		python3.PyEval_RestoreThread(state)
		python3.Py_Finalize()
	}, nil
}

func timeSeriesToSupervised(data *dataframe.DataFrame, nIn, nOut int) *dataframe.DataFrame {
	if nIn == 0 && nOut == 0 {
		return dataframe.NewDataFrame(data.Series...)
	}

	series := data.Series
	nVars := len(series)
	var inAndOut []dataframe.Series
	// input sequence (t-n, ... t-1)
	for i := nIn; i > 0; i-- {
		for j := 0; j < nVars; j++ {
			name := fmt.Sprintf("%s(t-%d)", series[j].Name(), i)
			nRows := series[j].NRows()
			start, end := 0, nRows-1-i
			copied := series[j].Copy(dataframe.Range{Start: &start, End: &end})
			copied.Rename(name)
			for k := 0; k < i; k++ {
				copied.Prepend(nil)
			}
			inAndOut = append(inAndOut, copied)
		}
	}
	// forecast sequence (t, t+1, ... t+n))
	for i := 0; i < nOut; i++ {
		for j := 0; j < nVars; j++ {
			name := series[j].Name()
			if i == 0 {
				name = fmt.Sprintf("%s(t)", name)
			} else {
				name = fmt.Sprintf("%s(t+%d)", name, i)
			}
			nRows := series[j].NRows()
			start, end := i, nRows-1
			copied := series[j].Copy(dataframe.Range{Start: &start, End: &end})
			copied.Rename(name)
			for k := 0; k < i; k++ {
				copied.Append(nil)
			}
			inAndOut = append(inAndOut, copied)
		}
	}
	return dataframe.NewDataFrame(inAndOut...)
}

func dropNaN(df *dataframe.DataFrame) *dataframe.DataFrame {
	iterator := df.ValuesIterator(dataframe.ValuesOptions{InitialRow: 0, Step: 1, DontReadLock: true})
	row, values, _ := iterator(dataframe.SeriesName)
	for ; row != nil; row, values, _ = iterator(dataframe.SeriesName) {
		for _, value := range values {
			if value == nil {
				df.Remove(*row)
			}
		}
	}
	return df
}

func trainTestSplit(ctx context.Context, df *dataframe.DataFrame, trainSize float64) (train, test *dataframe.DataFrame, err error) {
	split := int(math.Ceil(float64(df.NRows()) * 0.8))
	trainFilterFn := dataframe.FilterDataFrameFn(func(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error) {
		if row < split {
			return dataframe.KEEP, nil
		}
		return dataframe.DROP, nil
	})
	testFilterFn := dataframe.FilterDataFrameFn(func(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error) {
		if row >= split {
			return dataframe.KEEP, nil
		}
		return dataframe.DROP, nil
	})

	trainFiltered, err := dataframe.Filter(ctx, df, trainFilterFn)
	if err != nil {
		return nil, nil, err
	}
	testFiltered, err := dataframe.Filter(ctx, df, testFilterFn)
	if err != nil {
		return nil, nil, err
	}

	return trainFiltered.(*dataframe.DataFrame), testFiltered.(*dataframe.DataFrame), nil
}

func exportCSV(ctx context.Context, path string, df *dataframe.DataFrame, start, end int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := exports.ExportToCSV(ctx, file, df, exports.CSVExportOptions{
		Range:     dataframe.Range{Start: &start, End: &end},
		Separator: ',',
	}); err != nil {
		return err
	}
	return file.Close()
}
