package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GetSubStringTest struct {
	Input    AcharyaEntity
	Expected string
}

type GetSubStringSuite struct {
	suite.Suite
	TestData        []GetSubStringTest
	TestInvalidData []GetSubStringTest
}

func (suite *GetSubStringSuite) SetupTest() {
	suite.TestData = []GetSubStringTest{
		{AcharyaEntity{0, 2, "Hello world"}, "He"},
		{AcharyaEntity{0, 2, "Hello world"}, "He"},
		{AcharyaEntity{0, 2, "Apple MacBook Pro (16-inch"}, "Ap"},
		{AcharyaEntity{0, 2, "नमस्ते दुनिया"}, "नम"},
		{AcharyaEntity{0, 2, "ਸਤਿ ਸ੍ਰੀ ਅਕਾਲ ਦੁਨਿਆ"}, "ਸਤ"},
		{AcharyaEntity{0, 2, "ഹലോ വേൾഡ്"}, "ഹല"},
		{AcharyaEntity{0, 2, "こんにちは世界"}, "こん"},
	}

	suite.TestInvalidData = []GetSubStringTest{
		{AcharyaEntity{-22, 10, "Negative start"}, ""},
		{AcharyaEntity{22, 10, "End pos smaller than start"}, ""},
		{AcharyaEntity{0, 500, "LENGTH OF THIS STRING IS SHORTER THAN END POS"}, ""},
	}
}

func (suite *GetSubStringSuite) TestValidData() {
	for _, v := range suite.TestData {
		subStr, err := GetSubString(v.Input.Name, v.Input.Begin, v.Input.End)
		suite.Nil(err, fmt.Sprintf("Error %v", v))
		suite.Equal(subStr, v.Expected, fmt.Sprintf("Error %v", v))
	}
}

func (suite *GetSubStringSuite) TestInValidData() {
	for _, v := range suite.TestInvalidData {
		subStr, err := GetSubString(v.Input.Name, v.Input.Begin, v.Input.End)
		suite.NotNil(err, fmt.Sprintf("Error %v", v))
		suite.Equal(subStr, v.Expected, fmt.Sprintf("Error %v", v))
	}
}

type GetEntitiesFromFileTest struct {
	Input    string
	Expected []string
}

type GetEntitiesFromFileSuite struct {
	suite.Suite
	TestData        []GetEntitiesFromFileTest
	TestDataInvalid []GetEntitiesFromFileTest
}

func (suite *GetEntitiesFromFileSuite) SetupTest() {
	suite.TestData = []GetEntitiesFromFileTest{
		{
			"./testData/news/annotation.conf",
			[]string{"Person", "Organization", "GPE", "Money"},
		},
		{
			"./testData/CoNLL-ST_2002/annotation.conf",
			[]string{"ORG", "PER", "LOC", "MISC"},
		},
	}
	suite.TestDataInvalid = []GetEntitiesFromFileTest{
		{
			"./testData/invalid-files/invalid-entities/annotation.conf",
			[]string{},
		},
	}
}

func (suite *GetEntitiesFromFileSuite) TestGetEntitiesFromFileValid() {
	for i, v := range suite.TestData {
		cDat, cErr := os.Open(suite.TestData[i].Input)
		suite.Nil(cErr)
		defer cDat.Close()
		output := GetEntitiesFromFile(cDat)
		suite.Equal(len(v.Expected), len(output))

		for _, ent := range v.Expected {
			_, ok := output[ent]
			suite.True(ok)
		}
	}
}

func (suite *GetEntitiesFromFileSuite) TestGetEntitiesFromFileInValid() {
	for i, v := range suite.TestDataInvalid {
		cDat, cErr := os.Open(suite.TestDataInvalid[i].Input)
		suite.Nil(cErr)
		defer cDat.Close()
		output := GetEntitiesFromFile(cDat)
		suite.Equal(len(v.Expected), len(output))
		for _, ent := range v.Expected {
			_, ok := output[ent]
			suite.False(ok)
		}
	}
}

type GeneratePathsTest struct {
	Input       string
	ExpectedAnn []string
	ExpectedTxt []string
}

type GeneratePathsSuite struct {
	suite.Suite
	TestData        []GeneratePathsTest
	TestDataInvalid []GeneratePathsTest
}

