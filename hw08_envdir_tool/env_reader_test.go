package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleReadDirFails(t *testing.T) {
	testCases := []struct {
		name      string
		dir       string
		wantError string
	}{
		{
			name:      "Передали_не_существующую_директорию",
			dir:       "/path/to/not_exist/directory",
			wantError: "stat dir: /path/to/not_exist/directory",
		},
		{
			name:      "Передали_не_директорию",
			dir:       "/dev/null",
			wantError: "/dev/null: is not directory",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := ReadDir(tc.dir)
			assert.Contains(t, err.Error(), tc.wantError)
		})
	}
}

func TestEnvironmentReadDir(t *testing.T) {
	type prepare func(dir string, envs *Environment)

	testCases := []struct {
		name     string
		wantEnvs *Environment
		prepare  prepare
	}{
		{
			name: "Позитивный",
			wantEnvs: &Environment{
				"FOO": EnvValue{
					Value:      "BAR",
					NeedRemove: false,
				},
				"BUZZ": EnvValue{
					Value:      "",
					NeedRemove: true, // true потому что в prepare в файл пишем значение как есть
				},
				"CHANGE": EnvValue{
					Value:      "CHANGE",
					NeedRemove: false,
				},
			},
			prepare: func(dir string, envs *Environment) {
				err := os.Mkdir(dir, os.ModePerm)
				if err != nil {
					t.Error(err)
					return
				}

				for key, env := range *envs {
					file, err := os.Create(fmt.Sprintf("%s/%s", dir, key))
					if err != nil {
						t.Error(err)
						return
					}

					_, err = file.WriteString(env.Value)
					if err != nil {
						t.Error(err)
						return
					}
				}
			},
		},
		{
			name: "Тест_обрезки_крайних_пробелов_и_табуляции",
			wantEnvs: &Environment{
				"FOO": EnvValue{
					Value:      "BAR",
					NeedRemove: false,
				},
				"BUZZ": EnvValue{
					Value:      "",
					NeedRemove: false, // true потому что в prepare в файл пишем значение как есть + nbsb + \t
				},
				"CHANGE": EnvValue{
					Value:      "CHANGE",
					NeedRemove: false,
				},
			},
			prepare: func(dir string, envs *Environment) {
				err := os.Mkdir(dir, os.ModePerm)
				if err != nil {
					t.Error(err)
					return
				}

				for key, env := range *envs {
					file, err := os.Create(fmt.Sprintf("%s/%s", dir, key))
					if err != nil {
						t.Error(err)
						return
					}

					_, err = file.WriteString(env.Value + "     \t")
					if err != nil {
						t.Error(err)
						return
					}
				}
			},
		},
		{
			name:     "Тест_пропуска_файла_с_=",
			wantEnvs: &Environment{},
			prepare: func(dir string, envs *Environment) {
				err := os.Mkdir(dir, os.ModePerm)
				if err != nil {
					t.Error(err)
					return
				}

				file, err := os.Create(fmt.Sprintf("%s/%s", dir, "test=item"))
				if err != nil {
					t.Error(err)
					return
				}

				_, err = file.WriteString("     \t")
				if err != nil {
					t.Error(err)
					return
				}
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			i := make([]byte, 1)
			_, err := rand.Read(i)
			assert.NoError(t, err)

			dir := fmt.Sprintf("/tmp/test_%v", i[0])
			tc.prepare(dir, tc.wantEnvs)
			defer os.RemoveAll(dir)

			resultEnvs, err := ReadDir(dir)
			assert.NoError(t, err)

			assert.Equal(t, *tc.wantEnvs, resultEnvs)
		})
	}
}
