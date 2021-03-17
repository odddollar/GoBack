# GoBack File Backup

GoBack is a simple file backup tool, similar to Microsoft's SyncToy. 

By giving GoBack a source and destination directory, it checks if each file in the source directory is in the destination directory. If it [the current file] is in the destination directory then nothing happens, if not then it is copied across. 

This means the program only copies across the differences, creating fast incremental backups where no files are ever deleted or overwritten in the destination directory.

## Usage

Run the following in a command prompt:

```
goback -s [PATH_TO_SOURCE_DIRECTORY] -d [PATH_TO_DESTINATION_DIRECTORY]
```

Provide the ```-p``` flag to print a table of actions taken for each file upon completion.
Provide ```-o``` to output this data to ```output.csv```.

## How it works

The program walks the source directory, gets the modification date of the current file and appends it to the file name in memory. It then checks if this file name exists in the destination directory. If not then it copies it across with the modification-date-appended filename. 

This means that when a file is modified, GoBack will recognise it as a completely new file when run and copy it across, thus not deleting or overwriting old versions of that file. 

This allows for going back (Haha, get it?) to previous versions and safely and incrementally backing up files.

## Build from source

To build from source, ensure that Go is installed.

Run: 

```
git clone https://github.com/odddollar/GoBack.git
cd GoBack
go build GoBack
```
