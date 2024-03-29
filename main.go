package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	flag "github.com/spf13/pflag"
)

var Version = "development"

func printVersion() {
	fmt.Printf("\n bratconverter version: %s\n", Version)
}

const (
	ErrNoEntities                           = "the conf file does not have an `[entities]` field or `[entities]` field is empty"
	ErrMultipleConfFilesFound               = "multiple `annotation.conf` files found"
	ErrDiscontinuosTextboundAnnNotSupported = "discontinuous text-bound annotations is not currently supported"

	ErrSubStrNegativeStartPos         = "start position should be a positive number, Received start position %d"
	ErrSubStrEndPosSmallerThanStart   = "end position should be greater than start position, Received end position %d"
	ErrSubStrEndposGreaterThanDataLen = "end position should be lesser than length of the txt data, Length of txt data: %d, End position: %d"

	ErrFilesNotExist = "%s file does not exist"

	ErrTxtAnnBadFormat = "text annotation is badly formatted"

	ErrBadFormat    = "file follows unknown format: "
	ErrBadFormatTab = "file follows unknown format: expected 3 properties separared by [tab]"

	ErrFlagFileAlreadyExists = "the output file already exists use `--force` or `-f` flag  to overwrite the file"

	InfoSuccessfullyGenFile = "successfully generated file: %s"

	ErrValidateNoAnnFiles         = "no annotation files specified in the input"
	ErrValidateNoTxtFiles         = "no txt files specified in the input"
	ErrValidateNoConfFile         = "no conf file specified in the input"
	ErrValidateEmptyFolder        = "received empty folder path"
	ErrValidateOutputFileNotFound = "force flag is provided but output file is not specified"

	ErrNoAnnNoTxtNotMatch        = "the number of annotation files should be equal to the number of txt files,\n Received Annotation Files: %s Length: %d,Txt Files: %s Length: %d"
	ErrAnnFileNotCorrespondToTxt = "expected annotation file: %s to correspond to: %s.txt Received: %s"
)

func exit1() {
	os.Exit(1)
}

type AcharyaEntity struct {
	Begin int
	End   int
	Name  string
}

type NumberAcharyaEntity struct {
	TxtAnnNo int
	Entity   AcharyaEntity
}

func GetSubString(originalString string, startPos, endPos int) (string, error) {

	if startPos < 0 {
		return "", fmt.Errorf(ErrSubStrNegativeStartPos, startPos)
	} else if endPos < startPos {
		return "", fmt.Errorf(ErrSubStrEndPosSmallerThanStart, endPos)
	} else if endPos > len(originalString) {
		return "", fmt.Errorf(ErrSubStrEndposGreaterThanDataLen, len(originalString), endPos)
	}

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
	return val, nil
}

func GetEntitiesFromFile(confFile *os.File) map[string]bool {

	scanner := bufio.NewScanner(confFile)
	scanner.Split(bufio.ScanLines)
	startScan := false
	entities := make(map[string]bool)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "[entities]") {
			startScan = true
			continue
		}
		if startScan {
			if len(strings.TrimSpace(scanner.Text())) == 0 {
				continue
			}
			if strings.HasPrefix(scanner.Text(), "[") {
				break
			} else if strings.HasPrefix(scanner.Text(), "#") {
				continue
			}
			entities[strings.TrimSpace(scanner.Text())] = true
		}
	}
	return entities
}

func GetSubDirectories(path string) ([]string, []string, error) {
	const dotAnnSuffix = ".ann"
	const dotTxtSuffix = ".txt"

	annConfCount := 0
	annMult := []string{}
	textMult := []string{}

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			switch {
			// .ann file should have a corresponding .txt file
			case strings.HasSuffix(path, dotAnnSuffix):
				if _, err := os.Stat(strings.TrimSuffix(path, dotAnnSuffix) + dotTxtSuffix); os.IsNotExist(err) {
					return fmt.Errorf(ErrFilesNotExist, strings.TrimSuffix(path, dotAnnSuffix)+dotTxtSuffix)
				}
				annMult = append(annMult, path)
				textMult = append(textMult, strings.TrimSuffix(path, dotAnnSuffix)+dotTxtSuffix)
			// .ann file should have a corresponding .txt file
			case strings.HasSuffix(path, dotTxtSuffix):
				if _, err := os.Stat(strings.TrimSuffix(path, dotTxtSuffix) + dotAnnSuffix); os.IsNotExist(err) {
					return fmt.Errorf(ErrFilesNotExist, strings.TrimSuffix(path, dotTxtSuffix)+dotAnnSuffix)
				}
			case strings.HasSuffix(path, "annotation.conf"):
				annConfCount++
			}

			if annConfCount > 1 {
				return errors.New(ErrMultipleConfFilesFound)
			}
			return nil
		})
	if err != nil {
		return []string{}, []string{}, err
	}
	return annMult, textMult, nil
}

