FROM amezin/debian-build-essentials:jessie
MAINTAINER allonsy

#supply the version number here. Do not prepend a 'v'
ENV version 3.1.0

RUN apt-get update && apt-get install -y git

CMD cd /root && \
git clone --branch v${version} https://github.com/Originate/git-town.git && \
cd /root/git-town/package/debian && \
./debian_build.sh && \
cp ../gittown_${version}_all.deb /shared
