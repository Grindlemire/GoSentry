set -x


mkdir -p target/usr/local/gosentry
mkdir -p target/usr/lib/systemd/system
mkdir -p target/var/log/gosentry
cd ..
cp GoSentry c.yml seelog.xml packaging/target/usr/local/gosentry/
cp gosentry.service packaging/target/usr/lib/systemd/system/
touch packaging/target/var/log/gosentry/gosentry.log
cd packaging


mkdir scripts
echo 'useradd -M -r gosentry' > scripts/beforeInstall.sh
echo 'userdel gosentry' > scripts/afterRemove.sh



fpm -f -s dir -t rpm -n gosentry -v 0.0.1 \
      --config-files usr/local/gosentry/c.yml \
      --config-files usr/local/gosentry/seelog.xml \
      --rpm-os linux \
      --rpm-user gosentry \
      --before-install scripts/beforeInstall.sh \
      --after-remove scripts/afterRemove.sh \
      -C target \
      -m grindlemire@github.com \
      usr/local/gosentry/GoSentry \
      usr/local/gosentry/c.yml \
      usr/local/gosentry/seelog.xml \
      usr/lib/systemd/system/gosentry.service \
      var/log/gosentry/gosentry.log



fpm -f -s dir -t deb -n gosentry -v 0.0.1 \
      --config-files usr/local/gosentry/c.yml \
      --config-files usr/local/gosentry/seelog.xml \
      --rpm-os linux \
      --rpm-user gosentry \
      --before-install scripts/beforeInstall.sh \
      --after-remove scripts/afterRemove.sh \
      -C target \
      -m grindlemire@github.com \
      usr/local/gosentry/GoSentry \
      usr/local/gosentry/c.yml \
      usr/local/gosentry/seelog.xml \
      usr/lib/systemd/system/gosentry.service \
      var/log/gosentry/gosentry.log
