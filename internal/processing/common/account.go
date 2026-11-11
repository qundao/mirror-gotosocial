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

package common

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"

	"code.superseriousbusiness.org/gopkg/log"
	"code.superseriousbusiness.org/gotosocial/internal/ap"
	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/internal/messages"
	"code.superseriousbusiness.org/gotosocial/internal/text"
	"code.superseriousbusiness.org/gotosocial/internal/typeutils"
	"code.superseriousbusiness.org/gotosocial/internal/util"
	"code.superseriousbusiness.org/gotosocial/internal/validate"
	"codeberg.org/gruf/go-iotools"
)

// GetTargetAccountBy fetches the target account with db load function, given the authorized (or, nil) requester's
// account. This returns an approprate gtserror.WithCode accounting (ha) for not found and visibility to requester.
func (p *Processor) GetTargetAccountBy(
	ctx context.Context,
	requester *gtsmodel.Account,
	getTargetFromDB func() (*gtsmodel.Account, error),
) (
	account *gtsmodel.Account,
	visible bool,
	errWithCode gtserror.WithCode,
) {
	// Fetch the target account from db.
	target, err := getTargetFromDB()
	if err != nil && !errors.Is(err, db.ErrNoEntries) {
		err := gtserror.Newf("error getting from db: %w", err)
		return nil, false, gtserror.NewErrorInternalError(err)
	}

	if target == nil {
		// DB loader could not find account in database.
		const text = "target account not found"
		return nil, false, gtserror.NewErrorNotFound(
			errors.New(text),
			text,
		)
	}

	// Check whether target account is visible to requesting account.
	visible, err = p.visFilter.AccountVisible(ctx, requester, target)
	if err != nil {
		err := gtserror.Newf("error checking visibility: %w", err)
		return nil, false, gtserror.NewErrorInternalError(err)
	}

	if requester != nil && visible {
		// Only refresh account if visible to requester,
		// and there is *authorized* requester to prevent
		// a possible DOS vector for unauthorized clients.
		latest, _, err := p.federator.RefreshAccount(ctx,
			requester,
			target,
			nil,
			nil,
		)
		if err != nil {
			log.Errorf(ctx, "error refreshing target %s: %v", target.URI, err)
			return target, visible, nil
		}

		// Set latest.
		target = latest
	}

	return target, visible, nil
}

// GetTargetAccountByID is a call-through to GetTargetAccountBy() using the db GetAccountByID() function.
func (p *Processor) GetTargetAccountByID(
	ctx context.Context,
	requester *gtsmodel.Account,
	targetID string,
) (
	account *gtsmodel.Account,
	visible bool,
	errWithCode gtserror.WithCode,
) {
	return p.GetTargetAccountBy(ctx, requester, func() (*gtsmodel.Account, error) {
		return p.state.DB.GetAccountByID(ctx, targetID)
	})
}

// GetVisibleTargetAccount calls GetTargetAccountByID(),
// but converts a non-visible result to not-found error.
func (p *Processor) GetVisibleTargetAccount(
	ctx context.Context,
	requester *gtsmodel.Account,
	targetID string,
) (
	account *gtsmodel.Account,
	errWithCode gtserror.WithCode,
) {
	// Fetch the target account by ID from the database.
	target, visible, errWithCode := p.GetTargetAccountByID(ctx,
		requester,
		targetID,
	)
	if errWithCode != nil {
		return nil, errWithCode
	}

	if !visible {
		// Pretend account doesn't exist if not visible.
		const text = "target account not found"
		return nil, gtserror.NewErrorNotFound(
			errors.New(text),
			text,
		)
	}

	return target, nil
}

