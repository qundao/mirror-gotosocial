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

import React, { useMemo } from "react";

import { useLocation } from "wouter";
import { useBaseUrl } from "../lib/navigation/util";
import { RelayConnection } from "../lib/types/relay";
import MutationButton from "./form/mutation-button";
import UsernameLozenge from "./username-lozenge";
import { useBoolInput, useTextInput, useValue } from "../lib/form";
import { Checkbox, TextInput } from "./form/inputs";
import useFormSubmit from "../lib/form/submit";
import RelayFlags from "./relayflags";
import { DateTimeMinute } from "./datetime";

interface RelayDetailFormProps {
	data: RelayConnection,
	verb: string,
	updateHook,
	deleteHook,
	createMatcherHook,
	deleteMatcherHook,
}

export default function RelayDetailForm({
	data: conn,
	verb,
	updateHook,
	deleteHook,
	createMatcherHook,
	deleteMatcherHook,
}: RelayDetailFormProps) {
	const [ _location, setLocation ] = useLocation();
	const baseUrl = useBaseUrl();
	const form = {
		id: useValue("id", conn.id),
		public: useBoolInput("public", { source: conn }),
		unlisted: useBoolInput("unlisted", { source: conn }),
		match_by_default: useBoolInput("match_by_default", { source: conn }),
		ignore_sensitive: useBoolInput("ignore_sensitive", { source: conn }),
		ignore_media: useBoolInput("ignore_media", { source: conn }),
		ignore_replies: useBoolInput("ignore_replies", { source: conn }),
	};
	const [ submit, result ] = useFormSubmit(form, updateHook());
	const [ removeTrigger, removeResult ] = deleteHook();

	return (
		<>
			<dl className="info-list">
				<div className="info-list-entry">
					<dt>Actor URI:</dt>
					<dd>{conn.relay_actor_uri}</dd>
				</div>
				<div className="info-list-entry">
					<dt>Created at:</dt>
					<dd><DateTimeMinute iso8601={conn.created_at}/></dd>
				</div>
				{
					conn.account_id &&
						<div className="info-list-entry">
							<dt>Created by:</dt>
							<dd><UsernameLozenge account={conn.account_id}/></dd>
						</div>
				}
				<div className="info-list-entry">
					<dt>Approved:</dt>
					<dd className={`text-cutoff ${conn.approved ? "relay-connection-approved" : "relay-connection-not-approved"}`}>{conn.approved ? "yes" : "not yet"}</dd>
				</div>
			</dl>
			<form onSubmit={submit}>
				<RelayFlags
					verb={verb}
					form_field_public={form.public}
					form_field_unlisted={form.unlisted}
					form_field_ignore_sensitive={form.ignore_sensitive}
					form_field_ignore_media={form.ignore_media}
					form_field_ignore_replies={form.ignore_replies}
					form_field_match_by_default={form.match_by_default}
				/>

				<div className="action-buttons row">
					<MutationButton
						label="Update"
						result={result}
						disabled={
							!form.public.hasChanged() &&
							!form.unlisted.hasChanged() &&
							!form.match_by_default.hasChanged() &&
							!form.ignore_sensitive.hasChanged() &&
							!form.ignore_media.hasChanged() &&
							!form.ignore_replies.hasChanged()
						}
					/>

					<MutationButton
						type="button"
						onClick={() => {
							removeTrigger(conn.id);
							setLocation(`~${baseUrl}/overview`);
						}}
						label="Delete"
						result={removeResult}
						className="button danger"
						showError={false}
						disabled={false}
					/>
				</div>
			</form>
			<RelayMatchers
				conn={conn}
				createMatcherHook={createMatcherHook}
				deleteMatcherHook={deleteMatcherHook}
			/>
		</>
	);
}

interface RelayMatchersProps {
	conn: RelayConnection,
	createMatcherHook,
	deleteMatcherHook,
}

