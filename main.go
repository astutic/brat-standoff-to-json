package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	LOG_NO_ENTITIES                              = "The conf file does not have an `[entities]` field or `[entities]` field is empty"
	LOG_MULTIPLE_CONF_FILES_FOUND                = "Multiple `annotation.conf` files found"
	LOG_FILES_NOT_EXIST                          = "%s File does not exist"
	LOG_DISCONTINUOS_TEXTBOUND_ANN_NOT_SUPPORTED = "Discontinuous text-bound annotations is not currently supported"
	LOG_SUCCESSFULLY_GENERATED_TEXTFILE          = "Successfully generated `%s` for %s\n"
	LOG_SUCCESSFULLY_UPDATED_TEXTFILE            = "Successfully generated `%s` for %s\n"
	LOG_NUMBER_OF_ANN_FILES_NOT_EQUAL_TO_TXT     = "The number of annotation files should be equal to the number of text files"
)

// Flag constants
const (
	logValNotSet              = "kindly set the value for -%s or --%s\n"
	logOutputFileNotSpecified = "overwrite flag is provided but output file is not specified"
)

func exit1() {
	os.Exit(1)
}

type AcharyaEntity struct {
	begin int
	end   int
	name  string
}

func getSubString(originalString string, startPos, endPos int) string {
	counter := 0
	val := ""
	var r rune
	for i, s := 0, 0; i <= len(originalString); i += s {
		r, s = utf8.DecodeRuneInString(originalString[i:])
		if r == '\r' {
			continue
		}
		if counter >= startPos {
			if counter >= endPos {
				break
			}
			val = val + string(r)
		}
		counter++
	}
	return val
}

func getEntities(confFile *os.File) map[string]bool {

	scannerC := bufio.NewScanner(confFile)
	scannerC.Split(bufio.ScanLines)
	startScan := false
	entities := make(map[string]bool)

	for scannerC.Scan() {
		if strings.Contains(scannerC.Text(), "[entities]") {
			startScan = true
			continue
		}
		if startScan {
			if len(strings.TrimSpace(scannerC.Text())) == 0 {
				continue
			}
			if strings.HasPrefix(scannerC.Text(), "[") {
				break
			} else if strings.HasPrefix(scannerC.Text(), "#") {
				continue
			}
			entities[strings.TrimSpace(scannerC.Text())] = true
		}
	}
	if len(entities) == 0 {
		fmt.Println(LOG_NO_ENTITIES)
		exit1()
	}
	return entities
}

func generatePaths(path string) ([]string, []string, error) {
	const dotAnn = ".ann"
	const dotTxt = ".txt"

	confCount := 0
	annMult := []string{}
	textMult := []string{}

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if confCount > 1 {
				fmt.Println(LOG_MULTIPLE_CONF_FILES_FOUND)
				exit1()
			}

			switch {
			case strings.HasSuffix(path, dotAnn):
				if _, err := os.Stat(strings.TrimSuffix(path, dotAnn) + dotTxt); os.IsNotExist(err) {
					fmt.Printf(LOG_FILES_NOT_EXIST, strings.TrimSuffix(path, dotAnn)+dotTxt)
					exit1()
				}
				annMult = append(annMult, path)
				textMult = append(textMult, strings.TrimSuffix(path, dotAnn)+dotTxt)
			case strings.HasSuffix(path, dotTxt):
				if _, err := os.Stat(strings.TrimSuffix(path, dotTxt) + dotAnn); os.IsNotExist(err) {
					fmt.Printf(LOG_FILES_NOT_EXIST, strings.TrimSuffix(path, dotTxt)+dotAnn)
					exit1()
				}
			case strings.HasSuffix(path, "annotation.conf"):
				confCount++
			}
			return nil
		})
	if err != nil {
		return []string{}, []string{}, err
	}
	return annMult, textMult, nil
}

func generateEntityMap(ent map[string]bool, aData *os.File) map[int]AcharyaEntity {
	scanner := bufio.NewScanner(aData)
	scanner.Split(bufio.ScanLines)

	entityMap := make(map[int]AcharyaEntity)
	count := 0

	for scanner.Scan() {
		// Uncomment the lines below to dispaly the ann file
		// fmt.Println(strings.Repeat("#", 30), "Annotations", strings.Repeat("#", 30))
		// fmt.Println(scanner.Text())
		if strings.HasPrefix(scanner.Text(), "T") {
			splitAnn := strings.Split(strings.Replace(scanner.Text(), "\t", ",", 2), ",")
			if strings.Contains(splitAnn[1], ";") {
				fmt.Println(LOG_DISCONTINUOS_TEXTBOUND_ANN_NOT_SUPPORTED)
				exit1()
			}
			entAndPos := strings.Split(splitAnn[1], " ")
			if ent[strings.TrimSpace(entAndPos[0])] {
				b, _ := strconv.Atoi(entAndPos[1])
				e, _ := strconv.Atoi(entAndPos[2])
				entityMap[count] = AcharyaEntity{b, e, entAndPos[0]}
			}
			count++
		}
	}

	return entityMap
}

