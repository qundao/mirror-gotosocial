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

package dereferencing_test

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"
	"code.superseriousbusiness.org/gotosocial/internal/ap"
	"code.superseriousbusiness.org/gotosocial/internal/config"
	"code.superseriousbusiness.org/gotosocial/internal/db"
	"code.superseriousbusiness.org/gotosocial/internal/federation/dereferencing"
	"code.superseriousbusiness.org/gotosocial/internal/filter/interaction"
	"code.superseriousbusiness.org/gotosocial/internal/filter/relay"
	"code.superseriousbusiness.org/gotosocial/internal/filter/visibility"
	"code.superseriousbusiness.org/gotosocial/internal/gtserror"
	"code.superseriousbusiness.org/gotosocial/internal/media"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type AccountTestSuite struct {
	DereferencerTestSuite
}

func (suite *AccountTestSuite) TestDereferenceGroup() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]

	groupURL := testrig.URLMustParse("https://unknown-instance.com/groups/some_group")
	group, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(
		ctx,
		fetchingAccount.Username,
		groupURL,
		false,
	)
	suite.NoError(err)
	suite.NotNil(group)

	// group values should be set
	suite.Equal("https://unknown-instance.com/groups/some_group", group.URI)
	suite.Equal("https://unknown-instance.com/@some_group", group.URL)
	suite.WithinDuration(time.Now(), group.FetchedAt, 5*time.Second)

	// group should be in the database
	dbGroup, err := testStructs.State.DB.GetAccountByURI(ctx, group.URI)
	suite.NoError(err)
	suite.Equal(group.ID, dbGroup.ID)
	suite.Equal(ap.ActorGroup, dbGroup.ActorType.String())
}

func (suite *AccountTestSuite) TestDereferenceService() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]

	serviceURL := testrig.URLMustParse("https://owncast.example.org/federation/user/rgh")
	service, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(
		ctx,
		fetchingAccount.Username,
		serviceURL,
		false,
	)
	suite.NoError(err)
	suite.NotNil(service)

	// service values should be set
	suite.Equal("https://owncast.example.org/federation/user/rgh", service.URI)
	suite.Equal("https://owncast.example.org/federation/user/rgh", service.URL)
	suite.WithinDuration(time.Now(), service.FetchedAt, 5*time.Second)

	// service should be in the database
	dbService, err := testStructs.State.DB.GetAccountByURI(ctx, service.URI)
	suite.NoError(err)
	suite.Equal(service.ID, dbService.ID)
	suite.Equal(ap.ActorService, dbService.ActorType.String())
	suite.Equal("example.org", dbService.Domain)
}

/*
	We shouldn't try webfingering or making http calls to dereference local accounts
	that might be passed into GetRemoteAccount for whatever reason, so these tests are
	here to make sure that such cases are (basically) short-circuit evaluated and given
	back as-is without trying to make any calls to one's own instance.
*/

func (suite *AccountTestSuite) TestDereferenceLocalAccountAsRemoteURL() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]
	// Account being dereferenced.
	targetAccount := suite.testAccounts["local_account_2"]

	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(
		ctx,
		fetchingAccount.Username,
		testrig.URLMustParse(targetAccount.URI),
		false,
	)
	suite.NoError(err)
	suite.NotNil(fetchedAccount)
	suite.Empty(fetchedAccount.Domain)
}

func (suite *AccountTestSuite) TestDereferenceLocalAccountAsRemoteURLNoSharedInboxYet() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]
	// Account being dereferenced.
	targetAccount := suite.testAccounts["local_account_2"]

	targetAccount.SharedInboxURI = nil
	if err := testStructs.State.DB.UpdateAccount(ctx, targetAccount); err != nil {
		suite.FailNow(err.Error())
	}

	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(
		ctx,
		fetchingAccount.Username,
		testrig.URLMustParse(targetAccount.URI),
		false,
	)
	suite.NoError(err)
	suite.NotNil(fetchedAccount)
	suite.Empty(fetchedAccount.Domain)
}

