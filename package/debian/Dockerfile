FROM amezin/debian-build-essentials:jessie
MAINTAINER allonsy

#supply the version number here. Do not prepend a 'v'
ENV version 4.0.1

RUN apt-get update && apt-get install -y git wget

CMD cd /root && \
git clone --branch v${version} https://github.com/Originate/git-town.git && \
cd /root/git-town/package/debian && \
wget https://github.com/Originate/git-town/releases/download/v4.0.1/git-town-linux-amd64 && \
mv git-town-linux-amd64 git-town && \
./debian_build.sh && \
cp ../gittown_${version}_amd64.deb /shared
