@echo off
set "fileName=main.exe"
set "fileURL=https://drive.google.com/file/d/1uEbRUeZpXXL2tIMSncsyMgEHorRUmMtF/view"

@REM net stop "H2J Monitor"
@REM ping 127.0.0.1 -n 3 > nul
cd /d "%~dp0"
curl -L -o "https://drive.google.com/file/d/1uEbRUeZpXXL2tIMSncsyMgEHorRUmMtF/view?usp=sharing" "main.exe"
@REM net start "H2J Monitor"