// GetAPIAccount fetches the appropriate API account
// model depending on whether requester = target.
func (p *Processor) GetAPIAccount(
	ctx context.Context,
	requester *gtsmodel.Account,
	target *gtsmodel.Account,
) (
	apiAcc *apimodel.Account,
	errWithCode gtserror.WithCode,
) {
	var err error

	if requester != nil && requester.ID == target.ID {
		// Only return sensitive account model _if_ requester = target.
		apiAcc, err = p.converter.AccountToAPIAccountSensitive(ctx, target)
	} else {
		// Else, fall back to returning the public account model.
		apiAcc, err = p.converter.AccountToAPIAccountPublic(ctx, target)
	}

	if err != nil {
		err := gtserror.Newf("error converting: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	return apiAcc, nil
}

// GetAPIAccountBlocked fetches the limited
// "blocked" account model for given target.
func (p *Processor) GetAPIAccountBlocked(
	ctx context.Context,
	targetAcc *gtsmodel.Account,
) (
	apiAcc *apimodel.Account,
	errWithCode gtserror.WithCode,
) {
	apiAccount, err := p.converter.AccountToAPIAccountBlocked(ctx, targetAcc)
	if err != nil {
		err := gtserror.Newf("error converting: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}
	return apiAccount, nil
}

// GetAPIAccountSensitive fetches the "sensitive" account model for the given target.
// *BE CAREFUL!* Only return a sensitive account if targetAcc == account making the request.
func (p *Processor) GetAPIAccountSensitive(
	ctx context.Context,
	targetAcc *gtsmodel.Account,
) (
	apiAcc *apimodel.Account,
	errWithCode gtserror.WithCode,
) {
	apiAccount, err := p.converter.AccountToAPIAccountSensitive(ctx, targetAcc)
	if err != nil {
		err := gtserror.Newf("error converting: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}
	return apiAccount, nil
}

// GetVisibleAPIAccounts converts an array of gtsmodel.Accounts (inputted by next function) into
// public API model accounts, checking first for visibility. Please note that all errors will be
// logged at ERROR level, but will not be returned. Callers are likely to run into show-stopping
// errors in the lead-up to this function, whereas calling this should not be a show-stopper.
func (p *Processor) GetVisibleAPIAccounts(
	ctx context.Context,
	requester *gtsmodel.Account,
	next func(int) *gtsmodel.Account,
	length int,
) []*apimodel.Account {
	return p.getVisibleAPIAccounts(ctx, 3, requester, next, length)
}

// GetVisibleAPIAccountsPaged is functionally equivalent to GetVisibleAPIAccounts(),
// except the accounts are returned as a converted slice of accounts as interface{}.
func (p *Processor) GetVisibleAPIAccountsPaged(
	ctx context.Context,
	requester *gtsmodel.Account,
	next func(int) *gtsmodel.Account,
	length int,
) []interface{} {
	accounts := p.getVisibleAPIAccounts(ctx, 3, requester, next, length)
	items := make([]interface{}, len(accounts))
	for i, account := range accounts {
		items[i] = account
	}
	return items
}

func (p *Processor) getVisibleAPIAccounts(
	ctx context.Context,
	calldepth int, // used to skip wrapping func above these's names
	requester *gtsmodel.Account,
	next func(int) *gtsmodel.Account,
	length int,
) []*apimodel.Account {
	// Start new log entry with
	// the above calling func's name.
	l := log.WithContext(ctx).
		WithField("caller", log.Caller(calldepth+1))

	// Preallocate slice according to expected length.
	accounts := make([]*apimodel.Account, 0, length)

	for i := 0; i < length; i++ {
		// Get next account.
		account := next(i)
		if account == nil {
			continue
		}

		// Check whether this account is visible to requesting account.
		visible, err := p.visFilter.AccountVisible(ctx, requester, account)
		if err != nil {
			l.Errorf("error checking account visibility: %v", err)
			continue
		}

		if !visible {
			// Not visible to requester.
			continue
		}

		// Convert the account to a public API model representation.
		apiAcc, err := p.converter.AccountToAPIAccountPublic(ctx, account)
		if err != nil {
			l.Errorf("error converting account: %v", err)
			continue
		}

		// Append API model to return slice.
		accounts = append(accounts, apiAcc)
	}

	return accounts
}

// Update processes the update of an account
// with the given form and textContentType
// (content type to use when parsing bio, etc.).
func (p *Processor) UpdateAccount(
	ctx context.Context,
	account *gtsmodel.Account,
	form *apimodel.UpdateAccountRequest,
	textContentType string,
) gtserror.WithCode {
	// Ensure account populated.
	if err := p.state.DB.PopulateAccount(ctx, account); err != nil {
		log.Errorf(ctx, "error(s) populating account, will continue: %s", err)
	}

	var (
		// Indicates that the account's
		// note, display name, and/or fields
		// have changed, and so emojis should
		// be re-parsed and updated as well.
		textChanged bool

		// DB columns on the account
		// that need to be updated.
		acctColumns []string

		// DB columns on the settings
		// that need to be updated.
		settingsColumns []string
	)

	// Account flags.

	if form.Discoverable != nil {
		account.Discoverable = form.Discoverable
		acctColumns = append(acctColumns, "discoverable")
	}

	if form.Indexable != nil {
		account.Indexable = form.Indexable
		acctColumns = append(acctColumns, "indexable")
	}

	if bot := form.Bot; bot != nil {
		if *bot {
			// Mark account as a Service.
			// See: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-service
			account.ActorType = gtsmodel.AccountActorTypeService
		} else {
			// Mark account as a Person.
			// See: https://www.w3.org/TR/activitystreams-vocabulary/#dfn-person
			account.ActorType = gtsmodel.AccountActorTypePerson
		}
		acctColumns = append(acctColumns, "actor_type")
	}

	if form.Locked != nil {
		account.Locked = form.Locked
		acctColumns = append(acctColumns, "locked")
	}

	// Account text fields.

	if form.DisplayName != nil {
		// Display name text
		// is changing.
		textChanged = true

		displayName := *form.DisplayName
		if err := validate.DisplayName(displayName); err != nil {
			return gtserror.NewErrorBadRequest(err, err.Error())
		}

		// HTML tags not allowed in display name.
		account.DisplayName = text.StripHTMLFromText(displayName)
		acctColumns = append(acctColumns, "display_name")
	}

	if form.Note != nil {
		// Note text is changing.
		textChanged = true

		note := *form.Note
		if err := validate.Note(note); err != nil {
			return gtserror.NewErrorBadRequest(err, err.Error())
		}

		// Store raw version of note
		// for now, we'll process
		// the proper version later.
		account.NoteRaw = note
		acctColumns = append(acctColumns, []string{
			"note",
			"note_raw",
		}...)
	}

	if form.FieldsAttributes != nil {
		// Field text is changing.
		textChanged = true

		if err := p.UpdateFields(
			account,
			*form.FieldsAttributes,
		); err != nil {
			return err
		}
		acctColumns = append(acctColumns, []string{
			"fields",
			"fields_raw",
		}...)
	}

	if textChanged {
		// Process display name, note, fields,
		// and any concomitant emoji changes.
		p.ProcessAccountText(ctx,
			account,
			textContentType,
		)
		acctColumns = append(acctColumns, "emojis")
	}

	// Account avatar + header.

	if form.AvatarDescription != nil {
		desc := text.StripHTMLFromText(*form.AvatarDescription)
		form.AvatarDescription = &desc
	}

	if form.Avatar != nil && form.Avatar.Size != 0 {
		avatarInfo, errWithCode := p.UpdateAvatar(ctx,
			account,
			form.Avatar,
			form.AvatarDescription,
		)
		if errWithCode != nil {
			return errWithCode
		}
		account.AvatarMediaAttachmentID = avatarInfo.ID
		account.AvatarMediaAttachment = avatarInfo
		acctColumns = append(acctColumns, "avatar_media_attachment_id")
	} else if form.AvatarDescription != nil && account.AvatarMediaAttachment != nil {
		// Update just existing description if possible.
		account.AvatarMediaAttachment.Description = *form.AvatarDescription
		if err := p.state.DB.UpdateAttachment(
			ctx,
			account.AvatarMediaAttachment,
			"description",
		); err != nil {
			err := gtserror.Newf("db error updating account avatar description: %w", err)
			return gtserror.NewErrorInternalError(err)
		}
	}

	if form.HeaderDescription != nil {
		desc := text.StripHTMLFromText(*form.HeaderDescription)
		form.HeaderDescription = util.Ptr(desc)
	}

	if form.Header != nil && form.Header.Size != 0 {
		headerInfo, errWithCode := p.UpdateHeader(ctx,
			account,
			form.Header,
			form.HeaderDescription,
		)
		if errWithCode != nil {
			return errWithCode
		}
		account.HeaderMediaAttachmentID = headerInfo.ID
		account.HeaderMediaAttachment = headerInfo
		acctColumns = append(acctColumns, "header_media_attachment_id")
	} else if form.HeaderDescription != nil && account.HeaderMediaAttachment != nil {
		// Update just existing description if possible.
		account.HeaderMediaAttachment.Description = *form.HeaderDescription
		if err := p.state.DB.UpdateAttachment(
			ctx,
			account.HeaderMediaAttachment,
			"description",
		); err != nil {
			err := gtserror.Newf("db error updating account avatar description: %w", err)
			return gtserror.NewErrorInternalError(err)
		}
	}

	// Account settings flags and web stuff.

	if form.WebVisibility != nil {
		switch apimodel.Visibility(*form.WebVisibility) {

		// Show none.
		case apimodel.VisibilityNone:
			account.HidesToPublicFromUnauthedWeb = util.Ptr(true)
			account.HidesCcPublicFromUnauthedWeb = util.Ptr(true)

		// Show public only (GtS default).
		case apimodel.VisibilityPublic:
			account.HidesToPublicFromUnauthedWeb = util.Ptr(false)
			account.HidesCcPublicFromUnauthedWeb = util.Ptr(true)

		// Show public and unlisted (Masto default).
		case apimodel.VisibilityUnlisted:
			account.HidesToPublicFromUnauthedWeb = util.Ptr(false)
			account.HidesCcPublicFromUnauthedWeb = util.Ptr(false)

		default:
			const text = "web_visibility must be one of public, unlisted, or none"
			err := errors.New(text)
			return gtserror.NewErrorBadRequest(err, text)
		}

		acctColumns = append(
			acctColumns,
			"hides_to_public_from_unauthed_web",
			"hides_cc_public_from_unauthed_web",
		)
	}

	if form.Source != nil {
		if form.Source.Language != nil {
			language, err := validate.Language(*form.Source.Language)
			if err != nil {
				return gtserror.NewErrorBadRequest(err, err.Error())
			}

			account.Settings.Language = language
			settingsColumns = append(settingsColumns, "language")
		}

		if form.Source.Sensitive != nil {
			account.Settings.Sensitive = form.Source.Sensitive
			settingsColumns = append(settingsColumns, "sensitive")
		}

		if form.Source.Privacy != nil {
			if err := validate.Privacy(*form.Source.Privacy); err != nil {
				return gtserror.NewErrorBadRequest(err, err.Error())
			}

			priv := apimodel.Visibility(*form.Source.Privacy)
			account.Settings.Privacy = typeutils.APIVisToVis(priv)
			settingsColumns = append(settingsColumns, "privacy")
		}

		if form.Source.StatusContentType != nil {
			if err := validate.StatusContentType(*form.Source.StatusContentType); err != nil {
				return gtserror.NewErrorBadRequest(err, err.Error())
			}

			account.Settings.StatusContentType = *form.Source.StatusContentType
			settingsColumns = append(settingsColumns, "status_content_type")
		}
	}

	if form.Theme != nil {
		theme := *form.Theme
		if theme == "" {
			// Empty is easy, just clear this.
			account.Settings.Theme = ""
		} else {
			// Theme was provided, check
			// against known available themes.
			if _, ok := p.themes.ByFileName[theme]; !ok {
				err := fmt.Errorf("theme %s not available on this instance, see /api/v1/accounts/themes for available themes", theme)
				return gtserror.NewErrorBadRequest(err, err.Error())
			}
			account.Settings.Theme = theme
		}
		settingsColumns = append(settingsColumns, "theme")
	}

	if form.CustomCSS != nil {
		customCSS := *form.CustomCSS
		if err := validate.CustomCSS(customCSS); err != nil {
			return gtserror.NewErrorBadRequest(err, err.Error())
		}

		account.Settings.CustomCSS = text.StripHTMLFromText(customCSS)
		settingsColumns = append(settingsColumns, "custom_css")
	}

	if form.EnableRSS != nil {
		account.Settings.EnableRSS = form.EnableRSS
		settingsColumns = append(settingsColumns, "enable_rss")
	}

	if form.HideCollections != nil {
		account.Settings.HideCollections = form.HideCollections
		settingsColumns = append(settingsColumns, "hide_collections")
	}

	if form.WebLayout != nil {
		webLayout := gtsmodel.ParseWebLayout(*form.WebLayout)
		if webLayout == gtsmodel.WebLayoutUnknown {
			const text = "web_layout must be one of microblog or gallery"
			err := errors.New(text)
			return gtserror.NewErrorBadRequest(err, text)
		}

		account.Settings.WebLayout = webLayout
		settingsColumns = append(settingsColumns, "web_layout")
	}

	if form.WebIncludeBoosts != nil {
		account.Settings.WebIncludeBoosts = form.WebIncludeBoosts
		settingsColumns = append(settingsColumns, "web_include_boosts")
	}

	// We've parsed + set everything, do
	// necessary database updates now.

	if len(acctColumns) > 0 {
		if err := p.state.DB.UpdateAccount(ctx, account, acctColumns...); err != nil {
			err := gtserror.Newf("db error updating account %s: %w", account.ID, err)
			return gtserror.NewErrorInternalError(err)
		}
	}

	if len(settingsColumns) > 0 {
		if err := p.state.DB.UpdateAccountSettings(ctx, account.Settings, settingsColumns...); err != nil {
			err := gtserror.Newf("db error updating account settings %s: %w", account.ID, err)
			return gtserror.NewErrorInternalError(err)
		}
	}

	// Send out Update message over the s2s (fedi) API.
	p.state.Workers.Client.Queue.Push(&messages.FromClientAPI{
		APObjectType:   ap.ActorPerson,
		APActivityType: ap.ActivityUpdate,
		GTSModel:       account,
		Origin:         account,
	})

	return nil
}

func (p *Processor) selectNoteFormatter(contentType string) text.FormatFunc {
	if contentType == "text/markdown" {
		return p.formatter.FromMarkdown
	}

	return p.formatter.FromPlain
}

// ProcessAccountText processes the raw versions of the given
// account's display name, note, and fields, and sets those
// processed versions on the account, while also updating the
// account's emojis entry based on the results of the processing.
func (p *Processor) ProcessAccountText(
	ctx context.Context,
	account *gtsmodel.Account,
	contentType string,
) {
	// Use map to deduplicate emojis by their ID.
	emojis := make(map[string]*gtsmodel.Emoji)

	// Retrieve display name emojis.
	for _, emoji := range p.formatter.FromPlainBasic(
		ctx,
		p.parseMention,
		account.ID,
		"",
		account.DisplayName,
	).Emojis {
		emojis[emoji.ID] = emoji
	}

	// Format + set note according to prefs.
	f := p.selectNoteFormatter(contentType)
	formatNoteResult := f(ctx, p.parseMention, account.ID, "", account.NoteRaw)
	account.Note = formatNoteResult.HTML

	// Retrieve note emojis.
	for _, emoji := range formatNoteResult.Emojis {
		emojis[emoji.ID] = emoji
	}

	// Process raw fields.
	account.Fields = make([]*gtsmodel.Field, 0, len(account.FieldsRaw))
	for _, fieldRaw := range account.FieldsRaw {
		field := &gtsmodel.Field{}

		// Name stays plain, but we still need to
		// see if there are any emojis set in it.
		field.Name = fieldRaw.Name
		for _, emoji := range p.formatter.FromPlainBasic(
			ctx,
			p.parseMention,
			account.ID,
			"",
			fieldRaw.Name,
		).Emojis {
			emojis[emoji.ID] = emoji
		}

		// Value can be HTML, but we don't want
		// to wrap the result in <p> tags.
		fieldFormatValueResult := p.formatter.FromPlainNoParagraph(ctx, p.parseMention, account.ID, "", fieldRaw.Value)
		field.Value = fieldFormatValueResult.HTML

		// Retrieve field emojis.
		for _, emoji := range fieldFormatValueResult.Emojis {
			emojis[emoji.ID] = emoji
		}

		// We're done, append the shiny new field.
		account.Fields = append(account.Fields, field)
	}

	// Update the account's emojis.
	emojisCount := len(emojis)
	account.Emojis = make([]*gtsmodel.Emoji, 0, emojisCount)
	account.EmojiIDs = make([]string, 0, emojisCount)

	for id, emoji := range emojis {
		account.Emojis = append(account.Emojis, emoji)
		account.EmojiIDs = append(account.EmojiIDs, id)
	}
}

// UpdateAvatar does the dirty work of checking the avatar
// part of an account update form, parsing and checking the
// media, and doing the necessary updates in the database
// for this to become the account's new avatar.
func (p *Processor) UpdateAvatar(
	ctx context.Context,
	account *gtsmodel.Account,
	avatar *multipart.FileHeader,
	description *string,
) (
	*gtsmodel.MediaAttachment,
	gtserror.WithCode,
) {
	// Get maximum supported local media size.
	maxsz := config.GetMediaLocalMaxSize()
	maxszInt64 := int64(maxsz) // #nosec G115 -- Already validated.

	// Ensure media within size bounds.
	if avatar.Size > maxszInt64 {
		text := fmt.Sprintf("media exceeds configured max size: %s", maxsz)
		return nil, gtserror.NewErrorBadRequest(errors.New(text), text)
	}

	// Open multipart file reader.
	mpfile, err := avatar.Open()
	if err != nil {
		err := gtserror.Newf("error opening multipart file: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Wrap the multipart file reader to ensure is limited to max.
	rc, _, _ := iotools.UpdateReadCloserLimit(mpfile, maxszInt64)

	// Write to instance storage.
	return p.StoreLocalMedia(ctx,
		account.ID,
		func(ctx context.Context) (reader io.ReadCloser, err error) {
			return rc, nil
		},
		media.AdditionalMediaInfo{
			Avatar:      util.Ptr(true),
			Description: description,
		},
	)
}

// UpdateHeader does the dirty work of checking the header
// part of an account update form, parsing and checking the
// media, and doing the necessary updates in the database
// for this to become the account's new header.
func (p *Processor) UpdateHeader(
	ctx context.Context,
	account *gtsmodel.Account,
	header *multipart.FileHeader,
	description *string,
) (
	*gtsmodel.MediaAttachment,
	gtserror.WithCode,
) {
	// Get maximum supported local media size.
	maxsz := config.GetMediaLocalMaxSize()
	maxszInt64 := int64(maxsz) // #nosec G115 -- Already validated.

	// Ensure media within size bounds.
	if header.Size > maxszInt64 {
		text := fmt.Sprintf("media exceeds configured max size: %s", maxsz)
		return nil, gtserror.NewErrorBadRequest(errors.New(text), text)
	}

	// Open multipart file reader.
	mpfile, err := header.Open()
	if err != nil {
		err := gtserror.Newf("error opening multipart file: %w", err)
		return nil, gtserror.NewErrorInternalError(err)
	}

	// Wrap the multipart file reader to ensure is limited to max.
	rc, _, _ := iotools.UpdateReadCloserLimit(mpfile, maxszInt64)

	// Write to instance storage.
	return p.StoreLocalMedia(ctx,
		account.ID,
		func(ctx context.Context) (reader io.ReadCloser, err error) {
			return rc, nil
		},
		media.AdditionalMediaInfo{
			Header:      util.Ptr(true),
			Description: description,
		},
	)
}

// UpdateFields sets FieldsRaw on the given
// account, and resets account.Fields to an
// empty slice, ready for further processing.
func (p *Processor) UpdateFields(
	account *gtsmodel.Account,
	fieldsAttributes []apimodel.UpdateField,
) gtserror.WithCode {
	var (
		fieldsLen = len(fieldsAttributes)
		fieldsRaw = make([]*gtsmodel.Field, 0, fieldsLen)
	)

	for _, updateField := range fieldsAttributes {
		if updateField.Name == nil || updateField.Value == nil {
			continue
		}

		var (
			name  string = *updateField.Name
			value string = *updateField.Value
		)

		if name == "" || value == "" {
			continue
		}

		// Sanitize raw field values.
		fieldRaw := &gtsmodel.Field{
			Name:  text.StripHTMLFromText(name),
			Value: text.StripHTMLFromText(value),
		}
		fieldsRaw = append(fieldsRaw, fieldRaw)
	}

	// Check length of parsed raw fields.
	if err := validate.ProfileFields(fieldsRaw); err != nil {
		return gtserror.NewErrorBadRequest(err, err.Error())
	}

	// OK, new raw fields are valid.
	account.FieldsRaw = fieldsRaw
	account.Fields = make([]*gtsmodel.Field, 0, fieldsLen)
	return nil
}