func (suite *GeneratePathsSuite) SetupTest() {
	suite.TestData = []GeneratePathsTest{
		{"./testData/CoNLL-ST_2002",
			[]string{"testData/CoNLL-ST_2002/esp/esp.train-doc-100.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-1400.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-423.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-46.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-503.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-896.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-919.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-978.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-989.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-118.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-123.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-134.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-181.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-184.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-236.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-251.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-27.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-46.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-75.ann"},
			[]string{"testData/CoNLL-ST_2002/esp/esp.train-doc-100.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-1400.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-423.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-46.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-503.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-896.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-919.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-978.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-989.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-118.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-123.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-134.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-181.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-184.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-236.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-251.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-27.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-46.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-75.txt"},
		},
	}

	suite.TestDataInvalid = []GeneratePathsTest{
		{"./testData/invalid-files/multiple-ann",
			[]string{},
			[]string{},
		},
		{"./testData/invalid-files/no-corresponding-ann",
			[]string{},
			[]string{},
		},
		{"./testData/invalid-files/no-corresponding-txt",
			[]string{},
			[]string{},
		},
	}
}

func (suite *GeneratePathsSuite) TestGetSubDirectories() {

	for _, v := range suite.TestData {
		anns, txts, err := GetSubDirectories(v.Input)
		suite.Nil(err)

		for i, ann := range anns {
			suite.Equal(v.ExpectedAnn[i], filepath.ToSlash(ann))
			suite.Equal(v.ExpectedTxt[i], filepath.ToSlash(txts[i]))
		}
	}
}

func (suite *GeneratePathsSuite) TestGetSubDirectoriesInvalid() {
	for _, v := range suite.TestDataInvalid {
		anns, txts, err := GetSubDirectories(v.Input)
		suite.Equal(anns, []string{})
		suite.NotNil(err)

		for i, ann := range anns {
			suite.Equal(v.ExpectedAnn[i], filepath.ToSlash(ann))
			suite.Equal(v.ExpectedTxt[i], filepath.ToSlash(txts[i]))
		}
	}
}

type GetTextAnnNoTest struct {
	Input    string
	Expected int
}

type GetTextAnnNoSuite struct {
	suite.Suite
	TestData        []GetTextAnnNoTest
	TestDataInvalid []GetTextAnnNoTest
}

func (suite *GetTextAnnNoSuite) SetupTest() {
	suite.TestData = []GetTextAnnNoTest{
		{
			"T1	Organization 0 4	Sony",
			1,
		},
		{
			"T671	Organization 0 4	Sony",
			671,
		},
	}
	suite.TestDataInvalid = []GetTextAnnNoTest{
		{
			"",
			0,
		},
		{
			"TINVALID	Organization 0 4	Sony",
			0,
		},
	}
}

func (suite *GetTextAnnNoSuite) TestGetTextAnnNo() {

	for _, v := range suite.TestData {
		number, err := GetTextAnnNum(v.Input)
		suite.Nil(err)
		suite.Equal(v.Expected, number)
	}

}

func (suite *GetTextAnnNoSuite) TestGetTextAnnNoInvalid() {

	for _, v := range suite.TestDataInvalid {
		number, err := GetTextAnnNum(v.Input)
		suite.NotNil(err, fmt.Sprintf("Error %v", v))
		suite.Equal(v.Expected, number, fmt.Sprintf("Error %v", v))
	}

}

type GenNumberEntityArrTest struct {
	Input struct {
		EntityMap   map[string]bool
		AnnFilePath string
	}
	Expected []NumberAcharyaEntity
}

type GenNumberEntityArrSuite struct {
	suite.Suite
	TestData        []GenNumberEntityArrTest
	TestDataInvalid []GenNumberEntityArrTest
}

