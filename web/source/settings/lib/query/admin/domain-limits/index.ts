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

import type {
	DomainLimit,
	DomainLimitCreateParams,
	DomainLimitUpdateParams,
	MappedDomainLimits
} from "../../../types/domain";
import { gtsApi } from "../../gts-api";
import { removeFromCacheOnMutation, replaceCacheOnMutation, updateCacheOnMutation } from "../../query-modifiers";
import { listToKeyedObject } from "../../transforms";

const extended = gtsApi.injectEndpoints({
	endpoints: (build) => ({
		domainLimits: build.query<MappedDomainLimits, void>({
			query: () => ({
				url: `/api/v1/admin/domain_limits`
			}),
			transformResponse: listToKeyedObject<DomainLimit>("domain"),
		}),

		createDomainLimit: build.mutation<MappedDomainLimits, DomainLimitCreateParams>({
			query: (formData) => ({
				method: "POST",
				url: `/api/v1/admin/domain_limits`,
				asForm: true,
				body: formData,
				discardEmpty: true
			}),
			transformResponse: listToKeyedObject<DomainLimit>("domain"),
			...replaceCacheOnMutation("domainLimits"),
		}),

		updateDomainLimit: build.mutation<DomainLimit, {id: string} & DomainLimitUpdateParams>({
			query: ({ id, ...formData}) => ({
				method: "PUT",
				url: `/api/v1/admin/domain_limits/${id}`,
				asForm: true,
				body: formData,
				discardEmpty: false
			}),
			...updateCacheOnMutation("domainLimits", {
				key: (_draft, newData) => {
					return newData.domain;
				}
			})
		}),

		removeDomainLimit: build.mutation<DomainLimit, string>({
			query: (id) => ({
				method: "DELETE",
				url: `/api/v1/admin/domain_limits/${id}`,
			}),
			...removeFromCacheOnMutation("domainLimits", {
				key: (_draft, newData) => {
					return newData.domain;
				}
			})
		}),
	}),
});

/**
 * GET all domain limits from `/api/v1/admin/domain_limits` .
 */
const useGetAllDomainLimitsQuery = extended.useDomainLimitsQuery;

/**
 * POST a new domain limit to `/api/v1/admin/domain_limits`.
 */
const useCreateDomainLimitMutation = extended.useCreateDomainLimitMutation;

/**
 * Update a domain limit by PUTing to `/api/v1/admin/domain_limits/{id}`.
 */
const useUpdateDomainLimitMutation = extended.useUpdateDomainLimitMutation;

/**
 * Delete a domain limit by DELETEing to `/api/v1/admin/domain_limits/{id}`.
 */
const useRemoveDomainLimitMutation = extended.useRemoveDomainLimitMutation;

export {
	useGetAllDomainLimitsQuery,
	useCreateDomainLimitMutation,
	useUpdateDomainLimitMutation,
	useRemoveDomainLimitMutation,
};
