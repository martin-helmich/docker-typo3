TYPOVER=6.2 7.6 8.7 9.1
DOCKERFILES=$(foreach subdir, $(TYPOVER), $(subdir)/Dockerfile)
DEPENDS=Dockerfile6-7.in Makefile

all: $(DOCKERFILES)

6.2/Dockerfile: Dockerfile6-7.in
	sed -e 's/PHPVER/5.6/' -e 's/TYPOVER/6.2/' $< > $@

7.6/Dockerfile: Dockerfile6-7.in
	sed -e 's/PHPVER/5.6/' -e 's/TYPOVER/7.6/' $< > $@

8.7/Dockerfile: Dockerfile8-9.in
	sed -e 's/PHPVER/7.2/' -e 's/TYPOVER/8.7/' $< > $@

9.1/Dockerfile: Dockerfile8-9.in
	sed -e 's/PHPVER/7.2/' -e 's/TYPOVER/9.1/' $< > $@

