@echo off
setlocal

SET EXE=C:\monsoon\arc_windows.exe
SET LOG_FILE=C:\monsoon\agent.log

nssm stop arc
nssm remove arc confirm
nssm install arc %EXE% server

nssm set arc Description "Monsoon cloud remote control agent"
nssm set arc DisplayName "Monsoon agent"
nssm set arc AppStdout %LOG_FILE%
nssm set arc AppStderr %LOG_FILE%
nssm set arc AppRotateFiles 1
nssm set arc AppRotateBytes 100000
nssm set arc AppRotateOnline 1
nssm set arc AppStopMethodSkip 6
nssm set arc AppStopMethodConsole 2000
