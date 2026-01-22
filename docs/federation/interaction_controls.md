# Interaction Controls

GoToSocial includes extensions to the ActivityStreams / ActivityPub protocol using the following types and properties to add interaction controls, interaction authorizations, and authorization validation:

- `interactionPolicy` object property.
  - Sub-policies `canLike`, `canReply`, `canAnnounce`.
  - Sub-policy properties `automaticApproval` and `manualApproval`.
- `LikeRequest`, `ReplyRequest`, and `AnnounceRequest` activities (aka interaction requests).
- `LikeAuthorization`, `ReplyAuthorization`, and `AnnounceAuthorization` objects (aka interaction authorizations).
- `interactingObject` and `interactionTarget` properties on authorizations.
- `likeAuthorization`, `replyAuthorization`, and `announceAuthorization` properties on posts (`Note`, `Question`, etc).

The json-ld `@context` document for these types and properties is hosted at `https://gotosocial.org/ns`.

The usage of these types and properties to provide users with interaction controls is described below.

!!! danger
    Interaction controls are an attempt to limit the harmful effects of unwanted replies and other interactions on a user's posts (eg., "reply guys").
    
    However, it is far from being sufficient for this purpose, as there are still many "out-of-band" ways that posts can be distributed or replied to beyond a user's initial wishes or intentions.
    
    For example, a user might create a post with a very strict interaction policy attached to it, only to find that other server softwares do not respect that policy, and users on other instances are having discussions and replying to the post *from their instance's perspective*. The original poster's instance will automatically drop these unwanted interactions from their view, but remote instances may still show them.
    
    Another example: someone might see a post that specifies nobody can reply to it, but screenshot the post, post the screenshot in their own new post, and tag the original post author in as a mention. Alternatively, they might just link to the URL of the post and tag the author in as a mention. In this case, they effectively "reply" to the post by creating a new thread.
    
    For better or worse, interaction controls can offer only a best-effort, partial, technological solution to what is more or less an issue of social behavior and boundaries.

!!! info "Deprecated `always` and `approvalRequired` properties"
    Previous versions of this document used the properties `always` and `approvalRequired`. These are now deprecated in favor of `automaticApproval` and `manualApproval`. GoToSocial versions before v0.20.0 send and receive only these deprecated properties. GoToSocial v0.20.0 sends and receives both the deprecated and the new properties. GoToSocial v0.21.0 onwards uses only the new properties. 

## Interaction Policy

GoToSocial uses the property `interactionPolicy` on posts, in order to indicate to remote instances what sort of interactions are (conditionally) permitted to be processed and stored, for any given post.

`interactionPolicy` is an object property attached to the post-like `Object`s `Note`, `Article`, `Question`, etc, with the following format:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    [... rest of the @context ...]
  ],
  "type": "Note",
  "interactionPolicy": {
    "canLike": {
      "automaticApproval": [ "zero_or_more_uris" ],
      "manualApproval": [ "zero_or_more_uris" ]
    },
    "canReply": {
      "automaticApproval": [ "zero_or_more_uris" ],
      "manualApproval": [ "zero_or_more_uris" ]
    },
    "canAnnounce": {
      "automaticApproval": [ "zero_or_more_uris" ],
      "manualApproval": [ "zero_or_more_uris" ]
    }
  },
  [... rest of the Note ...],
}
```

In the `interactionPolicy` object:

- `canLike` is a sub-policy which indicates who is permitted to create a `Like` with the post URI as the `object` of the `Like`.
- `canReply` is a sub-policy which indicates who is permitted to create a post with `inReplyTo` set to the URI/ID of the post.
- `canAnnounce` is a sub-policy which indicates who is permitted to create an `Announce` with the post URI/ID as the `object` of the `Announce`. 

And:

- `automaticApproval` denotes ActivityPub URIs/IDs of `Actor`s or `Collection`s of `Actor`s who will receive automated approval from the post author when creating an interaction targeting a post.
- `manualApproval` denotes ActivityPub URIs/IDs of `Actor`s or `Collection`s of `Actor`s who will receive approval from the post author at the author's own discretion; this means they may not receive approval at all, or their interaction may be rejected.

Valid URI entries in `automaticApproval` and `manualApproval` include:

- the magic ActivityStreams Public URI `https://www.w3.org/ns/activitystreams#Public`
- the URIs of the post creator's `Following` and/or `Followers` collections
- individual Actor URIs

For example:

```json
[
    "https://www.w3.org/ns/activitystreams#Public",
    "https://example.org/users/someone/followers",
    "https://example.org/users/someone/following",
    "https://example.org/users/boobslover6969",
    "https://boobs.example.com/users/someone_on_a_different_instance"
]
```

### Defaults per sub-policy

When an interaction policy is only *partially* defined (eg., only `canReply` is set, `canLike` or `canAnnounce` keys are not set), then implementations should make the following assumptions for each sub-policy in the `interactionPolicy` object that is *undefined*.

