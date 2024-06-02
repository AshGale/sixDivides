@echo off
setlocal EnableDelayedExpansion

:: Define the lists
set "PLATFORMS=windows darwin linux js"
set "ARCHITECTURE=amd64 amd64 arm64 wasm"

echo Starting Building...

:: Convert lists to arrays
set i=0
for %%p in (%PLATFORMS%) do (
    set "PLATFORM[!i!]=%%p"
    set /a i+=1
)

set i=0
for %%a in (%ARCHITECTURE%) do (
    set "ARCH[!i!]=%%a"
    set /a i+=1
)

@REM :: Get the current date and time
@REM for /f "tokens=1-3 delims=/ " %%a in ("%date%") do (
@REM     set "month=%%a"
@REM     set "day=%%b"
@REM     set "year=%%c"
@REM )
@REM for /f "tokens=1-2 delims=: " %%a in ("%time%") do (
@REM     set "hour=%%a"
@REM     set "minute=%%b"
@REM )

@REM :: Format date and time for filename
@REM set "DATE_TIME=!year!-!month!-!day!_!hour!-!minute!"

:: Get the length of the arrays (assuming both arrays are of the same length)
set /a i-=1
set "length=!i!"

:: Loop over the list of platforms, and build
for /L %%i in (0,1,!length!) do (
    set GOOS=!PLATFORM[%%i]!
    set GOARCH=!ARCH[%%i]!

    :: Determine the binary extension based on the platform
    set "BIN_EXT="
    if "!GOOS!"=="windows" (
        set "BIN_EXT=.exe"
    ) else if "!GOOS!"=="darwin" (
        set "BIN_EXT="
    ) else if "!GOOS!"=="linux" (
        set "BIN_EXT="
    ) else if "!GOOS!"=="js" (
        set "BIN_EXT=.wasm"
    )

    echo   Platform: !GOOS!, Architecture: !GOARCH!, Extension: !BIN_EXT!

   go build -o "builds\!GOOS!-!GOARCH!!BIN_EXT!" main.go
)

endlocal

echo Finished Building!

pause
