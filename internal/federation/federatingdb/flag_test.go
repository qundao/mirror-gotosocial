package federatingdb_test

import (
	"context"
	"encoding/json"
	"slices"
	"testing"

	"code.superseriousbusiness.org/activity/streams"
	"code.superseriousbusiness.org/activity/streams/vocab"
	"code.superseriousbusiness.org/gotosocial/internal/gtsmodel"
	"code.superseriousbusiness.org/gotosocial/internal/util"
	"code.superseriousbusiness.org/gotosocial/testrig"
	"github.com/stretchr/testify/suite"
)

type FlagTestSuite struct {
	FederatingDBTestSuite
}

func (suite *FlagTestSuite) TestCreateFlag() {
	// Set up our test structs + tear down on finish.
	testStructs := testrig.SetupTestStructs(rMediaPath, rTemplatePath)
	defer testrig.TearDownTestStructs(testStructs)

	// Clean up test context when done.
	ctx, cncl := context.WithCancel(suite.T().Context())
	defer cncl()

	// Create a misskey-style flag where the
	// reported status URI is inside the "content" field.
	reportedAccount := suite.testAccounts["local_account_1"]
	reportingAccount := suite.testAccounts["remote_account_1"]
	reportedStatus := suite.testStatuses["local_account_1_status_1"]
	reportContent := `Note: ` + reportedStatus.URL + `\n-----\nban this sick filth ⛔`
	const reportURI = "http://fossbros-anonymous.io/db22128d-884e-4358-9935-6a7c3940535d"

	raw := `{
  "@context": "https://www.w3.org/ns/activitystreams",
  "actor": "` + reportingAccount.URI + `",
  "content": "` + reportContent + `",
  "id": "` + reportURI + `",
  "object": "` + reportedAccount.URI + `",
  "type": "Flag"
}`

	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		suite.FailNow(err.Error())
	}

	t, err := streams.ToType(suite.T().Context(), m)
	if err != nil {
		suite.FailNow(err.Error())
	}

	flag := t.(vocab.ActivityStreamsFlag)
	if err := testStructs.Federator.FederatingDB().Flag(
		createTestContext(
			ctx,
			reportingAccount,
			reportedAccount,
		),
		flag,
	); err != nil {
		suite.FailNow(err.Error())
	}

	// Wait for report to
	// appear in the database.
	var report *gtsmodel.Report
	if !testrig.WaitFor(func() bool {
		// Get reports by reporting
		// account targeting reported account.
		reports, err := testStructs.State.DB.GetReports(ctx,
			util.Ptr(false),
			reportingAccount.ID,
			reportedAccount.ID,
			nil,
		)
		if err != nil {
			return false
		}

		// Get index of report
		// we're interested in.
		i := slices.IndexFunc(
			reports,
			func(r *gtsmodel.Report) bool {
				return r.URI == reportURI
			},
		)
		if i == -1 {
			return false
		}

		// Set the report.
		report = reports[i]
		return true
	}) {
		suite.FailNow("waiting for report")
	}

	// We should have properly extracted the
	// status ID from the misskey-format report.
	suite.Equal(reportedStatus.ID, report.StatusIDs[0])
}

func TestFlagTestSuite(t *testing.T) {
	suite.Run(t, new(FlagTestSuite))
}
