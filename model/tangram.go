package model

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/DataDog/go-python3"
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

func (t Tangram) Train(ctx context.Context, module *python3.PyObject, data prom.Metric, duration time.Duration) (Forecasts, error) {
	samples := dataframe.NewSeriesFloat64("samples", nil)
	timeSeries := dataframe.NewDataFrame(samples)
	iteration := int(duration.Minutes())
	wantedFrequency := time.Second

	var (
		prevTime    time.Time
		granularity time.Duration
	)
	for i, value := range data.Values {
		if i == 0 {
			prevTime = data.Values[0].Time
		}
		if i == 1 {
			granularity = value.Time.Sub(prevTime)
		}
		if i > 0 && value.Time.Sub(prevTime) < wantedFrequency {
			continue
		}
		prevTime = value.Time
		samples.Append(value.Value)
	}

	executable, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(executable)
	trainingPath := path.Join(dir, "train.csv")
	// testingPath := path.Join(dir, "test.csv")
	tangramPath := path.Join(dir, "model.tangram")
	configPath := path.Join(dir, "tangram.json")

	numX, numY := 6, 1
	var columnNames []string
	for i := numX; i > 0; i-- {
		columnNames = append(columnNames, fmt.Sprintf("samples(t-%d)", i))
	}

	t.logger.Info("start model training", zap.String("training_data", trainingPath), zap.String("tangram_model", tangramPath))
	start := time.Now()

	supervised := timeSeriesToSupervised(timeSeries, numX, numY)
	nanDropped := dropNaN(supervised)

	// trainX, testX, err := trainTestSplit(ctx, nanDropped, 0.8)
	// if err != nil {
	// 	return nil, err
	// }
	// if err := exportCSV(ctx, trainingPath, trainX, 0, -1); err != nil {
	// 	return nil, err
	// }
	// if err := exportCSV(ctx, testingPath, testX, 0, -2); err != nil {
	// 	return nil, err
	// }

	if err := exportCSV(ctx, trainingPath, nanDropped, 0, -1); err != nil {
		return nil, err
	}

	cmdName := "tangram"
	cmdOpts := []string{
		"train",
		"--no-progress",
		// "--file-train", trainingPath,
		// "--file-test", testingPath,
		"--file", trainingPath,
		"--config", configPath,
		"--output", tangramPath,
		"--target", "samples(t)",
	}
	if err := exec.Command(cmdName, cmdOpts...).Run(); err != nil {
		return nil, err
	}

	// testingRow := testX.Row(testX.NRows()-1, true, dataframe.SeriesName)
	testingRow := nanDropped.Row(nanDropped.NRows()-1, true, dataframe.SeriesName)
	predictInput := make(tangram.PredictInput)
	for _, name := range columnNames {
		predictInput[name] = testingRow[name]
	}

	model, err := tangram.LoadModelFromPath(tangramPath, nil)
	if err != nil {
		return nil, err
	}
	defer model.Destroy()

	iteration = len(data.Values)

	var (
		predicted []float64
		forecasts Forecasts
	)
	timestamp := prevTime.Unix()
	gapInBetween := int64(granularity.Seconds())
	for i := 0; i < iteration; i++ {
		output := model.PredictOne(predictInput, nil)
		pred, ok := output.(tangram.RegressionPredictOutput)
		if !ok {
			return nil, fmt.Errorf("predict output %T is not a RegressionPredictOutput", output)
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

		timestamp += int64(i+1) * gapInBetween
		stdDevX2 := stat.StdDev(predicted[:i], nil) * 2
		forecasts = append(forecasts, Forecast{
			Timestamp: timestamp,
			Values: map[Column]float64{
				Yhat:      current,
				YhatUpper: current + stdDevX2,
				YhatLower: current - stdDevX2,
			},
		})
	}

	t.logger.Info("training completed",
		zap.Stringer("duration", time.Since(start)), zap.Int("predictions", len(predicted)))

	return forecasts, nil
}