!!! tip "Future extensions with different defaults"
    Note that **the below list is not exhaustive**, and extensions to `interactionPolicy` may wish to define **different defaults** for other types of interaction, for example `canQuote`.

### `canLike`

If `canLike` is missing on an `interactionPolicy`, or the value of `canLike` is `null` or `{}`, then implementations should assume:

```json
"canLike": {
  "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
}
```

In other words, the default is **anyone who can see the post can like it**.

### `canReply`

If `canReply` is missing on an `interactionPolicy`, or the value of `canReply` is `null` or `{}`, then implementations should assume:

```json
"canReply": {
  "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
}
```

In other words, the default is **anyone who can see the post can reply to it**.

### `canAnnounce`

If `canAnnounce` is missing on an `interactionPolicy`, or the value of `canAnnounce` is `null` or `{}`, then implementations should assume:

```json
"canAnnounce": {
  "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
}
```

In other words, the default is **anyone who can see the post can announce it**.

### Default / fallback `interactionPolicy`

When the `interactionPolicy` property is not present at all on a post, or the `interactionPolicy` key is set but its value resolves to `null` or `{}`, implementations can assume the following implicit, default `interactionPolicy` for that post (which follows the [defaults per sub-policy](#defaults-per-sub-policy) described above):

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    [... rest of the @context ...]
  ],
  "type": "Note",
  "interactionPolicy": {
    "canLike": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    },
    "canReply": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    },
    "canAnnounce": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    }
  },
  [... rest of the Note ...]
}
```

As implied by the lack of any `manualApproval` property in any of the sub-policies, the default value for `manualApproval` is an empty array.

This default `interactionPolicy` reflects the de facto interaction policy of ActivityPub server softwares that are not interaction policy aware. That is to say, it is exactly what servers that are not interaction policy aware *already assume* about interaction permissions.

!!! info "Actors can only ever interact with a post they are permitted to see"
    Note that even when assuming a default `interactionPolicy` for a post, the **visibility** of a post must still be accounted for by looking at the `to`, `cc`, and/or `audience` properties, to ensure that actors who cannot *see* a post also cannot *interact* with it. Eg., if a post is addressed to followers-only, and the default `interactionPolicy` is assumed, then someone who does not follow the post creator should still *not* be able to see *or* interact with it.

!!! tip
    As is standard across AP implementations, implementations will likely still wish to limit `Announce` actities targeting the post to only the author themself if the post is addressed to followers-only.

### Indicating that verification is required / not required per sub-policy

Because few server softwares have implemented interaction policies at the time of writing, it is necessary to provide a method by which implementing servers can indicate to one another that they are both **aware of** and **will enforce** interaction policies as described below in the section [Requesting, Obtaining, and Proving Interaction Authorization](#requesting-obtaining-and-proving-interaction-authorization).

This indication of interaction policy participation is done via a server explicitly setting `interactionPolicy` and its sub-policies on outgoing posts, instead of relying on the implicit defaults described above when `interactionPolicy` is not set.

That is, **by setting `interactionPolicy.*` on a post, an instance indicates to other instances that they will enforce validation of interactions for each sub-policy that is explicitly set.**

This means that implementations must always explicitly set all sub-policies on an `interactionPolicy` for which they have implemented interaction controls themselves, and with which they would like other servers to comply, *even when the values do not differ from the [implicit defaults](#defaults-per-sub-policy)*.

For example, if a server understands and wishes to enforce the `canLike`, `canReply`, and `canAnnounce` sub-policies (as is the case with GoToSocial), then they should explicitly set those sub-policies on an outgoing post. This indicates to remote servers that the origin server knows about interaction controls, does enforcement, and knows how to handle appropriate `Reject` / `Accept` messages for each sub-policy.

Another example: if a server only implements the `canReply` interaction sub-policy, but not `canLike` or `canAnnounce`, then they should always set `interactionPolicy.canReply`, and leave the other two sub-policies out of the `interactionPolicy` to indicate that they cannot understand or enforce them.

This means of indicating participation in interaction policies through the absence of presence of keys was designed so that the large majority of servers that *do not* set `interactionPolicy` at all, because they have not (yet) implemented it, do not need to change their behavior. Servers that do implement `interactionPolicy` can understand, by the absence of the `interactionPolicy` key on a post, that the origin server is not `interactionPolicy` aware, and behave accordingly.

### Specifying Nobody

To specify that **nobody** can perform an interaction on a post **except** for its author (who is always permitted), implementations should set the `automaticApproval` array to **just the URI of the post author**, and either omit `manualApproval` or leave it empty.

For example, the following `canLike` value indicates that nobody can `Like` the post it is attached to except for the post author:

```json
[... rest of the interaction policy ...],
"canLike": {
  "automaticApproval": ["the_activitypub_uri_of_the_post_author"]
},
[... rest of the interaction policy ...]
```

Another example. The following `interactionPolicy` on a post by `https://example.org/users/someone` indicates that anyone can like the post, but nobody but the author can reply or announce:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    [... rest of the @context ...]
  ],
  "type": "Note",
  "interactionPolicy": {
    "canLike": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    },
    "canReply": {
      "automaticApproval": ["https://example.org/users/someone"]
    },
    "canAnnounce": {
      "automaticApproval": ["https://example.org/users/someone"]
    }
  },
  [... rest of the Note ...]
}
```

!!! note
    To avoid mischief, GoToSocial makes implicit assumptions about who can/can't interact, even if a policy specifies nobody. See [implicit assumptions](#implicit-assumptions).

### Implicit Assumptions

For common-sense safety reasons, GoToSocial makes, and will always apply, two implicit assumptions about interaction policies.

#### 1. Mentioned + replied-to actors can always reply

Actors mentioned in, or replied to by, a post should **ALWAYS** be able to reply to that post without requiring approval, regardless of the post visiblity and the `interactionPolicy`, **UNLESS** the post that mentioned or replied to them is itself currently pending approval.

This is to prevent a would-be harasser from mentioning someone in an abusive post, and leaving no recourse to the mentioned user to reply.

As such, when sending out interaction policies, GoToSocial will **ALWAYS** add the URIs of mentioned users to the `canReply.automaticApproval` array, unless they are already covered by the ActivityStreams magic public URI.

Likewise, when enforcing received interaction policies, GoToSocial will **ALWAYS** behave as though the URIs of mentioned users were present in the `canReply.automaticApproval` array, even if they weren't.

#### 2. An actor can always interact in any way with their own post

**Secondly**, an actor should **ALWAYS** be able to reply to their own post, like their own post, and boost their own post without requiring approval, **UNLESS** that post is itself currently pending approval.

As such, when sending out interaction policies, GoToSocial will **ALWAYS** add the URI of the post author to the `canLike.automaticApproval`, `canReply.automaticApproval`, and `canAnnounce.automaticApproval` arrays, **UNLESS** they are already covered by the ActivityStreams magic public URI.

Likewise, when enforcing received interaction policies, GoToSocial will **ALWAYS** behave as though the URI of the post author themself is present in each `automaticApproval` field, even if it wasn't.

### Conflicting / Duplicate Values

In cases where a user is present in a Collection URI, and is *also* targeted explicitly by URI, the **more specific value** takes precedence.

For example:

```json
[... rest of the interaction policy ...],
"canReply": {
  "automaticApproval": ["https://example.org/users/someone"],
  "manualApproval": ["https://www.w3.org/ns/activitystreams#Public"]
},
[... rest of the interaction policy ...]
```

Here, `@someone@example.org` is present in `automaticApproval`, and is also implicitly present in the magic ActivityStreams Public collection in `manualApproval`. In this case, they can always reply, as the `automaticApproval` value is more explicit.

Another example:

```json
[... rest of the interaction policy ...],
"canReply": {
  "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"],
  "manualApproval": ["https://example.org/users/someone"]
},
[... rest of the interaction policy ...]
```

Here, `@someone@example.org` is present in `manualApproval`, but is also implicitly present in the magic ActivityStreams Public collection in `automaticApproval`. In this case everyone can reply without approval, **except** for `@someone@example.org`, who requires approval.

In case the **exact same** URI is present in both `automaticApproval` and `manualApproval`, the **highest level of permission** takes precedence (ie., a URI in `automaticApproval` takes precedence over the same URI in `manualApproval`).

### Examples

Here's some examples of what interaction policies allow users to do.

#### 1. Limiting scope of a conversation

In this example, the user `@the_mighty_zork` wants to begin a conversation with the users `@booblover6969` and `@hodor`.

To avoid the discussion being derailed by others, they want replies to their post by users other than the three participants to be permitted only if they're approved by `@the_mighty_zork`.

Furthermore, they want to limit the boosting / `Announce`ing of their post to only their own followers, and to the three conversation participants.

However, anyone should be able to `Like` the post by `@the_mighty_zork`.

This can be achieved with the following `interactionPolicy`, which is attached to a post with visibility level public:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    [... rest of the @context ...]
  ],
  "type": "Note",
  "interactionPolicy": {
    "canLike": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    },
    "canReply": {
      "automaticApproval": [
        "https://example.org/users/the_mighty_zork",
        "https://example.org/users/booblover6969",
        "https://example.org/users/hodor"
      ],
      "manualApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    },
    "canAnnounce": {
      "automaticApproval": [
        "https://example.org/users/the_mighty_zork",
        "https://example.org/users/the_mighty_zork/followers",
        "https://example.org/users/booblover6969",
        "https://example.org/users/hodor"
      ]
    }
  },
  [... rest of the Note ...]
}
```

