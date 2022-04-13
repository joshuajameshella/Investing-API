rm -rf dist
mkdir dist
env GOOS=linux go build -ldflags="-s -w" -o main .
zip BuyPosition.zip main
mv BuyPosition.zip ./dist/
rm main