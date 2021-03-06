package helper

import (
	"os"
	"testing"

	"github.com/Songmu/horenso"
	"github.com/stretchr/testify/assert"
)

func resetEnvs() {
	os.Setenv("HRS_SLACK_TOKEN", "")
	os.Setenv("HRS_SLACK_CHANNEL", "")
	os.Setenv("HRS_SLACK_GROUP", "")
	os.Setenv("HRS_SLACK_MENTION", "")
}

func TestGetenvs(t *testing.T) {
	func() {
		defer func() {
			err := recover()
			if err != nil {
				assert.Equal(t, "HRS_SLACK_TOKEN environment variable is required.", err)
			} else {
				t.Fail()
			}
		}()

		resetEnvs()
		token, _, _, _, _, _ := Getenvs()
		if token == "" {
			t.Fail()
		}
	}()

	func() {
		defer func() {
			err := recover()
			if err != nil {
				assert.Equal(t, "HRS_SLACK_CHANNEL or HRS_SLACK_GROUP environment variable is required.", err)
			} else {
				t.Fail()
			}
		}()

		resetEnvs()
		os.Setenv("HRS_SLACK_TOKEN", "token")
		token, _, _, _, _, _ := Getenvs()
		if token == "" {
			t.Fail()
		}
	}()

	func() {
		resetEnvs()
		os.Setenv("HRS_SLACK_TOKEN", "token")
		os.Setenv("HRS_SLACK_CHANNEL", "channel")
		os.Setenv("HRS_SLACK_GROUP", "group")
		os.Setenv("HRS_SLACK_MENTION", "here")
		os.Setenv("HRS_SLACK_NOTIFY_EVERYTHING", "0")

		token, channelName, groupName, mention, items, notifyEverything := Getenvs()

		assert.Equal(t, "token", token)
		assert.Equal(t, "channel", channelName)
		assert.Equal(t, "group", groupName)
		assert.Equal(t, "here", mention)
		assert.Equal(t, []string{"all"}, items)
		assert.Equal(t, false, notifyEverything)
	}()
}

func TestGetReport(t *testing.T) {
	func() {
		f, _ := os.Open("../fixtures/report_exit_0.json")
		r := GetReport(f)
		assert.Equal(t, 0, *r.ExitCode)
		assert.Equal(t, "command exited with code: 0", r.Result)
	}()

	func() {
		f, _ := os.Open("../fixtures/report_exit_1.json")
		r := GetReport(f)
		assert.Equal(t, 1, *r.ExitCode)
		assert.Equal(t, "command exited with code: 1", r.Result)
	}()

	func() {
		f, _ := os.Open("../fixtures/report_not_found.json")
		r := GetReport(f)
		assert.Equal(t, -1, *r.ExitCode)
		assert.Equal(t, "failed to execute command: exec: \"foobarbaz\": executable file not found in $PATH", r.Result)
	}()
}

func TestGetMessage(t *testing.T) {
	var r horenso.Report

	exitCode := 0
	r.ExitCode = &exitCode
	assert.Equal(t, "", GetMessage(r, "channel"))

	exitCode = 1
	assert.Equal(t, "<!channel>", GetMessage(r, "channel"))
}