func GetTextAnnNum(ann string) (int, error) {
	if len(ann) > 0 {
		annSplit := strings.Split(ann, "\t")
		if len(annSplit[0]) > 1 {
			noStr := annSplit[0][1:]
			return strconv.Atoi(noStr)
		}
	}
	return 0, errors.New(ErrTxtAnnBadFormat)
}

func GenNumberEntityArr(entFromConf map[string]bool, aData *os.File) ([]NumberAcharyaEntity, error) {
	scanner := bufio.NewScanner(aData)
	scanner.Split(bufio.ScanLines)

	numberEntityArr := []NumberAcharyaEntity{}

	for scanner.Scan() {
		// Uncomment the lines below to dispaly the ann file
		// fmt.Println(strings.Repeat("#", 30), "Annotations", strings.Repeat("#", 30))
		// fmt.Println(scanner.Text())
		if strings.HasPrefix(scanner.Text(), "T") {
			splitAnn := strings.Split(scanner.Text(), "\t")
			if len(splitAnn) == 3 {
				if strings.Contains(splitAnn[1], ";") {
					return []NumberAcharyaEntity{}, errors.New(ErrDiscontinuosTextboundAnnNotSupported)
				}
				entAndPos := strings.Split(splitAnn[1], " ")
				if (len(entAndPos)) == 3 {
					if entFromConf[strings.TrimSpace(entAndPos[0])] {
						b, err := strconv.Atoi(entAndPos[1])
						if err != nil {
							return []NumberAcharyaEntity{}, err
						}
						e, err := strconv.Atoi(entAndPos[2])
						if err != nil {
							return []NumberAcharyaEntity{}, err
						}

						annotationNo, err := GetTextAnnNum(scanner.Text())
						if err != nil {
							return []NumberAcharyaEntity{}, err
						}

						numberEntityArr = append(numberEntityArr, NumberAcharyaEntity{annotationNo, AcharyaEntity{b, e, entAndPos[0]}})
					}
				} else {
					return numberEntityArr, errors.New(ErrBadFormat)
				}
			} else {
				return numberEntityArr, errors.New(ErrBadFormatTab)
			}
		}
	}

	return numberEntityArr, nil
}

func GenerateAcharyaAndStandoff(tData string, numberAcharyaEnt []NumberAcharyaEntity) (string, string, error) {
	standoff := ""
	// It is necessary to marshal string as to avoid problems by escape sequences
	escapedStr, err := json.Marshal(tData)
	if err != nil {
		return "", "", err
	}

	acharya := fmt.Sprintf("{\"Data\":%s,\"Entities\":[", fmt.Sprintf("%s", escapedStr))

	for _, v := range numberAcharyaEnt {
		str, err := GetSubString(tData, v.Entity.Begin, v.Entity.End)
		if err != nil {
			return "", "", err
		}
		standoff = standoff + fmt.Sprintf("T%d\t%s %d %d\t%s\n", v.TxtAnnNo, v.Entity.Name, v.Entity.Begin, v.Entity.End, str)
		acharya = acharya + fmt.Sprintf("[%d,%d,\"%s\"],", v.Entity.Begin, v.Entity.End, v.Entity.Name)
	}

	standoff = strings.TrimSuffix(standoff, "\n")
	acharya = strings.TrimSuffix(acharya, ",")
	acharya = strings.ReplaceAll(acharya, "\n", "\\n")
	acharya = acharya + "]}\n"

	return acharya, standoff, nil
}

func handleOutput(outputFile, acharya string, overWrite bool) error {
	if !overWrite {
		if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
			return errors.New(ErrFlagFileAlreadyExists)
		}
	}

	f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.WriteString(acharya); err != nil {
		return err
	}

	return nil
}

