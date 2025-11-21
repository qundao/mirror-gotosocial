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

import typia from "typia";
import { Links } from "parse-link-header";

export const validateDomainPerms = typia.createValidate<DomainPermission[]>();

/**
 * A single domain action entry, containing common
 * fields across limits, blocks, and allows.
 */
export interface DomainAction {
	id?: string;
	domain: string;
	created_at?: string;
	private_comment?: string;
	public_comment?: string;
	created_by?: string;
}

/**
 * Type of a domain action.
 */
export type DomainActionType = "limit" | "block" | "allow";

/**
 * Uppercase version of DomainActionType.
 */
export type DomainActionTypeUpper = "Limit" | "Block" | "Allow";

/**
 * Domain actions mapped to an Object where the Object
 * keys are the "domain" value of each DomainAction.
 */
export interface MappedDomainActions {
	[key: string]: DomainAction;
}

/**
 * Policy to apply to media files originating from the limited domain.
 */
export type DomainLimitMediaPolicy = "no_action" | "mark_sensitive" | "reject";

/**
 * Policy to apply to follow (requests) originating from the limited domain.
 */
export type DomainLimitFollowsPolicy = "no_action" | "manual_approval" | "reject_non_mutual" | "reject_all";

/**
 * Policy to apply to statuses of non-followed accounts on the limited domain.
 */
export type DomainLimitStatusesPolicy = "no_action" | "filter_warn" | "filter_hide";

/**
 * Policy to apply to non-followed accounts on the limited domain.
 */
export type DomainLimitAccountsPolicy = "no_action" | "mute";

/**
 * DomainLimit is a domain action that enforces specific policies
 * for media, follows, statuses, accounts, and content warning.
 */
export interface DomainLimit extends DomainAction {
	/**
	 * Policy to apply to media files originating from the limited domain.
	 */
	media_policy: DomainLimitMediaPolicy;
	
	/**
	 * Policy to apply to follow (requests) originating from the limited domain.
	 */
	follows_policy: DomainLimitFollowsPolicy;
	
	/**
	 * Policy to apply to statuses of non-followed accounts on the limited domain.
	 */
	statuses_policy: DomainLimitStatusesPolicy;
	
	/**
	 * Policy to apply to non-followed accounts on the limited domain.
	 */
	accounts_policy: DomainLimitAccountsPolicy;
	
	/**
	 * Content warning to prepend to statuses originating from the limited domain.
	 */
	content_warning?: string;
}

/**
 * Domain limits mapped to an Object where the Object
 * keys are the "domain" value of each domain limit.
 */
export interface MappedDomainLimits {
	[key: string]: DomainLimit;
}

export interface DomainLimitUpdateParams {
	/**
	 * Policy to apply to media files originating from the limited domain.
	 */
	media_policy?: DomainLimitMediaPolicy;
	
	/**
	 * Policy to apply to follow (requests) originating from the limited domain.
	 */
	follows_policy?: DomainLimitFollowsPolicy;
	
	/**
	 * Policy to apply to statuses of non-followed accounts on the limited domain.
	 */
	statuses_policy?: DomainLimitStatusesPolicy;
	
	/**
	 * Policy to apply to non-followed accounts on the limited domain.
	 */
	accounts_policy?: DomainLimitAccountsPolicy;
	
	/**
	 * Content warning to prepend to statuses originating from the limited domain.
	 */
	content_warning?: string;

	/**
	 * (Optionally) publicly stated reason for limiting the domain.
	 */
	public_comment?: string;

	/**
	 * Privately stated reason for limiting the domain.
	 */
	private_comment?: string;
}

export interface DomainLimitCreateParams extends DomainLimitUpdateParams {
	/**
	 * The hostname of the domain.
	 */
	domain: string;
}

/**
 * Type of a domain permission entry.
 */
export type DomainPermissionType = "block" | "allow" | "draft" | "exclude";

/**
 * A single domain permission entry of type block, allow, draft, or exclude.
 */
export interface DomainPermission extends DomainAction {
	obfuscate?: boolean;
	comment?: string;
	subscription_id?: string;

	// Fields that should be stripped before
	// sending a domain permission via the API.

	permission_type?: DomainPermissionType;
	key?: string;
	suggest?: string;
	valid?: boolean;
	checked?: boolean;
	commentType?: string;
	replace_private_comment?: boolean;
	replace_public_comment?: boolean;
}

/**
 * Domain permissions mapped to an Object where the Object
 * keys are the "domain" value of each DomainPerm.
 */
