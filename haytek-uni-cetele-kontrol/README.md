# Cetele

haytek sıtacer üni çetele kontrolü yapar.

### sunucuya dosya gönderme

scp -i ~/.ssh/cetele.pem -r /Users/talha/go/src/github.com/haytek-uni-cetele-kontrol/main ubuntu@3.89.107.238:/home/ubuntu

### linux exe çıktısı

GOOS=linux GOARCH=amd64 go build main.go

### system start

sudo systemctl start cetele.service

### set env var to systemd
sudo systemctl set-environment IS_DEVELOPMENT=true  || makine reboot olursa tekrar girmek gerekiyor.

### systom stop

sudo systemctl stop cetele.service

## baglan

ssh -i ~/.ssh/cetele.pem ubuntu@ec2-3-89-107-238.compute-1.amazonaws.com

## postgres

sudo su postgres --- psql

## last logs

systemctl -l status cetele

# bug fix todo

- refresh token ile token aldığında write protectable yapıyor. onu düzeltmem lazım -----  
- token.json dosyasını ubuntunun chown yapması lazım -------
- scheduled messagelerde sleep olsa da ticker orada tickelemeye devam ediyor. sleep bittiğinde arkada 23 * 6 tane tick bekliyor. o yüzden erken atıyor mesajı. ticker stop ile durdurup sleep bitince yine başlatmalı ya da doğru zamanı buldurup 24 saat uyutmalı (şimdilik resetleyerek yaptım çünkü bu channel karın ağrısnı daha kullanamıyorum) -------
- mesajlarda saati utc farklı gösteriyor. ona da europe istanbul vermek lazım. bunu yazacağıma yapabilrdim ama boşver yarın yaparım. ------ 


# TODO
- alınan tarihi trimspace falan yapmak lazım.(daha yapılmadı gerekmedi de aslında. bakarız) ----- 
- şu production olayını gerçekten çözmem lazım paso config değiştirmek midemi bulandırdı. -----
- panic catch yapılmalı. uygulama patladığında bana log göndermeli. 
- günlük haftalık aylık raporlar oluşturulmalı.
- gruba bunu atsın ama özelden sadece kişinin kendi özetini atsın -------- 
- biraz generic düşüne
- #N/A gelme durumu da var.



# GENEL GÜNCELLEME
- program adeti negatif ise program yapılmadı  olarak kabul ediliyor ve error mesajları da formatlanıp gönderiliyor.
- program adeti rakam olmayanlar direkt kabul edilmiyor ve yapılmadı sayılıyor.
- sheet'se özet sayfası eklendi dünün özetini gösteriyor fakat bu dün kavramı neye kime göre şuan belli değil. aws için bugün mü dün mü falan filan onlar test edilecek.
- refresh token'i refresh yapılabiliyor mu veya bu cetele app'i publish etmek paralı mı değil mi falan o da belli değil yoksa 1 haftada expire oluyor bo. çukuru


=METNEÇEVİR(SAYIYAÇEVİR(B4); "yyyy-mm-dd")
=QUERY(B2:F34;"select B,C,D,E,F where B = date '"& METNEÇEVİR(SAYIYAÇEVİR(BUGÜN()); "yyyy-mm-dd") &"' ";1)