#### 2. Long solo thread

In this example, the user `@the_mighty_zork` wants to write a long solo thread.

They don't mind if people boost and like posts in the thread, but they don't want to get any replies because they don't have the energy to moderate the discussion; they just want to vent by throwing their thoughts out there.

This can be achieved by setting the following `interactionPolicy` on every post in the thread:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    [... rest of the @context ...]
  ],
  "type": "Note",
  "interactionPolicy": {
    "canLike": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    },
    "canReply": {
      "automaticApproval": ["https://example.org/users/the_mighty_zork"]
    },
    "canAnnounce": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    }
  },
  [... rest of the Note ...]
}
```

Here, anyone is allowed to like or boost, but nobody is permitted to reply (except `@the_mighty_zork` themself).

#### 3. Completely open

In this example, `@the_mighty_zork` wants to write a completely open post that can be replied to, boosted, or liked by anyone who can see it:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    [... rest of the @context ...]
  ],
  "type": "Note",
  "interactionPolicy": {
    "canLike": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    },
    "canReply": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    },
    "canAnnounce": {
      "automaticApproval": ["https://www.w3.org/ns/activitystreams#Public"]
    }
  },
  [... rest of the Note ...]
}
```

### Subsequent Replies / Scope Widening

Each subsequent reply in a conversation will have its own interaction policy, chosen by the user who created the reply. In other words, the entire *conversation* or *thread* is not controlled by one `interactionPolicy`, but the policy can differ for each subsequent post in a thread, as set by the post author.