func (suite *AccountTestSuite) TestDereferenceLocalAccountAsUsername() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]
	// Account being dereferenced.
	targetAccount := suite.testAccounts["local_account_2"]

	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(
		ctx,
		fetchingAccount.Username,
		testrig.URLMustParse(targetAccount.URI),
		false,
	)
	suite.NoError(err)
	suite.NotNil(fetchedAccount)
	suite.Empty(fetchedAccount.Domain)
}

func (suite *AccountTestSuite) TestDereferenceLocalAccountAsUsernameDomain() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]
	// Account being dereferenced.
	targetAccount := suite.testAccounts["local_account_2"]

	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(
		ctx,
		fetchingAccount.Username,
		testrig.URLMustParse(targetAccount.URI),
		false,
	)
	suite.NoError(err)
	suite.NotNil(fetchedAccount)
	suite.Empty(fetchedAccount.Domain)
}

func (suite *AccountTestSuite) TestDereferenceLocalAccountAsUsernameDomainAndURL() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]
	// Account being dereferenced.
	targetAccount := suite.testAccounts["local_account_2"]

	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByUsernameDomain(
		ctx,
		fetchingAccount.Username,
		targetAccount.Username,
		config.GetHost(),
	)
	suite.NoError(err)
	suite.NotNil(fetchedAccount)
	suite.Empty(fetchedAccount.Domain)
}

func (suite *AccountTestSuite) TestDereferenceLocalAccountWithUnknownUsername() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]

	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByUsernameDomain(
		ctx,
		fetchingAccount.Username,
		"thisaccountdoesnotexist",
		config.GetHost(),
	)
	suite.True(gtserror.IsUnretrievable(err))
	suite.EqualError(err, db.ErrNoEntries.Error())
	suite.Nil(fetchedAccount)
}

func (suite *AccountTestSuite) TestDereferenceLocalAccountWithUnknownUsernameDomain() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]

	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByUsernameDomain(
		ctx,
		fetchingAccount.Username,
		"thisaccountdoesnotexist",
		"localhost:8080",
	)
	suite.True(gtserror.IsUnretrievable(err))
	suite.EqualError(err, db.ErrNoEntries.Error())
	suite.Nil(fetchedAccount)
}

func (suite *AccountTestSuite) TestDereferenceLocalAccountWithUnknownUserURI() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]

	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(
		ctx,
		fetchingAccount.Username,
		testrig.URLMustParse("http://localhost:8080/users/thisaccountdoesnotexist"),
		false,
	)
	suite.True(gtserror.IsUnretrievable(err))
	suite.EqualError(err, db.ErrNoEntries.Error())
	suite.Nil(fetchedAccount)
}

func (suite *AccountTestSuite) TestDereferenceLocalAccountByRedirect() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]
	// Account being dereferenced.
	targetAccount := suite.testAccounts["local_account_2"]

	// Convert the target account to ActivityStreams model for dereference.
	targetAccountable, err := testStructs.TypeConverter.AccountToAS(ctx, targetAccount)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Serialize to "raw" JSON map for response.
	raw, err := ap.Serialize(targetAccountable)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Finally serialize to actual bytes.
	json := testrig.MustJSONBytes(raw)

	// Replace test HTTP client with one that always returns the target account AS model.
	testStructs.HTTPClient = testrig.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			Status:        http.StatusText(http.StatusOK),
			StatusCode:    http.StatusOK,
			ContentLength: int64(len(json)),
			Header:        http.Header{"Content-Type": {"application/activity+json"}},
			Body:          io.NopCloser(bytes.NewReader(json)),
			Request:       &http.Request{URL: testrig.URLMustParse(targetAccount.URI)},
		}, nil
	}, "")

	// Update dereferencer to use new test HTTP client.
	testStructs.Federator.Dereferencer = dereferencing.NewDereferencer(
		testStructs.State,
		testStructs.TypeConverter,
		testrig.NewTestTransportController(testStructs.State, testStructs.HTTPClient),
		visibility.NewFilter(testStructs.State),
		interaction.NewFilter(testStructs.State),
		relay.NewFilter(testStructs.State),
		media.NewManager(testStructs.State),
	)

	// Use any old input test URI, this doesn't actually matter what it is.
	uri := testrig.URLMustParse("https://this-will-be-redirected.butts/")

	// Try dereference the test URI, since it correctly redirects to us it should return our account.
	account, accountable, err := testStructs.Federator.Dereferencer.GetAccountByURI(ctx, fetchingAccount.Username, uri, false)
	suite.NoError(err)
	suite.Nil(accountable)
	suite.NotNil(account)
	suite.Equal(targetAccount.ID, account.ID)
}

