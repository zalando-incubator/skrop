readonly VIPS_VERSION="8.4.5"

readonly IS_UBUNTU=$(cat /etc/*-release | grep -o -m 1 ubuntu)
readonly IS_ALPINE=$(cat /etc/*-release | grep -o -m 1 alpine)

function install_on_alpine {
  apk add --update \
    ca-certificates \
    wget \
    build-base \
    glib-dev \
    libxml2-dev \
    libjpeg-turbo-dev \
    libexif-dev \
    tiff-dev \
    libgsf-dev \
    libpng-dev \
  && wget http://www.vips.ecs.soton.ac.uk/supported/current/vips-${VIPS_VERSION}.tar.gz \
  && tar -zxvf vips-${VIPS_VERSION}.tar.gz \
  && cd vips-${VIPS_VERSION}/ \
  && ./configure \
    --prefix=/usr \
    --disable-debug \
    --disable-static \
    --disable-introspection \
    --disable-dependency-tracking \
    --without-python \
    --without-orc \
    --without-fftw \
  && make \
  && make install \
  && cd ../ \
  && rm -rf vips-${VIPS_VERSION}/ \
  && rm vips-${VIPS_VERSION}.tar.gz
}

function install_on_ubuntu {
  sudo apt-get -qq update \
  && sudo apt-get install -y --no-install-recommends \
    ca-certificates \
    wget \
    build-essential \
    glib2.0-dev \
    libxml2-dev \
    libjpeg-turbo8-dev \
    libexif-dev \
    libtiff5-dev \
    libgsf-1-dev \
    libpng-dev \
  && wget http://www.vips.ecs.soton.ac.uk/supported/current/vips-${VIPS_VERSION}.tar.gz \
  && tar -zxvf vips-${VIPS_VERSION}.tar.gz \
  && cd vips-${VIPS_VERSION}/ \
  && ./configure \
    --prefix=/usr \
    --disable-debug \
    --disable-static \
    --disable-introspection \
    --disable-dependency-tracking \
    --without-python \
    --without-orc \
    --without-fftw \
  make \
  sudo make install
}

if [[ ! -z "$IS_UBUNTU" ]]; then
  install_on_ubuntu
elif [[ ! -z "$IS_ALPINE" ]]; then
  install_on_alpine
else
  echo "Unsupported operating system!"
  exit 1
fi
