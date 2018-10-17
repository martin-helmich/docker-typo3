TYPOVER=6.2 7.6 8.7 9.1 9.2 9.3 9.4 9.5
DOCKERFILES=$(foreach subdir, $(TYPOVER), $(subdir)/Dockerfile)
DEPENDS=Dockerfile.in Makefile

all: $(DOCKERFILES)

6.2/Dockerfile: Dockerfile.in
	sed -e 's/PHPVER/5.6/' -e 's/TYPOVER/6.2/' $< > $@

7.6/Dockerfile: Dockerfile.in
	sed -e 's/PHPVER/5.6/' -e 's/TYPOVER/7.6/' $< > $@

8.7/Dockerfile: Dockerfile.in
	sed -e 's/PHPVER/7.2/' -e 's/TYPOVER/8.7/' $< > $@

9.1/Dockerfile: Dockerfile.in
	sed -e 's/PHPVER/7.2/' -e 's/TYPOVER/9.1/' $< > $@

9.2/Dockerfile: Dockerfile.in
	sed -e 's/PHPVER/7.2/' -e 's/TYPOVER/9.2/' $< > $@

9.3/Dockerfile: Dockerfile.in
	sed -e 's/PHPVER/7.2/' -e 's/TYPOVER/9.3/' $< > $@

9.4/Dockerfile: Dockerfile.in
	sed -e 's/PHPVER/7.2/' -e 's/TYPOVER/9.4/' $< > $@

9.5/Dockerfile: Dockerfile.in
	sed -e 's/PHPVER/7.2/' -e 's/TYPOVER/9.5/' $< > $@
