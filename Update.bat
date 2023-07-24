@echo off
net stop "H2J Monitor"
ping 127.0.0.1 -n 3 > nul
cd /d "%~dp0"
curl -LJO https://github.com/HenriqueJST/H2J_Monitor/raw/main/src/main.exe
net start "H2J Monitor"