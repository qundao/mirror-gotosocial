# Domain Limits

GoToSocial allows you to create "domain limits" to modify how your instance handles posts, follow requests, media, and accounts from remote domains.

This gives you the power to "fine-tune" federation with a problematic domain, without having to resort to putting a full [domain block](./domain_blocks.md) in place to cut off federation entirely.

For example, you can use a domain limit to mute all accounts on a given domain except for ones people on your instance follow, and/or to mark all media from a given domain as "sensitive", etc.

!!! tip
    When you create a domain limit, it extends to all subdomains as well, so limiting 'example.com' also limits 'social.example.com'.

You can view, create, and remove domain limits using the [instance admin panel](./settings.md#domain-limits).

Each domain limit has five components that you can tweak to tune federation with a limited domain:

- Content warning
- Media policy
- Follows policy
- Statuses policy
- Accounts policy

## Content Warning

Any text that you set as a content warning will be added to each post from the limited domain as part of the content warning / subject field when that post is viewed (in the home or public timeline, for example) using a client app.

If the post already has a content warning, any text set as the limit content warning will be prepended to the existing content warning with a semicolon.

!!! tip
    Setting a content warning will also have the effect of marking all posts (and attachments) from the limited domain as sensitive.

## Media Policy

You can apply a media policy in order to change whether and how your instance processes media attachments from the limited domain.

<dl>
	<dt>No limit</dt>
	<dd>Media will be processed as normal.</dd>
	<dt>Mark sensitive</dt>
	<dd>
		Media will be processed as normal.
		<br/>However, all post attachments from the limited domain will be marked sensitive.
	</dd>
	<dt>Reject</dt>
	<dd>
		No media from the limited domain will be downloaded, processed, or stored.
		<br/>This includes emoji, avatars, headers, and media attachments.
		<br/>Posts will contain a link to view attached media on the remote instance.
	</dd>
</dl>

## Follows Policy

You can apply a follows policy to determine how follows from the limited domain are processed.

Any restrictions put in place here apply even when a follow targets an unlocked account.

Note that this only applies to new follows from the moment you apply the policy, existing relationships will not be affected.

Accounts on this instance will still be able to follow accounts from the limited domain as normal.

<dl>
	<dt>No limit</dt>
	<dd>Follows will be processed as normal.</dd>
	<dt>Manual approval</dt>
	<dd>All follows originating from the limited domain will require manual approval.</dd>
	<dt>Reject non-mutual</dt>
	<dd>
		Each follow originating from the limited domain will be automatically rejected unless it is a "follow-back" or "mutual" follow.
		<br/>For example, if user A on this instance <em>does</em> already follows or follow-requests user B from the limited domain, user B will be able to send a follow (request) to user A as normal.
		<br/>However, if user A on this instance <em>does not</em> already follow or follow-request user B from the limited domain, any attempt by user B to follow user A will be automatically rejected.
	</dd>
	<dt>Reject</dt>
	<dd>Any follows originating from the limited domain will be automatically rejected.</dd>
</dl>

## Statuses Policy

You can apply a statuses policy to determine how statuses aka posts from the limited domain are processed.
This only applies to non-followed accounts. For example, if user A from this instance follows user B from the limited domain, user A will see user B's posts as normal.
However, if user A on this instance does not follow user B from the limited domain, user B's posts will be filtered from user A's perspective.

## Accounts Policy

You can apply an accounts policy to mute/silence accounts from the limited domain by default.
This only applies to non-followed accounts. For example, if user A from this instance follows user B from the limited domain, user B will not be muted.
However if user A from this instance does not follow user B from the limited domain, user B will be muted.
