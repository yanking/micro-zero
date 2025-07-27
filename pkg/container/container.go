package container

import (
	"context"
	"errors"
	"fmt"
	"github.com/yanking/micro-zero/pkg/contract"
	"github.com/yanking/micro-zero/pkg/log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var _ contract.IContainer = (*Container)(nil)

type Container struct {
	name string
	// +optional
	shutdownOverallTimeout time.Duration
	components             []contract.Component
	mu                     sync.Mutex
	wg                     sync.WaitGroup
}

func (c *Container) Name() string {
	return c.name
}

// Option defines optional parameters for initializing the application
// structure.
type Option func(*Container)

func WithShutdownOverallTimeout(t time.Duration) Option {
	return func(c *Container) {
		c.shutdownOverallTimeout = t
	}
}

func New(name string, opts ...Option) *Container {
	c := &Container{
		name:       name,
		components: make([]contract.Component, 0),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func (c *Container) Register(cp contract.Component) {
	c.components = append(c.components, cp)
	log.Infof("container: register component: %s", c.Name())
}

func (c *Container) Run() error {
	// Create a context that gets cancelled on OS signal (SIGINT, SIGTERM)
	// or if any critical component's Start method calls the 'appStopCauseError' function.
	appCtx, appStopCauseError := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appStopCauseError() // Important to release resources

	// Start all registered container components concurrently
	for _, comp := range c.components {
		cp := comp // Capture range variable for goroutine
		c.wg.Add(1)
		go func(cp contract.Component) {
			defer c.wg.Done()
			log.Infof("container: starting component: %s", cp.Name())
			if err := cp.Start(appCtx); err != nil { // Pass the container's cancellable context
				log.Errorf("container: error starting component %s: %v. initiating application shutdown.", cp.Name(), err)
				appStopCauseError() // Trigger shutdown for all other components
			}
			// If Start is blocking and succeeds, it will run until appCtx is cancelled.
			// If Start is non-blocking, this goroutine will exit after Start returns.
			// We need to ensure Start methods handle appCtx cancellation properly if they are long-running.
		}(cp)
	}
	log.Infof("container: all components initiated for start. application %s is running.", c.name)

	// Block here until the appCtx is cancelled
	<-appCtx.Done()
	errCause := context.Cause(appCtx)
	if errCause != nil && !errors.Is(errCause, context.Canceled) && !errors.Is(errCause, context.DeadlineExceeded) { // context.Canceled from signal is normal
		log.Infof("container: shutdown initiated due to: %v", errCause)
	} else {
		log.Infof("container: shutdown signal received or context cancelled normally.")
	}

	// --- Graceful Shutdown Procedure ---
	log.Infof("container: initiating graceful stop of application %s...", c.name)

	// Create a new context for the shutdown procedure itself, with a timeout.
	// This timeout is for the *entire* shutdown sequence of all components.
	stopCtx, cancelStopCtx := context.WithTimeout(context.Background(), c.shutdownOverallTimeout)
	defer cancelStopCtx()

	// Stop components in reverse order of registration (LIFO).
	// This assumes a simple dependency order; more complex apps might need a dependency graph.
	for i := len(c.components) - 1; i >= 0; i-- {
		comp := c.components[i]
		log.Infof("container: attempting to stop component: %s", comp.Name())

		// We pass stopCtx to each component's Stop method.
		// The component's Stop method should respect this context's deadline.
		if err := comp.Stop(stopCtx); err != nil {
			log.Errorf("container: error stopping component %s: %v", comp.Name(), err)
		} else {
			log.Infof("container: component %s stopped successfully.", comp.Name())
		}
	}

	// Wait for all goroutines launched by `go func() { c.Start(appCtx) }` to complete.
	// This is crucial because if a component's Start() errors out and calls appStopCauseError(),
	// its goroutine might still be running. We need to ensure all these initial goroutines
	// have finished before declaring the container fully stopped.
	log.Infof("container: waiting for initial component start goroutines to complete...")
	c.wg.Wait()
	log.Infof("container: all initial component start goroutines have completed.")

	log.Infof("container: application %s stopped gracefully.", c.name)

	// Check if shutdown was due to an error from context.Cause or normal signal.
	// signal.NotifyContext cancels with context.Canceled on signal.
	// If it was another error, that might be the one to return from Run().
	if errCause != nil && !errors.Is(errCause, context.Canceled) {
		return fmt.Errorf("application shutdown due to error: %w", errCause)
	}
	return nil // Graceful shutdown completed
}
