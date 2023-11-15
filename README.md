## How to install and run on desktops

```
mkdir bubbles
cd bubbles
go mod init examle.com/m
go run github.com/drevald/bubbles
```

## How to build for Android

```
git clone https://github.com/drevald/bubbles
cd bubbles
go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile bind -target android -javapkg com.flumine.bubbles -o ./mobile/android/bubbles/bubbles.aar ./mobile
```

and run the Android Studio project in `./mobile/android`.

## How to build for iOS

```
git clone https://github.com/drevald/bubbles
cd bubbles
go run github.com/hajimehoshi/ebiten/v2/cmd/ebitenmobile bind -target ios -o ./mobile/ios/Mobile.xcframework ./mobile
```

and run the Xcode project in `./mobile/ios`.