func (suite *AccountTestSuite) TestDereferenceMasqueradingLocalAccount() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]
	// Account being dereferenced.
	targetAccount := suite.testAccounts["local_account_2"]

	// Convert the target account to ActivityStreams model for dereference.
	targetAccountable, err := testStructs.TypeConverter.AccountToAS(ctx, targetAccount)
	if err != nil {
		suite.FailNow(err.Error())
	}

	// Serialize to "raw" JSON map for response.
	raw, err := ap.Serialize(targetAccountable)
	suite.NoError(err)

	// Finally serialize to actual bytes.
	json := testrig.MustJSONBytes(raw)

	// Use any old input test URI, this doesn't actually matter what it is.
	uri := testrig.URLMustParse("https://this-will-be-redirected.butts/")

	// Replace test HTTP client with one that returns OUR account, but at their URI endpoint.
	testStructs.HTTPClient = testrig.NewMockHTTPClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			Status:        http.StatusText(http.StatusOK),
			StatusCode:    http.StatusOK,
			ContentLength: int64(len(json)),
			Header:        http.Header{"Content-Type": {"application/activity+json"}},
			Body:          io.NopCloser(bytes.NewReader(json)),
			Request:       &http.Request{URL: uri},
		}, nil
	}, "")

	// Update dereferencer to use new test HTTP client.
	testStructs.Federator.Dereferencer = dereferencing.NewDereferencer(
		testStructs.State,
		testStructs.TypeConverter,
		testrig.NewTestTransportController(testStructs.State, testStructs.HTTPClient),
		visibility.NewFilter(testStructs.State),
		interaction.NewFilter(testStructs.State),
		relay.NewFilter(testStructs.State),
		media.NewManager(testStructs.State),
	)

	// Try dereference the test URI, since it correctly redirects to us it should return our account.
	account, accountable, err := testStructs.Federator.Dereferencer.GetAccountByURI(ctx, fetchingAccount.Username, uri, false)
	suite.NotNil(err)
	suite.Nil(account)
	suite.Nil(accountable)
}

func (suite *AccountTestSuite) TestDereferenceRemoteAccountWithNonMatchingURI() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Account to dereference on behalf of.
	fetchingAccount := suite.testAccounts["local_account_1"]

	const (
		remoteURI    = "https://turnip.farm/users/turniplover6969"
		remoteAltURI = "https://turnip.farm/users/turniphater420"
	)

	// Create a copy of this remote account at alternative URI.
	remotePerson := testStructs.HTTPClient.TestRemotePeople[remoteURI]
	testStructs.HTTPClient.TestRemotePeople[remoteAltURI] = remotePerson

	// Attempt to fetch account at alternative URI, it should fail!
	fetchedAccount, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(
		ctx,
		fetchingAccount.Username,
		testrig.URLMustParse(remoteAltURI),
		false,
	)
	suite.Equal(err.Error(), fmt.Sprintf("enrichAccount: account uri %s does not match %s", remoteURI, remoteAltURI))
	suite.Nil(fetchedAccount)
}

