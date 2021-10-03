package main

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestGetSubString(t *testing.T) {
	type params struct {
		originalString string
		startPos       int
		endPos         int
	}

	tests := []struct {
		input    *params
		expected string
	}{
		{&params{"Hello world", 0, 2}, "He"},
		{&params{"Apple MacBook Pro (16-inch", 0, 2}, "Ap"},
		{&params{"नमस्ते दुनिया", 0, 2}, "नम"},
		{&params{"ਸਤਿ ਸ੍ਰੀ ਅਕਾਲ ਦੁਨਿਆ", 0, 2}, "ਸਤ"},
		{&params{"ഹലോ വേൾഡ്", 0, 2}, "ഹല"},
		{&params{"こんにちは世界", 0, 2}, "こん"},
	}
	for _, testCase := range tests {
		output := getSubString(testCase.input.originalString, testCase.input.startPos, testCase.input.endPos)
		if output != testCase.expected {
			log.Errorf("Test FAILED: input: %+v, expected: %s, received: %s", testCase.input, testCase.expected, output)
		}
	}
}

func TestGetEntities(t *testing.T) {

	confPaths := []string{
		"./testData/news/annotation.conf",
		"./testData/CoNLL-ST_2002/annotation.conf",
	}

	entities := [][]string{
		{"Person", "Organization", "GPE", "Money", "Inv1"},
		{"ORG", "PER", "LOC", "MISC", "Inv2"},
	}
	invalid := []string{"Inv1", "Inv2"}

	for i, v := range entities {
		cDat, cErr := os.Open(confPaths[i])
		if cErr != nil {
			log.Fatal(cErr)
		}
		defer cDat.Close()
		output := getEntities(cDat)
		for _, e := range v {
			_, ok := output[e]
			if !ok {
				if e != invalid[0] && e != invalid[1] {
					log.Errorf("Test FAILED: %s was not found in the annotation file at %s", e, confPaths[i])
				}
			}
		}
	}
}

func TestGeneratePaths(t *testing.T) {

	conll2002Ann := []string{"testData/CoNLL-ST_2002/esp/esp.train-doc-100.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-1400.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-423.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-46.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-503.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-896.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-919.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-978.ann", "testData/CoNLL-ST_2002/esp/esp.train-doc-989.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-118.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-123.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-134.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-181.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-184.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-236.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-251.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-27.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-46.ann", "testData/CoNLL-ST_2002/ned/ned.train-doc-75.ann"}
	conll2002Txt := []string{"testData/CoNLL-ST_2002/esp/esp.train-doc-100.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-1400.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-423.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-46.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-503.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-896.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-919.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-978.txt", "testData/CoNLL-ST_2002/esp/esp.train-doc-989.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-118.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-123.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-134.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-181.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-184.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-236.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-251.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-27.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-46.txt", "testData/CoNLL-ST_2002/ned/ned.train-doc-75.txt"}

	anns, txts, _ := generatePaths("./testData/CoNLL-ST_2002")

	for i, v := range conll2002Ann {

		if v != filepath.ToSlash(anns[i]) {
			log.Errorf("Test FAILED: Annotation file path %s is not equal to %s", strings.Replace(anns[i], "/", "\\", -1), v)
		}
		if conll2002Txt[i] != filepath.ToSlash(txts[i]) {
			log.Errorf("Test FAILED: Text file path %s is not equal to %s", strings.Replace(txts[i], "/", "\\", -1), v)
		}

	}

}

func TestGenerateEntityMap(t *testing.T) {

	expected := []AcharyaEntity{{418, 426, "Organization"},
		{456, 468, "Money"},
		{0, 0, ""},
		{473, 496, "Person"},
		{511, 535, "Person"},
		{540, 545, "Organization"},
		{549, 560, "GPE"},
	}

	cDat, cErr := os.Open("./testData/news/annotation.conf")
	if cErr != nil {
		log.Fatal(cErr)
	}
	defer cDat.Close()
	aDat, aErr := os.Open("./testData/news/000-introduction.ann")
	if aErr != nil {
		log.Fatal(aErr)
	}
	defer aDat.Close()
	ent := getEntities(cDat)
	entMap := generateEntityMap(ent, aDat)

	for i, v := range expected {
		if entMap[i] != v {
			log.Errorf("Test FAILED: entity at postion %d does not match (%+v != %+v) ", i, entMap[i], v)
		}
	}

}

