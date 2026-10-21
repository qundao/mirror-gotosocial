#!/bin/sh

set -eu

EXPECT=$(cat << "EOF"
{
    "account-domain": "peepee",
    "accounts-allow-custom-css": true,
    "accounts-custom-css-length": 5000,
    "accounts-max-profile-fields": 8,
    "accounts-reason-required": false,
    "accounts-registration-backlog-limit": 100,
    "accounts-registration-daily-limit": 50,
    "accounts-registration-open": true,
    "advanced-cookies-samesite": "strict",
    "advanced-csp-extra-uris": [],
    "advanced-header-filter-mode": "block",
    "advanced-rate-limit-exceptions": [
        "192.0.2.0/24",
        "127.0.0.1/32"
    ],
    "advanced-rate-limit-requests": 6969,
    "advanced-sender-multiplier": -1,
    "advanced-throttling-multiplier": -1,
    "advanced-throttling-retry-after": 10000000000,
    "application-name": "gts",
    "bind-address": "127.0.0.1",
    "cache-account-mem-ratio": 5,
    "cache-account-note-mem-ratio": 1,
    "cache-account-settings-mem-ratio": 0.1,
    "cache-account-stats-mem-ratio": 2,
    "cache-application-mem-ratio": 0.1,
    "cache-block-ids-mem-ratio": 3,
    "cache-block-mem-ratio": 2,
    "cache-boost-of-ids-mem-ratio": 3,
    "cache-client-mem-ratio": 0.1,
    "cache-conversation-last-status-ids-mem-ratio": 2,
    "cache-conversation-mem-ratio": 1,
    "cache-domain-permission-draft-mem-ratio": 0.5,
    "cache-domain-permission-limit-mem-ratio": 0.5,
    "cache-domain-permission-subscription-mem-ratio": 0.5,
    "cache-emoji-category-mem-ratio": 0.1,
    "cache-emoji-mem-ratio": 3,
    "cache-federation-error-mem-ratio": 0.2,
    "cache-filter-ids-mem-ratio": 2,
    "cache-filter-keyword-mem-ratio": 0.5,
    "cache-filter-mem-ratio": 0.5,
    "cache-filter-status-mem-ratio": 0.5,
    "cache-follow-ids-mem-ratio": 4,
    "cache-follow-mem-ratio": 2,
    "cache-follow-request-ids-mem-ratio": 2,
    "cache-follow-request-mem-ratio": 2,
    "cache-following-tag-ids-mem-ratio": 2,
    "cache-home-account-ids-mem-ratio": 2,
    "cache-home-timeline-size": 800,
    "cache-home-timeline-timeout": 21600000000000,
    "cache-in-reply-to-ids-mem-ratio": 3,
    "cache-instance-mem-ratio": 1,
    "cache-interaction-request-mem-ratio": 1,
    "cache-list-ids-mem-ratio": 2,
    "cache-list-mem-ratio": 1,
    "cache-list-timeline-size": 800,
    "cache-list-timeline-timeout": 7200000000000,
    "cache-listed-ids-mem-ratio": 2,
    "cache-marker-mem-ratio": 0.5,
    "cache-media-mem-ratio": 4,
    "cache-memory-target": "100MiB",
    "cache-mention-mem-ratio": 2,
    "cache-move-mem-ratio": 0.1,
    "cache-mutes-mem-ratio": 2,
    "cache-notification-mem-ratio": 2,
    "cache-notification-timeline-size": 200,
    "cache-notification-timeline-timeout": 21600000000000,
    "cache-poll-mem-ratio": 1,
    "cache-poll-vote-ids-mem-ratio": 2,
    "cache-poll-vote-mem-ratio": 2,
    "cache-relay-actor-mem-ratio": 0.1,
    "cache-relay-matcher-mem-ratio": 0.1,
    "cache-relay-push-ids-mem-ratio": 0.1,
    "cache-relay-push-mem-ratio": 0.1,
    "cache-relay-subscription-mem-ratio": 0.1,
    "cache-report-mem-ratio": 1,
    "cache-s3-object-info": 0,
    "cache-scheduled-status-mem-ratio": 4,
    "cache-sin-bin-status-mem-ratio": 0.5,
    "cache-status-bookmark-ids-mem-ratio": 2,
    "cache-status-bookmark-mem-ratio": 0.5,
    "cache-status-edit-mem-ratio": 2,
    "cache-status-fave-ids-mem-ratio": 3,
    "cache-status-fave-mem-ratio": 2,
    "cache-status-filter-mem-ratio": 7,
    "cache-status-mem-ratio": 5,
    "cache-status-pinned-ids-mem-ratio": 3,
    "cache-tag-mem-ratio": 2,
    "cache-tag-timeline-size": 400,
    "cache-tag-timeline-timeout": 600000000000,
    "cache-thread-mute-mem-ratio": 0.2,
    "cache-token-mem-ratio": 0.75,
    "cache-tombstone-mem-ratio": 0.5,
    "cache-user-mem-ratio": 0.25,
    "cache-user-mute-ids-mem-ratio": 3,
    "cache-user-mute-mem-ratio": 2,
    "cache-visibility-mem-ratio": 2,
    "cache-web-push-subscription-ids-mem-ratio": 1,
    "cache-web-push-subscription-mem-ratio": 1,
    "cache-webfinger-mem-ratio": 0.1,
    "config-path": "internal/config/testdata/test.yaml",
    "db-address": ":memory:",
    "db-database": "gotosocial_prod",
    "db-max-open-conns-multiplier": 3,
    "db-password": "hunter2",
    "db-port": 6969,
    "db-postgres-connection-string": "",
    "db-sqlite-busy-timeout": 1000000000,
    "db-sqlite-cache-size": "0B",
    "db-sqlite-journal-mode": "DELETE",
    "db-sqlite-synchronous": "FULL",
    "db-tls-ca-cert": "",
    "db-tls-mode": "disable",
    "db-type": "sqlite",
    "db-user": "sex-haver",
    "dry-run": true,
    "email": "",
    "host": "example.com",
    "http-client-allow-ips": null,
    "http-client-block-ips": null,
    "http-client-disable-keep-alives": false,
    "http-client-idle-conn-timeout": 300000000000,
    "http-client-insecure-outgoing": false,
    "http-client-max-conns-per-host": 6,
    "http-client-max-idle-conns": 0,
    "http-client-max-idle-conns-per-host": 0,
    "http-client-read-buffer-size": "0B",
    "http-client-response-header-timeout": 0,
    "http-client-timeout": 60000000000,
    "http-client-tls-handshake-timeout": 0,
    "http-client-tls-insecure-skip-verify": false,
    "http-client-write-buffer-size": "0B",
    "http-server-idle-timeout": 30000000000,
    "http-server-max-concurrent-streams": 0,
    "http-server-max-decoder-header-table-size": "0B",
    "http-server-max-encoder-header-table-size": "0B",
    "http-server-max-header-bytes": "0B",
    "http-server-max-multipart-memory": "40.0MiB",
    "http-server-max-read-frame-size": "0B",
    "http-server-max-receive-buffer-per-connection": "0B",
    "http-server-max-receive-buffer-per-stream": "0B",
    "http-server-ping-timeout": 0,
    "http-server-read-header-timeout": 0,
    "http-server-read-timeout": 60000000000,
    "http-server-send-ping-timeout": 0,
    "http-server-use-h2c": false,
    "http-server-write-byte-timeout": 0,
    "http-server-write-timeout": 30000000000,
    "instance-allow-backdating-statuses": true,
    "instance-deliver-to-shared-inboxes": false,
    "instance-directory-mode": "open",
    "instance-expose-allowlist": true,
    "instance-expose-allowlist-web": true,
    "instance-expose-blocklist": true,
    "instance-expose-blocklist-web": true,
    "instance-expose-custom-emojis": true,
    "instance-expose-peers": true,
    "instance-expose-public-timeline": true,
    "instance-federation-mode": "allowlist",
    "instance-federation-spam-filter": true,
    "instance-inject-mastodon-version": true,
    "instance-languages": [
        "nl",
        "en-GB"
    ],
    "instance-robots-allow-indexing": true,
    "instance-stats-mode": "baffle",
    "instance-subscriptions-process-cron": "0 23 * * *",
    "landing-page-user": "admin",
    "letsencrypt-cert-dir": "/gotosocial/storage/certs",
    "letsencrypt-email-address": "",
    "letsencrypt-enabled": true,
    "letsencrypt-port": 80,
    "local-only": false,
    "log-client-ip": false,
    "log-db-queries": true,
    "log-format": "json",
    "log-level": "info",
    "log-timestamp-format": "banana",
    "media-cleanup-cron": "0 0 * * *",
    "media-description-max-chars": 5000,
    "media-description-min-chars": 69,
    "media-emoji-local-max-size": "420B",
    "media-emoji-remote-max-size": "420B",
    "media-ffmpeg-pool-size": 8,
    "media-image-size-hint": "5.00MiB",
    "media-local-max-size": "420B",
    "media-remote-cache-duration": "1 week",
    "media-remote-max-size": "420B",
    "media-thumb-max-pixels": 42069,
    "media-video-size-hint": "40.0MiB",
    "metrics-enabled": false,
    "oidc-admin-groups": [
        "steamy"
    ],
    "oidc-allowed-groups": [
        "sloths"
    ],
    "oidc-client-id": "1234",
    "oidc-client-secret": "shhhh its a secret",
    "oidc-enabled": true,
    "oidc-idp-name": "sex-haver",
    "oidc-issuer": "whoknows",
    "oidc-link-existing": true,
    "oidc-scopes": [
        "read",
        "write"
    ],
    "oidc-skip-verification": true,
    "password": "",
    "path": "",
    "port": 6969,
    "protocol": "http",
    "remote-only": false,
    "request-id-header": "X-Trace-Id",
    "scheduled-statuses-max-daily": 25,
    "scheduled-statuses-max-total": 300,
    "skip-db-setup": false,
    "skip-db-teardown": false,
    "smtp-disclose-recipients": true,
    "smtp-from": "queen.rip.in.piss@terfisland.org",
    "smtp-from-display-name": "yes queen",
    "smtp-host": "example.com",
    "smtp-password": "hunter2",
    "smtp-port": 4269,
    "smtp-username": "sex-haver",
    "software-version": "",
    "statuses-cleanup-cron": "0 1 * * 0",
    "statuses-cleanup-remote-older-than": "0 sec",
    "statuses-max-chars": 69,
    "statuses-media-max-files": 1,
    "statuses-poll-max-options": 1,
    "statuses-poll-option-max-chars": 50,
    "storage-backend": "local",
    "storage-local-base-path": "/root/store",
    "storage-s3-access-key": "minio",
    "storage-s3-bucket": "gts",
    "storage-s3-bucket-lookup": "auto",
    "storage-s3-endpoint": "localhost:9000",
    "storage-s3-key-prefix": "",
    "storage-s3-proxy": true,
    "storage-s3-redirect-url": "",
    "storage-s3-region": "us-east-2",
    "storage-s3-secret-key": "miniostorage",
    "storage-s3-use-ssl": false,
    "syslog-address": "127.0.0.1:6969",
    "syslog-enabled": true,
    "syslog-mirror": true,
    "syslog-msg-length": 2048,
    "syslog-protocol": "udp",
    "tls-certificate-chain": "",
    "tls-certificate-key": "",
    "tracing-enabled": false,
    "trusted-proxies": [
        "127.0.0.1/32",
        "127.0.1.0/24"
    ],
    "username": "",
    "web-asset-base-dir": "/root",
    "web-template-base-dir": "/root"
}
EOF
)