func (suite *AccountTestSuite) TestDereferenceRemoteAccountWithUnexpectedKeyChange() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	fetchingAcc := suite.testAccounts["local_account_1"]
	remoteURI := "https://turnip.farm/users/turniplover6969"

	// Fetch the remote account to load into the database.
	remoteAcc, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(ctx,
		fetchingAcc.Username,
		testrig.URLMustParse(remoteURI),
		false,
	)
	suite.NoError(err)
	suite.NotNil(remoteAcc)

	// Mark account as requiring a refetch.
	remoteAcc.FetchedAt = time.Time{}
	err = testStructs.State.DB.UpdateAccount(ctx, remoteAcc, "fetched_at")
	suite.NoError(err)

	// Update remote to have an unexpected different key.
	remotePerson := testStructs.HTTPClient.TestRemotePeople[remoteURI]
	setPublicKey(remotePerson,
		remoteURI,
		fetchingAcc.PublicKeyURI+".unique",
		fetchingAcc.PublicKey,
	)

	// Force refresh account expecting key change error.
	_, _, err = testStructs.Federator.Dereferencer.RefreshAccount(ctx,
		fetchingAcc.Username,
		remoteAcc,
		nil,
		nil,
	)
	suite.Equal(err.Error(), fmt.Sprintf("RefreshAccount: enrichAccount: account %s pubkey has changed (key rotation required?)", remoteURI))
}

func (suite *AccountTestSuite) TestDereferenceRemoteAccountWithExpectedKeyChange() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	fetchingAcc := suite.testAccounts["local_account_1"]
	remoteURI := "https://turnip.farm/users/turniplover6969"

	// Fetch the remote account to load into the database.
	remoteAcc, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(ctx,
		fetchingAcc.Username,
		testrig.URLMustParse(remoteURI),
		false,
	)
	suite.NoError(err)
	suite.NotNil(remoteAcc)

	// Expire the remote account's public key.
	remoteAcc.PublicKeyExpiresAt = time.Now()
	remoteAcc.FetchedAt = time.Time{} // force fetch
	err = testStructs.State.DB.UpdateAccount(ctx, remoteAcc, "fetched_at", "public_key_expires_at")
	suite.NoError(err)

	// Update remote to have a different stored public key.
	remotePerson := testStructs.HTTPClient.TestRemotePeople[remoteURI]
	setPublicKey(remotePerson,
		remoteURI,
		fetchingAcc.PublicKeyURI+".unique",
		fetchingAcc.PublicKey,
	)

	// Refresh account expecting a succesful refresh with changed keys!
	updatedAcc, apAcc, err := testStructs.Federator.Dereferencer.RefreshAccount(ctx,
		fetchingAcc.Username,
		remoteAcc,
		nil,
		nil,
	)
	suite.NoError(err)
	suite.NotNil(apAcc)
	suite.True(updatedAcc.PublicKey.Equal(fetchingAcc.PublicKey))
}

func (suite *AccountTestSuite) TestRefreshFederatedRemoteAccountWithKeyChange() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	fetchingAcc := suite.testAccounts["local_account_1"]
	remoteURI := "https://turnip.farm/users/turniplover6969"

	// Fetch the remote account to load into the database.
	remoteAcc, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(ctx,
		fetchingAcc.Username,
		testrig.URLMustParse(remoteURI),
		false,
	)
	suite.NoError(err)
	suite.NotNil(remoteAcc)

	// Update remote to have a different stored public key.
	remotePerson := testStructs.HTTPClient.TestRemotePeople[remoteURI]
	setPublicKey(remotePerson,
		remoteURI,
		fetchingAcc.PublicKeyURI+".unique",
		fetchingAcc.PublicKey,
	)

	// Refresh account expecting a succesful refresh with changed keys!
	// By passing in the remote person model this indicates that the data
	// was received via the federator, which should trust any key change.
	updatedAcc, apAcc, err := testStructs.Federator.Dereferencer.RefreshAccount(ctx,
		fetchingAcc.Username,
		remoteAcc,
		remotePerson,
		nil,
	)
	suite.NoError(err)
	suite.NotNil(apAcc)
	suite.True(updatedAcc.PublicKey.Equal(fetchingAcc.PublicKey))
}

