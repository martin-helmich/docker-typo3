TYPOVER=6.2 7.6 8.0
DOCKERFILES=$(foreach subdir, $(TYPOVER), $(subdir)/Dockerfile)
DEPENDS=Dockerfile.in Makefile

all: $(DOCKERFILES)

6.2/Dockerfile: $(DEPENDS)
	sed -e 's/PHPVER/5.6/' -e 's/TYPOVER/6.2/' $< > $@

7.6/Dockerfile: $(DEPENDS)
	sed -e 's/PHPVER/5.6/' -e 's/TYPOVER/7.6/' $< > $@

8.0/Dockerfile: $(DEPENDS)
	sed -e 's/PHPVER/7/' -e 's/TYPOVER/8/' $< > $@

