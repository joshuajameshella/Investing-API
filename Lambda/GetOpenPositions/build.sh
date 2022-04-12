rm -rf dist
mkdir dist
env GOOS=linux go build -ldflags="-s -w" -o main .
zip GetOpenPositions.zip main
mv GetOpenPositions.zip ./dist/
rm main