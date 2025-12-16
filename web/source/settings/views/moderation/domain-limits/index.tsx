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
import { useGetAllDomainLimitsQuery } from "../../../lib/query/admin/domain-limits";
import Loading from "../../../components/loading";
import { Error } from "../../../components/error";
import { FilterableDomainActionsList } from "../../../components/filterable-domain-actions-list";

export default function DomainLimitsOverview() {
	const {
		data,
		isLoading,
		isFetching,
		isError,
		error,
	} = useGetAllDomainLimitsQuery();

	if (isLoading || isFetching) {
		return <Loading/>;
	}

	if (isError) {
		return <Error error={error} />;
	}

	if (!data) {
		throw "no data";
	}

	return (
		<div className={`domain-limits`}>
			<div className="form-section-docs">
				<h1>Domain Limits</h1>
				<p>
					Domain limits can be used to modify how your instance handles posts, follow requests, media, and accounts from remote domains.
					This allows you to "fine-tune" federation with a domain without having to resort to putting a full domain block in place to cut
					off federation entirely. For example, you can use a domain limit to mute all accounts on a given domain except for ones people
					on your instance follow, and/or to mark all media from a given domain as "sensitive", etc. When you create a domain limit, it
					extends to all subdomains as well, so limiting 'example.com' also limits 'social.example.com'. 
					<br/>
					<a
						href="https://docs.gotosocial.org/en/stable/admin/domain_limits/"
						target="_blank"
						className="docslink"
						rel="noreferrer"
					>
						Learn more about domain limits (opens in a new tab)
					</a>
				</p>
			</div>
			<FilterableDomainActionsList
				data={data}
				type={"limit"}
				typeUpper={"Limit"}
				submitToLocation={filter => `/${filter}`}
				linkToLocation={entry => `/${entry.domain}`}
			/>
		</div>
	);
}