func (suite *GenNumberEntityArrSuite) SetupTest() {

	entitiesMap := make(map[string]bool)
	entitiesMap["Organization"] = true
	entitiesMap["Money"] = true
	entitiesMap["Transfer-money"] = true
	entitiesMap["Person"] = true
	entitiesMap["GPE"] = true

	type TestInput struct {
		EntityMap   map[string]bool
		AnnFilePath string
	}

	suite.TestData = []GenNumberEntityArrTest{
		{
			Input: TestInput{
				entitiesMap,
				"./testData/news/000-introduction.ann",
			},
			Expected: []NumberAcharyaEntity{{TxtAnnNo: 1, Entity: AcharyaEntity{Begin: 418, End: 426, Name: "Organization"}}, {TxtAnnNo: 2, Entity: AcharyaEntity{Begin: 456, End: 468, Name: "Money"}}, {TxtAnnNo: 3, Entity: AcharyaEntity{Begin: 443, End: 449, Name: "Transfer-money"}}, {TxtAnnNo: 4, Entity: AcharyaEntity{Begin: 473, End: 496, Name: "Person"}}, {TxtAnnNo: 5, Entity: AcharyaEntity{Begin: 511, End: 535, Name: "Person"}}, {TxtAnnNo: 6, Entity: AcharyaEntity{Begin: 540, End: 545, Name: "Organization"}}, {TxtAnnNo: 7, Entity: AcharyaEntity{Begin: 549, End: 560, Name: "GPE"}}},
		},
	}

	suite.TestDataInvalid = []GenNumberEntityArrTest{
		{
			Input: TestInput{
				entitiesMap,
				"./testData/invalid-files/dicontinous-text-bound-annotations/030-login.ann",
			},
			Expected: []NumberAcharyaEntity{},
		},
		{
			Input: TestInput{
				entitiesMap,
				"./testData/invalid-files/invalid-ann-atoi/pos1.ann",
			},
			Expected: []NumberAcharyaEntity{},
		},
		{
			Input: TestInput{
				entitiesMap,
				"./testData/invalid-files/invalid-ann-atoi/pos2.ann",
			},
			Expected: []NumberAcharyaEntity{},
		},
		{
			Input: TestInput{
				entitiesMap,
				"./testData/invalid-files/invalid-ann-space/030-login.ann",
			},
			Expected: []NumberAcharyaEntity{},
		},
		{
			Input: TestInput{
				entitiesMap,
				"./testData/invalid-files/invalid-ann-tab/030-login.ann",
			},
			Expected: []NumberAcharyaEntity{},
		},
		{
			Input: TestInput{
				entitiesMap,
				"./testData/invalid-files/invalid-ann-tabNumber/030-login.ann",
			},
			Expected: []NumberAcharyaEntity{},
		},
	}

}

func (suite *GenNumberEntityArrSuite) TestGenNumberEntityArr() {
	for _, v := range suite.TestData {
		annFile, aErr := os.Open(v.Input.AnnFilePath)
		suite.Nil(aErr)

		noEnt, err := GenNumberEntityArr(v.Input.EntityMap, annFile)
		suite.Nil(err)
		suite.Equal(noEnt, v.Expected)
	}
}

func (suite *GenNumberEntityArrSuite) TestGenNumberEntityArrInvalid() {
	for _, v := range suite.TestDataInvalid {
		annFile, aErr := os.Open(v.Input.AnnFilePath)
		suite.Nil(aErr)

		noEnt, err := GenNumberEntityArr(v.Input.EntityMap, annFile)
		suite.NotNil(err)
		suite.Equal(noEnt, v.Expected)
	}
}

type GenerateAcharyaAndStandoffTest struct {
	Input struct {
		Data             string
		NumberAcharyaEnt []NumberAcharyaEntity
	}
	Expected struct {
		Acharya  string
		Standoff string
	}
}

type GenerateAcharyaAndStandoffSuite struct {
	suite.Suite
	TestData        []GenerateAcharyaAndStandoffTest
	TestDataInvalid []GenerateAcharyaAndStandoffTest
}

func (suite *GenerateAcharyaAndStandoffSuite) SetupTest() {

	type TestInput struct {
		Data             string
		NumberAcharyaEnt []NumberAcharyaEntity
	}

	type TestExpected struct {
		Acharya  string
		Standoff string
	}

	suite.TestData = []GenerateAcharyaAndStandoffTest{
		{Input: TestInput{
			"Hello world",
			[]NumberAcharyaEntity{{TxtAnnNo: 1, Entity: AcharyaEntity{Begin: 0, End: 1, Name: "Organization"}}},
		},
			Expected: TestExpected{"{\"Data\":\"Hello world\",\"Entities\":[[0,1,\"Organization\"]]}\n", "T1\tOrganization 0 1\tH"},
		},
	}

	suite.TestDataInvalid = []GenerateAcharyaAndStandoffTest{
		{Input: TestInput{
			"End position is larger",
			[]NumberAcharyaEntity{{TxtAnnNo: 1, Entity: AcharyaEntity{Begin: 0, End: 1000, Name: "Organization"}}},
		},
			Expected: TestExpected{"", ""},
		},
	}
}

