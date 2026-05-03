// Package engine implements the resource tick engine and game loop computation
// for afk-x. All functions are pure — no I/O, no global mutable state.
// The only Bubbletea dependency is in tick.go (TickMsg type + TickCmd function),
// which bridges engine timing to the Bubbletea command system.
//
// Dependency rules: engine imports game and content only.
// engine is never imported by game, content, save, or combat.
package engine
