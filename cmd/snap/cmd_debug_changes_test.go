// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2019 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main_test

import (
	"io/ioutil"
	"path/filepath"

	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/cmd/snap"
)

var stateJSON = []byte(`
{
	"last-task-id": 999,
	"last-change-id": 99,

	"data": {
		"snaps": {}
	},
	"changes": {
		"1": {
			"id": "1",
			"kind": "install-snap",
			"summary": "install a snap",
			"status": 0,
			"data": {"snap-names": ["a"]},
			"task-ids": ["11","12"]
		},
		"2": {
			"id": "2",
			"kind": "revert-snap",
			"summary": "revert c snap",
			"status": 0,
			"data": {"snap-names": ["c"]},
			"task-ids": ["21","31"]
		}
	},
	"tasks": {
		"11": {
				"id": "11",
				"change": "1",
				"kind": "download-snap",
				"summary": "Download snap a from channel edge",
				"status": 4,
				"data": {"snap-setup": {
						"channel": "edge",
						"flags": 1
				}},
				"halt-tasks": ["12"]
		},
		"12": {"id": "12", "change": "1", "kind": "some-other-task"},
		"21": {
				"id": "21",
				"change": "2",
				"kind": "download-snap",
				"summary": "Download snap b from channel beta",
				"status": 4,
				"data": {"snap-setup": {
						"channel": "beta",
						"flags": 2
				}},
				"halt-tasks": ["12"]
		},
		"31": {
				"id": "31",
				"change": "2",
				"kind": "prepare-snap",
				"summary": "Prepare snap c",
				"status": 4,
				"data": {"snap-setup": {
						"channel": "stable",
						"flags": 1073741828
				}},
				"halt-tasks": ["12"],
				"log": ["logline1", "logline2"]
		}
	}
}
`)

func (s *SnapSuite) TestDebugChanges(c *C) {
	dir := c.MkDir()
	stateFile := filepath.Join(dir, "test-state.json")
	c.Assert(ioutil.WriteFile(stateFile, stateJSON, 0644), IsNil)

	rest, err := main.Parser(main.Client()).ParseArgs([]string{"debug", "changes", stateFile})
	c.Assert(err, IsNil)
	c.Assert(rest, DeepEquals, []string{})
	c.Check(s.Stdout(), Matches,
		"ID   Status  Spawn                 Ready                 Label         Summary\n"+
			"1    Do      0001-01-01T00:00:00Z  0001-01-01T00:00:00Z  install-snap  install a snap\n"+
			"2    Done    0001-01-01T00:00:00Z  0001-01-01T00:00:00Z  revert-snap   revert c snap\n")
	c.Check(s.Stderr(), Equals, "")
}

func (s *SnapSuite) TestDebugChangesMissingState(c *C) {
	_, err := main.Parser(main.Client()).ParseArgs([]string{"debug", "changes", "/missing-state.json"})
	c.Check(err, ErrorMatches, "cannot read the state file: open /missing-state.json: no such file or directory")
}
