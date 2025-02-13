// This file is part of DiceDB.
// Copyright (C) 2024 DiceDB (dicedb.io).
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package websocket

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHINCRBY(t *testing.T) {
	exec := NewWebsocketCommandExecutor()
	conn := exec.ConnectToServer()
	DeleteKey(t, conn, exec, "key")
	DeleteKey(t, conn, exec, "keys")
	DeleteKey(t, conn, exec, "new-key")

	testCases := []struct {
		name   string
		cmds   []string
		expect []interface{}
		delays []time.Duration
	}{

		{
			name:   "HINCRBY on non-existing key",
			cmds:   []string{"HINCRBY key field1 10"},
			expect: []interface{}{float64(10)},
			delays: []time.Duration{0},
		},
		{
			name:   "HINCRBY on existing key",
			cmds:   []string{"HINCRBY key field1 5"},
			expect: []interface{}{float64(15)},
			delays: []time.Duration{0},
		},
		{
			name:   "HINCRBY on non-integer value",
			cmds:   []string{"HSET keys field value", "HINCRBY keys field 1"},
			expect: []interface{}{float64(1), "ERR hash value is not an integer"},
			delays: []time.Duration{0, 0, 0},
		},
		{
			name:   "HINCRBY on non-hashmap key",
			cmds:   []string{"SET key value", "HINCRBY key value 10"},
			expect: []interface{}{"OK", "WRONGTYPE Operation against a key holding the wrong kind of value"},
			delays: []time.Duration{0, 0, 0},
		},
		{
			name:   "HINCRBY overflow",
			cmds:   []string{"HSET new-key value 9000000000000000000", "HINCRBY new-key value 1000000000000000000"},
			expect: []interface{}{float64(1), "ERR increment or decrement would overflow"},
			delays: []time.Duration{0, 0, 0},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			for i, cmd := range tc.cmds {
				if tc.delays[i] > 0 {
					time.Sleep(tc.delays[i])
				}
				result, err := exec.FireCommandAndReadResponse(conn, cmd)
				assert.Nil(t, err)
				assert.Equal(t, tc.expect[i], result, "Value mismatch for cmd %s", cmd)
			}
		})
	}
}
