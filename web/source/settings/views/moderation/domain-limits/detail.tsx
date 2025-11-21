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
import { useDomainFromParams } from "../../../lib/util/domain";
import { useBaseUrl } from "../../../lib/navigation/util";
import { useLocation, useParams, useSearch } from "wouter";
import {
	useCreateDomainLimitMutation,
	useGetAllDomainLimitsQuery,
	useRemoveDomainLimitMutation,
	useUpdateDomainLimitMutation
} from "../../../lib/query/admin/domain-limits";
import Loading from "../../../components/loading";
import { Error } from "../../../components/error";
import BackButton from "../../../components/back-button";
import { DomainLimit } from "../../../lib/types/domain";
import UsernameLozenge from "../../../components/username-lozenge";
import { formDomainValidator } from "../../../lib/util/formvalidators";
import { useRadioInput, useTextInput } from "../../../lib/form";
import { FormSubmitEvent } from "../../../lib/form/types";
import useFormSubmit from "../../../lib/form/submit";
import { RadioGroup, TextArea, TextInput } from "../../../components/form/inputs";
import MutationButton from "../../../components/form/mutation-button";

export default function DomainLimitView() {
	const baseUrl = useBaseUrl();
	const search = useSearch();

	// Parse domain from routing params.
	const params = useParams();
	let domain = useDomainFromParams(params, search);
	if (domain === undefined) {
		throw "undefined domain";
	}
	
	// Normalize / decode domain
	// (it may be URL-encoded).
	domain = decodeURIComponent(domain);
	
	// Fetch all domain limits.
	const {
		data,
		isLoading,
		isFetching,
		isError,
		error,
	} = useGetAllDomainLimitsQuery();

	if (isLoading || isFetching) {
		return <Loading />;
	}

	if (isError) {
		return <Error error={error} />;
	}

	if (data === undefined) {
		throw "undefined data";
	}

	// Check if we already have a
	// limit in place for this domain.
	const existingLimit = data[domain];

	return (
		<div className="domain-limit-details">
			<h1><BackButton to={`~${baseUrl}`} /> <span>Domain limit for {domain}</span></h1>
			{ existingLimit
				? <DomainLimitDetails limit={existingLimit} />
				: <span>No stored limit yet, you can add one below:</span>
			}
			<CreateOrUpdateDomainLimit
				defaultDomain={domain}
				limit={existingLimit}
			/>
		</div>
	);
}

function DomainLimitDetails({ limit }: { limit: DomainLimit }) {
	const baseUrl = useBaseUrl();
	const [ location ] = useLocation();

	const created = useMemo(() => {
		if (limit.created_at) {
			return new Date(limit.created_at).toDateString();
		}
		return "unknown";
	}, [limit.created_at]);

	return (
		<dl className="info-list">
			<div className="info-list-entry">
				<dt>Created</dt>
				<dd><time dateTime={limit.created_at}>{created}</time></dd>
			</div>
			<div className="info-list-entry">
				<dt>Created By</dt>
				<dd>
					<UsernameLozenge
						account={limit.created_by}
						linkTo={`~/settings/moderation/accounts/${limit.created_by}`}
						backLocation={`~${baseUrl}${location}`}
					/>
				</dd>
			</div>
			<div className="info-list-entry">
				<dt>Domain</dt>
				<dd>{limit.domain}</dd>
			</div>
		</dl>
	);
}

interface CreateOrUpdateDomainLimitProps {
	defaultDomain: string;
	limit?: DomainLimit;
}

