package coverage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/devyoujin/gococo/internal/utils"
	"github.com/stretchr/testify/suite"
)

type ManagerTestSuite struct {
	suite.Suite
	manager           ManagerInterface
	executor          *utils.MockCommandExecutor
	mergedCoverageDir string
	coverageProfile   string
	rootDir           string
}

func (suite *ManagerTestSuite) SetupSuite() {
	suite.executor = utils.NewMockCommandExecutor(suite.T())
	suite.rootDir = "test_data"
	suite.mergedCoverageDir = filepath.Join(suite.rootDir, "merged")
	suite.coverageProfile = filepath.Join(suite.rootDir, "coverage.out")
	suite.manager = NewManager(suite.executor, suite.rootDir, suite.mergedCoverageDir, suite.coverageProfile)
}

func (suite *ManagerTestSuite) SetupSubTest() {
	_ = os.Mkdir(suite.rootDir, 0755)
}

func (suite *ManagerTestSuite) TearDownSubTest() {
	_ = os.RemoveAll(suite.rootDir)
}

func (suite *ManagerTestSuite) TestFindModules() {
	type testCase struct {
		name        string
		givenRunner func()
		thenRunner  func(modules []module, err error)
	}
	testCases := []testCase{
		{
			name: "should find one module when there one module exists",
			givenRunner: func() {
				pathA := filepath.Join(suite.rootDir, "serviceA")
				_ = os.MkdirAll(pathA, 0755)
				_ = os.WriteFile(filepath.Join(pathA, "go.mod"), []byte("module serviceA"), 0644)
			},
			thenRunner: func(modules []module, err error) {
				suite.NoError(err)
				suite.Equal(1, len(modules))
				suite.Equal("serviceA", modules[0].name)
			},
		},
		{
			name: "should find two modules when two modules exist",
			givenRunner: func() {
				pathA := filepath.Join(suite.rootDir, "serviceA")
				_ = os.MkdirAll(pathA, 0755)
				_ = os.WriteFile(filepath.Join(pathA, "go.mod"), []byte("module serviceA"), 0644)
				pathB := filepath.Join(suite.rootDir, "serviceB")
				_ = os.MkdirAll(pathB, 0755)
				_ = os.WriteFile(filepath.Join(pathB, "go.mod"), []byte("module serviceB"), 0644)
			},
			thenRunner: func(modules []module, err error) {
				suite.NoError(err)
				suite.Equal(2, len(modules))
				suite.Equal("serviceA", modules[0].name)
				suite.Equal("serviceB", modules[1].name)
			},
		},
		{
			name: "should find one module when one module and one folder given",
			givenRunner: func() {
				pathA := filepath.Join(suite.rootDir, "serviceA")
				_ = os.MkdirAll(pathA, 0755)
				_ = os.WriteFile(filepath.Join(pathA, "go.mod"), []byte("module serviceA"), 0644)
				pathB := filepath.Join(suite.rootDir, "serviceB")
				_ = os.MkdirAll(pathB, 0755)
			},
			thenRunner: func(modules []module, err error) {
				suite.NoError(err)
				suite.Equal(1, len(modules))
				suite.Equal("serviceA", modules[0].name)
			},
		},
		{
			name: "should return empty when no modules exist",
			givenRunner: func() {
			},
			thenRunner: func(modules []module, err error) {
				suite.NoError(err)
				suite.Equal(0, len(modules))
			},
		},
	}

	for _, test := range testCases {
		suite.Run(test.name, func() {
			test.givenRunner()
			modules, err := suite.manager.FindGoModules()
			test.thenRunner(modules, err)
		})
	}
}

func (suite *ManagerTestSuite) TestGenerateCoverages() {
	type testCase struct {
		name          string
		mockRunner    func()
		modules       []module
		expectedError error
	}
	testCases := []testCase{
		{
			name: "should generate coverage for one module",
			mockRunner: func() {
				cmd := exec.Command("go", "test", "-cover", "./...", "-test.gocoverdir="+suite.mergedCoverageDir)
				cmd.Dir = filepath.Join(suite.rootDir, "serviceA")
				suite.executor.EXPECT().Run(cmd).Return(nil).Once()
			},
			modules: []module{
				{name: "serviceA", path: filepath.Join(suite.rootDir, "serviceA")},
			},
			expectedError: nil,
		},
		{
			name: "should generate coverages for two modules",
			mockRunner: func() {
				cmd := exec.Command("go", "test", "-cover", "./...", "-test.gocoverdir="+suite.mergedCoverageDir)
				cmd.Dir = filepath.Join(suite.rootDir, "serviceA")
				suite.executor.EXPECT().Run(cmd).Return(nil).Once()
				cmd = exec.Command("go", "test", "-cover", "./...", "-test.gocoverdir="+suite.mergedCoverageDir)
				cmd.Dir = filepath.Join(suite.rootDir, "serviceB")
				suite.executor.EXPECT().Run(cmd).Return(nil).Once()
			},
			modules: []module{
				{name: "serviceA", path: filepath.Join(suite.rootDir, "serviceA")},
				{name: "serviceB", path: filepath.Join(suite.rootDir, "serviceB")},
			},
			expectedError: nil,
		},
		{
			name:          "should return no error when no modules exist",
			mockRunner:    func() {},
			modules:       nil,
			expectedError: nil,
		},
		{
			name: "should return error when failed to run command",
			mockRunner: func() {
				cmd := exec.Command("go", "test", "-cover", "./...", "-test.gocoverdir="+suite.mergedCoverageDir)
				cmd.Dir = filepath.Join(suite.rootDir, "serviceA")
				suite.executor.EXPECT().Run(cmd).Return(fmt.Errorf("failed")).Once()
			},
			modules: []module{
				{name: "serviceA", path: filepath.Join(suite.rootDir, "serviceA")},
			},
			expectedError: fmt.Errorf("failed to run tests in %s: failed", filepath.Join(suite.rootDir, "serviceA")),
		},
	}

	for _, test := range testCases {
		suite.Run(test.name, func() {
			test.mockRunner()
			err := suite.manager.GenerateCoverages(test.modules)
			if test.expectedError != nil {
				suite.Equal(err.Error(), test.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *ManagerTestSuite) TestGenerateCoverProfile() {
	type testCase struct {
		name          string
		mockRunner    func()
		expectedError error
	}
	testCases := []testCase{
		{
			name: "should generate coverage profile",
			mockRunner: func() {
				cmd := exec.Command("go", "tool", "covdata", "textfmt", "-i="+suite.mergedCoverageDir, "-o="+suite.coverageProfile)
				suite.executor.EXPECT().Run(cmd).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "should return error when failed to run command",
			mockRunner: func() {
				cmd := exec.Command("go", "tool", "covdata", "textfmt", "-i="+suite.mergedCoverageDir, "-o="+suite.coverageProfile)
				suite.executor.EXPECT().Run(cmd).Return(fmt.Errorf("failed")).Once()
			},
			expectedError: fmt.Errorf("failed to generate coverage profile: failed"),
		},
	}

	for _, test := range testCases {
		suite.Run(test.name, func() {
			test.mockRunner()
			err := suite.manager.GenerateCoverProfile()
			if test.expectedError != nil {
				suite.Equal(err.Error(), test.expectedError.Error())
			} else {
				suite.NoError(err)
			}
		})
	}
}

func TestManagerTestSuite(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
