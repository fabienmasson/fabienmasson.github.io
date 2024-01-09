## Browser support

You can open the XML file directly.
The `<?xml-stylesheet type="text/xsl" href="cv.xsl"?>` element allow the browser to display the resulting html content.

## Generate HTML

```
xsltproc cv.xsl fabien.xml > index.html
``` 

## Generate PDF

Install weasyprint (https://doc.courtbouillon.org/weasyprint/stable/first_steps.html#installation)
```
weasyprint index.html /mnt/c/temp/cv.pdf
```



