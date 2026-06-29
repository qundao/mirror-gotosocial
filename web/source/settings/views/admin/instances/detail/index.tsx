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

import React, { ReactNode } from "react";

import { useClearInstanceDeliveryErrorsMutation, useGetInstanceQuery } from "../../../../lib/query/admin";
import FormWithData from "../../../../lib/form/form-with-data";
import { AdminInstance } from "../../../../lib/types/instance";
import { useLocation, useParams } from "wouter";
import { useBaseUrl } from "../../../../lib/navigation/util";
import BackButton from "../../../../components/back-button";
import MutationButton from "../../../../components/form/mutation-button";
import { DateTimeMinute, DateTimeSecond } from "../../../../components/datetime";

export default function InstanceDetail() {
	const params: { instanceID: string } = useParams();
	const baseUrl = useBaseUrl();
	const backLocation: String = history.state?.backLocation ?? `~${baseUrl}`;

	return (
		<div className="instance-detail">
			<h1><BackButton to={backLocation} /> Instance Details</h1>
			<FormWithData
				dataQuery={useGetInstanceQuery}
				queryArg={params.instanceID}
				DataForm={InstanceDetailForm}
			/>
		</div>
	);
}

function InstanceDetailForm({ data: instance }: { data: AdminInstance }) {
	const domain = instance.domain;
	const software = instance.software ?? "unknown";

	return (
		<>
			<dl className="info-list">
				<div className="info-list-entry">
					<dt>Domain:</dt>
					<dd>
						<a
							href={`https://${domain}`}
							target="_blank"
							rel="noreferrer"
						>
							<i className="fa fa-fw fa-external-link" aria-hidden="true"></i> {domain} (opens in a new tab)
						</a>
					</dd>
				</div>

				<div className="info-list-entry">
					<dt>Software:</dt>
					<dd>{software}</dd>
				</div>

				<div className="info-list-entry">
					<dt>First seen:</dt>
					<dd><DateTimeMinute iso8601={instance.first_seen}/></dd>
				</div>
				<div className="info-list-entry">
					<dt>Latest successful delivery:</dt>
					<dd><DateTimeMinute iso8601={instance.latest_successful_delivery}/></dd>
				</div>
			</dl>
			<InstanceDeliveryErrors data={instance} />
		</>
	);
}

function InstanceDeliveryErrors({ data: instance }: { data: AdminInstance }): ReactNode {
	const baseUrl = useBaseUrl();
	const backLocation = `~${baseUrl}/${instance.id}`;
	const [ clearDeliveryErrors, clearDeliveryErrorsResult ] = useClearInstanceDeliveryErrorsMutation();
	
	if (!instance.delivery_errors) {
		return null;
	}
	
	return (
		<>
			<div className="form-section-docs">
				<h3>Recent Delivery Errors</h3>
				<p>
					This section shows the 20 most recent delivery errors since the latest successful delivery of an activity
					to an inbox on this instance (if ever). If the instance appears to have gone offline permanently, you may
					wish to block it using the <DomainBlocksLink domain={instance.domain} backLocation={backLocation} />.
					This cleans up accounts and statuses from the decommissioned instance stored in your database, and avoids the
					risk of later federating with an instance created by a baddie masquerading as the original instance owner.
				</p>
			</div>
			<MutationButton
				label={"Clear delivery errors"}
				title={"Clear delivery errors"}
				type="button"
				className="button danger"
				onClick={(e) => {
					e.preventDefault();
					clearDeliveryErrors(instance.id);
				}}
				disabled={false}
				showError={false}
				result={clearDeliveryErrorsResult}
			/>
			<dl className="info-list delivery-errors-list">
				{ instance.delivery_errors.map((err, i) => {
					return (
						<div className="info-list-entry" key={i}>
							<dt><DateTimeSecond iso8601={err.time} /></dt>
							<dd>{err.error}</dd>
						</div>
					);
				}) }
			</dl>
		</>
	);
}

function DomainBlocksLink({ domain, backLocation }: { domain: string, backLocation: string }) {
	const [ _location, setLocation ] = useLocation();
	const linkTo = `~/settings/moderation/domain-permissions/blocks/${domain}`;
	const onClick = () => {
		setLocation(linkTo, {
			// Store the back location in history so
			// it can be used to return to this page.
			state: { backLocation: backLocation }
		});
	};
	
	return (
		<span
			className="domain-blocks-link pseudolink"
			onClick={onClick}
			onKeyDown={e => e.key === "Enter" && onClick()}
			role="link"
			tabIndex={0}
		>
			domain blocks page
		</span>
	);
}