Unfortunately, this means that even with `interactionPolicy` in place, the scope of a thread can inadvertently widen beyond the intention of the author of the first post in the thread.

For instance, in [example 1](#example-1---limiting-scope-of-a-conversation) above, `@the_mighty_zork` specifies in the first post a `canReply.automaticApproval` value of

```json
[
  "https://example.org/users/the_mighty_zork",
  "https://example.org/users/booblover6969",
  "https://example.org/users/hodor"
]
```

In a subsequent reply, either accidentally or on purpose `@booblover6969` sets the `canReply.automaticApproval` value to:

```json
[
  "https://www.w3.org/ns/activitystreams#Public"
]
```

This widens the scope of the conversation, as now anyone can reply to `@booblover6969`'s post, and possibly also tag `@the_mighty_zork` in that reply.

To avoid this issue, it is recommended that remote instances prevent users from being able to widen scope (exact mechanism of doing this TBD).

It is also a good idea for instances to consider any interaction with a post-like `Object` that is itself currently pending approval, as also pending approval. 

In other words, instances should mark all children interactions below a pending-approval parent as also pending approval, no matter what the interaction policy on the parent would ordinarily allow.

This avoids situations where someone could reply to a post, then, even if their reply is pending approval, they could reply *to their own reply* and have that marked as permitted (since as author, they would normally have [implicit permission to reply](#implicit-assumptions)).

## Requesting, Obtaining, and Proving Interaction Authorization

The [interaction policy](#interaction-policy) section described the shape of interaction policies, assumed defaults, assumptions, and examples.

This section describes how servers that are interaction policy aware can request approval for an interaction, how servers should send acceptance or rejection of a requested interaction, and how third party servers can verify that authorization to interact with a post has been granted to the interactor from the interactee.

### Interaction Request Types

The `LikeRequest`, `ReplyRequest`, and `AnnounceRequest` activity types (together called "interaction requests") are the basis of the mechanism whereby actors can *politely* request permission to perform an interaction.

In an interaction request activity, `object` is the URI of the post being interacted with, and `instrument` contains the desired interaction (either a post, a `Like`, or an `Announce`).

### When to Send an Interaction Request

There are two cases in which an actor should politely request approval for an interaction by sending an interaction request to the inbox of another actor who has created a post that has an `interactionPolicy` set:

#### Case 1: Manual Approval is Required

An actor should polite request approval for an interaction when their desired interaction falls into the `manualApproval` category of the target post's interaction policy.

For example, actor `@boobslover6969@boobs.example.com` should send a `ReplyRequest` if they wish to reply to a post with the following interaction policy, which specifies that manual approval is required for all replies:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    [... rest of the @context ...]
  ],
  "type": "Note",
  "interactionPolicy": {
    "canReply": {
      "automaticApproval": [
        "https://example.org/users/someone"
      ],
      "manualApproval": [
        "https://www.w3.org/ns/activitystreams#Public"
      ]
    }
  },
  [... rest of the Note ...]
}
  ```

#### Case 2: Approval is Conditional on Presence in a Followers or Following Collection

An actor should politely request approval for an interaction when the interacting actor is (conditionally) permitted to interact with a post because they follow or are followed by the post author.

For example, if actor `@boobslover6969@boobs.example.com` follows ``@someone@example.org`, they should send a `ReplyRequest` if they wish to reply to a post with the following interaction policy:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    [... rest of the @context ...]
  ],
  "type": "Note",
  "interactionPolicy": {
    "canReply": {
      "automaticApproval": [
        "https://example.org/users/someone",
        "https://example.org/users/someone/followers"
      ]
    }
  },
  [... rest of the Note ...]
}
```

!!! info
    This case is necessary because in a case where an actor on Server A is followed by an actor on Server B, Server C has no way to validate that the actor on Server B is a follower of the actor on Server A, without iterating through their followers collection (which may not be available).
    
    Rather than having third party servers try to iterate through followers/following collections, it is simpler to have the actor from Server A send a polite interaction request, wait for approval, and then transmit that approval as proof that Server A agrees that the actor from Server B is indeed a follower and is therefore permitted to interact.

### Sending an Interaction Request and Waiting for Approval

Say that `@boobslover6969@boobs.example.com` wants to reply to a post by`@someone@example.org`, and they fall into one of the above cases, they should send a `ReplyRequest` of the following form:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    "https://www.w3.org/ns/activitystreams",
    [... rest of the @context ...]
  ],
  "type": "ReplyRequest",
  "id": "https://boobs.example.com/users/boobslover6969/reply_requests/01KFGBRPGDCTAMB5MWVW700Y11",
  "actor": "https://boobs.example.com/users/boobslover6969",
  "to": "https://example.org/users/someone",
  "object": "https://example.org/users/someone/statuses/01KFGBJN586W5EZFHE16R7FMJ1",
  "instrument": {
    "type": "Note",
    "id": "https://boobs.example.com/users/boobslover6969/statuses/01KFGBRPFJJCP0FM4039H4VDVQ",
    "attributedTo": "https://boobs.example.com/users/boobslover6969",
    "inReplyTo": "https://example.org/users/someone/statuses/01KFGBJN586W5EZFHE16R7FMJ1",
    "to": [
      "https://www.w3.org/ns/activitystreams#Public",
    ],
    "cc": [
      "https://boobs.example.com/users/boobslover6969/followers",
      "https://example.org/users/someone"
    ],
    "content": "<p>hey what's up</p>",
    [... rest of the Note ...]
  }
}
```

