######################################
# Build the Goop (Go Object-Oriented #
# Programming) package		     #
#				     #
# By Scott Pakin <pakin@lanl.gov>    #
######################################

include $(GOROOT)/src/Make.inc

VERSION=1.0

TARG=goop

GOFILES=\
	goop.go\

DISTFILES=\
	$(GOFILES)

include $(GOROOT)/src/Make.pkg

# ---------------------------------------------------------------------------

FULLNAME=goop-$(VERSION)

dist: $(DISTFILES)
	mkdir $(FULLNAME)
	cp $(DISTFILES) $(FULLNAME)
	tar -czf $(FULLNAME).tar.gz $(FULLNAME)
	$(RM) -r $(FULLNAME)
	tar -tzvf $(FULLNAME).tar.gz