func TestGenerateAcharyaAndConvert(t *testing.T) {

	expected := `{"Data":"Welcome to the Brat Rapid Annotation Tool (brat) tutorial!\n\nbrat is a web-based tool for structured text annotation and visualization. The easiest way to explain what this means is by example: see the following sentence illustrating various types of annotation. Take a moment to study this example, moving your mouse cursor over some of the annotations. Hold the cursor still over an annotation for more detail.\n\n\n1 ) Citibank was involved in moving about $100 million for Raul Salinas de Gortari, brother of a former Mexican president, to banks in Switzerland.\n\n\nIf this example seems complicated, don't panic! This tutorial will present the key features of brat interactively, with each document presenting one or a few features. If you follow this brief tutorial, you'll be able to understand and create annotations such as those above in no time.\n\nTry moving to the next document now by clicking on the arrow to the right on the blue bar at the top left corner of the page.\n","Entities":[[418,426,"Organization"],[456,468,"Money"],[473,496,"Person"],[511,535,"Person"],[540,545,"Organization"],[549,560,"GPE"]]}` + "\n"
	cDat, cErr := os.Open("./testData/news/annotation.conf")
	if cErr != nil {
		log.Fatal(cErr)
	}
	defer cDat.Close()
	aDat, aErr := os.Open("./testData/news/000-introduction.ann")
	if aErr != nil {
		log.Fatal(aErr)
	}
	defer aDat.Close()

	tDat, tErr := os.Open("./testData/news/000-introduction.txt")
	if tErr != nil {
		log.Fatal(tErr)
	}
	defer tDat.Close()

	tData := ""
	scannerD := bufio.NewScanner(tDat)
	scannerD.Split(bufio.ScanLines)

	for scannerD.Scan() {
		tData = tData + scannerD.Text() + "\n"
	}

	ent := getEntities(cDat)
	entMap := generateEntityMap(ent, aDat)
	acharya, _ := generateAcharyaAndConvert(string(tData), ent, entMap)

	if acharya != expected {
		log.Errorf("Test FAILED: generated acharya format does not match with expected.\n Generated:%s \n\n\n Expected:%s ", acharya, expected)
	}

}

func TestHandleOutput(t *testing.T) {
	e := os.Remove("./testData/file_gen_by_handleOutputTest.jsonl")
	if e != nil {
		log.Fatal(e)
	}
	acharya := `{"Data":"Welcome to the Brat Rapid Annotation Tool (brat) tutorial!\n\nbrat is a web-based tool for structured text annotation and visualization. The easiest way to explain what this means is by example: see the following sentence illustrating various types of annotation. Take a moment to study this example, moving your mouse cursor over some of the annotations. Hold the cursor still over an annotation for more detail.\n\n\n1 ) Citibank was involved in moving about $100 million for Raul Salinas de Gortari, brother of a former Mexican president, to banks in Switzerland.\n\n\nIf this example seems complicated, don't panic! This tutorial will present the key features of brat interactively, with each document presenting one or a few features. If you follow this brief tutorial, you'll be able to understand and create annotations such as those above in no time.\n\nTry moving to the next document now by clicking on the arrow to the right on the blue bar at the top left corner of the page.\n","Entities":[[418,426,"Organization"],[456,468,"Money"],[473,496,"Person"],[511,535,"Person"],[540,545,"Organization"],[549,560,"GPE"]]}` + "\n"
	handleOutput("./testData/file_gen_by_handleOutputTest.jsonl", acharya, "not loaded txt", false, 0)
	achDat, achErr := os.Open("./testData/file_gen_by_handleOutputTest.jsonl")
	if achErr != nil {
		log.Fatal(achErr)
	}

	achData := ""
	scannerD := bufio.NewScanner(achDat)
	scannerD.Split(bufio.ScanLines)

	for scannerD.Scan() {
		achData = achData + scannerD.Text() + "\n"
	}

	if achData != acharya {
		log.Errorf("Test FAILED: content in the file ./testData/test.jsonl do not match with Expected:%s Got:%s", acharya, achData)
	}

}

func TestHandleMain(t *testing.T) {

	handleMain("./testData/news", "", "./testData/file_gen_by_handleMainTest.jsonl", "", "", true)

	genDat, genErr := os.Open("./testData/file_gen_by_handleMainTest.jsonl")
	if genErr != nil {
		log.Fatal(genErr)
	}

	genData := ""
	scanner := bufio.NewScanner(genDat)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		genData = genData + scanner.Text() + "\n"
	}

	cmpDat, cmpErr := os.Open("./testData/test_cmp.jsonl")
	if cmpErr != nil {
		log.Fatal(genErr)
	}

	cmpData := ""
	scanner = bufio.NewScanner(cmpDat)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		cmpData = cmpData + scanner.Text() + "\n"
	}

	if genData != cmpData {
		log.Errorf("Test FAILED: The generated file and the reference file don't match \n\n Expected: %s \n\n Got: %s", cmpData, genData)
	}

}
