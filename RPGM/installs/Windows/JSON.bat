@echo off

rem checks if node.js exists or not
rem https://nodejs.org/en/download/

if not exist "%ProgramFiles%\nodejs\node.exe" (
   echo "Node.js not found. Please install it first."
   echo "Go to the node.js website: https://nodejs.org/en"
   exit /b 1
)

rem install json-translator
rem https://www.npmjs.com/package/@parvineyvazov/json-translator

npm install --global @parvineyvazov/json-translator

echo "json-translator installed successfully"

rem usage example jsontt 

jsontt  -h

pause
