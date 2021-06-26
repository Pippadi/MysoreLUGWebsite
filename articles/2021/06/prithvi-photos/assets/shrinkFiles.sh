for file in originals/* ; do
	imagename=$(echo $file | sed -e s,.jpg,, | sed -e s,originals/,,)
	magick convert $file -resize 800x800 shrunk/$imagename.webp
	echo Created shrunk/$imagename.webp
done