func generateAcharyaAndConvert(tData string, ent map[string]bool, entMap map[int]AcharyaEntity) (string, string) {
	conv := ""
	// It is necessary to marshal string as to avoid problems by escape sequences
	escapedStr, err := json.Marshal(tData)
	if err != nil {
		fmt.Println(err)
		exit1()
	}

	acharya := fmt.Sprintf("{\"Data\":%s,\"Entities\":[", fmt.Sprintf("%s", escapedStr))
	j := 0
	c := 0
	for {
		ent, ok := entMap[j]
		if ok {
			conv = conv + fmt.Sprintf("T%d\t%s %d %d\t%s\n", j+1, ent.name, ent.begin, ent.end, getSubString(tData, ent.begin, ent.end))
			acharya = acharya + fmt.Sprintf("[%d,%d,\"%s\"],", ent.begin, ent.end, ent.name)
			c++
		}
		// loop until all entities in a are fetched
		if c >= len(entMap) {
			conv = strings.TrimSuffix(conv, "\n")
			acharya = strings.TrimSuffix(acharya, ",")
			acharya = strings.ReplaceAll(acharya, "\n", "\\n")
			acharya = acharya + "]}\n"
			break
		}
		j++
	}
	return acharya, conv
}

func handleOutput(outputFile, acharya, textfileName string, overWrite bool, currentParseIndex int) {
	if outputFile == "" {
		fmt.Println(acharya)
	} else {

		// asking for the first time
		if currentParseIndex == 0 && !overWrite {
			if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
				fmt.Println("The output file already exists use `--overwrite flag to overwrite the file`")
				exit1()
			}
		}

		if currentParseIndex == 0 && overWrite {
			if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
				e := os.Remove(outputFile)
				if e != nil {
					fmt.Println(e)
					exit1()
				}
			}

		}

		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			fmt.Println(err)
			exit1()
		}
		if _, err = f.WriteString(acharya); err != nil {
			fmt.Println(err)
			exit1()
		}
		if currentParseIndex == 0 {
			fmt.Printf(LOG_SUCCESSFULLY_GENERATED_TEXTFILE, outputFile, strings.TrimSpace(textfileName))
		} else {
			fmt.Printf(LOG_SUCCESSFULLY_UPDATED_TEXTFILE, outputFile, strings.TrimSpace(textfileName))
		}
	}
}

func handleMain(all, conf, opfile, anns, texts string, overwrite bool) {
	annMult := []string{}
	textMult := []string{}
	var err error
	if all != "" {
		annMult, textMult, err = generatePaths(all)

		if err != nil {
			fmt.Println(err)
			exit1()
		}
	}

	var confPath string

	if all != "" {
		confPath = all + "/annotation.conf"
	} else {
		confPath = conf
	}

	cDat, cErr := os.Open(confPath)
	if cErr != nil {
		fmt.Println(cErr)
		exit1()
	}
	defer cDat.Close()

	ent := getEntities(cDat)

	if all == "" {
		annMult = strings.Split(anns, ",")
		textMult = strings.Split(texts, ",")
	}

	if len(annMult) != len(textMult) {
		fmt.Println(LOG_NUMBER_OF_ANN_FILES_NOT_EQUAL_TO_TXT)
		exit1()
	}

	for i := range annMult {
		aData, aErr := os.Open(strings.TrimSpace(annMult[i]))
		if aErr != nil {
			fmt.Println(aErr.Error())
			exit1()
		}
		defer aData.Close()

		tDat, tErr := os.Open(strings.TrimSpace(textMult[i]))
		if tErr != nil {
			fmt.Println(tErr.Error())
			exit1()
		}
		defer tDat.Close()

		tData := ""
		scannerD := bufio.NewScanner(tDat)
		scannerD.Split(bufio.ScanLines)

		for scannerD.Scan() {
			// add `\n` since they are removed by `scannerD.Split(bufio.ScanLines)`
			tData = tData + scannerD.Text() + "\n"
		}

		entityMap := generateEntityMap(ent, aData)

		acharya, _ := generateAcharyaAndConvert(string(tData), ent, entityMap)

		handleOutput(opfile, acharya, textMult[i], overwrite, i)

		// Uncomment the below lines to show the converted files
		// fmt.Println(strings.Repeat("#", 30), "Convert", strings.Repeat("#", 30))
		// fmt.Println(conv)

		// uncomment the  below lines to write an output.ann file
		// err := ioutil.WriteFile("./output.ann", []byte(conv), 0644)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// fmt.Println("Successfully `output.ann` in the current diretory")
	}
}

func main() {
	all := flag.String("all", "", "Path to the folder containing the collection")
	anns := flag.String("ann", "", "Comma sepeartad locations of the annotation files (.ann) in correct order")
	texts := flag.String("text", "", "Comma sepeartad locations of the text files (.txt) in correct order")
	conf := flag.String("conf", "", "Location of the annotation configuration file (annotation.conf)")
	oFileName := flag.String("output", "", "Name of the output file to be generated")
	overWrite := flag.Bool("overwrite", false, "If you wish to overwrite the generated file then set overwrite to true")

	flag.Parse()

	if *all == "" {
		flag.VisitAll(func(f *flag.Flag) {
			if f.Value.String() == "" {
				if (f.Name != "output") && (f.Name != "all") {
					fmt.Printf(logValNotSet, f.Name, f.Name)
					exit1()
				}
			}
		})
	}

	if *overWrite && *oFileName == "" {
		fmt.Println(logOutputFileNotSpecified)
		exit1()
	}

	handleMain(*all, *conf, *oFileName, *anns, *texts, *overWrite)

}
