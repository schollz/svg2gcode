
svg2gcode:
	go build -v

camel.gcode: svg2gcode
	convert src/examples/camel.jpg -resize 300x -background White -gravity center -threshold 60% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out camel -x 0 -y 0 --width 45 --height 45 --simplify 0.005 --consolidate 0.07 --max-length 0.03 --animate
	rm 1.tga
	rm potrace.svg
	
person.gcode: svg2gcode
	convert src/examples/person1.jpg -resize 300x -background White -gravity center -threshold 60% 1.tga
	autotrace -output-file potrace.svg --output-format svg --centerline 1.tga
	./svg2gcode convert --debug --png --in potrace.svg --out person -x 0 -y 0 --width 45 --height 45 --simplify 0.005 --consolidate 0.07 --max-length 0.03 --animate
	rm 1.tga
	rm potrace.svg
