package content

// EnemyDefinition describes a daemon enemy type that appears on tower floors.
type EnemyDefinition struct {
	ID          string
	Name        string
	Type        string
	Description string
	FloorMin    int
	FloorMax    int // 0 = no upper limit
}

// Enemies is the ordered list of daemon enemy tiers by floor range.
var Enemies = []EnemyDefinition{
	{
		ID:          "ghost_process",
		Name:        "Ghost Process",
		Type:        "ORPHAN",
		Description: "A dangling reference still executing after its parent was killed.",
		FloorMin:    1,
		FloorMax:    3,
	},
	{
		ID:          "zombie_daemon",
		Name:        "Zombie Daemon",
		Type:        "UNDEAD",
		Description: "Terminated but unreleased from the process table. Waiting for reaping.",
		FloorMin:    4,
		FloorMax:    6,
	},
	{
		ID:          "race_condition",
		Name:        "Race Condition",
		Type:        "CONCURRENT",
		Description: "Two threads competing for the same lock. One will survive.",
		FloorMin:    7,
		FloorMax:    10,
	},
	{
		ID:          "memory_leak",
		Name:        "Memory Leak",
		Type:        "ENTROPIC",
		Description: "Slow. Growing. Consuming resources without bound.",
		FloorMin:    11,
		FloorMax:    15,
	},
	{
		ID:          "kernel_panic",
		Name:        "Kernel Panic",
		Type:        "FATAL",
		Description: "Unrecoverable error state given malevolent form.",
		FloorMin:    16,
		FloorMax:    0,
	},
}

// EnemiesForFloor returns all enemy definitions active on the given floor.
func EnemiesForFloor(floor int) []EnemyDefinition {
	var result []EnemyDefinition
	for _, e := range Enemies {
		if floor >= e.FloorMin && (e.FloorMax == 0 || floor <= e.FloorMax) {
			result = append(result, e)
		}
	}
	return result
}
