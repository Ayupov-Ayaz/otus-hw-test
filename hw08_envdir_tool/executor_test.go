package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	helloKey = "HELLO"
	fooKey   = "FOO"
	unsetKey = "UNSET"
	emptyKey = "EMPTY"
	userKey  = "USER"
	barKey   = "BAR"
	addedKey = "ADDED"
	//
	bar   = "bar"
	hello = "hello"
	foo   = "foo\nwith new line"
	added = "added"
	user  = "tommy"
)

func Test_getOsEnvironments(t *testing.T) {
	resp := getOsEnvironments()
	require.NotNil(t, resp)
	const notContains = "="

	for k, v := range resp {
		require.NotContains(t, k, notContains)
		require.NotContains(t, v, notContains)
	}
}

func Test_joinEnvironments(t *testing.T) {
	osEnvs := map[string]string{
		helloKey: "SHOULD_REPLACE",
		fooKey:   "SHOULD_REPLACE",
		unsetKey: "SHOULD_REPLACE",
		emptyKey: "SHOULD_BE_EMPTY",
		userKey:  user,
	}

	appEnvs := Environment{
		helloKey: NewEnv(hello),
		barKey:   NewEnv(bar),
		fooKey:   NewEnv(foo),
		unsetKey: NewEnv(""),
		addedKey: NewEnv(added),
		emptyKey: NewEnv(""),
	}

	joinEnvironments(osEnvs, appEnvs)
	require.Equal(t, 5, len(osEnvs))

	require.Equal(t, osEnvs[helloKey], hello)
	require.Equal(t, osEnvs[barKey], bar)
	require.Equal(t, osEnvs[userKey], user)
	require.Equal(t, osEnvs[addedKey], added)
	require.Equal(t, osEnvs[fooKey], foo)

	removedKeys := []string{
		unsetKey, emptyKey,
	}

	for _, key := range removedKeys {
		_, ok := osEnvs[key]
		require.False(t, ok)
	}
}

func Test_prepareEnvironments(t *testing.T) {
	cmdEnvs := map[string]string{
		userKey:  user,
		fooKey:   foo,
		helloKey: hello,
	}

	resp := prepareEnvironments(cmdEnvs)
	require.Equal(t, 3, len(resp))

	expResp := []string{
		userKey + "=" + user,
		fooKey + "=" + foo,
		helloKey + "=" + hello,
	}

	for _, exp := range expResp {
		var ok bool
		for _, got := range resp {
			if exp == got {
				ok = true
			}
		}

		require.True(t, ok)
	}
}
