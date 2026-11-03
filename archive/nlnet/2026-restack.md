# NLnet Grant Application - Restack 2026

This document is the application on behalf of GoToSocial for funding from the [NLnet](https://nlnet.nl) [Restack fund](https://nlnet.nl/restack/), November 2026. See [here](https://nlnet.nl/propose/).

## Proposal title

GoToSocial – QoL and compatibility

## Project website(s) / repositories

https://codeberg.org/superseriousbusiness/gotosocial / https://docs.gotosocial.org

## Proposed effort

> Please provide a breakdown clarifying how the requested budget would be used.
> ·	Which tasks will you undertake, and (roughly) how much time/effort will each take?
> ·	What hourly/daily/monthly rate did you use?
> ·	Will there be any expenses?
> (max. 4000 characters)

In 2025-2026, we received NLnet funding to improve performance, streamline storage + database requirements, and enhance connectivity of GoToSocial (GtS) with things like relay connections, remote availability monitoring, etc, on the way to a v1.0 release late 2026/early 2027. We would like to use the Restack fund to continue to support post-v1 GtS by building Fediverse features that users and admins want as part of their social media stack (quote posts, link previews, easier account management, etc) while also updating our libraries and APIs to ensure that GoToSocial remains compatible with other Fediverse server implementations using the latest FEPs and RFCs.

We are asking €100k for two full-time developers for one year, ie, €50k each per year, or roughly €27 each per hour. After taxes, this comes to about €2400 each per month. This is the same as what we received for our 2025 application, and represents a liveable developer wage. We will use this to cover our living costs: gas, electricity, internet, food, etc. The funding also gives enough breathing room to do "invisible" work like bug fixes, optimizations, refactoring, and release coordination.

Proposed tasks:

- Mastodon-style quote posts: allow GtS users to receive, interact with, and create privacy-respecting quotes. ~10 weeks.
- OpenGraph link previews: call links in posts, generate preview cards using OpenGraph and other meta tags, store previews in the database, and serve them to users via the client API + web frontend. 6-8 weeks.
- Design and implement authorized push, a security mechanism both related and opposite to the existing authorized fetch. While authorized fetch prevents bad actors *retrieving* data, authorized push would prevent bad actors from *being sent* data, for example after a domain renewal lapse.
- Account portability via ActivityPub JSON data export: allow GtS users to export their account data to archive + ActivityPub JSON format, for backup and/or ingestion by other ActivityPub services. 8-10 weeks.
- Local post autodeletion: allow local users to configure autodeletion of their posts after x amount of time, for privacy + hygiene purposes. ~8 weeks.
- HTTP Signature library rewrite: to ensure continued compatibility of GtS in the fedi ecosystem, GtS's HTTP signature library needs rewriting to support the latest HTTP signature RFC. (Other projects written in Go will also be able to benefit from this rewrite!) ~10 weeks.
- Instance-wide announcements: allow admins to create and pin announcements to all users on an instance, and if desired email those announcement to users as well. ~6 weeks.
- "Likes" and "Shares" counts for federated posts: send and ingest `Likes` and `Shares` collections with `totalItems` count on posts, serve these stats via the API to avoid the appearance of "dead fedi". ~6 weeks.
- Accept remote username changes: for forward compatibility with softwares like Akkoma that (plan to) allow username changes, GtS needs to lessen reliance on unchanging `preferredUsername` in application logic. ~8 weeks.
- Instance webring support for better discoverability of instances: allow GtS instance admins to belong to webrings of instances by fetching webring list + showing next/prev instance links on their instance. 6-8 weeks.
- Account management improvements: self-serve password reset and account deletion; admin account creation, confirmation, and updating; unsuspend previously suspended accounts. ~8 weeks.
- Improved delivery failure rescheduling: fix existing queue persistence issues and add longer-term delivery attempts. ~8 weeks.

Taken together, these tasks represent a large push forward for GoToSocial features, bringing it closer to parity with more "full-featured" implementations like Mastodon, while we continue to "optimize downwards" to allow GtS to fit snugly into a resource-constrained open social stack.

> Please compare your own project with other relevant efforts (both existing and historical). What is new, more thorough or otherwise different? Have you considered contributing to any of the others? (max. 4000 characters)

In 2026 ActivityPub continues to be a popular protocol, with a proliferation of AP implementations like Mastodon, Akkoma, Lemmy, Misskey (and its many forks), Pixelfed, WriteFreely, Loops, Wordpress, WAFRN, and more, each with its own focus and purpose. As strong as many of these implementations are, each has its downsides.

For instance, Mastodon and Pixelfed both lean more towards replicating existing corporate social media efforts (X/Bluesky and Instagram respectively), which is offputting to some. These two implementations are also complex to install and maintain, which puts many would-be admins off hosting their own servers, and leads to a concentration of users on developer-run servers, which breaks the promise of decentralization.

Other ActivityPub-enabled microblogging software like Snac2 have fewer moving parts and lower system requirements, which helps people deploy single-user instances. However, their implementation of the Mastodon client API (the de-facto client API of the fediverse, for better or worse), misses some features, and the barebones admin functionality is not user-friendly for people unaccustomed to using the command line.

Misskey and its derived projects, additionally, have a heavy front-end UI which is great when it works smoothly, but can result in a poor experience on low-bandwidth connections or low-powered devices.

GoToSocial's continued aim is to strike a balance between featureful but resource-intensive AP implementations like Mastodon, and lighter and simpler AP implementations like Snac2.

Firstly, GtS is small, easy to deploy, and very well documented. It uses ~300MiB of system memory on average, so it can be deployed on tiny VPSs and single-board computers. People who want to run GoToSocial don't need to install or maintain anything else on the host, as it has no dependencies on things like ffmpeg or on-system SQLite: everything is virtualized inside the single binary using WASM / Wazero. This also makes it very easy for third parties like Yunohost to package and distribute.

Secondly, the admin / settings panel offers admins and users the ability to easily customize their profiles and adjust the way that their instance looks, feels, and federates, and to handle day-to-day admin tasks like reviewing reports, blocking or allowing other instances, etc. The web view of profiles is rendered in simple HTML + CSS; Javascript is not required, but is used for progressive enhancement if available. This makes it feasible to view a GoToSocial profile on a mobile device with poor connection, on a lower-end laptop etc, without fans spinning up.

Thirdly, GtS prioritizes strong safety features such as strict block implementation, always-on HTTP signature authentication, interaction controls, block- and allowlist subscriptions, domain limits, and allowing users to choose what visibility of level of posts can be viewed on their profile (if any), and whether or not reblogs should be shown. These features have been very popular with privacy-conscious Fediverse denizens who want to communicate with groups of friends and don't necessarily want their posts to be made viewable by as wide an audience as possible.

Fourthly, we have been -- and continue to be -- fastidious with our Mastodon client API implementation, which means that GoToSocial can be used with a wide variety of clients that provide different experiences to the user.

> What technical challenges do you expect during the project? (max. 4000 characters)

By now, GoToSocial is a fairly stable and performant project, and Kim and I have a proven track record of implementing features without introducing lots of new bugs and performance issues. Nevertheless, we do foresee some challenges with the set of tasks described above:

- HTTP Signature rewrite: comprehensive tests will be needed to ensure the library is fully compliant with the latest RFC and avoid buggy edge cases. Since not everyone is using the new RFC yet, we foresee having to implement "double knocking" whereby we first try signing/reading signatures using the new RFC, then fall back to the previous more widespread one.
- Remote username changes: we wrote GtS with an assumption that usernames do not change, and modifying our internal application logic to not rely on this assumption will require careful refactoring if we don't want to introduce new bugs.
- Instance webring support: there's a lot to think about here, for example what format should the webring host use to return info, how often should a GtS instance call that host, how to handle revocations, etc.
- Quote posts: happily, we've done a lot of work already in the GtS codebase to implement interaction controls, and quotes are just another form of interaction. Nevertheless, we need to be able to handle things like quote revocations, only displaying quoted posts to authorized users, quotes of quotes, etc.
- Improved delivery failure re-scheduling will require experimentation to a balance between re-attempted delivery over long periods while preventing queue clogging. This will also require more intelligent heuristics for job selection to maintain instance responsiveness.
- Local post autodeletion: we'll have to be wary here when introducing new indexes on the statuses (and statuses_stats) table to not introduce performance regressions: we'll need to design our table schema and indexes carefully to make the autodeleter run smoothly, and not bloat the statuses table by indexing on unnecessary columns.
- Account JSON data export and archiving: this will be an interesting challenge because while the archive gathering will be triggered by the user via the client API (easy enough), we will have to deliver the archive by gathering it together asynchronously in the background, before notifying the user (via email, most likely) that their download is ready. We will have to ensure we don't tank performance while the archiving of media and posts is taking place, and also ensure that the link we return to the user is not exposed to the wider internet, as it will contain very private data!
- Authorized push: we foresee the possibility of having to create a new Activity type for this, eg., a `Ping` activity that can be delivered to a remote instance's inbox, while the pinging instance awaits receipt of a signed `Pong` activity in return (exact technical details tbd). This may require a new FEP, and standards documentation for other ActivityPub servers to implement if they so-wish. If the implementation by necessity causes compatibility issues, it will be configuration gated similar to Mastodon's authorized fetch rollout.

> What is the ecosystem around your own project, and how will you engage relevant actors? For example, does your project depend on other projects/people? And who would be the main users? (max. 2000 characters)

The main users of GoToSocial are privacy-conscious Fediverse folks who want to deploy their own small or single-user Fedi instances with the least amount of fuss and expense. We aim to provide stable, featureful, customizable, reliable, and secure software to these people, as part of our goal of diversifying the Fediverse ecosystem, and making it simple to get started with self-hosting an instance. Last time we checked around 2024/2025, there were about 2000 or so GoToSocial instances on the Fediverse, but it is difficult to accurately measure this number as we do not use phone-home telemetry in GoToSocial, and we provide opt-in features that make it difficult for crawlers to monitor GoToSocial instances reliably.

With regard to other ActivityPub implementations on the Fediverse, we spend a lot of time making sure that GoToSocial is compatible on a protocol-level with various Fediverse projects, reporting bugs to maintainers and fixing our bugs where necessary. In that sense, we are dependent on collaboration with the wider Fediverse, and we contribute back to it by trying to make GoToSocial a "good citizen" in terms of federation. On that note, in 2026 we are working with Mastodon on their FASP work for Remote Media Storage, and we've previously worked alongside Mastodon on finalizing interaction policy protocol details.

For GoToSocial-specific extensions to ActivityPub, we've diligently documented what we've done so far, and exposed a GoToSocial namespace so that remote softwares can easily incorporate GtS extensions if they want to: https://docs.gotosocial.org/en/latest/federation/interaction_controls/, https://gotosocial.org/ns. Thanks to this, other projects have now begun to adopt GoToSocial-style interaction policies (WAFRN, Pixelfed and Loops, Bonfire). In 2027, we intend to continue our work to get GoToSocial's interaction policies codified as a W3 recommendation, so that a wide spread of fedi projects can continue to benefit from them.

## Background

> What makes you the right person(s) for this effort? Do you have relevant experience in this or other projects/organisations? And if so, can you tell us a bit about your contributions? (max. 2000 characters)

Kim and I have been working professionally on GoToSocial for the last 3 years and 5 years respectively (since the infancy of the project). Over these years we have put many thousands of hours of work into GtS: writing features code and documentation, fixing bugs, communicating with contributors and doing code review, refactoring, deploying infrastructure for project builds and discussion, doing project planning, and answering user questions. We have shown over the years that we work brilliantly well together as a team, and given our experience and track record we're 100% confident that we can deliver on the tasks above.

Aside from our work on GoToSocial, we also maintain a fork of the Activity library (https://codeberg.org/superseriousbusiness/activity), a fork of the standalone Mastodon frontend customized for GoToSocial (https://codeberg.org/superseriousbusiness/masto-fe-standalone, deployed at https://masto-fe.superseriousbusiness.org), and Kim in particular maintains a large amount of libraries used by the project (https://codeberg.org/gruf), particularly go-ffmpreg (https://codeberg.org/gruf/go-ffmpreg).

> Does (or did) the project have other funding sources? (max. 1000 characters)

GoToSocial has received three NLnet grants previously, one in 2022-2023, for 50k EUR, another one in 2024-2025, for 75k EUR, and one in 2025-2026, for 100k EUR. These grants went towards paying (rising) living costs for myself and Kim (rent, groceries, utilities, taxes, etc), while we both worked full time on the project. In addition to NLnet as our main funding source, GoToSocial has a crowdfunding page on OpenCollective (https://opencollective.com/gotosocial) that currently makes about 1200 EUR per month, and we receive odd donations on Liberapay (https://liberapay.com/gotosocial) averaging about 50 EUR / week.
