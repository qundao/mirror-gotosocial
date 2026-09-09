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

package media_test

import (
	"io"
	"net/http"
	"path"
	"testing"

	apimodel "code.superseriousbusiness.org/gotosocial/internal/api/model"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type GetFileTestSuite struct {
	MediaStandardTestSuite
}

func (suite *GetFileTestSuite) TestGetRemoteFileCached() {
	ctx := suite.T().Context()

	testAttachment := suite.testAttachments["remote_account_1_status_1_attachment_1"]
	fileName := path.Base(testAttachment.File.Path)
	requestingAccount := suite.testAccounts["local_account_1"]

	content, errWithCode := suite.mediaProcessor.GetFile(ctx, requestingAccount, &apimodel.GetContentRequestForm{
		AccountID: testAttachment.AccountID,
		MediaType: string(media.TypeAttachment),
		MediaSize: string(media.SizeOriginal),
		FileName:  fileName,
	})

	suite.NoError(errWithCode)
	suite.NotNil(content)
	b, err := io.ReadAll(content.Content)
	suite.NoError(err)

	if closer, ok := content.Content.(io.Closer); ok {
		suite.NoError(closer.Close())
	}

	suite.Equal(suite.testRemoteAttachments[testAttachment.RemoteURL].Data, b)
	suite.Equal(suite.testRemoteAttachments[testAttachment.RemoteURL].ContentType, content.ContentType)
	suite.EqualValues(len(suite.testRemoteAttachments[testAttachment.RemoteURL].Data), content.ContentLength)
}

func (suite *GetFileTestSuite) TestGetRemoteFileUncached() {
	ctx := suite.T().Context()

	// uncache the file from local
	testAttachment := suite.testAttachments["remote_account_1_status_1_attachment_1"]

	err := suite.storage.Delete(ctx, testAttachment.File.Path)
	suite.NoError(err)
	err = suite.storage.Delete(ctx, testAttachment.Thumbnail.Path)
	suite.NoError(err)

	fileName := path.Base(testAttachment.File.Path)
	testAttachment.File.Path = ""
	testAttachment.Thumbnail.Path = ""

	err = suite.db.UpdateAttachment(ctx, testAttachment, "thumbnail_path", "file_path")
	suite.NoError(err)

	// now fetch it
	requestingAccount := suite.testAccounts["local_account_1"]
	content, errWithCode := suite.mediaProcessor.GetFile(ctx, requestingAccount, &apimodel.GetContentRequestForm{
		AccountID: testAttachment.AccountID,
		MediaType: string(media.TypeAttachment),
		MediaSize: string(media.SizeOriginal),
		FileName:  fileName,
	})
	suite.NoError(errWithCode)
	suite.NotNil(content)

	b, err := io.ReadAll(content.Content)
	suite.NoError(err)
	suite.NoError(content.Content.Close())

	suite.Equal(suite.testRemoteAttachments[testAttachment.RemoteURL].Data, b)
	suite.Equal(suite.testRemoteAttachments[testAttachment.RemoteURL].ContentType, content.ContentType)
	suite.EqualValues(len(suite.testRemoteAttachments[testAttachment.RemoteURL].Data), content.ContentLength)

	// the attachment should be updated in the database
	var dbAttachment *gtsmodel.MediaAttachment
	if !testrig.WaitFor(func() bool {
		dbAttachment, err = suite.db.GetAttachmentByID(ctx, testAttachment.ID)
		return dbAttachment != nil
	}) {
		suite.FailNow("timed out waiting for updated attachment")
	}

	suite.NoError(err)
	suite.True(dbAttachment.Cached())

	// the file should be back in storage at the same path as before
	refreshedBytes, err := suite.storage.Get(ctx, dbAttachment.File.Path)
	suite.NoError(err)
	suite.Equal(suite.testRemoteAttachments[testAttachment.RemoteURL].Data, refreshedBytes)
}

func (suite *GetFileTestSuite) TestGetRemoteFileUncachedInterrupted() {
	ctx := suite.T().Context()

	// uncache the file from local
	testAttachment := suite.testAttachments["remote_account_1_status_1_attachment_1"]

	err := suite.storage.Delete(ctx, testAttachment.File.Path)
	suite.NoError(err)
	err = suite.storage.Delete(ctx, testAttachment.Thumbnail.Path)
	suite.NoError(err)

	fileName := path.Base(testAttachment.File.Path)
	testAttachment.File.Path = ""
	testAttachment.Thumbnail.Path = ""

	err = suite.db.UpdateAttachment(ctx, testAttachment, "thumbnail_path", "file_path")
	suite.NoError(err)

	// now fetch it
	requestingAccount := suite.testAccounts["local_account_1"]
	content, errWithCode := suite.mediaProcessor.GetFile(ctx, requestingAccount, &apimodel.GetContentRequestForm{
		AccountID: testAttachment.AccountID,
		MediaType: string(media.TypeAttachment),
		MediaSize: string(media.SizeOriginal),
		FileName:  fileName,
	})
	suite.NoError(errWithCode)
	suite.NotNil(content)

	_, err = io.CopyN(io.Discard, content.Content, 1024)
	suite.NoError(err)

	err = content.Content.Close()
	suite.NoError(err)

	// the attachment should still be updated in the database even though the caller hung up
	var dbAttachment *gtsmodel.MediaAttachment
	if !testrig.WaitFor(func() bool {
		dbAttachment, _ = suite.db.GetAttachmentByID(ctx, testAttachment.ID)
		return dbAttachment.Cached()
	}) {
		suite.FailNow("timed out waiting for attachment to be updated")
	}

	// the file should be back in storage at the same path as before
	refreshedBytes, err := suite.storage.Get(ctx, dbAttachment.File.Path)
	suite.NoError(err)
	suite.Equal(suite.testRemoteAttachments[testAttachment.RemoteURL].Data, refreshedBytes)
}

func (suite *GetFileTestSuite) TestGetRemoteFileThumbnailUncached() {
	ctx := suite.T().Context()
	testAttachment := suite.testAttachments["remote_account_1_status_1_attachment_1"]

	// fetch the existing thumbnail bytes from storage first
	thumbnailBytes, err := suite.storage.Get(ctx, testAttachment.Thumbnail.Path)
	suite.NoError(err)

	// uncache the file from local
	err = suite.storage.Delete(ctx, testAttachment.File.Path)
	suite.NoError(err)
	err = suite.storage.Delete(ctx, testAttachment.Thumbnail.Path)
	suite.NoError(err)

	fileName := path.Base(testAttachment.File.Path)
	testAttachment.File.Path = ""
	testAttachment.Thumbnail.Path = ""

	err = suite.db.UpdateAttachment(ctx, testAttachment, "thumbnail_path", "file_path")
	suite.NoError(err)

	// now fetch the thumbnail
	requestingAccount := suite.testAccounts["local_account_1"]
	content, errWithCode := suite.mediaProcessor.GetFile(ctx, requestingAccount, &apimodel.GetContentRequestForm{
		AccountID: testAttachment.AccountID,
		MediaType: string(media.TypeAttachment),
		MediaSize: string(media.SizeSmall),
		FileName:  fileName,
	})
	suite.NoError(errWithCode)
	suite.NotNil(content)

	b, err := io.ReadAll(content.Content)
	suite.NoError(err)
	suite.NoError(content.Content.Close())

	suite.Equal("image/jpeg", content.ContentType)
	testrig.AssertSimilarImages(suite.T(), b, thumbnailBytes)
}

func (suite *GetFileTestSuite) TestGetRemoteFileRejectMedia() {
	ctx := suite.T().Context()

	testAttachment := suite.testAttachments["remote_account_1_status_1_attachment_1"]
	fileName := path.Base(testAttachment.File.Path)

	// Update the existing domain limit on
	// fossbros-anonymous.io to reject media.
	limit := new(gtsmodel.DomainLimit)
	*limit = *suite.testDomainLimits["fossbros-anonymous.io"]
	limit.MediaPolicy = gtsmodel.MediaPolicyReject
	if err := suite.state.DB.UpdateDomainLimit(ctx, limit, "media_policy"); err != nil {
		suite.FailNow(err.Error())
	}

	// Try to get the content. This should return a 404 since
	// we're now rejecting media from fossbros-anonymous.io.
	content, errWithCode := suite.mediaProcessor.GetFile(ctx, nil, &apimodel.GetContentRequestForm{
		AccountID: testAttachment.AccountID,
		MediaType: string(media.TypeAttachment),
		MediaSize: string(media.SizeOriginal),
		FileName:  fileName,
	})

	suite.EqualError(errWithCode, "GetFile: rejecting media from fossbros-anonymous.io")
	suite.Equal(http.StatusNotFound, errWithCode.Code())
	suite.Equal("Not Found", errWithCode.Safe())
	suite.Nil(content)
}

func TestGetFileTestSuite(t *testing.T) {
	suite.Run(t, &GetFileTestSuite{})
}