function RelayMatchers({
	conn,
	createMatcherHook,
	deleteMatcherHook,
}: RelayMatchersProps) {
	const form = {
		id: useValue("id", conn.id),
		keyword: useTextInput("keyword"),
		whole_word: useBoolInput("whole_word"),
		exclude: useBoolInput("exclude"),
	};
	const [ formSubmit, result ] = useFormSubmit(form, createMatcherHook());
	const [ remove, removeResult ] = deleteMatcherHook();

	// Link to appropriate page depending on if
	// this is a relay subscription or a relay push.
	//
	// We know account id is only ever defined
	// for relay subscriptions so check for that.
	const docsLink = useMemo(() => {
		if (conn.account_id) {
			// This is a relay subscription.
			return "https://docs.gotosocial.org/en/stable/admin/relay_subscriptions/#relay-matchers";
		} else {
			// Must be a relay push.
			return "https://docs.gotosocial.org/en/stable/user_guide/relay_pushes/#relay-matchers";
		}
	}, [conn.account_id]);

	return (
		<>
			<div className="form-section-docs">
				<h3>Matchers</h3>
				<p>
					You can add relay matchers to this connection to give granular control over which posts are relayed.
					<br/><br/>If the relay connection <em>does not</em> match posts by default, posts will only be relayed if their content is matched by a matcher. If you create no matchers, nothing will be relayed by the connection.
					<br/><br/>Conversely, if the relay connection <em>does</em> match posts by default, you can use exclude matchers to <em>prevent</em> posts from being relayed, based on their content. If you create no exclude matchers, everything will be relayed.
					<br/><br/>Regardless of whether the relay connection does or does not match posts by default, exclude matchers will prevent posts from being relayed, even if they would otherwise be matched (ie., exclude matchers take priority).
				</p>
				<a
					href={docsLink}
					target="_blank"
					className="docslink"
					rel="noreferrer"
				>
					Learn more about relay matchers (opens in a new tab)
				</a>
			</div>
			<form
				onSubmit={formSubmit}
				// Prevent password managers
				// trying to fill in fields.
				autoComplete="off"
			>
				<TextInput
					field={form.keyword}
					label="Keyword (case insensitive)"
					placeholder="#SomeHashtag"
					spellCheck="false"
					autoCapitalize="none"
				/>

				<Checkbox
					label={"Match whole word; if unchecked, allow matching word fragments"}
					field={form.whole_word}
				/>

				<Checkbox
					label={"Exclude posts matched by this matcher, instead of including them"}
					field={form.exclude}
				/>

				<MutationButton
					label="Create matcher"
					result={result}
					disabled={form.keyword.value == ""}
				/>
			</form>
			{ conn.matchers.length !== 0 &&
				<>
					<h4>Active matchers</h4>
					<ol className="matchers list">
						{ conn.matchers.map(matcher => {
							const label = `"${matcher.keyword}"; ${matcher.whole_word ? "whole word match" : "partial match"}; ${matcher.exclude ? "exclude matches" : "include matches"}`;
							return (
								<li
									className="entry"
									id={matcher.id}
									key={matcher.id}
									aria-label={label}
									title={label}
								>
									<div className="relay-flags-icons">
										{ matcher.whole_word
											? <>
												<div title="whole word match">
													<i className="fa fa-fw fa-text-width" aria-hidden="true"></i>
													<span className="sr-only">whole word match</span>
												</div>
											</>
											: <>
												<div title="partial word match">
													<i className="fa fa-fw fa-i-cursor" aria-hidden="true"></i>
													<span className="sr-only">partial word match</span>
												</div>
											</>
										}
										{ matcher.exclude
											? <>
												<div title="exclude matches">
													<i className="fa fa-fw fa-close" aria-hidden="true"></i>
													<span className="sr-only">exclude matches</span>
												</div>
											</>
											: <>
												<div title="include matches">
													<i className="fa fa-fw fa-check" aria-hidden="true"></i>
													<span className="sr-only">include matches</span>
												</div>
											</>
										}
									</div>
									<div className="relay-matcher-keyword">{matcher.keyword}</div>
									<MutationButton
										label="Delete"
										type="button"
										className="button danger"
										onClick={(e) => {
											e.preventDefault();
											remove({id: conn.id, matcherID: matcher.id});
										}}
										disabled={false}
										showError={false}
										result={removeResult}
									/>
								</li>
							);
						}) }
					</ol>
				</>
			}
		</>
	);
}