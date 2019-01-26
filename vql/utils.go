/*
   Velociraptor - Hunting Evil
   Copyright (C) 2019 Velocidex Innovations.

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package vql

import (
	"encoding/json"

	"github.com/pkg/errors"
	actions_proto "www.velocidex.com/golang/velociraptor/actions/proto"
	vfilter "www.velocidex.com/golang/vfilter"
)

func ExtractRows(vql_response *actions_proto.VQLResponse) ([]vfilter.Row, error) {
	result := []vfilter.Row{}
	var rows []map[string]interface{}
	err := json.Unmarshal([]byte(vql_response.Response), &rows)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, row := range rows {
		item := vfilter.NewDict()
		for k, v := range row {
			item.Set(k, v)
		}
		result = append(result, item)
	}

	return result, nil
}

func RowToDict(scope *vfilter.Scope, row vfilter.Row) *vfilter.Dict {
	// If the row is already a dict nothing to do:
	result, ok := row.(*vfilter.Dict)
	if ok {
		return result
	}

	result = vfilter.NewDict()
	for _, column := range scope.GetMembers(row) {
		value, pres := scope.Associative(row, column)
		if pres {
			result.Set(column, value)
		}
	}

	return result
}
