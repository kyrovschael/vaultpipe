package env

import (
	"sync"
	"time"
)

// Checkpoint records a named snapshot at a point in time, allowing callers to
// compare the current environment against a previously saved baseline.
type Checkpoint struct {
	mu      sync.RWMutex
	name    string
	snap    Snapshot
	recorded time.Time
}

// NewCheckpoint creates a Checkpoint capturing snap under the given name.
func NewCheckpoint(name string, snap Snapshot) *Checkpoint {
	return &Checkpoint{
		name:     name,
		snap:     snap.Clone(),
		recorded: time.Now(),
	}
}

// Name returns the label given to this checkpoint.
func (c *Checkpoint) Name() string { return c.name }

// RecordedAt returns when the checkpoint was taken.
func (c *Checkpoint) RecordedAt() time.Time { return c.recorded }

// Snapshot returns a copy of the captured snapshot.
func (c *Checkpoint) Snapshot() Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.snap.Clone()
}

// Update replaces the stored snapshot with a new one and refreshes the timestamp.
func (c *Checkpoint) Update(snap Snapshot) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snap = snap.Clone()
	c.recorded = time.Now()
}

// DiffFrom returns the diff between the checkpoint and current, using opts.
func (c *Checkpoint) DiffFrom(current Snapshot, opts DiffOptions) []DiffEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return DiffSnapshots(c.snap, current, opts)
}
