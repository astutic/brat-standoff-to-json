**[ Windows PowerShell ]**
If you want to use the Brat → JSONL converter and If Brat Standoff **contains non English characters** Then its advised to set the following in PowerShell first

```powershell
$OutputEncoding = [console]::InputEncoding = [console]::OutputEncoding = New-Object System.Text.UTF8Encoding
```

#### Features that are currently unsupported:

- [Discontinuous text-bound annotations](https://brat.nlplab.org/standoff.html "https://brat.nlplab.org/standoff.html")
- [Advanced entity configuration](https://brat.nlplab.org/configuration.html#tool-configuration "https://brat.nlplab.org/configuration.html#tool-configuration")

## Using brat Standoff Converter

```bash
 git clone https://github.com/astutic/bratStandoffConverter.git
```

OR

Download a release from [here](https://github.com/astutic/bratStandoffConverter/releases)

Then run the file using go OR use the executable

## Examples

### Generates acharya format for files in a specific directory and logs it to the console

```bash
go run main.go -p "./path/to/the/collection"
```

OR

```bash
bratconverter -p "./path/to/the/collection"
```

##### example

```bash
go run main.go -p "./testData/news"
```

OR

```bash
bratconverter  -p "./testData/news"
```

### Generate an output file

```bash
go run main.go -p "./path/to/the/collection" --output "path/output-file-name"
```

OR

```bash
bratconverter  -p "./path/to/the/collection" --output "path/output-file-name"
```

##### example

> The command below will generate an output file named **acharyaFormat.jsonl** in the current directory

```bash
go run main.go -p "./testData/news" --output "./acharyaFormat.jsonl"
```

OR

```bash
bratconverter  -p "./testData/news" --output "./acharyaFormat.jsonl"
```

### Generating for specific files

! **NOTE** the order of the .ann files an .txt files should be the same  
`go run main.go --ann "file1.ann,file2.ann" --text "file1.txt,file2.txt" --conf "file.conf"`

##### example

```bash
go run main.go --ann "testData/news/000-introduction.ann,testData/news/040-text_span_annotation.ann" --text "testData/news/000-introduction.txt,testData/news/040-text_span_annotation.txt" --conf "testData/news/annotation.conf"
```

OR

```bash
bratconverter  --ann "testData/news/000-introduction.ann,testData/news/040-text_span_annotation.ann" --text "testData/news/000-introduction.txt,testData/news/040-text_span_annotation.txt" --conf "testData/news/annotation.conf"
```

## Commands

| Command    | Short hand | Type   | Description                                                               | Default value |
| ---------- | ---------- | ------ | ------------------------------------------------------------------------- | ------------- |
| folderPath | p          | string | Path to the folder containing the collection                              |
| ann        | a          | string | Comma sepeartad locations of the annotation files (.ann) in correct order |
| txt        | t          | string | Comma sepeartad locations of the text files (.txt) in correct order       |
| conf       | c          | string | Location of the annotation configuration file (annotation.conf)           |
| output     | o          | string | Name of the output file to be generated                                   |
| force      | f          | bool   | If you wish to overwrite the generated file then set force to true        | false         |

## Original data displayed in brat

![Original data displayed in brat](./docs/images/brat_ui.png "Brat UI")

## Data from Brat converted to Acharya format

![Brat data displayed in Acharya](./docs/images/brat_to_Acharya_ui.png "Acharya UI")
