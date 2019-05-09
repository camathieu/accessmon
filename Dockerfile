#####
####
### Accessmon - Docker file
##
#

# Let's start with a fresh alpine build
FROM alpine

# Some generic information
MAINTAINER Charles-Antoine Mathieu

# Build and copy
ADD cmd/accessmon /usr/bin/
RUN touch /tmp/access.log

# Launch it
CMD /usr/bin/accessmon