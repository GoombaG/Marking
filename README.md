# Windows Setup

## wkhtmltopdf
wkhtmltopdf needs to be downloaded from https://wkhtmltopdf.org/downloads.html -- an open source project hosted on github at https://github.com/wkhtmltopdf. It should be installed in the default location at C:\Program Files\

This is needed to run the program.

## Go
Install go following the instructions listed here: https://go.dev/doc/install

This is needed to build the program.

## Obtaining the program
After installing Go and wkhtmltopdf, a zip file of the Marking Program should be downloaded (green Code button at the top right of the page).

Unzip the folder, navigate to “main” in the command prompt at "Marking-main\Marking-main\main”, and then type “go build”. This will generate an executable called “main.exe” which can run the program.

This can be moved anywhere and renamed.

## Running the program
Place your data spreadsheet in the same location as the executable. After double clicking the executable to launch, go to http://localhost:8080 in Chrome.
