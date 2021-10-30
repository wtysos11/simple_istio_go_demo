package main

// calculate execution time of fft
// need to install: go get github.com/mjibson/go-dsp
// source of code: https://pkg.go.dev/github.com/mjibson/go-dsp/fft#FFT

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"github.com/kataras/iris/v12"
	"math"
	"strconv"
	"time"
	prometheusMiddleware "github.com/iris-contrib/middleware/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func calculate_fft(numSamples int, freq float64) time.Duration{
	start := time.Now()
	// Equation 3-10.
	x := func(n int, freq float64) float64 {
		wave0 := 10 * math.Sin(2.0 * math.Pi * float64(n) * freq / 8.0)
		wave1 := 0.5 * math.Sin(2*math.Pi*float64(n) * freq /4.0+3.0*math.Pi/4.0)
		return wave0 + wave1
	}

	// Discretize our function by sampling at 8 points.
	a := make([]float64, numSamples)
	for i := 0; i < numSamples; i++ {
		a[i] = x(i, freq)
	}

	_ = fft.FFTReal(a)
	return time.Since(start)
}

func main(){
	app := iris.New()
	m := prometheusMiddleware.New("serviceName", 0.001,0.002,0.005,0.01,0.02,0.05,0.1,0.2,0.5)
	app.Use(m.ServeHTTP)
	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		m.ServeHTTP(ctx)
		ctx.Writef("Not Found")
	})
	app.Get("/metrics", iris.FromStd(promhttp.Handler()))

	app.Post("/fft/real/{numSamples}", func(ctx iris.Context) {
		fmt.Println("fft_real")
		numSampleStr := ctx.Params().Get("numSamples")
		numSample, err := strconv.ParseInt(numSampleStr,10, 64)
		if err != nil{
			ctx.StatusCode(iris.StatusInternalServerError)
		}
		ctx.Writef("Length is %d", numSample)
		elapseTime := calculate_fft(int(numSample),0.01)
		fmt.Println(elapseTime)
	})
	app.Post("/fft", func(ctx iris.Context) {
		fmt.Println("fft")
		numSamplesStr := ctx.URLParamDefault("numSamples", "10000")
		freqStr := ctx.URLParamDefault("freq","0.1") // shortcut for ctx.Request().URL.Query().Get("lastname")


		numSample, err := strconv.ParseInt(numSamplesStr,10, 64)
		if err != nil{
			ctx.StatusCode(302)
			return
		}
		freq, err := strconv.ParseFloat(freqStr, 32)
		if err != nil{
			ctx.StatusCode(302)
			return
		}
		elapseTime := calculate_fft(int(numSample),freq)
		fmt.Println(numSample, freq)
		fmt.Println(elapseTime)

		ctx.Writef("Hello %s %s", numSamplesStr, freqStr)
	})


	app.Listen(":8080")
}
