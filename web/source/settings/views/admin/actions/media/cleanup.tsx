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
import { useMediaCleanupMutation } from "../../../../lib/query/admin/actions";

export default function Cleanup({}) {
	const media_remote_cache_duration = useTextInput("media_remote_cache_duration", { defaultValue: "1 day" });

	const [mediaCleanup, mediaCleanupResult] = useMediaCleanupMutation();

	function submitCleanup(e) {
		e.preventDefault();
		mediaCleanup(media_remote_cache_duration.value);
	}
	return (
		<form onSubmit={submitCleanup}>
			<div className="form-section-docs">
				<h2>Cleanup</h2>
				<p>
					Cleanup (by removing from storage) remote media, headers, avatars, and emojis
					older than the given duration string (<code>1 second</code>, <code>1 day</code>,
					<code>1 week</code>, etc).
					<br/>
					If you specify <code>0 sec</code> here, then the value of your runtime config
					variable <code>media_remote_cache_duration</code> will be used instead.
					<br/>
					To clear as much media as possible, pass something like <code>1 sec</code>.
					<br/>
					If the remote instance is still online, any media removed
					from storage in this way will be recached when needed.
				</p>
				<a
					href="https://docs.gotosocial.org/en/stable/admin/media_caching/"
					target="_blank"
					className="docslink"
					rel="noreferrer"
				>
					Learn more about media caching + cleanup (opens in a new tab)
				</a>
			</div>
			<TextInput
				field={media_remote_cache_duration}
				label="Duration"
				placeholder="1 day"
			/>
			<MutationButton
				disabled={!media_remote_cache_duration.value}
				label="Cleanup"
				result={mediaCleanupResult}
			/>
		</form>
	);
}
