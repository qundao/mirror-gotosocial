// THIS IS A GENERATED FILE, DO NOT EDIT BY HAND
// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package config

import (
	"fmt"
	"time"

	"code.superseriousbusiness.org/gotosocial/internal/language"
	"codeberg.org/gruf/go-bytesize"
	"github.com/spf13/cast"
	"github.com/spf13/pflag"
)

const (
	DatabasePostgresPortFlag                      = "port"
	DatabasePostgresUserFlag                      = "user"
	DatabasePostgresPasswordFlag                  = "password"
	DatabasePostgresDatabaseFlag                  = "database"
	DatabasePostgresTLSModeFlag                   = "tls-mode"
	DatabasePostgresTLSCACertFlag                 = "tls-ca-cert"
	DatabasePostgresConnectionStringFlag          = "postgres-connection-string"
	DatabaseSQLiteJournalModeFlag                 = "db-sqlite-journal-mode"
	DatabaseSQLiteSynchronousFlag                 = "db-sqlite-synchronous"
	DatabaseSQLiteCacheSizeFlag                   = "db-sqlite-cache-size"
	DatabaseSQLiteBusyTimeoutFlag                 = "db-sqlite-busy-timeout"
	DatabaseTypeFlag                              = "db-type"
	DatabaseAddressFlag                           = "db-address"
	DatabaseMaxOpenConnsMultiplierFlag            = "db-max-open-conns-multiplier"
	AdvancedRateLimitRequestsFlag                 = "advanced-rate-limit-requests"
	AdvancedRateLimitExceptionsFlag               = "advanced-rate-limit-exceptions"
	AdvancedThrottlingMultiplierFlag              = "advanced-throttling-multiplier"
	AdvancedThrottlingRetryAfterFlag              = "advanced-throttling-retry-after"
	AdvancedCookiesSamesiteFlag                   = "advanced-cookies-samesite"
	AdvancedSenderMultiplierFlag                  = "advanced-sender-multiplier"
	AdvancedCSPExtraURIsFlag                      = "advanced-csp-extra-uris"
	AdvancedHeaderFilterModeFlag                  = "advanced-header-filter-mode"
	HTTPServerMaxMultipartMemoryFlag              = "http-server-max-multipart-memory"
	HTTPServerUseH2CFlag                          = "http-server-use-h2c"
	HTTPServerReadTimeoutFlag                     = "http-server-read-timeout"
	HTTPServerReadHeaderTimeoutFlag               = "http-server-read-header-timeout"
	HTTPServerWriteTimeoutFlag                    = "http-server-write-timeout"
	HTTPServerIdleTimeoutFlag                     = "http-server-idle-timeout"
	HTTPServerMaxHeaderBytesFlag                  = "http-server-max-header-bytes"
	HTTPServerMaxConcurrentStreamsFlag            = "http-server-max-concurrent-streams"
	HTTPServerMaxDecoderHeaderTableSizeFlag       = "http-server-max-decoder-header-table-size"
	HTTPServerMaxEncoderHeaderTableSizeFlag       = "http-server-max-encoder-header-table-size"
	HTTPServerMaxReadFrameSizeFlag                = "http-server-max-read-frame-size"
	HTTPServerMaxReceiveBufferPerConnectionFlag   = "http-server-max-receive-buffer-per-connection"
	HTTPServerMaxReceiveBufferPerStreamFlag       = "http-server-max-receive-buffer-per-stream"
	HTTPServerSendPingTimeoutFlag                 = "http-server-send-ping-timeout"
	HTTPServerPingTimeoutFlag                     = "http-server-ping-timeout"
	HTTPServerWriteByteTimeoutFlag                = "http-server-write-byte-timeout"
	HTTPClientAllowIPsFlag                        = "http-client-allow-ips"
	HTTPClientBlockIPsFlag                        = "http-client-block-ips"
	HTTPClientTimeoutFlag                         = "http-client-timeout"
	HTTPClientTLSInsecureSkipVerifyFlag           = "http-client-tls-insecure-skip-verify"
	HTTPClientInsecureOutgoingFlag                = "http-client-insecure-outgoing"
	HTTPClientDisableKeepAlivesFlag               = "http-client-disable-keep-alives"
	HTTPClientMaxIdleConnsFlag                    = "http-client-max-idle-conns"
	HTTPClientMaxIdleConnsPerHostFlag             = "http-client-max-idle-conns-per-host"
	HTTPClientMaxConnsPerHostFlag                 = "http-client-max-conns-per-host"
	HTTPClientIdleConnTimeoutFlag                 = "http-client-idle-conn-timeout"
	HTTPClientTLSHandshakeTimeoutFlag             = "http-client-tls-handshake-timeout"
	HTTPClientResponseHeaderTimeoutFlag           = "http-client-response-header-timeout"
	HTTPClientReadBufferSizeFlag                  = "http-client-read-buffer-size"
	HTTPClientWriteBufferSizeFlag                 = "http-client-write-buffer-size"
	MediaDescriptionMinCharsFlag                  = "media-description-min-chars"
	MediaDescriptionMaxCharsFlag                  = "media-description-max-chars"
	MediaEmojiLocalMaxSizeFlag                    = "media-emoji-local-max-size"
	MediaEmojiRemoteMaxSizeFlag                   = "media-emoji-remote-max-size"
	MediaImageSizeHintFlag                        = "media-image-size-hint"
	MediaVideoSizeHintFlag                        = "media-video-size-hint"
	MediaLocalMaxSizeFlag                         = "media-local-max-size"
	MediaRemoteMaxSizeFlag                        = "media-remote-max-size"
	MediaFfmpegPoolSizeFlag                       = "media-ffmpeg-pool-size"
	MediaThumbMaxPixelsFlag                       = "media-thumb-max-pixels"
	CacheS3ObjectInfoFlag                         = "cache-s3-object-info"
	CacheHomeTimelineSizeFlag                     = "cache-home-timeline-size"
	CacheListTimelineSizeFlag                     = "cache-list-timeline-size"
	CacheTagTimelineSizeFlag                      = "cache-tag-timeline-size"
	CacheNotificationTimelineSizeFlag             = "cache-notification-timeline-size"
	CacheHomeTimelineTimeoutFlag                  = "cache-home-timeline-timeout"
	CacheListTimelineTimeoutFlag                  = "cache-list-timeline-timeout"
	CacheTagTimelineTimeoutFlag                   = "cache-tag-timeline-timeout"
	CacheNotificationTimelineTimeoutFlag          = "cache-notification-timeline-timeout"
	CacheMemoryTargetFlag                         = "cache-memory-target"
	CacheAccountMemRatioFlag                      = "cache-account-mem-ratio"
	CacheAccountNoteMemRatioFlag                  = "cache-account-note-mem-ratio"
	CacheAccountSettingsMemRatioFlag              = "cache-account-settings-mem-ratio"
	CacheAccountStatsMemRatioFlag                 = "cache-account-stats-mem-ratio"
	CacheApplicationMemRatioFlag                  = "cache-application-mem-ratio"
	CacheBlockMemRatioFlag                        = "cache-block-mem-ratio"
	CacheBlockIDsMemRatioFlag                     = "cache-block-ids-mem-ratio"
	CacheBoostOfIDsMemRatioFlag                   = "cache-boost-of-ids-mem-ratio"
	CacheClientMemRatioFlag                       = "cache-client-mem-ratio"
	CacheConversationMemRatioFlag                 = "cache-conversation-mem-ratio"
	CacheConversationLastStatusIDsMemRatioFlag    = "cache-conversation-last-status-ids-mem-ratio"
	CacheDomainPermissionDraftMemRatioFlag        = "cache-domain-permission-draft-mem-ratio"
	CacheDomainLimitMemRatioFlag                  = "cache-domain-permission-limit-mem-ratio"
	CacheDomainPermissionSubscriptionMemRatioFlag = "cache-domain-permission-subscription-mem-ratio"
	CacheEmojiMemRatioFlag                        = "cache-emoji-mem-ratio"
	CacheEmojiCategoryMemRatioFlag                = "cache-emoji-category-mem-ratio"
	CacheFederationErrorMemRatioFlag              = "cache-federation-error-mem-ratio"
	CacheFilterMemRatioFlag                       = "cache-filter-mem-ratio"
	CacheFilterIDsMemRatioFlag                    = "cache-filter-ids-mem-ratio"
	CacheFilterKeywordMemRatioFlag                = "cache-filter-keyword-mem-ratio"
	CacheFilterStatusMemRatioFlag                 = "cache-filter-status-mem-ratio"
	CacheFollowMemRatioFlag                       = "cache-follow-mem-ratio"
	CacheFollowIDsMemRatioFlag                    = "cache-follow-ids-mem-ratio"
	CacheFollowRequestMemRatioFlag                = "cache-follow-request-mem-ratio"
	CacheFollowRequestIDsMemRatioFlag             = "cache-follow-request-ids-mem-ratio"
	CacheFollowingTagIDsMemRatioFlag              = "cache-following-tag-ids-mem-ratio"
	CacheHomeAccountIDsMemRatioFlag               = "cache-home-account-ids-mem-ratio"
	CacheInReplyToIDsMemRatioFlag                 = "cache-in-reply-to-ids-mem-ratio"
	CacheInstanceMemRatioFlag                     = "cache-instance-mem-ratio"
	CacheInteractionRequestMemRatioFlag           = "cache-interaction-request-mem-ratio"
	CacheListMemRatioFlag                         = "cache-list-mem-ratio"
	CacheListIDsMemRatioFlag                      = "cache-list-ids-mem-ratio"
	CacheListedIDsMemRatioFlag                    = "cache-listed-ids-mem-ratio"
	CacheMarkerMemRatioFlag                       = "cache-marker-mem-ratio"
	CacheMediaMemRatioFlag                        = "cache-media-mem-ratio"
	CacheMentionMemRatioFlag                      = "cache-mention-mem-ratio"
	CacheMoveMemRatioFlag                         = "cache-move-mem-ratio"
	CacheNotificationMemRatioFlag                 = "cache-notification-mem-ratio"
	CachePollMemRatioFlag                         = "cache-poll-mem-ratio"
	CachePollVoteMemRatioFlag                     = "cache-poll-vote-mem-ratio"
	CachePollVoteIDsMemRatioFlag                  = "cache-poll-vote-ids-mem-ratio"
	CacheReportMemRatioFlag                       = "cache-report-mem-ratio"
	CacheRelayActorMemRatioFlag                   = "cache-relay-actor-mem-ratio"
	CacheRelayMatcherMemRatioFlag                 = "cache-relay-matcher-mem-ratio"
	CacheRelayPushMemRatioFlag                    = "cache-relay-push-mem-ratio"
	CacheRelayPushIDsMemRatioFlag                 = "cache-relay-push-ids-mem-ratio"
	CacheRelaySubscriptionMemRatioFlag            = "cache-relay-subscription-mem-ratio"
	CacheScheduledStatusMemRatioFlag              = "cache-scheduled-status-mem-ratio"
	CacheSinBinStatusMemRatioFlag                 = "cache-sin-bin-status-mem-ratio"
	CacheStatusMemRatioFlag                       = "cache-status-mem-ratio"
	CacheStatusBookmarkMemRatioFlag               = "cache-status-bookmark-mem-ratio"
	CacheStatusBookmarkIDsMemRatioFlag            = "cache-status-bookmark-ids-mem-ratio"
	CacheStatusEditMemRatioFlag                   = "cache-status-edit-mem-ratio"
	CacheStatusFaveMemRatioFlag                   = "cache-status-fave-mem-ratio"
	CacheStatusFaveIDsMemRatioFlag                = "cache-status-fave-ids-mem-ratio"
	CacheStatusPinnedIDsMemRatioFlag              = "cache-status-pinned-ids-mem-ratio"
	CacheTagMemRatioFlag                          = "cache-tag-mem-ratio"
	CacheThreadMuteMemRatioFlag                   = "cache-thread-mute-mem-ratio"
	CacheTokenMemRatioFlag                        = "cache-token-mem-ratio"
	CacheTombstoneMemRatioFlag                    = "cache-tombstone-mem-ratio"
	CacheUserMemRatioFlag                         = "cache-user-mem-ratio"
	CacheUserMuteMemRatioFlag                     = "cache-user-mute-mem-ratio"
	CacheUserMuteIDsMemRatioFlag                  = "cache-user-mute-ids-mem-ratio"
	CacheWebfingerMemRatioFlag                    = "cache-webfinger-mem-ratio"
	CacheWebPushSubscriptionMemRatioFlag          = "cache-web-push-subscription-mem-ratio"
	CacheWebPushSubscriptionIDsMemRatioFlag       = "cache-web-push-subscription-ids-mem-ratio"
	CacheMutesMemRatioFlag                        = "cache-mutes-mem-ratio"
	CacheStatusFilterMemRatioFlag                 = "cache-status-filter-mem-ratio"
	CacheVisibilityMemRatioFlag                   = "cache-visibility-mem-ratio"
	LogLevelFlag                                  = "log-level"
	LogFormatFlag                                 = "log-format"
	LogTimestampFormatFlag                        = "log-timestamp-format"
	LogDbQueriesFlag                              = "log-db-queries"
	LogClientIPFlag                               = "log-client-ip"
	RequestIDHeaderFlag                           = "request-id-header"
	ConfigPathFlag                                = "config-path"
	ApplicationNameFlag                           = "application-name"
	LandingPageUserFlag                           = "landing-page-user"
	HostFlag                                      = "host"
	AccountDomainFlag                             = "account-domain"
	ProtocolFlag                                  = "protocol"
	BindAddressFlag                               = "bind-address"
	PortFlag                                      = "port"
	TrustedProxiesFlag                            = "trusted-proxies"
	SoftwareVersionFlag                           = "software-version"
	WebTemplateBaseDirFlag                        = "web-template-base-dir"
	WebAssetBaseDirFlag                           = "web-asset-base-dir"
	InstanceFederationModeFlag                    = "instance-federation-mode"
	InstanceFederationSpamFilterFlag              = "instance-federation-spam-filter"
	InstanceExposePeersFlag                       = "instance-expose-peers"
	InstanceExposeBlocklistFlag                   = "instance-expose-blocklist"
	InstanceExposeBlocklistWebFlag                = "instance-expose-blocklist-web"
	InstanceExposeAllowlistFlag                   = "instance-expose-allowlist"
	InstanceExposeAllowlistWebFlag                = "instance-expose-allowlist-web"
	InstanceExposePublicTimelineFlag              = "instance-expose-public-timeline"
	InstanceExposeCustomEmojisFlag                = "instance-expose-custom-emojis"
	InstanceDirectoryModeFlag                     = "instance-directory-mode"
	InstanceDeliverToSharedInboxesFlag            = "instance-deliver-to-shared-inboxes"
	InstanceInjectMastodonVersionFlag             = "instance-inject-mastodon-version"
	InstanceLanguagesFlag                         = "instance-languages"
	InstanceStatsModeFlag                         = "instance-stats-mode"
	InstanceAllowBackdatingStatusesFlag           = "instance-allow-backdating-statuses"
	InstanceRobotsAllowIndexingFlag               = "instance-robots-allow-indexing"
	AccountsRegistrationOpenFlag                  = "accounts-registration-open"
	AccountsReasonRequiredFlag                    = "accounts-reason-required"
	AccountsRegistrationDailyLimitFlag            = "accounts-registration-daily-limit"
	AccountsRegistrationBacklogLimitFlag          = "accounts-registration-backlog-limit"
	AccountsAllowCustomCSSFlag                    = "accounts-allow-custom-css"
	AccountsCustomCSSLengthFlag                   = "accounts-custom-css-length"
	AccountsMaxProfileFieldsFlag                  = "accounts-max-profile-fields"
	StorageBackendFlag                            = "storage-backend"
	StorageLocalBasePathFlag                      = "storage-local-base-path"
	StorageS3EndpointFlag                         = "storage-s3-endpoint"
	StorageS3AccessKeyFlag                        = "storage-s3-access-key"
	StorageS3SecretKeyFlag                        = "storage-s3-secret-key"
	StorageS3UseSSLFlag                           = "storage-s3-use-ssl"
	StorageS3BucketNameFlag                       = "storage-s3-bucket"
	StorageS3ProxyFlag                            = "storage-s3-proxy"
	StorageS3RedirectURLFlag                      = "storage-s3-redirect-url"
	StorageS3BucketLookupFlag                     = "storage-s3-bucket-lookup"
	StorageS3KeyPrefixFlag                        = "storage-s3-key-prefix"
	StorageS3RegionFlag                           = "storage-s3-region"
	StatusesMaxCharsFlag                          = "statuses-max-chars"
	StatusesPollMaxOptionsFlag                    = "statuses-poll-max-options"
	StatusesPollOptionMaxCharsFlag                = "statuses-poll-option-max-chars"
	StatusesMediaMaxFilesFlag                     = "statuses-media-max-files"
	ScheduledStatusesMaxTotalFlag                 = "scheduled-statuses-max-total"
	ScheduledStatusesMaxDailyFlag                 = "scheduled-statuses-max-daily"
	LetsEncryptEnabledFlag                        = "letsencrypt-enabled"
	LetsEncryptPortFlag                           = "letsencrypt-port"
	LetsEncryptCertDirFlag                        = "letsencrypt-cert-dir"
	LetsEncryptEmailAddressFlag                   = "letsencrypt-email-address"
	TLSCertificateChainFlag                       = "tls-certificate-chain"
	TLSCertificateKeyFlag                         = "tls-certificate-key"
	OIDCEnabledFlag                               = "oidc-enabled"
	OIDCIdpNameFlag                               = "oidc-idp-name"
	OIDCSkipVerificationFlag                      = "oidc-skip-verification"
	OIDCIssuerFlag                                = "oidc-issuer"
	OIDCClientIDFlag                              = "oidc-client-id"
	OIDCClientSecretFlag                          = "oidc-client-secret"
	OIDCScopesFlag                                = "oidc-scopes"
	OIDCLinkExistingFlag                          = "oidc-link-existing"
	OIDCAllowedGroupsFlag                         = "oidc-allowed-groups"
	OIDCAdminGroupsFlag                           = "oidc-admin-groups"
	TracingEnabledFlag                            = "tracing-enabled"
	MetricsEnabledFlag                            = "metrics-enabled"
	SMTPHostFlag                                  = "smtp-host"
	SMTPPortFlag                                  = "smtp-port"
	SMTPUsernameFlag                              = "smtp-username"
	SMTPPasswordFlag                              = "smtp-password"
	SMTPFromFlag                                  = "smtp-from"
	SMTPFromDisplayNameFlag                       = "smtp-from-display-name"
	SMTPDiscloseRecipientsFlag                    = "smtp-disclose-recipients"
	SyslogEnabledFlag                             = "syslog-enabled"
	SyslogProtocolFlag                            = "syslog-protocol"
	SyslogAddressFlag                             = "syslog-address"
	SyslogMirrorFlag                              = "syslog-mirror"
	SyslogMsgLengthFlag                           = "syslog-msg-length"
	AdminAccountUsernameFlag                      = "username"
	AdminAccountEmailFlag                         = "email"
	AdminAccountPasswordFlag                      = "password"
	AdminTransPathFlag                            = "path"
	AdminMediaPruneDryRunFlag                     = "dry-run"
	AdminMediaListLocalOnlyFlag                   = "local-only"
	AdminMediaListRemoteOnlyFlag                  = "remote-only"
	TestrigSkipDBSetupFlag                        = "skip-db-setup"
	TestrigSkipDBTeardownFlag                     = "skip-db-teardown"
)

func (cfg *PostgresConfiguration) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	flags.Uint16(joinFlag(prefix, "port"), cfg.Port, "Database port")
	flags.String(joinFlag(prefix, "user"), cfg.User, "Database username")
	flags.String(joinFlag(prefix, "password"), cfg.Password, "Database password")
	flags.String(joinFlag(prefix, "database"), cfg.Database, "Database name")
	flags.String(joinFlag(prefix, "tls-mode"), cfg.TLSMode, "Database tls mode")
	flags.String(joinFlag(prefix, "tls-ca-cert"), cfg.TLSCACert, "Path to CA cert for db tls connection")
	flags.String(joinFlag(prefix, "postgres-connection-string"), cfg.ConnectionString, "Full Database URL for connection to postgres")
}

func (cfg *SQLiteConfiguration) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	flags.String(joinFlag(prefix, "journal-mode"), cfg.JournalMode, "Sqlite only: see https://www.sqlite.org/pragma.html#pragma_journal_mode")
	flags.String(joinFlag(prefix, "synchronous"), cfg.Synchronous, "Sqlite only: see https://www.sqlite.org/pragma.html#pragma_synchronous")
	flags.String(joinFlag(prefix, "cache-size"), cfg.CacheSize.String(), "Sqlite only: see https://www.sqlite.org/pragma.html#pragma_cache_size")
	flags.Duration(joinFlag(prefix, "busy-timeout"), cfg.BusyTimeout, "Sqlite only: see https://www.sqlite.org/pragma.html#pragma_busy_timeout")
}

func (cfg *DatabaseConfiguration) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	cfg.Postgres.RegisterFlags(joinFlag(prefix, ""), flags)
	cfg.SQLite.RegisterFlags(joinFlag(prefix, "sqlite"), flags)
	flags.String(joinFlag(prefix, "type"), cfg.Type, "Database type: eg., postgres")
	flags.String(joinFlag(prefix, "address"), cfg.Address, "Database ipv4 address, hostname, or filename")
	flags.Int(joinFlag(prefix, "max-open-conns-multiplier"), cfg.MaxOpenConnsMultiplier, "Multiplier to use per cpu for max open database connections. 0 or less is normalized to 1.")
}

func (cfg *RateLimitConfig) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	flags.Int(joinFlag(prefix, "requests"), cfg.Requests, "Amount of HTTP requests to permit within a 5 minute window. 0 or less turns rate limiting off.")
	flags.StringSlice(joinFlag(prefix, "exceptions"), cfg.Exceptions.Strings(), "Slice of CIDRs to exclude from rate limit restrictions.")
}

func (cfg *ThrottlingConfig) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	flags.Int(joinFlag(prefix, "multiplier"), cfg.Multiplier, "Multiplier to use per cpu for http request throttling. 0 or less turns throttling off.")
	flags.Duration(joinFlag(prefix, "retry-after"), cfg.RetryAfter, "Retry-After duration response to send for throttled requests.")
}

func (cfg *AdvancedConfig) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	cfg.RateLimit.RegisterFlags(joinFlag(prefix, "rate-limit"), flags)
	cfg.Throttling.RegisterFlags(joinFlag(prefix, "throttling"), flags)
	flags.String(joinFlag(prefix, "cookies-samesite"), cfg.CookiesSamesite, "'strict' or 'lax', see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie/SameSite")
	flags.Int(joinFlag(prefix, "sender-multiplier"), cfg.SenderMultiplier, "Multiplier to use per cpu for batching outgoing fedi messages. 0 or less turns batching off (not recommended).")
	flags.StringSlice(joinFlag(prefix, "csp-extra-uris"), cfg.CSPExtraURIs, "Additional URIs to allow when building content-security-policy for media + images.")
	flags.String(joinFlag(prefix, "header-filter-mode"), cfg.HeaderFilterMode, "Set incoming request header filtering mode.")
}

func (cfg *HTTPServerConfiguration) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	flags.String(joinFlag(prefix, "max-multipart-memory"), cfg.MaxMultipartMemory.String(), "")
	flags.Bool(joinFlag(prefix, "use-h2c"), cfg.UseH2C, "")
	flags.Duration(joinFlag(prefix, "read-timeout"), cfg.ReadTimeout, "")
	flags.Duration(joinFlag(prefix, "read-header-timeout"), cfg.ReadHeaderTimeout, "")
	flags.Duration(joinFlag(prefix, "write-timeout"), cfg.WriteTimeout, "")
	flags.Duration(joinFlag(prefix, "idle-timeout"), cfg.IdleTimeout, "")
	flags.String(joinFlag(prefix, "max-header-bytes"), cfg.MaxHeaderBytes.String(), "")
	flags.Int(joinFlag(prefix, "max-concurrent-streams"), cfg.MaxConcurrentStreams, "")
	flags.String(joinFlag(prefix, "max-decoder-header-table-size"), cfg.MaxDecoderHeaderTableSize.String(), "")
	flags.String(joinFlag(prefix, "max-encoder-header-table-size"), cfg.MaxEncoderHeaderTableSize.String(), "")
	flags.String(joinFlag(prefix, "max-read-frame-size"), cfg.MaxReadFrameSize.String(), "")
	flags.String(joinFlag(prefix, "max-receive-buffer-per-connection"), cfg.MaxReceiveBufferPerConnection.String(), "")
	flags.String(joinFlag(prefix, "max-receive-buffer-per-stream"), cfg.MaxReceiveBufferPerStream.String(), "")
	flags.Duration(joinFlag(prefix, "send-ping-timeout"), cfg.SendPingTimeout, "")
	flags.Duration(joinFlag(prefix, "ping-timeout"), cfg.PingTimeout, "")
	flags.Duration(joinFlag(prefix, "write-byte-timeout"), cfg.WriteByteTimeout, "")
}

func (cfg *HTTPClientConfiguration) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	flags.StringSlice(joinFlag(prefix, "allow-ips"), cfg.AllowIPs.Strings(), "")
	flags.StringSlice(joinFlag(prefix, "block-ips"), cfg.BlockIPs.Strings(), "")
	flags.Duration(joinFlag(prefix, "timeout"), cfg.Timeout, "")
	flags.Bool(joinFlag(prefix, "tls-insecure-skip-verify"), cfg.TLSInsecureSkipVerify, "")
	flags.Bool(joinFlag(prefix, "insecure-outgoing"), cfg.InsecureOutgoing, "")
	flags.Bool(joinFlag(prefix, "disable-keep-alives"), cfg.DisableKeepAlives, "")
	flags.Int(joinFlag(prefix, "max-idle-conns"), cfg.MaxIdleConns, "")
	flags.Int(joinFlag(prefix, "max-idle-conns-per-host"), cfg.MaxIdleConnsPerHost, "")
	flags.Int(joinFlag(prefix, "max-conns-per-host"), cfg.MaxConnsPerHost, "")
	flags.Duration(joinFlag(prefix, "idle-conn-timeout"), cfg.IdleConnTimeout, "")
	flags.Duration(joinFlag(prefix, "tls-handshake-timeout"), cfg.TLSHandshakeTimeout, "")
	flags.Duration(joinFlag(prefix, "response-header-timeout"), cfg.ResponseHeaderTimeout, "")
	flags.String(joinFlag(prefix, "read-buffer-size"), cfg.ReadBufferSize.String(), "")
	flags.String(joinFlag(prefix, "write-buffer-size"), cfg.WriteBufferSize.String(), "")
}

func (cfg *MediaConfiguration) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	flags.Int(joinFlag(prefix, "description-min-chars"), cfg.DescriptionMinChars, "Min required chars for an image description")
	flags.Int(joinFlag(prefix, "description-max-chars"), cfg.DescriptionMaxChars, "Max permitted chars for an image description")
	flags.String(joinFlag(prefix, "emoji-local-max-size"), cfg.EmojiLocalMaxSize.String(), "Max size in bytes of emojis uploaded to this instance via the admin API.")
	flags.String(joinFlag(prefix, "emoji-remote-max-size"), cfg.EmojiRemoteMaxSize.String(), "Max size in bytes of emojis to download from other instances.")
	flags.String(joinFlag(prefix, "image-size-hint"), cfg.ImageSizeHint.String(), "Size in bytes of max image size referred to on /api/v_/instance endpoints (else, local max size)")
	flags.String(joinFlag(prefix, "video-size-hint"), cfg.VideoSizeHint.String(), "Size in bytes of max video size referred to on /api/v_/instance endpoints (else, local max size)")
	flags.String(joinFlag(prefix, "local-max-size"), cfg.LocalMaxSize.String(), "Max size in bytes of media uploaded to this instance via API")
	flags.String(joinFlag(prefix, "remote-max-size"), cfg.RemoteMaxSize.String(), "Max size in bytes of media to download from other instances")
	flags.Int(joinFlag(prefix, "ffmpeg-pool-size"), cfg.FfmpegPoolSize, "Number of instances of the embedded ffmpeg WASM binary to add to the media processing pool. 0 or less uses GOMAXPROCS.")
	flags.Int(joinFlag(prefix, "thumb-max-pixels"), cfg.ThumbMaxPixels, "Max size in pixels of any one dimension of a thumbnail (as input media ratio is preserved).")
}

func (cfg *CacheConfiguration) RegisterFlags(prefix string, flags *pflag.FlagSet) {
	flags.Uint32(joinFlag(prefix, "s3-object-info"), cfg.S3ObjectInfo, "Enables caching of S3 object information in the storage driver to reduce S3 calls, value is cache capacity.")
	flags.Uint32(joinFlag(prefix, "home-timeline-size"), cfg.HomeTimelineSize, "Per-user home timeline cache length, in number of posts. (minimum = 100)")
	flags.Uint32(joinFlag(prefix, "list-timeline-size"), cfg.ListTimelineSize, "Per-list timeline cache length, in number of posts. (minimum = 100)")
	flags.Uint32(joinFlag(prefix, "tag-timeline-size"), cfg.TagTimelineSize, "Per-tag timeline cache length, in number of posts. (minimum = 50)")
	flags.Uint32(joinFlag(prefix, "notification-timeline-size"), cfg.NotificationTimelineSize, "Per-user notification timeline cache length, in number of notifications. (minimum = 50)")
	flags.Duration(joinFlag(prefix, "home-timeline-timeout"), cfg.HomeTimelineTimeout, "Duration before any one home timeline cache is unloaded from memory. Values <= 0 disable unloading.")
	flags.Duration(joinFlag(prefix, "list-timeline-timeout"), cfg.ListTimelineTimeout, "Duration before any one list timeline cache is unloaded from memory. Values <= 0 disable unloading.")
	flags.Duration(joinFlag(prefix, "tag-timeline-timeout"), cfg.TagTimelineTimeout, "Duration before any one tag timeline cache is unloaded from memory. Values <= 0 disable unloading.")
	flags.Duration(joinFlag(prefix, "notification-timeline-timeout"), cfg.NotificationTimelineTimeout, "Duration before any one notification timeline cache is unloaded from memory. Values <= disable unloading.")
	flags.String(joinFlag(prefix, "memory-target"), cfg.MemoryTarget.String(), "")
	flags.Float64(joinFlag(prefix, "account-mem-ratio"), cfg.AccountMemRatio, "")
	flags.Float64(joinFlag(prefix, "account-note-mem-ratio"), cfg.AccountNoteMemRatio, "")
	flags.Float64(joinFlag(prefix, "account-settings-mem-ratio"), cfg.AccountSettingsMemRatio, "")
	flags.Float64(joinFlag(prefix, "account-stats-mem-ratio"), cfg.AccountStatsMemRatio, "")
	flags.Float64(joinFlag(prefix, "application-mem-ratio"), cfg.ApplicationMemRatio, "")
	flags.Float64(joinFlag(prefix, "block-mem-ratio"), cfg.BlockMemRatio, "")
	flags.Float64(joinFlag(prefix, "block-ids-mem-ratio"), cfg.BlockIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "boost-of-ids-mem-ratio"), cfg.BoostOfIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "client-mem-ratio"), cfg.ClientMemRatio, "")
	flags.Float64(joinFlag(prefix, "conversation-mem-ratio"), cfg.ConversationMemRatio, "")
	flags.Float64(joinFlag(prefix, "conversation-last-status-ids-mem-ratio"), cfg.ConversationLastStatusIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "domain-permission-draft-mem-ratio"), cfg.DomainPermissionDraftMemRatio, "")
	flags.Float64(joinFlag(prefix, "domain-permission-limit-mem-ratio"), cfg.DomainLimitMemRatio, "")
	flags.Float64(joinFlag(prefix, "domain-permission-subscription-mem-ratio"), cfg.DomainPermissionSubscriptionMemRatio, "")
	flags.Float64(joinFlag(prefix, "emoji-mem-ratio"), cfg.EmojiMemRatio, "")
	flags.Float64(joinFlag(prefix, "emoji-category-mem-ratio"), cfg.EmojiCategoryMemRatio, "")
	flags.Float64(joinFlag(prefix, "federation-error-mem-ratio"), cfg.FederationErrorMemRatio, "")
	flags.Float64(joinFlag(prefix, "filter-mem-ratio"), cfg.FilterMemRatio, "")
	flags.Float64(joinFlag(prefix, "filter-ids-mem-ratio"), cfg.FilterIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "filter-keyword-mem-ratio"), cfg.FilterKeywordMemRatio, "")
	flags.Float64(joinFlag(prefix, "filter-status-mem-ratio"), cfg.FilterStatusMemRatio, "")
	flags.Float64(joinFlag(prefix, "follow-mem-ratio"), cfg.FollowMemRatio, "")
	flags.Float64(joinFlag(prefix, "follow-ids-mem-ratio"), cfg.FollowIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "follow-request-mem-ratio"), cfg.FollowRequestMemRatio, "")
	flags.Float64(joinFlag(prefix, "follow-request-ids-mem-ratio"), cfg.FollowRequestIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "following-tag-ids-mem-ratio"), cfg.FollowingTagIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "home-account-ids-mem-ratio"), cfg.HomeAccountIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "in-reply-to-ids-mem-ratio"), cfg.InReplyToIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "instance-mem-ratio"), cfg.InstanceMemRatio, "")
	flags.Float64(joinFlag(prefix, "interaction-request-mem-ratio"), cfg.InteractionRequestMemRatio, "")
	flags.Float64(joinFlag(prefix, "list-mem-ratio"), cfg.ListMemRatio, "")
	flags.Float64(joinFlag(prefix, "list-ids-mem-ratio"), cfg.ListIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "listed-ids-mem-ratio"), cfg.ListedIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "marker-mem-ratio"), cfg.MarkerMemRatio, "")
	flags.Float64(joinFlag(prefix, "media-mem-ratio"), cfg.MediaMemRatio, "")
	flags.Float64(joinFlag(prefix, "mention-mem-ratio"), cfg.MentionMemRatio, "")
	flags.Float64(joinFlag(prefix, "move-mem-ratio"), cfg.MoveMemRatio, "")
	flags.Float64(joinFlag(prefix, "notification-mem-ratio"), cfg.NotificationMemRatio, "")
	flags.Float64(joinFlag(prefix, "poll-mem-ratio"), cfg.PollMemRatio, "")
	flags.Float64(joinFlag(prefix, "poll-vote-mem-ratio"), cfg.PollVoteMemRatio, "")
	flags.Float64(joinFlag(prefix, "poll-vote-ids-mem-ratio"), cfg.PollVoteIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "report-mem-ratio"), cfg.ReportMemRatio, "")
	flags.Float64(joinFlag(prefix, "relay-actor-mem-ratio"), cfg.RelayActorMemRatio, "")
	flags.Float64(joinFlag(prefix, "relay-matcher-mem-ratio"), cfg.RelayMatcherMemRatio, "")
	flags.Float64(joinFlag(prefix, "relay-push-mem-ratio"), cfg.RelayPushMemRatio, "")
	flags.Float64(joinFlag(prefix, "relay-push-ids-mem-ratio"), cfg.RelayPushIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "relay-subscription-mem-ratio"), cfg.RelaySubscriptionMemRatio, "")
	flags.Float64(joinFlag(prefix, "scheduled-status-mem-ratio"), cfg.ScheduledStatusMemRatio, "")
	flags.Float64(joinFlag(prefix, "sin-bin-status-mem-ratio"), cfg.SinBinStatusMemRatio, "")
	flags.Float64(joinFlag(prefix, "status-mem-ratio"), cfg.StatusMemRatio, "")
	flags.Float64(joinFlag(prefix, "status-bookmark-mem-ratio"), cfg.StatusBookmarkMemRatio, "")
	flags.Float64(joinFlag(prefix, "status-bookmark-ids-mem-ratio"), cfg.StatusBookmarkIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "status-edit-mem-ratio"), cfg.StatusEditMemRatio, "")
	flags.Float64(joinFlag(prefix, "status-fave-mem-ratio"), cfg.StatusFaveMemRatio, "")
	flags.Float64(joinFlag(prefix, "status-fave-ids-mem-ratio"), cfg.StatusFaveIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "status-pinned-ids-mem-ratio"), cfg.StatusPinnedIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "tag-mem-ratio"), cfg.TagMemRatio, "")
	flags.Float64(joinFlag(prefix, "thread-mute-mem-ratio"), cfg.ThreadMuteMemRatio, "")
	flags.Float64(joinFlag(prefix, "token-mem-ratio"), cfg.TokenMemRatio, "")
	flags.Float64(joinFlag(prefix, "tombstone-mem-ratio"), cfg.TombstoneMemRatio, "")
	flags.Float64(joinFlag(prefix, "user-mem-ratio"), cfg.UserMemRatio, "")
	flags.Float64(joinFlag(prefix, "user-mute-mem-ratio"), cfg.UserMuteMemRatio, "")
	flags.Float64(joinFlag(prefix, "user-mute-ids-mem-ratio"), cfg.UserMuteIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "webfinger-mem-ratio"), cfg.WebfingerMemRatio, "")
	flags.Float64(joinFlag(prefix, "web-push-subscription-mem-ratio"), cfg.WebPushSubscriptionMemRatio, "")
	flags.Float64(joinFlag(prefix, "web-push-subscription-ids-mem-ratio"), cfg.WebPushSubscriptionIDsMemRatio, "")
	flags.Float64(joinFlag(prefix, "mutes-mem-ratio"), cfg.MutesMemRatio, "")
	flags.Float64(joinFlag(prefix, "status-filter-mem-ratio"), cfg.StatusFilterMemRatio, "")
	flags.Float64(joinFlag(prefix, "visibility-mem-ratio"), cfg.VisibilityMemRatio, "")
}

func (cfg *Configuration) RegisterFlags(flags *pflag.FlagSet) {
	cfg.Database.RegisterFlags("db", flags)
	cfg.Advanced.RegisterFlags("advanced", flags)
	cfg.HTTPServer.RegisterFlags("http-server", flags)
	cfg.HTTPClient.RegisterFlags("http-client", flags)
	cfg.Media.RegisterFlags("media", flags)
	cfg.Cache.RegisterFlags("cache", flags)
	flags.String("log-level", cfg.LogLevel, "Log level to run at: [trace, debug, info, warn, fatal]")
	flags.String("log-format", cfg.LogFormat, "Log output format: [logfmt, json]")
	flags.String("log-timestamp-format", cfg.LogTimestampFormat, "Format to use for the log timestamp, as supported by Go's time.Layout")
	flags.Bool("log-db-queries", cfg.LogDbQueries, "Log all database queries at TRACE level (regardless of current log-level)")
	flags.Bool("log-client-ip", cfg.LogClientIP, "Include the client IP in logs")
	flags.String("request-id-header", cfg.RequestIDHeader, "Header to extract the Request ID from. Eg.,'X-Request-Id'.")
	flags.String("config-path", cfg.ConfigPath, "Path to a file containing gotosocial configuration. Values set in this file will be overwritten by values set as env vars or arguments")
	flags.String("application-name", cfg.ApplicationName, "Name of the application, used in various places internally")
	flags.String("landing-page-user", cfg.LandingPageUser, "the user that should be shown on the instance's landing page")
	flags.String("host", cfg.Host, "Hostname to use for the server (eg., example.org, gotosocial.whatever.com). DO NOT change this on a server that's already run!")
	flags.String("account-domain", cfg.AccountDomain, "Domain to use in account names (eg., example.org, whatever.com). If not set, will default to the setting for host. DO NOT change this on a server that's already run!")
	flags.String("protocol", cfg.Protocol, "Protocol to use for the REST api of the server (only use http if you are debugging; https should be used even if running behind a reverse proxy!)")
	flags.String("bind-address", cfg.BindAddress, "Bind address to use for the GoToSocial server (eg., 0.0.0.0, 172.138.0.9, [::], localhost). For ipv6, enclose the address in square brackets, eg [2001:db8::fed1]. Default binds to all interfaces.")
	flags.Int("port", cfg.Port, "Port to use for GoToSocial. Change this to 443 if you're running the binary directly on the host machine.")
	flags.StringSlice("trusted-proxies", cfg.TrustedProxies.Strings(), "Proxies to trust when parsing x-forwarded headers into real IPs.")
	flags.String("software-version", cfg.SoftwareVersion, "")
	flags.String("web-template-base-dir", cfg.WebTemplateBaseDir, "Basedir for html templating files for rendering pages and composing emails.")
	flags.String("web-asset-base-dir", cfg.WebAssetBaseDir, "Directory to serve static assets from, accessible at example.org/assets/")
	flags.String("instance-federation-mode", cfg.InstanceFederationMode, "Set instance federation mode.")
	flags.Bool("instance-federation-spam-filter", cfg.InstanceFederationSpamFilter, "Enable basic spam filter heuristics for messages coming from other instances, and drop messages identified as spam")
	flags.Bool("instance-expose-peers", cfg.InstanceExposePeers, "Allow unauthenticated users to query /api/v1/instance/peers?filter=open")
	flags.Bool("instance-expose-blocklist", cfg.InstanceExposeBlocklist, "Expose list of blocked domains via web UI, and allow unauthenticated users to query /api/v1/instance/peers?filter=blocked and /api/v1/instance/domain_blocks")
	flags.Bool("instance-expose-blocklist-web", cfg.InstanceExposeBlocklistWeb, "Expose list of explicitly blocked domains as webpage on /about/domain_blocks")
	flags.Bool("instance-expose-allowlist", cfg.InstanceExposeAllowlist, "Expose list of allowed domains via web UI, and allow unauthenticated users to query /api/v1/instance/peers?filter=allowed and /api/v1/instance/domain_allows")
	flags.Bool("instance-expose-allowlist-web", cfg.InstanceExposeAllowlistWeb, "Expose list of explicitly allowed domains as webpage on /about/domain_allows")
	flags.Bool("instance-expose-public-timeline", cfg.InstanceExposePublicTimeline, "Allow unauthenticated users to query /api/v1/timelines/public")
	flags.Bool("instance-expose-custom-emojis", cfg.InstanceExposeCustomEmojis, "Allow unauthenticated access to /api/v1/custom_emojis")
	flags.String("instance-directory-mode", cfg.InstanceDirectoryMode.String(), "Customize if and how the instance accounts directory is served: one of '' or 'off', 'webonly', or 'open'")
	flags.Bool("instance-deliver-to-shared-inboxes", cfg.InstanceDeliverToSharedInboxes, "Deliver federated messages to shared inboxes, if they're available.")
	flags.Bool("instance-inject-mastodon-version", cfg.InstanceInjectMastodonVersion, "This injects a Mastodon compatible version in /api/v1/instance to help Mastodon clients that use that version for feature detection")
	flags.StringSlice("instance-languages", cfg.InstanceLanguages.Strings(), "BCP47 language tags for the instance. Used to indicate the preferred languages of instance residents (in order from most-preferred to least-preferred).")
	flags.String("instance-stats-mode", cfg.InstanceStatsMode, "Allows you to customize the way stats are served to crawlers: one of '', 'serve', 'zero', 'baffle'. Home page stats remain unchanged.")
	flags.Bool("instance-allow-backdating-statuses", cfg.InstanceAllowBackdatingStatuses, "Allow local accounts to backdate statuses using the scheduled_at param to /api/v1/statuses")
	flags.Bool("instance-robots-allow-indexing", cfg.InstanceRobotsAllowIndexing, "Return robots headers and meta tags that allow search engine indexing of instance home page, directory (if enabled), and accounts that have opted in to being discoverable.")
	flags.Bool("accounts-registration-open", cfg.AccountsRegistrationOpen, "Allow anyone to submit an account signup request. If false, server will be invite-only.")
	flags.Bool("accounts-reason-required", cfg.AccountsReasonRequired, "Do new account signups require a reason to be submitted on registration?")
	flags.Int("accounts-registration-daily-limit", cfg.AccountsRegistrationDailyLimit, "Limit amount of approved account sign-ups allowed per 24hrs before registration is closed. 0 or less = no limit.")
	flags.Int("accounts-registration-backlog-limit", cfg.AccountsRegistrationBacklogLimit, "Limit how big the 'accounts pending approval' queue can grow before registration is closed. 0 or less = no limit.")
	flags.Bool("accounts-allow-custom-css", cfg.AccountsAllowCustomCSS, "Allow accounts to enable custom CSS for their profile pages and statuses.")
	flags.Int("accounts-custom-css-length", cfg.AccountsCustomCSSLength, "Maximum permitted length (characters) of custom CSS for accounts.")
	flags.Int("accounts-max-profile-fields", cfg.AccountsMaxProfileFields, "Maximum number of profile fields allowed for each account.")
	flags.String("storage-backend", cfg.StorageBackend, "Storage backend to use for media attachments")
	flags.String("storage-local-base-path", cfg.StorageLocalBasePath, "Full path to an already-created directory where gts should store/retrieve media files. Subfolders will be created within this dir.")
	flags.String("storage-s3-endpoint", cfg.StorageS3Endpoint, "S3 Endpoint URL (e.g 'minio.example.org:9000')")
	flags.String("storage-s3-access-key", cfg.StorageS3AccessKey, "S3 Access Key")
	flags.String("storage-s3-secret-key", cfg.StorageS3SecretKey, "S3 Secret Key")
	flags.Bool("storage-s3-use-ssl", cfg.StorageS3UseSSL, "Use SSL for S3 connections. Only set this to 'false' when testing locally")
	flags.String("storage-s3-bucket", cfg.StorageS3BucketName, "Place blobs in this bucket")
	flags.Bool("storage-s3-proxy", cfg.StorageS3Proxy, "Proxy S3 contents through GoToSocial instead of redirecting to a presigned URL")
	flags.String("storage-s3-redirect-url", cfg.StorageS3RedirectURL, "Custom URL to use for redirecting S3 media links. If set, this will be used instead of the S3 bucket URL.")
	flags.String("storage-s3-bucket-lookup", cfg.StorageS3BucketLookup, "S3 bucket lookup type to use. Can be 'auto', 'dns' or 'path'. Defaults to 'auto'.")
	flags.String("storage-s3-key-prefix", cfg.StorageS3KeyPrefix, "Prefix to use for S3 keys. This is useful for separating multiple instances sharing the same S3 bucket.")
	flags.String("storage-s3-region", cfg.StorageS3Region, "Region to use for S3.")
	flags.Int("statuses-max-chars", cfg.StatusesMaxChars, "Max permitted characters for posted statuses, including content warning")
	flags.Int("statuses-poll-max-options", cfg.StatusesPollMaxOptions, "Max amount of options permitted on a poll")
	flags.Int("statuses-poll-option-max-chars", cfg.StatusesPollOptionMaxChars, "Max amount of characters for a poll option")
	flags.Int("statuses-media-max-files", cfg.StatusesMediaMaxFiles, "Maximum number of media files/attachments per status")
	flags.Int("scheduled-statuses-max-total", cfg.ScheduledStatusesMaxTotal, "Maximum number of scheduled statuses per user")
	flags.Int("scheduled-statuses-max-daily", cfg.ScheduledStatusesMaxDaily, "Maximum number of scheduled statuses per user for a single day")
	flags.Bool("letsencrypt-enabled", cfg.LetsEncryptEnabled, "Enable letsencrypt TLS certs for this server. If set to true, then cert dir also needs to be set (or take the default).")
	flags.Int("letsencrypt-port", cfg.LetsEncryptPort, "Port to listen on for letsencrypt certificate challenges. Must not be the same as the GtS webserver/API port.")
	flags.String("letsencrypt-cert-dir", cfg.LetsEncryptCertDir, "Directory to store acquired letsencrypt certificates.")
	flags.String("letsencrypt-email-address", cfg.LetsEncryptEmailAddress, "Email address to use when requesting letsencrypt certs. Will receive updates on cert expiry etc.")
	flags.String("tls-certificate-chain", cfg.TLSCertificateChain, "Filesystem path to the certificate chain including any intermediate CAs and the TLS public key")
	flags.String("tls-certificate-key", cfg.TLSCertificateKey, "Filesystem path to the TLS private key")
	flags.Bool("oidc-enabled", cfg.OIDCEnabled, "Enabled OIDC authorization for this instance. If set to true, then the other OIDC flags must also be set.")
	flags.String("oidc-idp-name", cfg.OIDCIdpName, "Name of the OIDC identity provider. Will be shown to the user when logging in.")
	flags.Bool("oidc-skip-verification", cfg.OIDCSkipVerification, "Skip verification of tokens returned by the OIDC provider. Should only be set to 'true' for testing purposes, never in a production environment!")
	flags.String("oidc-issuer", cfg.OIDCIssuer, "Address of the OIDC issuer. Should be the web address, including protocol, at which the issuer can be reached. Eg., 'https://example.org/auth'")
	flags.String("oidc-client-id", cfg.OIDCClientID, "ClientID of GoToSocial, as registered with the OIDC provider.")
	flags.String("oidc-client-secret", cfg.OIDCClientSecret, "ClientSecret of GoToSocial, as registered with the OIDC provider.")
	flags.StringSlice("oidc-scopes", cfg.OIDCScopes, "OIDC scopes.")
	flags.Bool("oidc-link-existing", cfg.OIDCLinkExisting, "link existing user accounts to OIDC logins based on the stored email value")
	flags.StringSlice("oidc-allowed-groups", cfg.OIDCAllowedGroups, "Membership of one of the listed groups allows access to GtS. If this is empty, all groups are allowed.")
	flags.StringSlice("oidc-admin-groups", cfg.OIDCAdminGroups, "Membership of one of the listed groups makes someone a GtS admin")
	flags.Bool("tracing-enabled", cfg.TracingEnabled, "Enable OTLP Tracing")
	flags.Bool("metrics-enabled", cfg.MetricsEnabled, "Enable OpenTelemetry based metrics support.")
	flags.String("smtp-host", cfg.SMTPHost, "Host of the smtp server. Eg., 'smtp.eu.mailgun.org'")
	flags.Int("smtp-port", cfg.SMTPPort, "Port of the smtp server. Eg., 587")
	flags.String("smtp-username", cfg.SMTPUsername, "Username to authenticate with the smtp server as. Eg., 'postmaster@mail.example.org'")
	flags.String("smtp-password", cfg.SMTPPassword, "Password to pass to the smtp server.")
	flags.String("smtp-from", cfg.SMTPFrom, "Address to use as the 'from' field of the email. Eg., 'gotosocial@example.org'")
	flags.String("smtp-from-display-name", cfg.SMTPFromDisplayName, "Optional display name to use in addition to 'from' address. Eg., 'Admin'")
	flags.Bool("smtp-disclose-recipients", cfg.SMTPDiscloseRecipients, "If true, email notifications sent to multiple recipients will be To'd to every recipient at once. If false, recipients will not be disclosed")
	flags.Bool("syslog-enabled", cfg.SyslogEnabled, "Enable the syslog logging hook. Logs will be mirrored to the configured destination.")
	flags.String("syslog-protocol", cfg.SyslogProtocol, "Protocol to use when directing logs to syslog. Leave empty to connect to local syslog.")
	flags.String("syslog-address", cfg.SyslogAddress, "Address:port to send syslog logs to. Leave empty to connect to local syslog.")
	flags.Bool("syslog-mirror", cfg.SyslogMirror, "When syslog is enabled, whether to mirror output syslog. Else, only outputs to syslog.")
	flags.Uint32("syslog-msg-length", cfg.SyslogMsgLength, "Truncates syslog messages beyond this length, defaults to 2048 according to rfc5424.")
}

func (cfg *PostgresConfiguration) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfgmap[joinFlag(prefix, "port")] = cfg.Port
	cfgmap[joinFlag(prefix, "user")] = cfg.User
	cfgmap[joinFlag(prefix, "password")] = cfg.Password
	cfgmap[joinFlag(prefix, "database")] = cfg.Database
	cfgmap[joinFlag(prefix, "tls-mode")] = cfg.TLSMode
	cfgmap[joinFlag(prefix, "tls-ca-cert")] = cfg.TLSCACert
	cfgmap[joinFlag(prefix, "postgres-connection-string")] = cfg.ConnectionString
}

func (cfg *SQLiteConfiguration) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfgmap[joinFlag(prefix, "journal-mode")] = cfg.JournalMode
	cfgmap[joinFlag(prefix, "synchronous")] = cfg.Synchronous
	cfgmap[joinFlag(prefix, "cache-size")] = cfg.CacheSize.String()
	cfgmap[joinFlag(prefix, "busy-timeout")] = cfg.BusyTimeout
}

func (cfg *DatabaseConfiguration) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfg.Postgres.MarshalIntoMap(joinFlag(prefix, ""), cfgmap)
	cfg.SQLite.MarshalIntoMap(joinFlag(prefix, "sqlite"), cfgmap)
	cfgmap[joinFlag(prefix, "type")] = cfg.Type
	cfgmap[joinFlag(prefix, "address")] = cfg.Address
	cfgmap[joinFlag(prefix, "max-open-conns-multiplier")] = cfg.MaxOpenConnsMultiplier
}

func (cfg *RateLimitConfig) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfgmap[joinFlag(prefix, "requests")] = cfg.Requests
	cfgmap[joinFlag(prefix, "exceptions")] = cfg.Exceptions.Strings()
}

func (cfg *ThrottlingConfig) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfgmap[joinFlag(prefix, "multiplier")] = cfg.Multiplier
	cfgmap[joinFlag(prefix, "retry-after")] = cfg.RetryAfter
}

func (cfg *AdvancedConfig) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfg.RateLimit.MarshalIntoMap(joinFlag(prefix, "rate-limit"), cfgmap)
	cfg.Throttling.MarshalIntoMap(joinFlag(prefix, "throttling"), cfgmap)
	cfgmap[joinFlag(prefix, "cookies-samesite")] = cfg.CookiesSamesite
	cfgmap[joinFlag(prefix, "sender-multiplier")] = cfg.SenderMultiplier
	cfgmap[joinFlag(prefix, "csp-extra-uris")] = cfg.CSPExtraURIs
	cfgmap[joinFlag(prefix, "header-filter-mode")] = cfg.HeaderFilterMode
}

func (cfg *HTTPServerConfiguration) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfgmap[joinFlag(prefix, "max-multipart-memory")] = cfg.MaxMultipartMemory.String()
	cfgmap[joinFlag(prefix, "use-h2c")] = cfg.UseH2C
	cfgmap[joinFlag(prefix, "read-timeout")] = cfg.ReadTimeout
	cfgmap[joinFlag(prefix, "read-header-timeout")] = cfg.ReadHeaderTimeout
	cfgmap[joinFlag(prefix, "write-timeout")] = cfg.WriteTimeout
	cfgmap[joinFlag(prefix, "idle-timeout")] = cfg.IdleTimeout
	cfgmap[joinFlag(prefix, "max-header-bytes")] = cfg.MaxHeaderBytes.String()
	cfgmap[joinFlag(prefix, "max-concurrent-streams")] = cfg.MaxConcurrentStreams
	cfgmap[joinFlag(prefix, "max-decoder-header-table-size")] = cfg.MaxDecoderHeaderTableSize.String()
	cfgmap[joinFlag(prefix, "max-encoder-header-table-size")] = cfg.MaxEncoderHeaderTableSize.String()
	cfgmap[joinFlag(prefix, "max-read-frame-size")] = cfg.MaxReadFrameSize.String()
	cfgmap[joinFlag(prefix, "max-receive-buffer-per-connection")] = cfg.MaxReceiveBufferPerConnection.String()
	cfgmap[joinFlag(prefix, "max-receive-buffer-per-stream")] = cfg.MaxReceiveBufferPerStream.String()
	cfgmap[joinFlag(prefix, "send-ping-timeout")] = cfg.SendPingTimeout
	cfgmap[joinFlag(prefix, "ping-timeout")] = cfg.PingTimeout
	cfgmap[joinFlag(prefix, "write-byte-timeout")] = cfg.WriteByteTimeout
}

func (cfg *HTTPClientConfiguration) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfgmap[joinFlag(prefix, "allow-ips")] = cfg.AllowIPs.Strings()
	cfgmap[joinFlag(prefix, "block-ips")] = cfg.BlockIPs.Strings()
	cfgmap[joinFlag(prefix, "timeout")] = cfg.Timeout
	cfgmap[joinFlag(prefix, "tls-insecure-skip-verify")] = cfg.TLSInsecureSkipVerify
	cfgmap[joinFlag(prefix, "insecure-outgoing")] = cfg.InsecureOutgoing
	cfgmap[joinFlag(prefix, "disable-keep-alives")] = cfg.DisableKeepAlives
	cfgmap[joinFlag(prefix, "max-idle-conns")] = cfg.MaxIdleConns
	cfgmap[joinFlag(prefix, "max-idle-conns-per-host")] = cfg.MaxIdleConnsPerHost
	cfgmap[joinFlag(prefix, "max-conns-per-host")] = cfg.MaxConnsPerHost
	cfgmap[joinFlag(prefix, "idle-conn-timeout")] = cfg.IdleConnTimeout
	cfgmap[joinFlag(prefix, "tls-handshake-timeout")] = cfg.TLSHandshakeTimeout
	cfgmap[joinFlag(prefix, "response-header-timeout")] = cfg.ResponseHeaderTimeout
	cfgmap[joinFlag(prefix, "read-buffer-size")] = cfg.ReadBufferSize.String()
	cfgmap[joinFlag(prefix, "write-buffer-size")] = cfg.WriteBufferSize.String()
}

func (cfg *MediaConfiguration) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfgmap[joinFlag(prefix, "description-min-chars")] = cfg.DescriptionMinChars
	cfgmap[joinFlag(prefix, "description-max-chars")] = cfg.DescriptionMaxChars
	cfgmap[joinFlag(prefix, "emoji-local-max-size")] = cfg.EmojiLocalMaxSize.String()
	cfgmap[joinFlag(prefix, "emoji-remote-max-size")] = cfg.EmojiRemoteMaxSize.String()
	cfgmap[joinFlag(prefix, "image-size-hint")] = cfg.ImageSizeHint.String()
	cfgmap[joinFlag(prefix, "video-size-hint")] = cfg.VideoSizeHint.String()
	cfgmap[joinFlag(prefix, "local-max-size")] = cfg.LocalMaxSize.String()
	cfgmap[joinFlag(prefix, "remote-max-size")] = cfg.RemoteMaxSize.String()
	cfgmap[joinFlag(prefix, "ffmpeg-pool-size")] = cfg.FfmpegPoolSize
	cfgmap[joinFlag(prefix, "thumb-max-pixels")] = cfg.ThumbMaxPixels
}

func (cfg *CacheConfiguration) MarshalIntoMap(prefix string, cfgmap map[string]any) {
	cfgmap[joinFlag(prefix, "s3-object-info")] = cfg.S3ObjectInfo
	cfgmap[joinFlag(prefix, "home-timeline-size")] = cfg.HomeTimelineSize
	cfgmap[joinFlag(prefix, "list-timeline-size")] = cfg.ListTimelineSize
	cfgmap[joinFlag(prefix, "tag-timeline-size")] = cfg.TagTimelineSize
	cfgmap[joinFlag(prefix, "notification-timeline-size")] = cfg.NotificationTimelineSize
	cfgmap[joinFlag(prefix, "home-timeline-timeout")] = cfg.HomeTimelineTimeout
	cfgmap[joinFlag(prefix, "list-timeline-timeout")] = cfg.ListTimelineTimeout
	cfgmap[joinFlag(prefix, "tag-timeline-timeout")] = cfg.TagTimelineTimeout
	cfgmap[joinFlag(prefix, "notification-timeline-timeout")] = cfg.NotificationTimelineTimeout
	cfgmap[joinFlag(prefix, "memory-target")] = cfg.MemoryTarget.String()
	cfgmap[joinFlag(prefix, "account-mem-ratio")] = cfg.AccountMemRatio
	cfgmap[joinFlag(prefix, "account-note-mem-ratio")] = cfg.AccountNoteMemRatio
	cfgmap[joinFlag(prefix, "account-settings-mem-ratio")] = cfg.AccountSettingsMemRatio
	cfgmap[joinFlag(prefix, "account-stats-mem-ratio")] = cfg.AccountStatsMemRatio
	cfgmap[joinFlag(prefix, "application-mem-ratio")] = cfg.ApplicationMemRatio
	cfgmap[joinFlag(prefix, "block-mem-ratio")] = cfg.BlockMemRatio
	cfgmap[joinFlag(prefix, "block-ids-mem-ratio")] = cfg.BlockIDsMemRatio
	cfgmap[joinFlag(prefix, "boost-of-ids-mem-ratio")] = cfg.BoostOfIDsMemRatio
	cfgmap[joinFlag(prefix, "client-mem-ratio")] = cfg.ClientMemRatio
	cfgmap[joinFlag(prefix, "conversation-mem-ratio")] = cfg.ConversationMemRatio
	cfgmap[joinFlag(prefix, "conversation-last-status-ids-mem-ratio")] = cfg.ConversationLastStatusIDsMemRatio
	cfgmap[joinFlag(prefix, "domain-permission-draft-mem-ratio")] = cfg.DomainPermissionDraftMemRatio
	cfgmap[joinFlag(prefix, "domain-permission-limit-mem-ratio")] = cfg.DomainLimitMemRatio
	cfgmap[joinFlag(prefix, "domain-permission-subscription-mem-ratio")] = cfg.DomainPermissionSubscriptionMemRatio
	cfgmap[joinFlag(prefix, "emoji-mem-ratio")] = cfg.EmojiMemRatio
	cfgmap[joinFlag(prefix, "emoji-category-mem-ratio")] = cfg.EmojiCategoryMemRatio
	cfgmap[joinFlag(prefix, "federation-error-mem-ratio")] = cfg.FederationErrorMemRatio
	cfgmap[joinFlag(prefix, "filter-mem-ratio")] = cfg.FilterMemRatio
	cfgmap[joinFlag(prefix, "filter-ids-mem-ratio")] = cfg.FilterIDsMemRatio
	cfgmap[joinFlag(prefix, "filter-keyword-mem-ratio")] = cfg.FilterKeywordMemRatio
	cfgmap[joinFlag(prefix, "filter-status-mem-ratio")] = cfg.FilterStatusMemRatio
	cfgmap[joinFlag(prefix, "follow-mem-ratio")] = cfg.FollowMemRatio
	cfgmap[joinFlag(prefix, "follow-ids-mem-ratio")] = cfg.FollowIDsMemRatio
	cfgmap[joinFlag(prefix, "follow-request-mem-ratio")] = cfg.FollowRequestMemRatio
	cfgmap[joinFlag(prefix, "follow-request-ids-mem-ratio")] = cfg.FollowRequestIDsMemRatio
	cfgmap[joinFlag(prefix, "following-tag-ids-mem-ratio")] = cfg.FollowingTagIDsMemRatio
	cfgmap[joinFlag(prefix, "home-account-ids-mem-ratio")] = cfg.HomeAccountIDsMemRatio
	cfgmap[joinFlag(prefix, "in-reply-to-ids-mem-ratio")] = cfg.InReplyToIDsMemRatio
	cfgmap[joinFlag(prefix, "instance-mem-ratio")] = cfg.InstanceMemRatio
	cfgmap[joinFlag(prefix, "interaction-request-mem-ratio")] = cfg.InteractionRequestMemRatio
	cfgmap[joinFlag(prefix, "list-mem-ratio")] = cfg.ListMemRatio
	cfgmap[joinFlag(prefix, "list-ids-mem-ratio")] = cfg.ListIDsMemRatio
	cfgmap[joinFlag(prefix, "listed-ids-mem-ratio")] = cfg.ListedIDsMemRatio
	cfgmap[joinFlag(prefix, "marker-mem-ratio")] = cfg.MarkerMemRatio
	cfgmap[joinFlag(prefix, "media-mem-ratio")] = cfg.MediaMemRatio
	cfgmap[joinFlag(prefix, "mention-mem-ratio")] = cfg.MentionMemRatio
	cfgmap[joinFlag(prefix, "move-mem-ratio")] = cfg.MoveMemRatio
	cfgmap[joinFlag(prefix, "notification-mem-ratio")] = cfg.NotificationMemRatio
	cfgmap[joinFlag(prefix, "poll-mem-ratio")] = cfg.PollMemRatio
	cfgmap[joinFlag(prefix, "poll-vote-mem-ratio")] = cfg.PollVoteMemRatio
	cfgmap[joinFlag(prefix, "poll-vote-ids-mem-ratio")] = cfg.PollVoteIDsMemRatio
	cfgmap[joinFlag(prefix, "report-mem-ratio")] = cfg.ReportMemRatio
	cfgmap[joinFlag(prefix, "relay-actor-mem-ratio")] = cfg.RelayActorMemRatio
	cfgmap[joinFlag(prefix, "relay-matcher-mem-ratio")] = cfg.RelayMatcherMemRatio
	cfgmap[joinFlag(prefix, "relay-push-mem-ratio")] = cfg.RelayPushMemRatio
	cfgmap[joinFlag(prefix, "relay-push-ids-mem-ratio")] = cfg.RelayPushIDsMemRatio
	cfgmap[joinFlag(prefix, "relay-subscription-mem-ratio")] = cfg.RelaySubscriptionMemRatio
	cfgmap[joinFlag(prefix, "scheduled-status-mem-ratio")] = cfg.ScheduledStatusMemRatio
	cfgmap[joinFlag(prefix, "sin-bin-status-mem-ratio")] = cfg.SinBinStatusMemRatio
	cfgmap[joinFlag(prefix, "status-mem-ratio")] = cfg.StatusMemRatio
	cfgmap[joinFlag(prefix, "status-bookmark-mem-ratio")] = cfg.StatusBookmarkMemRatio
	cfgmap[joinFlag(prefix, "status-bookmark-ids-mem-ratio")] = cfg.StatusBookmarkIDsMemRatio
	cfgmap[joinFlag(prefix, "status-edit-mem-ratio")] = cfg.StatusEditMemRatio
	cfgmap[joinFlag(prefix, "status-fave-mem-ratio")] = cfg.StatusFaveMemRatio
	cfgmap[joinFlag(prefix, "status-fave-ids-mem-ratio")] = cfg.StatusFaveIDsMemRatio
	cfgmap[joinFlag(prefix, "status-pinned-ids-mem-ratio")] = cfg.StatusPinnedIDsMemRatio
	cfgmap[joinFlag(prefix, "tag-mem-ratio")] = cfg.TagMemRatio
	cfgmap[joinFlag(prefix, "thread-mute-mem-ratio")] = cfg.ThreadMuteMemRatio
	cfgmap[joinFlag(prefix, "token-mem-ratio")] = cfg.TokenMemRatio
	cfgmap[joinFlag(prefix, "tombstone-mem-ratio")] = cfg.TombstoneMemRatio
	cfgmap[joinFlag(prefix, "user-mem-ratio")] = cfg.UserMemRatio
	cfgmap[joinFlag(prefix, "user-mute-mem-ratio")] = cfg.UserMuteMemRatio
	cfgmap[joinFlag(prefix, "user-mute-ids-mem-ratio")] = cfg.UserMuteIDsMemRatio
	cfgmap[joinFlag(prefix, "webfinger-mem-ratio")] = cfg.WebfingerMemRatio
	cfgmap[joinFlag(prefix, "web-push-subscription-mem-ratio")] = cfg.WebPushSubscriptionMemRatio
	cfgmap[joinFlag(prefix, "web-push-subscription-ids-mem-ratio")] = cfg.WebPushSubscriptionIDsMemRatio
	cfgmap[joinFlag(prefix, "mutes-mem-ratio")] = cfg.MutesMemRatio
	cfgmap[joinFlag(prefix, "status-filter-mem-ratio")] = cfg.StatusFilterMemRatio
	cfgmap[joinFlag(prefix, "visibility-mem-ratio")] = cfg.VisibilityMemRatio
}

func (cfg *Configuration) MarshalMap() map[string]any {
	cfgmap := make(map[string]any, 239)
	cfg.MarshalIntoMap(cfgmap)
	return cfgmap
}

func (cfg *Configuration) MarshalIntoMap(cfgmap map[string]any) {
	cfg.Database.MarshalIntoMap("db", cfgmap)
	cfg.Advanced.MarshalIntoMap("advanced", cfgmap)
	cfg.HTTPServer.MarshalIntoMap("http-server", cfgmap)
	cfg.HTTPClient.MarshalIntoMap("http-client", cfgmap)
	cfg.Media.MarshalIntoMap("media", cfgmap)
	cfg.Cache.MarshalIntoMap("cache", cfgmap)
	cfgmap["log-level"] = cfg.LogLevel
	cfgmap["log-format"] = cfg.LogFormat
	cfgmap["log-timestamp-format"] = cfg.LogTimestampFormat
	cfgmap["log-db-queries"] = cfg.LogDbQueries
	cfgmap["log-client-ip"] = cfg.LogClientIP
	cfgmap["request-id-header"] = cfg.RequestIDHeader
	cfgmap["config-path"] = cfg.ConfigPath
	cfgmap["application-name"] = cfg.ApplicationName
	cfgmap["landing-page-user"] = cfg.LandingPageUser
	cfgmap["host"] = cfg.Host
	cfgmap["account-domain"] = cfg.AccountDomain
	cfgmap["protocol"] = cfg.Protocol
	cfgmap["bind-address"] = cfg.BindAddress
	cfgmap["port"] = cfg.Port
	cfgmap["trusted-proxies"] = cfg.TrustedProxies.Strings()
	cfgmap["software-version"] = cfg.SoftwareVersion
	cfgmap["web-template-base-dir"] = cfg.WebTemplateBaseDir
	cfgmap["web-asset-base-dir"] = cfg.WebAssetBaseDir
	cfgmap["instance-federation-mode"] = cfg.InstanceFederationMode
	cfgmap["instance-federation-spam-filter"] = cfg.InstanceFederationSpamFilter
	cfgmap["instance-expose-peers"] = cfg.InstanceExposePeers
	cfgmap["instance-expose-blocklist"] = cfg.InstanceExposeBlocklist
	cfgmap["instance-expose-blocklist-web"] = cfg.InstanceExposeBlocklistWeb
	cfgmap["instance-expose-allowlist"] = cfg.InstanceExposeAllowlist
	cfgmap["instance-expose-allowlist-web"] = cfg.InstanceExposeAllowlistWeb
	cfgmap["instance-expose-public-timeline"] = cfg.InstanceExposePublicTimeline
	cfgmap["instance-expose-custom-emojis"] = cfg.InstanceExposeCustomEmojis
	cfgmap["instance-directory-mode"] = cfg.InstanceDirectoryMode.String()
	cfgmap["instance-deliver-to-shared-inboxes"] = cfg.InstanceDeliverToSharedInboxes
	cfgmap["instance-inject-mastodon-version"] = cfg.InstanceInjectMastodonVersion
	cfgmap["instance-languages"] = cfg.InstanceLanguages.Strings()
	cfgmap["instance-stats-mode"] = cfg.InstanceStatsMode
	cfgmap["instance-allow-backdating-statuses"] = cfg.InstanceAllowBackdatingStatuses
	cfgmap["instance-robots-allow-indexing"] = cfg.InstanceRobotsAllowIndexing
	cfgmap["accounts-registration-open"] = cfg.AccountsRegistrationOpen
	cfgmap["accounts-reason-required"] = cfg.AccountsReasonRequired
	cfgmap["accounts-registration-daily-limit"] = cfg.AccountsRegistrationDailyLimit
	cfgmap["accounts-registration-backlog-limit"] = cfg.AccountsRegistrationBacklogLimit
	cfgmap["accounts-allow-custom-css"] = cfg.AccountsAllowCustomCSS
	cfgmap["accounts-custom-css-length"] = cfg.AccountsCustomCSSLength
	cfgmap["accounts-max-profile-fields"] = cfg.AccountsMaxProfileFields
	cfgmap["storage-backend"] = cfg.StorageBackend
	cfgmap["storage-local-base-path"] = cfg.StorageLocalBasePath
	cfgmap["storage-s3-endpoint"] = cfg.StorageS3Endpoint
	cfgmap["storage-s3-access-key"] = cfg.StorageS3AccessKey
	cfgmap["storage-s3-secret-key"] = cfg.StorageS3SecretKey
	cfgmap["storage-s3-use-ssl"] = cfg.StorageS3UseSSL
	cfgmap["storage-s3-bucket"] = cfg.StorageS3BucketName
	cfgmap["storage-s3-proxy"] = cfg.StorageS3Proxy
	cfgmap["storage-s3-redirect-url"] = cfg.StorageS3RedirectURL
	cfgmap["storage-s3-bucket-lookup"] = cfg.StorageS3BucketLookup
	cfgmap["storage-s3-key-prefix"] = cfg.StorageS3KeyPrefix
	cfgmap["storage-s3-region"] = cfg.StorageS3Region
	cfgmap["statuses-max-chars"] = cfg.StatusesMaxChars
	cfgmap["statuses-poll-max-options"] = cfg.StatusesPollMaxOptions
	cfgmap["statuses-poll-option-max-chars"] = cfg.StatusesPollOptionMaxChars
	cfgmap["statuses-media-max-files"] = cfg.StatusesMediaMaxFiles
	cfgmap["scheduled-statuses-max-total"] = cfg.ScheduledStatusesMaxTotal
	cfgmap["scheduled-statuses-max-daily"] = cfg.ScheduledStatusesMaxDaily
	cfgmap["letsencrypt-enabled"] = cfg.LetsEncryptEnabled
	cfgmap["letsencrypt-port"] = cfg.LetsEncryptPort
	cfgmap["letsencrypt-cert-dir"] = cfg.LetsEncryptCertDir
	cfgmap["letsencrypt-email-address"] = cfg.LetsEncryptEmailAddress
	cfgmap["tls-certificate-chain"] = cfg.TLSCertificateChain
	cfgmap["tls-certificate-key"] = cfg.TLSCertificateKey
	cfgmap["oidc-enabled"] = cfg.OIDCEnabled
	cfgmap["oidc-idp-name"] = cfg.OIDCIdpName
	cfgmap["oidc-skip-verification"] = cfg.OIDCSkipVerification
	cfgmap["oidc-issuer"] = cfg.OIDCIssuer
	cfgmap["oidc-client-id"] = cfg.OIDCClientID
	cfgmap["oidc-client-secret"] = cfg.OIDCClientSecret
	cfgmap["oidc-scopes"] = cfg.OIDCScopes
	cfgmap["oidc-link-existing"] = cfg.OIDCLinkExisting
	cfgmap["oidc-allowed-groups"] = cfg.OIDCAllowedGroups
	cfgmap["oidc-admin-groups"] = cfg.OIDCAdminGroups
	cfgmap["tracing-enabled"] = cfg.TracingEnabled
	cfgmap["metrics-enabled"] = cfg.MetricsEnabled
	cfgmap["smtp-host"] = cfg.SMTPHost
	cfgmap["smtp-port"] = cfg.SMTPPort
	cfgmap["smtp-username"] = cfg.SMTPUsername
	cfgmap["smtp-password"] = cfg.SMTPPassword
	cfgmap["smtp-from"] = cfg.SMTPFrom
	cfgmap["smtp-from-display-name"] = cfg.SMTPFromDisplayName
	cfgmap["smtp-disclose-recipients"] = cfg.SMTPDiscloseRecipients
	cfgmap["syslog-enabled"] = cfg.SyslogEnabled
	cfgmap["syslog-protocol"] = cfg.SyslogProtocol
	cfgmap["syslog-address"] = cfg.SyslogAddress
	cfgmap["syslog-mirror"] = cfg.SyslogMirror
	cfgmap["syslog-msg-length"] = cfg.SyslogMsgLength
	cfgmap["username"] = cfg.AdminAccountUsername
	cfgmap["email"] = cfg.AdminAccountEmail
	cfgmap["password"] = cfg.AdminAccountPassword
	cfgmap["path"] = cfg.AdminTransPath
	cfgmap["dry-run"] = cfg.AdminMediaPruneDryRun
	cfgmap["local-only"] = cfg.AdminMediaListLocalOnly
	cfgmap["remote-only"] = cfg.AdminMediaListRemoteOnly
	cfgmap["skip-db-setup"] = cfg.TestrigSkipDBSetup
	cfgmap["skip-db-teardown"] = cfg.TestrigSkipDBTeardown
}

func (cfg *PostgresConfiguration) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if ival, ok := cfgmap[joinFlag(prefix, "port")]; ok {
		var err error
		cfg.Port, err = cast.ToUint16E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "uint16", joinFlag(prefix, "port"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "user")]; ok {
		var err error
		cfg.User, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "user"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "password")]; ok {
		var err error
		cfg.Password, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "password"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "database")]; ok {
		var err error
		cfg.Database, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "database"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "tls-mode")]; ok {
		var err error
		cfg.TLSMode, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "tls-mode"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "tls-ca-cert")]; ok {
		var err error
		cfg.TLSCACert, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "tls-ca-cert"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "postgres-connection-string")]; ok {
		var err error
		cfg.ConnectionString, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "postgres-connection-string"), err)
		}
	}

	return nil
}

func (cfg *SQLiteConfiguration) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if ival, ok := cfgmap[joinFlag(prefix, "journal-mode")]; ok {
		var err error
		cfg.JournalMode, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "journal-mode"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "synchronous")]; ok {
		var err error
		cfg.Synchronous, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "synchronous"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "cache-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "cache-size"), err)
		}
		cfg.CacheSize = 0x0
		if err := cfg.CacheSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "cache-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "busy-timeout")]; ok {
		var err error
		cfg.BusyTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "busy-timeout"), err)
		}
	}

	return nil
}

func (cfg *DatabaseConfiguration) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if err := cfg.Postgres.UnmarshalMap(joinFlag(prefix, ""), cfgmap); err != nil {
		return err
	}

	if err := cfg.SQLite.UnmarshalMap(joinFlag(prefix, "sqlite"), cfgmap); err != nil {
		return err
	}

	if ival, ok := cfgmap[joinFlag(prefix, "type")]; ok {
		var err error
		cfg.Type, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "type"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "address")]; ok {
		var err error
		cfg.Address, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "address"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-open-conns-multiplier")]; ok {
		var err error
		cfg.MaxOpenConnsMultiplier, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "max-open-conns-multiplier"), err)
		}
	}

	return nil
}

func (cfg *RateLimitConfig) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if ival, ok := cfgmap[joinFlag(prefix, "requests")]; ok {
		var err error
		cfg.Requests, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "requests"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "exceptions")]; ok {
		t, err := toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> []string for '%s': %w", ival, joinFlag(prefix, "exceptions"), err)
		}
		cfg.Exceptions = IPPrefixes{}
		for _, in := range t {
			if err := cfg.Exceptions.Set(in); err != nil {
				return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "exceptions"), err)
			}
		}
	}

	return nil
}

func (cfg *ThrottlingConfig) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if ival, ok := cfgmap[joinFlag(prefix, "multiplier")]; ok {
		var err error
		cfg.Multiplier, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "multiplier"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "retry-after")]; ok {
		var err error
		cfg.RetryAfter, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "retry-after"), err)
		}
	}

	return nil
}

func (cfg *AdvancedConfig) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if err := cfg.RateLimit.UnmarshalMap(joinFlag(prefix, "rate-limit"), cfgmap); err != nil {
		return err
	}

	if err := cfg.Throttling.UnmarshalMap(joinFlag(prefix, "throttling"), cfgmap); err != nil {
		return err
	}

	if ival, ok := cfgmap[joinFlag(prefix, "cookies-samesite")]; ok {
		var err error
		cfg.CookiesSamesite, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "cookies-samesite"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "sender-multiplier")]; ok {
		var err error
		cfg.SenderMultiplier, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "sender-multiplier"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "csp-extra-uris")]; ok {
		var err error
		cfg.CSPExtraURIs, err = toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "[]string", joinFlag(prefix, "csp-extra-uris"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "header-filter-mode")]; ok {
		var err error
		cfg.HeaderFilterMode, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", joinFlag(prefix, "header-filter-mode"), err)
		}
	}

	return nil
}

func (cfg *HTTPServerConfiguration) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if ival, ok := cfgmap[joinFlag(prefix, "max-multipart-memory")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "max-multipart-memory"), err)
		}
		cfg.MaxMultipartMemory = 0x0
		if err := cfg.MaxMultipartMemory.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "max-multipart-memory"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "use-h2c")]; ok {
		var err error
		cfg.UseH2C, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", joinFlag(prefix, "use-h2c"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "read-timeout")]; ok {
		var err error
		cfg.ReadTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "read-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "read-header-timeout")]; ok {
		var err error
		cfg.ReadHeaderTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "read-header-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "write-timeout")]; ok {
		var err error
		cfg.WriteTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "write-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "idle-timeout")]; ok {
		var err error
		cfg.IdleTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "idle-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-header-bytes")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "max-header-bytes"), err)
		}
		cfg.MaxHeaderBytes = 0x0
		if err := cfg.MaxHeaderBytes.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "max-header-bytes"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-concurrent-streams")]; ok {
		var err error
		cfg.MaxConcurrentStreams, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "max-concurrent-streams"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-decoder-header-table-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "max-decoder-header-table-size"), err)
		}
		cfg.MaxDecoderHeaderTableSize = 0x0
		if err := cfg.MaxDecoderHeaderTableSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "max-decoder-header-table-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-encoder-header-table-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "max-encoder-header-table-size"), err)
		}
		cfg.MaxEncoderHeaderTableSize = 0x0
		if err := cfg.MaxEncoderHeaderTableSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "max-encoder-header-table-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-read-frame-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "max-read-frame-size"), err)
		}
		cfg.MaxReadFrameSize = 0x0
		if err := cfg.MaxReadFrameSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "max-read-frame-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-receive-buffer-per-connection")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "max-receive-buffer-per-connection"), err)
		}
		cfg.MaxReceiveBufferPerConnection = 0x0
		if err := cfg.MaxReceiveBufferPerConnection.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "max-receive-buffer-per-connection"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-receive-buffer-per-stream")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "max-receive-buffer-per-stream"), err)
		}
		cfg.MaxReceiveBufferPerStream = 0x0
		if err := cfg.MaxReceiveBufferPerStream.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "max-receive-buffer-per-stream"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "send-ping-timeout")]; ok {
		var err error
		cfg.SendPingTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "send-ping-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "ping-timeout")]; ok {
		var err error
		cfg.PingTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "ping-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "write-byte-timeout")]; ok {
		var err error
		cfg.WriteByteTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "write-byte-timeout"), err)
		}
	}

	return nil
}

func (cfg *HTTPClientConfiguration) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if ival, ok := cfgmap[joinFlag(prefix, "allow-ips")]; ok {
		t, err := toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> []string for '%s': %w", ival, joinFlag(prefix, "allow-ips"), err)
		}
		cfg.AllowIPs = IPPrefixes{}
		for _, in := range t {
			if err := cfg.AllowIPs.Set(in); err != nil {
				return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "allow-ips"), err)
			}
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "block-ips")]; ok {
		t, err := toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> []string for '%s': %w", ival, joinFlag(prefix, "block-ips"), err)
		}
		cfg.BlockIPs = IPPrefixes{}
		for _, in := range t {
			if err := cfg.BlockIPs.Set(in); err != nil {
				return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "block-ips"), err)
			}
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "timeout")]; ok {
		var err error
		cfg.Timeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "tls-insecure-skip-verify")]; ok {
		var err error
		cfg.TLSInsecureSkipVerify, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", joinFlag(prefix, "tls-insecure-skip-verify"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "insecure-outgoing")]; ok {
		var err error
		cfg.InsecureOutgoing, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", joinFlag(prefix, "insecure-outgoing"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "disable-keep-alives")]; ok {
		var err error
		cfg.DisableKeepAlives, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", joinFlag(prefix, "disable-keep-alives"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-idle-conns")]; ok {
		var err error
		cfg.MaxIdleConns, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "max-idle-conns"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-idle-conns-per-host")]; ok {
		var err error
		cfg.MaxIdleConnsPerHost, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "max-idle-conns-per-host"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "max-conns-per-host")]; ok {
		var err error
		cfg.MaxConnsPerHost, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "max-conns-per-host"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "idle-conn-timeout")]; ok {
		var err error
		cfg.IdleConnTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "idle-conn-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "tls-handshake-timeout")]; ok {
		var err error
		cfg.TLSHandshakeTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "tls-handshake-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "response-header-timeout")]; ok {
		var err error
		cfg.ResponseHeaderTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "response-header-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "read-buffer-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "read-buffer-size"), err)
		}
		cfg.ReadBufferSize = 0x0
		if err := cfg.ReadBufferSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "read-buffer-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "write-buffer-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "write-buffer-size"), err)
		}
		cfg.WriteBufferSize = 0x0
		if err := cfg.WriteBufferSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "write-buffer-size"), err)
		}
	}

	return nil
}

func (cfg *MediaConfiguration) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if ival, ok := cfgmap[joinFlag(prefix, "description-min-chars")]; ok {
		var err error
		cfg.DescriptionMinChars, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "description-min-chars"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "description-max-chars")]; ok {
		var err error
		cfg.DescriptionMaxChars, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "description-max-chars"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "emoji-local-max-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "emoji-local-max-size"), err)
		}
		cfg.EmojiLocalMaxSize = 0x0
		if err := cfg.EmojiLocalMaxSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "emoji-local-max-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "emoji-remote-max-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "emoji-remote-max-size"), err)
		}
		cfg.EmojiRemoteMaxSize = 0x0
		if err := cfg.EmojiRemoteMaxSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "emoji-remote-max-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "image-size-hint")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "image-size-hint"), err)
		}
		cfg.ImageSizeHint = 0x0
		if err := cfg.ImageSizeHint.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "image-size-hint"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "video-size-hint")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "video-size-hint"), err)
		}
		cfg.VideoSizeHint = 0x0
		if err := cfg.VideoSizeHint.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "video-size-hint"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "local-max-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "local-max-size"), err)
		}
		cfg.LocalMaxSize = 0x0
		if err := cfg.LocalMaxSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "local-max-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "remote-max-size")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "remote-max-size"), err)
		}
		cfg.RemoteMaxSize = 0x0
		if err := cfg.RemoteMaxSize.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "remote-max-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "ffmpeg-pool-size")]; ok {
		var err error
		cfg.FfmpegPoolSize, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "ffmpeg-pool-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "thumb-max-pixels")]; ok {
		var err error
		cfg.ThumbMaxPixels, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", joinFlag(prefix, "thumb-max-pixels"), err)
		}
	}

	return nil
}

func (cfg *CacheConfiguration) UnmarshalMap(prefix string, cfgmap map[string]any) error {
	if ival, ok := cfgmap[joinFlag(prefix, "s3-object-info")]; ok {
		var err error
		cfg.S3ObjectInfo, err = cast.ToUint32E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "uint32", joinFlag(prefix, "s3-object-info"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "home-timeline-size")]; ok {
		var err error
		cfg.HomeTimelineSize, err = cast.ToUint32E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "uint32", joinFlag(prefix, "home-timeline-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "list-timeline-size")]; ok {
		var err error
		cfg.ListTimelineSize, err = cast.ToUint32E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "uint32", joinFlag(prefix, "list-timeline-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "tag-timeline-size")]; ok {
		var err error
		cfg.TagTimelineSize, err = cast.ToUint32E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "uint32", joinFlag(prefix, "tag-timeline-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "notification-timeline-size")]; ok {
		var err error
		cfg.NotificationTimelineSize, err = cast.ToUint32E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "uint32", joinFlag(prefix, "notification-timeline-size"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "home-timeline-timeout")]; ok {
		var err error
		cfg.HomeTimelineTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "home-timeline-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "list-timeline-timeout")]; ok {
		var err error
		cfg.ListTimelineTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "list-timeline-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "tag-timeline-timeout")]; ok {
		var err error
		cfg.TagTimelineTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "tag-timeline-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "notification-timeline-timeout")]; ok {
		var err error
		cfg.NotificationTimelineTimeout, err = cast.ToDurationE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "time.Duration", joinFlag(prefix, "notification-timeline-timeout"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "memory-target")]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, joinFlag(prefix, "memory-target"), err)
		}
		cfg.MemoryTarget = 0x0
		if err := cfg.MemoryTarget.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, joinFlag(prefix, "memory-target"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "account-mem-ratio")]; ok {
		var err error
		cfg.AccountMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "account-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "account-note-mem-ratio")]; ok {
		var err error
		cfg.AccountNoteMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "account-note-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "account-settings-mem-ratio")]; ok {
		var err error
		cfg.AccountSettingsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "account-settings-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "account-stats-mem-ratio")]; ok {
		var err error
		cfg.AccountStatsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "account-stats-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "application-mem-ratio")]; ok {
		var err error
		cfg.ApplicationMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "application-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "block-mem-ratio")]; ok {
		var err error
		cfg.BlockMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "block-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "block-ids-mem-ratio")]; ok {
		var err error
		cfg.BlockIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "block-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "boost-of-ids-mem-ratio")]; ok {
		var err error
		cfg.BoostOfIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "boost-of-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "client-mem-ratio")]; ok {
		var err error
		cfg.ClientMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "client-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "conversation-mem-ratio")]; ok {
		var err error
		cfg.ConversationMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "conversation-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "conversation-last-status-ids-mem-ratio")]; ok {
		var err error
		cfg.ConversationLastStatusIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "conversation-last-status-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "domain-permission-draft-mem-ratio")]; ok {
		var err error
		cfg.DomainPermissionDraftMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "domain-permission-draft-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "domain-permission-limit-mem-ratio")]; ok {
		var err error
		cfg.DomainLimitMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "domain-permission-limit-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "domain-permission-subscription-mem-ratio")]; ok {
		var err error
		cfg.DomainPermissionSubscriptionMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "domain-permission-subscription-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "emoji-mem-ratio")]; ok {
		var err error
		cfg.EmojiMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "emoji-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "emoji-category-mem-ratio")]; ok {
		var err error
		cfg.EmojiCategoryMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "emoji-category-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "federation-error-mem-ratio")]; ok {
		var err error
		cfg.FederationErrorMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "federation-error-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "filter-mem-ratio")]; ok {
		var err error
		cfg.FilterMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "filter-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "filter-ids-mem-ratio")]; ok {
		var err error
		cfg.FilterIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "filter-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "filter-keyword-mem-ratio")]; ok {
		var err error
		cfg.FilterKeywordMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "filter-keyword-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "filter-status-mem-ratio")]; ok {
		var err error
		cfg.FilterStatusMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "filter-status-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "follow-mem-ratio")]; ok {
		var err error
		cfg.FollowMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "follow-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "follow-ids-mem-ratio")]; ok {
		var err error
		cfg.FollowIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "follow-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "follow-request-mem-ratio")]; ok {
		var err error
		cfg.FollowRequestMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "follow-request-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "follow-request-ids-mem-ratio")]; ok {
		var err error
		cfg.FollowRequestIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "follow-request-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "following-tag-ids-mem-ratio")]; ok {
		var err error
		cfg.FollowingTagIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "following-tag-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "home-account-ids-mem-ratio")]; ok {
		var err error
		cfg.HomeAccountIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "home-account-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "in-reply-to-ids-mem-ratio")]; ok {
		var err error
		cfg.InReplyToIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "in-reply-to-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "instance-mem-ratio")]; ok {
		var err error
		cfg.InstanceMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "instance-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "interaction-request-mem-ratio")]; ok {
		var err error
		cfg.InteractionRequestMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "interaction-request-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "list-mem-ratio")]; ok {
		var err error
		cfg.ListMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "list-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "list-ids-mem-ratio")]; ok {
		var err error
		cfg.ListIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "list-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "listed-ids-mem-ratio")]; ok {
		var err error
		cfg.ListedIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "listed-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "marker-mem-ratio")]; ok {
		var err error
		cfg.MarkerMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "marker-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "media-mem-ratio")]; ok {
		var err error
		cfg.MediaMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "media-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "mention-mem-ratio")]; ok {
		var err error
		cfg.MentionMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "mention-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "move-mem-ratio")]; ok {
		var err error
		cfg.MoveMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "move-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "notification-mem-ratio")]; ok {
		var err error
		cfg.NotificationMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "notification-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "poll-mem-ratio")]; ok {
		var err error
		cfg.PollMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "poll-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "poll-vote-mem-ratio")]; ok {
		var err error
		cfg.PollVoteMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "poll-vote-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "poll-vote-ids-mem-ratio")]; ok {
		var err error
		cfg.PollVoteIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "poll-vote-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "report-mem-ratio")]; ok {
		var err error
		cfg.ReportMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "report-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "relay-actor-mem-ratio")]; ok {
		var err error
		cfg.RelayActorMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "relay-actor-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "relay-matcher-mem-ratio")]; ok {
		var err error
		cfg.RelayMatcherMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "relay-matcher-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "relay-push-mem-ratio")]; ok {
		var err error
		cfg.RelayPushMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "relay-push-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "relay-push-ids-mem-ratio")]; ok {
		var err error
		cfg.RelayPushIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "relay-push-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "relay-subscription-mem-ratio")]; ok {
		var err error
		cfg.RelaySubscriptionMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "relay-subscription-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "scheduled-status-mem-ratio")]; ok {
		var err error
		cfg.ScheduledStatusMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "scheduled-status-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "sin-bin-status-mem-ratio")]; ok {
		var err error
		cfg.SinBinStatusMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "sin-bin-status-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "status-mem-ratio")]; ok {
		var err error
		cfg.StatusMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "status-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "status-bookmark-mem-ratio")]; ok {
		var err error
		cfg.StatusBookmarkMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "status-bookmark-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "status-bookmark-ids-mem-ratio")]; ok {
		var err error
		cfg.StatusBookmarkIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "status-bookmark-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "status-edit-mem-ratio")]; ok {
		var err error
		cfg.StatusEditMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "status-edit-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "status-fave-mem-ratio")]; ok {
		var err error
		cfg.StatusFaveMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "status-fave-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "status-fave-ids-mem-ratio")]; ok {
		var err error
		cfg.StatusFaveIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "status-fave-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "status-pinned-ids-mem-ratio")]; ok {
		var err error
		cfg.StatusPinnedIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "status-pinned-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "tag-mem-ratio")]; ok {
		var err error
		cfg.TagMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "tag-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "thread-mute-mem-ratio")]; ok {
		var err error
		cfg.ThreadMuteMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "thread-mute-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "token-mem-ratio")]; ok {
		var err error
		cfg.TokenMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "token-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "tombstone-mem-ratio")]; ok {
		var err error
		cfg.TombstoneMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "tombstone-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "user-mem-ratio")]; ok {
		var err error
		cfg.UserMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "user-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "user-mute-mem-ratio")]; ok {
		var err error
		cfg.UserMuteMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "user-mute-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "user-mute-ids-mem-ratio")]; ok {
		var err error
		cfg.UserMuteIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "user-mute-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "webfinger-mem-ratio")]; ok {
		var err error
		cfg.WebfingerMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "webfinger-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "web-push-subscription-mem-ratio")]; ok {
		var err error
		cfg.WebPushSubscriptionMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "web-push-subscription-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "web-push-subscription-ids-mem-ratio")]; ok {
		var err error
		cfg.WebPushSubscriptionIDsMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "web-push-subscription-ids-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "mutes-mem-ratio")]; ok {
		var err error
		cfg.MutesMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "mutes-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "status-filter-mem-ratio")]; ok {
		var err error
		cfg.StatusFilterMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "status-filter-mem-ratio"), err)
		}
	}

	if ival, ok := cfgmap[joinFlag(prefix, "visibility-mem-ratio")]; ok {
		var err error
		cfg.VisibilityMemRatio, err = cast.ToFloat64E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "float64", joinFlag(prefix, "visibility-mem-ratio"), err)
		}
	}

	return nil
}

func (cfg *Configuration) UnmarshalMap(cfgmap map[string]any) error {
	if err := cfg.Database.UnmarshalMap("db", cfgmap); err != nil {
		return err
	}

	if err := cfg.Advanced.UnmarshalMap("advanced", cfgmap); err != nil {
		return err
	}

	if err := cfg.HTTPServer.UnmarshalMap("http-server", cfgmap); err != nil {
		return err
	}

	if err := cfg.HTTPClient.UnmarshalMap("http-client", cfgmap); err != nil {
		return err
	}

	if err := cfg.Media.UnmarshalMap("media", cfgmap); err != nil {
		return err
	}

	if err := cfg.Cache.UnmarshalMap("cache", cfgmap); err != nil {
		return err
	}

	if ival, ok := cfgmap["log-level"]; ok {
		var err error
		cfg.LogLevel, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "log-level", err)
		}
	}

	if ival, ok := cfgmap["log-format"]; ok {
		var err error
		cfg.LogFormat, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "log-format", err)
		}
	}

	if ival, ok := cfgmap["log-timestamp-format"]; ok {
		var err error
		cfg.LogTimestampFormat, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "log-timestamp-format", err)
		}
	}

	if ival, ok := cfgmap["log-db-queries"]; ok {
		var err error
		cfg.LogDbQueries, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "log-db-queries", err)
		}
	}

	if ival, ok := cfgmap["log-client-ip"]; ok {
		var err error
		cfg.LogClientIP, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "log-client-ip", err)
		}
	}

	if ival, ok := cfgmap["request-id-header"]; ok {
		var err error
		cfg.RequestIDHeader, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "request-id-header", err)
		}
	}

	if ival, ok := cfgmap["config-path"]; ok {
		var err error
		cfg.ConfigPath, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "config-path", err)
		}
	}

	if ival, ok := cfgmap["application-name"]; ok {
		var err error
		cfg.ApplicationName, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "application-name", err)
		}
	}

	if ival, ok := cfgmap["landing-page-user"]; ok {
		var err error
		cfg.LandingPageUser, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "landing-page-user", err)
		}
	}

	if ival, ok := cfgmap["host"]; ok {
		var err error
		cfg.Host, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "host", err)
		}
	}

	if ival, ok := cfgmap["account-domain"]; ok {
		var err error
		cfg.AccountDomain, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "account-domain", err)
		}
	}

	if ival, ok := cfgmap["protocol"]; ok {
		var err error
		cfg.Protocol, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "protocol", err)
		}
	}

	if ival, ok := cfgmap["bind-address"]; ok {
		var err error
		cfg.BindAddress, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "bind-address", err)
		}
	}

	if ival, ok := cfgmap["port"]; ok {
		var err error
		cfg.Port, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "port", err)
		}
	}

	if ival, ok := cfgmap["trusted-proxies"]; ok {
		t, err := toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> []string for '%s': %w", ival, "trusted-proxies", err)
		}
		cfg.TrustedProxies = IPPrefixes{}
		for _, in := range t {
			if err := cfg.TrustedProxies.Set(in); err != nil {
				return fmt.Errorf("error parsing %#v for '%s': %w", ival, "trusted-proxies", err)
			}
		}
	}

	if ival, ok := cfgmap["software-version"]; ok {
		var err error
		cfg.SoftwareVersion, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "software-version", err)
		}
	}

	if ival, ok := cfgmap["web-template-base-dir"]; ok {
		var err error
		cfg.WebTemplateBaseDir, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "web-template-base-dir", err)
		}
	}

	if ival, ok := cfgmap["web-asset-base-dir"]; ok {
		var err error
		cfg.WebAssetBaseDir, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "web-asset-base-dir", err)
		}
	}

	if ival, ok := cfgmap["instance-federation-mode"]; ok {
		var err error
		cfg.InstanceFederationMode, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "instance-federation-mode", err)
		}
	}

	if ival, ok := cfgmap["instance-federation-spam-filter"]; ok {
		var err error
		cfg.InstanceFederationSpamFilter, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-federation-spam-filter", err)
		}
	}

	if ival, ok := cfgmap["instance-expose-peers"]; ok {
		var err error
		cfg.InstanceExposePeers, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-expose-peers", err)
		}
	}

	if ival, ok := cfgmap["instance-expose-blocklist"]; ok {
		var err error
		cfg.InstanceExposeBlocklist, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-expose-blocklist", err)
		}
	}

	if ival, ok := cfgmap["instance-expose-blocklist-web"]; ok {
		var err error
		cfg.InstanceExposeBlocklistWeb, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-expose-blocklist-web", err)
		}
	}

	if ival, ok := cfgmap["instance-expose-allowlist"]; ok {
		var err error
		cfg.InstanceExposeAllowlist, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-expose-allowlist", err)
		}
	}

	if ival, ok := cfgmap["instance-expose-allowlist-web"]; ok {
		var err error
		cfg.InstanceExposeAllowlistWeb, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-expose-allowlist-web", err)
		}
	}

	if ival, ok := cfgmap["instance-expose-public-timeline"]; ok {
		var err error
		cfg.InstanceExposePublicTimeline, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-expose-public-timeline", err)
		}
	}

	if ival, ok := cfgmap["instance-expose-custom-emojis"]; ok {
		var err error
		cfg.InstanceExposeCustomEmojis, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-expose-custom-emojis", err)
		}
	}

	if ival, ok := cfgmap["instance-directory-mode"]; ok {
		t, err := cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> string for '%s': %w", ival, "instance-directory-mode", err)
		}
		cfg.InstanceDirectoryMode = 0
		if err := cfg.InstanceDirectoryMode.Set(t); err != nil {
			return fmt.Errorf("error parsing %#v for '%s': %w", ival, "instance-directory-mode", err)
		}
	}

	if ival, ok := cfgmap["instance-deliver-to-shared-inboxes"]; ok {
		var err error
		cfg.InstanceDeliverToSharedInboxes, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-deliver-to-shared-inboxes", err)
		}
	}

	if ival, ok := cfgmap["instance-inject-mastodon-version"]; ok {
		var err error
		cfg.InstanceInjectMastodonVersion, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-inject-mastodon-version", err)
		}
	}

	if ival, ok := cfgmap["instance-languages"]; ok {
		t, err := toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> []string for '%s': %w", ival, "instance-languages", err)
		}
		cfg.InstanceLanguages = language.Languages{}
		for _, in := range t {
			if err := cfg.InstanceLanguages.Set(in); err != nil {
				return fmt.Errorf("error parsing %#v for '%s': %w", ival, "instance-languages", err)
			}
		}
	}

	if ival, ok := cfgmap["instance-stats-mode"]; ok {
		var err error
		cfg.InstanceStatsMode, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "instance-stats-mode", err)
		}
	}

	if ival, ok := cfgmap["instance-allow-backdating-statuses"]; ok {
		var err error
		cfg.InstanceAllowBackdatingStatuses, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-allow-backdating-statuses", err)
		}
	}

	if ival, ok := cfgmap["instance-robots-allow-indexing"]; ok {
		var err error
		cfg.InstanceRobotsAllowIndexing, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "instance-robots-allow-indexing", err)
		}
	}

	if ival, ok := cfgmap["accounts-registration-open"]; ok {
		var err error
		cfg.AccountsRegistrationOpen, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "accounts-registration-open", err)
		}
	}

	if ival, ok := cfgmap["accounts-reason-required"]; ok {
		var err error
		cfg.AccountsReasonRequired, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "accounts-reason-required", err)
		}
	}

	if ival, ok := cfgmap["accounts-registration-daily-limit"]; ok {
		var err error
		cfg.AccountsRegistrationDailyLimit, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "accounts-registration-daily-limit", err)
		}
	}

	if ival, ok := cfgmap["accounts-registration-backlog-limit"]; ok {
		var err error
		cfg.AccountsRegistrationBacklogLimit, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "accounts-registration-backlog-limit", err)
		}
	}

	if ival, ok := cfgmap["accounts-allow-custom-css"]; ok {
		var err error
		cfg.AccountsAllowCustomCSS, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "accounts-allow-custom-css", err)
		}
	}

	if ival, ok := cfgmap["accounts-custom-css-length"]; ok {
		var err error
		cfg.AccountsCustomCSSLength, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "accounts-custom-css-length", err)
		}
	}

	if ival, ok := cfgmap["accounts-max-profile-fields"]; ok {
		var err error
		cfg.AccountsMaxProfileFields, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "accounts-max-profile-fields", err)
		}
	}

	if ival, ok := cfgmap["storage-backend"]; ok {
		var err error
		cfg.StorageBackend, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-backend", err)
		}
	}

	if ival, ok := cfgmap["storage-local-base-path"]; ok {
		var err error
		cfg.StorageLocalBasePath, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-local-base-path", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-endpoint"]; ok {
		var err error
		cfg.StorageS3Endpoint, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-s3-endpoint", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-access-key"]; ok {
		var err error
		cfg.StorageS3AccessKey, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-s3-access-key", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-secret-key"]; ok {
		var err error
		cfg.StorageS3SecretKey, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-s3-secret-key", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-use-ssl"]; ok {
		var err error
		cfg.StorageS3UseSSL, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "storage-s3-use-ssl", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-bucket"]; ok {
		var err error
		cfg.StorageS3BucketName, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-s3-bucket", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-proxy"]; ok {
		var err error
		cfg.StorageS3Proxy, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "storage-s3-proxy", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-redirect-url"]; ok {
		var err error
		cfg.StorageS3RedirectURL, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-s3-redirect-url", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-bucket-lookup"]; ok {
		var err error
		cfg.StorageS3BucketLookup, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-s3-bucket-lookup", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-key-prefix"]; ok {
		var err error
		cfg.StorageS3KeyPrefix, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-s3-key-prefix", err)
		}
	}

	if ival, ok := cfgmap["storage-s3-region"]; ok {
		var err error
		cfg.StorageS3Region, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "storage-s3-region", err)
		}
	}

	if ival, ok := cfgmap["statuses-max-chars"]; ok {
		var err error
		cfg.StatusesMaxChars, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "statuses-max-chars", err)
		}
	}

	if ival, ok := cfgmap["statuses-poll-max-options"]; ok {
		var err error
		cfg.StatusesPollMaxOptions, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "statuses-poll-max-options", err)
		}
	}

	if ival, ok := cfgmap["statuses-poll-option-max-chars"]; ok {
		var err error
		cfg.StatusesPollOptionMaxChars, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "statuses-poll-option-max-chars", err)
		}
	}

	if ival, ok := cfgmap["statuses-media-max-files"]; ok {
		var err error
		cfg.StatusesMediaMaxFiles, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "statuses-media-max-files", err)
		}
	}

	if ival, ok := cfgmap["scheduled-statuses-max-total"]; ok {
		var err error
		cfg.ScheduledStatusesMaxTotal, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "scheduled-statuses-max-total", err)
		}
	}

	if ival, ok := cfgmap["scheduled-statuses-max-daily"]; ok {
		var err error
		cfg.ScheduledStatusesMaxDaily, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "scheduled-statuses-max-daily", err)
		}
	}

	if ival, ok := cfgmap["letsencrypt-enabled"]; ok {
		var err error
		cfg.LetsEncryptEnabled, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "letsencrypt-enabled", err)
		}
	}

	if ival, ok := cfgmap["letsencrypt-port"]; ok {
		var err error
		cfg.LetsEncryptPort, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "letsencrypt-port", err)
		}
	}

	if ival, ok := cfgmap["letsencrypt-cert-dir"]; ok {
		var err error
		cfg.LetsEncryptCertDir, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "letsencrypt-cert-dir", err)
		}
	}

	if ival, ok := cfgmap["letsencrypt-email-address"]; ok {
		var err error
		cfg.LetsEncryptEmailAddress, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "letsencrypt-email-address", err)
		}
	}

	if ival, ok := cfgmap["tls-certificate-chain"]; ok {
		var err error
		cfg.TLSCertificateChain, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "tls-certificate-chain", err)
		}
	}

	if ival, ok := cfgmap["tls-certificate-key"]; ok {
		var err error
		cfg.TLSCertificateKey, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "tls-certificate-key", err)
		}
	}

	if ival, ok := cfgmap["oidc-enabled"]; ok {
		var err error
		cfg.OIDCEnabled, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "oidc-enabled", err)
		}
	}

	if ival, ok := cfgmap["oidc-idp-name"]; ok {
		var err error
		cfg.OIDCIdpName, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "oidc-idp-name", err)
		}
	}

	if ival, ok := cfgmap["oidc-skip-verification"]; ok {
		var err error
		cfg.OIDCSkipVerification, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "oidc-skip-verification", err)
		}
	}

	if ival, ok := cfgmap["oidc-issuer"]; ok {
		var err error
		cfg.OIDCIssuer, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "oidc-issuer", err)
		}
	}

	if ival, ok := cfgmap["oidc-client-id"]; ok {
		var err error
		cfg.OIDCClientID, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "oidc-client-id", err)
		}
	}

	if ival, ok := cfgmap["oidc-client-secret"]; ok {
		var err error
		cfg.OIDCClientSecret, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "oidc-client-secret", err)
		}
	}

	if ival, ok := cfgmap["oidc-scopes"]; ok {
		var err error
		cfg.OIDCScopes, err = toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "[]string", "oidc-scopes", err)
		}
	}

	if ival, ok := cfgmap["oidc-link-existing"]; ok {
		var err error
		cfg.OIDCLinkExisting, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "oidc-link-existing", err)
		}
	}

	if ival, ok := cfgmap["oidc-allowed-groups"]; ok {
		var err error
		cfg.OIDCAllowedGroups, err = toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "[]string", "oidc-allowed-groups", err)
		}
	}

	if ival, ok := cfgmap["oidc-admin-groups"]; ok {
		var err error
		cfg.OIDCAdminGroups, err = toStringSlice(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "[]string", "oidc-admin-groups", err)
		}
	}

	if ival, ok := cfgmap["tracing-enabled"]; ok {
		var err error
		cfg.TracingEnabled, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "tracing-enabled", err)
		}
	}

	if ival, ok := cfgmap["metrics-enabled"]; ok {
		var err error
		cfg.MetricsEnabled, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "metrics-enabled", err)
		}
	}

	if ival, ok := cfgmap["smtp-host"]; ok {
		var err error
		cfg.SMTPHost, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "smtp-host", err)
		}
	}

	if ival, ok := cfgmap["smtp-port"]; ok {
		var err error
		cfg.SMTPPort, err = cast.ToIntE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "int", "smtp-port", err)
		}
	}

	if ival, ok := cfgmap["smtp-username"]; ok {
		var err error
		cfg.SMTPUsername, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "smtp-username", err)
		}
	}

	if ival, ok := cfgmap["smtp-password"]; ok {
		var err error
		cfg.SMTPPassword, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "smtp-password", err)
		}
	}

	if ival, ok := cfgmap["smtp-from"]; ok {
		var err error
		cfg.SMTPFrom, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "smtp-from", err)
		}
	}

	if ival, ok := cfgmap["smtp-from-display-name"]; ok {
		var err error
		cfg.SMTPFromDisplayName, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "smtp-from-display-name", err)
		}
	}

	if ival, ok := cfgmap["smtp-disclose-recipients"]; ok {
		var err error
		cfg.SMTPDiscloseRecipients, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "smtp-disclose-recipients", err)
		}
	}

	if ival, ok := cfgmap["syslog-enabled"]; ok {
		var err error
		cfg.SyslogEnabled, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "syslog-enabled", err)
		}
	}

	if ival, ok := cfgmap["syslog-protocol"]; ok {
		var err error
		cfg.SyslogProtocol, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "syslog-protocol", err)
		}
	}

	if ival, ok := cfgmap["syslog-address"]; ok {
		var err error
		cfg.SyslogAddress, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "syslog-address", err)
		}
	}

	if ival, ok := cfgmap["syslog-mirror"]; ok {
		var err error
		cfg.SyslogMirror, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "syslog-mirror", err)
		}
	}

	if ival, ok := cfgmap["syslog-msg-length"]; ok {
		var err error
		cfg.SyslogMsgLength, err = cast.ToUint32E(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "uint32", "syslog-msg-length", err)
		}
	}

	if ival, ok := cfgmap["username"]; ok {
		var err error
		cfg.AdminAccountUsername, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "username", err)
		}
	}

	if ival, ok := cfgmap["email"]; ok {
		var err error
		cfg.AdminAccountEmail, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "email", err)
		}
	}

	if ival, ok := cfgmap["password"]; ok {
		var err error
		cfg.AdminAccountPassword, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "password", err)
		}
	}

	if ival, ok := cfgmap["path"]; ok {
		var err error
		cfg.AdminTransPath, err = cast.ToStringE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "string", "path", err)
		}
	}

	if ival, ok := cfgmap["dry-run"]; ok {
		var err error
		cfg.AdminMediaPruneDryRun, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "dry-run", err)
		}
	}

	if ival, ok := cfgmap["local-only"]; ok {
		var err error
		cfg.AdminMediaListLocalOnly, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "local-only", err)
		}
	}

	if ival, ok := cfgmap["remote-only"]; ok {
		var err error
		cfg.AdminMediaListRemoteOnly, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "remote-only", err)
		}
	}

	if ival, ok := cfgmap["skip-db-setup"]; ok {
		var err error
		cfg.TestrigSkipDBSetup, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "skip-db-setup", err)
		}
	}

	if ival, ok := cfgmap["skip-db-teardown"]; ok {
		var err error
		cfg.TestrigSkipDBTeardown, err = cast.ToBoolE(ival)
		if err != nil {
			return fmt.Errorf("error casting %#v -> %s for '%s': %w", ival, "bool", "skip-db-teardown", err)
		}
	}

	return nil
}

// GetDatabasePostgresPort safely fetches the Configuration value for state's 'Database.Postgres.Port' field
func (st *ConfigState) GetDatabasePostgresPort() (v uint16) {
	return st.config.Database.Postgres.Port
}

// SetDatabasePostgresPort safely sets the Configuration value for state's 'Database.Postgres.Port' field
func (st *ConfigState) SetDatabasePostgresPort(v uint16) {
	st.config.Database.Postgres.Port = v
	st.reloadToViper()
}

// GetDatabasePostgresPort safely fetches the value for global configuration 'Database.Postgres.Port' field
func GetDatabasePostgresPort() uint16 { return global.GetDatabasePostgresPort() }

// SetDatabasePostgresPort safely sets the value for global configuration 'Database.Postgres.Port' field
func SetDatabasePostgresPort(v uint16) { global.SetDatabasePostgresPort(v) }

// GetDatabasePostgresUser safely fetches the Configuration value for state's 'Database.Postgres.User' field
func (st *ConfigState) GetDatabasePostgresUser() (v string) {
	return st.config.Database.Postgres.User
}

// SetDatabasePostgresUser safely sets the Configuration value for state's 'Database.Postgres.User' field
func (st *ConfigState) SetDatabasePostgresUser(v string) {
	st.config.Database.Postgres.User = v
	st.reloadToViper()
}

// GetDatabasePostgresUser safely fetches the value for global configuration 'Database.Postgres.User' field
func GetDatabasePostgresUser() string { return global.GetDatabasePostgresUser() }

// SetDatabasePostgresUser safely sets the value for global configuration 'Database.Postgres.User' field
func SetDatabasePostgresUser(v string) { global.SetDatabasePostgresUser(v) }

// GetDatabasePostgresPassword safely fetches the Configuration value for state's 'Database.Postgres.Password' field
func (st *ConfigState) GetDatabasePostgresPassword() (v string) {
	return st.config.Database.Postgres.Password
}

// SetDatabasePostgresPassword safely sets the Configuration value for state's 'Database.Postgres.Password' field
func (st *ConfigState) SetDatabasePostgresPassword(v string) {
	st.config.Database.Postgres.Password = v
	st.reloadToViper()
}

// GetDatabasePostgresPassword safely fetches the value for global configuration 'Database.Postgres.Password' field
func GetDatabasePostgresPassword() string { return global.GetDatabasePostgresPassword() }

// SetDatabasePostgresPassword safely sets the value for global configuration 'Database.Postgres.Password' field
func SetDatabasePostgresPassword(v string) { global.SetDatabasePostgresPassword(v) }

// GetDatabasePostgresDatabase safely fetches the Configuration value for state's 'Database.Postgres.Database' field
func (st *ConfigState) GetDatabasePostgresDatabase() (v string) {
	return st.config.Database.Postgres.Database
}

// SetDatabasePostgresDatabase safely sets the Configuration value for state's 'Database.Postgres.Database' field
func (st *ConfigState) SetDatabasePostgresDatabase(v string) {
	st.config.Database.Postgres.Database = v
	st.reloadToViper()
}

// GetDatabasePostgresDatabase safely fetches the value for global configuration 'Database.Postgres.Database' field
func GetDatabasePostgresDatabase() string { return global.GetDatabasePostgresDatabase() }

// SetDatabasePostgresDatabase safely sets the value for global configuration 'Database.Postgres.Database' field
func SetDatabasePostgresDatabase(v string) { global.SetDatabasePostgresDatabase(v) }

// GetDatabasePostgresTLSMode safely fetches the Configuration value for state's 'Database.Postgres.TLSMode' field
func (st *ConfigState) GetDatabasePostgresTLSMode() (v string) {
	return st.config.Database.Postgres.TLSMode
}

// SetDatabasePostgresTLSMode safely sets the Configuration value for state's 'Database.Postgres.TLSMode' field
func (st *ConfigState) SetDatabasePostgresTLSMode(v string) {
	st.config.Database.Postgres.TLSMode = v
	st.reloadToViper()
}

// GetDatabasePostgresTLSMode safely fetches the value for global configuration 'Database.Postgres.TLSMode' field
func GetDatabasePostgresTLSMode() string { return global.GetDatabasePostgresTLSMode() }

// SetDatabasePostgresTLSMode safely sets the value for global configuration 'Database.Postgres.TLSMode' field
func SetDatabasePostgresTLSMode(v string) { global.SetDatabasePostgresTLSMode(v) }

// GetDatabasePostgresTLSCACert safely fetches the Configuration value for state's 'Database.Postgres.TLSCACert' field
func (st *ConfigState) GetDatabasePostgresTLSCACert() (v string) {
	return st.config.Database.Postgres.TLSCACert
}

// SetDatabasePostgresTLSCACert safely sets the Configuration value for state's 'Database.Postgres.TLSCACert' field
func (st *ConfigState) SetDatabasePostgresTLSCACert(v string) {
	st.config.Database.Postgres.TLSCACert = v
	st.reloadToViper()
}

// GetDatabasePostgresTLSCACert safely fetches the value for global configuration 'Database.Postgres.TLSCACert' field
func GetDatabasePostgresTLSCACert() string { return global.GetDatabasePostgresTLSCACert() }

// SetDatabasePostgresTLSCACert safely sets the value for global configuration 'Database.Postgres.TLSCACert' field
func SetDatabasePostgresTLSCACert(v string) { global.SetDatabasePostgresTLSCACert(v) }

// GetDatabasePostgresConnectionString safely fetches the Configuration value for state's 'Database.Postgres.ConnectionString' field
func (st *ConfigState) GetDatabasePostgresConnectionString() (v string) {
	return st.config.Database.Postgres.ConnectionString
}

// SetDatabasePostgresConnectionString safely sets the Configuration value for state's 'Database.Postgres.ConnectionString' field
func (st *ConfigState) SetDatabasePostgresConnectionString(v string) {
	st.config.Database.Postgres.ConnectionString = v
	st.reloadToViper()
}

// GetDatabasePostgresConnectionString safely fetches the value for global configuration 'Database.Postgres.ConnectionString' field
func GetDatabasePostgresConnectionString() string {
	return global.GetDatabasePostgresConnectionString()
}

// SetDatabasePostgresConnectionString safely sets the value for global configuration 'Database.Postgres.ConnectionString' field
func SetDatabasePostgresConnectionString(v string) { global.SetDatabasePostgresConnectionString(v) }

// GetDatabaseSQLiteJournalMode safely fetches the Configuration value for state's 'Database.SQLite.JournalMode' field
func (st *ConfigState) GetDatabaseSQLiteJournalMode() (v string) {
	return st.config.Database.SQLite.JournalMode
}

// SetDatabaseSQLiteJournalMode safely sets the Configuration value for state's 'Database.SQLite.JournalMode' field
func (st *ConfigState) SetDatabaseSQLiteJournalMode(v string) {
	st.config.Database.SQLite.JournalMode = v
	st.reloadToViper()
}

// GetDatabaseSQLiteJournalMode safely fetches the value for global configuration 'Database.SQLite.JournalMode' field
func GetDatabaseSQLiteJournalMode() string { return global.GetDatabaseSQLiteJournalMode() }

// SetDatabaseSQLiteJournalMode safely sets the value for global configuration 'Database.SQLite.JournalMode' field
func SetDatabaseSQLiteJournalMode(v string) { global.SetDatabaseSQLiteJournalMode(v) }

// GetDatabaseSQLiteSynchronous safely fetches the Configuration value for state's 'Database.SQLite.Synchronous' field
func (st *ConfigState) GetDatabaseSQLiteSynchronous() (v string) {
	return st.config.Database.SQLite.Synchronous
}

// SetDatabaseSQLiteSynchronous safely sets the Configuration value for state's 'Database.SQLite.Synchronous' field
func (st *ConfigState) SetDatabaseSQLiteSynchronous(v string) {
	st.config.Database.SQLite.Synchronous = v
	st.reloadToViper()
}

// GetDatabaseSQLiteSynchronous safely fetches the value for global configuration 'Database.SQLite.Synchronous' field
func GetDatabaseSQLiteSynchronous() string { return global.GetDatabaseSQLiteSynchronous() }

// SetDatabaseSQLiteSynchronous safely sets the value for global configuration 'Database.SQLite.Synchronous' field
func SetDatabaseSQLiteSynchronous(v string) { global.SetDatabaseSQLiteSynchronous(v) }

// GetDatabaseSQLiteCacheSize safely fetches the Configuration value for state's 'Database.SQLite.CacheSize' field
func (st *ConfigState) GetDatabaseSQLiteCacheSize() (v bytesize.Size) {
	return st.config.Database.SQLite.CacheSize
}

// SetDatabaseSQLiteCacheSize safely sets the Configuration value for state's 'Database.SQLite.CacheSize' field
func (st *ConfigState) SetDatabaseSQLiteCacheSize(v bytesize.Size) {
	st.config.Database.SQLite.CacheSize = v
	st.reloadToViper()
}

// GetDatabaseSQLiteCacheSize safely fetches the value for global configuration 'Database.SQLite.CacheSize' field
func GetDatabaseSQLiteCacheSize() bytesize.Size { return global.GetDatabaseSQLiteCacheSize() }

// SetDatabaseSQLiteCacheSize safely sets the value for global configuration 'Database.SQLite.CacheSize' field
func SetDatabaseSQLiteCacheSize(v bytesize.Size) { global.SetDatabaseSQLiteCacheSize(v) }

// GetDatabaseSQLiteBusyTimeout safely fetches the Configuration value for state's 'Database.SQLite.BusyTimeout' field
func (st *ConfigState) GetDatabaseSQLiteBusyTimeout() (v time.Duration) {
	return st.config.Database.SQLite.BusyTimeout
}

// SetDatabaseSQLiteBusyTimeout safely sets the Configuration value for state's 'Database.SQLite.BusyTimeout' field
func (st *ConfigState) SetDatabaseSQLiteBusyTimeout(v time.Duration) {
	st.config.Database.SQLite.BusyTimeout = v
	st.reloadToViper()
}

// GetDatabaseSQLiteBusyTimeout safely fetches the value for global configuration 'Database.SQLite.BusyTimeout' field
func GetDatabaseSQLiteBusyTimeout() time.Duration { return global.GetDatabaseSQLiteBusyTimeout() }

// SetDatabaseSQLiteBusyTimeout safely sets the value for global configuration 'Database.SQLite.BusyTimeout' field
func SetDatabaseSQLiteBusyTimeout(v time.Duration) { global.SetDatabaseSQLiteBusyTimeout(v) }

// GetDatabaseType safely fetches the Configuration value for state's 'Database.Type' field
func (st *ConfigState) GetDatabaseType() (v string) {
	return st.config.Database.Type
}

// SetDatabaseType safely sets the Configuration value for state's 'Database.Type' field
func (st *ConfigState) SetDatabaseType(v string) {
	st.config.Database.Type = v
	st.reloadToViper()
}

// GetDatabaseType safely fetches the value for global configuration 'Database.Type' field
func GetDatabaseType() string { return global.GetDatabaseType() }

// SetDatabaseType safely sets the value for global configuration 'Database.Type' field
func SetDatabaseType(v string) { global.SetDatabaseType(v) }

// GetDatabaseAddress safely fetches the Configuration value for state's 'Database.Address' field
func (st *ConfigState) GetDatabaseAddress() (v string) {
	return st.config.Database.Address
}

// SetDatabaseAddress safely sets the Configuration value for state's 'Database.Address' field
func (st *ConfigState) SetDatabaseAddress(v string) {
	st.config.Database.Address = v
	st.reloadToViper()
}

// GetDatabaseAddress safely fetches the value for global configuration 'Database.Address' field
func GetDatabaseAddress() string { return global.GetDatabaseAddress() }

// SetDatabaseAddress safely sets the value for global configuration 'Database.Address' field
func SetDatabaseAddress(v string) { global.SetDatabaseAddress(v) }

// GetDatabaseMaxOpenConnsMultiplier safely fetches the Configuration value for state's 'Database.MaxOpenConnsMultiplier' field
func (st *ConfigState) GetDatabaseMaxOpenConnsMultiplier() (v int) {
	return st.config.Database.MaxOpenConnsMultiplier
}

// SetDatabaseMaxOpenConnsMultiplier safely sets the Configuration value for state's 'Database.MaxOpenConnsMultiplier' field
func (st *ConfigState) SetDatabaseMaxOpenConnsMultiplier(v int) {
	st.config.Database.MaxOpenConnsMultiplier = v
	st.reloadToViper()
}

// GetDatabaseMaxOpenConnsMultiplier safely fetches the value for global configuration 'Database.MaxOpenConnsMultiplier' field
func GetDatabaseMaxOpenConnsMultiplier() int { return global.GetDatabaseMaxOpenConnsMultiplier() }

// SetDatabaseMaxOpenConnsMultiplier safely sets the value for global configuration 'Database.MaxOpenConnsMultiplier' field
func SetDatabaseMaxOpenConnsMultiplier(v int) { global.SetDatabaseMaxOpenConnsMultiplier(v) }

// GetAdvancedRateLimitRequests safely fetches the Configuration value for state's 'Advanced.RateLimit.Requests' field
func (st *ConfigState) GetAdvancedRateLimitRequests() (v int) {
	return st.config.Advanced.RateLimit.Requests
}

// SetAdvancedRateLimitRequests safely sets the Configuration value for state's 'Advanced.RateLimit.Requests' field
func (st *ConfigState) SetAdvancedRateLimitRequests(v int) {
	st.config.Advanced.RateLimit.Requests = v
	st.reloadToViper()
}

// GetAdvancedRateLimitRequests safely fetches the value for global configuration 'Advanced.RateLimit.Requests' field
func GetAdvancedRateLimitRequests() int { return global.GetAdvancedRateLimitRequests() }

// SetAdvancedRateLimitRequests safely sets the value for global configuration 'Advanced.RateLimit.Requests' field
func SetAdvancedRateLimitRequests(v int) { global.SetAdvancedRateLimitRequests(v) }

// GetAdvancedRateLimitExceptions safely fetches the Configuration value for state's 'Advanced.RateLimit.Exceptions' field
func (st *ConfigState) GetAdvancedRateLimitExceptions() (v IPPrefixes) {
	return st.config.Advanced.RateLimit.Exceptions
}

// SetAdvancedRateLimitExceptions safely sets the Configuration value for state's 'Advanced.RateLimit.Exceptions' field
func (st *ConfigState) SetAdvancedRateLimitExceptions(v IPPrefixes) {
	st.config.Advanced.RateLimit.Exceptions = v
	st.reloadToViper()
}

// GetAdvancedRateLimitExceptions safely fetches the value for global configuration 'Advanced.RateLimit.Exceptions' field
func GetAdvancedRateLimitExceptions() IPPrefixes { return global.GetAdvancedRateLimitExceptions() }

// SetAdvancedRateLimitExceptions safely sets the value for global configuration 'Advanced.RateLimit.Exceptions' field
func SetAdvancedRateLimitExceptions(v IPPrefixes) { global.SetAdvancedRateLimitExceptions(v) }

// GetAdvancedThrottlingMultiplier safely fetches the Configuration value for state's 'Advanced.Throttling.Multiplier' field
func (st *ConfigState) GetAdvancedThrottlingMultiplier() (v int) {
	return st.config.Advanced.Throttling.Multiplier
}

// SetAdvancedThrottlingMultiplier safely sets the Configuration value for state's 'Advanced.Throttling.Multiplier' field
func (st *ConfigState) SetAdvancedThrottlingMultiplier(v int) {
	st.config.Advanced.Throttling.Multiplier = v
	st.reloadToViper()
}

// GetAdvancedThrottlingMultiplier safely fetches the value for global configuration 'Advanced.Throttling.Multiplier' field
func GetAdvancedThrottlingMultiplier() int { return global.GetAdvancedThrottlingMultiplier() }

// SetAdvancedThrottlingMultiplier safely sets the value for global configuration 'Advanced.Throttling.Multiplier' field
func SetAdvancedThrottlingMultiplier(v int) { global.SetAdvancedThrottlingMultiplier(v) }

// GetAdvancedThrottlingRetryAfter safely fetches the Configuration value for state's 'Advanced.Throttling.RetryAfter' field
func (st *ConfigState) GetAdvancedThrottlingRetryAfter() (v time.Duration) {
	return st.config.Advanced.Throttling.RetryAfter
}

// SetAdvancedThrottlingRetryAfter safely sets the Configuration value for state's 'Advanced.Throttling.RetryAfter' field
func (st *ConfigState) SetAdvancedThrottlingRetryAfter(v time.Duration) {
	st.config.Advanced.Throttling.RetryAfter = v
	st.reloadToViper()
}

// GetAdvancedThrottlingRetryAfter safely fetches the value for global configuration 'Advanced.Throttling.RetryAfter' field
func GetAdvancedThrottlingRetryAfter() time.Duration { return global.GetAdvancedThrottlingRetryAfter() }

// SetAdvancedThrottlingRetryAfter safely sets the value for global configuration 'Advanced.Throttling.RetryAfter' field
func SetAdvancedThrottlingRetryAfter(v time.Duration) { global.SetAdvancedThrottlingRetryAfter(v) }

// GetAdvancedCookiesSamesite safely fetches the Configuration value for state's 'Advanced.CookiesSamesite' field
func (st *ConfigState) GetAdvancedCookiesSamesite() (v string) {
	return st.config.Advanced.CookiesSamesite
}

// SetAdvancedCookiesSamesite safely sets the Configuration value for state's 'Advanced.CookiesSamesite' field
func (st *ConfigState) SetAdvancedCookiesSamesite(v string) {
	st.config.Advanced.CookiesSamesite = v
	st.reloadToViper()
}

// GetAdvancedCookiesSamesite safely fetches the value for global configuration 'Advanced.CookiesSamesite' field
func GetAdvancedCookiesSamesite() string { return global.GetAdvancedCookiesSamesite() }

// SetAdvancedCookiesSamesite safely sets the value for global configuration 'Advanced.CookiesSamesite' field
func SetAdvancedCookiesSamesite(v string) { global.SetAdvancedCookiesSamesite(v) }

// GetAdvancedSenderMultiplier safely fetches the Configuration value for state's 'Advanced.SenderMultiplier' field
func (st *ConfigState) GetAdvancedSenderMultiplier() (v int) {
	return st.config.Advanced.SenderMultiplier
}

// SetAdvancedSenderMultiplier safely sets the Configuration value for state's 'Advanced.SenderMultiplier' field
func (st *ConfigState) SetAdvancedSenderMultiplier(v int) {
	st.config.Advanced.SenderMultiplier = v
	st.reloadToViper()
}

// GetAdvancedSenderMultiplier safely fetches the value for global configuration 'Advanced.SenderMultiplier' field
func GetAdvancedSenderMultiplier() int { return global.GetAdvancedSenderMultiplier() }

// SetAdvancedSenderMultiplier safely sets the value for global configuration 'Advanced.SenderMultiplier' field
func SetAdvancedSenderMultiplier(v int) { global.SetAdvancedSenderMultiplier(v) }

// GetAdvancedCSPExtraURIs safely fetches the Configuration value for state's 'Advanced.CSPExtraURIs' field
func (st *ConfigState) GetAdvancedCSPExtraURIs() (v []string) {
	return st.config.Advanced.CSPExtraURIs
}

// SetAdvancedCSPExtraURIs safely sets the Configuration value for state's 'Advanced.CSPExtraURIs' field
func (st *ConfigState) SetAdvancedCSPExtraURIs(v []string) {
	st.config.Advanced.CSPExtraURIs = v
	st.reloadToViper()
}

// GetAdvancedCSPExtraURIs safely fetches the value for global configuration 'Advanced.CSPExtraURIs' field
func GetAdvancedCSPExtraURIs() []string { return global.GetAdvancedCSPExtraURIs() }

// SetAdvancedCSPExtraURIs safely sets the value for global configuration 'Advanced.CSPExtraURIs' field
func SetAdvancedCSPExtraURIs(v []string) { global.SetAdvancedCSPExtraURIs(v) }

// GetAdvancedHeaderFilterMode safely fetches the Configuration value for state's 'Advanced.HeaderFilterMode' field
func (st *ConfigState) GetAdvancedHeaderFilterMode() (v string) {
	return st.config.Advanced.HeaderFilterMode
}

// SetAdvancedHeaderFilterMode safely sets the Configuration value for state's 'Advanced.HeaderFilterMode' field
func (st *ConfigState) SetAdvancedHeaderFilterMode(v string) {
	st.config.Advanced.HeaderFilterMode = v
	st.reloadToViper()
}

// GetAdvancedHeaderFilterMode safely fetches the value for global configuration 'Advanced.HeaderFilterMode' field
func GetAdvancedHeaderFilterMode() string { return global.GetAdvancedHeaderFilterMode() }

// SetAdvancedHeaderFilterMode safely sets the value for global configuration 'Advanced.HeaderFilterMode' field
func SetAdvancedHeaderFilterMode(v string) { global.SetAdvancedHeaderFilterMode(v) }

// GetHTTPServerMaxMultipartMemory safely fetches the Configuration value for state's 'HTTPServer.MaxMultipartMemory' field
func (st *ConfigState) GetHTTPServerMaxMultipartMemory() (v bytesize.Size) {
	return st.config.HTTPServer.MaxMultipartMemory
}

// SetHTTPServerMaxMultipartMemory safely sets the Configuration value for state's 'HTTPServer.MaxMultipartMemory' field
func (st *ConfigState) SetHTTPServerMaxMultipartMemory(v bytesize.Size) {
	st.config.HTTPServer.MaxMultipartMemory = v
	st.reloadToViper()
}

// GetHTTPServerMaxMultipartMemory safely fetches the value for global configuration 'HTTPServer.MaxMultipartMemory' field
func GetHTTPServerMaxMultipartMemory() bytesize.Size { return global.GetHTTPServerMaxMultipartMemory() }

// SetHTTPServerMaxMultipartMemory safely sets the value for global configuration 'HTTPServer.MaxMultipartMemory' field
func SetHTTPServerMaxMultipartMemory(v bytesize.Size) { global.SetHTTPServerMaxMultipartMemory(v) }

// GetHTTPServerUseH2C safely fetches the Configuration value for state's 'HTTPServer.UseH2C' field
func (st *ConfigState) GetHTTPServerUseH2C() (v bool) {
	return st.config.HTTPServer.UseH2C
}

// SetHTTPServerUseH2C safely sets the Configuration value for state's 'HTTPServer.UseH2C' field
func (st *ConfigState) SetHTTPServerUseH2C(v bool) {
	st.config.HTTPServer.UseH2C = v
	st.reloadToViper()
}

// GetHTTPServerUseH2C safely fetches the value for global configuration 'HTTPServer.UseH2C' field
func GetHTTPServerUseH2C() bool { return global.GetHTTPServerUseH2C() }

// SetHTTPServerUseH2C safely sets the value for global configuration 'HTTPServer.UseH2C' field
func SetHTTPServerUseH2C(v bool) { global.SetHTTPServerUseH2C(v) }

// GetHTTPServerReadTimeout safely fetches the Configuration value for state's 'HTTPServer.ReadTimeout' field
func (st *ConfigState) GetHTTPServerReadTimeout() (v time.Duration) {
	return st.config.HTTPServer.ReadTimeout
}

// SetHTTPServerReadTimeout safely sets the Configuration value for state's 'HTTPServer.ReadTimeout' field
func (st *ConfigState) SetHTTPServerReadTimeout(v time.Duration) {
	st.config.HTTPServer.ReadTimeout = v
	st.reloadToViper()
}

// GetHTTPServerReadTimeout safely fetches the value for global configuration 'HTTPServer.ReadTimeout' field
func GetHTTPServerReadTimeout() time.Duration { return global.GetHTTPServerReadTimeout() }

// SetHTTPServerReadTimeout safely sets the value for global configuration 'HTTPServer.ReadTimeout' field
func SetHTTPServerReadTimeout(v time.Duration) { global.SetHTTPServerReadTimeout(v) }

// GetHTTPServerReadHeaderTimeout safely fetches the Configuration value for state's 'HTTPServer.ReadHeaderTimeout' field
func (st *ConfigState) GetHTTPServerReadHeaderTimeout() (v time.Duration) {
	return st.config.HTTPServer.ReadHeaderTimeout
}

// SetHTTPServerReadHeaderTimeout safely sets the Configuration value for state's 'HTTPServer.ReadHeaderTimeout' field
func (st *ConfigState) SetHTTPServerReadHeaderTimeout(v time.Duration) {
	st.config.HTTPServer.ReadHeaderTimeout = v
	st.reloadToViper()
}

// GetHTTPServerReadHeaderTimeout safely fetches the value for global configuration 'HTTPServer.ReadHeaderTimeout' field
func GetHTTPServerReadHeaderTimeout() time.Duration { return global.GetHTTPServerReadHeaderTimeout() }

// SetHTTPServerReadHeaderTimeout safely sets the value for global configuration 'HTTPServer.ReadHeaderTimeout' field
func SetHTTPServerReadHeaderTimeout(v time.Duration) { global.SetHTTPServerReadHeaderTimeout(v) }

// GetHTTPServerWriteTimeout safely fetches the Configuration value for state's 'HTTPServer.WriteTimeout' field
func (st *ConfigState) GetHTTPServerWriteTimeout() (v time.Duration) {
	return st.config.HTTPServer.WriteTimeout
}

// SetHTTPServerWriteTimeout safely sets the Configuration value for state's 'HTTPServer.WriteTimeout' field
func (st *ConfigState) SetHTTPServerWriteTimeout(v time.Duration) {
	st.config.HTTPServer.WriteTimeout = v
	st.reloadToViper()
}

// GetHTTPServerWriteTimeout safely fetches the value for global configuration 'HTTPServer.WriteTimeout' field
func GetHTTPServerWriteTimeout() time.Duration { return global.GetHTTPServerWriteTimeout() }

// SetHTTPServerWriteTimeout safely sets the value for global configuration 'HTTPServer.WriteTimeout' field
func SetHTTPServerWriteTimeout(v time.Duration) { global.SetHTTPServerWriteTimeout(v) }

// GetHTTPServerIdleTimeout safely fetches the Configuration value for state's 'HTTPServer.IdleTimeout' field
func (st *ConfigState) GetHTTPServerIdleTimeout() (v time.Duration) {
	return st.config.HTTPServer.IdleTimeout
}

// SetHTTPServerIdleTimeout safely sets the Configuration value for state's 'HTTPServer.IdleTimeout' field
func (st *ConfigState) SetHTTPServerIdleTimeout(v time.Duration) {
	st.config.HTTPServer.IdleTimeout = v
	st.reloadToViper()
}

// GetHTTPServerIdleTimeout safely fetches the value for global configuration 'HTTPServer.IdleTimeout' field
func GetHTTPServerIdleTimeout() time.Duration { return global.GetHTTPServerIdleTimeout() }

// SetHTTPServerIdleTimeout safely sets the value for global configuration 'HTTPServer.IdleTimeout' field
func SetHTTPServerIdleTimeout(v time.Duration) { global.SetHTTPServerIdleTimeout(v) }

// GetHTTPServerMaxHeaderBytes safely fetches the Configuration value for state's 'HTTPServer.MaxHeaderBytes' field
func (st *ConfigState) GetHTTPServerMaxHeaderBytes() (v bytesize.Size) {
	return st.config.HTTPServer.MaxHeaderBytes
}

// SetHTTPServerMaxHeaderBytes safely sets the Configuration value for state's 'HTTPServer.MaxHeaderBytes' field
func (st *ConfigState) SetHTTPServerMaxHeaderBytes(v bytesize.Size) {
	st.config.HTTPServer.MaxHeaderBytes = v
	st.reloadToViper()
}

// GetHTTPServerMaxHeaderBytes safely fetches the value for global configuration 'HTTPServer.MaxHeaderBytes' field
func GetHTTPServerMaxHeaderBytes() bytesize.Size { return global.GetHTTPServerMaxHeaderBytes() }

// SetHTTPServerMaxHeaderBytes safely sets the value for global configuration 'HTTPServer.MaxHeaderBytes' field
func SetHTTPServerMaxHeaderBytes(v bytesize.Size) { global.SetHTTPServerMaxHeaderBytes(v) }

// GetHTTPServerMaxConcurrentStreams safely fetches the Configuration value for state's 'HTTPServer.MaxConcurrentStreams' field
func (st *ConfigState) GetHTTPServerMaxConcurrentStreams() (v int) {
	return st.config.HTTPServer.MaxConcurrentStreams
}

// SetHTTPServerMaxConcurrentStreams safely sets the Configuration value for state's 'HTTPServer.MaxConcurrentStreams' field
func (st *ConfigState) SetHTTPServerMaxConcurrentStreams(v int) {
	st.config.HTTPServer.MaxConcurrentStreams = v
	st.reloadToViper()
}

// GetHTTPServerMaxConcurrentStreams safely fetches the value for global configuration 'HTTPServer.MaxConcurrentStreams' field
func GetHTTPServerMaxConcurrentStreams() int { return global.GetHTTPServerMaxConcurrentStreams() }

// SetHTTPServerMaxConcurrentStreams safely sets the value for global configuration 'HTTPServer.MaxConcurrentStreams' field
func SetHTTPServerMaxConcurrentStreams(v int) { global.SetHTTPServerMaxConcurrentStreams(v) }

// GetHTTPServerMaxDecoderHeaderTableSize safely fetches the Configuration value for state's 'HTTPServer.MaxDecoderHeaderTableSize' field
func (st *ConfigState) GetHTTPServerMaxDecoderHeaderTableSize() (v bytesize.Size) {
	return st.config.HTTPServer.MaxDecoderHeaderTableSize
}

// SetHTTPServerMaxDecoderHeaderTableSize safely sets the Configuration value for state's 'HTTPServer.MaxDecoderHeaderTableSize' field
func (st *ConfigState) SetHTTPServerMaxDecoderHeaderTableSize(v bytesize.Size) {
	st.config.HTTPServer.MaxDecoderHeaderTableSize = v
	st.reloadToViper()
}

// GetHTTPServerMaxDecoderHeaderTableSize safely fetches the value for global configuration 'HTTPServer.MaxDecoderHeaderTableSize' field
func GetHTTPServerMaxDecoderHeaderTableSize() bytesize.Size {
	return global.GetHTTPServerMaxDecoderHeaderTableSize()
}

// SetHTTPServerMaxDecoderHeaderTableSize safely sets the value for global configuration 'HTTPServer.MaxDecoderHeaderTableSize' field
func SetHTTPServerMaxDecoderHeaderTableSize(v bytesize.Size) {
	global.SetHTTPServerMaxDecoderHeaderTableSize(v)
}

// GetHTTPServerMaxEncoderHeaderTableSize safely fetches the Configuration value for state's 'HTTPServer.MaxEncoderHeaderTableSize' field
func (st *ConfigState) GetHTTPServerMaxEncoderHeaderTableSize() (v bytesize.Size) {
	return st.config.HTTPServer.MaxEncoderHeaderTableSize
}

// SetHTTPServerMaxEncoderHeaderTableSize safely sets the Configuration value for state's 'HTTPServer.MaxEncoderHeaderTableSize' field
func (st *ConfigState) SetHTTPServerMaxEncoderHeaderTableSize(v bytesize.Size) {
	st.config.HTTPServer.MaxEncoderHeaderTableSize = v
	st.reloadToViper()
}

// GetHTTPServerMaxEncoderHeaderTableSize safely fetches the value for global configuration 'HTTPServer.MaxEncoderHeaderTableSize' field
func GetHTTPServerMaxEncoderHeaderTableSize() bytesize.Size {
	return global.GetHTTPServerMaxEncoderHeaderTableSize()
}

// SetHTTPServerMaxEncoderHeaderTableSize safely sets the value for global configuration 'HTTPServer.MaxEncoderHeaderTableSize' field
func SetHTTPServerMaxEncoderHeaderTableSize(v bytesize.Size) {
	global.SetHTTPServerMaxEncoderHeaderTableSize(v)
}

// GetHTTPServerMaxReadFrameSize safely fetches the Configuration value for state's 'HTTPServer.MaxReadFrameSize' field
func (st *ConfigState) GetHTTPServerMaxReadFrameSize() (v bytesize.Size) {
	return st.config.HTTPServer.MaxReadFrameSize
}

// SetHTTPServerMaxReadFrameSize safely sets the Configuration value for state's 'HTTPServer.MaxReadFrameSize' field
func (st *ConfigState) SetHTTPServerMaxReadFrameSize(v bytesize.Size) {
	st.config.HTTPServer.MaxReadFrameSize = v
	st.reloadToViper()
}

// GetHTTPServerMaxReadFrameSize safely fetches the value for global configuration 'HTTPServer.MaxReadFrameSize' field
func GetHTTPServerMaxReadFrameSize() bytesize.Size { return global.GetHTTPServerMaxReadFrameSize() }

// SetHTTPServerMaxReadFrameSize safely sets the value for global configuration 'HTTPServer.MaxReadFrameSize' field
func SetHTTPServerMaxReadFrameSize(v bytesize.Size) { global.SetHTTPServerMaxReadFrameSize(v) }

// GetHTTPServerMaxReceiveBufferPerConnection safely fetches the Configuration value for state's 'HTTPServer.MaxReceiveBufferPerConnection' field
func (st *ConfigState) GetHTTPServerMaxReceiveBufferPerConnection() (v bytesize.Size) {
	return st.config.HTTPServer.MaxReceiveBufferPerConnection
}

// SetHTTPServerMaxReceiveBufferPerConnection safely sets the Configuration value for state's 'HTTPServer.MaxReceiveBufferPerConnection' field
func (st *ConfigState) SetHTTPServerMaxReceiveBufferPerConnection(v bytesize.Size) {
	st.config.HTTPServer.MaxReceiveBufferPerConnection = v
	st.reloadToViper()
}

// GetHTTPServerMaxReceiveBufferPerConnection safely fetches the value for global configuration 'HTTPServer.MaxReceiveBufferPerConnection' field
func GetHTTPServerMaxReceiveBufferPerConnection() bytesize.Size {
	return global.GetHTTPServerMaxReceiveBufferPerConnection()
}

// SetHTTPServerMaxReceiveBufferPerConnection safely sets the value for global configuration 'HTTPServer.MaxReceiveBufferPerConnection' field
func SetHTTPServerMaxReceiveBufferPerConnection(v bytesize.Size) {
	global.SetHTTPServerMaxReceiveBufferPerConnection(v)
}

// GetHTTPServerMaxReceiveBufferPerStream safely fetches the Configuration value for state's 'HTTPServer.MaxReceiveBufferPerStream' field
func (st *ConfigState) GetHTTPServerMaxReceiveBufferPerStream() (v bytesize.Size) {
	return st.config.HTTPServer.MaxReceiveBufferPerStream
}

// SetHTTPServerMaxReceiveBufferPerStream safely sets the Configuration value for state's 'HTTPServer.MaxReceiveBufferPerStream' field
func (st *ConfigState) SetHTTPServerMaxReceiveBufferPerStream(v bytesize.Size) {
	st.config.HTTPServer.MaxReceiveBufferPerStream = v
	st.reloadToViper()
}

// GetHTTPServerMaxReceiveBufferPerStream safely fetches the value for global configuration 'HTTPServer.MaxReceiveBufferPerStream' field
func GetHTTPServerMaxReceiveBufferPerStream() bytesize.Size {
	return global.GetHTTPServerMaxReceiveBufferPerStream()
}

// SetHTTPServerMaxReceiveBufferPerStream safely sets the value for global configuration 'HTTPServer.MaxReceiveBufferPerStream' field
func SetHTTPServerMaxReceiveBufferPerStream(v bytesize.Size) {
	global.SetHTTPServerMaxReceiveBufferPerStream(v)
}

// GetHTTPServerSendPingTimeout safely fetches the Configuration value for state's 'HTTPServer.SendPingTimeout' field
func (st *ConfigState) GetHTTPServerSendPingTimeout() (v time.Duration) {
	return st.config.HTTPServer.SendPingTimeout
}

// SetHTTPServerSendPingTimeout safely sets the Configuration value for state's 'HTTPServer.SendPingTimeout' field
func (st *ConfigState) SetHTTPServerSendPingTimeout(v time.Duration) {
	st.config.HTTPServer.SendPingTimeout = v
	st.reloadToViper()
}

// GetHTTPServerSendPingTimeout safely fetches the value for global configuration 'HTTPServer.SendPingTimeout' field
func GetHTTPServerSendPingTimeout() time.Duration { return global.GetHTTPServerSendPingTimeout() }

// SetHTTPServerSendPingTimeout safely sets the value for global configuration 'HTTPServer.SendPingTimeout' field
func SetHTTPServerSendPingTimeout(v time.Duration) { global.SetHTTPServerSendPingTimeout(v) }

// GetHTTPServerPingTimeout safely fetches the Configuration value for state's 'HTTPServer.PingTimeout' field
func (st *ConfigState) GetHTTPServerPingTimeout() (v time.Duration) {
	return st.config.HTTPServer.PingTimeout
}

// SetHTTPServerPingTimeout safely sets the Configuration value for state's 'HTTPServer.PingTimeout' field
func (st *ConfigState) SetHTTPServerPingTimeout(v time.Duration) {
	st.config.HTTPServer.PingTimeout = v
	st.reloadToViper()
}

// GetHTTPServerPingTimeout safely fetches the value for global configuration 'HTTPServer.PingTimeout' field
func GetHTTPServerPingTimeout() time.Duration { return global.GetHTTPServerPingTimeout() }

// SetHTTPServerPingTimeout safely sets the value for global configuration 'HTTPServer.PingTimeout' field
func SetHTTPServerPingTimeout(v time.Duration) { global.SetHTTPServerPingTimeout(v) }

// GetHTTPServerWriteByteTimeout safely fetches the Configuration value for state's 'HTTPServer.WriteByteTimeout' field
func (st *ConfigState) GetHTTPServerWriteByteTimeout() (v time.Duration) {
	return st.config.HTTPServer.WriteByteTimeout
}

// SetHTTPServerWriteByteTimeout safely sets the Configuration value for state's 'HTTPServer.WriteByteTimeout' field
func (st *ConfigState) SetHTTPServerWriteByteTimeout(v time.Duration) {
	st.config.HTTPServer.WriteByteTimeout = v
	st.reloadToViper()
}

// GetHTTPServerWriteByteTimeout safely fetches the value for global configuration 'HTTPServer.WriteByteTimeout' field
func GetHTTPServerWriteByteTimeout() time.Duration { return global.GetHTTPServerWriteByteTimeout() }

// SetHTTPServerWriteByteTimeout safely sets the value for global configuration 'HTTPServer.WriteByteTimeout' field
func SetHTTPServerWriteByteTimeout(v time.Duration) { global.SetHTTPServerWriteByteTimeout(v) }

// GetHTTPClientAllowIPs safely fetches the Configuration value for state's 'HTTPClient.AllowIPs' field
func (st *ConfigState) GetHTTPClientAllowIPs() (v IPPrefixes) {
	return st.config.HTTPClient.AllowIPs
}

// SetHTTPClientAllowIPs safely sets the Configuration value for state's 'HTTPClient.AllowIPs' field
func (st *ConfigState) SetHTTPClientAllowIPs(v IPPrefixes) {
	st.config.HTTPClient.AllowIPs = v
	st.reloadToViper()
}

// GetHTTPClientAllowIPs safely fetches the value for global configuration 'HTTPClient.AllowIPs' field
func GetHTTPClientAllowIPs() IPPrefixes { return global.GetHTTPClientAllowIPs() }

// SetHTTPClientAllowIPs safely sets the value for global configuration 'HTTPClient.AllowIPs' field
func SetHTTPClientAllowIPs(v IPPrefixes) { global.SetHTTPClientAllowIPs(v) }

// GetHTTPClientBlockIPs safely fetches the Configuration value for state's 'HTTPClient.BlockIPs' field
func (st *ConfigState) GetHTTPClientBlockIPs() (v IPPrefixes) {
	return st.config.HTTPClient.BlockIPs
}

// SetHTTPClientBlockIPs safely sets the Configuration value for state's 'HTTPClient.BlockIPs' field
func (st *ConfigState) SetHTTPClientBlockIPs(v IPPrefixes) {
	st.config.HTTPClient.BlockIPs = v
	st.reloadToViper()
}

// GetHTTPClientBlockIPs safely fetches the value for global configuration 'HTTPClient.BlockIPs' field
func GetHTTPClientBlockIPs() IPPrefixes { return global.GetHTTPClientBlockIPs() }

// SetHTTPClientBlockIPs safely sets the value for global configuration 'HTTPClient.BlockIPs' field
func SetHTTPClientBlockIPs(v IPPrefixes) { global.SetHTTPClientBlockIPs(v) }

// GetHTTPClientTimeout safely fetches the Configuration value for state's 'HTTPClient.Timeout' field
func (st *ConfigState) GetHTTPClientTimeout() (v time.Duration) {
	return st.config.HTTPClient.Timeout
}

// SetHTTPClientTimeout safely sets the Configuration value for state's 'HTTPClient.Timeout' field
func (st *ConfigState) SetHTTPClientTimeout(v time.Duration) {
	st.config.HTTPClient.Timeout = v
	st.reloadToViper()
}

// GetHTTPClientTimeout safely fetches the value for global configuration 'HTTPClient.Timeout' field
func GetHTTPClientTimeout() time.Duration { return global.GetHTTPClientTimeout() }

// SetHTTPClientTimeout safely sets the value for global configuration 'HTTPClient.Timeout' field
func SetHTTPClientTimeout(v time.Duration) { global.SetHTTPClientTimeout(v) }

// GetHTTPClientTLSInsecureSkipVerify safely fetches the Configuration value for state's 'HTTPClient.TLSInsecureSkipVerify' field
func (st *ConfigState) GetHTTPClientTLSInsecureSkipVerify() (v bool) {
	return st.config.HTTPClient.TLSInsecureSkipVerify
}

// SetHTTPClientTLSInsecureSkipVerify safely sets the Configuration value for state's 'HTTPClient.TLSInsecureSkipVerify' field
func (st *ConfigState) SetHTTPClientTLSInsecureSkipVerify(v bool) {
	st.config.HTTPClient.TLSInsecureSkipVerify = v
	st.reloadToViper()
}

// GetHTTPClientTLSInsecureSkipVerify safely fetches the value for global configuration 'HTTPClient.TLSInsecureSkipVerify' field
func GetHTTPClientTLSInsecureSkipVerify() bool { return global.GetHTTPClientTLSInsecureSkipVerify() }

// SetHTTPClientTLSInsecureSkipVerify safely sets the value for global configuration 'HTTPClient.TLSInsecureSkipVerify' field
func SetHTTPClientTLSInsecureSkipVerify(v bool) { global.SetHTTPClientTLSInsecureSkipVerify(v) }

// GetHTTPClientInsecureOutgoing safely fetches the Configuration value for state's 'HTTPClient.InsecureOutgoing' field
func (st *ConfigState) GetHTTPClientInsecureOutgoing() (v bool) {
	return st.config.HTTPClient.InsecureOutgoing
}

// SetHTTPClientInsecureOutgoing safely sets the Configuration value for state's 'HTTPClient.InsecureOutgoing' field
func (st *ConfigState) SetHTTPClientInsecureOutgoing(v bool) {
	st.config.HTTPClient.InsecureOutgoing = v
	st.reloadToViper()
}

// GetHTTPClientInsecureOutgoing safely fetches the value for global configuration 'HTTPClient.InsecureOutgoing' field
func GetHTTPClientInsecureOutgoing() bool { return global.GetHTTPClientInsecureOutgoing() }

// SetHTTPClientInsecureOutgoing safely sets the value for global configuration 'HTTPClient.InsecureOutgoing' field
func SetHTTPClientInsecureOutgoing(v bool) { global.SetHTTPClientInsecureOutgoing(v) }

// GetHTTPClientDisableKeepAlives safely fetches the Configuration value for state's 'HTTPClient.DisableKeepAlives' field
func (st *ConfigState) GetHTTPClientDisableKeepAlives() (v bool) {
	return st.config.HTTPClient.DisableKeepAlives
}

// SetHTTPClientDisableKeepAlives safely sets the Configuration value for state's 'HTTPClient.DisableKeepAlives' field
func (st *ConfigState) SetHTTPClientDisableKeepAlives(v bool) {
	st.config.HTTPClient.DisableKeepAlives = v
	st.reloadToViper()
}

// GetHTTPClientDisableKeepAlives safely fetches the value for global configuration 'HTTPClient.DisableKeepAlives' field
func GetHTTPClientDisableKeepAlives() bool { return global.GetHTTPClientDisableKeepAlives() }

// SetHTTPClientDisableKeepAlives safely sets the value for global configuration 'HTTPClient.DisableKeepAlives' field
func SetHTTPClientDisableKeepAlives(v bool) { global.SetHTTPClientDisableKeepAlives(v) }

// GetHTTPClientMaxIdleConns safely fetches the Configuration value for state's 'HTTPClient.MaxIdleConns' field
func (st *ConfigState) GetHTTPClientMaxIdleConns() (v int) {
	return st.config.HTTPClient.MaxIdleConns
}

// SetHTTPClientMaxIdleConns safely sets the Configuration value for state's 'HTTPClient.MaxIdleConns' field
func (st *ConfigState) SetHTTPClientMaxIdleConns(v int) {
	st.config.HTTPClient.MaxIdleConns = v
	st.reloadToViper()
}

// GetHTTPClientMaxIdleConns safely fetches the value for global configuration 'HTTPClient.MaxIdleConns' field
func GetHTTPClientMaxIdleConns() int { return global.GetHTTPClientMaxIdleConns() }

// SetHTTPClientMaxIdleConns safely sets the value for global configuration 'HTTPClient.MaxIdleConns' field
func SetHTTPClientMaxIdleConns(v int) { global.SetHTTPClientMaxIdleConns(v) }

// GetHTTPClientMaxIdleConnsPerHost safely fetches the Configuration value for state's 'HTTPClient.MaxIdleConnsPerHost' field
func (st *ConfigState) GetHTTPClientMaxIdleConnsPerHost() (v int) {
	return st.config.HTTPClient.MaxIdleConnsPerHost
}

// SetHTTPClientMaxIdleConnsPerHost safely sets the Configuration value for state's 'HTTPClient.MaxIdleConnsPerHost' field
func (st *ConfigState) SetHTTPClientMaxIdleConnsPerHost(v int) {
	st.config.HTTPClient.MaxIdleConnsPerHost = v
	st.reloadToViper()
}

// GetHTTPClientMaxIdleConnsPerHost safely fetches the value for global configuration 'HTTPClient.MaxIdleConnsPerHost' field
func GetHTTPClientMaxIdleConnsPerHost() int { return global.GetHTTPClientMaxIdleConnsPerHost() }

// SetHTTPClientMaxIdleConnsPerHost safely sets the value for global configuration 'HTTPClient.MaxIdleConnsPerHost' field
func SetHTTPClientMaxIdleConnsPerHost(v int) { global.SetHTTPClientMaxIdleConnsPerHost(v) }

// GetHTTPClientMaxConnsPerHost safely fetches the Configuration value for state's 'HTTPClient.MaxConnsPerHost' field
func (st *ConfigState) GetHTTPClientMaxConnsPerHost() (v int) {
	return st.config.HTTPClient.MaxConnsPerHost
}

// SetHTTPClientMaxConnsPerHost safely sets the Configuration value for state's 'HTTPClient.MaxConnsPerHost' field
func (st *ConfigState) SetHTTPClientMaxConnsPerHost(v int) {
	st.config.HTTPClient.MaxConnsPerHost = v
	st.reloadToViper()
}

// GetHTTPClientMaxConnsPerHost safely fetches the value for global configuration 'HTTPClient.MaxConnsPerHost' field
func GetHTTPClientMaxConnsPerHost() int { return global.GetHTTPClientMaxConnsPerHost() }

// SetHTTPClientMaxConnsPerHost safely sets the value for global configuration 'HTTPClient.MaxConnsPerHost' field
func SetHTTPClientMaxConnsPerHost(v int) { global.SetHTTPClientMaxConnsPerHost(v) }

// GetHTTPClientIdleConnTimeout safely fetches the Configuration value for state's 'HTTPClient.IdleConnTimeout' field
func (st *ConfigState) GetHTTPClientIdleConnTimeout() (v time.Duration) {
	return st.config.HTTPClient.IdleConnTimeout
}

// SetHTTPClientIdleConnTimeout safely sets the Configuration value for state's 'HTTPClient.IdleConnTimeout' field
func (st *ConfigState) SetHTTPClientIdleConnTimeout(v time.Duration) {
	st.config.HTTPClient.IdleConnTimeout = v
	st.reloadToViper()
}

// GetHTTPClientIdleConnTimeout safely fetches the value for global configuration 'HTTPClient.IdleConnTimeout' field
func GetHTTPClientIdleConnTimeout() time.Duration { return global.GetHTTPClientIdleConnTimeout() }

// SetHTTPClientIdleConnTimeout safely sets the value for global configuration 'HTTPClient.IdleConnTimeout' field
func SetHTTPClientIdleConnTimeout(v time.Duration) { global.SetHTTPClientIdleConnTimeout(v) }

// GetHTTPClientTLSHandshakeTimeout safely fetches the Configuration value for state's 'HTTPClient.TLSHandshakeTimeout' field
func (st *ConfigState) GetHTTPClientTLSHandshakeTimeout() (v time.Duration) {
	return st.config.HTTPClient.TLSHandshakeTimeout
}

// SetHTTPClientTLSHandshakeTimeout safely sets the Configuration value for state's 'HTTPClient.TLSHandshakeTimeout' field
func (st *ConfigState) SetHTTPClientTLSHandshakeTimeout(v time.Duration) {
	st.config.HTTPClient.TLSHandshakeTimeout = v
	st.reloadToViper()
}

// GetHTTPClientTLSHandshakeTimeout safely fetches the value for global configuration 'HTTPClient.TLSHandshakeTimeout' field
func GetHTTPClientTLSHandshakeTimeout() time.Duration {
	return global.GetHTTPClientTLSHandshakeTimeout()
}

// SetHTTPClientTLSHandshakeTimeout safely sets the value for global configuration 'HTTPClient.TLSHandshakeTimeout' field
func SetHTTPClientTLSHandshakeTimeout(v time.Duration) { global.SetHTTPClientTLSHandshakeTimeout(v) }

// GetHTTPClientResponseHeaderTimeout safely fetches the Configuration value for state's 'HTTPClient.ResponseHeaderTimeout' field
func (st *ConfigState) GetHTTPClientResponseHeaderTimeout() (v time.Duration) {
	return st.config.HTTPClient.ResponseHeaderTimeout
}

// SetHTTPClientResponseHeaderTimeout safely sets the Configuration value for state's 'HTTPClient.ResponseHeaderTimeout' field
func (st *ConfigState) SetHTTPClientResponseHeaderTimeout(v time.Duration) {
	st.config.HTTPClient.ResponseHeaderTimeout = v
	st.reloadToViper()
}

// GetHTTPClientResponseHeaderTimeout safely fetches the value for global configuration 'HTTPClient.ResponseHeaderTimeout' field
func GetHTTPClientResponseHeaderTimeout() time.Duration {
	return global.GetHTTPClientResponseHeaderTimeout()
}

// SetHTTPClientResponseHeaderTimeout safely sets the value for global configuration 'HTTPClient.ResponseHeaderTimeout' field
func SetHTTPClientResponseHeaderTimeout(v time.Duration) {
	global.SetHTTPClientResponseHeaderTimeout(v)
}

// GetHTTPClientReadBufferSize safely fetches the Configuration value for state's 'HTTPClient.ReadBufferSize' field
func (st *ConfigState) GetHTTPClientReadBufferSize() (v bytesize.Size) {
	return st.config.HTTPClient.ReadBufferSize
}

// SetHTTPClientReadBufferSize safely sets the Configuration value for state's 'HTTPClient.ReadBufferSize' field
func (st *ConfigState) SetHTTPClientReadBufferSize(v bytesize.Size) {
	st.config.HTTPClient.ReadBufferSize = v
	st.reloadToViper()
}

// GetHTTPClientReadBufferSize safely fetches the value for global configuration 'HTTPClient.ReadBufferSize' field
func GetHTTPClientReadBufferSize() bytesize.Size { return global.GetHTTPClientReadBufferSize() }

// SetHTTPClientReadBufferSize safely sets the value for global configuration 'HTTPClient.ReadBufferSize' field
func SetHTTPClientReadBufferSize(v bytesize.Size) { global.SetHTTPClientReadBufferSize(v) }

// GetHTTPClientWriteBufferSize safely fetches the Configuration value for state's 'HTTPClient.WriteBufferSize' field
func (st *ConfigState) GetHTTPClientWriteBufferSize() (v bytesize.Size) {
	return st.config.HTTPClient.WriteBufferSize
}

// SetHTTPClientWriteBufferSize safely sets the Configuration value for state's 'HTTPClient.WriteBufferSize' field
func (st *ConfigState) SetHTTPClientWriteBufferSize(v bytesize.Size) {
	st.config.HTTPClient.WriteBufferSize = v
	st.reloadToViper()
}

// GetHTTPClientWriteBufferSize safely fetches the value for global configuration 'HTTPClient.WriteBufferSize' field
func GetHTTPClientWriteBufferSize() bytesize.Size { return global.GetHTTPClientWriteBufferSize() }

// SetHTTPClientWriteBufferSize safely sets the value for global configuration 'HTTPClient.WriteBufferSize' field
func SetHTTPClientWriteBufferSize(v bytesize.Size) { global.SetHTTPClientWriteBufferSize(v) }

// GetMediaDescriptionMinChars safely fetches the Configuration value for state's 'Media.DescriptionMinChars' field
func (st *ConfigState) GetMediaDescriptionMinChars() (v int) {
	return st.config.Media.DescriptionMinChars
}

// SetMediaDescriptionMinChars safely sets the Configuration value for state's 'Media.DescriptionMinChars' field
func (st *ConfigState) SetMediaDescriptionMinChars(v int) {
	st.config.Media.DescriptionMinChars = v
	st.reloadToViper()
}

// GetMediaDescriptionMinChars safely fetches the value for global configuration 'Media.DescriptionMinChars' field
func GetMediaDescriptionMinChars() int { return global.GetMediaDescriptionMinChars() }

// SetMediaDescriptionMinChars safely sets the value for global configuration 'Media.DescriptionMinChars' field
func SetMediaDescriptionMinChars(v int) { global.SetMediaDescriptionMinChars(v) }

// GetMediaDescriptionMaxChars safely fetches the Configuration value for state's 'Media.DescriptionMaxChars' field
func (st *ConfigState) GetMediaDescriptionMaxChars() (v int) {
	return st.config.Media.DescriptionMaxChars
}

// SetMediaDescriptionMaxChars safely sets the Configuration value for state's 'Media.DescriptionMaxChars' field
func (st *ConfigState) SetMediaDescriptionMaxChars(v int) {
	st.config.Media.DescriptionMaxChars = v
	st.reloadToViper()
}

// GetMediaDescriptionMaxChars safely fetches the value for global configuration 'Media.DescriptionMaxChars' field
func GetMediaDescriptionMaxChars() int { return global.GetMediaDescriptionMaxChars() }

// SetMediaDescriptionMaxChars safely sets the value for global configuration 'Media.DescriptionMaxChars' field
func SetMediaDescriptionMaxChars(v int) { global.SetMediaDescriptionMaxChars(v) }

// GetMediaEmojiLocalMaxSize safely fetches the Configuration value for state's 'Media.EmojiLocalMaxSize' field
func (st *ConfigState) GetMediaEmojiLocalMaxSize() (v bytesize.Size) {
	return st.config.Media.EmojiLocalMaxSize
}

// SetMediaEmojiLocalMaxSize safely sets the Configuration value for state's 'Media.EmojiLocalMaxSize' field
func (st *ConfigState) SetMediaEmojiLocalMaxSize(v bytesize.Size) {
	st.config.Media.EmojiLocalMaxSize = v
	st.reloadToViper()
}

// GetMediaEmojiLocalMaxSize safely fetches the value for global configuration 'Media.EmojiLocalMaxSize' field
func GetMediaEmojiLocalMaxSize() bytesize.Size { return global.GetMediaEmojiLocalMaxSize() }

// SetMediaEmojiLocalMaxSize safely sets the value for global configuration 'Media.EmojiLocalMaxSize' field
func SetMediaEmojiLocalMaxSize(v bytesize.Size) { global.SetMediaEmojiLocalMaxSize(v) }

// GetMediaEmojiRemoteMaxSize safely fetches the Configuration value for state's 'Media.EmojiRemoteMaxSize' field
func (st *ConfigState) GetMediaEmojiRemoteMaxSize() (v bytesize.Size) {
	return st.config.Media.EmojiRemoteMaxSize
}

// SetMediaEmojiRemoteMaxSize safely sets the Configuration value for state's 'Media.EmojiRemoteMaxSize' field
func (st *ConfigState) SetMediaEmojiRemoteMaxSize(v bytesize.Size) {
	st.config.Media.EmojiRemoteMaxSize = v
	st.reloadToViper()
}

// GetMediaEmojiRemoteMaxSize safely fetches the value for global configuration 'Media.EmojiRemoteMaxSize' field
func GetMediaEmojiRemoteMaxSize() bytesize.Size { return global.GetMediaEmojiRemoteMaxSize() }

// SetMediaEmojiRemoteMaxSize safely sets the value for global configuration 'Media.EmojiRemoteMaxSize' field
func SetMediaEmojiRemoteMaxSize(v bytesize.Size) { global.SetMediaEmojiRemoteMaxSize(v) }

// GetMediaImageSizeHint safely fetches the Configuration value for state's 'Media.ImageSizeHint' field
func (st *ConfigState) GetMediaImageSizeHint() (v bytesize.Size) {
	return st.config.Media.ImageSizeHint
}

// SetMediaImageSizeHint safely sets the Configuration value for state's 'Media.ImageSizeHint' field
func (st *ConfigState) SetMediaImageSizeHint(v bytesize.Size) {
	st.config.Media.ImageSizeHint = v
	st.reloadToViper()
}

// GetMediaImageSizeHint safely fetches the value for global configuration 'Media.ImageSizeHint' field
func GetMediaImageSizeHint() bytesize.Size { return global.GetMediaImageSizeHint() }

// SetMediaImageSizeHint safely sets the value for global configuration 'Media.ImageSizeHint' field
func SetMediaImageSizeHint(v bytesize.Size) { global.SetMediaImageSizeHint(v) }

// GetMediaVideoSizeHint safely fetches the Configuration value for state's 'Media.VideoSizeHint' field
func (st *ConfigState) GetMediaVideoSizeHint() (v bytesize.Size) {
	return st.config.Media.VideoSizeHint
}

// SetMediaVideoSizeHint safely sets the Configuration value for state's 'Media.VideoSizeHint' field
func (st *ConfigState) SetMediaVideoSizeHint(v bytesize.Size) {
	st.config.Media.VideoSizeHint = v
	st.reloadToViper()
}

// GetMediaVideoSizeHint safely fetches the value for global configuration 'Media.VideoSizeHint' field
func GetMediaVideoSizeHint() bytesize.Size { return global.GetMediaVideoSizeHint() }

// SetMediaVideoSizeHint safely sets the value for global configuration 'Media.VideoSizeHint' field
func SetMediaVideoSizeHint(v bytesize.Size) { global.SetMediaVideoSizeHint(v) }

// GetMediaLocalMaxSize safely fetches the Configuration value for state's 'Media.LocalMaxSize' field
func (st *ConfigState) GetMediaLocalMaxSize() (v bytesize.Size) {
	return st.config.Media.LocalMaxSize
}

// SetMediaLocalMaxSize safely sets the Configuration value for state's 'Media.LocalMaxSize' field
func (st *ConfigState) SetMediaLocalMaxSize(v bytesize.Size) {
	st.config.Media.LocalMaxSize = v
	st.reloadToViper()
}

// GetMediaLocalMaxSize safely fetches the value for global configuration 'Media.LocalMaxSize' field
func GetMediaLocalMaxSize() bytesize.Size { return global.GetMediaLocalMaxSize() }

// SetMediaLocalMaxSize safely sets the value for global configuration 'Media.LocalMaxSize' field
func SetMediaLocalMaxSize(v bytesize.Size) { global.SetMediaLocalMaxSize(v) }

// GetMediaRemoteMaxSize safely fetches the Configuration value for state's 'Media.RemoteMaxSize' field
func (st *ConfigState) GetMediaRemoteMaxSize() (v bytesize.Size) {
	return st.config.Media.RemoteMaxSize
}

// SetMediaRemoteMaxSize safely sets the Configuration value for state's 'Media.RemoteMaxSize' field
func (st *ConfigState) SetMediaRemoteMaxSize(v bytesize.Size) {
	st.config.Media.RemoteMaxSize = v
	st.reloadToViper()
}

// GetMediaRemoteMaxSize safely fetches the value for global configuration 'Media.RemoteMaxSize' field
func GetMediaRemoteMaxSize() bytesize.Size { return global.GetMediaRemoteMaxSize() }

// SetMediaRemoteMaxSize safely sets the value for global configuration 'Media.RemoteMaxSize' field
func SetMediaRemoteMaxSize(v bytesize.Size) { global.SetMediaRemoteMaxSize(v) }

// GetMediaFfmpegPoolSize safely fetches the Configuration value for state's 'Media.FfmpegPoolSize' field
func (st *ConfigState) GetMediaFfmpegPoolSize() (v int) {
	return st.config.Media.FfmpegPoolSize
}

// SetMediaFfmpegPoolSize safely sets the Configuration value for state's 'Media.FfmpegPoolSize' field
func (st *ConfigState) SetMediaFfmpegPoolSize(v int) {
	st.config.Media.FfmpegPoolSize = v
	st.reloadToViper()
}

// GetMediaFfmpegPoolSize safely fetches the value for global configuration 'Media.FfmpegPoolSize' field
func GetMediaFfmpegPoolSize() int { return global.GetMediaFfmpegPoolSize() }

// SetMediaFfmpegPoolSize safely sets the value for global configuration 'Media.FfmpegPoolSize' field
func SetMediaFfmpegPoolSize(v int) { global.SetMediaFfmpegPoolSize(v) }

// GetMediaThumbMaxPixels safely fetches the Configuration value for state's 'Media.ThumbMaxPixels' field
func (st *ConfigState) GetMediaThumbMaxPixels() (v int) {
	return st.config.Media.ThumbMaxPixels
}

// SetMediaThumbMaxPixels safely sets the Configuration value for state's 'Media.ThumbMaxPixels' field
func (st *ConfigState) SetMediaThumbMaxPixels(v int) {
	st.config.Media.ThumbMaxPixels = v
	st.reloadToViper()
}

// GetMediaThumbMaxPixels safely fetches the value for global configuration 'Media.ThumbMaxPixels' field
func GetMediaThumbMaxPixels() int { return global.GetMediaThumbMaxPixels() }

// SetMediaThumbMaxPixels safely sets the value for global configuration 'Media.ThumbMaxPixels' field
func SetMediaThumbMaxPixels(v int) { global.SetMediaThumbMaxPixels(v) }

// GetCacheS3ObjectInfo safely fetches the Configuration value for state's 'Cache.S3ObjectInfo' field
func (st *ConfigState) GetCacheS3ObjectInfo() (v uint32) {
	return st.config.Cache.S3ObjectInfo
}

// SetCacheS3ObjectInfo safely sets the Configuration value for state's 'Cache.S3ObjectInfo' field
func (st *ConfigState) SetCacheS3ObjectInfo(v uint32) {
	st.config.Cache.S3ObjectInfo = v
	st.reloadToViper()
}

// GetCacheS3ObjectInfo safely fetches the value for global configuration 'Cache.S3ObjectInfo' field
func GetCacheS3ObjectInfo() uint32 { return global.GetCacheS3ObjectInfo() }

// SetCacheS3ObjectInfo safely sets the value for global configuration 'Cache.S3ObjectInfo' field
func SetCacheS3ObjectInfo(v uint32) { global.SetCacheS3ObjectInfo(v) }

// GetCacheHomeTimelineSize safely fetches the Configuration value for state's 'Cache.HomeTimelineSize' field
func (st *ConfigState) GetCacheHomeTimelineSize() (v uint32) {
	return st.config.Cache.HomeTimelineSize
}

// SetCacheHomeTimelineSize safely sets the Configuration value for state's 'Cache.HomeTimelineSize' field
func (st *ConfigState) SetCacheHomeTimelineSize(v uint32) {
	st.config.Cache.HomeTimelineSize = v
	st.reloadToViper()
}

// GetCacheHomeTimelineSize safely fetches the value for global configuration 'Cache.HomeTimelineSize' field
func GetCacheHomeTimelineSize() uint32 { return global.GetCacheHomeTimelineSize() }

// SetCacheHomeTimelineSize safely sets the value for global configuration 'Cache.HomeTimelineSize' field
func SetCacheHomeTimelineSize(v uint32) { global.SetCacheHomeTimelineSize(v) }

// GetCacheListTimelineSize safely fetches the Configuration value for state's 'Cache.ListTimelineSize' field
func (st *ConfigState) GetCacheListTimelineSize() (v uint32) {
	return st.config.Cache.ListTimelineSize
}

// SetCacheListTimelineSize safely sets the Configuration value for state's 'Cache.ListTimelineSize' field
func (st *ConfigState) SetCacheListTimelineSize(v uint32) {
	st.config.Cache.ListTimelineSize = v
	st.reloadToViper()
}

// GetCacheListTimelineSize safely fetches the value for global configuration 'Cache.ListTimelineSize' field
func GetCacheListTimelineSize() uint32 { return global.GetCacheListTimelineSize() }

// SetCacheListTimelineSize safely sets the value for global configuration 'Cache.ListTimelineSize' field
func SetCacheListTimelineSize(v uint32) { global.SetCacheListTimelineSize(v) }

// GetCacheTagTimelineSize safely fetches the Configuration value for state's 'Cache.TagTimelineSize' field
func (st *ConfigState) GetCacheTagTimelineSize() (v uint32) {
	return st.config.Cache.TagTimelineSize
}

// SetCacheTagTimelineSize safely sets the Configuration value for state's 'Cache.TagTimelineSize' field
func (st *ConfigState) SetCacheTagTimelineSize(v uint32) {
	st.config.Cache.TagTimelineSize = v
	st.reloadToViper()
}

// GetCacheTagTimelineSize safely fetches the value for global configuration 'Cache.TagTimelineSize' field
func GetCacheTagTimelineSize() uint32 { return global.GetCacheTagTimelineSize() }

// SetCacheTagTimelineSize safely sets the value for global configuration 'Cache.TagTimelineSize' field
func SetCacheTagTimelineSize(v uint32) { global.SetCacheTagTimelineSize(v) }

// GetCacheNotificationTimelineSize safely fetches the Configuration value for state's 'Cache.NotificationTimelineSize' field
func (st *ConfigState) GetCacheNotificationTimelineSize() (v uint32) {
	return st.config.Cache.NotificationTimelineSize
}

// SetCacheNotificationTimelineSize safely sets the Configuration value for state's 'Cache.NotificationTimelineSize' field
func (st *ConfigState) SetCacheNotificationTimelineSize(v uint32) {
	st.config.Cache.NotificationTimelineSize = v
	st.reloadToViper()
}

// GetCacheNotificationTimelineSize safely fetches the value for global configuration 'Cache.NotificationTimelineSize' field
func GetCacheNotificationTimelineSize() uint32 { return global.GetCacheNotificationTimelineSize() }

// SetCacheNotificationTimelineSize safely sets the value for global configuration 'Cache.NotificationTimelineSize' field
func SetCacheNotificationTimelineSize(v uint32) { global.SetCacheNotificationTimelineSize(v) }

// GetCacheHomeTimelineTimeout safely fetches the Configuration value for state's 'Cache.HomeTimelineTimeout' field
func (st *ConfigState) GetCacheHomeTimelineTimeout() (v time.Duration) {
	return st.config.Cache.HomeTimelineTimeout
}

// SetCacheHomeTimelineTimeout safely sets the Configuration value for state's 'Cache.HomeTimelineTimeout' field
func (st *ConfigState) SetCacheHomeTimelineTimeout(v time.Duration) {
	st.config.Cache.HomeTimelineTimeout = v
	st.reloadToViper()
}

// GetCacheHomeTimelineTimeout safely fetches the value for global configuration 'Cache.HomeTimelineTimeout' field
func GetCacheHomeTimelineTimeout() time.Duration { return global.GetCacheHomeTimelineTimeout() }

// SetCacheHomeTimelineTimeout safely sets the value for global configuration 'Cache.HomeTimelineTimeout' field
func SetCacheHomeTimelineTimeout(v time.Duration) { global.SetCacheHomeTimelineTimeout(v) }

// GetCacheListTimelineTimeout safely fetches the Configuration value for state's 'Cache.ListTimelineTimeout' field
func (st *ConfigState) GetCacheListTimelineTimeout() (v time.Duration) {
	return st.config.Cache.ListTimelineTimeout
}

// SetCacheListTimelineTimeout safely sets the Configuration value for state's 'Cache.ListTimelineTimeout' field
func (st *ConfigState) SetCacheListTimelineTimeout(v time.Duration) {
	st.config.Cache.ListTimelineTimeout = v
	st.reloadToViper()
}

// GetCacheListTimelineTimeout safely fetches the value for global configuration 'Cache.ListTimelineTimeout' field
func GetCacheListTimelineTimeout() time.Duration { return global.GetCacheListTimelineTimeout() }

// SetCacheListTimelineTimeout safely sets the value for global configuration 'Cache.ListTimelineTimeout' field
func SetCacheListTimelineTimeout(v time.Duration) { global.SetCacheListTimelineTimeout(v) }

// GetCacheTagTimelineTimeout safely fetches the Configuration value for state's 'Cache.TagTimelineTimeout' field
func (st *ConfigState) GetCacheTagTimelineTimeout() (v time.Duration) {
	return st.config.Cache.TagTimelineTimeout
}

// SetCacheTagTimelineTimeout safely sets the Configuration value for state's 'Cache.TagTimelineTimeout' field
func (st *ConfigState) SetCacheTagTimelineTimeout(v time.Duration) {
	st.config.Cache.TagTimelineTimeout = v
	st.reloadToViper()
}

// GetCacheTagTimelineTimeout safely fetches the value for global configuration 'Cache.TagTimelineTimeout' field
func GetCacheTagTimelineTimeout() time.Duration { return global.GetCacheTagTimelineTimeout() }

// SetCacheTagTimelineTimeout safely sets the value for global configuration 'Cache.TagTimelineTimeout' field
func SetCacheTagTimelineTimeout(v time.Duration) { global.SetCacheTagTimelineTimeout(v) }

// GetCacheNotificationTimelineTimeout safely fetches the Configuration value for state's 'Cache.NotificationTimelineTimeout' field
func (st *ConfigState) GetCacheNotificationTimelineTimeout() (v time.Duration) {
	return st.config.Cache.NotificationTimelineTimeout
}

// SetCacheNotificationTimelineTimeout safely sets the Configuration value for state's 'Cache.NotificationTimelineTimeout' field
func (st *ConfigState) SetCacheNotificationTimelineTimeout(v time.Duration) {
	st.config.Cache.NotificationTimelineTimeout = v
	st.reloadToViper()
}

// GetCacheNotificationTimelineTimeout safely fetches the value for global configuration 'Cache.NotificationTimelineTimeout' field
func GetCacheNotificationTimelineTimeout() time.Duration {
	return global.GetCacheNotificationTimelineTimeout()
}

// SetCacheNotificationTimelineTimeout safely sets the value for global configuration 'Cache.NotificationTimelineTimeout' field
func SetCacheNotificationTimelineTimeout(v time.Duration) {
	global.SetCacheNotificationTimelineTimeout(v)
}

// GetCacheMemoryTarget safely fetches the Configuration value for state's 'Cache.MemoryTarget' field
func (st *ConfigState) GetCacheMemoryTarget() (v bytesize.Size) {
	return st.config.Cache.MemoryTarget
}

// SetCacheMemoryTarget safely sets the Configuration value for state's 'Cache.MemoryTarget' field
func (st *ConfigState) SetCacheMemoryTarget(v bytesize.Size) {
	st.config.Cache.MemoryTarget = v
	st.reloadToViper()
}

// GetCacheMemoryTarget safely fetches the value for global configuration 'Cache.MemoryTarget' field
func GetCacheMemoryTarget() bytesize.Size { return global.GetCacheMemoryTarget() }

// SetCacheMemoryTarget safely sets the value for global configuration 'Cache.MemoryTarget' field
func SetCacheMemoryTarget(v bytesize.Size) { global.SetCacheMemoryTarget(v) }

// GetCacheAccountMemRatio safely fetches the Configuration value for state's 'Cache.AccountMemRatio' field
func (st *ConfigState) GetCacheAccountMemRatio() (v float64) {
	return st.config.Cache.AccountMemRatio
}

// SetCacheAccountMemRatio safely sets the Configuration value for state's 'Cache.AccountMemRatio' field
func (st *ConfigState) SetCacheAccountMemRatio(v float64) {
	st.config.Cache.AccountMemRatio = v
	st.reloadToViper()
}

// GetCacheAccountMemRatio safely fetches the value for global configuration 'Cache.AccountMemRatio' field
func GetCacheAccountMemRatio() float64 { return global.GetCacheAccountMemRatio() }

// SetCacheAccountMemRatio safely sets the value for global configuration 'Cache.AccountMemRatio' field
func SetCacheAccountMemRatio(v float64) { global.SetCacheAccountMemRatio(v) }

// GetCacheAccountNoteMemRatio safely fetches the Configuration value for state's 'Cache.AccountNoteMemRatio' field
func (st *ConfigState) GetCacheAccountNoteMemRatio() (v float64) {
	return st.config.Cache.AccountNoteMemRatio
}

// SetCacheAccountNoteMemRatio safely sets the Configuration value for state's 'Cache.AccountNoteMemRatio' field
func (st *ConfigState) SetCacheAccountNoteMemRatio(v float64) {
	st.config.Cache.AccountNoteMemRatio = v
	st.reloadToViper()
}

// GetCacheAccountNoteMemRatio safely fetches the value for global configuration 'Cache.AccountNoteMemRatio' field
func GetCacheAccountNoteMemRatio() float64 { return global.GetCacheAccountNoteMemRatio() }

// SetCacheAccountNoteMemRatio safely sets the value for global configuration 'Cache.AccountNoteMemRatio' field
func SetCacheAccountNoteMemRatio(v float64) { global.SetCacheAccountNoteMemRatio(v) }

// GetCacheAccountSettingsMemRatio safely fetches the Configuration value for state's 'Cache.AccountSettingsMemRatio' field
func (st *ConfigState) GetCacheAccountSettingsMemRatio() (v float64) {
	return st.config.Cache.AccountSettingsMemRatio
}

// SetCacheAccountSettingsMemRatio safely sets the Configuration value for state's 'Cache.AccountSettingsMemRatio' field
func (st *ConfigState) SetCacheAccountSettingsMemRatio(v float64) {
	st.config.Cache.AccountSettingsMemRatio = v
	st.reloadToViper()
}

// GetCacheAccountSettingsMemRatio safely fetches the value for global configuration 'Cache.AccountSettingsMemRatio' field
func GetCacheAccountSettingsMemRatio() float64 { return global.GetCacheAccountSettingsMemRatio() }

// SetCacheAccountSettingsMemRatio safely sets the value for global configuration 'Cache.AccountSettingsMemRatio' field
func SetCacheAccountSettingsMemRatio(v float64) { global.SetCacheAccountSettingsMemRatio(v) }

// GetCacheAccountStatsMemRatio safely fetches the Configuration value for state's 'Cache.AccountStatsMemRatio' field
func (st *ConfigState) GetCacheAccountStatsMemRatio() (v float64) {
	return st.config.Cache.AccountStatsMemRatio
}

// SetCacheAccountStatsMemRatio safely sets the Configuration value for state's 'Cache.AccountStatsMemRatio' field
func (st *ConfigState) SetCacheAccountStatsMemRatio(v float64) {
	st.config.Cache.AccountStatsMemRatio = v
	st.reloadToViper()
}

// GetCacheAccountStatsMemRatio safely fetches the value for global configuration 'Cache.AccountStatsMemRatio' field
func GetCacheAccountStatsMemRatio() float64 { return global.GetCacheAccountStatsMemRatio() }

// SetCacheAccountStatsMemRatio safely sets the value for global configuration 'Cache.AccountStatsMemRatio' field
func SetCacheAccountStatsMemRatio(v float64) { global.SetCacheAccountStatsMemRatio(v) }

// GetCacheApplicationMemRatio safely fetches the Configuration value for state's 'Cache.ApplicationMemRatio' field
func (st *ConfigState) GetCacheApplicationMemRatio() (v float64) {
	return st.config.Cache.ApplicationMemRatio
}

// SetCacheApplicationMemRatio safely sets the Configuration value for state's 'Cache.ApplicationMemRatio' field
func (st *ConfigState) SetCacheApplicationMemRatio(v float64) {
	st.config.Cache.ApplicationMemRatio = v
	st.reloadToViper()
}

// GetCacheApplicationMemRatio safely fetches the value for global configuration 'Cache.ApplicationMemRatio' field
func GetCacheApplicationMemRatio() float64 { return global.GetCacheApplicationMemRatio() }

// SetCacheApplicationMemRatio safely sets the value for global configuration 'Cache.ApplicationMemRatio' field
func SetCacheApplicationMemRatio(v float64) { global.SetCacheApplicationMemRatio(v) }

// GetCacheBlockMemRatio safely fetches the Configuration value for state's 'Cache.BlockMemRatio' field
func (st *ConfigState) GetCacheBlockMemRatio() (v float64) {
	return st.config.Cache.BlockMemRatio
}

// SetCacheBlockMemRatio safely sets the Configuration value for state's 'Cache.BlockMemRatio' field
func (st *ConfigState) SetCacheBlockMemRatio(v float64) {
	st.config.Cache.BlockMemRatio = v
	st.reloadToViper()
}

// GetCacheBlockMemRatio safely fetches the value for global configuration 'Cache.BlockMemRatio' field
func GetCacheBlockMemRatio() float64 { return global.GetCacheBlockMemRatio() }

// SetCacheBlockMemRatio safely sets the value for global configuration 'Cache.BlockMemRatio' field
func SetCacheBlockMemRatio(v float64) { global.SetCacheBlockMemRatio(v) }

// GetCacheBlockIDsMemRatio safely fetches the Configuration value for state's 'Cache.BlockIDsMemRatio' field
func (st *ConfigState) GetCacheBlockIDsMemRatio() (v float64) {
	return st.config.Cache.BlockIDsMemRatio
}

// SetCacheBlockIDsMemRatio safely sets the Configuration value for state's 'Cache.BlockIDsMemRatio' field
func (st *ConfigState) SetCacheBlockIDsMemRatio(v float64) {
	st.config.Cache.BlockIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheBlockIDsMemRatio safely fetches the value for global configuration 'Cache.BlockIDsMemRatio' field
func GetCacheBlockIDsMemRatio() float64 { return global.GetCacheBlockIDsMemRatio() }

// SetCacheBlockIDsMemRatio safely sets the value for global configuration 'Cache.BlockIDsMemRatio' field
func SetCacheBlockIDsMemRatio(v float64) { global.SetCacheBlockIDsMemRatio(v) }

// GetCacheBoostOfIDsMemRatio safely fetches the Configuration value for state's 'Cache.BoostOfIDsMemRatio' field
func (st *ConfigState) GetCacheBoostOfIDsMemRatio() (v float64) {
	return st.config.Cache.BoostOfIDsMemRatio
}

// SetCacheBoostOfIDsMemRatio safely sets the Configuration value for state's 'Cache.BoostOfIDsMemRatio' field
func (st *ConfigState) SetCacheBoostOfIDsMemRatio(v float64) {
	st.config.Cache.BoostOfIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheBoostOfIDsMemRatio safely fetches the value for global configuration 'Cache.BoostOfIDsMemRatio' field
func GetCacheBoostOfIDsMemRatio() float64 { return global.GetCacheBoostOfIDsMemRatio() }

// SetCacheBoostOfIDsMemRatio safely sets the value for global configuration 'Cache.BoostOfIDsMemRatio' field
func SetCacheBoostOfIDsMemRatio(v float64) { global.SetCacheBoostOfIDsMemRatio(v) }

// GetCacheClientMemRatio safely fetches the Configuration value for state's 'Cache.ClientMemRatio' field
func (st *ConfigState) GetCacheClientMemRatio() (v float64) {
	return st.config.Cache.ClientMemRatio
}

// SetCacheClientMemRatio safely sets the Configuration value for state's 'Cache.ClientMemRatio' field
func (st *ConfigState) SetCacheClientMemRatio(v float64) {
	st.config.Cache.ClientMemRatio = v
	st.reloadToViper()
}

// GetCacheClientMemRatio safely fetches the value for global configuration 'Cache.ClientMemRatio' field
func GetCacheClientMemRatio() float64 { return global.GetCacheClientMemRatio() }

// SetCacheClientMemRatio safely sets the value for global configuration 'Cache.ClientMemRatio' field
func SetCacheClientMemRatio(v float64) { global.SetCacheClientMemRatio(v) }

// GetCacheConversationMemRatio safely fetches the Configuration value for state's 'Cache.ConversationMemRatio' field
func (st *ConfigState) GetCacheConversationMemRatio() (v float64) {
	return st.config.Cache.ConversationMemRatio
}

// SetCacheConversationMemRatio safely sets the Configuration value for state's 'Cache.ConversationMemRatio' field
func (st *ConfigState) SetCacheConversationMemRatio(v float64) {
	st.config.Cache.ConversationMemRatio = v
	st.reloadToViper()
}

// GetCacheConversationMemRatio safely fetches the value for global configuration 'Cache.ConversationMemRatio' field
func GetCacheConversationMemRatio() float64 { return global.GetCacheConversationMemRatio() }

// SetCacheConversationMemRatio safely sets the value for global configuration 'Cache.ConversationMemRatio' field
func SetCacheConversationMemRatio(v float64) { global.SetCacheConversationMemRatio(v) }

// GetCacheConversationLastStatusIDsMemRatio safely fetches the Configuration value for state's 'Cache.ConversationLastStatusIDsMemRatio' field
func (st *ConfigState) GetCacheConversationLastStatusIDsMemRatio() (v float64) {
	return st.config.Cache.ConversationLastStatusIDsMemRatio
}

// SetCacheConversationLastStatusIDsMemRatio safely sets the Configuration value for state's 'Cache.ConversationLastStatusIDsMemRatio' field
func (st *ConfigState) SetCacheConversationLastStatusIDsMemRatio(v float64) {
	st.config.Cache.ConversationLastStatusIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheConversationLastStatusIDsMemRatio safely fetches the value for global configuration 'Cache.ConversationLastStatusIDsMemRatio' field
func GetCacheConversationLastStatusIDsMemRatio() float64 {
	return global.GetCacheConversationLastStatusIDsMemRatio()
}

// SetCacheConversationLastStatusIDsMemRatio safely sets the value for global configuration 'Cache.ConversationLastStatusIDsMemRatio' field
func SetCacheConversationLastStatusIDsMemRatio(v float64) {
	global.SetCacheConversationLastStatusIDsMemRatio(v)
}

// GetCacheDomainPermissionDraftMemRatio safely fetches the Configuration value for state's 'Cache.DomainPermissionDraftMemRatio' field
func (st *ConfigState) GetCacheDomainPermissionDraftMemRatio() (v float64) {
	return st.config.Cache.DomainPermissionDraftMemRatio
}

// SetCacheDomainPermissionDraftMemRatio safely sets the Configuration value for state's 'Cache.DomainPermissionDraftMemRatio' field
func (st *ConfigState) SetCacheDomainPermissionDraftMemRatio(v float64) {
	st.config.Cache.DomainPermissionDraftMemRatio = v
	st.reloadToViper()
}

// GetCacheDomainPermissionDraftMemRatio safely fetches the value for global configuration 'Cache.DomainPermissionDraftMemRatio' field
func GetCacheDomainPermissionDraftMemRatio() float64 {
	return global.GetCacheDomainPermissionDraftMemRatio()
}

// SetCacheDomainPermissionDraftMemRatio safely sets the value for global configuration 'Cache.DomainPermissionDraftMemRatio' field
func SetCacheDomainPermissionDraftMemRatio(v float64) {
	global.SetCacheDomainPermissionDraftMemRatio(v)
}

// GetCacheDomainLimitMemRatio safely fetches the Configuration value for state's 'Cache.DomainLimitMemRatio' field
func (st *ConfigState) GetCacheDomainLimitMemRatio() (v float64) {
	return st.config.Cache.DomainLimitMemRatio
}

// SetCacheDomainLimitMemRatio safely sets the Configuration value for state's 'Cache.DomainLimitMemRatio' field
func (st *ConfigState) SetCacheDomainLimitMemRatio(v float64) {
	st.config.Cache.DomainLimitMemRatio = v
	st.reloadToViper()
}

// GetCacheDomainLimitMemRatio safely fetches the value for global configuration 'Cache.DomainLimitMemRatio' field
func GetCacheDomainLimitMemRatio() float64 { return global.GetCacheDomainLimitMemRatio() }

// SetCacheDomainLimitMemRatio safely sets the value for global configuration 'Cache.DomainLimitMemRatio' field
func SetCacheDomainLimitMemRatio(v float64) { global.SetCacheDomainLimitMemRatio(v) }

// GetCacheDomainPermissionSubscriptionMemRatio safely fetches the Configuration value for state's 'Cache.DomainPermissionSubscriptionMemRatio' field
func (st *ConfigState) GetCacheDomainPermissionSubscriptionMemRatio() (v float64) {
	return st.config.Cache.DomainPermissionSubscriptionMemRatio
}

// SetCacheDomainPermissionSubscriptionMemRatio safely sets the Configuration value for state's 'Cache.DomainPermissionSubscriptionMemRatio' field
func (st *ConfigState) SetCacheDomainPermissionSubscriptionMemRatio(v float64) {
	st.config.Cache.DomainPermissionSubscriptionMemRatio = v
	st.reloadToViper()
}

// GetCacheDomainPermissionSubscriptionMemRatio safely fetches the value for global configuration 'Cache.DomainPermissionSubscriptionMemRatio' field
func GetCacheDomainPermissionSubscriptionMemRatio() float64 {
	return global.GetCacheDomainPermissionSubscriptionMemRatio()
}

// SetCacheDomainPermissionSubscriptionMemRatio safely sets the value for global configuration 'Cache.DomainPermissionSubscriptionMemRatio' field
func SetCacheDomainPermissionSubscriptionMemRatio(v float64) {
	global.SetCacheDomainPermissionSubscriptionMemRatio(v)
}

// GetCacheEmojiMemRatio safely fetches the Configuration value for state's 'Cache.EmojiMemRatio' field
func (st *ConfigState) GetCacheEmojiMemRatio() (v float64) {
	return st.config.Cache.EmojiMemRatio
}

// SetCacheEmojiMemRatio safely sets the Configuration value for state's 'Cache.EmojiMemRatio' field
func (st *ConfigState) SetCacheEmojiMemRatio(v float64) {
	st.config.Cache.EmojiMemRatio = v
	st.reloadToViper()
}

// GetCacheEmojiMemRatio safely fetches the value for global configuration 'Cache.EmojiMemRatio' field
func GetCacheEmojiMemRatio() float64 { return global.GetCacheEmojiMemRatio() }

// SetCacheEmojiMemRatio safely sets the value for global configuration 'Cache.EmojiMemRatio' field
func SetCacheEmojiMemRatio(v float64) { global.SetCacheEmojiMemRatio(v) }

// GetCacheEmojiCategoryMemRatio safely fetches the Configuration value for state's 'Cache.EmojiCategoryMemRatio' field
func (st *ConfigState) GetCacheEmojiCategoryMemRatio() (v float64) {
	return st.config.Cache.EmojiCategoryMemRatio
}

// SetCacheEmojiCategoryMemRatio safely sets the Configuration value for state's 'Cache.EmojiCategoryMemRatio' field
func (st *ConfigState) SetCacheEmojiCategoryMemRatio(v float64) {
	st.config.Cache.EmojiCategoryMemRatio = v
	st.reloadToViper()
}

// GetCacheEmojiCategoryMemRatio safely fetches the value for global configuration 'Cache.EmojiCategoryMemRatio' field
func GetCacheEmojiCategoryMemRatio() float64 { return global.GetCacheEmojiCategoryMemRatio() }

// SetCacheEmojiCategoryMemRatio safely sets the value for global configuration 'Cache.EmojiCategoryMemRatio' field
func SetCacheEmojiCategoryMemRatio(v float64) { global.SetCacheEmojiCategoryMemRatio(v) }

// GetCacheFederationErrorMemRatio safely fetches the Configuration value for state's 'Cache.FederationErrorMemRatio' field
func (st *ConfigState) GetCacheFederationErrorMemRatio() (v float64) {
	return st.config.Cache.FederationErrorMemRatio
}

// SetCacheFederationErrorMemRatio safely sets the Configuration value for state's 'Cache.FederationErrorMemRatio' field
func (st *ConfigState) SetCacheFederationErrorMemRatio(v float64) {
	st.config.Cache.FederationErrorMemRatio = v
	st.reloadToViper()
}

// GetCacheFederationErrorMemRatio safely fetches the value for global configuration 'Cache.FederationErrorMemRatio' field
func GetCacheFederationErrorMemRatio() float64 { return global.GetCacheFederationErrorMemRatio() }

// SetCacheFederationErrorMemRatio safely sets the value for global configuration 'Cache.FederationErrorMemRatio' field
func SetCacheFederationErrorMemRatio(v float64) { global.SetCacheFederationErrorMemRatio(v) }

// GetCacheFilterMemRatio safely fetches the Configuration value for state's 'Cache.FilterMemRatio' field
func (st *ConfigState) GetCacheFilterMemRatio() (v float64) {
	return st.config.Cache.FilterMemRatio
}

// SetCacheFilterMemRatio safely sets the Configuration value for state's 'Cache.FilterMemRatio' field
func (st *ConfigState) SetCacheFilterMemRatio(v float64) {
	st.config.Cache.FilterMemRatio = v
	st.reloadToViper()
}

// GetCacheFilterMemRatio safely fetches the value for global configuration 'Cache.FilterMemRatio' field
func GetCacheFilterMemRatio() float64 { return global.GetCacheFilterMemRatio() }

// SetCacheFilterMemRatio safely sets the value for global configuration 'Cache.FilterMemRatio' field
func SetCacheFilterMemRatio(v float64) { global.SetCacheFilterMemRatio(v) }

// GetCacheFilterIDsMemRatio safely fetches the Configuration value for state's 'Cache.FilterIDsMemRatio' field
func (st *ConfigState) GetCacheFilterIDsMemRatio() (v float64) {
	return st.config.Cache.FilterIDsMemRatio
}

// SetCacheFilterIDsMemRatio safely sets the Configuration value for state's 'Cache.FilterIDsMemRatio' field
func (st *ConfigState) SetCacheFilterIDsMemRatio(v float64) {
	st.config.Cache.FilterIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheFilterIDsMemRatio safely fetches the value for global configuration 'Cache.FilterIDsMemRatio' field
func GetCacheFilterIDsMemRatio() float64 { return global.GetCacheFilterIDsMemRatio() }

// SetCacheFilterIDsMemRatio safely sets the value for global configuration 'Cache.FilterIDsMemRatio' field
func SetCacheFilterIDsMemRatio(v float64) { global.SetCacheFilterIDsMemRatio(v) }

// GetCacheFilterKeywordMemRatio safely fetches the Configuration value for state's 'Cache.FilterKeywordMemRatio' field
func (st *ConfigState) GetCacheFilterKeywordMemRatio() (v float64) {
	return st.config.Cache.FilterKeywordMemRatio
}

// SetCacheFilterKeywordMemRatio safely sets the Configuration value for state's 'Cache.FilterKeywordMemRatio' field
func (st *ConfigState) SetCacheFilterKeywordMemRatio(v float64) {
	st.config.Cache.FilterKeywordMemRatio = v
	st.reloadToViper()
}

// GetCacheFilterKeywordMemRatio safely fetches the value for global configuration 'Cache.FilterKeywordMemRatio' field
func GetCacheFilterKeywordMemRatio() float64 { return global.GetCacheFilterKeywordMemRatio() }

// SetCacheFilterKeywordMemRatio safely sets the value for global configuration 'Cache.FilterKeywordMemRatio' field
func SetCacheFilterKeywordMemRatio(v float64) { global.SetCacheFilterKeywordMemRatio(v) }

// GetCacheFilterStatusMemRatio safely fetches the Configuration value for state's 'Cache.FilterStatusMemRatio' field
func (st *ConfigState) GetCacheFilterStatusMemRatio() (v float64) {
	return st.config.Cache.FilterStatusMemRatio
}

// SetCacheFilterStatusMemRatio safely sets the Configuration value for state's 'Cache.FilterStatusMemRatio' field
func (st *ConfigState) SetCacheFilterStatusMemRatio(v float64) {
	st.config.Cache.FilterStatusMemRatio = v
	st.reloadToViper()
}

// GetCacheFilterStatusMemRatio safely fetches the value for global configuration 'Cache.FilterStatusMemRatio' field
func GetCacheFilterStatusMemRatio() float64 { return global.GetCacheFilterStatusMemRatio() }

// SetCacheFilterStatusMemRatio safely sets the value for global configuration 'Cache.FilterStatusMemRatio' field
func SetCacheFilterStatusMemRatio(v float64) { global.SetCacheFilterStatusMemRatio(v) }

// GetCacheFollowMemRatio safely fetches the Configuration value for state's 'Cache.FollowMemRatio' field
func (st *ConfigState) GetCacheFollowMemRatio() (v float64) {
	return st.config.Cache.FollowMemRatio
}

// SetCacheFollowMemRatio safely sets the Configuration value for state's 'Cache.FollowMemRatio' field
func (st *ConfigState) SetCacheFollowMemRatio(v float64) {
	st.config.Cache.FollowMemRatio = v
	st.reloadToViper()
}

// GetCacheFollowMemRatio safely fetches the value for global configuration 'Cache.FollowMemRatio' field
func GetCacheFollowMemRatio() float64 { return global.GetCacheFollowMemRatio() }

// SetCacheFollowMemRatio safely sets the value for global configuration 'Cache.FollowMemRatio' field
func SetCacheFollowMemRatio(v float64) { global.SetCacheFollowMemRatio(v) }

// GetCacheFollowIDsMemRatio safely fetches the Configuration value for state's 'Cache.FollowIDsMemRatio' field
func (st *ConfigState) GetCacheFollowIDsMemRatio() (v float64) {
	return st.config.Cache.FollowIDsMemRatio
}

// SetCacheFollowIDsMemRatio safely sets the Configuration value for state's 'Cache.FollowIDsMemRatio' field
func (st *ConfigState) SetCacheFollowIDsMemRatio(v float64) {
	st.config.Cache.FollowIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheFollowIDsMemRatio safely fetches the value for global configuration 'Cache.FollowIDsMemRatio' field
func GetCacheFollowIDsMemRatio() float64 { return global.GetCacheFollowIDsMemRatio() }

// SetCacheFollowIDsMemRatio safely sets the value for global configuration 'Cache.FollowIDsMemRatio' field
func SetCacheFollowIDsMemRatio(v float64) { global.SetCacheFollowIDsMemRatio(v) }

// GetCacheFollowRequestMemRatio safely fetches the Configuration value for state's 'Cache.FollowRequestMemRatio' field
func (st *ConfigState) GetCacheFollowRequestMemRatio() (v float64) {
	return st.config.Cache.FollowRequestMemRatio
}

// SetCacheFollowRequestMemRatio safely sets the Configuration value for state's 'Cache.FollowRequestMemRatio' field
func (st *ConfigState) SetCacheFollowRequestMemRatio(v float64) {
	st.config.Cache.FollowRequestMemRatio = v
	st.reloadToViper()
}

// GetCacheFollowRequestMemRatio safely fetches the value for global configuration 'Cache.FollowRequestMemRatio' field
func GetCacheFollowRequestMemRatio() float64 { return global.GetCacheFollowRequestMemRatio() }

// SetCacheFollowRequestMemRatio safely sets the value for global configuration 'Cache.FollowRequestMemRatio' field
func SetCacheFollowRequestMemRatio(v float64) { global.SetCacheFollowRequestMemRatio(v) }

// GetCacheFollowRequestIDsMemRatio safely fetches the Configuration value for state's 'Cache.FollowRequestIDsMemRatio' field
func (st *ConfigState) GetCacheFollowRequestIDsMemRatio() (v float64) {
	return st.config.Cache.FollowRequestIDsMemRatio
}

// SetCacheFollowRequestIDsMemRatio safely sets the Configuration value for state's 'Cache.FollowRequestIDsMemRatio' field
func (st *ConfigState) SetCacheFollowRequestIDsMemRatio(v float64) {
	st.config.Cache.FollowRequestIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheFollowRequestIDsMemRatio safely fetches the value for global configuration 'Cache.FollowRequestIDsMemRatio' field
func GetCacheFollowRequestIDsMemRatio() float64 { return global.GetCacheFollowRequestIDsMemRatio() }

// SetCacheFollowRequestIDsMemRatio safely sets the value for global configuration 'Cache.FollowRequestIDsMemRatio' field
func SetCacheFollowRequestIDsMemRatio(v float64) { global.SetCacheFollowRequestIDsMemRatio(v) }

// GetCacheFollowingTagIDsMemRatio safely fetches the Configuration value for state's 'Cache.FollowingTagIDsMemRatio' field
func (st *ConfigState) GetCacheFollowingTagIDsMemRatio() (v float64) {
	return st.config.Cache.FollowingTagIDsMemRatio
}

// SetCacheFollowingTagIDsMemRatio safely sets the Configuration value for state's 'Cache.FollowingTagIDsMemRatio' field
func (st *ConfigState) SetCacheFollowingTagIDsMemRatio(v float64) {
	st.config.Cache.FollowingTagIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheFollowingTagIDsMemRatio safely fetches the value for global configuration 'Cache.FollowingTagIDsMemRatio' field
func GetCacheFollowingTagIDsMemRatio() float64 { return global.GetCacheFollowingTagIDsMemRatio() }

// SetCacheFollowingTagIDsMemRatio safely sets the value for global configuration 'Cache.FollowingTagIDsMemRatio' field
func SetCacheFollowingTagIDsMemRatio(v float64) { global.SetCacheFollowingTagIDsMemRatio(v) }

// GetCacheHomeAccountIDsMemRatio safely fetches the Configuration value for state's 'Cache.HomeAccountIDsMemRatio' field
func (st *ConfigState) GetCacheHomeAccountIDsMemRatio() (v float64) {
	return st.config.Cache.HomeAccountIDsMemRatio
}

// SetCacheHomeAccountIDsMemRatio safely sets the Configuration value for state's 'Cache.HomeAccountIDsMemRatio' field
func (st *ConfigState) SetCacheHomeAccountIDsMemRatio(v float64) {
	st.config.Cache.HomeAccountIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheHomeAccountIDsMemRatio safely fetches the value for global configuration 'Cache.HomeAccountIDsMemRatio' field
func GetCacheHomeAccountIDsMemRatio() float64 { return global.GetCacheHomeAccountIDsMemRatio() }

// SetCacheHomeAccountIDsMemRatio safely sets the value for global configuration 'Cache.HomeAccountIDsMemRatio' field
func SetCacheHomeAccountIDsMemRatio(v float64) { global.SetCacheHomeAccountIDsMemRatio(v) }

// GetCacheInReplyToIDsMemRatio safely fetches the Configuration value for state's 'Cache.InReplyToIDsMemRatio' field
func (st *ConfigState) GetCacheInReplyToIDsMemRatio() (v float64) {
	return st.config.Cache.InReplyToIDsMemRatio
}

// SetCacheInReplyToIDsMemRatio safely sets the Configuration value for state's 'Cache.InReplyToIDsMemRatio' field
func (st *ConfigState) SetCacheInReplyToIDsMemRatio(v float64) {
	st.config.Cache.InReplyToIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheInReplyToIDsMemRatio safely fetches the value for global configuration 'Cache.InReplyToIDsMemRatio' field
func GetCacheInReplyToIDsMemRatio() float64 { return global.GetCacheInReplyToIDsMemRatio() }

// SetCacheInReplyToIDsMemRatio safely sets the value for global configuration 'Cache.InReplyToIDsMemRatio' field
func SetCacheInReplyToIDsMemRatio(v float64) { global.SetCacheInReplyToIDsMemRatio(v) }

// GetCacheInstanceMemRatio safely fetches the Configuration value for state's 'Cache.InstanceMemRatio' field
func (st *ConfigState) GetCacheInstanceMemRatio() (v float64) {
	return st.config.Cache.InstanceMemRatio
}

// SetCacheInstanceMemRatio safely sets the Configuration value for state's 'Cache.InstanceMemRatio' field
func (st *ConfigState) SetCacheInstanceMemRatio(v float64) {
	st.config.Cache.InstanceMemRatio = v
	st.reloadToViper()
}

// GetCacheInstanceMemRatio safely fetches the value for global configuration 'Cache.InstanceMemRatio' field
func GetCacheInstanceMemRatio() float64 { return global.GetCacheInstanceMemRatio() }

// SetCacheInstanceMemRatio safely sets the value for global configuration 'Cache.InstanceMemRatio' field
func SetCacheInstanceMemRatio(v float64) { global.SetCacheInstanceMemRatio(v) }

// GetCacheInteractionRequestMemRatio safely fetches the Configuration value for state's 'Cache.InteractionRequestMemRatio' field
func (st *ConfigState) GetCacheInteractionRequestMemRatio() (v float64) {
	return st.config.Cache.InteractionRequestMemRatio
}

// SetCacheInteractionRequestMemRatio safely sets the Configuration value for state's 'Cache.InteractionRequestMemRatio' field
func (st *ConfigState) SetCacheInteractionRequestMemRatio(v float64) {
	st.config.Cache.InteractionRequestMemRatio = v
	st.reloadToViper()
}

// GetCacheInteractionRequestMemRatio safely fetches the value for global configuration 'Cache.InteractionRequestMemRatio' field
func GetCacheInteractionRequestMemRatio() float64 { return global.GetCacheInteractionRequestMemRatio() }

// SetCacheInteractionRequestMemRatio safely sets the value for global configuration 'Cache.InteractionRequestMemRatio' field
func SetCacheInteractionRequestMemRatio(v float64) { global.SetCacheInteractionRequestMemRatio(v) }

// GetCacheListMemRatio safely fetches the Configuration value for state's 'Cache.ListMemRatio' field
func (st *ConfigState) GetCacheListMemRatio() (v float64) {
	return st.config.Cache.ListMemRatio
}

// SetCacheListMemRatio safely sets the Configuration value for state's 'Cache.ListMemRatio' field
func (st *ConfigState) SetCacheListMemRatio(v float64) {
	st.config.Cache.ListMemRatio = v
	st.reloadToViper()
}

// GetCacheListMemRatio safely fetches the value for global configuration 'Cache.ListMemRatio' field
func GetCacheListMemRatio() float64 { return global.GetCacheListMemRatio() }

// SetCacheListMemRatio safely sets the value for global configuration 'Cache.ListMemRatio' field
func SetCacheListMemRatio(v float64) { global.SetCacheListMemRatio(v) }

// GetCacheListIDsMemRatio safely fetches the Configuration value for state's 'Cache.ListIDsMemRatio' field
func (st *ConfigState) GetCacheListIDsMemRatio() (v float64) {
	return st.config.Cache.ListIDsMemRatio
}

// SetCacheListIDsMemRatio safely sets the Configuration value for state's 'Cache.ListIDsMemRatio' field
func (st *ConfigState) SetCacheListIDsMemRatio(v float64) {
	st.config.Cache.ListIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheListIDsMemRatio safely fetches the value for global configuration 'Cache.ListIDsMemRatio' field
func GetCacheListIDsMemRatio() float64 { return global.GetCacheListIDsMemRatio() }

// SetCacheListIDsMemRatio safely sets the value for global configuration 'Cache.ListIDsMemRatio' field
func SetCacheListIDsMemRatio(v float64) { global.SetCacheListIDsMemRatio(v) }

// GetCacheListedIDsMemRatio safely fetches the Configuration value for state's 'Cache.ListedIDsMemRatio' field
func (st *ConfigState) GetCacheListedIDsMemRatio() (v float64) {
	return st.config.Cache.ListedIDsMemRatio
}

// SetCacheListedIDsMemRatio safely sets the Configuration value for state's 'Cache.ListedIDsMemRatio' field
func (st *ConfigState) SetCacheListedIDsMemRatio(v float64) {
	st.config.Cache.ListedIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheListedIDsMemRatio safely fetches the value for global configuration 'Cache.ListedIDsMemRatio' field
func GetCacheListedIDsMemRatio() float64 { return global.GetCacheListedIDsMemRatio() }

// SetCacheListedIDsMemRatio safely sets the value for global configuration 'Cache.ListedIDsMemRatio' field
func SetCacheListedIDsMemRatio(v float64) { global.SetCacheListedIDsMemRatio(v) }

// GetCacheMarkerMemRatio safely fetches the Configuration value for state's 'Cache.MarkerMemRatio' field
func (st *ConfigState) GetCacheMarkerMemRatio() (v float64) {
	return st.config.Cache.MarkerMemRatio
}

// SetCacheMarkerMemRatio safely sets the Configuration value for state's 'Cache.MarkerMemRatio' field
func (st *ConfigState) SetCacheMarkerMemRatio(v float64) {
	st.config.Cache.MarkerMemRatio = v
	st.reloadToViper()
}

// GetCacheMarkerMemRatio safely fetches the value for global configuration 'Cache.MarkerMemRatio' field
func GetCacheMarkerMemRatio() float64 { return global.GetCacheMarkerMemRatio() }

// SetCacheMarkerMemRatio safely sets the value for global configuration 'Cache.MarkerMemRatio' field
func SetCacheMarkerMemRatio(v float64) { global.SetCacheMarkerMemRatio(v) }

// GetCacheMediaMemRatio safely fetches the Configuration value for state's 'Cache.MediaMemRatio' field
func (st *ConfigState) GetCacheMediaMemRatio() (v float64) {
	return st.config.Cache.MediaMemRatio
}

// SetCacheMediaMemRatio safely sets the Configuration value for state's 'Cache.MediaMemRatio' field
func (st *ConfigState) SetCacheMediaMemRatio(v float64) {
	st.config.Cache.MediaMemRatio = v
	st.reloadToViper()
}

// GetCacheMediaMemRatio safely fetches the value for global configuration 'Cache.MediaMemRatio' field
func GetCacheMediaMemRatio() float64 { return global.GetCacheMediaMemRatio() }

// SetCacheMediaMemRatio safely sets the value for global configuration 'Cache.MediaMemRatio' field
func SetCacheMediaMemRatio(v float64) { global.SetCacheMediaMemRatio(v) }

// GetCacheMentionMemRatio safely fetches the Configuration value for state's 'Cache.MentionMemRatio' field
func (st *ConfigState) GetCacheMentionMemRatio() (v float64) {
	return st.config.Cache.MentionMemRatio
}

// SetCacheMentionMemRatio safely sets the Configuration value for state's 'Cache.MentionMemRatio' field
func (st *ConfigState) SetCacheMentionMemRatio(v float64) {
	st.config.Cache.MentionMemRatio = v
	st.reloadToViper()
}

// GetCacheMentionMemRatio safely fetches the value for global configuration 'Cache.MentionMemRatio' field
func GetCacheMentionMemRatio() float64 { return global.GetCacheMentionMemRatio() }

// SetCacheMentionMemRatio safely sets the value for global configuration 'Cache.MentionMemRatio' field
func SetCacheMentionMemRatio(v float64) { global.SetCacheMentionMemRatio(v) }

// GetCacheMoveMemRatio safely fetches the Configuration value for state's 'Cache.MoveMemRatio' field
func (st *ConfigState) GetCacheMoveMemRatio() (v float64) {
	return st.config.Cache.MoveMemRatio
}

// SetCacheMoveMemRatio safely sets the Configuration value for state's 'Cache.MoveMemRatio' field
func (st *ConfigState) SetCacheMoveMemRatio(v float64) {
	st.config.Cache.MoveMemRatio = v
	st.reloadToViper()
}

// GetCacheMoveMemRatio safely fetches the value for global configuration 'Cache.MoveMemRatio' field
func GetCacheMoveMemRatio() float64 { return global.GetCacheMoveMemRatio() }

// SetCacheMoveMemRatio safely sets the value for global configuration 'Cache.MoveMemRatio' field
func SetCacheMoveMemRatio(v float64) { global.SetCacheMoveMemRatio(v) }

// GetCacheNotificationMemRatio safely fetches the Configuration value for state's 'Cache.NotificationMemRatio' field
func (st *ConfigState) GetCacheNotificationMemRatio() (v float64) {
	return st.config.Cache.NotificationMemRatio
}

// SetCacheNotificationMemRatio safely sets the Configuration value for state's 'Cache.NotificationMemRatio' field
func (st *ConfigState) SetCacheNotificationMemRatio(v float64) {
	st.config.Cache.NotificationMemRatio = v
	st.reloadToViper()
}

// GetCacheNotificationMemRatio safely fetches the value for global configuration 'Cache.NotificationMemRatio' field
func GetCacheNotificationMemRatio() float64 { return global.GetCacheNotificationMemRatio() }

// SetCacheNotificationMemRatio safely sets the value for global configuration 'Cache.NotificationMemRatio' field
func SetCacheNotificationMemRatio(v float64) { global.SetCacheNotificationMemRatio(v) }

// GetCachePollMemRatio safely fetches the Configuration value for state's 'Cache.PollMemRatio' field
func (st *ConfigState) GetCachePollMemRatio() (v float64) {
	return st.config.Cache.PollMemRatio
}

// SetCachePollMemRatio safely sets the Configuration value for state's 'Cache.PollMemRatio' field
func (st *ConfigState) SetCachePollMemRatio(v float64) {
	st.config.Cache.PollMemRatio = v
	st.reloadToViper()
}

// GetCachePollMemRatio safely fetches the value for global configuration 'Cache.PollMemRatio' field
func GetCachePollMemRatio() float64 { return global.GetCachePollMemRatio() }

// SetCachePollMemRatio safely sets the value for global configuration 'Cache.PollMemRatio' field
func SetCachePollMemRatio(v float64) { global.SetCachePollMemRatio(v) }

// GetCachePollVoteMemRatio safely fetches the Configuration value for state's 'Cache.PollVoteMemRatio' field
func (st *ConfigState) GetCachePollVoteMemRatio() (v float64) {
	return st.config.Cache.PollVoteMemRatio
}

// SetCachePollVoteMemRatio safely sets the Configuration value for state's 'Cache.PollVoteMemRatio' field
func (st *ConfigState) SetCachePollVoteMemRatio(v float64) {
	st.config.Cache.PollVoteMemRatio = v
	st.reloadToViper()
}

// GetCachePollVoteMemRatio safely fetches the value for global configuration 'Cache.PollVoteMemRatio' field
func GetCachePollVoteMemRatio() float64 { return global.GetCachePollVoteMemRatio() }

// SetCachePollVoteMemRatio safely sets the value for global configuration 'Cache.PollVoteMemRatio' field
func SetCachePollVoteMemRatio(v float64) { global.SetCachePollVoteMemRatio(v) }

// GetCachePollVoteIDsMemRatio safely fetches the Configuration value for state's 'Cache.PollVoteIDsMemRatio' field
func (st *ConfigState) GetCachePollVoteIDsMemRatio() (v float64) {
	return st.config.Cache.PollVoteIDsMemRatio
}

// SetCachePollVoteIDsMemRatio safely sets the Configuration value for state's 'Cache.PollVoteIDsMemRatio' field
func (st *ConfigState) SetCachePollVoteIDsMemRatio(v float64) {
	st.config.Cache.PollVoteIDsMemRatio = v
	st.reloadToViper()
}

// GetCachePollVoteIDsMemRatio safely fetches the value for global configuration 'Cache.PollVoteIDsMemRatio' field
func GetCachePollVoteIDsMemRatio() float64 { return global.GetCachePollVoteIDsMemRatio() }

// SetCachePollVoteIDsMemRatio safely sets the value for global configuration 'Cache.PollVoteIDsMemRatio' field
func SetCachePollVoteIDsMemRatio(v float64) { global.SetCachePollVoteIDsMemRatio(v) }

// GetCacheReportMemRatio safely fetches the Configuration value for state's 'Cache.ReportMemRatio' field
func (st *ConfigState) GetCacheReportMemRatio() (v float64) {
	return st.config.Cache.ReportMemRatio
}

// SetCacheReportMemRatio safely sets the Configuration value for state's 'Cache.ReportMemRatio' field
func (st *ConfigState) SetCacheReportMemRatio(v float64) {
	st.config.Cache.ReportMemRatio = v
	st.reloadToViper()
}

// GetCacheReportMemRatio safely fetches the value for global configuration 'Cache.ReportMemRatio' field
func GetCacheReportMemRatio() float64 { return global.GetCacheReportMemRatio() }

// SetCacheReportMemRatio safely sets the value for global configuration 'Cache.ReportMemRatio' field
func SetCacheReportMemRatio(v float64) { global.SetCacheReportMemRatio(v) }

// GetCacheRelayActorMemRatio safely fetches the Configuration value for state's 'Cache.RelayActorMemRatio' field
func (st *ConfigState) GetCacheRelayActorMemRatio() (v float64) {
	return st.config.Cache.RelayActorMemRatio
}

// SetCacheRelayActorMemRatio safely sets the Configuration value for state's 'Cache.RelayActorMemRatio' field
func (st *ConfigState) SetCacheRelayActorMemRatio(v float64) {
	st.config.Cache.RelayActorMemRatio = v
	st.reloadToViper()
}

// GetCacheRelayActorMemRatio safely fetches the value for global configuration 'Cache.RelayActorMemRatio' field
func GetCacheRelayActorMemRatio() float64 { return global.GetCacheRelayActorMemRatio() }

// SetCacheRelayActorMemRatio safely sets the value for global configuration 'Cache.RelayActorMemRatio' field
func SetCacheRelayActorMemRatio(v float64) { global.SetCacheRelayActorMemRatio(v) }

// GetCacheRelayMatcherMemRatio safely fetches the Configuration value for state's 'Cache.RelayMatcherMemRatio' field
func (st *ConfigState) GetCacheRelayMatcherMemRatio() (v float64) {
	return st.config.Cache.RelayMatcherMemRatio
}

// SetCacheRelayMatcherMemRatio safely sets the Configuration value for state's 'Cache.RelayMatcherMemRatio' field
func (st *ConfigState) SetCacheRelayMatcherMemRatio(v float64) {
	st.config.Cache.RelayMatcherMemRatio = v
	st.reloadToViper()
}

// GetCacheRelayMatcherMemRatio safely fetches the value for global configuration 'Cache.RelayMatcherMemRatio' field
func GetCacheRelayMatcherMemRatio() float64 { return global.GetCacheRelayMatcherMemRatio() }

// SetCacheRelayMatcherMemRatio safely sets the value for global configuration 'Cache.RelayMatcherMemRatio' field
func SetCacheRelayMatcherMemRatio(v float64) { global.SetCacheRelayMatcherMemRatio(v) }

// GetCacheRelayPushMemRatio safely fetches the Configuration value for state's 'Cache.RelayPushMemRatio' field
func (st *ConfigState) GetCacheRelayPushMemRatio() (v float64) {
	return st.config.Cache.RelayPushMemRatio
}

// SetCacheRelayPushMemRatio safely sets the Configuration value for state's 'Cache.RelayPushMemRatio' field
func (st *ConfigState) SetCacheRelayPushMemRatio(v float64) {
	st.config.Cache.RelayPushMemRatio = v
	st.reloadToViper()
}

// GetCacheRelayPushMemRatio safely fetches the value for global configuration 'Cache.RelayPushMemRatio' field
func GetCacheRelayPushMemRatio() float64 { return global.GetCacheRelayPushMemRatio() }

// SetCacheRelayPushMemRatio safely sets the value for global configuration 'Cache.RelayPushMemRatio' field
func SetCacheRelayPushMemRatio(v float64) { global.SetCacheRelayPushMemRatio(v) }

// GetCacheRelayPushIDsMemRatio safely fetches the Configuration value for state's 'Cache.RelayPushIDsMemRatio' field
func (st *ConfigState) GetCacheRelayPushIDsMemRatio() (v float64) {
	return st.config.Cache.RelayPushIDsMemRatio
}

// SetCacheRelayPushIDsMemRatio safely sets the Configuration value for state's 'Cache.RelayPushIDsMemRatio' field
func (st *ConfigState) SetCacheRelayPushIDsMemRatio(v float64) {
	st.config.Cache.RelayPushIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheRelayPushIDsMemRatio safely fetches the value for global configuration 'Cache.RelayPushIDsMemRatio' field
func GetCacheRelayPushIDsMemRatio() float64 { return global.GetCacheRelayPushIDsMemRatio() }

// SetCacheRelayPushIDsMemRatio safely sets the value for global configuration 'Cache.RelayPushIDsMemRatio' field
func SetCacheRelayPushIDsMemRatio(v float64) { global.SetCacheRelayPushIDsMemRatio(v) }

// GetCacheRelaySubscriptionMemRatio safely fetches the Configuration value for state's 'Cache.RelaySubscriptionMemRatio' field
func (st *ConfigState) GetCacheRelaySubscriptionMemRatio() (v float64) {
	return st.config.Cache.RelaySubscriptionMemRatio
}

// SetCacheRelaySubscriptionMemRatio safely sets the Configuration value for state's 'Cache.RelaySubscriptionMemRatio' field
func (st *ConfigState) SetCacheRelaySubscriptionMemRatio(v float64) {
	st.config.Cache.RelaySubscriptionMemRatio = v
	st.reloadToViper()
}

// GetCacheRelaySubscriptionMemRatio safely fetches the value for global configuration 'Cache.RelaySubscriptionMemRatio' field
func GetCacheRelaySubscriptionMemRatio() float64 { return global.GetCacheRelaySubscriptionMemRatio() }

// SetCacheRelaySubscriptionMemRatio safely sets the value for global configuration 'Cache.RelaySubscriptionMemRatio' field
func SetCacheRelaySubscriptionMemRatio(v float64) { global.SetCacheRelaySubscriptionMemRatio(v) }

// GetCacheScheduledStatusMemRatio safely fetches the Configuration value for state's 'Cache.ScheduledStatusMemRatio' field
func (st *ConfigState) GetCacheScheduledStatusMemRatio() (v float64) {
	return st.config.Cache.ScheduledStatusMemRatio
}

// SetCacheScheduledStatusMemRatio safely sets the Configuration value for state's 'Cache.ScheduledStatusMemRatio' field
func (st *ConfigState) SetCacheScheduledStatusMemRatio(v float64) {
	st.config.Cache.ScheduledStatusMemRatio = v
	st.reloadToViper()
}

// GetCacheScheduledStatusMemRatio safely fetches the value for global configuration 'Cache.ScheduledStatusMemRatio' field
func GetCacheScheduledStatusMemRatio() float64 { return global.GetCacheScheduledStatusMemRatio() }

// SetCacheScheduledStatusMemRatio safely sets the value for global configuration 'Cache.ScheduledStatusMemRatio' field
func SetCacheScheduledStatusMemRatio(v float64) { global.SetCacheScheduledStatusMemRatio(v) }

// GetCacheSinBinStatusMemRatio safely fetches the Configuration value for state's 'Cache.SinBinStatusMemRatio' field
func (st *ConfigState) GetCacheSinBinStatusMemRatio() (v float64) {
	return st.config.Cache.SinBinStatusMemRatio
}

// SetCacheSinBinStatusMemRatio safely sets the Configuration value for state's 'Cache.SinBinStatusMemRatio' field
func (st *ConfigState) SetCacheSinBinStatusMemRatio(v float64) {
	st.config.Cache.SinBinStatusMemRatio = v
	st.reloadToViper()
}

// GetCacheSinBinStatusMemRatio safely fetches the value for global configuration 'Cache.SinBinStatusMemRatio' field
func GetCacheSinBinStatusMemRatio() float64 { return global.GetCacheSinBinStatusMemRatio() }

// SetCacheSinBinStatusMemRatio safely sets the value for global configuration 'Cache.SinBinStatusMemRatio' field
func SetCacheSinBinStatusMemRatio(v float64) { global.SetCacheSinBinStatusMemRatio(v) }

// GetCacheStatusMemRatio safely fetches the Configuration value for state's 'Cache.StatusMemRatio' field
func (st *ConfigState) GetCacheStatusMemRatio() (v float64) {
	return st.config.Cache.StatusMemRatio
}

// SetCacheStatusMemRatio safely sets the Configuration value for state's 'Cache.StatusMemRatio' field
func (st *ConfigState) SetCacheStatusMemRatio(v float64) {
	st.config.Cache.StatusMemRatio = v
	st.reloadToViper()
}

// GetCacheStatusMemRatio safely fetches the value for global configuration 'Cache.StatusMemRatio' field
func GetCacheStatusMemRatio() float64 { return global.GetCacheStatusMemRatio() }

// SetCacheStatusMemRatio safely sets the value for global configuration 'Cache.StatusMemRatio' field
func SetCacheStatusMemRatio(v float64) { global.SetCacheStatusMemRatio(v) }

// GetCacheStatusBookmarkMemRatio safely fetches the Configuration value for state's 'Cache.StatusBookmarkMemRatio' field
func (st *ConfigState) GetCacheStatusBookmarkMemRatio() (v float64) {
	return st.config.Cache.StatusBookmarkMemRatio
}

// SetCacheStatusBookmarkMemRatio safely sets the Configuration value for state's 'Cache.StatusBookmarkMemRatio' field
func (st *ConfigState) SetCacheStatusBookmarkMemRatio(v float64) {
	st.config.Cache.StatusBookmarkMemRatio = v
	st.reloadToViper()
}

// GetCacheStatusBookmarkMemRatio safely fetches the value for global configuration 'Cache.StatusBookmarkMemRatio' field
func GetCacheStatusBookmarkMemRatio() float64 { return global.GetCacheStatusBookmarkMemRatio() }

// SetCacheStatusBookmarkMemRatio safely sets the value for global configuration 'Cache.StatusBookmarkMemRatio' field
func SetCacheStatusBookmarkMemRatio(v float64) { global.SetCacheStatusBookmarkMemRatio(v) }

// GetCacheStatusBookmarkIDsMemRatio safely fetches the Configuration value for state's 'Cache.StatusBookmarkIDsMemRatio' field
func (st *ConfigState) GetCacheStatusBookmarkIDsMemRatio() (v float64) {
	return st.config.Cache.StatusBookmarkIDsMemRatio
}

// SetCacheStatusBookmarkIDsMemRatio safely sets the Configuration value for state's 'Cache.StatusBookmarkIDsMemRatio' field
func (st *ConfigState) SetCacheStatusBookmarkIDsMemRatio(v float64) {
	st.config.Cache.StatusBookmarkIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheStatusBookmarkIDsMemRatio safely fetches the value for global configuration 'Cache.StatusBookmarkIDsMemRatio' field
func GetCacheStatusBookmarkIDsMemRatio() float64 { return global.GetCacheStatusBookmarkIDsMemRatio() }

// SetCacheStatusBookmarkIDsMemRatio safely sets the value for global configuration 'Cache.StatusBookmarkIDsMemRatio' field
func SetCacheStatusBookmarkIDsMemRatio(v float64) { global.SetCacheStatusBookmarkIDsMemRatio(v) }

// GetCacheStatusEditMemRatio safely fetches the Configuration value for state's 'Cache.StatusEditMemRatio' field
func (st *ConfigState) GetCacheStatusEditMemRatio() (v float64) {
	return st.config.Cache.StatusEditMemRatio
}

// SetCacheStatusEditMemRatio safely sets the Configuration value for state's 'Cache.StatusEditMemRatio' field
func (st *ConfigState) SetCacheStatusEditMemRatio(v float64) {
	st.config.Cache.StatusEditMemRatio = v
	st.reloadToViper()
}

// GetCacheStatusEditMemRatio safely fetches the value for global configuration 'Cache.StatusEditMemRatio' field
func GetCacheStatusEditMemRatio() float64 { return global.GetCacheStatusEditMemRatio() }

// SetCacheStatusEditMemRatio safely sets the value for global configuration 'Cache.StatusEditMemRatio' field
func SetCacheStatusEditMemRatio(v float64) { global.SetCacheStatusEditMemRatio(v) }

// GetCacheStatusFaveMemRatio safely fetches the Configuration value for state's 'Cache.StatusFaveMemRatio' field
func (st *ConfigState) GetCacheStatusFaveMemRatio() (v float64) {
	return st.config.Cache.StatusFaveMemRatio
}

// SetCacheStatusFaveMemRatio safely sets the Configuration value for state's 'Cache.StatusFaveMemRatio' field
func (st *ConfigState) SetCacheStatusFaveMemRatio(v float64) {
	st.config.Cache.StatusFaveMemRatio = v
	st.reloadToViper()
}

// GetCacheStatusFaveMemRatio safely fetches the value for global configuration 'Cache.StatusFaveMemRatio' field
func GetCacheStatusFaveMemRatio() float64 { return global.GetCacheStatusFaveMemRatio() }

// SetCacheStatusFaveMemRatio safely sets the value for global configuration 'Cache.StatusFaveMemRatio' field
func SetCacheStatusFaveMemRatio(v float64) { global.SetCacheStatusFaveMemRatio(v) }

// GetCacheStatusFaveIDsMemRatio safely fetches the Configuration value for state's 'Cache.StatusFaveIDsMemRatio' field
func (st *ConfigState) GetCacheStatusFaveIDsMemRatio() (v float64) {
	return st.config.Cache.StatusFaveIDsMemRatio
}

// SetCacheStatusFaveIDsMemRatio safely sets the Configuration value for state's 'Cache.StatusFaveIDsMemRatio' field
func (st *ConfigState) SetCacheStatusFaveIDsMemRatio(v float64) {
	st.config.Cache.StatusFaveIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheStatusFaveIDsMemRatio safely fetches the value for global configuration 'Cache.StatusFaveIDsMemRatio' field
func GetCacheStatusFaveIDsMemRatio() float64 { return global.GetCacheStatusFaveIDsMemRatio() }

// SetCacheStatusFaveIDsMemRatio safely sets the value for global configuration 'Cache.StatusFaveIDsMemRatio' field
func SetCacheStatusFaveIDsMemRatio(v float64) { global.SetCacheStatusFaveIDsMemRatio(v) }

// GetCacheStatusPinnedIDsMemRatio safely fetches the Configuration value for state's 'Cache.StatusPinnedIDsMemRatio' field
func (st *ConfigState) GetCacheStatusPinnedIDsMemRatio() (v float64) {
	return st.config.Cache.StatusPinnedIDsMemRatio
}

// SetCacheStatusPinnedIDsMemRatio safely sets the Configuration value for state's 'Cache.StatusPinnedIDsMemRatio' field
func (st *ConfigState) SetCacheStatusPinnedIDsMemRatio(v float64) {
	st.config.Cache.StatusPinnedIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheStatusPinnedIDsMemRatio safely fetches the value for global configuration 'Cache.StatusPinnedIDsMemRatio' field
func GetCacheStatusPinnedIDsMemRatio() float64 { return global.GetCacheStatusPinnedIDsMemRatio() }

// SetCacheStatusPinnedIDsMemRatio safely sets the value for global configuration 'Cache.StatusPinnedIDsMemRatio' field
func SetCacheStatusPinnedIDsMemRatio(v float64) { global.SetCacheStatusPinnedIDsMemRatio(v) }

// GetCacheTagMemRatio safely fetches the Configuration value for state's 'Cache.TagMemRatio' field
func (st *ConfigState) GetCacheTagMemRatio() (v float64) {
	return st.config.Cache.TagMemRatio
}

// SetCacheTagMemRatio safely sets the Configuration value for state's 'Cache.TagMemRatio' field
func (st *ConfigState) SetCacheTagMemRatio(v float64) {
	st.config.Cache.TagMemRatio = v
	st.reloadToViper()
}

// GetCacheTagMemRatio safely fetches the value for global configuration 'Cache.TagMemRatio' field
func GetCacheTagMemRatio() float64 { return global.GetCacheTagMemRatio() }

// SetCacheTagMemRatio safely sets the value for global configuration 'Cache.TagMemRatio' field
func SetCacheTagMemRatio(v float64) { global.SetCacheTagMemRatio(v) }

// GetCacheThreadMuteMemRatio safely fetches the Configuration value for state's 'Cache.ThreadMuteMemRatio' field
func (st *ConfigState) GetCacheThreadMuteMemRatio() (v float64) {
	return st.config.Cache.ThreadMuteMemRatio
}

// SetCacheThreadMuteMemRatio safely sets the Configuration value for state's 'Cache.ThreadMuteMemRatio' field
func (st *ConfigState) SetCacheThreadMuteMemRatio(v float64) {
	st.config.Cache.ThreadMuteMemRatio = v
	st.reloadToViper()
}

// GetCacheThreadMuteMemRatio safely fetches the value for global configuration 'Cache.ThreadMuteMemRatio' field
func GetCacheThreadMuteMemRatio() float64 { return global.GetCacheThreadMuteMemRatio() }

// SetCacheThreadMuteMemRatio safely sets the value for global configuration 'Cache.ThreadMuteMemRatio' field
func SetCacheThreadMuteMemRatio(v float64) { global.SetCacheThreadMuteMemRatio(v) }

// GetCacheTokenMemRatio safely fetches the Configuration value for state's 'Cache.TokenMemRatio' field
func (st *ConfigState) GetCacheTokenMemRatio() (v float64) {
	return st.config.Cache.TokenMemRatio
}

// SetCacheTokenMemRatio safely sets the Configuration value for state's 'Cache.TokenMemRatio' field
func (st *ConfigState) SetCacheTokenMemRatio(v float64) {
	st.config.Cache.TokenMemRatio = v
	st.reloadToViper()
}

// GetCacheTokenMemRatio safely fetches the value for global configuration 'Cache.TokenMemRatio' field
func GetCacheTokenMemRatio() float64 { return global.GetCacheTokenMemRatio() }

// SetCacheTokenMemRatio safely sets the value for global configuration 'Cache.TokenMemRatio' field
func SetCacheTokenMemRatio(v float64) { global.SetCacheTokenMemRatio(v) }

// GetCacheTombstoneMemRatio safely fetches the Configuration value for state's 'Cache.TombstoneMemRatio' field
func (st *ConfigState) GetCacheTombstoneMemRatio() (v float64) {
	return st.config.Cache.TombstoneMemRatio
}

// SetCacheTombstoneMemRatio safely sets the Configuration value for state's 'Cache.TombstoneMemRatio' field
func (st *ConfigState) SetCacheTombstoneMemRatio(v float64) {
	st.config.Cache.TombstoneMemRatio = v
	st.reloadToViper()
}

// GetCacheTombstoneMemRatio safely fetches the value for global configuration 'Cache.TombstoneMemRatio' field
func GetCacheTombstoneMemRatio() float64 { return global.GetCacheTombstoneMemRatio() }

// SetCacheTombstoneMemRatio safely sets the value for global configuration 'Cache.TombstoneMemRatio' field
func SetCacheTombstoneMemRatio(v float64) { global.SetCacheTombstoneMemRatio(v) }

// GetCacheUserMemRatio safely fetches the Configuration value for state's 'Cache.UserMemRatio' field
func (st *ConfigState) GetCacheUserMemRatio() (v float64) {
	return st.config.Cache.UserMemRatio
}

// SetCacheUserMemRatio safely sets the Configuration value for state's 'Cache.UserMemRatio' field
func (st *ConfigState) SetCacheUserMemRatio(v float64) {
	st.config.Cache.UserMemRatio = v
	st.reloadToViper()
}

// GetCacheUserMemRatio safely fetches the value for global configuration 'Cache.UserMemRatio' field
func GetCacheUserMemRatio() float64 { return global.GetCacheUserMemRatio() }

// SetCacheUserMemRatio safely sets the value for global configuration 'Cache.UserMemRatio' field
func SetCacheUserMemRatio(v float64) { global.SetCacheUserMemRatio(v) }

// GetCacheUserMuteMemRatio safely fetches the Configuration value for state's 'Cache.UserMuteMemRatio' field
func (st *ConfigState) GetCacheUserMuteMemRatio() (v float64) {
	return st.config.Cache.UserMuteMemRatio
}

// SetCacheUserMuteMemRatio safely sets the Configuration value for state's 'Cache.UserMuteMemRatio' field
func (st *ConfigState) SetCacheUserMuteMemRatio(v float64) {
	st.config.Cache.UserMuteMemRatio = v
	st.reloadToViper()
}

// GetCacheUserMuteMemRatio safely fetches the value for global configuration 'Cache.UserMuteMemRatio' field
func GetCacheUserMuteMemRatio() float64 { return global.GetCacheUserMuteMemRatio() }

// SetCacheUserMuteMemRatio safely sets the value for global configuration 'Cache.UserMuteMemRatio' field
func SetCacheUserMuteMemRatio(v float64) { global.SetCacheUserMuteMemRatio(v) }

// GetCacheUserMuteIDsMemRatio safely fetches the Configuration value for state's 'Cache.UserMuteIDsMemRatio' field
func (st *ConfigState) GetCacheUserMuteIDsMemRatio() (v float64) {
	return st.config.Cache.UserMuteIDsMemRatio
}

// SetCacheUserMuteIDsMemRatio safely sets the Configuration value for state's 'Cache.UserMuteIDsMemRatio' field
func (st *ConfigState) SetCacheUserMuteIDsMemRatio(v float64) {
	st.config.Cache.UserMuteIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheUserMuteIDsMemRatio safely fetches the value for global configuration 'Cache.UserMuteIDsMemRatio' field
func GetCacheUserMuteIDsMemRatio() float64 { return global.GetCacheUserMuteIDsMemRatio() }

// SetCacheUserMuteIDsMemRatio safely sets the value for global configuration 'Cache.UserMuteIDsMemRatio' field
func SetCacheUserMuteIDsMemRatio(v float64) { global.SetCacheUserMuteIDsMemRatio(v) }

// GetCacheWebfingerMemRatio safely fetches the Configuration value for state's 'Cache.WebfingerMemRatio' field
func (st *ConfigState) GetCacheWebfingerMemRatio() (v float64) {
	return st.config.Cache.WebfingerMemRatio
}

// SetCacheWebfingerMemRatio safely sets the Configuration value for state's 'Cache.WebfingerMemRatio' field
func (st *ConfigState) SetCacheWebfingerMemRatio(v float64) {
	st.config.Cache.WebfingerMemRatio = v
	st.reloadToViper()
}

// GetCacheWebfingerMemRatio safely fetches the value for global configuration 'Cache.WebfingerMemRatio' field
func GetCacheWebfingerMemRatio() float64 { return global.GetCacheWebfingerMemRatio() }

// SetCacheWebfingerMemRatio safely sets the value for global configuration 'Cache.WebfingerMemRatio' field
func SetCacheWebfingerMemRatio(v float64) { global.SetCacheWebfingerMemRatio(v) }

// GetCacheWebPushSubscriptionMemRatio safely fetches the Configuration value for state's 'Cache.WebPushSubscriptionMemRatio' field
func (st *ConfigState) GetCacheWebPushSubscriptionMemRatio() (v float64) {
	return st.config.Cache.WebPushSubscriptionMemRatio
}

// SetCacheWebPushSubscriptionMemRatio safely sets the Configuration value for state's 'Cache.WebPushSubscriptionMemRatio' field
func (st *ConfigState) SetCacheWebPushSubscriptionMemRatio(v float64) {
	st.config.Cache.WebPushSubscriptionMemRatio = v
	st.reloadToViper()
}

// GetCacheWebPushSubscriptionMemRatio safely fetches the value for global configuration 'Cache.WebPushSubscriptionMemRatio' field
func GetCacheWebPushSubscriptionMemRatio() float64 {
	return global.GetCacheWebPushSubscriptionMemRatio()
}

// SetCacheWebPushSubscriptionMemRatio safely sets the value for global configuration 'Cache.WebPushSubscriptionMemRatio' field
func SetCacheWebPushSubscriptionMemRatio(v float64) { global.SetCacheWebPushSubscriptionMemRatio(v) }

// GetCacheWebPushSubscriptionIDsMemRatio safely fetches the Configuration value for state's 'Cache.WebPushSubscriptionIDsMemRatio' field
func (st *ConfigState) GetCacheWebPushSubscriptionIDsMemRatio() (v float64) {
	return st.config.Cache.WebPushSubscriptionIDsMemRatio
}

// SetCacheWebPushSubscriptionIDsMemRatio safely sets the Configuration value for state's 'Cache.WebPushSubscriptionIDsMemRatio' field
func (st *ConfigState) SetCacheWebPushSubscriptionIDsMemRatio(v float64) {
	st.config.Cache.WebPushSubscriptionIDsMemRatio = v
	st.reloadToViper()
}

// GetCacheWebPushSubscriptionIDsMemRatio safely fetches the value for global configuration 'Cache.WebPushSubscriptionIDsMemRatio' field
func GetCacheWebPushSubscriptionIDsMemRatio() float64 {
	return global.GetCacheWebPushSubscriptionIDsMemRatio()
}

// SetCacheWebPushSubscriptionIDsMemRatio safely sets the value for global configuration 'Cache.WebPushSubscriptionIDsMemRatio' field
func SetCacheWebPushSubscriptionIDsMemRatio(v float64) {
	global.SetCacheWebPushSubscriptionIDsMemRatio(v)
}

// GetCacheMutesMemRatio safely fetches the Configuration value for state's 'Cache.MutesMemRatio' field
func (st *ConfigState) GetCacheMutesMemRatio() (v float64) {
	return st.config.Cache.MutesMemRatio
}

// SetCacheMutesMemRatio safely sets the Configuration value for state's 'Cache.MutesMemRatio' field
func (st *ConfigState) SetCacheMutesMemRatio(v float64) {
	st.config.Cache.MutesMemRatio = v
	st.reloadToViper()
}

// GetCacheMutesMemRatio safely fetches the value for global configuration 'Cache.MutesMemRatio' field
func GetCacheMutesMemRatio() float64 { return global.GetCacheMutesMemRatio() }

// SetCacheMutesMemRatio safely sets the value for global configuration 'Cache.MutesMemRatio' field
func SetCacheMutesMemRatio(v float64) { global.SetCacheMutesMemRatio(v) }

// GetCacheStatusFilterMemRatio safely fetches the Configuration value for state's 'Cache.StatusFilterMemRatio' field
func (st *ConfigState) GetCacheStatusFilterMemRatio() (v float64) {
	return st.config.Cache.StatusFilterMemRatio
}

// SetCacheStatusFilterMemRatio safely sets the Configuration value for state's 'Cache.StatusFilterMemRatio' field
func (st *ConfigState) SetCacheStatusFilterMemRatio(v float64) {
	st.config.Cache.StatusFilterMemRatio = v
	st.reloadToViper()
}

// GetCacheStatusFilterMemRatio safely fetches the value for global configuration 'Cache.StatusFilterMemRatio' field
func GetCacheStatusFilterMemRatio() float64 { return global.GetCacheStatusFilterMemRatio() }

// SetCacheStatusFilterMemRatio safely sets the value for global configuration 'Cache.StatusFilterMemRatio' field
func SetCacheStatusFilterMemRatio(v float64) { global.SetCacheStatusFilterMemRatio(v) }

// GetCacheVisibilityMemRatio safely fetches the Configuration value for state's 'Cache.VisibilityMemRatio' field
func (st *ConfigState) GetCacheVisibilityMemRatio() (v float64) {
	return st.config.Cache.VisibilityMemRatio
}

// SetCacheVisibilityMemRatio safely sets the Configuration value for state's 'Cache.VisibilityMemRatio' field
func (st *ConfigState) SetCacheVisibilityMemRatio(v float64) {
	st.config.Cache.VisibilityMemRatio = v
	st.reloadToViper()
}

// GetCacheVisibilityMemRatio safely fetches the value for global configuration 'Cache.VisibilityMemRatio' field
func GetCacheVisibilityMemRatio() float64 { return global.GetCacheVisibilityMemRatio() }

// SetCacheVisibilityMemRatio safely sets the value for global configuration 'Cache.VisibilityMemRatio' field
func SetCacheVisibilityMemRatio(v float64) { global.SetCacheVisibilityMemRatio(v) }

// GetLogLevel safely fetches the Configuration value for state's 'LogLevel' field
func (st *ConfigState) GetLogLevel() (v string) {
	return st.config.LogLevel
}

// SetLogLevel safely sets the Configuration value for state's 'LogLevel' field
func (st *ConfigState) SetLogLevel(v string) {
	st.config.LogLevel = v
	st.reloadToViper()
}

// GetLogLevel safely fetches the value for global configuration 'LogLevel' field
func GetLogLevel() string { return global.GetLogLevel() }

// SetLogLevel safely sets the value for global configuration 'LogLevel' field
func SetLogLevel(v string) { global.SetLogLevel(v) }

// GetLogFormat safely fetches the Configuration value for state's 'LogFormat' field
func (st *ConfigState) GetLogFormat() (v string) {
	return st.config.LogFormat
}

// SetLogFormat safely sets the Configuration value for state's 'LogFormat' field
func (st *ConfigState) SetLogFormat(v string) {
	st.config.LogFormat = v
	st.reloadToViper()
}

// GetLogFormat safely fetches the value for global configuration 'LogFormat' field
func GetLogFormat() string { return global.GetLogFormat() }

// SetLogFormat safely sets the value for global configuration 'LogFormat' field
func SetLogFormat(v string) { global.SetLogFormat(v) }

// GetLogTimestampFormat safely fetches the Configuration value for state's 'LogTimestampFormat' field
func (st *ConfigState) GetLogTimestampFormat() (v string) {
	return st.config.LogTimestampFormat
}

// SetLogTimestampFormat safely sets the Configuration value for state's 'LogTimestampFormat' field
func (st *ConfigState) SetLogTimestampFormat(v string) {
	st.config.LogTimestampFormat = v
	st.reloadToViper()
}

// GetLogTimestampFormat safely fetches the value for global configuration 'LogTimestampFormat' field
func GetLogTimestampFormat() string { return global.GetLogTimestampFormat() }

// SetLogTimestampFormat safely sets the value for global configuration 'LogTimestampFormat' field
func SetLogTimestampFormat(v string) { global.SetLogTimestampFormat(v) }

// GetLogDbQueries safely fetches the Configuration value for state's 'LogDbQueries' field
func (st *ConfigState) GetLogDbQueries() (v bool) {
	return st.config.LogDbQueries
}

// SetLogDbQueries safely sets the Configuration value for state's 'LogDbQueries' field
func (st *ConfigState) SetLogDbQueries(v bool) {
	st.config.LogDbQueries = v
	st.reloadToViper()
}

// GetLogDbQueries safely fetches the value for global configuration 'LogDbQueries' field
func GetLogDbQueries() bool { return global.GetLogDbQueries() }

// SetLogDbQueries safely sets the value for global configuration 'LogDbQueries' field
func SetLogDbQueries(v bool) { global.SetLogDbQueries(v) }

// GetLogClientIP safely fetches the Configuration value for state's 'LogClientIP' field
func (st *ConfigState) GetLogClientIP() (v bool) {
	return st.config.LogClientIP
}

// SetLogClientIP safely sets the Configuration value for state's 'LogClientIP' field
func (st *ConfigState) SetLogClientIP(v bool) {
	st.config.LogClientIP = v
	st.reloadToViper()
}

// GetLogClientIP safely fetches the value for global configuration 'LogClientIP' field
func GetLogClientIP() bool { return global.GetLogClientIP() }

// SetLogClientIP safely sets the value for global configuration 'LogClientIP' field
func SetLogClientIP(v bool) { global.SetLogClientIP(v) }

// GetRequestIDHeader safely fetches the Configuration value for state's 'RequestIDHeader' field
func (st *ConfigState) GetRequestIDHeader() (v string) {
	return st.config.RequestIDHeader
}

// SetRequestIDHeader safely sets the Configuration value for state's 'RequestIDHeader' field
func (st *ConfigState) SetRequestIDHeader(v string) {
	st.config.RequestIDHeader = v
	st.reloadToViper()
}

// GetRequestIDHeader safely fetches the value for global configuration 'RequestIDHeader' field
func GetRequestIDHeader() string { return global.GetRequestIDHeader() }

// SetRequestIDHeader safely sets the value for global configuration 'RequestIDHeader' field
func SetRequestIDHeader(v string) { global.SetRequestIDHeader(v) }

// GetConfigPath safely fetches the Configuration value for state's 'ConfigPath' field
func (st *ConfigState) GetConfigPath() (v string) {
	return st.config.ConfigPath
}

// SetConfigPath safely sets the Configuration value for state's 'ConfigPath' field
func (st *ConfigState) SetConfigPath(v string) {
	st.config.ConfigPath = v
	st.reloadToViper()
}

// GetConfigPath safely fetches the value for global configuration 'ConfigPath' field
func GetConfigPath() string { return global.GetConfigPath() }

// SetConfigPath safely sets the value for global configuration 'ConfigPath' field
func SetConfigPath(v string) { global.SetConfigPath(v) }

// GetApplicationName safely fetches the Configuration value for state's 'ApplicationName' field
func (st *ConfigState) GetApplicationName() (v string) {
	return st.config.ApplicationName
}

// SetApplicationName safely sets the Configuration value for state's 'ApplicationName' field
func (st *ConfigState) SetApplicationName(v string) {
	st.config.ApplicationName = v
	st.reloadToViper()
}

// GetApplicationName safely fetches the value for global configuration 'ApplicationName' field
func GetApplicationName() string { return global.GetApplicationName() }

// SetApplicationName safely sets the value for global configuration 'ApplicationName' field
func SetApplicationName(v string) { global.SetApplicationName(v) }

// GetLandingPageUser safely fetches the Configuration value for state's 'LandingPageUser' field
func (st *ConfigState) GetLandingPageUser() (v string) {
	return st.config.LandingPageUser
}

// SetLandingPageUser safely sets the Configuration value for state's 'LandingPageUser' field
func (st *ConfigState) SetLandingPageUser(v string) {
	st.config.LandingPageUser = v
	st.reloadToViper()
}

// GetLandingPageUser safely fetches the value for global configuration 'LandingPageUser' field
func GetLandingPageUser() string { return global.GetLandingPageUser() }

// SetLandingPageUser safely sets the value for global configuration 'LandingPageUser' field
func SetLandingPageUser(v string) { global.SetLandingPageUser(v) }

// GetHost safely fetches the Configuration value for state's 'Host' field
func (st *ConfigState) GetHost() (v string) {
	return st.config.Host
}

// SetHost safely sets the Configuration value for state's 'Host' field
func (st *ConfigState) SetHost(v string) {
	st.config.Host = v
	st.reloadToViper()
}

// GetHost safely fetches the value for global configuration 'Host' field
func GetHost() string { return global.GetHost() }

// SetHost safely sets the value for global configuration 'Host' field
func SetHost(v string) { global.SetHost(v) }

// GetAccountDomain safely fetches the Configuration value for state's 'AccountDomain' field
func (st *ConfigState) GetAccountDomain() (v string) {
	return st.config.AccountDomain
}

// SetAccountDomain safely sets the Configuration value for state's 'AccountDomain' field
func (st *ConfigState) SetAccountDomain(v string) {
	st.config.AccountDomain = v
	st.reloadToViper()
}

// GetAccountDomain safely fetches the value for global configuration 'AccountDomain' field
func GetAccountDomain() string { return global.GetAccountDomain() }

// SetAccountDomain safely sets the value for global configuration 'AccountDomain' field
func SetAccountDomain(v string) { global.SetAccountDomain(v) }

// GetProtocol safely fetches the Configuration value for state's 'Protocol' field
func (st *ConfigState) GetProtocol() (v string) {
	return st.config.Protocol
}

// SetProtocol safely sets the Configuration value for state's 'Protocol' field
func (st *ConfigState) SetProtocol(v string) {
	st.config.Protocol = v
	st.reloadToViper()
}

// GetProtocol safely fetches the value for global configuration 'Protocol' field
func GetProtocol() string { return global.GetProtocol() }

// SetProtocol safely sets the value for global configuration 'Protocol' field
func SetProtocol(v string) { global.SetProtocol(v) }

// GetBindAddress safely fetches the Configuration value for state's 'BindAddress' field
func (st *ConfigState) GetBindAddress() (v string) {
	return st.config.BindAddress
}

// SetBindAddress safely sets the Configuration value for state's 'BindAddress' field
func (st *ConfigState) SetBindAddress(v string) {
	st.config.BindAddress = v
	st.reloadToViper()
}

// GetBindAddress safely fetches the value for global configuration 'BindAddress' field
func GetBindAddress() string { return global.GetBindAddress() }

// SetBindAddress safely sets the value for global configuration 'BindAddress' field
func SetBindAddress(v string) { global.SetBindAddress(v) }

// GetPort safely fetches the Configuration value for state's 'Port' field
func (st *ConfigState) GetPort() (v int) {
	return st.config.Port
}

// SetPort safely sets the Configuration value for state's 'Port' field
func (st *ConfigState) SetPort(v int) {
	st.config.Port = v
	st.reloadToViper()
}

// GetPort safely fetches the value for global configuration 'Port' field
func GetPort() int { return global.GetPort() }

// SetPort safely sets the value for global configuration 'Port' field
func SetPort(v int) { global.SetPort(v) }

// GetTrustedProxies safely fetches the Configuration value for state's 'TrustedProxies' field
func (st *ConfigState) GetTrustedProxies() (v IPPrefixes) {
	return st.config.TrustedProxies
}

// SetTrustedProxies safely sets the Configuration value for state's 'TrustedProxies' field
func (st *ConfigState) SetTrustedProxies(v IPPrefixes) {
	st.config.TrustedProxies = v
	st.reloadToViper()
}

// GetTrustedProxies safely fetches the value for global configuration 'TrustedProxies' field
func GetTrustedProxies() IPPrefixes { return global.GetTrustedProxies() }

// SetTrustedProxies safely sets the value for global configuration 'TrustedProxies' field
func SetTrustedProxies(v IPPrefixes) { global.SetTrustedProxies(v) }

// GetSoftwareVersion safely fetches the Configuration value for state's 'SoftwareVersion' field
func (st *ConfigState) GetSoftwareVersion() (v string) {
	return st.config.SoftwareVersion
}

// SetSoftwareVersion safely sets the Configuration value for state's 'SoftwareVersion' field
func (st *ConfigState) SetSoftwareVersion(v string) {
	st.config.SoftwareVersion = v
	st.reloadToViper()
}

// GetSoftwareVersion safely fetches the value for global configuration 'SoftwareVersion' field
func GetSoftwareVersion() string { return global.GetSoftwareVersion() }

// SetSoftwareVersion safely sets the value for global configuration 'SoftwareVersion' field
func SetSoftwareVersion(v string) { global.SetSoftwareVersion(v) }

// GetWebTemplateBaseDir safely fetches the Configuration value for state's 'WebTemplateBaseDir' field
func (st *ConfigState) GetWebTemplateBaseDir() (v string) {
	return st.config.WebTemplateBaseDir
}

// SetWebTemplateBaseDir safely sets the Configuration value for state's 'WebTemplateBaseDir' field
func (st *ConfigState) SetWebTemplateBaseDir(v string) {
	st.config.WebTemplateBaseDir = v
	st.reloadToViper()
}

// GetWebTemplateBaseDir safely fetches the value for global configuration 'WebTemplateBaseDir' field
func GetWebTemplateBaseDir() string { return global.GetWebTemplateBaseDir() }

// SetWebTemplateBaseDir safely sets the value for global configuration 'WebTemplateBaseDir' field
func SetWebTemplateBaseDir(v string) { global.SetWebTemplateBaseDir(v) }

// GetWebAssetBaseDir safely fetches the Configuration value for state's 'WebAssetBaseDir' field
func (st *ConfigState) GetWebAssetBaseDir() (v string) {
	return st.config.WebAssetBaseDir
}

// SetWebAssetBaseDir safely sets the Configuration value for state's 'WebAssetBaseDir' field
func (st *ConfigState) SetWebAssetBaseDir(v string) {
	st.config.WebAssetBaseDir = v
	st.reloadToViper()
}

// GetWebAssetBaseDir safely fetches the value for global configuration 'WebAssetBaseDir' field
func GetWebAssetBaseDir() string { return global.GetWebAssetBaseDir() }

// SetWebAssetBaseDir safely sets the value for global configuration 'WebAssetBaseDir' field
func SetWebAssetBaseDir(v string) { global.SetWebAssetBaseDir(v) }

// GetInstanceFederationMode safely fetches the Configuration value for state's 'InstanceFederationMode' field
func (st *ConfigState) GetInstanceFederationMode() (v string) {
	return st.config.InstanceFederationMode
}

// SetInstanceFederationMode safely sets the Configuration value for state's 'InstanceFederationMode' field
func (st *ConfigState) SetInstanceFederationMode(v string) {
	st.config.InstanceFederationMode = v
	st.reloadToViper()
}

// GetInstanceFederationMode safely fetches the value for global configuration 'InstanceFederationMode' field
func GetInstanceFederationMode() string { return global.GetInstanceFederationMode() }

// SetInstanceFederationMode safely sets the value for global configuration 'InstanceFederationMode' field
func SetInstanceFederationMode(v string) { global.SetInstanceFederationMode(v) }

// GetInstanceFederationSpamFilter safely fetches the Configuration value for state's 'InstanceFederationSpamFilter' field
func (st *ConfigState) GetInstanceFederationSpamFilter() (v bool) {
	return st.config.InstanceFederationSpamFilter
}

// SetInstanceFederationSpamFilter safely sets the Configuration value for state's 'InstanceFederationSpamFilter' field
func (st *ConfigState) SetInstanceFederationSpamFilter(v bool) {
	st.config.InstanceFederationSpamFilter = v
	st.reloadToViper()
}

// GetInstanceFederationSpamFilter safely fetches the value for global configuration 'InstanceFederationSpamFilter' field
func GetInstanceFederationSpamFilter() bool { return global.GetInstanceFederationSpamFilter() }

// SetInstanceFederationSpamFilter safely sets the value for global configuration 'InstanceFederationSpamFilter' field
func SetInstanceFederationSpamFilter(v bool) { global.SetInstanceFederationSpamFilter(v) }

// GetInstanceExposePeers safely fetches the Configuration value for state's 'InstanceExposePeers' field
func (st *ConfigState) GetInstanceExposePeers() (v bool) {
	return st.config.InstanceExposePeers
}

// SetInstanceExposePeers safely sets the Configuration value for state's 'InstanceExposePeers' field
func (st *ConfigState) SetInstanceExposePeers(v bool) {
	st.config.InstanceExposePeers = v
	st.reloadToViper()
}

// GetInstanceExposePeers safely fetches the value for global configuration 'InstanceExposePeers' field
func GetInstanceExposePeers() bool { return global.GetInstanceExposePeers() }

// SetInstanceExposePeers safely sets the value for global configuration 'InstanceExposePeers' field
func SetInstanceExposePeers(v bool) { global.SetInstanceExposePeers(v) }

// GetInstanceExposeBlocklist safely fetches the Configuration value for state's 'InstanceExposeBlocklist' field
func (st *ConfigState) GetInstanceExposeBlocklist() (v bool) {
	return st.config.InstanceExposeBlocklist
}

// SetInstanceExposeBlocklist safely sets the Configuration value for state's 'InstanceExposeBlocklist' field
func (st *ConfigState) SetInstanceExposeBlocklist(v bool) {
	st.config.InstanceExposeBlocklist = v
	st.reloadToViper()
}

// GetInstanceExposeBlocklist safely fetches the value for global configuration 'InstanceExposeBlocklist' field
func GetInstanceExposeBlocklist() bool { return global.GetInstanceExposeBlocklist() }

// SetInstanceExposeBlocklist safely sets the value for global configuration 'InstanceExposeBlocklist' field
func SetInstanceExposeBlocklist(v bool) { global.SetInstanceExposeBlocklist(v) }

// GetInstanceExposeBlocklistWeb safely fetches the Configuration value for state's 'InstanceExposeBlocklistWeb' field
func (st *ConfigState) GetInstanceExposeBlocklistWeb() (v bool) {
	return st.config.InstanceExposeBlocklistWeb
}

// SetInstanceExposeBlocklistWeb safely sets the Configuration value for state's 'InstanceExposeBlocklistWeb' field
func (st *ConfigState) SetInstanceExposeBlocklistWeb(v bool) {
	st.config.InstanceExposeBlocklistWeb = v
	st.reloadToViper()
}

// GetInstanceExposeBlocklistWeb safely fetches the value for global configuration 'InstanceExposeBlocklistWeb' field
func GetInstanceExposeBlocklistWeb() bool { return global.GetInstanceExposeBlocklistWeb() }

// SetInstanceExposeBlocklistWeb safely sets the value for global configuration 'InstanceExposeBlocklistWeb' field
func SetInstanceExposeBlocklistWeb(v bool) { global.SetInstanceExposeBlocklistWeb(v) }

// GetInstanceExposeAllowlist safely fetches the Configuration value for state's 'InstanceExposeAllowlist' field
func (st *ConfigState) GetInstanceExposeAllowlist() (v bool) {
	return st.config.InstanceExposeAllowlist
}

// SetInstanceExposeAllowlist safely sets the Configuration value for state's 'InstanceExposeAllowlist' field
func (st *ConfigState) SetInstanceExposeAllowlist(v bool) {
	st.config.InstanceExposeAllowlist = v
	st.reloadToViper()
}

// GetInstanceExposeAllowlist safely fetches the value for global configuration 'InstanceExposeAllowlist' field
func GetInstanceExposeAllowlist() bool { return global.GetInstanceExposeAllowlist() }

// SetInstanceExposeAllowlist safely sets the value for global configuration 'InstanceExposeAllowlist' field
func SetInstanceExposeAllowlist(v bool) { global.SetInstanceExposeAllowlist(v) }

// GetInstanceExposeAllowlistWeb safely fetches the Configuration value for state's 'InstanceExposeAllowlistWeb' field
func (st *ConfigState) GetInstanceExposeAllowlistWeb() (v bool) {
	return st.config.InstanceExposeAllowlistWeb
}

// SetInstanceExposeAllowlistWeb safely sets the Configuration value for state's 'InstanceExposeAllowlistWeb' field
func (st *ConfigState) SetInstanceExposeAllowlistWeb(v bool) {
	st.config.InstanceExposeAllowlistWeb = v
	st.reloadToViper()
}

// GetInstanceExposeAllowlistWeb safely fetches the value for global configuration 'InstanceExposeAllowlistWeb' field
func GetInstanceExposeAllowlistWeb() bool { return global.GetInstanceExposeAllowlistWeb() }

// SetInstanceExposeAllowlistWeb safely sets the value for global configuration 'InstanceExposeAllowlistWeb' field
func SetInstanceExposeAllowlistWeb(v bool) { global.SetInstanceExposeAllowlistWeb(v) }

// GetInstanceExposePublicTimeline safely fetches the Configuration value for state's 'InstanceExposePublicTimeline' field
func (st *ConfigState) GetInstanceExposePublicTimeline() (v bool) {
	return st.config.InstanceExposePublicTimeline
}

// SetInstanceExposePublicTimeline safely sets the Configuration value for state's 'InstanceExposePublicTimeline' field
func (st *ConfigState) SetInstanceExposePublicTimeline(v bool) {
	st.config.InstanceExposePublicTimeline = v
	st.reloadToViper()
}

// GetInstanceExposePublicTimeline safely fetches the value for global configuration 'InstanceExposePublicTimeline' field
func GetInstanceExposePublicTimeline() bool { return global.GetInstanceExposePublicTimeline() }

// SetInstanceExposePublicTimeline safely sets the value for global configuration 'InstanceExposePublicTimeline' field
func SetInstanceExposePublicTimeline(v bool) { global.SetInstanceExposePublicTimeline(v) }

// GetInstanceExposeCustomEmojis safely fetches the Configuration value for state's 'InstanceExposeCustomEmojis' field
func (st *ConfigState) GetInstanceExposeCustomEmojis() (v bool) {
	return st.config.InstanceExposeCustomEmojis
}

// SetInstanceExposeCustomEmojis safely sets the Configuration value for state's 'InstanceExposeCustomEmojis' field
func (st *ConfigState) SetInstanceExposeCustomEmojis(v bool) {
	st.config.InstanceExposeCustomEmojis = v
	st.reloadToViper()
}

// GetInstanceExposeCustomEmojis safely fetches the value for global configuration 'InstanceExposeCustomEmojis' field
func GetInstanceExposeCustomEmojis() bool { return global.GetInstanceExposeCustomEmojis() }

// SetInstanceExposeCustomEmojis safely sets the value for global configuration 'InstanceExposeCustomEmojis' field
func SetInstanceExposeCustomEmojis(v bool) { global.SetInstanceExposeCustomEmojis(v) }

// GetInstanceDirectoryMode safely fetches the Configuration value for state's 'InstanceDirectoryMode' field
func (st *ConfigState) GetInstanceDirectoryMode() (v InstanceDirectoryMode) {
	return st.config.InstanceDirectoryMode
}

// SetInstanceDirectoryMode safely sets the Configuration value for state's 'InstanceDirectoryMode' field
func (st *ConfigState) SetInstanceDirectoryMode(v InstanceDirectoryMode) {
	st.config.InstanceDirectoryMode = v
	st.reloadToViper()
}

// GetInstanceDirectoryMode safely fetches the value for global configuration 'InstanceDirectoryMode' field
func GetInstanceDirectoryMode() InstanceDirectoryMode { return global.GetInstanceDirectoryMode() }

// SetInstanceDirectoryMode safely sets the value for global configuration 'InstanceDirectoryMode' field
func SetInstanceDirectoryMode(v InstanceDirectoryMode) { global.SetInstanceDirectoryMode(v) }

// GetInstanceDeliverToSharedInboxes safely fetches the Configuration value for state's 'InstanceDeliverToSharedInboxes' field
func (st *ConfigState) GetInstanceDeliverToSharedInboxes() (v bool) {
	return st.config.InstanceDeliverToSharedInboxes
}

// SetInstanceDeliverToSharedInboxes safely sets the Configuration value for state's 'InstanceDeliverToSharedInboxes' field
func (st *ConfigState) SetInstanceDeliverToSharedInboxes(v bool) {
	st.config.InstanceDeliverToSharedInboxes = v
	st.reloadToViper()
}

// GetInstanceDeliverToSharedInboxes safely fetches the value for global configuration 'InstanceDeliverToSharedInboxes' field
func GetInstanceDeliverToSharedInboxes() bool { return global.GetInstanceDeliverToSharedInboxes() }

// SetInstanceDeliverToSharedInboxes safely sets the value for global configuration 'InstanceDeliverToSharedInboxes' field
func SetInstanceDeliverToSharedInboxes(v bool) { global.SetInstanceDeliverToSharedInboxes(v) }

// GetInstanceInjectMastodonVersion safely fetches the Configuration value for state's 'InstanceInjectMastodonVersion' field
func (st *ConfigState) GetInstanceInjectMastodonVersion() (v bool) {
	return st.config.InstanceInjectMastodonVersion
}

// SetInstanceInjectMastodonVersion safely sets the Configuration value for state's 'InstanceInjectMastodonVersion' field
func (st *ConfigState) SetInstanceInjectMastodonVersion(v bool) {
	st.config.InstanceInjectMastodonVersion = v
	st.reloadToViper()
}

// GetInstanceInjectMastodonVersion safely fetches the value for global configuration 'InstanceInjectMastodonVersion' field
func GetInstanceInjectMastodonVersion() bool { return global.GetInstanceInjectMastodonVersion() }

// SetInstanceInjectMastodonVersion safely sets the value for global configuration 'InstanceInjectMastodonVersion' field
func SetInstanceInjectMastodonVersion(v bool) { global.SetInstanceInjectMastodonVersion(v) }

// GetInstanceLanguages safely fetches the Configuration value for state's 'InstanceLanguages' field
func (st *ConfigState) GetInstanceLanguages() (v language.Languages) {
	return st.config.InstanceLanguages
}

// SetInstanceLanguages safely sets the Configuration value for state's 'InstanceLanguages' field
func (st *ConfigState) SetInstanceLanguages(v language.Languages) {
	st.config.InstanceLanguages = v
	st.reloadToViper()
}

// GetInstanceLanguages safely fetches the value for global configuration 'InstanceLanguages' field
func GetInstanceLanguages() language.Languages { return global.GetInstanceLanguages() }

// SetInstanceLanguages safely sets the value for global configuration 'InstanceLanguages' field
func SetInstanceLanguages(v language.Languages) { global.SetInstanceLanguages(v) }

// GetInstanceStatsMode safely fetches the Configuration value for state's 'InstanceStatsMode' field
func (st *ConfigState) GetInstanceStatsMode() (v string) {
	return st.config.InstanceStatsMode
}

// SetInstanceStatsMode safely sets the Configuration value for state's 'InstanceStatsMode' field
func (st *ConfigState) SetInstanceStatsMode(v string) {
	st.config.InstanceStatsMode = v
	st.reloadToViper()
}

// GetInstanceStatsMode safely fetches the value for global configuration 'InstanceStatsMode' field
func GetInstanceStatsMode() string { return global.GetInstanceStatsMode() }

// SetInstanceStatsMode safely sets the value for global configuration 'InstanceStatsMode' field
func SetInstanceStatsMode(v string) { global.SetInstanceStatsMode(v) }

// GetInstanceAllowBackdatingStatuses safely fetches the Configuration value for state's 'InstanceAllowBackdatingStatuses' field
func (st *ConfigState) GetInstanceAllowBackdatingStatuses() (v bool) {
	return st.config.InstanceAllowBackdatingStatuses
}

// SetInstanceAllowBackdatingStatuses safely sets the Configuration value for state's 'InstanceAllowBackdatingStatuses' field
func (st *ConfigState) SetInstanceAllowBackdatingStatuses(v bool) {
	st.config.InstanceAllowBackdatingStatuses = v
	st.reloadToViper()
}

// GetInstanceAllowBackdatingStatuses safely fetches the value for global configuration 'InstanceAllowBackdatingStatuses' field
func GetInstanceAllowBackdatingStatuses() bool { return global.GetInstanceAllowBackdatingStatuses() }

// SetInstanceAllowBackdatingStatuses safely sets the value for global configuration 'InstanceAllowBackdatingStatuses' field
func SetInstanceAllowBackdatingStatuses(v bool) { global.SetInstanceAllowBackdatingStatuses(v) }

// GetInstanceRobotsAllowIndexing safely fetches the Configuration value for state's 'InstanceRobotsAllowIndexing' field
func (st *ConfigState) GetInstanceRobotsAllowIndexing() (v bool) {
	return st.config.InstanceRobotsAllowIndexing
}

// SetInstanceRobotsAllowIndexing safely sets the Configuration value for state's 'InstanceRobotsAllowIndexing' field
func (st *ConfigState) SetInstanceRobotsAllowIndexing(v bool) {
	st.config.InstanceRobotsAllowIndexing = v
	st.reloadToViper()
}

// GetInstanceRobotsAllowIndexing safely fetches the value for global configuration 'InstanceRobotsAllowIndexing' field
func GetInstanceRobotsAllowIndexing() bool { return global.GetInstanceRobotsAllowIndexing() }

// SetInstanceRobotsAllowIndexing safely sets the value for global configuration 'InstanceRobotsAllowIndexing' field
func SetInstanceRobotsAllowIndexing(v bool) { global.SetInstanceRobotsAllowIndexing(v) }

// GetAccountsRegistrationOpen safely fetches the Configuration value for state's 'AccountsRegistrationOpen' field
func (st *ConfigState) GetAccountsRegistrationOpen() (v bool) {
	return st.config.AccountsRegistrationOpen
}

// SetAccountsRegistrationOpen safely sets the Configuration value for state's 'AccountsRegistrationOpen' field
func (st *ConfigState) SetAccountsRegistrationOpen(v bool) {
	st.config.AccountsRegistrationOpen = v
	st.reloadToViper()
}

// GetAccountsRegistrationOpen safely fetches the value for global configuration 'AccountsRegistrationOpen' field
func GetAccountsRegistrationOpen() bool { return global.GetAccountsRegistrationOpen() }

// SetAccountsRegistrationOpen safely sets the value for global configuration 'AccountsRegistrationOpen' field
func SetAccountsRegistrationOpen(v bool) { global.SetAccountsRegistrationOpen(v) }

// GetAccountsReasonRequired safely fetches the Configuration value for state's 'AccountsReasonRequired' field
func (st *ConfigState) GetAccountsReasonRequired() (v bool) {
	return st.config.AccountsReasonRequired
}

// SetAccountsReasonRequired safely sets the Configuration value for state's 'AccountsReasonRequired' field
func (st *ConfigState) SetAccountsReasonRequired(v bool) {
	st.config.AccountsReasonRequired = v
	st.reloadToViper()
}

// GetAccountsReasonRequired safely fetches the value for global configuration 'AccountsReasonRequired' field
func GetAccountsReasonRequired() bool { return global.GetAccountsReasonRequired() }

// SetAccountsReasonRequired safely sets the value for global configuration 'AccountsReasonRequired' field
func SetAccountsReasonRequired(v bool) { global.SetAccountsReasonRequired(v) }

// GetAccountsRegistrationDailyLimit safely fetches the Configuration value for state's 'AccountsRegistrationDailyLimit' field
func (st *ConfigState) GetAccountsRegistrationDailyLimit() (v int) {
	return st.config.AccountsRegistrationDailyLimit
}

// SetAccountsRegistrationDailyLimit safely sets the Configuration value for state's 'AccountsRegistrationDailyLimit' field
func (st *ConfigState) SetAccountsRegistrationDailyLimit(v int) {
	st.config.AccountsRegistrationDailyLimit = v
	st.reloadToViper()
}

// GetAccountsRegistrationDailyLimit safely fetches the value for global configuration 'AccountsRegistrationDailyLimit' field
func GetAccountsRegistrationDailyLimit() int { return global.GetAccountsRegistrationDailyLimit() }

// SetAccountsRegistrationDailyLimit safely sets the value for global configuration 'AccountsRegistrationDailyLimit' field
func SetAccountsRegistrationDailyLimit(v int) { global.SetAccountsRegistrationDailyLimit(v) }

// GetAccountsRegistrationBacklogLimit safely fetches the Configuration value for state's 'AccountsRegistrationBacklogLimit' field
func (st *ConfigState) GetAccountsRegistrationBacklogLimit() (v int) {
	return st.config.AccountsRegistrationBacklogLimit
}

// SetAccountsRegistrationBacklogLimit safely sets the Configuration value for state's 'AccountsRegistrationBacklogLimit' field
func (st *ConfigState) SetAccountsRegistrationBacklogLimit(v int) {
	st.config.AccountsRegistrationBacklogLimit = v
	st.reloadToViper()
}

// GetAccountsRegistrationBacklogLimit safely fetches the value for global configuration 'AccountsRegistrationBacklogLimit' field
func GetAccountsRegistrationBacklogLimit() int { return global.GetAccountsRegistrationBacklogLimit() }

// SetAccountsRegistrationBacklogLimit safely sets the value for global configuration 'AccountsRegistrationBacklogLimit' field
func SetAccountsRegistrationBacklogLimit(v int) { global.SetAccountsRegistrationBacklogLimit(v) }

// GetAccountsAllowCustomCSS safely fetches the Configuration value for state's 'AccountsAllowCustomCSS' field
func (st *ConfigState) GetAccountsAllowCustomCSS() (v bool) {
	return st.config.AccountsAllowCustomCSS
}

// SetAccountsAllowCustomCSS safely sets the Configuration value for state's 'AccountsAllowCustomCSS' field
func (st *ConfigState) SetAccountsAllowCustomCSS(v bool) {
	st.config.AccountsAllowCustomCSS = v
	st.reloadToViper()
}

// GetAccountsAllowCustomCSS safely fetches the value for global configuration 'AccountsAllowCustomCSS' field
func GetAccountsAllowCustomCSS() bool { return global.GetAccountsAllowCustomCSS() }

// SetAccountsAllowCustomCSS safely sets the value for global configuration 'AccountsAllowCustomCSS' field
func SetAccountsAllowCustomCSS(v bool) { global.SetAccountsAllowCustomCSS(v) }

// GetAccountsCustomCSSLength safely fetches the Configuration value for state's 'AccountsCustomCSSLength' field
func (st *ConfigState) GetAccountsCustomCSSLength() (v int) {
	return st.config.AccountsCustomCSSLength
}

// SetAccountsCustomCSSLength safely sets the Configuration value for state's 'AccountsCustomCSSLength' field
func (st *ConfigState) SetAccountsCustomCSSLength(v int) {
	st.config.AccountsCustomCSSLength = v
	st.reloadToViper()
}

// GetAccountsCustomCSSLength safely fetches the value for global configuration 'AccountsCustomCSSLength' field
func GetAccountsCustomCSSLength() int { return global.GetAccountsCustomCSSLength() }

// SetAccountsCustomCSSLength safely sets the value for global configuration 'AccountsCustomCSSLength' field
func SetAccountsCustomCSSLength(v int) { global.SetAccountsCustomCSSLength(v) }

// GetAccountsMaxProfileFields safely fetches the Configuration value for state's 'AccountsMaxProfileFields' field
func (st *ConfigState) GetAccountsMaxProfileFields() (v int) {
	return st.config.AccountsMaxProfileFields
}

// SetAccountsMaxProfileFields safely sets the Configuration value for state's 'AccountsMaxProfileFields' field
func (st *ConfigState) SetAccountsMaxProfileFields(v int) {
	st.config.AccountsMaxProfileFields = v
	st.reloadToViper()
}

// GetAccountsMaxProfileFields safely fetches the value for global configuration 'AccountsMaxProfileFields' field
func GetAccountsMaxProfileFields() int { return global.GetAccountsMaxProfileFields() }

// SetAccountsMaxProfileFields safely sets the value for global configuration 'AccountsMaxProfileFields' field
func SetAccountsMaxProfileFields(v int) { global.SetAccountsMaxProfileFields(v) }

// GetStorageBackend safely fetches the Configuration value for state's 'StorageBackend' field
func (st *ConfigState) GetStorageBackend() (v string) {
	return st.config.StorageBackend
}

// SetStorageBackend safely sets the Configuration value for state's 'StorageBackend' field
func (st *ConfigState) SetStorageBackend(v string) {
	st.config.StorageBackend = v
	st.reloadToViper()
}

// GetStorageBackend safely fetches the value for global configuration 'StorageBackend' field
func GetStorageBackend() string { return global.GetStorageBackend() }

// SetStorageBackend safely sets the value for global configuration 'StorageBackend' field
func SetStorageBackend(v string) { global.SetStorageBackend(v) }

// GetStorageLocalBasePath safely fetches the Configuration value for state's 'StorageLocalBasePath' field
func (st *ConfigState) GetStorageLocalBasePath() (v string) {
	return st.config.StorageLocalBasePath
}

// SetStorageLocalBasePath safely sets the Configuration value for state's 'StorageLocalBasePath' field
func (st *ConfigState) SetStorageLocalBasePath(v string) {
	st.config.StorageLocalBasePath = v
	st.reloadToViper()
}

// GetStorageLocalBasePath safely fetches the value for global configuration 'StorageLocalBasePath' field
func GetStorageLocalBasePath() string { return global.GetStorageLocalBasePath() }

// SetStorageLocalBasePath safely sets the value for global configuration 'StorageLocalBasePath' field
func SetStorageLocalBasePath(v string) { global.SetStorageLocalBasePath(v) }

// GetStorageS3Endpoint safely fetches the Configuration value for state's 'StorageS3Endpoint' field
func (st *ConfigState) GetStorageS3Endpoint() (v string) {
	return st.config.StorageS3Endpoint
}

// SetStorageS3Endpoint safely sets the Configuration value for state's 'StorageS3Endpoint' field
func (st *ConfigState) SetStorageS3Endpoint(v string) {
	st.config.StorageS3Endpoint = v
	st.reloadToViper()
}

// GetStorageS3Endpoint safely fetches the value for global configuration 'StorageS3Endpoint' field
func GetStorageS3Endpoint() string { return global.GetStorageS3Endpoint() }

// SetStorageS3Endpoint safely sets the value for global configuration 'StorageS3Endpoint' field
func SetStorageS3Endpoint(v string) { global.SetStorageS3Endpoint(v) }

// GetStorageS3AccessKey safely fetches the Configuration value for state's 'StorageS3AccessKey' field
func (st *ConfigState) GetStorageS3AccessKey() (v string) {
	return st.config.StorageS3AccessKey
}

// SetStorageS3AccessKey safely sets the Configuration value for state's 'StorageS3AccessKey' field
func (st *ConfigState) SetStorageS3AccessKey(v string) {
	st.config.StorageS3AccessKey = v
	st.reloadToViper()
}

// GetStorageS3AccessKey safely fetches the value for global configuration 'StorageS3AccessKey' field
func GetStorageS3AccessKey() string { return global.GetStorageS3AccessKey() }

// SetStorageS3AccessKey safely sets the value for global configuration 'StorageS3AccessKey' field
func SetStorageS3AccessKey(v string) { global.SetStorageS3AccessKey(v) }

// GetStorageS3SecretKey safely fetches the Configuration value for state's 'StorageS3SecretKey' field
func (st *ConfigState) GetStorageS3SecretKey() (v string) {
	return st.config.StorageS3SecretKey
}

// SetStorageS3SecretKey safely sets the Configuration value for state's 'StorageS3SecretKey' field
func (st *ConfigState) SetStorageS3SecretKey(v string) {
	st.config.StorageS3SecretKey = v
	st.reloadToViper()
}

// GetStorageS3SecretKey safely fetches the value for global configuration 'StorageS3SecretKey' field
func GetStorageS3SecretKey() string { return global.GetStorageS3SecretKey() }

// SetStorageS3SecretKey safely sets the value for global configuration 'StorageS3SecretKey' field
func SetStorageS3SecretKey(v string) { global.SetStorageS3SecretKey(v) }

// GetStorageS3UseSSL safely fetches the Configuration value for state's 'StorageS3UseSSL' field
func (st *ConfigState) GetStorageS3UseSSL() (v bool) {
	return st.config.StorageS3UseSSL
}

// SetStorageS3UseSSL safely sets the Configuration value for state's 'StorageS3UseSSL' field
func (st *ConfigState) SetStorageS3UseSSL(v bool) {
	st.config.StorageS3UseSSL = v
	st.reloadToViper()
}

// GetStorageS3UseSSL safely fetches the value for global configuration 'StorageS3UseSSL' field
func GetStorageS3UseSSL() bool { return global.GetStorageS3UseSSL() }

// SetStorageS3UseSSL safely sets the value for global configuration 'StorageS3UseSSL' field
func SetStorageS3UseSSL(v bool) { global.SetStorageS3UseSSL(v) }

// GetStorageS3BucketName safely fetches the Configuration value for state's 'StorageS3BucketName' field
func (st *ConfigState) GetStorageS3BucketName() (v string) {
	return st.config.StorageS3BucketName
}

// SetStorageS3BucketName safely sets the Configuration value for state's 'StorageS3BucketName' field
func (st *ConfigState) SetStorageS3BucketName(v string) {
	st.config.StorageS3BucketName = v
	st.reloadToViper()
}

// GetStorageS3BucketName safely fetches the value for global configuration 'StorageS3BucketName' field
func GetStorageS3BucketName() string { return global.GetStorageS3BucketName() }

// SetStorageS3BucketName safely sets the value for global configuration 'StorageS3BucketName' field
func SetStorageS3BucketName(v string) { global.SetStorageS3BucketName(v) }

// GetStorageS3Proxy safely fetches the Configuration value for state's 'StorageS3Proxy' field
func (st *ConfigState) GetStorageS3Proxy() (v bool) {
	return st.config.StorageS3Proxy
}

// SetStorageS3Proxy safely sets the Configuration value for state's 'StorageS3Proxy' field
func (st *ConfigState) SetStorageS3Proxy(v bool) {
	st.config.StorageS3Proxy = v
	st.reloadToViper()
}

// GetStorageS3Proxy safely fetches the value for global configuration 'StorageS3Proxy' field
func GetStorageS3Proxy() bool { return global.GetStorageS3Proxy() }

// SetStorageS3Proxy safely sets the value for global configuration 'StorageS3Proxy' field
func SetStorageS3Proxy(v bool) { global.SetStorageS3Proxy(v) }

// GetStorageS3RedirectURL safely fetches the Configuration value for state's 'StorageS3RedirectURL' field
func (st *ConfigState) GetStorageS3RedirectURL() (v string) {
	return st.config.StorageS3RedirectURL
}

// SetStorageS3RedirectURL safely sets the Configuration value for state's 'StorageS3RedirectURL' field
func (st *ConfigState) SetStorageS3RedirectURL(v string) {
	st.config.StorageS3RedirectURL = v
	st.reloadToViper()
}

// GetStorageS3RedirectURL safely fetches the value for global configuration 'StorageS3RedirectURL' field
func GetStorageS3RedirectURL() string { return global.GetStorageS3RedirectURL() }

// SetStorageS3RedirectURL safely sets the value for global configuration 'StorageS3RedirectURL' field
func SetStorageS3RedirectURL(v string) { global.SetStorageS3RedirectURL(v) }

// GetStorageS3BucketLookup safely fetches the Configuration value for state's 'StorageS3BucketLookup' field
func (st *ConfigState) GetStorageS3BucketLookup() (v string) {
	return st.config.StorageS3BucketLookup
}

// SetStorageS3BucketLookup safely sets the Configuration value for state's 'StorageS3BucketLookup' field
func (st *ConfigState) SetStorageS3BucketLookup(v string) {
	st.config.StorageS3BucketLookup = v
	st.reloadToViper()
}

// GetStorageS3BucketLookup safely fetches the value for global configuration 'StorageS3BucketLookup' field
func GetStorageS3BucketLookup() string { return global.GetStorageS3BucketLookup() }

// SetStorageS3BucketLookup safely sets the value for global configuration 'StorageS3BucketLookup' field
func SetStorageS3BucketLookup(v string) { global.SetStorageS3BucketLookup(v) }

// GetStorageS3KeyPrefix safely fetches the Configuration value for state's 'StorageS3KeyPrefix' field
func (st *ConfigState) GetStorageS3KeyPrefix() (v string) {
	return st.config.StorageS3KeyPrefix
}

// SetStorageS3KeyPrefix safely sets the Configuration value for state's 'StorageS3KeyPrefix' field
func (st *ConfigState) SetStorageS3KeyPrefix(v string) {
	st.config.StorageS3KeyPrefix = v
	st.reloadToViper()
}

// GetStorageS3KeyPrefix safely fetches the value for global configuration 'StorageS3KeyPrefix' field
func GetStorageS3KeyPrefix() string { return global.GetStorageS3KeyPrefix() }

// SetStorageS3KeyPrefix safely sets the value for global configuration 'StorageS3KeyPrefix' field
func SetStorageS3KeyPrefix(v string) { global.SetStorageS3KeyPrefix(v) }

// GetStorageS3Region safely fetches the Configuration value for state's 'StorageS3Region' field
func (st *ConfigState) GetStorageS3Region() (v string) {
	return st.config.StorageS3Region
}

// SetStorageS3Region safely sets the Configuration value for state's 'StorageS3Region' field
func (st *ConfigState) SetStorageS3Region(v string) {
	st.config.StorageS3Region = v
	st.reloadToViper()
}

// GetStorageS3Region safely fetches the value for global configuration 'StorageS3Region' field
func GetStorageS3Region() string { return global.GetStorageS3Region() }

// SetStorageS3Region safely sets the value for global configuration 'StorageS3Region' field
func SetStorageS3Region(v string) { global.SetStorageS3Region(v) }

// GetStatusesMaxChars safely fetches the Configuration value for state's 'StatusesMaxChars' field
func (st *ConfigState) GetStatusesMaxChars() (v int) {
	return st.config.StatusesMaxChars
}

// SetStatusesMaxChars safely sets the Configuration value for state's 'StatusesMaxChars' field
func (st *ConfigState) SetStatusesMaxChars(v int) {
	st.config.StatusesMaxChars = v
	st.reloadToViper()
}

// GetStatusesMaxChars safely fetches the value for global configuration 'StatusesMaxChars' field
func GetStatusesMaxChars() int { return global.GetStatusesMaxChars() }

// SetStatusesMaxChars safely sets the value for global configuration 'StatusesMaxChars' field
func SetStatusesMaxChars(v int) { global.SetStatusesMaxChars(v) }

// GetStatusesPollMaxOptions safely fetches the Configuration value for state's 'StatusesPollMaxOptions' field
func (st *ConfigState) GetStatusesPollMaxOptions() (v int) {
	return st.config.StatusesPollMaxOptions
}

// SetStatusesPollMaxOptions safely sets the Configuration value for state's 'StatusesPollMaxOptions' field
func (st *ConfigState) SetStatusesPollMaxOptions(v int) {
	st.config.StatusesPollMaxOptions = v
	st.reloadToViper()
}

// GetStatusesPollMaxOptions safely fetches the value for global configuration 'StatusesPollMaxOptions' field
func GetStatusesPollMaxOptions() int { return global.GetStatusesPollMaxOptions() }

// SetStatusesPollMaxOptions safely sets the value for global configuration 'StatusesPollMaxOptions' field
func SetStatusesPollMaxOptions(v int) { global.SetStatusesPollMaxOptions(v) }

// GetStatusesPollOptionMaxChars safely fetches the Configuration value for state's 'StatusesPollOptionMaxChars' field
func (st *ConfigState) GetStatusesPollOptionMaxChars() (v int) {
	return st.config.StatusesPollOptionMaxChars
}

// SetStatusesPollOptionMaxChars safely sets the Configuration value for state's 'StatusesPollOptionMaxChars' field
func (st *ConfigState) SetStatusesPollOptionMaxChars(v int) {
	st.config.StatusesPollOptionMaxChars = v
	st.reloadToViper()
}

// GetStatusesPollOptionMaxChars safely fetches the value for global configuration 'StatusesPollOptionMaxChars' field
func GetStatusesPollOptionMaxChars() int { return global.GetStatusesPollOptionMaxChars() }

// SetStatusesPollOptionMaxChars safely sets the value for global configuration 'StatusesPollOptionMaxChars' field
func SetStatusesPollOptionMaxChars(v int) { global.SetStatusesPollOptionMaxChars(v) }

// GetStatusesMediaMaxFiles safely fetches the Configuration value for state's 'StatusesMediaMaxFiles' field
func (st *ConfigState) GetStatusesMediaMaxFiles() (v int) {
	return st.config.StatusesMediaMaxFiles
}

// SetStatusesMediaMaxFiles safely sets the Configuration value for state's 'StatusesMediaMaxFiles' field
func (st *ConfigState) SetStatusesMediaMaxFiles(v int) {
	st.config.StatusesMediaMaxFiles = v
	st.reloadToViper()
}

// GetStatusesMediaMaxFiles safely fetches the value for global configuration 'StatusesMediaMaxFiles' field
func GetStatusesMediaMaxFiles() int { return global.GetStatusesMediaMaxFiles() }

// SetStatusesMediaMaxFiles safely sets the value for global configuration 'StatusesMediaMaxFiles' field
func SetStatusesMediaMaxFiles(v int) { global.SetStatusesMediaMaxFiles(v) }

// GetScheduledStatusesMaxTotal safely fetches the Configuration value for state's 'ScheduledStatusesMaxTotal' field
func (st *ConfigState) GetScheduledStatusesMaxTotal() (v int) {
	return st.config.ScheduledStatusesMaxTotal
}

// SetScheduledStatusesMaxTotal safely sets the Configuration value for state's 'ScheduledStatusesMaxTotal' field
func (st *ConfigState) SetScheduledStatusesMaxTotal(v int) {
	st.config.ScheduledStatusesMaxTotal = v
	st.reloadToViper()
}

// GetScheduledStatusesMaxTotal safely fetches the value for global configuration 'ScheduledStatusesMaxTotal' field
func GetScheduledStatusesMaxTotal() int { return global.GetScheduledStatusesMaxTotal() }

// SetScheduledStatusesMaxTotal safely sets the value for global configuration 'ScheduledStatusesMaxTotal' field
func SetScheduledStatusesMaxTotal(v int) { global.SetScheduledStatusesMaxTotal(v) }

// GetScheduledStatusesMaxDaily safely fetches the Configuration value for state's 'ScheduledStatusesMaxDaily' field
func (st *ConfigState) GetScheduledStatusesMaxDaily() (v int) {
	return st.config.ScheduledStatusesMaxDaily
}

// SetScheduledStatusesMaxDaily safely sets the Configuration value for state's 'ScheduledStatusesMaxDaily' field
func (st *ConfigState) SetScheduledStatusesMaxDaily(v int) {
	st.config.ScheduledStatusesMaxDaily = v
	st.reloadToViper()
}

// GetScheduledStatusesMaxDaily safely fetches the value for global configuration 'ScheduledStatusesMaxDaily' field
func GetScheduledStatusesMaxDaily() int { return global.GetScheduledStatusesMaxDaily() }

// SetScheduledStatusesMaxDaily safely sets the value for global configuration 'ScheduledStatusesMaxDaily' field
func SetScheduledStatusesMaxDaily(v int) { global.SetScheduledStatusesMaxDaily(v) }

// GetLetsEncryptEnabled safely fetches the Configuration value for state's 'LetsEncryptEnabled' field
func (st *ConfigState) GetLetsEncryptEnabled() (v bool) {
	return st.config.LetsEncryptEnabled
}

// SetLetsEncryptEnabled safely sets the Configuration value for state's 'LetsEncryptEnabled' field
func (st *ConfigState) SetLetsEncryptEnabled(v bool) {
	st.config.LetsEncryptEnabled = v
	st.reloadToViper()
}

// GetLetsEncryptEnabled safely fetches the value for global configuration 'LetsEncryptEnabled' field
func GetLetsEncryptEnabled() bool { return global.GetLetsEncryptEnabled() }

// SetLetsEncryptEnabled safely sets the value for global configuration 'LetsEncryptEnabled' field
func SetLetsEncryptEnabled(v bool) { global.SetLetsEncryptEnabled(v) }

// GetLetsEncryptPort safely fetches the Configuration value for state's 'LetsEncryptPort' field
func (st *ConfigState) GetLetsEncryptPort() (v int) {
	return st.config.LetsEncryptPort
}

// SetLetsEncryptPort safely sets the Configuration value for state's 'LetsEncryptPort' field
func (st *ConfigState) SetLetsEncryptPort(v int) {
	st.config.LetsEncryptPort = v
	st.reloadToViper()
}

// GetLetsEncryptPort safely fetches the value for global configuration 'LetsEncryptPort' field
func GetLetsEncryptPort() int { return global.GetLetsEncryptPort() }

// SetLetsEncryptPort safely sets the value for global configuration 'LetsEncryptPort' field
func SetLetsEncryptPort(v int) { global.SetLetsEncryptPort(v) }

// GetLetsEncryptCertDir safely fetches the Configuration value for state's 'LetsEncryptCertDir' field
func (st *ConfigState) GetLetsEncryptCertDir() (v string) {
	return st.config.LetsEncryptCertDir
}

// SetLetsEncryptCertDir safely sets the Configuration value for state's 'LetsEncryptCertDir' field
func (st *ConfigState) SetLetsEncryptCertDir(v string) {
	st.config.LetsEncryptCertDir = v
	st.reloadToViper()
}

// GetLetsEncryptCertDir safely fetches the value for global configuration 'LetsEncryptCertDir' field
func GetLetsEncryptCertDir() string { return global.GetLetsEncryptCertDir() }

// SetLetsEncryptCertDir safely sets the value for global configuration 'LetsEncryptCertDir' field
func SetLetsEncryptCertDir(v string) { global.SetLetsEncryptCertDir(v) }

// GetLetsEncryptEmailAddress safely fetches the Configuration value for state's 'LetsEncryptEmailAddress' field
func (st *ConfigState) GetLetsEncryptEmailAddress() (v string) {
	return st.config.LetsEncryptEmailAddress
}

// SetLetsEncryptEmailAddress safely sets the Configuration value for state's 'LetsEncryptEmailAddress' field
func (st *ConfigState) SetLetsEncryptEmailAddress(v string) {
	st.config.LetsEncryptEmailAddress = v
	st.reloadToViper()
}

// GetLetsEncryptEmailAddress safely fetches the value for global configuration 'LetsEncryptEmailAddress' field
func GetLetsEncryptEmailAddress() string { return global.GetLetsEncryptEmailAddress() }

// SetLetsEncryptEmailAddress safely sets the value for global configuration 'LetsEncryptEmailAddress' field
func SetLetsEncryptEmailAddress(v string) { global.SetLetsEncryptEmailAddress(v) }

// GetTLSCertificateChain safely fetches the Configuration value for state's 'TLSCertificateChain' field
func (st *ConfigState) GetTLSCertificateChain() (v string) {
	return st.config.TLSCertificateChain
}

// SetTLSCertificateChain safely sets the Configuration value for state's 'TLSCertificateChain' field
func (st *ConfigState) SetTLSCertificateChain(v string) {
	st.config.TLSCertificateChain = v
	st.reloadToViper()
}

// GetTLSCertificateChain safely fetches the value for global configuration 'TLSCertificateChain' field
func GetTLSCertificateChain() string { return global.GetTLSCertificateChain() }

// SetTLSCertificateChain safely sets the value for global configuration 'TLSCertificateChain' field
func SetTLSCertificateChain(v string) { global.SetTLSCertificateChain(v) }

// GetTLSCertificateKey safely fetches the Configuration value for state's 'TLSCertificateKey' field
func (st *ConfigState) GetTLSCertificateKey() (v string) {
	return st.config.TLSCertificateKey
}

// SetTLSCertificateKey safely sets the Configuration value for state's 'TLSCertificateKey' field
func (st *ConfigState) SetTLSCertificateKey(v string) {
	st.config.TLSCertificateKey = v
	st.reloadToViper()
}

// GetTLSCertificateKey safely fetches the value for global configuration 'TLSCertificateKey' field
func GetTLSCertificateKey() string { return global.GetTLSCertificateKey() }

// SetTLSCertificateKey safely sets the value for global configuration 'TLSCertificateKey' field
func SetTLSCertificateKey(v string) { global.SetTLSCertificateKey(v) }

// GetOIDCEnabled safely fetches the Configuration value for state's 'OIDCEnabled' field
func (st *ConfigState) GetOIDCEnabled() (v bool) {
	return st.config.OIDCEnabled
}

// SetOIDCEnabled safely sets the Configuration value for state's 'OIDCEnabled' field
func (st *ConfigState) SetOIDCEnabled(v bool) {
	st.config.OIDCEnabled = v
	st.reloadToViper()
}

// GetOIDCEnabled safely fetches the value for global configuration 'OIDCEnabled' field
func GetOIDCEnabled() bool { return global.GetOIDCEnabled() }

// SetOIDCEnabled safely sets the value for global configuration 'OIDCEnabled' field
func SetOIDCEnabled(v bool) { global.SetOIDCEnabled(v) }

// GetOIDCIdpName safely fetches the Configuration value for state's 'OIDCIdpName' field
func (st *ConfigState) GetOIDCIdpName() (v string) {
	return st.config.OIDCIdpName
}

// SetOIDCIdpName safely sets the Configuration value for state's 'OIDCIdpName' field
func (st *ConfigState) SetOIDCIdpName(v string) {
	st.config.OIDCIdpName = v
	st.reloadToViper()
}

// GetOIDCIdpName safely fetches the value for global configuration 'OIDCIdpName' field
func GetOIDCIdpName() string { return global.GetOIDCIdpName() }

// SetOIDCIdpName safely sets the value for global configuration 'OIDCIdpName' field
func SetOIDCIdpName(v string) { global.SetOIDCIdpName(v) }

// GetOIDCSkipVerification safely fetches the Configuration value for state's 'OIDCSkipVerification' field
func (st *ConfigState) GetOIDCSkipVerification() (v bool) {
	return st.config.OIDCSkipVerification
}

// SetOIDCSkipVerification safely sets the Configuration value for state's 'OIDCSkipVerification' field
func (st *ConfigState) SetOIDCSkipVerification(v bool) {
	st.config.OIDCSkipVerification = v
	st.reloadToViper()
}

// GetOIDCSkipVerification safely fetches the value for global configuration 'OIDCSkipVerification' field
func GetOIDCSkipVerification() bool { return global.GetOIDCSkipVerification() }

// SetOIDCSkipVerification safely sets the value for global configuration 'OIDCSkipVerification' field
func SetOIDCSkipVerification(v bool) { global.SetOIDCSkipVerification(v) }

// GetOIDCIssuer safely fetches the Configuration value for state's 'OIDCIssuer' field
func (st *ConfigState) GetOIDCIssuer() (v string) {
	return st.config.OIDCIssuer
}

// SetOIDCIssuer safely sets the Configuration value for state's 'OIDCIssuer' field
func (st *ConfigState) SetOIDCIssuer(v string) {
	st.config.OIDCIssuer = v
	st.reloadToViper()
}

// GetOIDCIssuer safely fetches the value for global configuration 'OIDCIssuer' field
func GetOIDCIssuer() string { return global.GetOIDCIssuer() }

// SetOIDCIssuer safely sets the value for global configuration 'OIDCIssuer' field
func SetOIDCIssuer(v string) { global.SetOIDCIssuer(v) }

// GetOIDCClientID safely fetches the Configuration value for state's 'OIDCClientID' field
func (st *ConfigState) GetOIDCClientID() (v string) {
	return st.config.OIDCClientID
}

// SetOIDCClientID safely sets the Configuration value for state's 'OIDCClientID' field
func (st *ConfigState) SetOIDCClientID(v string) {
	st.config.OIDCClientID = v
	st.reloadToViper()
}

// GetOIDCClientID safely fetches the value for global configuration 'OIDCClientID' field
func GetOIDCClientID() string { return global.GetOIDCClientID() }

// SetOIDCClientID safely sets the value for global configuration 'OIDCClientID' field
func SetOIDCClientID(v string) { global.SetOIDCClientID(v) }

// GetOIDCClientSecret safely fetches the Configuration value for state's 'OIDCClientSecret' field
func (st *ConfigState) GetOIDCClientSecret() (v string) {
	return st.config.OIDCClientSecret
}

// SetOIDCClientSecret safely sets the Configuration value for state's 'OIDCClientSecret' field
func (st *ConfigState) SetOIDCClientSecret(v string) {
	st.config.OIDCClientSecret = v
	st.reloadToViper()
}

// GetOIDCClientSecret safely fetches the value for global configuration 'OIDCClientSecret' field
func GetOIDCClientSecret() string { return global.GetOIDCClientSecret() }

// SetOIDCClientSecret safely sets the value for global configuration 'OIDCClientSecret' field
func SetOIDCClientSecret(v string) { global.SetOIDCClientSecret(v) }

// GetOIDCScopes safely fetches the Configuration value for state's 'OIDCScopes' field
func (st *ConfigState) GetOIDCScopes() (v []string) {
	return st.config.OIDCScopes
}

// SetOIDCScopes safely sets the Configuration value for state's 'OIDCScopes' field
func (st *ConfigState) SetOIDCScopes(v []string) {
	st.config.OIDCScopes = v
	st.reloadToViper()
}

// GetOIDCScopes safely fetches the value for global configuration 'OIDCScopes' field
func GetOIDCScopes() []string { return global.GetOIDCScopes() }

// SetOIDCScopes safely sets the value for global configuration 'OIDCScopes' field
func SetOIDCScopes(v []string) { global.SetOIDCScopes(v) }

// GetOIDCLinkExisting safely fetches the Configuration value for state's 'OIDCLinkExisting' field
func (st *ConfigState) GetOIDCLinkExisting() (v bool) {
	return st.config.OIDCLinkExisting
}

// SetOIDCLinkExisting safely sets the Configuration value for state's 'OIDCLinkExisting' field
func (st *ConfigState) SetOIDCLinkExisting(v bool) {
	st.config.OIDCLinkExisting = v
	st.reloadToViper()
}

// GetOIDCLinkExisting safely fetches the value for global configuration 'OIDCLinkExisting' field
func GetOIDCLinkExisting() bool { return global.GetOIDCLinkExisting() }

// SetOIDCLinkExisting safely sets the value for global configuration 'OIDCLinkExisting' field
func SetOIDCLinkExisting(v bool) { global.SetOIDCLinkExisting(v) }

// GetOIDCAllowedGroups safely fetches the Configuration value for state's 'OIDCAllowedGroups' field
func (st *ConfigState) GetOIDCAllowedGroups() (v []string) {
	return st.config.OIDCAllowedGroups
}

// SetOIDCAllowedGroups safely sets the Configuration value for state's 'OIDCAllowedGroups' field
func (st *ConfigState) SetOIDCAllowedGroups(v []string) {
	st.config.OIDCAllowedGroups = v
	st.reloadToViper()
}

// GetOIDCAllowedGroups safely fetches the value for global configuration 'OIDCAllowedGroups' field
func GetOIDCAllowedGroups() []string { return global.GetOIDCAllowedGroups() }

// SetOIDCAllowedGroups safely sets the value for global configuration 'OIDCAllowedGroups' field
func SetOIDCAllowedGroups(v []string) { global.SetOIDCAllowedGroups(v) }

// GetOIDCAdminGroups safely fetches the Configuration value for state's 'OIDCAdminGroups' field
func (st *ConfigState) GetOIDCAdminGroups() (v []string) {
	return st.config.OIDCAdminGroups
}

// SetOIDCAdminGroups safely sets the Configuration value for state's 'OIDCAdminGroups' field
func (st *ConfigState) SetOIDCAdminGroups(v []string) {
	st.config.OIDCAdminGroups = v
	st.reloadToViper()
}

// GetOIDCAdminGroups safely fetches the value for global configuration 'OIDCAdminGroups' field
func GetOIDCAdminGroups() []string { return global.GetOIDCAdminGroups() }

// SetOIDCAdminGroups safely sets the value for global configuration 'OIDCAdminGroups' field
func SetOIDCAdminGroups(v []string) { global.SetOIDCAdminGroups(v) }

// GetTracingEnabled safely fetches the Configuration value for state's 'TracingEnabled' field
func (st *ConfigState) GetTracingEnabled() (v bool) {
	return st.config.TracingEnabled
}

// SetTracingEnabled safely sets the Configuration value for state's 'TracingEnabled' field
func (st *ConfigState) SetTracingEnabled(v bool) {
	st.config.TracingEnabled = v
	st.reloadToViper()
}

// GetTracingEnabled safely fetches the value for global configuration 'TracingEnabled' field
func GetTracingEnabled() bool { return global.GetTracingEnabled() }

// SetTracingEnabled safely sets the value for global configuration 'TracingEnabled' field
func SetTracingEnabled(v bool) { global.SetTracingEnabled(v) }

// GetMetricsEnabled safely fetches the Configuration value for state's 'MetricsEnabled' field
func (st *ConfigState) GetMetricsEnabled() (v bool) {
	return st.config.MetricsEnabled
}

// SetMetricsEnabled safely sets the Configuration value for state's 'MetricsEnabled' field
func (st *ConfigState) SetMetricsEnabled(v bool) {
	st.config.MetricsEnabled = v
	st.reloadToViper()
}

// GetMetricsEnabled safely fetches the value for global configuration 'MetricsEnabled' field
func GetMetricsEnabled() bool { return global.GetMetricsEnabled() }

// SetMetricsEnabled safely sets the value for global configuration 'MetricsEnabled' field
func SetMetricsEnabled(v bool) { global.SetMetricsEnabled(v) }

// GetSMTPHost safely fetches the Configuration value for state's 'SMTPHost' field
func (st *ConfigState) GetSMTPHost() (v string) {
	return st.config.SMTPHost
}

// SetSMTPHost safely sets the Configuration value for state's 'SMTPHost' field
func (st *ConfigState) SetSMTPHost(v string) {
	st.config.SMTPHost = v
	st.reloadToViper()
}

// GetSMTPHost safely fetches the value for global configuration 'SMTPHost' field
func GetSMTPHost() string { return global.GetSMTPHost() }

// SetSMTPHost safely sets the value for global configuration 'SMTPHost' field
func SetSMTPHost(v string) { global.SetSMTPHost(v) }

// GetSMTPPort safely fetches the Configuration value for state's 'SMTPPort' field
func (st *ConfigState) GetSMTPPort() (v int) {
	return st.config.SMTPPort
}

// SetSMTPPort safely sets the Configuration value for state's 'SMTPPort' field
func (st *ConfigState) SetSMTPPort(v int) {
	st.config.SMTPPort = v
	st.reloadToViper()
}

// GetSMTPPort safely fetches the value for global configuration 'SMTPPort' field
func GetSMTPPort() int { return global.GetSMTPPort() }

// SetSMTPPort safely sets the value for global configuration 'SMTPPort' field
func SetSMTPPort(v int) { global.SetSMTPPort(v) }

// GetSMTPUsername safely fetches the Configuration value for state's 'SMTPUsername' field
func (st *ConfigState) GetSMTPUsername() (v string) {
	return st.config.SMTPUsername
}

// SetSMTPUsername safely sets the Configuration value for state's 'SMTPUsername' field
func (st *ConfigState) SetSMTPUsername(v string) {
	st.config.SMTPUsername = v
	st.reloadToViper()
}

// GetSMTPUsername safely fetches the value for global configuration 'SMTPUsername' field
func GetSMTPUsername() string { return global.GetSMTPUsername() }

// SetSMTPUsername safely sets the value for global configuration 'SMTPUsername' field
func SetSMTPUsername(v string) { global.SetSMTPUsername(v) }

// GetSMTPPassword safely fetches the Configuration value for state's 'SMTPPassword' field
func (st *ConfigState) GetSMTPPassword() (v string) {
	return st.config.SMTPPassword
}

// SetSMTPPassword safely sets the Configuration value for state's 'SMTPPassword' field
func (st *ConfigState) SetSMTPPassword(v string) {
	st.config.SMTPPassword = v
	st.reloadToViper()
}

// GetSMTPPassword safely fetches the value for global configuration 'SMTPPassword' field
func GetSMTPPassword() string { return global.GetSMTPPassword() }

// SetSMTPPassword safely sets the value for global configuration 'SMTPPassword' field
func SetSMTPPassword(v string) { global.SetSMTPPassword(v) }

// GetSMTPFrom safely fetches the Configuration value for state's 'SMTPFrom' field
func (st *ConfigState) GetSMTPFrom() (v string) {
	return st.config.SMTPFrom
}

// SetSMTPFrom safely sets the Configuration value for state's 'SMTPFrom' field
func (st *ConfigState) SetSMTPFrom(v string) {
	st.config.SMTPFrom = v
	st.reloadToViper()
}

// GetSMTPFrom safely fetches the value for global configuration 'SMTPFrom' field
func GetSMTPFrom() string { return global.GetSMTPFrom() }

// SetSMTPFrom safely sets the value for global configuration 'SMTPFrom' field
func SetSMTPFrom(v string) { global.SetSMTPFrom(v) }

// GetSMTPFromDisplayName safely fetches the Configuration value for state's 'SMTPFromDisplayName' field
func (st *ConfigState) GetSMTPFromDisplayName() (v string) {
	return st.config.SMTPFromDisplayName
}

// SetSMTPFromDisplayName safely sets the Configuration value for state's 'SMTPFromDisplayName' field
func (st *ConfigState) SetSMTPFromDisplayName(v string) {
	st.config.SMTPFromDisplayName = v
	st.reloadToViper()
}

// GetSMTPFromDisplayName safely fetches the value for global configuration 'SMTPFromDisplayName' field
func GetSMTPFromDisplayName() string { return global.GetSMTPFromDisplayName() }

// SetSMTPFromDisplayName safely sets the value for global configuration 'SMTPFromDisplayName' field
func SetSMTPFromDisplayName(v string) { global.SetSMTPFromDisplayName(v) }

// GetSMTPDiscloseRecipients safely fetches the Configuration value for state's 'SMTPDiscloseRecipients' field
func (st *ConfigState) GetSMTPDiscloseRecipients() (v bool) {
	return st.config.SMTPDiscloseRecipients
}

// SetSMTPDiscloseRecipients safely sets the Configuration value for state's 'SMTPDiscloseRecipients' field
func (st *ConfigState) SetSMTPDiscloseRecipients(v bool) {
	st.config.SMTPDiscloseRecipients = v
	st.reloadToViper()
}

// GetSMTPDiscloseRecipients safely fetches the value for global configuration 'SMTPDiscloseRecipients' field
func GetSMTPDiscloseRecipients() bool { return global.GetSMTPDiscloseRecipients() }

// SetSMTPDiscloseRecipients safely sets the value for global configuration 'SMTPDiscloseRecipients' field
func SetSMTPDiscloseRecipients(v bool) { global.SetSMTPDiscloseRecipients(v) }

// GetSyslogEnabled safely fetches the Configuration value for state's 'SyslogEnabled' field
func (st *ConfigState) GetSyslogEnabled() (v bool) {
	return st.config.SyslogEnabled
}

// SetSyslogEnabled safely sets the Configuration value for state's 'SyslogEnabled' field
func (st *ConfigState) SetSyslogEnabled(v bool) {
	st.config.SyslogEnabled = v
	st.reloadToViper()
}

// GetSyslogEnabled safely fetches the value for global configuration 'SyslogEnabled' field
func GetSyslogEnabled() bool { return global.GetSyslogEnabled() }

// SetSyslogEnabled safely sets the value for global configuration 'SyslogEnabled' field
func SetSyslogEnabled(v bool) { global.SetSyslogEnabled(v) }

// GetSyslogProtocol safely fetches the Configuration value for state's 'SyslogProtocol' field
func (st *ConfigState) GetSyslogProtocol() (v string) {
	return st.config.SyslogProtocol
}

// SetSyslogProtocol safely sets the Configuration value for state's 'SyslogProtocol' field
func (st *ConfigState) SetSyslogProtocol(v string) {
	st.config.SyslogProtocol = v
	st.reloadToViper()
}

// GetSyslogProtocol safely fetches the value for global configuration 'SyslogProtocol' field
func GetSyslogProtocol() string { return global.GetSyslogProtocol() }

// SetSyslogProtocol safely sets the value for global configuration 'SyslogProtocol' field
func SetSyslogProtocol(v string) { global.SetSyslogProtocol(v) }

// GetSyslogAddress safely fetches the Configuration value for state's 'SyslogAddress' field
func (st *ConfigState) GetSyslogAddress() (v string) {
	return st.config.SyslogAddress
}

// SetSyslogAddress safely sets the Configuration value for state's 'SyslogAddress' field
func (st *ConfigState) SetSyslogAddress(v string) {
	st.config.SyslogAddress = v
	st.reloadToViper()
}

// GetSyslogAddress safely fetches the value for global configuration 'SyslogAddress' field
func GetSyslogAddress() string { return global.GetSyslogAddress() }

// SetSyslogAddress safely sets the value for global configuration 'SyslogAddress' field
func SetSyslogAddress(v string) { global.SetSyslogAddress(v) }

// GetSyslogMirror safely fetches the Configuration value for state's 'SyslogMirror' field
func (st *ConfigState) GetSyslogMirror() (v bool) {
	return st.config.SyslogMirror
}

// SetSyslogMirror safely sets the Configuration value for state's 'SyslogMirror' field
func (st *ConfigState) SetSyslogMirror(v bool) {
	st.config.SyslogMirror = v
	st.reloadToViper()
}

// GetSyslogMirror safely fetches the value for global configuration 'SyslogMirror' field
func GetSyslogMirror() bool { return global.GetSyslogMirror() }

// SetSyslogMirror safely sets the value for global configuration 'SyslogMirror' field
func SetSyslogMirror(v bool) { global.SetSyslogMirror(v) }

// GetSyslogMsgLength safely fetches the Configuration value for state's 'SyslogMsgLength' field
func (st *ConfigState) GetSyslogMsgLength() (v uint32) {
	return st.config.SyslogMsgLength
}

// SetSyslogMsgLength safely sets the Configuration value for state's 'SyslogMsgLength' field
func (st *ConfigState) SetSyslogMsgLength(v uint32) {
	st.config.SyslogMsgLength = v
	st.reloadToViper()
}

// GetSyslogMsgLength safely fetches the value for global configuration 'SyslogMsgLength' field
func GetSyslogMsgLength() uint32 { return global.GetSyslogMsgLength() }

// SetSyslogMsgLength safely sets the value for global configuration 'SyslogMsgLength' field
func SetSyslogMsgLength(v uint32) { global.SetSyslogMsgLength(v) }

// GetAdminAccountUsername safely fetches the Configuration value for state's 'AdminAccountUsername' field
func (st *ConfigState) GetAdminAccountUsername() (v string) {
	return st.config.AdminAccountUsername
}

// SetAdminAccountUsername safely sets the Configuration value for state's 'AdminAccountUsername' field
func (st *ConfigState) SetAdminAccountUsername(v string) {
	st.config.AdminAccountUsername = v
	st.reloadToViper()
}

// GetAdminAccountUsername safely fetches the value for global configuration 'AdminAccountUsername' field
func GetAdminAccountUsername() string { return global.GetAdminAccountUsername() }

// SetAdminAccountUsername safely sets the value for global configuration 'AdminAccountUsername' field
func SetAdminAccountUsername(v string) { global.SetAdminAccountUsername(v) }

// GetAdminAccountEmail safely fetches the Configuration value for state's 'AdminAccountEmail' field
func (st *ConfigState) GetAdminAccountEmail() (v string) {
	return st.config.AdminAccountEmail
}

// SetAdminAccountEmail safely sets the Configuration value for state's 'AdminAccountEmail' field
func (st *ConfigState) SetAdminAccountEmail(v string) {
	st.config.AdminAccountEmail = v
	st.reloadToViper()
}

// GetAdminAccountEmail safely fetches the value for global configuration 'AdminAccountEmail' field
func GetAdminAccountEmail() string { return global.GetAdminAccountEmail() }

// SetAdminAccountEmail safely sets the value for global configuration 'AdminAccountEmail' field
func SetAdminAccountEmail(v string) { global.SetAdminAccountEmail(v) }

// GetAdminAccountPassword safely fetches the Configuration value for state's 'AdminAccountPassword' field
func (st *ConfigState) GetAdminAccountPassword() (v string) {
	return st.config.AdminAccountPassword
}

// SetAdminAccountPassword safely sets the Configuration value for state's 'AdminAccountPassword' field
func (st *ConfigState) SetAdminAccountPassword(v string) {
	st.config.AdminAccountPassword = v
	st.reloadToViper()
}

// GetAdminAccountPassword safely fetches the value for global configuration 'AdminAccountPassword' field
func GetAdminAccountPassword() string { return global.GetAdminAccountPassword() }

// SetAdminAccountPassword safely sets the value for global configuration 'AdminAccountPassword' field
func SetAdminAccountPassword(v string) { global.SetAdminAccountPassword(v) }

// GetAdminTransPath safely fetches the Configuration value for state's 'AdminTransPath' field
func (st *ConfigState) GetAdminTransPath() (v string) {
	return st.config.AdminTransPath
}

// SetAdminTransPath safely sets the Configuration value for state's 'AdminTransPath' field
func (st *ConfigState) SetAdminTransPath(v string) {
	st.config.AdminTransPath = v
	st.reloadToViper()
}

// GetAdminTransPath safely fetches the value for global configuration 'AdminTransPath' field
func GetAdminTransPath() string { return global.GetAdminTransPath() }

// SetAdminTransPath safely sets the value for global configuration 'AdminTransPath' field
func SetAdminTransPath(v string) { global.SetAdminTransPath(v) }

// GetAdminMediaPruneDryRun safely fetches the Configuration value for state's 'AdminMediaPruneDryRun' field
func (st *ConfigState) GetAdminMediaPruneDryRun() (v bool) {
	return st.config.AdminMediaPruneDryRun
}

// SetAdminMediaPruneDryRun safely sets the Configuration value for state's 'AdminMediaPruneDryRun' field
func (st *ConfigState) SetAdminMediaPruneDryRun(v bool) {
	st.config.AdminMediaPruneDryRun = v
	st.reloadToViper()
}

// GetAdminMediaPruneDryRun safely fetches the value for global configuration 'AdminMediaPruneDryRun' field
func GetAdminMediaPruneDryRun() bool { return global.GetAdminMediaPruneDryRun() }

// SetAdminMediaPruneDryRun safely sets the value for global configuration 'AdminMediaPruneDryRun' field
func SetAdminMediaPruneDryRun(v bool) { global.SetAdminMediaPruneDryRun(v) }

// GetAdminMediaListLocalOnly safely fetches the Configuration value for state's 'AdminMediaListLocalOnly' field
func (st *ConfigState) GetAdminMediaListLocalOnly() (v bool) {
	return st.config.AdminMediaListLocalOnly
}

// SetAdminMediaListLocalOnly safely sets the Configuration value for state's 'AdminMediaListLocalOnly' field
func (st *ConfigState) SetAdminMediaListLocalOnly(v bool) {
	st.config.AdminMediaListLocalOnly = v
	st.reloadToViper()
}

// GetAdminMediaListLocalOnly safely fetches the value for global configuration 'AdminMediaListLocalOnly' field
func GetAdminMediaListLocalOnly() bool { return global.GetAdminMediaListLocalOnly() }

// SetAdminMediaListLocalOnly safely sets the value for global configuration 'AdminMediaListLocalOnly' field
func SetAdminMediaListLocalOnly(v bool) { global.SetAdminMediaListLocalOnly(v) }

// GetAdminMediaListRemoteOnly safely fetches the Configuration value for state's 'AdminMediaListRemoteOnly' field
func (st *ConfigState) GetAdminMediaListRemoteOnly() (v bool) {
	return st.config.AdminMediaListRemoteOnly
}

// SetAdminMediaListRemoteOnly safely sets the Configuration value for state's 'AdminMediaListRemoteOnly' field
func (st *ConfigState) SetAdminMediaListRemoteOnly(v bool) {
	st.config.AdminMediaListRemoteOnly = v
	st.reloadToViper()
}

// GetAdminMediaListRemoteOnly safely fetches the value for global configuration 'AdminMediaListRemoteOnly' field
func GetAdminMediaListRemoteOnly() bool { return global.GetAdminMediaListRemoteOnly() }

// SetAdminMediaListRemoteOnly safely sets the value for global configuration 'AdminMediaListRemoteOnly' field
func SetAdminMediaListRemoteOnly(v bool) { global.SetAdminMediaListRemoteOnly(v) }

// GetTestrigSkipDBSetup safely fetches the Configuration value for state's 'TestrigSkipDBSetup' field
func (st *ConfigState) GetTestrigSkipDBSetup() (v bool) {
	return st.config.TestrigSkipDBSetup
}

// SetTestrigSkipDBSetup safely sets the Configuration value for state's 'TestrigSkipDBSetup' field
func (st *ConfigState) SetTestrigSkipDBSetup(v bool) {
	st.config.TestrigSkipDBSetup = v
	st.reloadToViper()
}

// GetTestrigSkipDBSetup safely fetches the value for global configuration 'TestrigSkipDBSetup' field
func GetTestrigSkipDBSetup() bool { return global.GetTestrigSkipDBSetup() }

// SetTestrigSkipDBSetup safely sets the value for global configuration 'TestrigSkipDBSetup' field
func SetTestrigSkipDBSetup(v bool) { global.SetTestrigSkipDBSetup(v) }

// GetTestrigSkipDBTeardown safely fetches the Configuration value for state's 'TestrigSkipDBTeardown' field
func (st *ConfigState) GetTestrigSkipDBTeardown() (v bool) {
	return st.config.TestrigSkipDBTeardown
}

// SetTestrigSkipDBTeardown safely sets the Configuration value for state's 'TestrigSkipDBTeardown' field
func (st *ConfigState) SetTestrigSkipDBTeardown(v bool) {
	st.config.TestrigSkipDBTeardown = v
	st.reloadToViper()
}

// GetTestrigSkipDBTeardown safely fetches the value for global configuration 'TestrigSkipDBTeardown' field
func GetTestrigSkipDBTeardown() bool { return global.GetTestrigSkipDBTeardown() }

// SetTestrigSkipDBTeardown safely sets the value for global configuration 'TestrigSkipDBTeardown' field
func SetTestrigSkipDBTeardown(v bool) { global.SetTestrigSkipDBTeardown(v) }

// GetTotalOfMemRatios safely fetches the combined value for all the state's mem ratio fields
func (st *ConfigState) GetTotalOfMemRatios() (total float64) {
	total += st.config.Cache.AccountMemRatio
	total += st.config.Cache.AccountNoteMemRatio
	total += st.config.Cache.AccountSettingsMemRatio
	total += st.config.Cache.AccountStatsMemRatio
	total += st.config.Cache.ApplicationMemRatio
	total += st.config.Cache.BlockMemRatio
	total += st.config.Cache.BlockIDsMemRatio
	total += st.config.Cache.BoostOfIDsMemRatio
	total += st.config.Cache.ClientMemRatio
	total += st.config.Cache.ConversationMemRatio
	total += st.config.Cache.ConversationLastStatusIDsMemRatio
	total += st.config.Cache.DomainPermissionDraftMemRatio
	total += st.config.Cache.DomainLimitMemRatio
	total += st.config.Cache.DomainPermissionSubscriptionMemRatio
	total += st.config.Cache.EmojiMemRatio
	total += st.config.Cache.EmojiCategoryMemRatio
	total += st.config.Cache.FederationErrorMemRatio
	total += st.config.Cache.FilterMemRatio
	total += st.config.Cache.FilterIDsMemRatio
	total += st.config.Cache.FilterKeywordMemRatio
	total += st.config.Cache.FilterStatusMemRatio
	total += st.config.Cache.FollowMemRatio
	total += st.config.Cache.FollowIDsMemRatio
	total += st.config.Cache.FollowRequestMemRatio
	total += st.config.Cache.FollowRequestIDsMemRatio
	total += st.config.Cache.FollowingTagIDsMemRatio
	total += st.config.Cache.HomeAccountIDsMemRatio
	total += st.config.Cache.InReplyToIDsMemRatio
	total += st.config.Cache.InstanceMemRatio
	total += st.config.Cache.InteractionRequestMemRatio
	total += st.config.Cache.ListMemRatio
	total += st.config.Cache.ListIDsMemRatio
	total += st.config.Cache.ListedIDsMemRatio
	total += st.config.Cache.MarkerMemRatio
	total += st.config.Cache.MediaMemRatio
	total += st.config.Cache.MentionMemRatio
	total += st.config.Cache.MoveMemRatio
	total += st.config.Cache.NotificationMemRatio
	total += st.config.Cache.PollMemRatio
	total += st.config.Cache.PollVoteMemRatio
	total += st.config.Cache.PollVoteIDsMemRatio
	total += st.config.Cache.ReportMemRatio
	total += st.config.Cache.RelayActorMemRatio
	total += st.config.Cache.RelayMatcherMemRatio
	total += st.config.Cache.RelayPushMemRatio
	total += st.config.Cache.RelayPushIDsMemRatio
	total += st.config.Cache.RelaySubscriptionMemRatio
	total += st.config.Cache.ScheduledStatusMemRatio
	total += st.config.Cache.SinBinStatusMemRatio
	total += st.config.Cache.StatusMemRatio
	total += st.config.Cache.StatusBookmarkMemRatio
	total += st.config.Cache.StatusBookmarkIDsMemRatio
	total += st.config.Cache.StatusEditMemRatio
	total += st.config.Cache.StatusFaveMemRatio
	total += st.config.Cache.StatusFaveIDsMemRatio
	total += st.config.Cache.StatusPinnedIDsMemRatio
	total += st.config.Cache.TagMemRatio
	total += st.config.Cache.ThreadMuteMemRatio
	total += st.config.Cache.TokenMemRatio
	total += st.config.Cache.TombstoneMemRatio
	total += st.config.Cache.UserMemRatio
	total += st.config.Cache.UserMuteMemRatio
	total += st.config.Cache.UserMuteIDsMemRatio
	total += st.config.Cache.WebfingerMemRatio
	total += st.config.Cache.WebPushSubscriptionMemRatio
	total += st.config.Cache.WebPushSubscriptionIDsMemRatio
	total += st.config.Cache.MutesMemRatio
	total += st.config.Cache.StatusFilterMemRatio
	total += st.config.Cache.VisibilityMemRatio
	return
}

// GetTotalOfMemRatios safely fetches the combined value for all the global state's mem ratio fields
func GetTotalOfMemRatios() (total float64) { return global.GetTotalOfMemRatios() }