func handleMain(fPath, annFiles, txtFiles, conf, opFile string, overwrite bool) error {
	annMult := []string{}
	textMult := []string{}
	var err error
	if fPath != "" {
		annMult, textMult, err = GetSubDirectories(fPath)
		if err != nil {
			return err
		}
	}

	var confPath string

	if fPath != "" {
		// If a folder path is provided then the annotation conf file should be present in the root of the folder
		confPath = fPath + "/annotation.conf"
	} else {
		confPath = conf
	}

	confFile, cErr := os.Open(confPath)
	if cErr != nil {
		return cErr
	}
	defer confFile.Close()

	entities := GetEntitiesFromFile(confFile)
	if len(entities) == 0 {
		return errors.New(ErrNoEntities)
	}

	if fPath == "" {
		annMult = strings.Split(annFiles, ",")
		textMult = strings.Split(txtFiles, ",")
	}

	generatedAcharya := ""

	for i := range annMult {
		annFile, aErr := os.Open(strings.TrimSpace(annMult[i]))
		if aErr != nil {
			return aErr
		}
		defer annFile.Close()

		txtFile, tErr := os.Open(strings.TrimSpace(textMult[i]))
		if tErr != nil {
			return tErr
		}
		defer txtFile.Close()

		txtFileData, err := ioutil.ReadAll(txtFile)
		if err != nil {
			return err
		}

		entityArr, err := GenNumberEntityArr(entities, annFile)
		if err != nil {
			return err
		}

		acharya, _, err := GenerateAcharyaAndStandoff(string(txtFileData), entityArr)
		if err != nil {
			return err
		}

		generatedAcharya = generatedAcharya + acharya

	}

	if opFile == "" {
		fmt.Println(generatedAcharya)
		return nil
	}

	err = handleOutput(opFile, generatedAcharya, overwrite)
	if err != nil {
		return err
	}

	fmt.Printf(InfoSuccessfullyGenFile, opFile)

	return nil
}

func ValidateFlags(fPath, annFiles, txtFiles, confFile, oFileName string, overWrite bool) error {
	if len(fPath) == 0 {
		switch {
		case IsEmptyString(annFiles):
			return errors.New(ErrValidateNoAnnFiles)
		case IsEmptyString(txtFiles):
			return errors.New(ErrValidateNoTxtFiles)
		case IsEmptyString(confFile):
			return errors.New(ErrValidateNoConfFile)
		}

		err := ValidateAnnAndTxt(annFiles, txtFiles)
		if err != nil {
			return err
		}
	} else if IsEmptyString(fPath) {
		return errors.New(ErrValidateEmptyFolder)
	}

	if overWrite && oFileName == "" {
		return errors.New(ErrValidateOutputFileNotFound)
	}

	return nil
}

func ValidateAnnAndTxt(ann, txt string) error {
	annArray := strings.Split(ann, ",")
	txtArray := strings.Split(txt, ",")

	if len(annArray) != len(txtArray) {
		return fmt.Errorf(ErrNoAnnNoTxtNotMatch, annArray, len(annArray), txtArray, len(txtArray))
	}

	for i, annPath := range annArray {
		annBaseName := strings.TrimSpace(filepath.Base(annPath))
		txtBaseName := strings.TrimSpace(filepath.Base(txtArray[i]))
		if strings.TrimSuffix(annBaseName, filepath.Ext(annBaseName))+".txt" != txtBaseName {
			return fmt.Errorf(ErrAnnFileNotCorrespondToTxt, annPath, strings.TrimSuffix(annBaseName, filepath.Ext(annBaseName)), txtArray[i])
		}
	}

	return nil
}

func IsEmptyString(s string) bool {
	return strings.TrimSpace(s) == "" || len(s) <= 0
}

func main() {
	folderPath := flag.StringP("folderPath", "p", "", "Path to the folder containing the collection")
	annFiles := flag.StringP("ann", "a", "", "Comma sepeartad locations of the annotation files (.ann) in correct order")
	txtFiles := flag.StringP("txt", "t", "", "Comma sepeartad locations of the text files (.txt) in correct order")
	confFile := flag.StringP("conf", "c", "", "Location of the annotation configuration file (annotation.conf)")
	oFileName := flag.StringP("output", "o", "", "Name of the output file to be generated")
	overWrite := flag.BoolP("force", "f", false, "If you wish to overwrite the generated file then set force to true")
	version := flag.BoolP("version", "v", false, "Print bratconverter version")

	flag.Parse()

	if *version {
		printVersion()
		exit1()
	}

	err := ValidateFlags(*folderPath, *annFiles, *txtFiles, *confFile, *oFileName, *overWrite)
	if err != nil {
		fmt.Println(err)
		exit1()
	}

	err = handleMain(*folderPath, *annFiles, *txtFiles, *confFile, *oFileName, *overWrite)
	if err != nil {
		fmt.Println(err)
		exit1()
	}
}
