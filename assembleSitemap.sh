#!/bin/bash

if test -z "$1" ; then
	echo "Usage: assembleSitemap.sh DESTINATION_FILE"
	exit 1
fi

sitemapfile=$1
today=$(date '+%F')
domain="plootarg.com"

echo '<?xml version="1.0" encoding="UTF-8"?>' > $sitemapfile
echo '<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">' >> $sitemapfile
echo "	<url>
		<loc>https://$domain/</loc>
		<lastmod>$today</lastmod>
	</url>" >> $sitemapfile

shopt -s globstar

for item in ** ; do
	if [ -d $item ] && [ -f $item/index.html ] ; then
		echo \
"	<url>
		<loc>https://$domain/$item/</loc>
		<lastmod>$today</lastmod>
	</url>" >> $sitemapfile
	fi
done

echo '</urlset>' >> $sitemapfile
