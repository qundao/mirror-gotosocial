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

import { useTextInput } from "../../../lib/form";
import { PageableList } from "../../../components/pageable-list";
import MutationButton from "../../../components/form/mutation-button";
import { useLocation, useSearch } from "wouter";
import { Select } from "../../../components/form/inputs";
import { useInvalidateTokenMutation, useLazySearchTokenInfoQuery } from "../../../lib/query/user/tokens";
import { TokenInfo } from "../../../lib/types/tokeninfo";
import WebsiteLink from "../../../components/website";
import { DateTimeMinute } from "../../../components/datetime";

export default function TokensSearchForm() {
	const [ location, setLocation ] = useLocation();
	const search = useSearch();
	const urlQueryParams = useMemo(() => new URLSearchParams(search), [search]);
	const [ searchTokenInfo, searchRes ] = useLazySearchTokenInfoQuery();

	// Populate search form using values from
	// urlQueryParams, to allow paging.
	const form = {
		limit: useTextInput("limit", { defaultValue: urlQueryParams.get("limit") ?? "25" }),
		order: useTextInput("order", { defaultValue: urlQueryParams.get("order") ?? "last_used"})
	};

	// On mount, trigger search.
	useEffect(() => {
		searchTokenInfo(Object.fromEntries(urlQueryParams), true);
	}, [urlQueryParams, searchTokenInfo]);

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
			} else if (typeof v.value === "string" && v.value.length === 0) {
				return null;
			}

			return [[k, v.value.toString()]];
		}).flatMap(kv => {
			// Remove any nulls.
			return kv !== null ? kv : [];
		});

		const searchParams = new URLSearchParams(entries);
		setLocation(location + "?" + searchParams.toString());
	}
	
	// Function to map an item to a list entry.
	function itemToEntry(tokenInfo: TokenInfo): ReactNode {
		return (
			<TokenInfoListEntry
				key={tokenInfo.id}
				tokenInfo={tokenInfo}
			/>
		);
	}

	return (
		<>
			<form
				onSubmit={submitQuery}
				// Prevent password managers
				// trying to fill in fields.
				autoComplete="off"
			>
				<Select
					field={form.order}
					label="Order results by last used (latest -> oldest), or creation time (newest -> oldest)"
					options={
						<>
							<option value="last_used">Last used</option>
							<option value="created">Creation time</option>
						</>
					}
				></Select>
				<Select
					field={form.limit}
					label="Items per page"
					options={
						<>
							<option value="25">25</option>
							<option value="50">50</option>
							<option value="75">75</option>
							<option value="100">100</option>
						</>
					}
				></Select>
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
				items={searchRes.data?.tokens}
				itemToEntry={itemToEntry}
				isError={searchRes.isError}
				error={searchRes.error}
				emptyMessage={<b>No tokens found.</b>}
				prevNextLinks={searchRes.data?.links}
			/>
		</>
	);
}

interface TokenInfoListEntryProps {
	tokenInfo: TokenInfo;
}

function TokenInfoListEntry({ tokenInfo }: TokenInfoListEntryProps) {
	const [ location, setLocation ] = useLocation();
	
	const onClick = (e) => {
		e.preventDefault();
		// When clicking on a token,
		// go to the detail view for it.
		setLocation(`/${tokenInfo.id}`, {
			// Store the back location in
			// history so the detail view
			// can use it to return here.
			state: { backLocation: location }
		});
	};

	const [ invalidate, invalidateResult ] = useInvalidateTokenMutation();

	return (
		<span
			className={"pseudolink token-info entry"}
			aria-label={`${tokenInfo.application.name}, scope: ${tokenInfo.scope}`}
			title={`${tokenInfo.application.name}, scope: ${tokenInfo.scope}`}
			onClick={onClick}
			onKeyDown={(e) => {
				if (e.key === "Enter") {
					e.preventDefault();
					onClick(e);
				}
			}}
			role="link"
			tabIndex={0}
			key={tokenInfo.id}
		>
			<dl className="info-list">
				{ tokenInfo.name && <>
					<div className="info-list-entry">
						<dt>Token name:</dt>
						<dd className="text-cutoff">{tokenInfo.name}</dd>
					</div>
				</>}
				<div className="info-list-entry">
					<dt>App name:</dt>
					<dd className="text-cutoff">{tokenInfo.application.name}</dd>
				</div>
				{ tokenInfo.application.website && 
					<div className="info-list-entry">
						<dt>App website:</dt>
						<dd className="text-cutoff">{WebsiteLink(tokenInfo.application.website)}</dd>
					</div>
				}
				<div className="info-list-entry">
					<dt>Scope:</dt>
					<dd className="text-cutoff monospace">{tokenInfo.scope}</dd>
				</div>
				<div className="info-list-entry">
					<dt>Last used:</dt>
					<dd className="text-cutoff"><DateTimeMinute iso8601={tokenInfo.last_used}/></dd>
				</div>
				<div className="info-list-entry">
					<dt>Created:</dt>
					<dd className="text-cutoff"><DateTimeMinute iso8601={tokenInfo.created_at}/></dd>
				</div>
			</dl>
			<div className="action-buttons">
				<MutationButton
					label={`Invalidate token`}
					title={`Invalidate token`}
					type="button"
					className="button danger"
					onClick={(e) => {
						e.preventDefault();
						e.stopPropagation();
						invalidate(tokenInfo.id);
					}}
					disabled={false}
					showError={true}
					result={invalidateResult}
				/>
			</div>
		</span>
	);
}
