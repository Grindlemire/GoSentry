set -x


mkdir -p target/usr/local/gosentry
mkdir -p target/usr/lib/systemd/system
cd ..
cp GoSentry c.yml seelog.xml packaging/target/usr/local/gosentry/
cp gosentry.service packaging/target/usr/lib/systemd/system/
cd packaging



fpm -f -s dir -t rpm -n gosentry -v 0.0.1 \
      --config-files usr/local/gosentry/c.yml \
      --config-files usr/local/gosentry/seelog.xml \
      --rpm-os linux \
      --rpm-user gosentry \
      -C target \
      -m grindlemire@github.com \
      usr/local/gosentry/GoSentry \
      usr/local/gosentry/c.yml \
      usr/local/gosentry/seelog.xml \
      usr/lib/systemd/system/gosentry.service



fpm -f -s dir -t deb -n gosentry -v 0.0.1 \
      --config-files usr/local/gosentry/c.yml \
      --config-files usr/local/gosentry/seelog.xml \
      --rpm-os linux \
      --rpm-user gosentry \
      -C target \
      -m grindlemire@github.com \
      usr/local/gosentry/GoSentry \
      usr/local/gosentry/c.yml \
      usr/local/gosentry/seelog.xml \
      usr/lib/systemd/system/gosentry.service