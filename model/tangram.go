package model

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/rocketlaunchr/dataframe-go"
	"github.com/tangramdotdev/tangram-go"
	"github.com/waltzofpearls/reckon/prom"
	"go.uber.org/zap"
	"gonum.org/v1/gonum/stat"
)

type Tangram struct {
	logger *zap.Logger
}

func NewTangram(lg *zap.Logger) Tangram {
	return Tangram{
		logger: lg,
	}
}

func (t Tangram) Train(ctx context.Context, data prom.Metric, duration time.Duration) (Forecasts, error) {
	tempDir, err := ioutil.TempDir("", "tangram")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create temp dir")
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	executable, err := os.Executable()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the current executable path")
	}
	execDir := filepath.Dir(executable)

	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a new v4 uuid")
	}
	configPath := filepath.Join(execDir, "model", "tangram", "config.json")
	trainingPath := filepath.Join(tempDir, uid.String()+".csv")
	tangramPath := filepath.Join(tempDir, uid.String()+".tangram")

	logger := t.logger.With(zap.String("metric_name", data.Name), zap.Any("metric_labels", data.Labels),
		zap.String("config", configPath), zap.String("training_data", trainingPath),
		zap.String("model_file", tangramPath), zap.Int("data_length", len(data.Values)),
		zap.String("want_duration", duration.String()))
	logger.Info("train model with data")

	samples := dataframe.NewSeriesFloat64("samples", nil)
	timeSeries := dataframe.NewDataFrame(samples)
	granularity := time.Minute
	iteration := int(duration.Minutes())

	var (
		prevTime   time.Time
		lastTime   time.Time
		stdDevVals []float64
	)
	for i, value := range data.Values {
		if i == 0 {
			prevTime = value.Time
		}
		if i > 0 && value.Time.Sub(prevTime) < granularity {
			continue
		}
		prevTime = value.Time
		samples.Append(value.Value)
		stdDevVals = append(stdDevVals, value.Value)
		lastTime = value.Time
	}

	numX, numY := 6, 1
	var columnNames []string
	for i := numX; i > 0; i-- {
		columnNames = append(columnNames, fmt.Sprintf("samples(t-%d)", i))
	}

	logger.Info("start tangram training",
		zap.String("training_data", trainingPath), zap.String("tangram_model", tangramPath))
	start := time.Now()

	supervised := timeSeriesToSupervised(timeSeries, numX, numY)
	nanDropped := dropNaN(supervised)

	if err := exportCSV(ctx, trainingPath, nanDropped, 0, -1); err != nil {
		return nil, errors.Wrap(err, "failed to export training data to csv file")
	}

	cmdName := "tangram"
	cmdOpts := []string{
		"train", "--no-progress",
		"--file", trainingPath,
		"--config", configPath,
		"--output", tangramPath,
		"--target", "samples(t)",
	}
	if err := exec.Command(cmdName, cmdOpts...).Run(); err != nil {
		return nil, errors.Wrap(err, "failed to execute tangram train command")
	}

	testingRow := nanDropped.Row(nanDropped.NRows()-1, true, dataframe.SeriesName)
	predictInput := make(tangram.PredictInput)
	for _, name := range columnNames {
		predictInput[name] = testingRow[name]
	}

	model, err := tangram.LoadModelFromPath(tangramPath, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load tangram model")
	}
	defer model.Destroy()

	var (
		predicted []float64
		forecasts Forecasts
	)
	timestamp := lastTime.Unix()
	gapInBetween := int64(granularity.Seconds())
	stdDevValsLen := len(stdDevVals)
	for i := 0; i < iteration; i++ {
		output := model.PredictOne(predictInput, nil)
		pred, ok := output.(tangram.RegressionPredictOutput)
		if !ok {
			return nil, errors.Errorf("predict output %T is not a RegressionPredictOutput", output)
		}

		current := float64(pred.Value)
		length := len(predicted)
		if length > 2 {
			secondLast := predicted[length-2]
			thirdLast := predicted[length-3]
			if secondLast == current {
				if secondLast > thirdLast {
					current += 0.000001
				} else {
					current -= 0.000001
				}
			}
		}

		for j := numX; j > 1; j-- {
			predictInput[fmt.Sprintf("samples(t-%d)", j)] = predictInput[fmt.Sprintf("samples(t-%d)", j-1)]
		}
		predictInput["samples(t-1)"] = current
		predicted = append(predicted, current)
		stdDevVals = append(stdDevVals, current)

		timestamp += gapInBetween
		stdDevX2 := stat.StdDev(stdDevVals[:stdDevValsLen+i], nil) * 2
		forecasts = append(forecasts, Forecast{
			Timestamp: timestamp,
			Values: map[Column]float64{
				Yhat:      current,
				YhatUpper: current + stdDevX2,
				YhatLower: current - stdDevX2,
			},
		})
	}

	logger.Info("tangram training completed",
		zap.Stringer("elasped", time.Since(start)), zap.Int("forecasts", len(forecasts)))

	return forecasts, nil
}