export interface MappedDomainPermissions {
	[key: string]: DomainPermission;
}

/**
 * Set of fields that should be stripped before
 * sending a domain permission via the API.
 */
const domainPermissionStripOnImport: Set<keyof DomainPermission> = new Set([
	"permission_type",
	"key",
	"suggest",
	"valid",
	"checked",
	"commentType",
	"replace_private_comment",
	"replace_public_comment",
]);

/**
 * Returns true if provided DomainPermission Object key is one
 * that should be stripped when importing a domain permission.
 * 
 * @param key 
 * @returns 
 */
export function stripOnImport(key: keyof DomainPermission) {
	return domainPermissionStripOnImport.has(key);
}

export interface ImportDomainPermissionsParams {
	domains: DomainPermission[];

	// Internal processing keys;
	// remove before serdes of form.
	obfuscate?: boolean;
	commentType?: string;
	permType: DomainPermissionType;
}

/**
 * Model domain permissions bulk export params.
 */
export interface ExportDomainPermissionsParams {
	permType: DomainPermissionType;
	action: "export" | "export-file";
	exportType: "json" | "csv" | "plain";
}

/**
 * Parameters for GET to /api/v1/admin/domain_permission_drafts.
 */
export interface DomainPermissionDraftsSearchParams {
	/**
	 * Show only drafts created by the given subscription ID.
	 */
	subscription_id?: string;
	/**
	 * Return only drafts that target the given domain.
	 */
	domain?: string;
	/**
	 * Filter on "block" or "allow" type drafts.
	 */
	permission_type?: DomainPermissionType;
	/**
	 * Return only items *OLDER* than the given max ID (for paging downwards).
	 * The item with the specified ID will not be included in the response.
	 */
	max_id?: string;
	/**
	 * Return only items *NEWER* than the given since ID.
	 * The item with the specified ID will not be included in the response.
	 */
	since_id?: string;
	/**
	 * Return only items immediately *NEWER* than the given min ID (for paging upwards).
	 * The item with the specified ID will not be included in the response.
	 */
	min_id?: string;
	/**
	 * Number of items to return.
	 */
	limit?: number;
}

export interface DomainPermissionDraftsSearchResp {
	drafts: DomainPermission[];
	links: Links | null;
}

export interface DomainPermissionDraftCreateParams {
	/**
	 * Domain to create the permission draft for.
	 */
	domain: string;
	/**
	 * Create a draft "allow" or a draft "block".
	 */
	permission_type: DomainPermissionType;
	/**
	 * Obfuscate the name of the domain when serving it publicly.
	 * Eg., `example.org` becomes something like `ex***e.org`.
	 */
	obfuscate?: boolean;
	/**
	 * Public comment about this domain permission. This will be displayed
	 * alongside the domain permission if you choose to share permissions.
	 */
	public_comment?: string;
	/**
	 * Private comment about this domain permission.
	 * Will only be shown to other admins, so this is a useful way of
	 * internally keeping track of why a certain domain ended up permissioned.
	 */
	private_comment?: string;
}

/**
 * Parameters for GET to /api/v1/admin/domain_permission_excludes.
 */
export interface DomainPermissionExcludesSearchParams {
	/**
	 * Return only excludes that target the given domain.
	 */
	domain?: string;
	/**
	 * Return only items *OLDER* than the given max ID (for paging downwards).
	 * The item with the specified ID will not be included in the response.
	 */
	max_id?: string;
	/**
	 * Return only items *NEWER* than the given since ID.
	 * The item with the specified ID will not be included in the response.
	 */
	since_id?: string;
	/**
	 * Return only items immediately *NEWER* than the given min ID (for paging upwards).
	 * The item with the specified ID will not be included in the response.
	 */
	min_id?: string;
	/**
	 * Number of items to return.
	 */
	limit?: number;
}

export interface DomainPermissionExcludesSearchResp {
	excludes: DomainPermission[];
	links: Links | null;
}

export interface DomainPermissionExcludeCreateParams {
	/**
	 * Domain to create the permission exclude for.
	 */
	domain: string;
	/**
	 * Private comment about this domain permission.
	 * Will only be shown to other admins, so this is a useful way of
	 * internally keeping track of why a certain domain ended up permissioned.
	 */
	private_comment?: string;
}

/**
 * Content type of a domain permission subscription.
 */
export type DomainPermissionSubscriptionContentType = "text/plain" | "text/csv" | "application/json";

/**
 * API model of one domain permission susbcription.
 */
