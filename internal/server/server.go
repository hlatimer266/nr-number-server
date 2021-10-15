package server

import (
	"context"
	"os"
	"sync"

	"github.com/hlatimer266/nr-number-server/internal/cache"
	"github.com/hlatimer266/nr-number-server/internal/connection"
	"github.com/hlatimer266/nr-number-server/internal/status"
	"github.com/hlatimer266/nr-number-server/internal/write"
)

func init() {
	os.Remove("number.log")
}

func Run(p string, ctx context.Context) {
	// create new listener and cache
	l := connection.NewListener(p)
	l.NumCache = cache.NewNumberCache()

	clientCtx, clientCancel := context.WithCancel(ctx)
	writeCtx, writeCancel := context.WithCancel(ctx)
	reportCtx, reportCancel := context.WithCancel(ctx)

	l.Ctx = clientCtx
	l.StopListener = clientCancel
	l.StopWriter = writeCancel
	l.StopReport = reportCancel

	defer l.Listen.Close()
	go l.WaitCtxFinish() // listen for context to finish

	l.AllWG = sync.WaitGroup{}
	l.AllWG.Add(3)
	go func() {
		err := l.ManageClient() // listen for client connetions
		if err != nil {
			return
		}
		l.AllWG.Done()
	}()

	go func() {
		write.Latest(writeCtx, l.NumCache) // start batch write to file
		l.AllWG.Done()
	}()

	go func() {
		status.ReportLatest(reportCtx, l.NumCache) // start status report on a different non-blocking thread
		l.AllWG.Done()
	}()

	l.AllWG.Wait()
	l.ClientWG.Wait()

}
