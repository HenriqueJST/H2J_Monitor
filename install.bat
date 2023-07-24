@echo off

>nul 2>&1 "%SYSTEMROOT%\system32\cacls.exe" "%SYSTEMROOT%\system32\config\system"
if '%errorlevel%' NEQ '0' (powershell start -verb runas '%0' & exit /b)

cd /d "%~dp0"

setlocal enabledelayedexpansion

set "DEFAULT_QLIK_LOGS_PATH_SAAS=C:/ProgramData/Qlik/DataTransfer/Log"
set "DEFAULT_QLIK_LOGS_PATH_ONPRIMESE=C:/ProgramData/Qlik/Sense/Log"
set "DEFAULT_ENGINEPORT=4747"
set "DEFAULT_CERTIFICATESPATH=C:/ProgramData/Qlik/Sense/Repository/Exported Certificates/.Local Certificates"

set /p ACTION="Deseja Instalar (I), Desinstalar (D) ou dar Update (U)? (I/D/U): "
if /i "%ACTION%"=="I" (
    curl -LJO https://github.com/HenriqueJST/H2J_Monitor/raw/main/src/main.exe
    cls
    main.exe -mode=install
    cls
    set /p "SAAS=SAAS (true or false): "
    cls
    set /p "CLIENTE_NAME=CLIENTE_NAME: "
    cls
    echo.
    echo Press Enter to set Default value: %DEFAULT_QLIK_LOGS_PATH_SAAS%
    set /p "QLIK_LOGS_PATH_SAAS=QLIK_LOGS_PATH_SAAS: "
    if "!QLIK_LOGS_PATH_SAAS!"=="" set "QLIK_LOGS_PATH_SAAS=%DEFAULT_QLIK_LOGS_PATH_SAAS%"
    cls
    echo.
    echo Press Enter to set Default value: %DEFAULT_QLIK_LOGS_PATH_ONPRIMESE%
    set /p "QLIK_LOGS_PATH_ONPRIMESE=QLIK_LOGS_PATH_ONPRIMESE: "
    if "!QLIK_LOGS_PATH_ONPRIMESE!"=="" set "QLIK_LOGS_PATH_ONPRIMESE=%DEFAULT_QLIK_LOGS_PATH_ONPRIMESE%"
    cls
    echo.
    echo Press Enter to set Default value: %DEFAULT_ENGINEPORT%
    set /p "ENGINEPORT=ENGINEPORT: "
    if "!ENGINEPORT!"=="" set "ENGINEPORT=%DEFAULT_ENGINEPORT%"
    cls
    set /p "ENGINEHOST=ENGINEHOST: "
    cls
    set /p "USERNAME_=USERNAME: "
    cls
    set /p "USERDIRECTORY=USERDIRECTORY: "
    cls
    echo.
    echo Press Enter to set Default value: %DEFAULT_CERTIFICATESPATH%
    set /p "CERTIFICATESPATH=CERTIFICATESPATH: "
    if "!CERTIFICATESPATH!"=="" set "CERTIFICATESPATH=%DEFAULT_CERTIFICATESPATH%"
    cd /d "%~dp0"
    xcopy "!CERTIFICATESPATH!/*" "certificados" /s /i
    cls

    (
        echo # true or false
        echo SAAS = !SAAS!
        echo CLIENTE_NAME = !CLIENTE_NAME!
        echo.
        echo # Default C:/ProgramData/Qlik/DataTransfer/Log
        echo QLIK_LOGS_PATH_SAAS = !QLIK_LOGS_PATH_SAAS!
        echo.
        echo # Default C:/ProgramData/Qlik/Sense/Log 
        echo QLIK_LOGS_PATH_ONPRIMESE = !QLIK_LOGS_PATH_ONPRIMESE!
        echo.
        echo STAGE_PATH = ./stage
        echo LOGS_PATH = ./logs
        echo.
        echo # Default Port is 4747
        echo ENGINEPORT = !ENGINEPORT!
        echo.
        echo ENGINEHOST = !ENGINEHOST!
        echo USERNAME = !USERNAME_!
        echo USERDIRECTORY = !USERDIRECTORY!
    ) > .env
    xcopy "!CERTIFICATESPATH!" "./certificados"
    cls
    echo ==============================================
    echo             Successfully installed
    echo ==============================================
    echo Pressione qualquer tecla para sair.    
    pause > nul
) else if /i "%ACTION%"=="D" (
        cd /d "%~dp0"
        echo Desinstalando......
        main.exe -mode=uninstall
        echo Pressione qualquer tecla para sair.
        pause > nul
   ) else if /i "%ACTION%"=="U" (
    net stop "H2J Monitor"
    cls
    ping 127.0.0.1 -n 3 > nul
    cls
    cd /d "%~dp0"
    curl -LJO https://github.com/HenriqueJST/H2J_Monitor/raw/main/src/main.exe
    net start "H2J Monitor"
    cls
    echo ==============================================
    echo             Update Successfully
    echo ==============================================
    echo Pressione qualquer tecla para sair.    
    pause > nul
) else (
    echo Opção inválida. Por favor, escolha I, D ou U.
)

