package mods

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
	"time"
)

// Helper to create a unique temporary Lua script
func createMockLuaScript(content string) (string, error) {
	tmpFile, err := os.CreateTemp("", "mock_plugin_*.lua")
	if err != nil {
		return "", err
	}
	_, err = tmpFile.Write([]byte(content))
	err = tmpFile.Close()
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), err
}

// Helper to initialize LuaModProvider
func setupLuaModProvider(t *testing.T, luaScript string) *LuaModProvider {
	tmpFile, err := createMockLuaScript(luaScript)
	assert.NoError(t, err)
	t.Cleanup(func() {
		err := os.Remove(tmpFile)
		if err != nil {
			return
		}
	})

	provider, err := NewLuaModProvider(tmpFile)
	assert.NoError(t, err)
	return provider
}

func TestLuaModProvider_GetMods(t *testing.T) {
	luaScript := `
		plugin = {
			GetInstalledMods = function()
				return {
					{ id = "mod1", name = "Test Mod", enabled = true, game_id = "game1" },
					{ id = "mod2", name = "Another Mod", enabled = false, game_id = "game2" }
				}
			end
		}
	`

	provider := setupLuaModProvider(t, luaScript)
	defer provider.Close()

	mods, err := provider.GetMods()
	assert.NoError(t, err)
	assert.Len(t, mods, 2)
	assert.Equal(t, "mod1", mods[0].ID)
	assert.Equal(t, "Test Mod", mods[0].Name)
	assert.True(t, mods[0].Enabled)
	assert.Equal(t, "game1", mods[0].GameID)
}

func TestLuaModProvider_GetModByID(t *testing.T) {
	luaScript := `
		plugin = {
			GetModByID = function(id)
				if id == "mod1" then
					return { id = "mod1", name = "Test Mod", enabled = true, game_id = "game1" }
				end
				return nil
			end
		}
	`

	provider := setupLuaModProvider(t, luaScript)
	defer provider.Close()

	mod, err := provider.GetModByID("mod1")
	assert.NoError(t, err)
	assert.NotNil(t, mod)
	assert.Equal(t, "mod1", mod.ID)
	assert.Equal(t, "Test Mod", mod.Name)
	assert.True(t, mod.Enabled)
	assert.Equal(t, "game1", mod.GameID)

	mod, err = provider.GetModByID("invalid")
	assert.NoError(t, err)
	assert.Nil(t, mod)
}

// Integration test for LuaModProvider
func TestLuaModProvider_Integration(t *testing.T) {
	// Realistic Lua script mimicking plugin behavior
	luaScript := `
		plugin = {
			GetInstalledMods = function()
				return {
					{ id = "mod1", name = "Test Mod", enabled = true, game_id = "game1" },
					{ id = "mod2", name = "Another Mod", enabled = false, game_id = "game2" }
				}
			end,
			GetModByID = function(id)
				if id == "mod1" then
					return { id = "mod1", name = "Test Mod", enabled = true, game_id = "game1" }
				end
				return nil
			end,
			AddMod = function(mod)
				-- Simulate adding a mod (no-op for testing)
				return true
			end,
			RemoveMod = function(id)
				-- Simulate removing a mod (no-op for testing)
				return true
			end
		}
	`

	// Write the Lua script to a temporary file
	tmpFile := "integration_plugin.lua"
	err := os.WriteFile(tmpFile, []byte(luaScript), 0644)
	assert.NoError(t, err)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Logf("Failed to remove temporary file: %v", err)
		}
	}(tmpFile)

	// Initialize LuaModProvider
	provider, err := NewLuaModProvider(tmpFile)
	assert.NoError(t, err)
	defer provider.Close()

	// Test GetInstalledMods
	mods, err := provider.GetMods()
	assert.NoError(t, err)
	assert.Len(t, mods, 2)
	assert.Equal(t, "mod1", mods[0].ID)
	assert.Equal(t, "Test Mod", mods[0].Name)

	// Test GetModByID
	mod, err := provider.GetModByID("mod1")
	assert.NoError(t, err)
	assert.NotNil(t, mod)
	assert.Equal(t, "mod1", mod.ID)

	// Test AddMod
	err = provider.AddMod(Mod{ID: "mod3", Name: "New Mod", Enabled: true, GameID: "game3"})
	assert.NoError(t, err)

	// Test RemoveMod
	err = provider.RemoveMod("mod1")
	assert.NoError(t, err)
}

func TestLuaModProvider_Performance(t *testing.T) {
	luaScript := `
		plugin = {
			GetInstalledMods = function()
				local mods = {}
				for i = 1, 1000 do
					table.insert(mods, { id = "mod" .. i, name = "Mod " .. i, enabled = true, game_id = "game" .. i })
				end
				return mods
			end
		}
	`

	provider := setupLuaModProvider(t, luaScript)
	defer provider.Close()

	// Measure execution time for GetInstalledMods
	start := time.Now()
	mods, err := provider.GetMods()
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.Len(t, mods, 1000)
	t.Logf("GetInstalledMods execution time: %v", duration)

	// Simulate high load with concurrent calls (each with its own provider)
	var wg sync.WaitGroup
	numRoutines := 100
	wg.Add(numRoutines)

	start = time.Now()
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			p := setupLuaModProvider(t, luaScript)
			defer p.Close()
			_, err := p.GetMods()
			assert.NoError(t, err)
		}()
	}
	wg.Wait()
	duration = time.Since(start)
	t.Logf("Concurrent GetInstalledMods execution time: %v", duration)
}
