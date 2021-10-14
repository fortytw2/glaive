package glaive

import (
	"context"
	"math/rand"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/fortytw2/glaive/buildinfo"
	"github.com/fortytw2/glaive/run"
	"github.com/fortytw2/lounge"
)

type Glaive struct {
	Log           lounge.Log
	Group         run.Group
	ListenAddress string
}

func (g *Glaive) AddRouter(h http.Handler) {
	ln, err := net.Listen("tcp", g.ListenAddress)
	if err != nil {
		g.Log.Errorf("failed listen on address %s: %s", g.ListenAddress, err)
		return
	}

	g.Group.Add(func() error {
		g.Log.Infof("boot sequence complete, listening on address: %s", g.ListenAddress)
		return http.Serve(ln, h)
	}, func(error) {
		ln.Close()
	})
}

// New initializes, logs default info, and returns a preset run.Group with exit conditions
func New() *Glaive {
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())

	log := lounge.NewDefaultLog()
	if !buildinfo.IsRelease() {
		log = lounge.NewDefaultLog(lounge.WithDebugEnabled())
		log.Infof("development build, for testing and development only")
	}

	totalCPUs, maxProcs := setGOMAXPROCS(log)

	tzName, offset := time.Now().Local().Zone()

	log.Infof("%s boot sequence initiated", buildinfo.ProjectName)
	log.Infof("timezone: %s %d", tzName, offset)

	if tzName != time.UTC.String() {
		log.Errorf("using system local timezone, not UTC. set TZ=UTC in the environment")
	}
	log.Infof("release version: %s", buildinfo.DisplayVersion)
	log.Infof("glaive version: %s", buildinfo.GlaiveVersion)
	log.Infof("go version: %s %d maxprocs", runtime.Version(), maxProcs)
	log.Infof("internal version: %s", buildinfo.GitVersion[0:12])
	log.Infof("platform: %s_%s %d cpus", runtime.GOOS, runtime.GOARCH, totalCPUs)

	printPlatformInformation(log)

	g := run.Group{}

	g.Add(run.SignalHandler(ctx, os.Interrupt))

	return &Glaive{
		Log:           log,
		ListenAddress: "127.0.0.1:3000",
		Group:         g,
	}
}

func setGOMAXPROCS(log lounge.Log) (int, int) {
	numCPU := runtime.NumCPU()
	maxProcsSetting := 0
	if numCPU > 8 {
		log.Infof("glaive: clamping GOMAXPROCS to 8")
		// clamp GOMAXPROCS so nothing looks weird
		runtime.GOMAXPROCS(8)
		maxProcsSetting = 8
	} else {
		// prevent this from being overridden by end users
		runtime.GOMAXPROCS(numCPU)
		maxProcsSetting = numCPU
	}

	return numCPU, maxProcsSetting
}
