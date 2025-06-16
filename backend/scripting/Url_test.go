package scripting

import (
	lua "github.com/yuin/gopher-lua"
	"net/url"
	"testing"
)

func SetupUrlTests() *lua.LState {
	l := lua.NewState()
	registerUrlType(l)
	return l
}

func TestUrlToLuaUserData(t *testing.T) {
	l := SetupUrlTests()
	defer l.Close()

	// Create a sample url.URL
	u, err := url.Parse("https://user:pass@example.com:8080/path?foo=bar&baz=qux#fragment")
	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}
	ud := UrlToLuaUserData(l, u)
	l.SetGlobal("testUrl", ud)

	// Test Lua access to url fields
	err = l.DoString(`
		assert(testUrl.scheme == "https")
		assert(testUrl.host == "example.com:8080")
		assert(testUrl.path == "/path")
		assert(testUrl.fragment == "fragment")
		assert(testUrl.rawQuery == "foo=bar&baz=qux")
		assert(testUrl.user.username == "user")
		assert(testUrl.user.password == "pass")
		assert(testUrl.query.foo == "bar")
		assert(testUrl.query.baz == "qux")
	`)
	if err != nil {
		t.Errorf("Lua test failed: %v", err)
	}
}
