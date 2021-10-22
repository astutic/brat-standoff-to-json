**[ Windows PowerShell ]** 
If you want to use the Brat → JSONL converter and If Brat Standoff **contains non English characters** Then its advised to set the following in PowerShell first 👇👇
```powershell
$OutputEncoding = [console]::InputEncoding = [console]::OutputEncoding = New-Object System.Text.UTF8Encoding
```
#### Features that are currently unsupported:
-   [Discontinuous text-bound annotations](https://brat.nlplab.org/standoff.html "https://brat.nlplab.org/standoff.html")
-   [Advanced entity configuration](https://brat.nlplab.org/configuration.html#tool-configuration "https://brat.nlplab.org/configuration.html#tool-configuration")

## Using brat Standoff Converter

```bash
 git clone https://github.com/astutic/bratStandoffConverter.git
 ```
OR 

Download a release from [here](https://github.com/astutic/bratStandoffConverter/releases)

Then run the file using go OR use the executable

#### Generates acharya format for files in a specific directory and logs it to the console

```bash
go run main.go --all "./path/to/the/collection"
```
OR
```bash
bratconverter --all "./path/to/the/collection"
```

example  
```bash
go run main.go --all "./testData/news"
```
OR
```bash
bratconverter  --all "./testData/news"
```


#### Generate an output file

```bash
go run main.go --all "./path/to/the/collection" --output "path/output-file-name"
```  
OR
```bash
bratconverter  --all "./path/to/the/collection" --output "path/output-file-name"
```  

example
> The command below will generate an output file named **acharyaFormat.jsonl** in the current directory

```bash
go run main.go --all "./testData/news" --output "./acharyaFormat.jsonl"
```
OR
```bash
bratconverter  --all "./testData/news" --output "./acharyaFormat.jsonl"
```

### Generating for specific files

! **NOTE** the order of the .ann files an .txt files should be the same  
`go run main.go --ann "file1.ann,file2.ann" --text "file1.txt,file2.txt" --conf "file.conf"`

example  
```bash
go run main.go --ann "testData/news/000-introduction.ann,testData/news/040-text_span_annotation.ann" --text "testData/news/000-introduction.txt,testData/news/040-text_span_annotation.txt" --conf "testData/news/annotation.conf"
```
OR
```bash
bratconverter  --ann "testData/news/000-introduction.ann,testData/news/040-text_span_annotation.ann" --text "testData/news/000-introduction.txt,testData/news/040-text_span_annotation.txt" --conf "testData/news/annotation.conf"
```

|Command   | Type  |     Description| Default value  |   
|---|---|---|---|
|-overwrite   | bool  | If you wish to overwrite the generated file then set overwrite to true  |  false |
|-all| string |Path to the folder containing the collection|   |
|-ann| string| Comma sepeartad locations of the annotation files (.ann) in correct order|   |
|-conf| string |Location of the annotation configuration file (annotation.conf)|   |
|-output| string| Name of the output file to be generated|   |
|-text| string| Comma sepeartad locations of the text files (.txt) in correct order|   |
