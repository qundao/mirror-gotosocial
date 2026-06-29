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

import React, { ReactNode, useEffect, useMemo } from "react";

import { useLazySearchInstancesQuery } from "../../../../lib/query/admin";
import { useBoolInput, useTextInput } from "../../../../lib/form";
import { PageableList } from "../../../../components/pageable-list";
import { Checkbox, Select, TextInput } from "../../../../components/form/inputs";
import MutationButton from "../../../../components/form/mutation-button";
import { useLocation, useSearch } from "wouter";
import { AdminInstance } from "../../../../lib/types/instance";
import { DateTimeMinute } from "../../../../components/datetime";

export function InstancesSearchForm() {
	const [ location, setLocation ] = useLocation();
	const search = useSearch();
	const urlQueryParams = useMemo(() => new URLSearchParams(search), [search]);
	const [ searchInstances, searchRes ] = useLazySearchInstancesQuery();

	// Populate search form using values from
	// urlQueryParams, to allow paging.
	const form = {
		domain: useTextInput("domain", { defaultValue: urlQueryParams.get("domain") ?? undefined }),
		order: useTextInput("order", { defaultValue: urlQueryParams.get("order") ?? undefined }),
		with_errors_only: useBoolInput("with_errors_only", { defaultValue: Boolean(urlQueryParams.get("with_errors_only")) ?? undefined }),
		limit: useTextInput("limit", { defaultValue: urlQueryParams.get("limit") ?? "40"}),
	};

	// On mount, trigger the search.
	useEffect(() => {
		searchInstances(Object.fromEntries(urlQueryParams), true);
	}, [urlQueryParams, searchInstances]);

	// Rather than triggering the search directly,
	// the "submit" button changes the location
	// based on form field params, and lets the
	// useEffect hook above actually do the search.
	function submitQuery(e) {
		e.preventDefault();
		
		// Parse query parameters.
		const entries = Object.entries(form).map(([k, v]) => {
			// Take only defined form fields.
			if (v.value === undefined) {
				return null;
			}
			if (typeof v.value === "string") {
				// Ignore 0 length strings.
				if (v.value.length === 0) {
					return null;
				}
				return [[k, v.value]];
			} else {
				if (!v.value) {
					// Not interested in false value
					// for "with_errors_only" as
					// false is the default anyway.
					return null;
				}
				return [[k, "true"]];
			}
		}).flatMap(kv => {
			// Remove any nulls.
			return kv || [];
		});

		const searchParams = new URLSearchParams(entries);
		setLocation(location + "?" + searchParams.toString());
	}

	// Location to return to when user clicks "back" on the detail view.
	const backLocation = location + (urlQueryParams.size > 0 ? `?${urlQueryParams}` : "");
	
	// Function to map an item to a list entry.
	function itemToEntry(instance: AdminInstance): ReactNode {
		return (
			<InstanceListEntry
				key={instance.id}
				instance={instance}
				linkTo={`/${instance.id}`}
				backLocation={backLocation}
			/>
		);
	}

	return (
		<>
			<form
				onSubmit={submitQuery}
				// Prevent password managers trying
				// to fill in username/email fields.
				autoComplete="off"
			>
				<TextInput
					field={form.domain}
					label={`Domain or first part of domain (without "https://" prefix)`}
					placeholder="example.org"
					autoCapitalize="none"
					spellCheck="false"
				/>
				<Select
					field={form.order}
					label="Order results by first seen (newest -> oldest), or alphabetical (a -> z)"
					options={
						<>
							<option value="first_seen">First seen</option>
							<option value="alphabetical">Alphabetical</option>
						</>
					}
				></Select>
				<Checkbox
					field={form.with_errors_only}
					label={"Show only instances with delivery errors"}
				/>
				<MutationButton
					disabled={false}
					label={"Search"}
					result={searchRes}
				/>
			</form>
			<PageableList
				isLoading={searchRes.isLoading}
				isFetching={searchRes.isFetching}
				isSuccess={searchRes.isSuccess}
				items={searchRes.data?.instances}
				itemToEntry={itemToEntry}
				isError={searchRes.isError}
				error={searchRes.error}
				emptyMessage={<b>No instances found that match your query.</b>}
				prevNextLinks={searchRes.data?.links}
			/>
		</>
	);
}

interface InstanceEntryProps {
	instance: AdminInstance;
	linkTo: string;
	backLocation: string;
}

function InstanceListEntry({ instance, linkTo, backLocation }: InstanceEntryProps) {
	const [ _location, setLocation ] = useLocation();
	
	const domain = instance.domain;
	const software = instance.software ?? "unknown";
	const deliveryErrors = instance.delivery_errors;

	const onClick = (e) => {
		e.preventDefault();
		// When clicking on a instance, direct
		// to the detail view for that instance.
		setLocation(linkTo, {
			// Store the back location in history so
			// the detail view can use it to return to
			// this page (including query parameters).
			state: { backLocation: backLocation }
		});
	};

	return (
		<span	
			className="pseudolink instance-info entry"
			aria-label={domain}
			title={domain}
			onClick={onClick}
			onKeyDown={(e) => {
				if (e.key === "Enter") {
					e.preventDefault();
					onClick(e);
				}
			}}
			role="link"
			tabIndex={0}
		>
			<h4 className="text-cutoff">{domain}</h4>
			<dl className="info-list">
				<div className="info-list-entry">
					<dt>Software:</dt>
					<dd className="text-cutoff">{software}</dd>
				</div>

				<div className="info-list-entry">
					<dt>First seen:</dt>
					<dd className="text-cutoff"><DateTimeMinute iso8601={instance.first_seen}/></dd>
				</div>

				<div className="info-list-entry">
					<dt>Latest successful delivery:</dt>
					<dd className="text-cutoff"><DateTimeMinute iso8601={instance.latest_successful_delivery}/></dd>
				</div>

				{ deliveryErrors &&
					<div className="info-list-entry delivery-errors">
						<dt>Delivery errors:</dt>
						<dd>{deliveryErrors.length}</dd>
					</div>
				}
			</dl>
		</span>
	);
}
