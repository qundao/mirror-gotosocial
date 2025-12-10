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

import { useTextInput } from "../../../../lib/form";
import { TextInput } from "../../../../components/form/inputs";
import MutationButton from "../../../../components/form/mutation-button";
import { useMediaPurgeMutation } from "../../../../lib/query/admin/actions";
import { formDomainValidator } from "../../../../lib/util/formvalidators";

export default function Purge({}) {
	const domainField = useTextInput("domain", {
		validator: formDomainValidator,
	});

	const [mediaPurge, mediaPurgeResult] = useMediaPurgeMutation();

	function submitPurge(e) {
		e.preventDefault();
		mediaPurge(domainField.value);
	}
    
	return (
		<form onSubmit={submitPurge}>
			<div className="form-section-docs">
				<h2>Purge</h2>
				<p>
					Purge remote media from the specified domain, including attachments, headers, avatars, and emojis.
					<br/>
					This will completely remove that domain's media from your storage.
					<br/>
					If the remote instance is still online, media will be refetched when needed.
					<br/>
					To prevent refetching of media, you can use a domain limit with a "reject" media policy.
				</p>
				<a
					href="https://docs.gotosocial.org/en/stable/admin/media_caching/"
					target="_blank"
					className="docslink"
					rel="noreferrer"
				>
					Learn more about media caching (opens in a new tab)
				</a>
			</div>
			<TextInput
				field={domainField}
				label={`Domain (without "https://" prefix)`}
				placeholder="example.org"
				autoCapitalize="none"
				spellCheck="false"
			/>
			<MutationButton
				disabled={!domainField.value}
				label="Purge"
				result={mediaPurgeResult}
			/>
		</form>
	);
}