# Set all the environment variables to
# ensure that these are parsed without panic
OUTPUT=$(GTS_LOG_LEVEL='info' \
GTS_LOG_TIMESTAMP_FORMAT="banana" \
GTS_LOG_DB_QUERIES=true \
GTS_LOG_CLIENT_IP=false \
GTS_LOG_FORMAT=json \
GTS_APPLICATION_NAME=gts \
GTS_LANDING_PAGE_USER=admin \
GTS_HOST=example.com \
GTS_ACCOUNT_DOMAIN='peepee' \
GTS_PROTOCOL=http \
GTS_BIND_ADDRESS='127.0.0.1' \
GTS_PORT=6969 \
GTS_TRUSTED_PROXIES='127.0.0.1/32,127.0.1.0/24' \
GTS_DB_TYPE='sqlite' \
GTS_DB_POSTGRES_CONNECTION_STRING='' \
GTS_DB_ADDRESS=':memory:' \
GTS_DB_PORT=6969 \
GTS_DB_USER='sex-haver' \
GTS_DB_PASSWORD='hunter2' \
GTS_DB_DATABASE='gotosocial_prod' \
GTS_DB_MAX_OPEN_CONNS_MULTIPLIER=3 \
GTS_DB_SQLITE_JOURNAL_MODE='DELETE' \
GTS_DB_SQLITE_SYNCHRONOUS='FULL' \
GTS_DB_SQLITE_CACHE_SIZE=0 \
GTS_DB_SQLITE_BUSY_TIMEOUT='1s' \
GTS_TLS_MODE='' \
GTS_DB_TLS_CA_CERT='' \
GTS_WEB_TEMPLATE_BASE_DIR='/root' \
GTS_WEB_ASSET_BASE_DIR='/root' \
GTS_INSTANCE_EXPOSE_PEERS=true \
GTS_INSTANCE_EXPOSE_BLOCKLIST=true \
GTS_INSTANCE_EXPOSE_BLOCKLIST_WEB=true \
GTS_INSTANCE_EXPOSE_ALLOWLIST=true \
GTS_INSTANCE_EXPOSE_ALLOWLIST_WEB=true \
GTS_INSTANCE_EXPOSE_CUSTOM_EMOJIS=true \
GTS_INSTANCE_EXPOSE_PUBLIC_TIMELINE=true \
GTS_INSTANCE_FEDERATION_MODE='allowlist' \
GTS_INSTANCE_FEDERATION_SPAM_FILTER=true \
GTS_INSTANCE_DELIVER_TO_SHARED_INBOXES=false \
GTS_INSTANCE_DIRECTORY_MODE="open" \
GTS_INSTANCE_INJECT_MASTODON_VERSION=true \
GTS_INSTANCE_LANGUAGES="nl,en-gb" \
GTS_INSTANCE_ROBOTS_ALLOW_INDEXING=true \
GTS_INSTANCE_STATS_MODE="baffle" \
GTS_ACCOUNTS_ALLOW_CUSTOM_CSS=true \
GTS_ACCOUNTS_CUSTOM_CSS_LENGTH=5000 \
GTS_ACCOUNTS_MAX_PROFILE_FIELDS=8 \
GTS_ACCOUNTS_REGISTRATION_BACKLOG_LIMIT=100 \
GTS_ACCOUNTS_REGISTRATION_DAILY_LIMIT=50 \
GTS_ACCOUNTS_REGISTRATION_OPEN=true \
GTS_ACCOUNTS_REASON_REQUIRED=false \
GTS_MEDIA_DESCRIPTION_MIN_CHARS=69 \
GTS_MEDIA_DESCRIPTION_MAX_CHARS=5000 \
GTS_MEDIA_IMAGE_SIZE_HINT='5MiB' \
GTS_MEDIA_LOCAL_MAX_SIZE=420 \
GTS_MEDIA_REMOTE_MAX_SIZE=420 \
GTS_MEDIA_EMOJI_LOCAL_MAX_SIZE=420 \
GTS_MEDIA_EMOJI_REMOTE_MAX_SIZE=420 \
GTS_MEDIA_FFMPEG_POOL_SIZE=8 \
GTS_MEDIA_VIDEO_SIZE_HINT='40MiB' \
GTS_MEDIA_THUMB_MAX_PIXELS=42069 \
GTS_METRICS_ENABLED=false \
GTS_STORAGE_BACKEND='local' \
GTS_STORAGE_LOCAL_BASE_PATH='/root/store' \
GTS_STORAGE_S3_ACCESS_KEY='minio' \
GTS_STORAGE_S3_SECRET_KEY='miniostorage' \
GTS_STORAGE_S3_ENDPOINT='localhost:9000' \
GTS_STORAGE_S3_USE_SSL='false' \
GTS_STORAGE_S3_BUCKET_LOOKUP='auto' \
GTS_STORAGE_S3_PROXY='true' \
GTS_STORAGE_S3_REDIRECT_URL='' \
GTS_STORAGE_S3_BUCKET='gts' \
GTS_STORAGE_S3_REGION='us-east-2' \
GTS_STATUSES_MAX_CHARS=69 \
GTS_STATUSES_CW_MAX_CHARS=420 \
GTS_STATUSES_POLL_MAX_OPTIONS=1 \
GTS_STATUSES_POLL_OPTIONS_MAX_CHARS=69 \
GTS_STATUSES_MEDIA_MAX_FILES=1 \
GTS_LETS_ENCRYPT_ENABLED=false \
GTS_LETS_ENCRYPT_PORT=8080 \
GTS_LETS_ENCRYPT_CERT_DIR='/root/certs' \
GTS_LETS_ENCRYPT_EMAIL_ADDRESS='le@example.com' \
GTS_OIDC_ENABLED=true \
GTS_OIDC_IDP_NAME='sex-haver' \
GTS_OIDC_SKIP_VERIFICATION=true \
GTS_OIDC_ISSUER='whoknows' \
GTS_OIDC_CLIENT_ID='1234' \
GTS_OIDC_CLIENT_SECRET='shhhh its a secret' \
GTS_OIDC_SCOPES='read,write' \
GTS_OIDC_LINK_EXISTING=true \
GTS_OIDC_ALLOWED_GROUPS='sloths' \
GTS_OIDC_ADMIN_GROUPS='steamy' \
GTS_SMTP_HOST='example.com' \
GTS_SMTP_PORT=4269 \
GTS_SMTP_USERNAME='sex-haver' \
GTS_SMTP_PASSWORD='hunter2' \
GTS_SMTP_FROM='queen.rip.in.piss@terfisland.org' \
GTS_SMTP_FROM_DISPLAY_NAME='yes queen' \
GTS_SMTP_DISCLOSE_RECIPIENTS=true \
GTS_SYSLOG_ENABLED=true \
GTS_SYSLOG_PROTOCOL='udp' \
GTS_SYSLOG_ADDRESS='127.0.0.1:6969' \
GTS_ADVANCED_COOKIES_SAMESITE='strict' \
GTS_ADVANCED_RATE_LIMIT_EXCEPTIONS="192.0.2.0/24,127.0.0.1/32" \
GTS_ADVANCED_RATE_LIMIT_REQUESTS=6969 \
GTS_ADVANCED_SENDER_MULTIPLIER=-1 \
GTS_ADVANCED_THROTTLING_MULTIPLIER=-1 \
GTS_ADVANCED_THROTTLING_RETRY_AFTER='10s' \
GTS_ADVANCED_HEADER_FILTER_MODE='block' \
GTS_REQUEST_ID_HEADER='X-Trace-Id' \
go run ./cmd/gotosocial/... --config-path internal/config/testdata/test.yaml debug config)

OUTPUT_OUT=$(mktemp)
echo "$OUTPUT" > "$OUTPUT_OUT"

EXPECT_OUT=$(mktemp)
echo "$EXPECT" > "$EXPECT_OUT"

DIFFCMD=$(command -v diff 2>&1)
if command -v jd >/dev/null 2>&1; then
    DIFFCMD=$(command -v jd 2>&1)
fi

if ! DIFF=$("$DIFFCMD" "$OUTPUT_OUT" "$EXPECT_OUT"); then
    echo "OUTPUT not equal EXPECTED"
    echo "$DIFF"
    exit 1
else
    echo "OK"
    exit 0
fi