function CreateOrUpdateDomainLimit({ defaultDomain, limit }: CreateOrUpdateDomainLimitProps) {
	const isExistingLimit = limit !== undefined;

	const form = {
		domain: useTextInput("domain", {
			source: limit,
			defaultValue: defaultDomain,
			validator: formDomainValidator,
		}),
		privateComment: useTextInput("private_comment", { source: limit }),
		publicComment: useTextInput("public_comment", { source: limit }),
		contentWarning: useTextInput("content_warning", { source: limit }),
		mediaPolicy: useRadioInput("media_policy", {
			source: limit,
			defaultValue: "no_action",
			options: {
				no_action: "No limit",
				mark_sensitive: "Mark sensitive",
				reject: "Reject",
			}
		}),
		followsPolicy: useRadioInput("follows_policy", {
			source: limit,
			defaultValue: "no_action",
			options: {
				no_action: "No limit",
				manual_approval: "Manual approval",
				reject_non_mutual: "Reject non-mutual",
				reject_all: "Reject all",
			}
		}),
		statusesPolicy: useRadioInput("statuses_policy", {
			source: limit,
			defaultValue: "no_action",
			options: {
				no_action: "No limit",
				filter_warn: "Apply a 'warn' filter",
				filter_hide: "Apply a 'hide' filter",
			},
		}),
		accountsPolicy: useRadioInput("accounts_policy", {
			source: limit,
			defaultValue: "no_action",
			options: {
				no_action: "No limit",
				mute: "Mute/silence",
			},
		}),
	};

	// Derive appropriate submit action depending on whether this limit
	// doesn't exist yet (create new) or already exists (update it).
	const [ removeTrigger, removeResult ] = useRemoveDomainLimitMutation();
	const [ createLimit, createLimitResult ] = useCreateDomainLimitMutation();
	const [ updateLimit, updateLimitResult ] = useUpdateDomainLimitMutation();
	const [ createOrUpdateTrigger, createOrUpdateResult ] = useMemo(() => {
		if (!isExistingLimit) {
			return [createLimit, createLimitResult];
		} else {
			return [updateLimit, updateLimitResult];
		}
	}, [
		isExistingLimit,
		createLimit, createLimitResult,
		updateLimit, updateLimitResult,
	]);

	const [submit, submitResult] = useFormSubmit(
		form,
		[ createOrUpdateTrigger, createOrUpdateResult ],
		{
			changedOnly: isExistingLimit,
			// If we're updating an existing limit,
			// insert the limit into the mutation
			// data before submitting. Otherwise just
			// return the mutationData unmodified.
			customizeMutationArgs: (mutationData) => {
				if (isExistingLimit) {
					return {
						id: limit?.id,
						...mutationData,
					};
				} else {
					return mutationData;
				}
			},
		},
	);

	const [location, setLocation] = useLocation();
	function onSubmit(e: FormSubmitEvent) {
		// Adding a new domain limit happens on a url like
		// "/settings/moderation/limits/example.org", but if
		// domain input changes, that doesn't match anymore
		// and causes issues later on so, before submitting
		// the form, silently change url, and THEN submit.
		if (!isExistingLimit) {
			let correctUrl = `/${form.domain.value}`;
			if (location != correctUrl) {
				setLocation(correctUrl);
			}
		}
		return submit(e);
	}

	return (
		<form onSubmit={onSubmit}>
			{ !isExistingLimit && 
				<TextInput
					field={form.domain}
					label="Domain"
					placeholder="example.com"
					autoCapitalize="none"
					spellCheck="false"
				/>
			}

			<TextArea
				field={form.privateComment}
				label="Private comment (shown to admins only)"
				autoCapitalize="sentences"
				rows={3}
			/>
			
			<TextArea
				field={form.publicComment}
				label="Public comment (shown to members of this instance via the instance info page, and on the web if enabled)"
				autoCapitalize="sentences"
				rows={3}
			/>

			<div className="form-section-docs">
				<h3 id="content-warning-label">Content Warning</h3>
				<p>
					Any text that you set here will be used as a content warning for posts originating from the limited domain.
					<br/>If the post already has a content warning, text set here will be prepended to the existing content warning with a semicolon.
					<br/>Setting a content warning here will also have the effect of marking all posts (and attachments) from the limited domain as sensitive.
				</p>
			</div>
			<TextInput
				field={form.contentWarning}
				aria-labelledby="content-warning-label"
				placeholder="Possibly NSFW or whatever"
				autoCapitalize="sentences"
			/>

			<div className="form-section-docs">
				<h3>Media Policy</h3>
				<p>You can apply a media policy in order to change whether and how your instance processes media attachments from the limited domain.</p>
				<details>
					<summary>Cheatsheet</summary>
					<dl>
						<div>
							<dt>No limit</dt>
							<dd>Media will be processed as normal.</dd>
						</div>
						<div>
							<dt>Mark sensitive</dt>
							<dd>
								Media will be processed as normal.
								<br/>However, all post attachments from the limited domain will be marked sensitive.
							</dd>
						</div>
						<div>
							<dt>Reject</dt>
							<dd>
								No media from the limited domain will be downloaded, processed, or stored.
								<br/>This includes emoji, avatars, headers, and media attachments.
								<br/>Posts will contain a link to view attached media on the remote instance.
							</dd>
						</div>
					</dl>
				</details>
			</div>
			<RadioGroup field={form.mediaPolicy} />

			<div className="form-section-docs">
				<h3>Follows Policy</h3>
				<p>
					You can apply a follows policy to determine how follows from the limited domain are processed.
					<br/>Any restrictions put in place here apply <em>even when a follow targets an unlocked account</em>.
					<br/>Note that this only applies to new follows from the moment you apply the policy, existing relationships will not be affected.
					<br/>Accounts on this instance will still be able to follow accounts from the limited domain as normal.
				</p>
				<details>
					<summary>Cheatsheet</summary>
					<dl>
						<div>
							<dt>No limit</dt>
							<dd>Follows will be processed as normal.</dd>
						</div>
						<div>
							<dt>Manual approval</dt>
							<dd>All follows originating from the limited domain will require manual approval.</dd>
						</div>
						<div>
							<dt>Reject non-mutual</dt>
							<dd>
								Each follow originating from the limited domain will be automatically rejected unless it is a "follow-back" or "mutual" follow.
								<br/>For example, if user A on this instance <em>does</em> already follows or follow-requests user B from the limited domain, user B will be able to send a follow (request) to user A as normal.
								<br/>However, if user A on this instance <em>does not</em> already follow or follow-request user B from the limited domain, any attempt by user B to follow user A will be automatically rejected.
							</dd>
						</div>
						<div>
							<dt>Reject</dt>
							<dd>Any follows originating from the limited domain will be automatically rejected.</dd>
						</div>
					</dl>
				</details>
			</div>
			<RadioGroup
				field={form.followsPolicy}
				aria-labelledby="follows-policy-label"	
			/>

			<div className="form-section-docs">
				<h3>Statuses Policy</h3>
				<p>
					You can apply a statuses policy to determine how statuses aka posts from the limited domain are processed.
					<br/>This only applies to non-followed accounts. For example, if user A from this instance follows user B from the limited domain, user A will see user B's posts as normal.
					<br/>However, if user A on this instance does <em>not</em> follow user B from the limited domain, user B's posts will be filtered from user A's perspective.
				</p>
			</div>
			<RadioGroup field={form.statusesPolicy}	/>

			<div className="form-section-docs">
				<h3>Accounts Policy</h3>
				<p>
					You can apply an accounts policy to mute/silence accounts from the limited domain by default.
					<br/>This only applies to non-followed accounts. For example, if user A from this instance follows user B from the limited domain, user B will not be muted.
					<br/>However if user A from this instance does <em>not</em> follow user B from the limited domain, user B will be muted.
				</p>
			</div>
			<RadioGroup	field={form.accountsPolicy} />

			<div className="action-buttons row">
				<MutationButton
					label={isExistingLimit ? "Update Limit" : "Limit"}
					result={submitResult}
					disabled={
						isExistingLimit &&
						!form.privateComment.hasChanged() &&
						!form.publicComment.hasChanged() &&
						!form.mediaPolicy.hasChanged() &&
						!form.followsPolicy.hasChanged() &&
						!form.statusesPolicy.hasChanged() &&
						!form.accountsPolicy.hasChanged() &&
						!form.contentWarning.hasChanged()
					}
				/>

				{ isExistingLimit &&
					<button
						type="button"
						onClick={() => removeTrigger(limit.id?? "")}
						className="button danger"
					>
						Remove Limit
					</button>
				}
			</div>

			<>
				{createOrUpdateResult.error && <Error error={createOrUpdateResult.error} />}
				{removeResult.error && <Error error={removeResult.error} />}
			</>
		</form>
	);
}