func (suite *AccountTestSuite) TestDereferenceRemoteAccountWithAvatarDescription() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	fetchingAcc := suite.testAccounts["local_account_1"]
	remoteURI := "https://shrimpnet.example.org/users/shrimp"
	description := "me scrolling fedi on a laptop, there's a monster ultra white and another fedi user on my right."

	// Fetch the remote account to load into the database.
	remoteAcc, _, err := testStructs.Federator.Dereferencer.GetAccountByURI(ctx,
		fetchingAcc.Username,
		testrig.URLMustParse(remoteURI),
		false,
	)
	suite.NoError(err)
	suite.NotNil(remoteAcc)

	suite.Equal(remoteAcc.AvatarMediaAttachment.Description, description)

	remotePerson := testStructs.HTTPClient.TestRemotePeople[remoteURI]

	description = strings.TrimSuffix(description, ".")

	icon := remotePerson.GetActivityStreamsIcon()
	image := icon.Begin().GetActivityStreamsImage()
	nameProp := streams.NewActivityStreamsNameProperty()
	nameProp.AppendXMLSchemaString(description)
	image.SetActivityStreamsName(nameProp)
	icon.SetActivityStreamsImage(0, image)
	remotePerson.SetActivityStreamsIcon(icon)

	updatedAcc, apAcc, err := testStructs.Federator.Dereferencer.RefreshAccount(ctx,
		fetchingAcc.Username,
		remoteAcc,
		remotePerson,
		nil,
	)
	suite.NoError(err)
	suite.NotNil(apAcc)

	// our account media fetches are
	// async, so wait until updated.
	testrig.WaitFor(func() bool {
		media, _ := testStructs.State.DB.GetAttachmentByID(ctx, updatedAcc.AvatarMediaAttachmentID)
		return media != nil && media.Description == description
	})
}

func TestAccountTestSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

func setPublicKey(person vocab.ActivityStreamsPerson, ownerURI, keyURI string, key *rsa.PublicKey) {
	profileIDURI, err := url.Parse(ownerURI)
	if err != nil {
		panic(err)
	}

	publicKeyURI, err := url.Parse(keyURI)
	if err != nil {
		panic(err)
	}

	publicKeyProp := streams.NewW3IDSecurityV1PublicKeyProperty()

	// create the public key
	publicKey := streams.NewW3IDSecurityV1PublicKey()

	// set ID for the public key
	publicKeyIDProp := streams.NewJSONLDIdProperty()
	publicKeyIDProp.SetIRI(publicKeyURI)
	publicKey.SetJSONLDId(publicKeyIDProp)

	// set owner for the public key
	publicKeyOwnerProp := streams.NewW3IDSecurityV1OwnerProperty()
	publicKeyOwnerProp.SetIRI(profileIDURI)
	publicKey.SetW3IDSecurityV1Owner(publicKeyOwnerProp)

	// set the pem key itself
	encodedPublicKey, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		panic(err)
	}
	publicKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: encodedPublicKey,
	})
	publicKeyPEMProp := streams.NewW3IDSecurityV1PublicKeyPemProperty()
	publicKeyPEMProp.Set(string(publicKeyBytes))
	publicKey.SetW3IDSecurityV1PublicKeyPem(publicKeyPEMProp)

	// append the public key to the public key property
	publicKeyProp.AppendW3IDSecurityV1PublicKey(publicKey)

	// set the public key property on the Person
	person.SetW3IDSecurityV1PublicKey(publicKeyProp)
}
