<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:template match="/cv">
    <html lang="fr">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>Curriculum <xsl:value-of select="name"/></title>
            <link rel="stylesheet" href="cv.css"/>
            <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css"  />        
        </head>
    <body>
        <main>
            <header>
                <div class="intro">
                    <h1><xsl:value-of select="name"/></h1>
                    <p><xsl:value-of select="summary"/></p>
                    <br/>
                    <p>
                        <a target="_blank">
                            <xsl:attribute name="href">
                                <xsl:value-of select="locationLink"/>
                            </xsl:attribute>
                            <i class="fa-solid fa-map-location"></i>  
                            <xsl:text> </xsl:text>
                            <xsl:value-of select="location"/>
                        </a>
                    </p>
                    <p>
                        <a>
                            <xsl:attribute name="href">
                                <xsl:text>mailto:</xsl:text>
                                <xsl:value-of select="contact/email"/>
                            </xsl:attribute>
                            <i class="fa-regular fa-envelope"></i>  
                            <xsl:text> </xsl:text>
                            <xsl:value-of select="contact/email"/>
                        </a>
                    </p>
                    <p>
                        <a>
                            <xsl:attribute name="href">
                                <xsl:text>tel:</xsl:text>
                                <xsl:value-of select="contact/telephone"/>
                            </xsl:attribute>
                            <i class="fa-solid fa-phone"></i>  
                            <xsl:text> </xsl:text>
                            <xsl:value-of select="contact/telephone"/>
                        </a>
                    </p>
                </div>
                <div class="portrait">
                    <img>
                        <xsl:attribute name="src">
                            <xsl:value-of select="avatarUrl"/>
                        </xsl:attribute>
                    </img>
                    <ul class="social">
                        <xsl:apply-templates select="contact/social"/>
                    </ul>
                </div>
            </header>

            <section id="presentation">
                <h2>Présentation</h2>
                <p><xsl:value-of select="about"/></p>
            </section>

            <div id="work">
                <h2>Expériences professionnelles</h2>
                <xsl:apply-templates select="work"/>
            </div>

            <section id="education">
                <h2>Formations</h2>
                <xsl:apply-templates select="education"/>
            </section>

            <section id="skills">
                <h2>Compétences</h2>
                <xsl:apply-templates select="skills"/>
            </section>

            <section id="projects">
                <h2>Projets personnels</h2>
                <div class="projects">
                    <xsl:apply-templates select="project"/>
                </div>
            </section>

        </main>
    </body>
    </html>
</xsl:template>

<xsl:template match="social">
  <li>
    <a target="_blank">
        <xsl:attribute name="href">
            <xsl:value-of select="@url"/>
        </xsl:attribute>
        <i>
            <xsl:attribute name="class">
                <xsl:text>fa-brands fa-</xsl:text>
                <xsl:value-of select="@name"/>
            </xsl:attribute>
        </i>
    </a>
  </li>
</xsl:template>

<xsl:template match="work">
  <section>
  <div class="work-head">
    <h3>
        <img>
            <xsl:attribute name="src">
                <xsl:text>img/</xsl:text>    
                <xsl:value-of select="logo"/>
            </xsl:attribute>
        </img>
        <xsl:text> </xsl:text>
        <xsl:value-of select="company"/>
    </h3>
    <p>
        <xsl:choose>
        <xsl:when test="end">
            <xsl:value-of select="start"/>
            <xsl:text> - </xsl:text>
            <xsl:value-of select="end"/>
        </xsl:when>
        <xsl:otherwise>
            <xsl:text>Depuis </xsl:text>
            <xsl:value-of select="start"/>
        </xsl:otherwise>
        </xsl:choose>
    </p>
  </div>
  <h6><xsl:value-of select="title"/></h6>
  <xsl:for-each select="mission">
    <p class="mission-title smaller">
        <xsl:value-of select="title"/>
    </p> 
    <xsl:for-each select="detail">
        <p class="smaller"><xsl:value-of select="."/></p>
    </xsl:for-each>
  </xsl:for-each>
  </section>
</xsl:template>

<xsl:template match="education">
  <div class="education-head">
    <h3>
        <xsl:value-of select="school"/>
    </h3>
    <p>
        <xsl:if test="end">
            <xsl:value-of select="start"/>
            <xsl:text> - </xsl:text>
            <xsl:value-of select="end"/>
        </xsl:if>
    </p>
  </div>
  <xsl:if test="city">
    <p>
        <i class="fa-solid fa-map-location"></i>  
        <xsl:text> </xsl:text>
        <xsl:value-of select="city"/>
    </p>
  </xsl:if>
  <xsl:for-each select="degree">
    <p class="smaller">
        <xsl:value-of select="."/>
        <xsl:if test="@certification or @link">
        <xsl:choose>
            <xsl:when test="@link">
                <xsl:text> (</xsl:text>
                <xsl:if test="@certification">
                    <xsl:value-of select="@certification"/>
                    <xsl:text> </xsl:text>
                </xsl:if>
                <a target="_blank">
                    <xsl:attribute name="href">
                    <xsl:value-of select="@link"/>
                    </xsl:attribute>
                    <i class="fa-solid fa-certificate"></i>
                </a>
                <xsl:text>)</xsl:text>   
                
            </xsl:when>
            <xsl:otherwise>
                <xsl:text> (</xsl:text>
                <xsl:value-of select="@certification"/>
                <xsl:text>)</xsl:text>    
            </xsl:otherwise>
        </xsl:choose>
        </xsl:if>
    </p>
  </xsl:for-each>

</xsl:template>

<xsl:template match="skills">
  <h3><xsl:value-of select="@title"/></h3>
  <xsl:for-each select="skill">
    <div class="skill-row">
        <p class="smaller" style="flex:0.5;"><xsl:value-of select="."/></p>
        <div class="gauge">
            <div>
            <xsl:attribute name="style">
            <xsl:text>width:</xsl:text>
            <xsl:value-of select="@level"/>
            <xsl:text>%;</xsl:text>
            </xsl:attribute>
            </div>
        </div>
    </div>
  </xsl:for-each>
</xsl:template>

<xsl:template match="project">
  <div>
    <h3><xsl:value-of select="title"/></h3>
    <p class="smaller"><xsl:value-of select="description"/></p>
    <p>
    <xsl:for-each select="tech">
        <span class="tag"><xsl:value-of select="."/></span>
    </xsl:for-each>
    </p>
  </div>
</xsl:template>


</xsl:stylesheet> 