In this `ReplyRequest`, the `object` is the URI of the post being interacted with, and the `instrument` field contains the interaction itself. Providing the interaction as `instrument` allows the target server to validate and store it.

Note that while the `Note` in this `ReplyRequest` `instrument` is addressed to public and cc'd to followers and `someone@example.org`, the `ReplyRequest` itself is addressed only to `@someone@example.org`, and should only be delivered to their inbox.

Once the `ReplyRequest` has been sent to `someone@example.org`'s inbox, then before displaying the reply in timelines, distributing it to followers, or assuming it to be now be included in the replies collection of the replied-to post, the server of `@boobslover6969@boobs.example.com` must wait for an `Accept` to be sent back from `example.org`. (This is much like how the `Accept` process works in response to `Follow`s.) In other words, the interaction can now be considered to be **pending**.

### `Accept`ing an Interaction Request

Once an interaction request has been received by a server, that server can approve the interaction by sending an `Accept` back to the inbox of the requester, referring back to the interaction request.

For example, say that `@boobslover6969@boobs.example.com` has sent a `ReplyRequest` to `@someone@example.org` as described [above](#sending-an-interaction-request-and-waiting-for-approval), in order to approve the reply, the server `example.org` should send an `Accept` back to `@boobslover6969@boobs.example.com`. This `Accept` should look like the following:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    "https://www.w3.org/ns/activitystreams"
  ],,
  "type": "Accept",
  "id": "https://example.org/users/someone/accepts/01KFGBRQBK8PPCQBE7K7Z22J5T",
  "actor": "https://example.org/users/someone",
  "to": "https://boobs.example.com/users/boobslover6969",
  "object": {
    "actor": "https://boobs.example.com/users/boobslover6969",
    "id": "https://boobs.example.com/users/boobslover6969/reply_requests/01KFGBRPGDCTAMB5MWVW700Y11",
    "instrument": "https://boobs.example.com/users/boobslover6969/statuses/01KFGBRPFJJCP0FM4039H4VDVQ",
    "object": "https://example.org/users/someone/statuses/01KFGBJN586W5EZFHE16R7FMJ1",
    "type": "ReplyRequest"
  },
  "result": "https://example.org/users/someone/authorizations/01KFGBRQBK8PPCQBE7K7Z22J5T"
}
```

In this `Accept`, the `object` is the interaction request that was sent and is now being accepted, and the `result` field is the URI of an authorization object that the server of `@boobslover6969@boobs.example.com` can then attach to the interaction to prove that it has been approved/authorized by `someone@example.org` (more on this shortly).

Once a server has delivered an `Accept` of an interaction to the interacting actor, the interaction should no longer be considered **pending**, it should be considered **approved** and can therefore be shown in timelines, shown on the web, highlighted in notifications, etc, as normal.

### Proving Interaction Authorization to Other Servers

Once an `Accept` has been received by a server in response to an interaction request, that server should consider the interaction to be **aproved**, as the interactor has received authorization to perform the interaction.

To prove that an interaction has been `Accepted`, and is **approved**, the interacting server should at this point distribute the interaction according to its addressing, **with the authorization attached in the relevant authorization field**.

At the time of writing, there are three authorization fields used by GoToSocial: `likeAuthorization`, `replyAuthorization`, and `announceAuthorization`. 

Following the above example, say that `@someone@example.org` has sent an `Accept` to `@boobslover6969@boobs.example.com` in response to their `ReplyRequest`, the server of `@boobslover6969@boobs.example.com` should now distribute the `Create` activity for the reply note to followers as normal, and include the value of the `result` field of the `Accept` in the `replyAuthorization` field of the note, as follows:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    "https://www.w3.org/ns/activitystreams"
    [... rest of the @context ...]
  ],
  "type": "Create",
  "object": {
    "type": "Note",
    "id": "https://boobs.example.com/users/boobslover6969/statuses/01KFGBRPFJJCP0FM4039H4VDVQ",
    "attributedTo": "https://boobs.example.com/users/boobslover6969",
    "inReplyTo": "https://example.org/users/someone/statuses/01KFGBJN586W5EZFHE16R7FMJ1",
    "to": [
      "https://www.w3.org/ns/activitystreams#Public",
    ],
    "cc": [
      "https://boobs.example.com/users/boobslover6969/followers",
      "https://example.org/users/someone"
    ],
    "content": "<p>hey what's up</p>",
    "replyAuthorization": "https://example.org/users/someone/authorizations/01KFGBRQBK8PPCQBE7K7Z22J5T"
    [... rest of the Note ...]
  },
  [... rest of the Create ...]
}
```

