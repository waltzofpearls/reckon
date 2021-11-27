package model

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/exports"
	"github.com/waltzofpearls/reckon/config"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
)

type Trainer interface {
	Train(context.Context, prom.Metric, time.Duration) (Forecasts, error)
}

const TypeProphet = "Prophet"
const TypeTangram = "Tangram"

func New(cf *config.Config, lg *zap.Logger, name string) (Trainer, error) {
	switch name {
	case TypeProphet:
		return NewProphet(cf, lg), nil
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

type Server struct {
	logger *zap.Logger
	signal chan os.Signal
}

func NewServer(lg *zap.Logger, sig chan os.Signal) Server {
	return Server{
		logger: lg,
		signal: sig,
	}
}

func (s Server) Start(ctx context.Context) func() error {
	return func() error {
		doneChan := make(chan bool)
	restart:
		s.logger.Info("starting python gRPC server...")
		cancelCtx, cancelFn := context.WithCancel(context.Background())
		subProc := exec.CommandContext(cancelCtx, "python", "model/server/main.py")
		subProc.Stdout = os.Stdout
		subProc.Stderr = os.Stderr
		go func() {
			select {
			case <-ctx.Done():
			case <-doneChan:
			}
			s.logger.Info("stopping python gRPC server...")
			subProc.Process.Signal(<-s.signal)
			time.AfterFunc(30*time.Second, cancelFn)
		}()
		if err := subProc.Run(); err != nil {
			switch err.(type) {
			case *exec.ExitError:
				select {
				case <-ctx.Done():
					s.logger.Info("stopped python gRPC server")
				default:
					doneChan <- true
					time.Sleep(5 * time.Second)
					goto restart
				}
			default:
				s.logger.Error("failed to run python gRPC server", zap.Error(err))
				return err
			}
		}
		return nil
	}
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
