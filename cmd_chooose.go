/*
 * Copyright (C) 2016 Eric Chou <zhssmail@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import "errors"

// args[0] is the profile to be choosen
func chooseHandler(args []string, data interface{}) (int, error) {
	if len(args) == 0 || args[0] == "" {
		return 0, errors.New("'choose' needs name or index of profile.")
	}
	p, index, err := doGetProfile(args[0], gConfig)
	if err != nil {
		return 0, err
	}

	printCmd("choose", p.Name)
	gConfig.Current = index
	gConfig.save()

	return 0, nil
}
