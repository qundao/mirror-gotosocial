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

import React from "react";

import { Link, useParams } from "wouter";
import Loading from "../../../components/loading";
import { useDomainAllowsQuery, useDomainBlocksQuery } from "../../../lib/query/admin/domain-permissions/get";
import type { MappedDomainPermissions, DomainPermissionType, DomainActionTypeUpper } from "../../../lib/types/domain";
import { NoArg } from "../../../lib/types/query";
import { useBaseUrl } from "../../../lib/navigation/util";
import { useCapitalize } from "../../../lib/util";
import { FilterableDomainActionsList } from "../../../components/filterable-domain-actions-list";

export default function DomainPermissionsOverview() {	
	const baseUrl = useBaseUrl();
	
	// Parse perm type from routing params.
	let params = useParams();
	if (params.permType !== "blocks" && params.permType !== "allows") {
		throw "unrecognized perm type " + params.permType;
	}
	// Safe to cast as we've already checked params.permType.
	const permType = params.permType.slice(0, -1) as DomainPermissionType;

	// Uppercase first letter of given permType.
	const permTypeUpper = useCapitalize(permType);

	// Fetch / wait for desired perms to load.
	const { data: blocks, isLoading: isLoadingBlocks } = useDomainBlocksQuery(NoArg, { skip: permType !== "block" });
	const { data: allows, isLoading: isLoadingAllows } = useDomainAllowsQuery(NoArg, { skip: permType !== "allow" });
	
	let data: MappedDomainPermissions | undefined;
	let isLoading: boolean;

	if (permType == "block") {
		data = blocks;
		isLoading = isLoadingBlocks;
	} else if (permType == "allow") {
		data = allows;
		isLoading = isLoadingAllows;
	} else {
		throw "unrecognized perm type " + permType;
	}

	if (isLoading || data === undefined) {
		return <Loading />;
	}
	
	return (
		<div className={`domain-${permType}`}>
			<div className="form-section-docs">
				<h1>Domain {permTypeUpper}s</h1>
				{ permType == "block" ? <BlockHelperText/> : <AllowHelperText/> }
			</div>
			<FilterableDomainActionsList
				data={data}
				type={permType}
				// Safe to cast as we've already checked permType.
				typeUpper={permTypeUpper as DomainActionTypeUpper}
				submitToLocation={filter => `/${permType}s/${filter}`}
				linkToLocation={entry => `/${permType}s/${entry.domain}`}
			/>
			<Link to={`~${baseUrl}/import-export`}>
				Or use the bulk import/export interface
			</Link>
		</div>
	);
}

function BlockHelperText() {
	return (
		<p>
			Blocking a domain blocks interaction between your instance, and all current and future accounts on
			instance(s) running on the blocked domain. Stored content will be removed, and no more data is sent to
			the remote server. This extends to all subdomains as well, so blocking 'example.com' also blocks 'social.example.com'.
			<br/>
			<a
				href="https://docs.gotosocial.org/en/latest/admin/domain_blocks/"
				target="_blank"
				className="docslink"
				rel="noreferrer"
			>
				Learn more about domain blocks (opens in a new tab)
			</a>
			<br/>
		</p>
	);
}

function AllowHelperText() {
	return (
		<p>
			Allowing a domain explicitly allows instance(s) running on that domain to interact with your instance.
			If you're running in allowlist mode, this is how you "allow" instances through.
			If you're running in blocklist mode (the default federation mode), you can use explicit domain allows
			to override domain blocks. In blocklist mode, explicitly allowed instances will be able to interact with
			your instance regardless of any domain blocks in place.  This extends to all subdomains as well, so allowing
			'example.com' also allows 'social.example.com'. This is useful when you're importing a block list but
			there are some domains on the list you don't want to block: just create an explicit allow for those domains
			before importing the list.
			<br/>
			<a
				href="https://docs.gotosocial.org/en/latest/admin/federation_modes/"
				target="_blank"
				className="docslink"
				rel="noreferrer"
			>
				Learn more about federation modes (opens in a new tab)
			</a>
		</p>
	);
}
