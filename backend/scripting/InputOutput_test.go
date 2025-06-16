package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"os"
	"testing"
)

func SetupInputOutputTests() *lua.LState {
	L := lua.NewState()
	LuaExtendIoTable(L)
	return L
}

func TestLuaGetFilesInDirectory(t *testing.T) {
	L := SetupInputOutputTests()
	defer L.Close()

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testFile1 := tempDir + "/test1.txt"
	testFile2 := tempDir + "/test2.txt"
	testFile3 := tempDir + "/test3.log"

	// Create test files
	if err := os.WriteFile(testFile1, []byte("Test file 1"), 0644); err != nil {
		t.Fatalf("Failed to create test file 1: %v", err)
	}
	if err := os.WriteFile(testFile2, []byte("Test file 2"), 0644); err != nil {
		t.Fatalf("Failed to create test file 2: %v", err)
	}
	if err := os.WriteFile(testFile3, []byte("Test file 3"), 0644); err != nil {
		t.Fatalf("Failed to create test file 3: %v", err)
	}
	println("Test files created in:", tempDir)

	script := `
		local files = io.getFilesInDirectory("` + tempDir + `", {"*.txt", "*.log"})
		assert(#files == 3, "Expected 3 files, got " .. #files)
		assert(files[1] == "` + testFile1 + `" or files[2] == "` + testFile1 + `" or files[3] == "` + testFile1 + `", "Missing ` + testFile1 + `")
        assert(files[1] == "` + testFile2 + `" or files[2] == "` + testFile2 + `" or files[3] == "` + testFile2 + `", "Missing ` + testFile2 + `")
		assert(files[1] == "` + testFile3 + `" or files[2] == "` + testFile3 + `" or files[3] == "` + testFile3 + `", "Missing ` + testFile3 + `")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua getFilesInDirectory failed: %v", err)
	}
}

func TestLuaGetFileName(t *testing.T) {
	L := SetupInputOutputTests()
	defer L.Close()

	// Create a temporary file for testing
	tempFile := t.TempDir() + "/testfile.txt"
	if err := os.WriteFile(tempFile, []byte("Test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	script := `
		local fileName = io.getFileName("` + tempFile + `")
		assert(fileName == "testfile.txt", "Expected 'testfile.txt', got " .. fileName)
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua getFileName failed: %v", err)
	}
}

func TestLuaGetFileContent(t *testing.T) {
	L := SetupInputOutputTests()
	defer L.Close()

	// Create a temporary file for testing
	tempFile := t.TempDir() + "/testfile.txt"
	content := "This is a test file."
	if err := os.WriteFile(tempFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	script := `
		local fileContent = io.getFileContent("` + tempFile + `")
        assert(fileContent ~= nil, "File content should not be nil")
		assert(fileContent == "` + content + `", "Expected '` + content + `', got '" .. fileContent .. "'")
	`

	if err := L.DoString(script); err != nil {
		t.Fatalf("Lua getFileContent failed: %v", err)
	}
}
