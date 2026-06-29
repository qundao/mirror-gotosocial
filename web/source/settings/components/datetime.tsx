/*
	GoToSocial
	Copyright (C) GoToSocial Authors admin@gotosocial.org
	SPDX-License-Identifier: AGPL-3.0-or-later

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

import React, { useMemo } from "react";

/**
 * DateTimeMinute parses the given iso8601 time string
 * and returns a time element with nicely formatted content
 * like `May 22, 2026, 12:52 PM`. If the given
 * string is empty or undefined, it will instead return
 * "unknown / never" inside a react fragment.
 * @param iso8601 
 * @returns 
 */
export function DateTimeMinute({iso8601}: {iso8601: string | undefined}) {
	const set = !(!iso8601);
	const v = useMemo(() => {
		if (!set) {
			return "unknown / never";
		}
		return dtMinuteFormat.format(new Date(iso8601));
	}, [set, iso8601]);

	return (
		<>
			{ iso8601 ? <time dateTime={iso8601}>{v}</time> : <>{v}</> } 
		</>
	);
}

const dtMinuteFormat = new Intl.DateTimeFormat(undefined, {
	dateStyle: "medium",
	timeStyle: "short",
});

/**
 * DateTimeSecond parses the given iso8601 time string
 * and returns a time element with nicely formatted content
 * like `May 22, 2026, 12:52:36 PM`. If the given
 * string is empty or undefined, it will instead return
 * "unknown / never" inside a react fragment.
 * @param iso8601 
 * @returns 
 */
export function DateTimeSecond({iso8601}: {iso8601: string | undefined}) {
	const set = !(!iso8601);
	const v = useMemo(() => {
		if (!set) {
			return "unknown / never";
		}
		return dtSecondFormat.format(new Date(iso8601));
	}, [set, iso8601]);

	return (
		<>
			{ iso8601 ? <time dateTime={iso8601}>{v}</time> : <>{v}</> } 
		</>
	);
}

const dtSecondFormat = new Intl.DateTimeFormat(undefined, {
	dateStyle: "medium",
	timeStyle: "medium",
});
