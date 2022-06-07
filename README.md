# MysoreLUGWebsite
A Website for the Mysore Linux Users' Group, currently at [plootarg.com](https://plootarg.com).

This site aims to be a medium for Linux users and geeks in general to express their thoughts and spread information.

The webpages themselves are designed to be very, very simple, in contrast to the bloated web of today. As of now, they use no JavaScript, and are completely static. There are no ads, sidebars, or other distractions.

## Contribution
To get started with an article, fork the repository and create a new folder for your article in `articles`, under the appropriate year and month folders.
There is a folder called `template-article` which will tell you where you need to keep what, and how to format your content.

Once your article reads and looks the way you want it, submit a pull request with your changes. They will be reviewed, and if satisfactory, merged.
The same goes for edits of others' articles, and changes to the format itself.

### Actually writing an article
These webpages are stored as static HTML. You could write your webpage in raw HTML, but that can get cumbersome very quickly.
For most articles that have no complex formatting requirements, writing in our markdown-like format can be the most convenient.

Formatting example:
```
# Template article
## This is a subtitle
---

### This is a heading

This is a paragraph.

#### This is a subheading

_This_ *is* `another` _paragraph_.
Here are some random characters- \` \_

\```
This is a code block.
\```

!../articles/template-article/assets/rpi.png
!width="70%" height="auto" alt="Raspberry Pi"
A Raspberry Pi

This is a [link](https://plootarg.com).
```

To run the `mdparse` script from the root folder of the repository to generate HTML:
```
mdparse/mdparse articles/2021/12/example-article/example-article.md mdparse/templates > articles/2021/12/example-article/index.html
```
`mdparse/templates` is the location of the HTML template files.
The script writes the HTML to `STDOUT`, so there is no need to create a new file to preview. Once final, the output can simply be redirected to a file as shown.
The example's HTML output can be seen [here](https://plootarg.com/articles/2021/12/template-article/Index.html).

The default templates in the `mdparse/templates` contain relative paths for items like the main stylesheet (`shared/css/main.css`). This means that you may want to change these paths depending on the location of your article.
By default, the relative paths are based on the assumption that the final HTML will reside in the `articles/<year>/<month>/<article-name>` directory (which is where articles ought to be).
