package app

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// InitLogger возвращает инстанс логгера.
func InitLogger() *log.Logger {
	mainLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	return mainLog
}

// GoAndWaitFnOrCtx запустить функцию и ждать либо завершения функции, либо контекста. В случае завершения контекста
// вызвать функцию done.
func GoAndWaitFnOrCtx(ctx context.Context, fn func() error, done func()) error {
	errCh := make(chan error)

	go func() {
		errCh <- fn()
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		done()

		return ctx.Err()
	}
}

// ErrorOSSignal обертка в виде ошибки над сигналом ОС для того, чтобы использовать SignalNotify в errgroup.
type ErrorOSSignal struct {
	signal os.Signal
}

func (e *ErrorOSSignal) Error() string {
	return e.signal.String()
}

// SignalNotify обработка событий выхода от ОС.
func SignalNotify(ctx context.Context) error {
	c := make(chan os.Signal, 1)

	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-c:
		return &ErrorOSSignal{signal: s}
	case <-ctx.Done():
		return ctx.Err()
	}
}

type RunFunc func(ctx context.Context) error

// RunParallel запускает параллельно функции.
func RunParallel(ctx context.Context, fns ...RunFunc) error {
	group, ctx := errgroup.WithContext(ctx)

	for _, fn := range fns {
		internalFn := fn

		group.Go(func() error {
			return internalFn(ctx)
		})
	}

	return group.Wait()
}

// Exit функция в зависимости от результат errFn выполнить os.Exit.
//
// Если errFn вернет ошибку не типа ErrorOSSignal, то os.Exit(1), иначе os.Exit(0).
func Exit(errFn func() error, l log.Logger) {
	var es *ErrorOSSignal

	err := errFn()

	isSignal := errors.As(err, &es)

	if err != nil && !isSignal {
		l.Println("app is stopped by error, %w", err)
		os.Exit(1)
	}

	l.Println("app is stopped", "signal", es.signal.String())
	os.Exit(0)
}