func (suite *GenerateAcharyaAndStandoffSuite) TestGenerateAcharyaAndStandoff() {
	for _, v := range suite.TestData {
		acharya, standoff, err := GenerateAcharyaAndStandoff(v.Input.Data, v.Input.NumberAcharyaEnt)
		suite.Nil(err)
		suite.Equal(v.Expected.Acharya, acharya)
		suite.Equal(v.Expected.Standoff, standoff)
	}
}

func (suite *GenerateAcharyaAndStandoffSuite) TestGenerateAcharyaAndStandoffInvalid() {
	for _, v := range suite.TestDataInvalid {
		acharya, standoff, err := GenerateAcharyaAndStandoff(v.Input.Data, v.Input.NumberAcharyaEnt)
		suite.NotNil(err)
		suite.Equal(v.Expected.Acharya, acharya)
		suite.Equal(v.Expected.Standoff, standoff)
	}
}

type ValidateFlagsTest struct {
	Input struct {
		FPath     string
		AnnFiles  string
		TxtFiles  string
		ConfFile  string
		OFileName string
		OverWrite bool
	}
}

type ValidateFlagsSuite struct {
	suite.Suite
	TestData        []ValidateFlagsTest
	TestDataInvalid []ValidateFlagsTest
}

func (suite *ValidateFlagsSuite) SetupTest() {

	type TestInput struct {
		FPath     string
		AnnFiles  string
		TxtFiles  string
		ConfFile  string
		OFileName string
		OverWrite bool
	}

	suite.TestData = []ValidateFlagsTest{
		{Input: TestInput{"./testData/", "a.ann,b.ann", "a.txt,b.txt", "./annotation.conf", "OfileName", true}},
		{Input: TestInput{"", "a.ann,b.ann", "a.txt,b.txt", "./annotation.conf", "OfileName", true}},
	}

	suite.TestDataInvalid = []ValidateFlagsTest{
		{Input: TestInput{" ", "a.ann,b.ann", "a.txt,b.txt", "./annotation.conf", "OfileName", true}},
		{Input: TestInput{"./testData/", "a.ann,b.ann", "a.txt,b.txt", "./annotation.conf", "", true}},
		{Input: TestInput{"", "", "a.txt,b.txt", "./annotation.conf", "", true}},
		{Input: TestInput{"", "a.ann,b.ann", "", "./annotation.conf", "", true}},
		{Input: TestInput{"", "a.ann,b.ann", "a.txt,b.txt", "", "", true}},

		{Input: TestInput{"", "a.ann,b.ann", "a.txt", "./annotation.conf", "", true}},
		{Input: TestInput{"", "a.ann,b.ann", "a.txt,c.txt", "./annotation.conf", "", true}},
	}
}

func (suite *ValidateFlagsSuite) TestValidateFlags() {
	for _, v := range suite.TestData {
		err := ValidateFlags(v.Input.FPath, v.Input.AnnFiles, v.Input.TxtFiles, v.Input.ConfFile, v.Input.OFileName, v.Input.OverWrite)
		suite.Nil(err)
	}
}

func (suite *ValidateFlagsSuite) TestValidateFlagsInvalid() {
	for _, v := range suite.TestDataInvalid {
		err := ValidateFlags(v.Input.FPath, v.Input.AnnFiles, v.Input.TxtFiles, v.Input.ConfFile, v.Input.OFileName, v.Input.OverWrite)
		suite.NotNil(err)
	}
}

type HandleMainTest struct {
	Input struct {
		FPath     string
		AnnFiles  string
		TxtFiles  string
		ConfFile  string
		OFileName string
		OverWrite bool
	}
}

type HandleMainTestSuite struct {
	suite.Suite
	TestData        []HandleMainTest
	TestDataInvalid []HandleMainTest
}

