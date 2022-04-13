rm -rf dist
mkdir dist
env GOOS=linux go build -ldflags="-s -w" -o main .
zip SellPosition.zip main
mv SellPosition.zip ./dist/
rm main