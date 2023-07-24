@echo off

echo Parando os serviços do Qlik Sense...

net stop "Qlik Sense Repository Database"
net stop "Qlik Sense Logging Service"
net stop "Qlik Sense Service Dispatcher"
net stop "Qlik Sense Proxy Service"
net stop "Qlik Sense Engine Service"
net stop "Qlik Sense Scheduler Service"
net stop "Qlik Sense Printing Service"
net stop "Qlik Sense Repository Service"

echo Aguardando 10 segundos...
ping 127.0.0.1 -n 11 > nul

echo Iniciando os serviços do Qlik Sense...

net start "Qlik Sense Repository Database"
net start "Qlik Sense Logging Service"
net start "Qlik Sense Service Dispatcher"
net start "Qlik Sense Repository Service"
net start "Qlik Sense Proxy Service"
net start "Qlik Sense Engine Service"
net start "Qlik Sense Scheduler Service"
net start "Qlik Sense Printing Service"

echo Concluído!