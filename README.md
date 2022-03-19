Golang tetris implementation powered by ebiten

Structure
./bubblemobile - package for android initializer
./game - package with game itself

To generate android package run:
cd bubblemobile
ebitenmobile bind -target android -javapkg bubbles.game -o bubbles.aar .