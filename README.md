# svg2gcode

Converts svg to gcode for pen plotters


## Requirements (optional)

If you want to convert a image to an svg, you will need imagemagick and autotrace.

```bash
sudo apt update
sudo apt install intltool imagemagick libmagickcore-dev pstoedit libpstoedit-dev autopoint

git clone https://github.com/autotrace/autotrace.git
cd autotrace
./autogen.sh
LD_LIBRARY_PATH=/usr/local/lib ./configure --prefix=/usr
make
sudo make install
```

(Windows users an just download `autotrace` from [here](https://github.com/scottvr/autotrace-win64-binaries/tree/master/bin)).

## Usage

```
git clone https://github.com/schollz/svg2gcode
cd svg2gcode
make camel.gcode
```

## Other similar repos

- https://github.com/mnk343/Single-Line-Portrait-Drawing
- https://github.com/davekch/linerizer
- https://github.com/javierbyte/pintr