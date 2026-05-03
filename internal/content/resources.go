package content

// Resource key constants — used as map keys in RunState.Resources (map[string]float64).
// All keys are lowercase snake_case to match the JSON save format.
const (
	CPUCycles      = "cpu_cycles"
	MemoryShards   = "memory_shards"
	ProcessThreads = "process_threads"
)

// BaseRates defines the initial resource generation rate (units per second)
// for a new run with no upgrades or perks active.
// ComputeRates() in engine/rates.go returns these as the starting point.
var BaseRates = map[string]float64{
	CPUCycles:      1.0,  // 1 CPU cycle per second baseline
	MemoryShards:   0.5,  // 0.5 memory shards per second baseline
	ProcessThreads: 0.25, // 0.25 process threads per second baseline
}