The presence of the `replyAuthorization` proves to other servers that the reply has been authorized by `example.org`.

### Interaction Authorization Objects

Servers that are capable of responding to `LikeRequest`, `ReplyRequest`, and `AnnounceRequest` types must serve the resulting interaction authorizations at the URI of the ID of the authorization, so that remote instances may dereference and verify them.

GoToSocial uses three interaction authorization object types: `LikeAuthorization`, `ReplyAuthorization`, and `AnnounceAuthorization`. These objects each have a property `interactionTarget` and `interactingObject`, which will be set to the URI of the `object` of an interaction request, and the URI of the `id` of the interaction, respectively.

For example, a `ReplyAuthorization` that authorizes a reply by `@boobslover6969@boobs.example.com` to a post by `@someone@example.org` would look as follows:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    "https://www.w3.org/ns/activitystreams"
  ],
  "attributedTo": "https://example.org/users/someone",
  "id": "https://example.org/users/someone/authorizations/01KFGBRQBK8PPCQBE7K7Z22J5T",
  "interactingObject": "https://boobs.example.org/users/boobslover6969/statuses/01KFGBRPFJJCP0FM4039H4VDVQ",
  "interactionTarget": "https://example.org/users/someone/statuses/01KFGBJN586W5EZFHE16R7FMJ1",
  "type": "ReplyAuthorization"
}
```

This interaction authorization must be served by `example.org` at the URI `https://example.org/users/someone/authorizations/01KFGBRQBK8PPCQBE7K7Z22J5T` so that other servers can dereference it in order to do verification.

### Verifying Interaction Authorization as a Third Server

When receiving an `Activity` and/or a `Create` of a post-like `Object` which represents an interaction targeting a post with an interaction policy set, servers that are neither the interactor nor the interactee (ie., third servers) should check the interaction policy of the targeted post to see whether the interaction is (conditionally) permitted.

