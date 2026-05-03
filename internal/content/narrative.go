package content

// BootText is the flavor text shown during the game's boot sequence (FR34).
// Terminal/computing aesthetic — ancient self-aware mainframe system logs.
const BootText = `SYSTEM INITIALIZE — THE TOWER
> Detecting hardware...          [OK]
> Mounting sector 0x00...        [OK]
> Daemon census: 1024 active processes
> WARNING: Unauthorized access detected on FLOOR 01
> Administrator terminal connected.
>
> Your resources are limited.
> The tower is not.
>
> Ascend.`

// milestoneTexts maps milestone floor numbers to their flavor text (FR35).
// Non-milestone floors are absent from this map — MilestoneText returns "" for them.
var milestoneTexts = map[int]string{
	5: `SECTOR 0x05 CLEARED
The daemon cluster on this floor was running a single process,
unchanged since the tower was seeded. Its last instruction: WAIT.
It has been waiting for 40 years.
You interrupted it.`,

	10: `SECTOR 0x0A CLEARED
You have cleared the first tenth of the lower stack.
The tower grows quieter here. Fewer daemons. More entropy.
The architects built this place to last.
They did not build it to be finished.`,

	20: `SECTOR 0x14 CLEARED
A recursive function. No base case.
It has been calling itself since initialization.
You broke the loop.
Somewhere in memory, a value that was never freed
has finally been collected.`,

	25: `SECTOR 0x19 CLEARED
ADMIN NOTE — ARCHIVED LOG:
"We built the tower to solve a problem we no longer remember.
The solution is still running."`,
}

// MilestoneText returns the flavor text for a milestone floor (FR35).
// Returns "" (empty string) if floor is not a milestone — callers should
// check for non-empty before displaying.
func MilestoneText(floor int) string {
	return milestoneTexts[floor]
}