func (suite *HandleMainTestSuite) SetupTest() {

	type TestInput struct {
		FPath     string
		AnnFiles  string
		TxtFiles  string
		ConfFile  string
		OFileName string
		OverWrite bool
	}

	suite.TestData = []HandleMainTest{
		{Input: TestInput{"./testData/news", "", "", "", "OfileName", true}},
		{Input: TestInput{"", "./testData/news/080-event-annotation.ann", "./testData/news/080-event-annotation.txt", "testData/news/annotation.conf", "OfileName", true}},
		{Input: TestInput{"", "./testData/news/080-event-annotation.ann", "./testData/news/080-event-annotation.txt", "testData/news/annotation.conf", "", true}},
	}

	suite.TestDataInvalid = []HandleMainTest{
		{Input: TestInput{" ", "a.ann,b.ann", "a.txt,b.txt", "./annotation.conf", "OfileName", true}},
		{Input: TestInput{"./testData/", "a.ann,b.ann", "a.txt,b.txt", "./annotation.conf", "", true}},
		{Input: TestInput{"", "", "a.txt,b.txt", "./annotation.conf", "", true}},
		{Input: TestInput{"", "a.ann,b.ann", "", "./annotation.conf", "", true}},
		{Input: TestInput{"", "a.ann,b.ann", "a.txt,b.txt", "", "", true}},

		{Input: TestInput{"", "a.ann,b.ann", "a.txt", "./annotation.conf", "", true}},
		{Input: TestInput{"", "a.ann,b.ann", "a.txt,c.txt", "./annotation.conf", "", true}},

		{Input: TestInput{"", "./testData/news/080-event-INVALID.ann", "./testData/news/080-event-annotation.txt", "testData/news/annotation.conf", "OfileName", true}},
		{Input: TestInput{"", "./testData/news/080-event-annotation.ann", "./testData/news/080-event-INVALID.txt", "testData/news/annotation.conf", "OfileName", true}},
		{Input: TestInput{"", "./testData/invalid-files/invalid-ann-space/030-login.ann", "./testData/news/080-event-annotation.txt", "testData/news/annotation.conf", "OfileName", true}},

		{Input: TestInput{"./testData/invalid-files/no-entities", "", "", "", "OfileName", true}},
	}
}

func (suite *HandleMainTestSuite) TestHandleMainTest() {
	for _, v := range suite.TestData {
		err := handleMain(v.Input.FPath, v.Input.AnnFiles, v.Input.TxtFiles, v.Input.ConfFile, v.Input.OFileName, v.Input.OverWrite)
		suite.Nil(err)
	}
}

func (suite *HandleMainTestSuite) TestHandleMainTestInvalid() {
	for _, v := range suite.TestDataInvalid {
		err := handleMain(v.Input.FPath, v.Input.AnnFiles, v.Input.TxtFiles, v.Input.ConfFile, v.Input.OFileName, v.Input.OverWrite)
		suite.NotNil(err, fmt.Sprint(v))
	}
}

func TestHandleOutput(t *testing.T) {
	e := os.Remove("./testData/file_gen_by_handleOutputTest.jsonl")
	assert.Nil(t, e)

	acharya := `{"Data":"Welcome to the Brat Rapid Annotation Tool","Entities":[[418,426,"Organization"]}` + "\n"

	err := handleOutput("./testData/file_gen_by_handleOutputTest.jsonl", acharya, false)
	assert.Nil(t, err)

	achDat, achErr := os.Open("./testData/file_gen_by_handleOutputTest.jsonl")
	assert.Nil(t, achErr)

	achData, err := ioutil.ReadAll(achDat)
	assert.Nil(t, err)

	assert.Equal(t, acharya, string(achData))

	err = handleOutput("./testData/file_gen_by_handleOutputTest.jsonl", acharya, false)
	assert.NotNil(t, err)

	err = handleOutput("./testData/file_gen_by_handleOutputTest.jsonl", acharya, true)
	assert.Nil(t, err, err)

}

func TestRunAllSuites(t *testing.T) {
	suite.Run(t, new(GetSubStringSuite))
	suite.Run(t, new(GetEntitiesFromFileSuite))
	suite.Run(t, new(GeneratePathsSuite))
	suite.Run(t, new(GetTextAnnNoSuite))
	suite.Run(t, new(GenNumberEntityArrSuite))
	suite.Run(t, new(GenerateAcharyaAndStandoffSuite))
	suite.Run(t, new(ValidateFlagsSuite))
	suite.Run(t, new(HandleMainTestSuite))

}
