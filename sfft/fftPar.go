package sfft

import "sync"

// FFT2Par is a parallel version of the FFT2
type FFT2Par struct {
	Transformers []*FFT2
}

// NewFFT2Par returns a new instance of the parallel FFT2. nr is the number of rows
// nc is the number of columns and nWork is the number of workers used to perform
// the FFTs. Note that both the number of rows and the number of columns has to be
// divisible by the number of workers
func NewFFT2Par(nr, nc, nWork int) *FFT2Par {
	if nr%nWork != 0 || nc%nWork != 0 {
		panic("fftpar: The number of rows and columns has to be divisible by the number of workers")
	}
	var ftPar FFT2Par
	ftPar.Transformers = make([]*FFT2, nWork)
	for i := 0; i < nWork; i++ {
		ftPar.Transformers[i] = NewFFT2(nr, nc)

		// Split rows and cols among the workers
		rowsPerWorker := nr / nWork
		colsPerWorkers := nc / nWork
		ftPar.Transformers[i].rows = ftPar.Transformers[i].rows[i*rowsPerWorker : (i+1)*rowsPerWorker]
		ftPar.Transformers[i].cols = ftPar.Transformers[i].cols[i*colsPerWorkers : (i+1)*colsPerWorkers]
	}
	return &ftPar
}

// FFT performs forward FFT
func (f *FFT2Par) FFT(data []complex128) []complex128 {
	var wg sync.WaitGroup
	for i := range f.Transformers {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transformers[num].RowTransform(data, f.Transformers[num].ftRow.Coefficients)
		}(i)
	}
	wg.Wait()

	for i := range f.Transformers {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transformers[num].ColTransform(data, f.Transformers[num].ftCol.Coefficients)
		}(i)
	}
	wg.Wait()
	return data
}

// IFFT performs backward FFT
func (f *FFT2Par) IFFT(data []complex128) []complex128 {
	var wg sync.WaitGroup
	for i := range f.Transformers {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transformers[num].RowTransform(data, f.Transformers[num].ftRow.Sequence)
		}(i)
	}
	wg.Wait()

	for i := range f.Transformers {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transformers[num].ColTransform(data, f.Transformers[num].ftCol.Sequence)
		}(i)
	}
	wg.Wait()
	return data
}

// Freq returns the frequency corresponding to index i in the array returned by FFT
func (f *FFT2Par) Freq(i int) []float64 {
	return f.Transformers[0].Freq(i)
}

// FFT3Par is type used to perform three dimensional FFTs with multiple workers
type FFT3Par struct {
	Transforms []*FFT3
}

// NewFFT3Par returns a new instance of FFT3Par. nr is the number of rows, nc is
// the number of columns, nd is the number of "planes" and nWorkers is the number
// of workers. Note that both the number of rows and the number of columns has to
// be divisible by the number of workers
func NewFFT3Par(nr, nc, nd, nWorkers int) *FFT3Par {
	if nr%nWorkers != 0 || nc%nWorkers != 0 {
		panic("fftpar: The number of rows and the number of columns must be divisible by the number of workers")
	}
	var ft FFT3Par
	ft.Transforms = make([]*FFT3, nWorkers)
	for i := 0; i < nWorkers; i++ {
		ft.Transforms[i] = NewFFT3(nr, nc, nd)
		rowsPerWorker := nr / nWorkers
		colsPerWorker := nc / nWorkers
		ft.Transforms[i].rows = ft.Transforms[i].rows[i*rowsPerWorker : (i+1)*rowsPerWorker]
		ft.Transforms[i].cols = ft.Transforms[i].cols[i*colsPerWorker : (i+1)*colsPerWorker]
	}
	return &ft
}

// FFT performs forward fourier transform
func (f *FFT3Par) FFT(data []complex128) []complex128 {
	var wg sync.WaitGroup
	for i := range f.Transforms {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transforms[num].RowTransform(data, f.Transforms[num].row.Coefficients)
		}(i)
	}
	wg.Wait()

	for i := range f.Transforms {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transforms[num].ColTransform(data, f.Transforms[num].col.Coefficients)
		}(i)
	}
	wg.Wait()

	for i := range f.Transforms {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transforms[num].DepthTransform(data, f.Transforms[num].depth.Coefficients)
		}(i)
	}
	wg.Wait()
	return data
}

// IFFT performs backward fourier transform
func (f *FFT3Par) IFFT(data []complex128) []complex128 {
	var wg sync.WaitGroup
	for i := range f.Transforms {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transforms[num].RowTransform(data, f.Transforms[num].row.Sequence)
		}(i)
	}
	wg.Wait()

	for i := range f.Transforms {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transforms[num].ColTransform(data, f.Transforms[num].col.Sequence)
		}(i)
	}
	wg.Wait()

	for i := range f.Transforms {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			f.Transforms[num].DepthTransform(data, f.Transforms[num].depth.Sequence)
		}(i)
	}
	wg.Wait()
	return data
}

// Freq returns the frequency corresponding to the i-th item in the array returned
// by FFT
func (f *FFT3Par) Freq(i int) []float64 {
	return f.Transforms[0].Freq(i)
}