export interface DomainPermissionSubscription {
	/**
	 * The ID of the domain permission subscription.
	 */
	id: string;
	/**
	 * The priority of the domain permission subscription.
	 */
	priority: number;
	/**
	 *  Time at which the subscription was created (ISO 8601 Datetime).
	 */
	created_at: string;
	/**
	 * Title of this subscription, as set by admin who created or updated it.
	 */
	title: string;
	/**
	 * The type of domain permission subscription (allow, block).
	 */
	permission_type: DomainPermissionType;
	/**
	 * If true, domain permissions arising from this subscription will be created as drafts that must be approved by a moderator to take effect.
	 * If false, domain permissions from this subscription will come into force immediately.
	 */
	as_draft: boolean;
	/**
	 * If true, this domain permission subscription will "adopt" domain permissions
	 * which already exist on the instance, and which meet the following conditions:
	 * 1) they have no subscription ID (ie., they're "orphaned") and 2) they are present
	 * in the subscribed list. Such orphaned domain permissions will be given this
	 * subscription's subscription ID value and be managed by this subscription.
	 */
	adopt_orphans: boolean;
	/**
	 * ID of the account that created this subscription.
	 */
	created_by: string;
	/**
	 * URI to call in order to fetch the permissions list.
	 */
	uri: string;
	/**
	 * MIME content type to use when parsing the permissions list.
	 */
	content_type: DomainPermissionSubscriptionContentType;
	/**
	 * (Optional) username to set for basic auth when doing a fetch of URI.
	 */
	fetch_username?: string;
	/**
	 * (Optional) password to set for basic auth when doing a fetch of URI.
	 */
	fetch_password?: string;
	/**
	 * Time at which the most recent fetch was attempted (ISO 8601 Datetime).
	 */
	fetched_at?: string;
	/**
	 *  Time of the most recent successful fetch (ISO 8601 Datetime).
	 */
	successfully_fetched_at?: string;
	/**
	 * If most recent fetch attempt failed, this field will contain an error message related to the fetch attempt.
	 */
	error?: string;
	/**
	 * Count of domain permission entries discovered at URI on last (successful) fetch.
	 */
	count: number;
}

/**
 * Parameters for GET to /api/v1/admin/domain_permission_subscriptions.
 */
export interface DomainPermissionSubscriptionSearchParams {
	/**
	 * Return only block or allow subscriptions.
	 */
	permission_type?: DomainPermissionType;
	/**
	 * Return only items *OLDER* than the given max ID (for paging downwards).
	 * The item with the specified ID will not be included in the response.
	 */
	max_id?: string;
	/**
	 * Return only items *NEWER* than the given since ID.
	 * The item with the specified ID will not be included in the response.
	 */
	since_id?: string;
	/**
	 * Return only items immediately *NEWER* than the given min ID (for paging upwards).
	 * The item with the specified ID will not be included in the response.
	 */
	min_id?: string;
	/**
	 * Number of items to return.
	 */
	limit?: number;
}

export interface DomainPermissionSubscriptionCreateUpdateParams {
	/**
	 * The priority of the domain permission subscription.
	 */
	priority?: number;
	/**
	 * Title of this subscription, as set by admin who created or updated it.
	 */
	title?: string;
	/**
	 * URI to call in order to fetch the permissions list.
	 */
	uri: string;
	/**
	 * MIME content type to use when parsing the permissions list.
	 */
	content_type: DomainPermissionSubscriptionContentType;
	/**
	 * If true, domain permissions arising from this subscription will be created as drafts that must be approved by a moderator to take effect.
	 * If false, domain permissions from this subscription will come into force immediately.
	 */
	as_draft?: boolean;
	/**
	 * If true, this domain permission subscription will "adopt" domain permissions
	 * which already exist on the instance, and which meet the following conditions:
	 * 1) they have no subscription ID (ie., they're "orphaned") and 2) they are present
	 * in the subscribed list. Such orphaned domain permissions will be given this
	 * subscription's subscription ID value and be managed by this subscription.
	 */
	adopt_orphans?: boolean;
	/**
	 * (Optional) username to set for basic auth when doing a fetch of URI.
	 */
	fetch_username?: string;
	/**
	 * (Optional) password to set for basic auth when doing a fetch of URI.
	 */
	fetch_password?: string;
	/**
	 * The type of domain permission subscription to create or update (allow, block).
	 */
	permission_type: DomainPermissionType;
}

export interface DomainPermissionSubscriptionsSearchResp {
	subs: DomainPermissionSubscription[];
	links: Links | null;
}
