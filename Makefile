
svg2gcode:
	go build -v

all: camel.gcode person.gcode horse.gcode flamingo.gcode dog.gcode 

camel.gcode: svg2gcode
	convert src/examples/camel.jpg -resize 300x -background White -gravity center -threshold 60% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out camel -x 0 -y 0 --width 90 --height 90 --simplify 0.005 --consolidate 0.07 --max-length 0.03 --animate
	rm 1.tga
	rm potrace.svg
	
person.gcode: svg2gcode
	convert src/examples/person1.jpg -resize 300x -background White -gravity center -threshold 60% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out person -x 0 -y 0 --width 90 --height 90 --simplify 0.005 --consolidate 0.07 --max-length 0.03 --animate
	rm 1.tga
	rm potrace.svg

horse.gcode: svg2gcode
	convert src/examples/horse1.jpg -resize 300x -background White -gravity center -threshold 60% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out horse -x 0 -y 0 --width 90 --height 90 --simplify 0.005 --consolidate 0.07 --max-length 0.03 --animate

flamingo.gcode: svg2gcode
	convert src/examples/flamingo.png 1.jpg
	convert 1.jpg -resize 300x -background White -gravity center -threshold 30% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out flamingo -x 0 -y 0 --width 90 --height 90 --simplify 0.005 --consolidate 0.15 --max-length 0.1 --animate

dog.gcode: svg2gcode
	convert src/examples/dog.png 1.jpg
	convert 1.jpg -resize 300x -background White -gravity center -threshold 60% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out dog -x 0 -y 0 --width 90 --height 90 --simplify 0.005 --consolidate 0.07 --max-length 0.03 --animate

penguin.gcode: svg2gcode
	convert src/examples/penguin.png 1.jpg
	convert 1.jpg -resize 250x -background White -gravity center -threshold 50% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out penguin -x 0 -y 0 --width 90 --height 90 --simplify 0.005 --consolidate 0.1 --max-length 0.2 --animate

people.gcode: svg2gcode
	convert src/examples/people.png 1.jpg
	convert 1.jpg -resize 500x -background White -gravity center -threshold 80% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out people -x 0 -y 0 --width 90 --height 90 --simplify 0.005 --consolidate 0.16 --max-length 0.05 --animate
