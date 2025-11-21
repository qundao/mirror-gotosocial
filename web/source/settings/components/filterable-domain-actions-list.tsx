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

import React, { FormEvent, useMemo } from "react";
import type {
	DomainAction,
	DomainActionType,
	DomainActionTypeUpper,
	MappedDomainActions
} from "../lib/types/domain";
import { Link, useLocation } from "wouter";
import { useTextInput } from "../lib/form";
import { matchSorter } from "match-sorter";
import { TextInput } from "./form/inputs";

export interface FilterableDomainActionsListProps {
    data: MappedDomainActions;
    type: DomainActionType;
    typeUpper: DomainActionTypeUpper;
	submitToLocation: (_filter: string) => string;
	linkToLocation: (_entry: DomainAction) => string;
}

export function FilterableDomainActionsList(props: FilterableDomainActionsListProps) {
	const {
		data,
		type,
		typeUpper,
		submitToLocation,
		linkToLocation,
	} = props;

	// Format items into a list.
	const items = useMemo(() => Object.values(data), [data]);

	const [_location, setLocation] = useLocation();
	const filterField = useTextInput("filter");

	function filterFormSubmit(e: FormEvent<HTMLFormElement>) {
		e.preventDefault();
		setLocation(submitToLocation(filter));
	}

	const filter = filterField.value ?? "";
	const filteredPerms = useMemo(() => matchSorter(items, filter, { keys: ["domain"] }), [items, filter]);
	const filtered = items.length - filteredPerms.length;

	const filterInfo = (
		<span>
			{items.length} {type}ed domain{items.length != 1 ? "s" : ""} {filtered > 0 && `(${filtered} filtered by search)`}
		</span>
	);

	const entries = filteredPerms.map((entry) => 
		<Link
			className="entry nounderline"
			key={entry.domain}
			to={linkToLocation(entry)}
		>
			<span id="domain">{entry.domain}</span>
			<span id="date">{new Date(entry.created_at ?? "").toLocaleString()}</span>
		</Link>
	);

	return (
		<div className="domain-actions-list">
			<form className="filter" role="search" onSubmit={filterFormSubmit}>
				<TextInput
					field={filterField}
					placeholder="example.org"
					label={`Search or add domain ${type}`}
				/>
				<button
					type="submit"
					disabled={
						filterField.value === undefined ||
						filterField.value.length == 0
					}
				>
					{typeUpper}&nbsp;{filter}
				</button>
			</form>
			<div>
				{filterInfo}
				<div className="list">
					<div className="entries scrolling">
						{entries}
					</div>
				</div>
			</div>
		</div>
	);
}