If the post being interacted with has a relevant sub-policy with an `automaticApproval` value equal to the ActivityPub public URI, then the interaction can already be considered **approved**, as in this case it is not necessary for the interactor to have sent an interaction request (and received an `Accept`) before distributing the interaction (see [When to Send an Interaction Request](#when-to-send-an-interaction-request)).

For example, if Server A receives a reply from Server B to a post from Server C, and the post from Server C has a `canReply.automaticApproval` value that includes `"https://www.w3.org/ns/activitystreams#Public"`, then there is no need for Server A to verify interaction authorization, as the replier from Server B is implicitly permitted to reply by this lenient policy.

On the other hand, if an interaction would require approval because [manual approval is required](#case-1-manual-approval-is-required), or [approval is conditional on the interactor's presence in a followers or following collection](#case-2-approval-is-conditional-on-presence-in-a-followers-or-following-collection), then the receiving server should check the interaction for a filled `likeAuthorization`, `replyAuthorization`, or `announceAuthorization` field. If the relevant field is not present or not complete, then the server should consider the interaction invalid and drop it.

If the relevant field for the interaction type is filled out, then the server should **verify** the interaction authorization by doing the following:

1. Validate that the host/domain of the authorization URI is equal to the host/domain of the author of the post being interacted with.
2. Dereference the authorization URI/ID to get the authorization object.
3. Validate the interaction authorization to ensure that it authorizes the interaction that the interactor claims it does.

Ie., if an interaction is a reply, then the `replyAuthorization` field must be set to the URI of a `ReplyAuthorization`. The `ReplyAuthorization` must be present on the server of the interactee, and it must have `interactingObject` set to the ID/URI of the reply, `interactionTarget` set to the replied-to post, and `attributedTo` set to the URI of the author of the post being replied to.

If the authorization is verified and valid, then the server doing the verification should consider the interaction to be **approved**, and is free to show it in timelines, on the web, etc.

If the interaction is forbidden by the interaction policy of the post it targets, or the supplied interaction authorization does not pass verification, then the interaction should be considered invalid and dropped.

### `Reject`ing an Interaction Request

If the target of an interaction request does not wish to approve and authorize the interaction, they may instead send out a `Reject` message in response to the interaction request.

For example, the following `Reject` object may be sent by `@someone@example.org` to inform `@boobslover6969@boobs.example.com` that their reply request has been rejected:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    "https://www.w3.org/ns/activitystreams"
  ],,
  "type": "Reject",
  "id": "https://example.org/users/someone/rejects/01KFGBRQBK8PPCQBE7K7Z22J5T",
  "actor": "https://example.org/users/someone",
  "to": "https://boobs.example.com/users/boobslover6969",
  "object": {
    "actor": "https://boobs.example.com/users/boobslover6969",
    "id": "https://boobs.example.com/users/boobslover6969/reply_requests/01KFGBRPGDCTAMB5MWVW700Y11",
    "instrument": "https://boobs.example.com/users/boobslover6969/statuses/01KFGBRPFJJCP0FM4039H4VDVQ",
    "object": "https://example.org/users/someone/statuses/01KFGBJN586W5EZFHE16R7FMJ1",
    "type": "ReplyRequest"
  }
}
```

If a server receives a `Reject` in response to an interaction request, then they should handle the rejection by removing the proposed interaction, not showing it in timelines etc, and possibly putting it in some kind of "rejected interaction" table.

## Optional But Recommended Behaviors

### Handling Interactions From Servers that are Not Interaction Controls Aware

Because not all ActivityPub server softwares implement interaction controls yet, it is necessary for servers that are interaction control aware to be able to handle interactions that an interaction policy *does* implicitly permit, but which do *not* go through the polite interaction request flow described above.

For example, say that Server A is interaction policy aware, and Actor A on Server A has set an interaction policy on a post which requires manual approval for all replies.

A second server, Server B is not interaction policy aware. Actor B from Server B wants to reply to the post of Actor A, and so their server sends a `Create` of a `Note` to Actor A's inbox, and to the inboxes of everyone who follows Actor B.

Following the guidance for [verifying authentication](#verifying-interaction-authorization-as-a-third-server), any third servers that receive the `Create` `Note` without `replyAuthorization` set should consider the interaction as invalid and drop the reply.

Unfortunately, this leads to situations where whole threads of conversation may be considered valid by non-interaction-controls-aware servers, but invalid by servers that are aware of (and enforce) interaction controls.

To remedy this, it is **very strongly recommended** that implementers of interaction controls provide a code path for handling so-called "impolite" interaction requests like the `Create` `Note` described above, which may otherwise be conditionally permitted but which have not gone through the request and approval flow.

GoToSocial, for example, stores so-called "impolite" interactions in the database, and shows them to the interactee as a pending interaction request that can be accepted or rejected, just as a normal interaction request would be. In this case, the "impolite" `Create` `Note` above from Server B would appear to Actor A on Server A *as though* it had been sent as a polite `ReplyRequest`.

If Actor A refuses the "impolite" interaction request, then the reply is dropped from Server A and will not appear in the replies collection of the post by Actor A.

However if Actor A approves the "impolite" interaction request, their server should send out the **slightly different kind of `Accept` described below**.

### Broadcasting `Accept`s for the Benefit of Third Servers

When sending out an `Accept` of an "impolite" interaction that was sent directly to the interactee, outside of the interaction request flow, it is **very strongly recommended** that implementers of interaction controls build and address the `Accept` slightly differently than they would if accepting a "polite" interaction request.

For example, instead of sending out an `Accept` of a polite request like this:

```json
{
  "@context": [
    "https://gotosocial.org/ns",
    "https://www.w3.org/ns/activitystreams"
  ],,
  "type": "Accept",
  "id": "https://example.org/users/someone/accepts/01KFGBRQBK8PPCQBE7K7Z22J5T",
  "actor": "https://example.org/users/someone",
  "to": "https://boobs.example.com/users/boobslover6969",
  "object": {
    "actor": "https://boobs.example.com/users/boobslover6969",
    "id": "https://boobs.example.com/users/boobslover6969/reply_requests/01KFGBRPGDCTAMB5MWVW700Y11",
    "instrument": "https://boobs.example.com/users/boobslover6969/statuses/01KFGBRPFJJCP0FM4039H4VDVQ",
    "object": "https://example.org/users/someone/statuses/01KFGBJN586W5EZFHE16R7FMJ1",
    "type": "ReplyRequest"
  },
  "result": "https://example.org/users/someone/authorizations/01KFGBRQBK8PPCQBE7K7Z22J5T"
}
```

An `Accept` of an "impolite" request would look like this:

```json
{
  "@context": "https://www.w3.org/ns/activitystreams",,
  "type": "Accept",
  "id": "https://example.org/users/someone/accepts/01KFGBRQBK8PPCQBE7K7Z22J5T",
  "actor": "https://goblin.technology/users/tobi",
  "to": "https://boobs.example.com/users/boobslover6969",
  "cc": [
    "https://www.w3.org/ns/activitystreams#Public",
    "https://example.org/users/someone/followers"
  ],
  "object": "https://boobs.example.com/users/boobslover6969/statuses/01KFGBRPFJJCP0FM4039H4VDVQ",
  "target": "https://example.org/users/someone/statuses/01KFGBJN586W5EZFHE16R7FMJ1",
  "result": "https://example.org/users/someone/authorizations/01KFGBRQBK8PPCQBE7K7Z22J5T"
}
```

In this alternative `Accept`of an impolite interaction, the `object` is the URI of the interaction (in this case, a reply), and the `target` is the post being interacted with. **Importantly, the `Accept` is addressed `cc` the followers collection of the interactee.**

When sent to the followers of `@someone@example.org`, this `Accept` provides enough information to those followers that their servers can validate that a reply from `@boobslover6969@boobs.example.com` has been accepted. They can then dereference the reply `https://boobs.example.com/users/boobslover6969/statuses/01KFGBRPFJJCP0FM4039H4VDVQ` and store it with an authorization URI of `https://example.org/users/someone/authorizations/01KFGBRQBK8PPCQBE7K7Z22J5T`.

By addressing `Accept`s of impolite interaction requests cc followers, implementers can ensure that parts of conversations from servers that are not interaction controls aware can still be shown on servers that *are* aware of and enforce interaction controls.

### Optional Inlining for Responses to Impolite Interactions

The following optional inlining may be done on `Accept`s and `Reject`s sent out in response to impolite interaction requests.

#### Type hinting: inline `object` property on `Accept` and `Reject`

If desired, implementers may partially expand/inline the `object` property of an `Accept` or `Reject` to hint to remote servers about the type of interaction being `Accept`ed or `Reject`ed. When inlining in this way, the `object`'s `type` and `id` must be defined at a minimum. For example:

```json
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "actor": "https://example.org/users/post_author",
  "cc": [
    "https://www.w3.org/ns/activitystreams#Public",
    "https://example.org/users/post_author/followers"
  ],
  "to": "https://somewhere.else.example.org/users/someone",
  "id": "https://example.org/users/post_author/activities/reject/01J0K2YXP9QCT5BE1JWQSAM3B6",
  "object": {
    "type": "Note",
    "id": "https://somewhere.else.example.org/users/someone/statuses/01J17XY2VXGMNNPH1XR7BG2524",
    [...]
  },
  "result": "https://example.org/users/post_author/approvals/01JMPS01E54DG9JCF2ZK3JDMXE",
  "type": "Accept"
}
```

#### Set `target` property on `Accept` and `Reject`

If desired, implementers may set the `target` property on outgoing `Accept` or `Reject` activities to the `id` of the post being interacted with, to make it easier for remote servers to understand the shape and relevance of the interaction that's being `Accept`ed or `Reject`ed.

For example, the following json object `Accept`s the attempt of `@someone@somewhere.else.example.org` to reply to a post by `@post_author@example.org` that has the id `https://example.org/users/post_author/statuses/01JJYV141Y5M4S65SC1XCP65NT`:

```json
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "actor": "https://example.org/users/post_author",
  "cc": [
    "https://www.w3.org/ns/activitystreams#Public",
    "https://example.org/users/post_author/followers"
  ],
  "to": "https://somewhere.else.example.org/users/someone",
  "id": "https://example.org/users/post_author/activities/reject/01J0K2YXP9QCT5BE1JWQSAM3B6",
  "object": "https://somewhere.else.example.org/users/someone/statuses/01J17XY2VXGMNNPH1XR7BG2524",
  "target": "https://example.org/users/post_author/statuses/01JJYV141Y5M4S65SC1XCP65NT",
  "result": "https://example.org/users/post_author/approvals/01JMPS01E54DG9JCF2ZK3JDMXE",
  "type": "Accept"
}
```

If desired, the `target` property can also be partially expanded/inlined to type hint about the post that was interacted with. When inlining in this way, the `target`'s `type` and `id` must be defined at a minimum. For example:

```json
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "actor": "https://example.org/users/post_author",
  "cc": [
    "https://www.w3.org/ns/activitystreams#Public",
    "https://example.org/users/post_author/followers"
  ],
  "to": "https://somewhere.else.example.org/users/someone",
  "id": "https://example.org/users/post_author/activities/reject/01J0K2YXP9QCT5BE1JWQSAM3B6",
  "object": "https://somewhere.else.example.org/users/someone/statuses/01J17XY2VXGMNNPH1XR7BG2524",
  "target": {
    "type": "Note",
    "id": "https://example.org/users/post_author/statuses/01JJYV141Y5M4S65SC1XCP65NT"
    [ ... ]
  },
  "result": "https://example.org/users/post_author/approvals/01JMPS01E54DG9JCF2ZK3JDMXE",
  "type": "Accept"
}
```
