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
import { useParams } from "wouter";
import { useBaseUrl } from "../../../lib/navigation/util";
import BackButton from "../../../components/back-button";
import { useAcceptRelayActorFollowRequestMutation, useRejectRelayActorFollowRequestMutation, useRelayActorBlocksQuery, useRelayActorCreateBlockMutation, useRelayActorFollowersQuery, useRelayActorFollowRequestsQuery, useRelayActorRemoveBlockMutation, useRelayActorRemoveFromFollowersMutation } from "../../../lib/query/admin/relay-actors";
import FormWithData from "../../../lib/form/form-with-data";
import { Error as ErrorC } from "../../../components/error";
import { Account } from "../../../lib/types/account";
import MutationButton from "../../../components/form/mutation-button";
import AccountsCollectionWithButtons from "../../../components/accounts-list-with-action-buttons";

export default function RelayActorManageRelationships() {
	const params: { relayActorId: string, relationshipType: string } = useParams();
	const baseUrl = useBaseUrl();
	const backLocation: String = history.state?.backLocation ?? `~${baseUrl}`;
	
	switch (true) {
		case params.relationshipType === "followers":
			return (
				<div className="manage-followers">
					<h1><BackButton to={backLocation} /> Followers</h1>
					<FormWithData
						dataQuery={useRelayActorFollowersQuery}
						queryArg={params.relayActorId}
						DataForm={RelayActorFollowers}
						{...{ collectionOwnerID: params.relayActorId }}
					/>
				</div>
			);
		case params.relationshipType === "follow_requests":
			return (
				<div className="manage-follow-requests">
					<h1><BackButton to={backLocation} /> Follow Requests</h1>
					<FormWithData
						dataQuery={useRelayActorFollowRequestsQuery}
						queryArg={params.relayActorId}
						DataForm={RelayActorFollowRequests}
						{...{ collectionOwnerID: params.relayActorId }}
					/>
				</div>
			);
		case params.relationshipType === "blocks":
			return (
				<div className="manage-blocks">
					<h1><BackButton to={backLocation} /> Blocks</h1>
					<FormWithData
						dataQuery={useRelayActorBlocksQuery}
						queryArg={params.relayActorId}
						DataForm={RelayActorBlocks}
						{...{ collectionOwnerID: params.relayActorId }}
					/>
				</div>
			);
		default:
			return <ErrorC error={new Error("unrecognized relationship type")} />;
	}
}

interface Props {
	collectionOwnerID: string;
	data: Account[];
}

function RelayActorFollowers({ collectionOwnerID, data: accounts }: Props) {
	const count = accounts.length;
	const [ remove, removeRes ] = useRelayActorRemoveFromFollowersMutation();
	const [ block, blockRes ] = useRelayActorCreateBlockMutation();

	// Present follow requests list
	// with Remove and Block buttons.
	const getButtons = (account) => {
		return (
			<>
				<MutationButton
					label="Remove follower"
					name="remove follower"
					onClick={e => {
						e.stopPropagation();
						e.preventDefault();
						remove({ relayActorID: collectionOwnerID, accountID: account.id });
					}}
					result={removeRes}
					disabled={false}
				/>
				<MutationButton
					className="danger"
					label="Block"
					name="block"
					onClick={e => {
						e.stopPropagation();
						e.preventDefault();
						block({ relayActorID: collectionOwnerID, accountID: account.id });
					}}
					result={blockRes}
					disabled={false}
				/>
			</>
		);
	};

	return (
		<>
			<p className="count">{count} follower{ count !== 1 && "s" }</p>
			<AccountsCollectionWithButtons
				accounts={accounts}
				getButtons={getButtons}
			/>
		</>
	);
}

function RelayActorFollowRequests({ collectionOwnerID, data: accounts }: Props) {
	const count = accounts.length;
	const [ accept, acceptRes ] = useAcceptRelayActorFollowRequestMutation();
	const [ reject, rejectRes ] = useRejectRelayActorFollowRequestMutation();
	const [ block, blockRes ] = useRelayActorCreateBlockMutation();

	// Present follow requests list with
	// Accept, Reject, and Block buttons.
	const getButtons = (account) => {
		return (
			<>
				<MutationButton
					label="Accept"
					name="accept"
					onClick={e => {
						e.stopPropagation();
						e.preventDefault();
						accept({ relayActorID: collectionOwnerID, accountID: account.id });
					}}
					result={acceptRes}
					disabled={false}
				/>
				<MutationButton
					label="Reject"
					name="reject"
					onClick={e => {
						e.stopPropagation();
						e.preventDefault();
						reject({ relayActorID: collectionOwnerID, accountID: account.id });
					}}
					result={rejectRes}
					disabled={false}
				/>
				<MutationButton
					className="danger"
					label="Block"
					name="block"
					onClick={e => {
						e.stopPropagation();
						e.preventDefault();
						block({ relayActorID: collectionOwnerID, accountID: account.id });
					}}
					result={blockRes}
					disabled={false}
				/>
			</>
		);
	};

	return (
		<>
			<p className="count">{count} pending follow request{ count !== 1 && "s" }</p>
			<AccountsCollectionWithButtons
				accounts={accounts}
				getButtons={getButtons}
			/>
		</>
	);
}

function RelayActorBlocks({ collectionOwnerID, data: accounts }: Props) {
	const count = accounts.length;
	const [ unblock, unblockRes ] = useRelayActorRemoveBlockMutation();
	// Present blocks list
	// with Unblock button.
	const getButtons = (account) => {
		return (
			<>
				<MutationButton
					label="Unblock"
					name="unblock"
					onClick={e => {
						e.stopPropagation();
						e.preventDefault();
						unblock({ relayActorID: collectionOwnerID, accountID: account.id });
					}}
					result={unblockRes}
					disabled={false}
				/>
			</>
		);
	};

	return (
		<>
			<p className="count">{count} block{ count !== 1 && "s" }</p>
			<AccountsCollectionWithButtons
				accounts={accounts}
				getButtons={getButtons}
			/>
		</>
	);
}
