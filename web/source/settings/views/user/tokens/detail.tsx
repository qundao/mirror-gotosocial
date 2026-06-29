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
import { useLocation, useParams } from "wouter";
import FormWithData from "../../../lib/form/form-with-data";
import BackButton from "../../../components/back-button";
import { useBaseUrl } from "../../../lib/navigation/util";
import { useGetTokenQuery, useInvalidateTokenMutation, useUpdateTokenMutation } from "../../../lib/query/user/tokens";
import { TokenInfo } from "../../../lib/types/tokeninfo";
import { useTextInput, useValue } from "../../../lib/form";
import useFormSubmit from "../../../lib/form/submit";
import MutationButton from "../../../components/form/mutation-button";
import { TextInput } from "../../../components/form/inputs";
import { DateTimeMinute } from "../../../components/datetime";
import WebsiteLink from "../../../components/website";

export default function TokenDetail({ }) {
	const params: { tokenId: string } = useParams();
	const baseUrl = useBaseUrl();
	const backLocation: String = history.state?.backLocation ?? `~${baseUrl}`;

	return (
		<div className="token-detail">
			<h1><BackButton to={backLocation}/> Token Details</h1>
			<FormWithData
				dataQuery={useGetTokenQuery}
				queryArg={params.tokenId}
				DataForm={TokenDetailForm}
			/>
		</div>
	);
}

function TokenDetailForm({ data: tokenInfo }: { data: TokenInfo, backLocation: string }) {	
	const [ _location, setLocation ] = useLocation();
	const baseUrl = useBaseUrl();

	const form = {
		id: useValue("id", tokenInfo.id),
		name: useTextInput("name", { source: tokenInfo })
	};
	const [ submit, result ] = useFormSubmit(form, useUpdateTokenMutation());
	const [ invalidate, invalidateResult ] = useInvalidateTokenMutation();

	return (
		<>
			<dl className="info-list">
				<div className="info-list-entry">
					<dt>App name:</dt>
					<dd>{tokenInfo.application.name}</dd>
				</div>
				{ tokenInfo.application.website && 
					<div className="info-list-entry">
						<dt>App website:</dt>
						<dd>{WebsiteLink(tokenInfo.application.website)}</dd>
					</div>
				}
				<div className="info-list-entry">
					<dt>Scope:</dt>
					<dd className="text-cutoff monospace">{tokenInfo.scope}</dd>
				</div>
				<div className="info-list-entry">
					<dt>Last used:</dt>
					<dd><DateTimeMinute iso8601={tokenInfo.last_used}/></dd>
				</div>
				<div className="info-list-entry">
					<dt>Created:</dt>
					<dd><DateTimeMinute iso8601={tokenInfo.created_at}/></dd>
				</div>
			</dl>
			<form onSubmit={submit}>
				<TextInput
					field={form.name}
					label="Name to use for this token"
					placeholder="Laptop firefox (windows)"
				/>
				<div className="action-buttons row">
					<MutationButton
						label="Update"
						result={result}
						disabled={!form.name.hasChanged()}
					/>
					<MutationButton
						label={"Invalidate token"}
						title={"Invalidate token"}
						type="button"
						className="button danger"
						onClick={() => {
							invalidate(tokenInfo.id);
							setLocation(`~${baseUrl}`);
						}}
						result={invalidateResult}
						showError={false}
						disabled={false}
					/>
				</div>
			</form>
		</>
	);
}
