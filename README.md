Golang tetris implementation powered by ebiten

Structure
./bubblemobile - package for android initializer
./game - package with game itself

To generate android package run:
cd bubblemobile
bubblemobile denis$ ebitenmobile bind -target android -javapkg com.flumine.bubbles -o bubbles.aar .